"""入库主流程（§7.4；可 CLI 调用）。

流程::

    遍历 data/raw/{ns}/**  →  loader.load_file
                          →  splitter.split_document
                          →  写 data/processed/{ns}.jsonl（可审计）
                          →  embedding.encode
                          →  Chroma upsert（metadata 带 §7.3 四字段）
                          →  清理本 namespace 的孤儿块（语料删改后不留残骸）
                          →  重载内存 BM25（P0-3 落地后生效）

三处必须一致，否则会**静默串味**：

1. Chroma 集合名（``Settings.collection_name``）
2. BM25 索引的 namespace 划分（**每个 namespace 一份，不可共用**，ADR-0005）
3. ``data/raw/`` 的子目录名

``chunk_id`` 是 Chroma 的 id，所以重复入库天然幂等：``upsert`` 覆盖同 id，
不产生重复块（§9 :777）。代价是语料**新增**文件会让排序靠后的文档整体
重编号，于是旧块变成孤儿——``prune`` 负责把它们删掉，否则库里会同时留着
新旧两套编号，检索结果出现近似重复。

CLI::

    python -m app.ingest.pipeline split  --input data/raw --output data/processed
    python -m app.ingest.pipeline ingest --input data/processed
"""

from __future__ import annotations

import argparse
import json
import sys
from pathlib import Path

from app.core.config import get_settings
from app.core.logging import get_logger
from app.ingest.loader import iter_raw_docs
from app.ingest.splitter import Chunk, split_document

logger = get_logger(__name__)

# §7.3 要求的四个元数据字段，缺一个引用回显或过滤检索就会退化
METADATA_FIELDS: tuple[str, ...] = ("source", "kb_namespace", "section_title", "chunk_id")

DEFAULT_EMBED_BATCH = 64


def _write_jsonl(path: Path, rows: list[dict]) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    with path.open("w", encoding="utf-8") as fh:
        for row in rows:
            fh.write(json.dumps(row, ensure_ascii=False) + "\n")


def _read_jsonl(path: Path) -> list[dict]:
    rows: list[dict] = []
    with path.open("r", encoding="utf-8") as fh:
        for lineno, line in enumerate(fh, 1):
            line = line.strip()
            if not line:
                continue
            try:
                rows.append(json.loads(line))
            except json.JSONDecodeError as exc:
                raise ValueError(f"{path}:{lineno} 不是合法 JSON：{exc}") from exc
    return rows


def split_corpus(
    *,
    raw_root: Path | str | None = None,
    out_dir: Path | str | None = None,
    namespaces: list[str] | None = None,
) -> dict[str, int]:
    """按 §7.3 切块，产出 ``data/processed/{ns}.jsonl``。

    Returns:
        ``{namespace: 块数}``。请求了但目录里没有语料的 namespace 记 0，
        这样调用方不用猜「是没跑还是跑空了」。
    """
    s = get_settings()
    raw_root = Path(raw_root) if raw_root else s.data_raw_dir
    out_dir = Path(out_dir) if out_dir else s.data_processed_dir

    ns_list = list(namespaces) if namespaces else list(s.kb_namespaces)
    stats: dict[str, int] = {ns: 0 for ns in ns_list}
    keep = set(ns_list)

    # iter_raw_docs 已是确定性顺序；这里只按 namespace 分桶，不重排
    buckets: dict[str, list[Chunk]] = {ns: [] for ns in ns_list}
    for doc in iter_raw_docs(raw_root):
        if doc.kb_namespace not in keep:
            continue
        buckets[doc.kb_namespace].append(doc)  # type: ignore[arg-type]

    for ns, docs in buckets.items():
        seq = 0
        chunks: list[Chunk] = []
        for doc in docs:  # type: ignore[assignment]
            produced = split_document(doc, start_seq=seq)
            seq += len(produced)
            chunks.extend(produced)
        _write_jsonl(out_dir / f"{ns}.jsonl", [c.to_dict() for c in chunks])
        stats[ns] = len(chunks)
        logger.info("切块 %-6s %2d 份文档 → %4d 块 → %s", ns, len(docs), len(chunks), out_dir / f"{ns}.jsonl")

    return stats


