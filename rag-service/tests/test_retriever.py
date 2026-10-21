"""混合检索测试（P0-3）。

最关键的一条：**跨构念域过滤必须两条通道单独验** ——
``typology`` 的 query 不得召回 ``symptom`` 的块。

向量通道靠 Chroma 的 ``MetadataFilter`` 在检索期过滤；BM25 通道没有过滤能力，
必须按 namespace 用独立索引（ADR-0005 / ADR-0009）。只验向量通道会漏掉 BM25 串味，
而 BM25 串味不报错，只是让引用「有点对但不对题」。

⚠️ namespace 是**构念域**（ADR-0008），不是测评名：``typology`` 覆盖 MBTI/DISC/
TrueColors/BemSexRoleInventory 四个测评，``symptom`` 覆盖 SCL-90 等 11 个。

BM25 有两套实现（ADR-0009：jieba 词级 vs LlamaIndex 官方汉字级），凡是与 BM25
有关的断言都对两套实现各跑一遍——只测默认实现等于没测另一套。

这些用例依赖已入库的语料。集合为空时整体跳过（先跑 ``scripts/ingest_cli.py``），
避免在没有语料的机器上把「环境没准备好」误报成「检索坏了」。
"""

from __future__ import annotations

import pytest

from app.rag.bm25 import namespace_index_stats
from app.rag.index import collection_count
from app.rag.retriever import RRF_K, reload_indexes, retrieve

pytestmark = pytest.mark.skipif(
    collection_count() == 0,
    reason="需要先入库语料：/Users/iron/rag/bin/python scripts/ingest_cli.py",
)

TYPO_QUERY = ["16型人格 INTJ 的认知功能与类型动力学 解释"]
TYPO_DOMAINS = ["typology", "common"]

SYMPTOM_QUERY = ["SCL-90 躯体化因子 得分偏高 说明什么"]
SYMPTOM_DOMAINS = ["symptom", "common"]

BACKENDS = ["jieba", "llamaindex"]


# ── 跨域过滤 ────────────────────────────────────────────────


def test_typology_query_does_not_recall_symptom_vector_channel():
    hits = retrieve(TYPO_QUERY, namespaces=TYPO_DOMAINS, top_k=10, mode="vector")
    assert hits, "向量通道应至少返回一条"
    assert {h.kb_namespace for h in hits} <= set(TYPO_DOMAINS)


def test_symptom_query_does_not_recall_typology_vector_channel():
    hits = retrieve(SYMPTOM_QUERY, namespaces=SYMPTOM_DOMAINS, top_k=10, mode="vector")
    assert hits
    assert {h.kb_namespace for h in hits} <= set(SYMPTOM_DOMAINS)


@pytest.mark.parametrize("backend", BACKENDS)
def test_typology_query_does_not_recall_symptom_bm25_channel(backend):
    hits = retrieve(
        TYPO_QUERY, namespaces=TYPO_DOMAINS, top_k=10, mode="bm25", bm25_backend=backend
    )
    assert hits, "BM25 通道应至少返回一条（该 query 与类型学语料有明显词重叠）"
    assert {h.kb_namespace for h in hits} <= set(TYPO_DOMAINS)


@pytest.mark.parametrize("backend", BACKENDS)
def test_symptom_query_does_not_recall_typology_bm25_channel(backend):
    hits = retrieve(
        SYMPTOM_QUERY, namespaces=SYMPTOM_DOMAINS, top_k=10, mode="bm25", bm25_backend=backend
    )
    assert hits
    assert {h.kb_namespace for h in hits} <= set(SYMPTOM_DOMAINS)


@pytest.mark.parametrize("backend", BACKENDS)
def test_hybrid_does_not_cross_domains(backend):
    """融合层不得把域外的块捞进来（融合只看排名，容易掩盖单通道的越界）。"""
    for query, domains in ((TYPO_QUERY, TYPO_DOMAINS), (SYMPTOM_QUERY, SYMPTOM_DOMAINS)):
        hits = retrieve(query, namespaces=domains, top_k=10, mode="hybrid", bm25_backend=backend)
        assert hits
        assert {h.kb_namespace for h in hits} <= set(domains)


# ── 每域独立索引 ────────────────────────────────────────────


