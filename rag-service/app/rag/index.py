"""Chroma 持久化客户端与集合管理。

集合名与 metadata 见 §8.2；元数据字段见 §8.3（:599）：

| 元数据 | 用途 |
|---|---|
| ``source`` / ``kb_namespace`` / ``section_title`` / ``chunk_id`` | 引用回显与过滤检索 |

``kb_namespace`` 的 where 过滤**只作用于向量通道**；BM25 通道没有过滤能力，
必须按 namespace 用独立索引（ADR-0005）。
"""

from __future__ import annotations

import threading
from typing import Any

from app.core.config import get_settings
from app.core.logging import get_logger

logger = get_logger(__name__)

_client: Any = None
_client_lock = threading.Lock()


def get_client():
    """持久化 Chroma 客户端（单例）。"""
    global _client

    if _client is None:
        with _client_lock:
            if _client is None:
                import chromadb
                from chromadb.config import Settings as ChromaSettings

                s = get_settings()
                s.chroma_path.mkdir(parents=True, exist_ok=True)

                _client = chromadb.PersistentClient(
                    path=str(s.chroma_path),
                    settings=ChromaSettings(anonymized_telemetry=False),
                )
                logger.info("Chroma 客户端就绪 | path=%s", s.chroma_path)
    return _client


def get_collection(*, create: bool = True):
    """取（或建）知识库集合。

    距离用 ``cosine``：``EmbeddingModel.encode`` 已做 L2 归一化，
    归一化向量的余弦相似度即内积，检索时无需再变换。
    """
    s = get_settings()
    client = get_client()

    if create:
        try:
            return client.get_or_create_collection(
                name=s.collection_name,
                metadata={"hnsw:space": "cosine"},
            )
        except Exception as exc:  # noqa: BLE001 —— 启动期要给出可读原因
            raise RuntimeError(
                f"创建/打开 Chroma 集合 {s.collection_name!r} 失败：{exc}。"
                f"若目录 {s.chroma_path} 里已有同名集合且维度不同，"
                "需先删除该目录重建。"
            ) from exc

    try:
        return client.get_collection(name=s.collection_name)
    except Exception:  # noqa: BLE001
        return None


def collection_count() -> int:
    """/health 用。集合不存在时返回 0，不抛错。"""
    try:
        col = get_collection(create=False)
        return 0 if col is None else int(col.count())
    except Exception as exc:  # noqa: BLE001
        logger.warning("读取集合计数失败：%s", exc)
        return 0


def health_ok() -> bool:
    """P0-1 验收项「Chroma 持久化目录可写」的可编程形态。"""
    try:
        get_settings().chroma_path.mkdir(parents=True, exist_ok=True)
        probe = get_settings().chroma_path / ".write_probe"
        probe.write_text("ok", encoding="utf-8")
        probe.unlink()
        get_client()
        return True
    except Exception as exc:  # noqa: BLE001
        logger.error("Chroma 健康检查失败：%s", exc)
        return False


# ── LlamaIndex 装配（ADR-0009）──────────────────────────────

_vector_store: Any = None
_vector_store_lock = threading.Lock()

_index: Any = None
_index_lock = threading.Lock()


def get_vector_store():
    """把现有 Chroma 集合包成 LlamaIndex 的 ``ChromaVectorStore``。

    ⚠️ **必须传 ``chroma_collection``**。``ChromaVectorStore`` 的默认分支
    （``chroma_collection is None``）会去 ``chromadb.HttpClient(host, port, ...)``
    再 ``get_or_create_collection(collection_name)``——那条路：

    1. 连不上我们的 ``PersistentClient``（一个走 HTTP，一个走本地目录）；
    2. 会丢掉 ``metadata={"hnsw:space": "cosine"}``，于是距离度量变成默认 L2，
       而我们的向量都是 L2 归一化的，L2 与 cosine 在这里单调等价但分数含义不同。

    另外它对距离→相似度的换算是 ``math.exp(-distance)``，不是旧实现的
    ``1 - distance``：排序不变，绝对分数不可跨版本比较。
    """
    global _vector_store

    if _vector_store is None:
        with _vector_store_lock:
            if _vector_store is None:
                from llama_index.vector_stores.chroma import ChromaVectorStore

                col = get_collection(create=False)
                if col is None:
                    return None
                _vector_store = ChromaVectorStore(chroma_collection=col)
                logger.info("ChromaVectorStore 就绪 | collection=%s", get_settings().collection_name)
    return _vector_store


def get_llama_index():
    """``VectorStoreIndex`` 单例。

    这个对象**不持有数据**——每次查询都从 ``ChromaVectorStore`` 读，所以入库后
    不需要重建它；需要重建的只有 BM25 的进程内索引（见 ``rag/bm25.py``）。
    """
    global _index

    if _index is None:
        with _index_lock:
            if _index is None:
                from llama_index.core import VectorStoreIndex

                from app.rag.embedding import get_llama_embedding

                store = get_vector_store()
                if store is None:
                    return None
                _index = VectorStoreIndex.from_vector_store(
                    store,
                    embed_model=get_llama_embedding(),
                )
                logger.info("VectorStoreIndex 就绪")
    return _index