def ingest_corpus(
    *,
    processed_dir: Path | str | None = None,
    namespaces: list[str] | None = None,
    prune: bool = True,
    batch_size: int = DEFAULT_EMBED_BATCH,
) -> dict[str, int]:
    """把 ``data/processed/*.jsonl`` 向量化并 upsert 进 Chroma。

    Args:
        prune: 删除该 namespace 下**不在本批**的旧 ``chunk_id``
            （语料删改后的残骸清理；留着会让检索结果出现近似重复）。

    Raises:
        ValueError: 同一批里 ``chunk_id`` 重复——那会让 upsert 静默互相覆盖。
    """
    from llama_index.core.schema import TextNode

    from app.rag.embedding import get_embedding_model
    from app.rag.index import get_collection, get_vector_store

    s = get_settings()
    processed_dir = Path(processed_dir) if processed_dir else s.data_processed_dir
    ns_list = list(namespaces) if namespaces else list(s.kb_namespaces)

    collection = get_collection(create=True)
    vector_store = get_vector_store()
    if vector_store is None:
        raise RuntimeError(
            "向量库不可用：Chroma 集合打不开，写库已中止（不要假装入库成功）"
        )
    model = get_embedding_model()
    stats: dict[str, int] = {}

    for ns in ns_list:
        path = processed_dir / f"{ns}.jsonl"
        if not path.is_file():
            logger.warning("跳过 %s：找不到 %s（先跑 split）", ns, path)
            stats[ns] = 0
            continue

        rows = _read_jsonl(path)
        if not rows:
            logger.warning("跳过 %s：%s 是空的", ns, path)
            stats[ns] = 0
            continue

        ids = [str(r["chunk_id"]) for r in rows]
        dupes = len(ids) - len(set(ids))
        if dupes:
            raise ValueError(
                f"{ns} 有 {dupes} 个重复 chunk_id，拒绝入库："
                "upsert 会静默互相覆盖，最终库里的块数对不上 jsonl"
            )

        missing = [i for i, r in enumerate(rows) if any(k not in r for k in METADATA_FIELDS)]
        if missing:
            raise ValueError(
                f"{ns} 的第 {missing[:5]} 行缺少 §7.3 要求的元数据字段"
                f"{METADATA_FIELDS}（§9 :775 验收要求每个 chunk 四字段齐全）"
            )

        texts = [str(r["text"]) for r in rows]
        metadatas = [{k: r[k] for k in METADATA_FIELDS} for r in rows]

        # 向量在这里一次算完再交给节点，而不是让 VectorStoreIndex.insert_nodes 逐个
        # 去算——后者是 1319 次单条 encode，比批量慢一个数量级。
        embeddings = model.encode(texts)

        nodes = [
            TextNode(
                id_=ids[i],
                text=texts[i],
                metadata=metadatas[i],
                embedding=list(embeddings[i]),
            )
            for i in range(len(rows))
        ]

        # ⚠️ ChromaVectorStore.add 内部走的是 Chroma 的 **add**，不是 upsert：
        # 同一个 id 写第二次会直接抛错（"IDs must be unique"）。没有 upsert 可用，
        # 所以「幂等」必须靠先删后写来实现——只删本批要写的那些 id，不碰其他块。
        existing_here = set(collection.get(where={"kb_namespace": ns})["ids"]) & set(ids)
        if existing_here:
            vector_store.delete_nodes(node_ids=sorted(existing_here))
            logger.info("%s 覆盖 %d 个已存在的块（先删后写）", ns, len(existing_here))

        for start in range(0, len(nodes), batch_size):
            stop = min(start + batch_size, len(nodes))
            vector_store.add(nodes[start:stop])

        if prune:
            existing = set(collection.get(where={"kb_namespace": ns})["ids"])
            stale = sorted(existing - set(ids))
            if stale:
                vector_store.delete_nodes(node_ids=stale)
                logger.info("%s 清理 %d 个孤儿块（旧编号）", ns, len(stale))

        stats[ns] = len(rows)
        logger.info("入库 %-6s %4d 块（集合共 %d）", ns, len(rows), collection.count())

    return stats