@pytest.mark.parametrize("backend", BACKENDS)
def test_bm25_index_is_per_namespace(backend):
    """不能出现跨 namespace 共用的 BM25 索引（ADR-0005）。"""
    _, typo_size, typo_ids = namespace_index_stats("typology", backend=backend)
    _, symp_size, symp_ids = namespace_index_stats("symptom", backend=backend)

    assert typo_size > 0 and symp_size > 0
    # chunk_id 前缀就是 namespace，两域之间不应有任何交集
    assert set(typo_ids).isdisjoint(symp_ids)
    assert all(cid.startswith("typology_") for cid in typo_ids)
    assert all(cid.startswith("symptom_") for cid in symp_ids)


@pytest.mark.parametrize("backend", BACKENDS)
def test_bm25_domains_partition_the_collection(backend):
    """各域索引的块数之和 = 集合总数；不重不漏。

    这里遍历 ``settings.kb_namespaces`` 而不是写死域名——写死的话，将来任何
    一个域补上语料都会让这条本该"守门"的用例变成"报错"，而它想守的是
    「不重复、不遗漏」，不是「哪几个域非空」。
    """
    from app.core.config import get_settings

    total = collection_count()
    seen: list[str] = []
    for ns in get_settings().kb_namespaces:
        _, _size, ids = namespace_index_stats(ns, backend=backend)
        # 空域是合法状态（语料还没写），不该让用例失败
        assert all(cid.startswith(f"{ns}_") for cid in ids), f"{ns} 的索引里混进了别的域的块"
        seen.extend(ids)

    assert len(seen) == total, f"各域块数合计 {len(seen)} != 集合总数 {total}"
    assert len(set(seen)) == total, "不同域的 chunk_id 不应重复"


# ── 融合 ────────────────────────────────────────────────────


@pytest.mark.parametrize("backend", BACKENDS)
def test_hybrid_only_returns_chunks_from_either_channel(backend):
    """融合是「两条通道的并集再排序」，不该凭空多出第三来源。"""
    vec = {
        h.chunk_id
        for h in retrieve(TYPO_QUERY, namespaces=TYPO_DOMAINS, top_k=15, mode="vector")
    }
    bm = {
        h.chunk_id
        for h in retrieve(
            TYPO_QUERY, namespaces=TYPO_DOMAINS, top_k=15, mode="bm25", bm25_backend=backend
        )
    }
    hybrid = retrieve(TYPO_QUERY, namespaces=TYPO_DOMAINS, top_k=5, mode="hybrid", bm25_backend=backend)
    assert hybrid
    assert {h.chunk_id for h in hybrid} <= (vec | bm)
    assert all(0 < h.fused_score <= 2.0 / RRF_K + 1e-9 for h in hybrid)


@pytest.mark.parametrize("backend", BACKENDS)
def test_hybrid_fuses_by_node_id(backend):
    """RRF 的去重键必须是 ``node_id``，否则「两通道都排第一」拿不到 2/k。

    LlamaIndex 原版用 ``node.hash`` 当键，而 ``hash = sha256(str(text) + str(metadata))``
    对 dict 键序敏感；同一块经 ``collection.query`` 与 ``collection.get`` 出来键序不同，
    于是 hash 不等、融合退化成拼列表。这个用例就是为了钉死那个修复。
    """
    vec = retrieve(TYPO_QUERY, namespaces=TYPO_DOMAINS, top_k=5, mode="vector", bm25_backend=backend)
    bm = retrieve(TYPO_QUERY, namespaces=TYPO_DOMAINS, top_k=5, mode="bm25", bm25_backend=backend)
    assert vec and bm

    both_first = vec[0].chunk_id == bm[0].chunk_id
    hybrid = retrieve(TYPO_QUERY, namespaces=TYPO_DOMAINS, top_k=5, mode="hybrid", bm25_backend=backend)
    top = hybrid[0]

    if both_first:
        assert abs(top.fused_score - 2.0 / RRF_K) < 1e-9, (
            f"两通道首位都是 {vec[0].chunk_id}，融合分应为 2/{RRF_K}，"
            f"实际 {top.fused_score}——说明去重键没对上"
        )
    else:
        assert abs(top.fused_score - 1.0 / RRF_K) < 1e-9


@pytest.mark.parametrize("backend", BACKENDS)
def test_results_are_deterministic(backend):
    a = retrieve(TYPO_QUERY, namespaces=TYPO_DOMAINS, top_k=5, mode="hybrid", bm25_backend=backend)
    b = retrieve(TYPO_QUERY, namespaces=TYPO_DOMAINS, top_k=5, mode="hybrid", bm25_backend=backend)
    assert [h.chunk_id for h in a] == [h.chunk_id for h in b]


