"""本地 BGE 向量模型封装（单例 + 显式预热）。

## 为什么必须预热

在开发机 `/Users/iron/rag` 内实测（BAAI/bge-small-zh-v1.5，512 维）：

| 阶段 | 耗时 |
|---|---|
| `import torch` + `import sentence_transformers` | **5.04 s** |
| `SentenceTransformer(path)` 冷 | 0.21 s |
| `SentenceTransformer(path)` 热 | 0.08 s |
| 加载后**第一次** `encode(3 条)` | **4.27 s** |
| 第二次 | 0.029 s |
| 第三次 | 0.008 s |

也就是说「模型加载」本身很便宜，真正贵的是**第一次前向**（线程池与算子首次
初始化）。ADR-0007 的同步快路径只有 **3 秒**，要是把首次 encode 留给第一个真实
请求，它一定超时——每次服务重启都会让第一个用户白等 3 秒再转异步。

所以 ``warmup()`` 必须在 lifespan 里跑完再开始接客。这也解释了为什么 §11 声称的
「Embedding 与检索 < 0.5 s」只对**稳态**成立。
"""

from __future__ import annotations

import threading
import time
from typing import Any, Iterable

from llama_index.core.embeddings import BaseEmbedding

from app.core.config import get_settings
from app.core.logging import get_logger

logger = get_logger(__name__)


def _detect_dimension(model: Any) -> int:
    """取真实向量维度。

    sentence-transformers 5.x 把 ``get_sentence_embedding_dimension()`` 改名成
    ``get_embedding_dimension()``；优先新名以免触发 FutureWarning，同时兼容旧版。
    """
    for attr in ("get_embedding_dimension", "get_sentence_embedding_dimension"):
        fn = getattr(model, attr, None)
        if callable(fn):
            return int(fn())
    raise RuntimeError(
        "无法从模型取到向量维度：既没有 get_embedding_dimension() "
        "也没有 get_sentence_embedding_dimension()"
    )


class EmbeddingModel:
    """进程内单例。线程安全：encode 上加锁，避免并发首次前向重复初始化。"""

    def __init__(self, model_path: str, expected_dim: int) -> None:
        self.model_path = model_path
        self.expected_dim = expected_dim
        self._model: Any = None
        self._dim: int | None = None
        self._lock = threading.Lock()

    # ── 生命周期 ────────────────────────────────────────────

    @property
    def loaded(self) -> bool:
        return self._model is not None

    @property
    def dim(self) -> int | None:
        return self._dim

    def load(self) -> None:
        """幂等加载 + 维度断言。维度不符即抛错，绝不静默产出错维向量。"""
        if self._model is not None:
            return

        with self._lock:
            if self._model is not None:
                return

            from sentence_transformers import SentenceTransformer

            t0 = time.perf_counter()
            model = SentenceTransformer(self.model_path)
            load_s = time.perf_counter() - t0

            dim = _detect_dimension(model)
            if dim != self.expected_dim:
                raise RuntimeError(
                    f"向量维度不符：模型 {self.model_path} 实际 {dim} 维，"
                    f"配置 EMBEDDING_DIM={self.expected_dim}。"
                    "维度不一致会让 Chroma 集合无法复用（已有集合会直接报错）。"
                )

            self._model = model
            self._dim = dim
            logger.info(
                "向量模型加载完成 | path=%s dim=%d 耗时=%.3fs",
                self.model_path,
                dim,
                load_s,
            )

    def warmup(self, *, rounds: int = 2) -> float:
        """跑 dummy encode 把首次前向的开销提前烧掉。

        返回首次 encode 的秒数（即被省掉的那部分延迟），便于启动日志对照。
        """
        self.load()
        assert self._model is not None

        samples = [
            "外倾者倾向于从社交互动中获取心理能量",
            "躯体化因子偏高可能表现为头痛与肌肉酸痛",
            "本结果仅供参考，不构成医学诊断",
        ]

        first_s = 0.0
        for i in range(rounds):
            t0 = time.perf_counter()
            self._model.encode(samples, normalize_embeddings=True)
            elapsed = time.perf_counter() - t0
            if i == 0:
                first_s = elapsed

        logger.info(
            "向量模型预热完成 | 首次 encode %d 条=%.3fs（本该由第一个用户承担）",
            len(samples),
            first_s,
        )
        return first_s

    # ── 推理 ────────────────────────────────────────────────

    def encode(self, texts: Iterable[str]) -> list[list[float]]:
        """把文本编码成**已 L2 归一化**的向量。

        归一化后余弦相似度等于内积，Chroma 集合用 ``hnsw:space=cosine`` 即可。
        """
        items = list(texts)
        if not items:
            return []

        self.load()
        assert self._model is not None

        with self._lock:
            vectors = self._model.encode(items, normalize_embeddings=True)

        # numpy.ndarray.tolist()；不引入 numpy 类型标注以免上层依赖 numpy
        return [list(map(float, row)) for row in vectors]


