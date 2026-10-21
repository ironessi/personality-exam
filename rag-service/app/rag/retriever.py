"""混合检索：Chroma 向量通道 + BM25 关键词通道，融合交给 LlamaIndex（ADR-0009）。

## 两条通道的过滤能力不对等 —— 这是最容易写错的地方

向量通道的 ``where={"kb_namespace": {"$in": ["typology", "common"]}}`` 是检索期过滤，
换一台机器、换一个请求都能变。BM25 没有这个能力：它是一组已建好的词频矩阵，过滤
只能发生在**建索引的时候**。所以 BM25 **按构念域各建一个独立索引**，检索时先按域选
索引，再在选中的索引内部排序——这段逻辑在 :mod:`app.rag.bm25`。

把多个域混进一个 BM25 索引、指望检索时过滤，会静默串味：类型学的 query 会召回症状
严重度的块。类型学与症状严重度的文档在措辞上高度重叠（都谈「情绪」「人际」「压力」），
所以这种串味不会报错，只会让引用看起来「有点对但不对题」。

``kb_namespace`` 的取值是**构念域**（ADR-0008），不是测评名——每个测评按
``test.domains`` 声明自己检索哪几个域（通常是「本域 + ``common``」）。

## 融合方式

两通道分数**不同量纲**（向量通道是 ``exp(-distance)`` ∈ (0, 1]，BM25 无上界），
直接加权求和需要逐域调参。这里用 **RRF（Reciprocal Rank Fusion）**——只吃排名不
吃分数：``rrf(chunk) += 1 / (k + rank)``。RRF 不需要归一化、不需要调权重，也不会
因为某个域的 BM25 分数整体偏大就压掉向量通道。

融合由 ``QueryFusionRetriever`` 完成，但它的**三个默认值都是陷阱**，必须全部显式覆盖：

| 参数 | 默认值 | 必须设成 | 不设的后果 |
|---|---|---|---|
| ``num_queries`` | ``4`` | ``1`` | 会调 LLM 生成 3 条查询变体——多 3 次 LLM 调用，同步快路径的 3 秒预算直接爆掉 |
| ``mode`` | ``SIMPLE`` | ``RECIPROCAL_RANK`` | SIMPLE 是直接拼列表，不去重、不看排名 |
| ``use_async`` | ``True`` | ``True`` | 保持默认即可，但要显式写出来，别让人以为是同步的 |

RRF 的 ``k = 60.0`` 是 LlamaIndex 内部**硬编码**的（源码引 Cormack 2009），不可配。
两条通道融合后的上界是 ``2/60``，不是 ``2/61``——rank 是 0 基。

## 状态

**P0-3 实现**（§9 :724）。验收要求**两条通道单独验**：``typology`` query 不召回
``symptom`` 块（§9 :783）。故 ``retrieve`` 保留 ``mode="vector"`` / ``"bm25"``，
把两条通道分开断言——只测 hybrid 无法证明是哪个通道串味。

## 一处诚实的能力缺失

``QueryFusionRetriever`` 融合后**只保留融合分**，不再暴露各通道的原始分（它按
``node.hash`` 去重累加，原始 ``NodeWithScore`` 被丢弃）。所以 hybrid 模式下
``vector_score`` / ``bm25_score`` 一律是 ``None``，只有 ``fused_score`` 有值。
要归因请用 ``mode="vector"`` / ``mode="bm25"`` 单独查，不要试图从 hybrid 结果反推。
"""

from __future__ import annotations

from dataclasses import dataclass, field
from typing import Any, Literal

from llama_index.core.llms import CompletionResponse, CustomLLM, LLMMetadata
from llama_index.core.retrievers import BaseRetriever
from llama_index.core.retrievers.fusion_retriever import FUSION_MODES, QueryFusionRetriever
from llama_index.core.schema import NodeWithScore, QueryBundle
from llama_index.core.vector_stores import MetadataFilter, MetadataFilters, FilterOperator

from app.core.config import get_settings
from app.core.logging import get_logger
from app.rag import bm25
from app.rag.index import get_llama_index

logger = get_logger(__name__)

#: RRF 常数**默认值**。真值以 ``settings.rrf_k`` 为准（这里只是给测试算期望分用）。
#:
#: ⚠️ **不要改回 60。** LlamaIndex 内部硬编码 60，那是给上千条候选设计的经验值；
#: 本服务每通道只取 ``top_k * 3`` 条（生产 15 条），60 会让 rank 的贡献落在
#: 1/60~1/75 这个极窄区间，融合退化成「谁在两个通道都露过面谁赢」。
#: 实测 k=60 → MRR 0.716，k=3 → 0.769，命中率与越域命中都不变。
#: 完整推导与故障样例见 ``app/core/config.py`` 的 ``rrf_k`` 注释。
RRF_K = 3