def reload_indexes() -> dict[str, int]:
    """只重载内存索引（``POST /api/v1/admin/reload`` 的实现）。

    BM25 的内存索引属 P0-3；它落地前这里如实返回 Chroma 的块数，
    而不是假装重载成功。Chroma 是持久化真源，进程重启后它是唯一依据。
    """
    from app.rag import retriever
    from app.rag.index import collection_count

    reloader = getattr(retriever, "reload_indexes", None)
    if callable(reloader):
        out = reloader()
        # 新版 reload_indexes 只清缓存、不回报计数，故兜底回落到 Chroma 真源
        return dict(out) if isinstance(out, dict) else {"chroma_total": collection_count()}

    logger.info("BM25 内存索引重载属 P0-3，暂以 Chroma 计数代替")
    return {"chroma_total": collection_count()}


def ingest_all(
    *,
    namespaces: list[str] | None = None,
    rebuild: bool = True,
) -> dict[str, int]:
    """全量入库。

    Args:
        rebuild: ``False`` 时**只切块不写库**（``--dry-run``），返回块数，
            用于在不触碰 Chroma 的前提下核对切块质量。

    Returns:
        ``{namespace: 块数}``。
    """
    stats = split_corpus(namespaces=namespaces)
    if not rebuild:
        logger.info("dry-run：只切块，未写 Chroma")
        return stats

    ingested = ingest_corpus(namespaces=namespaces, prune=True)
    reload_indexes()
    return ingested


def _resolve(path_str: str | None) -> Path | None:
    return Path(path_str).expanduser().resolve() if path_str else None


def main(argv: list[str] | None = None) -> int:
    """§7.4 的命令行入口。"""
    parser = argparse.ArgumentParser(prog="python -m app.ingest.pipeline", description="知识库切块 / 入库（P0-2）")
    sub = parser.add_subparsers(dest="command", required=True)

    p_split = sub.add_parser("split", help="只切块，产出 data/processed/*.jsonl")
    p_split.add_argument("--input", dest="raw_root", help="原始语料根目录，默认 <rag-service>/data/raw")
    p_split.add_argument("--output", dest="out_dir", help="切块输出目录，默认 <rag-service>/data/processed")
    p_split.add_argument("--ns", action="append", help="只处理指定 namespace，可重复")

    p_ingest = sub.add_parser("ingest", help="向量化并写入 Chroma")
    p_ingest.add_argument("--input", dest="processed_dir", help="切块目录，默认 <rag-service>/data/processed")
    p_ingest.add_argument("--collection", dest="collection", help="覆盖集合名（默认取 .env 的 COLLECTION_NAME）")
    p_ingest.add_argument("--ns", action="append", help="只处理指定 namespace，可重复")
    p_ingest.add_argument("--no-prune", action="store_true", help="保留不再产出的旧 chunk_id")

    args = parser.parse_args(argv)

    if args.collection:
        # 只影响本次进程的集合名；Settings 是 lru_cache 单例，直接改实例
        get_settings().collection_name = args.collection

    if args.command == "split":
        stats = split_corpus(
            raw_root=_resolve(args.raw_root),
            out_dir=_resolve(args.out_dir),
            namespaces=args.ns,
        )
    else:
        stats = ingest_corpus(
            processed_dir=_resolve(args.processed_dir),
            namespaces=args.ns,
            prune=not args.no_prune,
        )
        reload_indexes()

    total = sum(stats.values())
    for ns, n in sorted(stats.items()):
        print(f"  {ns:8} {n:5d}")
    print(f"  {'合计':8} {total:5d}")
    return 0 if total else 1


if __name__ == "__main__":
    sys.exit(main())