@pytest.mark.parametrize("backend", BACKENDS)
def test_reload_indexes_clears_cache_then_rebuilds(backend):
    _, before_size, before_ids = namespace_index_stats("typology", backend=backend)
    reload_indexes()
    _, after_size, after_ids = namespace_index_stats("typology", backend=backend)
    assert after_size == before_size
    assert after_ids == before_ids, "重建结果应与原来一致"


# ── 退化输入 ────────────────────────────────────────────────


def test_empty_namespaces_returns_nothing():
    assert retrieve(TYPO_QUERY, namespaces=[], top_k=5) == []
    assert retrieve([], namespaces=TYPO_DOMAINS, top_k=5) == []


def test_blank_query_returns_nothing():
    """向量通道会给空格算出一个「合法」的向量，所以空白查询必须在入口滤掉。"""
    assert retrieve(["", "   ", "\n"], namespaces=TYPO_DOMAINS, top_k=5) == []


def test_module_importable():
    assert callable(retrieve)


# ── RRF 常数 ────────────────────────────────────────────────


def test_rrf_k_is_tuned_for_short_candidate_lists(monkeypatch):
    """RRF 的 k 必须与候选池规模同量级，否则「两通道都平庸」会压过「某通道第一」。

    这是 2026-10 用真实 DeepSeek 跑端到端时发现的问题：查询
    「MBTI 外倾 内倾 维度 高分 低分 含义」，不相关的《心理咨询如何起效》
    在 vector 排 #10、BM25 排 #9，而真正该出的《MBTI维度强度解读》BM25 排 #1。
    k=60 时前者 1/69+1/68=0.0292 > 后者 1/60=0.0167 —— 排序完全反了；
    k=3 时翻转为 0.25 vs 0.148。

    单测用合成节点钉死这个数学事实，不依赖语料内容（语料会改，公式不会）。

    注意两点：
      1. ``_reciprocal_rerank_fusion`` 会**就地覆盖** ``nws.score`` 为融合分，
         所以两次调用必须各自构造新节点，不能复用同一个结果字典。
      2. 断言的是「excellent 相对 mediocre 的名次」，不是「谁是第一名」——
         列表里还有别的 rank-0 填充节点与 excellent 同分，靠 node_id 兜底排序，
         那属于实现细节，不是本用例要保的性质。
    """
    from llama_index.core.schema import NodeWithScore, TextNode

    from app.rag.retriever import _NodeIdFusionRetriever

    def build() -> dict:
        def nws(node_id: str) -> NodeWithScore:
            return NodeWithScore(node=TextNode(id_=node_id, text=node_id), score=1.0)

        # 通道 A：excellent 排 #1（0 基 rank 0），mediocre 排 #10
        # 通道 B：excellent 不在列表里，mediocre 排 #9
        return {
            ("a", 1): [nws("excellent")]
            + [nws(f"a{i}") for i in range(1, 10)]
            + [nws("mediocre")],
            ("b", 1): [nws(f"b{i}") for i in range(0, 9)] + [nws("mediocre")],
        }

    def rank_of(fused, node_id: str) -> int:
        return [n.node.node_id for n in fused].index(node_id)

    # k=60（官方默认）：mediocre 胜出——这正是要修的故障
    monkeypatch.setattr(
        "app.rag.retriever.get_settings", lambda: type("S", (), {"rrf_k": 60})()
    )
    fused60 = _NodeIdFusionRetriever._reciprocal_rerank_fusion(None, build())
    assert rank_of(fused60, "mediocre") < rank_of(fused60, "excellent"), (
        "k=60 下应当复现「双通道平庸胜出」，否则本用例的前提变了"
    )

    # 默认 k=3：excellent 胜出
    monkeypatch.setattr(
        "app.rag.retriever.get_settings", lambda: type("S", (), {"rrf_k": 3})()
    )
    fused3 = _NodeIdFusionRetriever._reciprocal_rerank_fusion(None, build())
    assert rank_of(fused3, "excellent") < rank_of(fused3, "mediocre"), (
        "k=3 下「单通道第一」必须胜出；若失败说明 rrf_k 被改大了"
    )


def test_default_rrf_k_matches_documented_value():
    """``RRF_K`` 只是测试算期望分用的符号，必须与 settings 默认值一致。"""
    from app.core.config import Settings

    assert RRF_K == Settings().rrf_k == 3