#: 候选池放大倍数。融合只让条目往下掉，不会让没进候选的冒出来，
#: 所以每条通道都要多取一些，否则「两通道都排第 6」的块永远进不了前 5。
CANDIDATE_MULTIPLIER = 3

RetrievalMode = Literal["hybrid", "vector", "bm25"]


@dataclass
class RetrievedChunk:
    """检索结果。字段与 §8.2 的契约一致，metadata 里带原始四个元数据字段。"""

    chunk_id: str
    text: str
    source: str
    kb_namespace: str
    section_title: str = ""
    vector_score: float | None = None
    bm25_score: float | None = None
    fused_score: float = 0.0
    metadata: dict = field(default_factory=dict)


class _NodeIdFusionRetriever(QueryFusionRetriever):
    """RRF 去重键换成 ``node_id``，其余完全用官方实现。

    ## 为什么必须覆盖

    ``QueryFusionRetriever._reciprocal_rerank_fusion`` 用 ``node_with_score.node.hash``
    当去重键，而 ``TextNode.hash`` 是::

        sha256(str(self.text) + str(self.metadata))

    ``str(dict)`` **对键序敏感**。而同一块从 Chroma 的两条读路径出来，键序并不一致：

    - ``collection.query(...)`` 返回 ``{'section_title':…, 'kb_namespace':…, 'source':…, 'chunk_id':…}``
    - ``collection.get(...)``（``ChromaVectorStore.get_nodes`` 走这条）返回
      ``{'source':…, 'kb_namespace':…, 'chunk_id':…, 'section_title':…}``

    于是**同一个 chunk 在两条通道里 hash 不同**，RRF 会把它们当成两个条目各投一票，
    「两通道都排第一」永远拿不到 ``2/60``，融合退化成「按名次拼列表」——比单通道好
    不了多少，而且完全无法奖励「两条通道同时认为它相关」。

    这里只改去重键（``node_id`` 就是 ``chunk_id``，天然稳定且唯一），打分公式、排序、
    截断全部沿用官方实现。**唯一偏离官方的是 ``k``**：官方硬编码 60，本服务用
    ``settings.rrf_k``（默认 3），原因见 ``app/core/config.py`` 里 ``rrf_k`` 的注释。
    同分时用 ``node_id`` 兜底排序，保证结果可复现。
    """

    def _reciprocal_rerank_fusion(
        self, results: dict[tuple[str, int], list[NodeWithScore]]
    ) -> list[NodeWithScore]:
        # 运行期以 settings 为准，便于 eval 脚本用环境变量做 A/B；
        # 取不到（例如单测直接构造本类）则退回模块默认值。
        try:
            k = float(get_settings().rrf_k)
        except Exception:  # noqa: BLE001
            k = float(RRF_K)

        fused_scores: dict[str, float] = {}
        id_to_node: dict[str, NodeWithScore] = {}

        for nodes_with_scores in results.values():
            for rank, node_with_score in enumerate(
                sorted(nodes_with_scores, key=lambda x: x.score or 0.0, reverse=True)
            ):
                key = node_with_score.node.node_id
                id_to_node[key] = node_with_score
                fused_scores[key] = fused_scores.get(key, 0.0) + 1.0 / (rank + k)

        ranked = sorted(fused_scores.items(), key=lambda kv: (-kv[1], kv[0]))

        out: list[NodeWithScore] = []
        for key, score in ranked:
            nws = id_to_node[key]
            nws.score = score
            out.append(nws)
        return out


def _namespace_filter(namespaces: list[str]) -> MetadataFilters:
    return MetadataFilters(
        filters=[
            MetadataFilter(
                key="kb_namespace",
                operator=FilterOperator.IN,
                value=list(namespaces),
            )
        ]
    )


