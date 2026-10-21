"""BM25 通道：两种实现，用 30 条 fixture 实测定夺（ADR-0009）。

## 为什么需要两种

LlamaIndex 自带的 ``BM25Retriever``（``llama-index-retrievers-bm25`` 0.8.0，底层是
``bm25s`` 而非 ``rank_bm25``）的构造参数里**没有 ``tokenizer``**，只有
``token_pattern = r"(?u)\\b\\w\\w+\\b"``。``\\w`` 在 Unicode 模式下匹配 CJK，而 ``\\b``
只在**串首尾**成立，所以一整句中文会被切成**一个 token**——BM25 于是退化成「整句
精确匹配」，中文检索直接失效。

``from_defaults`` 里确实有个 ``tokenizer`` 参数，但它已被标为 deprecated，且只是被
拼进 ``token_pattern`` 正则；把 jieba 的分词结果塞进正则不可行。

所以有两种诚实的选择，本模块把两种都实现，由 ``scripts/eval_retrieval.py`` 实测：

| 实现 | 分词 | 优点 | 代价 |
|---|---|---|---|
| ``llamaindex`` | 汉字级 ``token_pattern`` | 全用官方组件 | 单字粒度，精度下降 |
| ``jieba`` | jieba 词级 | 与现状一致（实测 bm25 通道 80%） | BM25 这层仍是自写 |

## 两种实现都保留「按构念域独立索引」

串味不报错，只会让引用「有点对但不对题」——所以每个 namespace 单独建，检索时合并。
（``llama-index-retrievers-bm25`` 支持 ``filters=MetadataFilters``，但它是**构造期**
预计算 ``corpus_weight_mask``，无法按请求变化；见 ``_build_llama_child`` 的注释。）
"""

from __future__ import annotations

import threading
from dataclasses import dataclass, field
from typing import Any, Literal

from llama_index.core.retrievers import BaseRetriever
from llama_index.core.schema import NodeWithScore, QueryBundle, TextNode
from pydantic import PrivateAttr

from app.core.config import get_settings
from app.core.logging import get_logger
from app.rag.index import get_collection

logger = get_logger(__name__)

Bm25Backend = Literal["jieba", "llamaindex"]

#: 汉字级分词的正则：一个汉字一个 token，ASCII 按词。
#: 走 ``bm25s.tokenize(token_pattern=...)``，故只能表达「切什么」，不能表达「怎么切」。
CJK_TOKEN_PATTERN = r"(?u)[\u4e00-\u9fff]|[A-Za-z0-9_]+"


def _tokenize(text: str) -> list[str]:
    """jieba 词级分词。jieba 首次加载约 0.4s，之后走自身缓存。"""
    import jieba

    return [tok for tok in jieba.lcut(text) if tok.strip()]


@dataclass
class _NamespaceNodes:
    """某个构念域的全部块，按 ``chunk_id`` 排序。

    排序必须固定：Chroma 不保证 ``get()`` 返回顺序稳定，顺序一变 BM25 的
    「同分文档」次序就变，评测不可复现。
    """

    namespace: str
    chunk_ids: list[str] = field(default_factory=list)
    texts: list[str] = field(default_factory=list)
    metadatas: list[dict] = field(default_factory=list)

    @property
    def size(self) -> int:
        return len(self.chunk_ids)

    def text_nodes(self) -> list[TextNode]:
        return [
            TextNode(id_=cid, text=text, metadata=meta)
            for cid, text, meta in zip(self.chunk_ids, self.texts, self.metadatas)
        ]


_nodes_cache: dict[str, _NamespaceNodes] = {}
_jieba_cache: dict[str, Any] = {}
_llama_cache: dict[str, BaseRetriever] = {}
_lock = threading.Lock()


def reload_cache() -> None:
    """丢弃全部 BM25 缓存，下次检索时按当前集合重建。

    ``app/rag/retriever.py::reload_indexes`` 会调用本函数——入库不重建索引，
    新增的块在 BM25 通道里永远搜不到，而向量通道却搜得到，两条通道会静默不一致。
    """
    with _lock:
        n = len(_jieba_cache) + len(_llama_cache)
        _nodes_cache.clear()
        _jieba_cache.clear()
        _llama_cache.clear()
    if n:
        logger.info("已清空 BM25 缓存（%d 个 namespace 索引已失效）", n)