_model_singleton: EmbeddingModel | None = None
_singleton_lock = threading.Lock()


def get_embedding_model() -> EmbeddingModel:
    """进程级单例。参数取自 Settings，改 .env 需重启。"""
    global _model_singleton

    if _model_singleton is None:
        with _singleton_lock:
            if _model_singleton is None:
                s = get_settings()
                _model_singleton = EmbeddingModel(
                    model_path=s.embedding_model_path,
                    expected_dim=s.embedding_dim,
                )
    return _model_singleton


# ── LlamaIndex 适配器（ADR-0009）────────────────────────────
#
# 为什么是「薄包装」而不是直接用 SentenceTransformerEmbedding：
#
# 1. **预热控制权**。ADR-0007 的同步快路径只有 3 秒，而首次 encode 实测 4.27s
#    （见模块开头表格）。必须由 lifespan 显式 `warmup()` 把这次开销提前烧掉。
#    `SentenceTransformerEmbedding` 在第一次调用时才 load，预热点就散进了请求路径。
# 2. **维度断言**。`EmbeddingModel.load()` 在维度不符时直接抛错；错误的维度会让
#    Chroma 集合彻底不可用（已有集合不会报错，只会返回错的结果）。
# 3. **单例**。上面那个单例是进程级的，多次构造 `SentenceTransformerEmbedding`
#    会重复占内存（torch 模型不是轻对象）。
#
# 所以这里只实现 LlamaIndex 要求的四个向量化钩子，实际计算全部转给单例。


class BGEEmbedding(BaseEmbedding):
    """把进程内 BGE 单例暴露成 LlamaIndex 的 ``BaseEmbedding``。"""

    @classmethod
    def class_name(cls) -> str:
        return "BGEEmbedding"

    def _get_query_embedding(self, query: str) -> list[float]:
        return get_embedding_model().encode([query])[0]

    def _get_text_embedding(self, text: str) -> list[float]:
        return get_embedding_model().encode([text])[0]

    def _get_text_embeddings(self, texts: list[str]) -> list[list[float]]:
        return get_embedding_model().encode(texts)

    async def _aget_query_embedding(self, query: str) -> list[float]:
        return self._get_query_embedding(query)

    async def _aget_text_embedding(self, text: str) -> list[float]:
        return self._get_text_embedding(text)

    async def _aget_text_embeddings(self, texts: list[str]) -> list[list[float]]:
        return self._get_text_embeddings(texts)


_llama_embedding: BGEEmbedding | None = None
_llama_embedding_lock = threading.Lock()


def get_llama_embedding() -> BGEEmbedding:
    """LlamaIndex 侧的单例。与 ``get_embedding_model()`` 指向同一个底层模型。"""
    global _llama_embedding

    if _llama_embedding is None:
        with _llama_embedding_lock:
            if _llama_embedding is None:
                s = get_settings()
                _llama_embedding = BGEEmbedding(model_name=s.embedding_model_path)
    return _llama_embedding