class _NoQueryRewriteLLM(CustomLLM):
    """占位 LLM：一旦被调用就抛错。

    ``QueryFusionRetriever.__init__`` 里有一行 ``resolve_llm(llm, ...) if llm else
    Settings.llm``——**即使 ``num_queries=1`` 也会去解析 LLM**。不传 ``llm`` 就会
    落到全局 ``Settings.llm``，而它的默认值是 OpenAI，于是报
    ``ValueError: No API key found for OpenAI``，检索直接失败。

    传一个「故意会炸」的 LLM 而不是 ``MockLLM``，是为了让约束变成可执行的：
    ``num_queries=1`` 时它一次都不会被调用；真有人改成 >1，这里会立刻抛错，
    而不是悄悄用假查询去检索、把错误藏进结果排序里。
    """

    @classmethod
    def class_name(cls) -> str:
        return "NoQueryRewriteLLM"

    @property
    def metadata(self) -> LLMMetadata:
        return LLMMetadata(
            context_window=1,
            num_output=1,
            model_name="no-query-rewrite",
            is_chat_model=False,
        )

    def complete(self, prompt: str, formatted: bool = False, **kwargs: Any) -> CompletionResponse:
        raise RuntimeError(
            "检索层不允许调用 LLM。查询改写（num_queries>1）会额外消耗 3 次 LLM 调用，"
            "并让 ADR-0007 的 3 秒同步预算失效。若确需改写，请先改 ADR。"
        )

    def stream_complete(
        self, prompt: str, formatted: bool = False, **kwargs: Any
    ) -> CompletionResponse:
        raise RuntimeError("检索层不允许调用 LLM。")


def _vector_retriever(namespaces: list[str], limit: int) -> BaseRetriever | None:
    index = get_llama_index()
    if index is None:
        return None

    from llama_index.core.indices.vector_store.retrievers import VectorIndexRetriever

    return VectorIndexRetriever(
        index=index,
        similarity_top_k=limit,
        filters=_namespace_filter(namespaces),
    )


def _bm25_retriever(
    namespaces: list[str],
    limit: int,
    backend: bm25.Bm25Backend | None = None,
) -> BaseRetriever | None:
    return bm25.build_retriever(namespaces, limit, backend=backend)


def _to_chunk(
    nws: NodeWithScore,
    *,
    vector_score: float | None = None,
    bm25_score: float | None = None,
    fused_score: float = 0.0,
) -> RetrievedChunk:
    node = nws.node
    # LlamaIndex 往 metadata 里写了几个自己的字段：``_node_content`` 是整份节点的
    # JSON 副本（正文再存一遍，体积翻倍）、``document_id``/``doc_id``/``ref_doc_id``
    # 是它的内部指针。这些一律不进业务的 ``metadata``——一旦透传出去，
    # 引用里会带上内部 JSON，而契约只认 ``source``/``kb_namespace``/``section_title``/``chunk_id``。
    meta = {
        k: v
        for k, v in (getattr(node, "metadata", None) or {}).items()
        if not k.startswith("_") and k not in {"document_id", "doc_id", "ref_doc_id"}
    }
    text = ""
    getter = getattr(node, "get_content", None)
    if callable(getter):
        try:
            text = getter() or ""
        except Exception:  # noqa: BLE001 —— 取不到正文不该让整次检索失败
            text = ""
    if not text:
        text = getattr(node, "text", "") or ""

    return RetrievedChunk(
        chunk_id=node.node_id,
        text=text,
        source=str(meta.get("source", "")),
        kb_namespace=str(meta.get("kb_namespace", "")),
        section_title=str(meta.get("section_title", "")),
        vector_score=vector_score,
        bm25_score=bm25_score,
        fused_score=fused_score,
        metadata=meta,
    )


def _run(retriever: BaseRetriever, query: str) -> list[NodeWithScore]:
    return retriever.retrieve(QueryBundle(query_str=query))


# ── 索引重建 ────────────────────────────────────────────────


def reload_indexes() -> None:
    """入库后必须调用：BM25 的进程内索引要重建。

    向量通道不用重建——``VectorStoreIndex`` 不持有数据，每次查询都从 Chroma 读。
    但 BM25 索引是进程内的一份词频矩阵，不重建的话新增块在 BM25 通道永远搜不到，
    而向量通道搜得到，两条通道静默不一致。``app/ingest/pipeline.py`` 会探测调用。
    """
    bm25.reload_cache()
    logger.info("检索索引已重建（BM25 缓存清空）")


# ── 对外入口 ────────────────────────────────────────────────