def _load_nodes(namespace: str) -> _NamespaceNodes:
    """从 Chroma 读出某构念域的全部块。以 Chroma 为唯一真源。"""
    cached = _nodes_cache.get(namespace)
    if cached is not None:
        return cached

    with _lock:
        cached = _nodes_cache.get(namespace)
        if cached is not None:
            return cached

        col = get_collection(create=False)
        if col is None:
            out = _NamespaceNodes(namespace=namespace)
            _nodes_cache[namespace] = out
            return out

        got = col.get(where={"kb_namespace": namespace}, include=["documents", "metadatas"])
        ids = list(got.get("ids") or [])
        docs = [d or "" for d in (got.get("documents") or [])]
        metas = [m or {} for m in (got.get("metadatas") or [])]

        order = sorted(range(len(ids)), key=lambda i: ids[i])
        out = _NamespaceNodes(
            namespace=namespace,
            chunk_ids=[ids[i] for i in order],
            texts=[docs[i] for i in order],
            metadatas=[metas[i] for i in order],
        )
        _nodes_cache[namespace] = out
        logger.info("BM25 语料就绪 | namespace=%s | %d 块", namespace, out.size)
        return out


# ── 实现一：jieba 词级 + rank_bm25 ──────────────────────────


class JiebaBm25Retriever(BaseRetriever):
    """单个构念域的 jieba 词级 BM25 检索器。"""

    _namespace: str = PrivateAttr()
    _chunk_ids: list[str] = PrivateAttr()
    _texts: list[str] = PrivateAttr()
    _metadatas: list[dict] = PrivateAttr()
    _bm25: Any = PrivateAttr()
    _top_k: int = PrivateAttr()

    def __init__(
        self,
        *,
        namespace: str,
        nodes: _NamespaceNodes,
        bm25: Any,
        top_k: int,
    ) -> None:
        super().__init__()
        self._namespace = namespace
        self._chunk_ids = nodes.chunk_ids
        self._texts = nodes.texts
        self._metadatas = nodes.metadatas
        self._bm25 = bm25
        self._top_k = top_k

    @property
    def namespace(self) -> str:
        return self._namespace

    @property
    def size(self) -> int:
        return len(self._chunk_ids)

    @property
    def chunk_ids(self) -> list[str]:
        return list(self._chunk_ids)

    def _retrieve(self, query_bundle: QueryBundle) -> list[NodeWithScore]:
        tokens = _tokenize(query_bundle.query_str)
        if not tokens or self._bm25 is None:
            return []

        scores = self._bm25.get_scores(tokens)
        order = sorted(range(len(scores)), key=lambda i: scores[i], reverse=True)

        out: list[NodeWithScore] = []
        for i in order:
            if scores[i] <= 0.0:
                break  # 已降序，后面全是 0 分（与该域无词重叠）
            out.append(
                NodeWithScore(
                    node=TextNode(
                        id_=self._chunk_ids[i],
                        text=self._texts[i],
                        metadata=self._metadatas[i],
                    ),
                    score=float(scores[i]),
                )
            )
            if len(out) >= self._top_k:
                break
        return out


def _build_jieba_child(namespace: str, top_k: int) -> JiebaBm25Retriever | None:
    from rank_bm25 import BM25Okapi

    nodes = _load_nodes(namespace)
    if nodes.size == 0:
        return None

    with _lock:
        bm25 = _jieba_cache.get(namespace)
        if bm25 is None:
            tokenized = [_tokenize(t) for t in nodes.texts]
            non_empty = [t for t in tokenized if t]
            bm25 = BM25Okapi(non_empty) if non_empty else None
            _jieba_cache[namespace] = bm25
            logger.info("jieba BM25 索引就绪 | namespace=%s | %d 块", namespace, nodes.size)

    return JiebaBm25Retriever(namespace=namespace, nodes=nodes, bm25=bm25, top_k=top_k)


# ── 实现二：LlamaIndex BM25Retriever（bm25s，汉字级）──────────