def retrieve(
    queries: list[str],
    *,
    namespaces: list[str] | None = None,
    top_k: int | None = None,
    mode: RetrievalMode = "hybrid",
    bm25_backend: bm25.Bm25Backend | None = None,
) -> list[RetrievedChunk]:
    """混合检索。

    Args:
        queries: 一到多条查询串。多条时每条独立检索再按 RRF 合并——同一个块被多条
            查询同时命中，说明它与整份档案多处相关，应该排更前。
        namespaces: 构念域列表，通常来自 ``test.domains``。``None`` 用全部 8 域。
        top_k: 返回条数，``None`` 取 ``settings.top_k``。
        mode: ``hybrid`` / ``vector`` / ``bm25``。后两者专供验收把通道分开断言。
        bm25_backend: 只给 ``scripts/eval_retrieval.py`` 做 A/B 实测用；``None`` 取
            ``settings.bm25_backend``。生产路径不该传它。
    """
    s = get_settings()
    k = top_k if top_k is not None else s.top_k
    # None = 调用方没指定，用全部域；[] = 调用方明确要求"不检索任何域"。
    # 两者不能混为一谈：曾经写成 `if domains: domains = ...` 的兜底，
    # 结果是空域退化成全库检索——越域过滤形同不存在，且不报错。
    domains = list(namespaces) if namespaces is not None else list(s.kb_namespaces)
    if not domains:
        return []
    limit = max(k * CANDIDATE_MULTIPLIER, k)

    # 空白查询没有检索意图。向量通道会给空格算出一个"合法"的向量并返回一堆结果，
    # 所以必须在入口滤掉，而不是指望下游各通道自己处理。
    queries = [q for q in queries if q and q.strip()]
    if not queries:
        return []

    if mode == "vector":
        retriever = _vector_retriever(domains, limit)
        if retriever is None:
            return []
        merged: dict[str, RetrievedChunk] = {}
        for q in queries:
            for nws in _run(retriever, q):
                chunk = _to_chunk(nws, vector_score=nws.score)
                prev = merged.get(chunk.chunk_id)
                if prev is None or (chunk.vector_score or 0.0) > (prev.vector_score or 0.0):
                    merged[chunk.chunk_id] = chunk
        ranked = sorted(merged.values(), key=lambda c: (-(c.vector_score or 0.0), c.chunk_id))
        return ranked[:k]

    if mode == "bm25":
        retriever = _bm25_retriever(domains, limit, bm25_backend)
        if retriever is None:
            return []
        merged = {}
        for q in queries:
            for nws in _run(retriever, q):
                chunk = _to_chunk(nws, bm25_score=nws.score)
                prev = merged.get(chunk.chunk_id)
                if prev is None or (chunk.bm25_score or 0.0) > (prev.bm25_score or 0.0):
                    merged[chunk.chunk_id] = chunk
        ranked = sorted(merged.values(), key=lambda c: (-(c.bm25_score or 0.0), c.chunk_id))
        return ranked[:k]

    # hybrid
    children: list[BaseRetriever] = []
    vec = _vector_retriever(domains, limit)
    if vec is not None:
        children.append(vec)
    bm = _bm25_retriever(domains, limit, bm25_backend)
    if bm is not None:
        children.append(bm)

    if not children:
        return []
    if len(children) == 1:
        # 只有一条通道可用（例如某一域还没语料），没必要套融合器
        only = children[0]
        merged = {}
        for q in queries:
            for nws in _run(only, q):
                chunk = _to_chunk(
                    nws,
                    vector_score=nws.score if only is vec else None,
                    bm25_score=nws.score if only is bm else None,
                    fused_score=float(nws.score or 0.0),
                )
                merged.setdefault(chunk.chunk_id, chunk)
        ranked = sorted(merged.values(), key=lambda c: (-c.fused_score, c.chunk_id))
        return ranked[:k]

    fusion = _NodeIdFusionRetriever(
        retrievers=children,
        mode=FUSION_MODES.RECIPROCAL_RANK,
        num_queries=1,  # 绝不能 >1：会调 LLM 生成查询变体
        similarity_top_k=k,
        use_async=True,
        # 必须显式传，否则构造时会去解析全局 Settings.llm（默认 OpenAI）并直接报错
        llm=_NoQueryRewriteLLM(),
    )

    merged: dict[str, RetrievedChunk] = {}
    for q in queries:
        for nws in _run(fusion, q):
            chunk = _to_chunk(nws, fused_score=float(nws.score or 0.0))
            prev = merged.get(chunk.chunk_id)
            if prev is None or chunk.fused_score > prev.fused_score:
                merged[chunk.chunk_id] = chunk

    ranked = sorted(merged.values(), key=lambda c: (-c.fused_score, c.chunk_id))
    return ranked[:k]