def _build_llama_child(namespace: str, top_k: int) -> BaseRetriever | None:
    """单构念域的官方 ``BM25Retriever``。

    这里刻意用「每域一个实例」而不是「一个实例 + ``filters``」：
    ``BM25Retriever`` 的 ``filters`` 在**构造时**就烘成 ``corpus_weight_mask``，
    检索时读的是这个固定掩码，无法按请求改写。要按请求换域就只能重建——
    而重建会重新分词整个语料（1319 块），比按域缓存贵得多。
    """
    from llama_index.retrievers.bm25 import BM25Retriever

    cached = _llama_cache.get(namespace)
    if cached is not None:
        return cached

    nodes = _load_nodes(namespace)
    if nodes.size == 0:
        return None

    with _lock:
        cached = _llama_cache.get(namespace)
        if cached is None:
            cached = BM25Retriever.from_defaults(
                nodes=nodes.text_nodes(),
                similarity_top_k=top_k,
                skip_stemming=True,  # 英文词干器对中文无意义，还会打乱汉字
                token_pattern=CJK_TOKEN_PATTERN,
            )
            _llama_cache[namespace] = cached
            logger.info("llamaindex BM25 索引就绪 | namespace=%s | %d 块", namespace, nodes.size)
        return cached


# ── 多域合并 ────────────────────────────────────────────────


class MergingBm25Retriever(BaseRetriever):
    """把多个构念域的独立索引合并成一个检索器。

    合并只看分数：各域用的是同一个 BM25 公式，所以分数同量纲，可以直接比。
    （这也是「必须每域独立建索引」的另一个理由——若混建一个索引，域内的
    idf 会被其他域的语料污染，分数就不再可比。）
    """

    _children: list[BaseRetriever] = PrivateAttr()
    _top_k: int = PrivateAttr()

    def __init__(self, *, children: list[BaseRetriever], top_k: int) -> None:
        super().__init__()
        self._children = children
        self._top_k = top_k

    @property
    def children_count(self) -> int:
        return len(self._children)

    def _retrieve(self, query_bundle: QueryBundle) -> list[NodeWithScore]:
        merged: dict[str, NodeWithScore] = {}
        for child in self._children:
            for nws in child.retrieve(query_bundle):
                cid = nws.node.node_id
                prev = merged.get(cid)
                if prev is None or (nws.score or 0.0) > (prev.score or 0.0):
                    merged[cid] = nws

        ranked = sorted(
            merged.values(),
            # 同分时用 chunk_id 兜底，保证结果可复现
            key=lambda n: (-(n.score or 0.0), n.node.node_id),
        )
        return ranked[: self._top_k]


def build_retriever(
    namespaces: list[str],
    top_k: int,
    *,
    backend: Bm25Backend | None = None,
) -> MergingBm25Retriever | None:
    """按构念域构造 BM25 检索器。``namespaces`` 里没有语料的域会被跳过。"""
    s = get_settings()
    chosen: Bm25Backend = backend or s.bm25_backend

    children: list[BaseRetriever] = []
    for ns in namespaces:
        child = (
            _build_llama_child(ns, top_k)
            if chosen == "llamaindex"
            else _build_jieba_child(ns, top_k)
        )
        if child is not None:
            children.append(child)

    if not children:
        return None
    return MergingBm25Retriever(children=children, top_k=top_k)


def namespace_index_stats(
    namespace: str,
    *,
    backend: Bm25Backend | None = None,
) -> tuple[str, int, list[str]]:
    """``(namespace, 块数, chunk_ids)``，给测试断言「每域独立、互不相交」用。"""
    s = get_settings()
    chosen: Bm25Backend = backend or s.bm25_backend

    if chosen == "llamaindex":
        child = _build_llama_child(namespace, 5)
        if child is None:
            return namespace, 0, []
        corpus = getattr(child, "corpus", []) or []
        ids = [str(tok.get("node_id", "")) for tok in corpus]
        return namespace, len(ids), ids

    child = _build_jieba_child(namespace, 5)
    if child is None:
        return namespace, 0, []
    return namespace, child.size, child.chunk_ids
