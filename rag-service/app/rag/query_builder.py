"""测评结果 → 检索 query（§8.1）。

## 构建器必须避开的三处陷阱

| 必须这样写 | 原因 | 写错会怎样 |
|---|---|---|
| 用 ``flags`` 里的 ``needs_referral`` 判断是否检索「何时就医」 | ``SCL90.go:173`` 的 ``NeedFurtherEval`` 已经算好（``TotalScore>=160 || len(HighRiskFactors)>0 || PositiveCount>43``） | 拿 ``level != "正常"`` 判断会**恒为真**——``SCL90.go:177/186/196/206`` 只产出 ``正常范围``/``轻度异常``/``明显异常``/``严重异常``，**根本没有「正常」**，正常用户也会被追加「何时就医」检索 |
| 用 ``measures[].kind`` 分发维度/因子模板 | 契约已把两者统一成单一数组（ADR-0001） | 按 ``test.code`` 分支写模板，会让 42 个测评各写一套，加一个量表改一处 |
| 用 ``measures[].strength`` 判断强度，``factor_level`` 只喂 prompt | ``strength ∈ [-1,1]`` 是数值，可比较可测 | 字符串比较（``level in ("中度","偏重","严重")``）既脆又难测，且 **``factors[].level`` 这个字段在 ``SCL90.go:9-14`` 的 ``Factor`` 结构体里根本不存在**（只有 ``Name/Items/Score/RawScore``），会直接 KeyError |

## 覆盖维度

按主文档要求，query 必须覆盖到**每一个 measure 的 label**，而不是只挑高的写。
覆盖不足会让 ``coverage`` 指标（有 claim 覆盖的 measure 数 / 总 measure 数）低于
0.8 的红线。

## 为什么 per-measure query 里**不注入测评名**（实测，勿改回去）

``req.test.name`` 形如 ``16型人格测试（MBTI）``。把它拼进每条 query 会给每条
query 注入「人格／测评／测试」这类**泛化词**，而 ``common`` 与 ``typology`` 里
恰好存在若干「测评工具总览」型 hub 块，它们同时枚举 16 型人格 / DISC / 抑郁焦虑
躯体化等几乎所有关键词，于是对任何测评类 query 都高匹配。

实测（top_k=3，域 ``["typology","common"]``，RRF k=3）：

| query | top-1 |
|---|---|
| ``16型人格测试（MBTI） 外倾 / 内倾 维度 高分与低分分别意味着什么`` | ❌ ``[common] 心理咨询如何起效.md``（fused 0.6667 = 两通道都排 #0） |
| ``外倾 与 内倾 维度 高分与低分分别意味着什么`` | ✅ ``[typology] MBTI维度强度解读.md`` |

同一条 MTBI 查询，两条通道的 top-4 向量分数是 0.7648 / 0.7499 / 0.7441 / 0.7381
—— **极差只有 0.027**，bge-small-zh（512 维）在这个密度下分不开 hub 块与正题块；
去掉测评名是把这个「人为抬高的 hub 匹配」消掉，而不是指望模型变准。

这条改动**不违反**主文档「query 模板里不许出现具体测评名」的本意：那条规则禁止的
是硬编码 ``"SCL-90"`` 之类字面量，不是要求必须注入 ``req.test.name``。
"""

from __future__ import annotations

from app.api.schemas import EnhanceRequest

#: 每个 measure 都必须产出至少一条 query，保证 coverage（ADR-0003）不至于因
#: 「没检索到该维度的知识」而塌方。query 模板里**不得出现具体测评名**——
#: 既不许硬编码 ``"SCL-90"`` 字面量，也**不要注入** ``req.test.name``（见下）。
_TEMPLATE_DIMENSION = "{label} 维度 高分与低分分别意味着什么"
_TEMPLATE_FACTOR = "{label} 因子 临床意义 解释"
_REFERRAL_QUERY = "何时需要寻求专业心理帮助"

#: ``label`` 里常见 ``"外倾 / 内倾"`` 这种斜杠并列写法。斜杠在中文 query 里
#: 切断了语义（分词器会把两侧当成独立符号），换成「与」后向量与 BM25 都更稳。
_LABEL_SEPARATORS = ("／", "/")


def _normalize_label(label: str) -> str:
    """把 ``"外倾 / 内倾"`` 归一成 ``"外倾 与 内倾"``。"""
    out = label
    for sep in _LABEL_SEPARATORS:
        out = out.replace(f" {sep} ", " 与 ").replace(sep, " 与 ")
    return " ".join(out.split())


def build_queries(req: EnhanceRequest) -> list[str]:
    """为一次增强生成一组检索 query。

    产出规则：
      * 每个 ``measure`` 一条 query——``kind=="dimension"`` 用维度模板，
        ``kind=="factor"`` 用因子模板；
      * ``req.needs_referral`` 为真时补一条转介 query（「何时寻求专业帮助」）；
      * 若最终为空（理论上不会，``measures`` 至少 1 项），退回一条整体解读。

    Args:
        req: 已通过 Pydantic 校验的请求。

    Returns:
        非空 query 列表。
    """
    queries: list[str] = []
    test = req.test.name

    for m in req.measures:
        label = _normalize_label(m.label)
        if m.kind == "dimension":
            queries.append(_TEMPLATE_DIMENSION.format(label=label))
        else:  # factor
            queries.append(_TEMPLATE_FACTOR.format(label=label))

    if req.needs_referral:
        queries.append(_REFERRAL_QUERY)

    # 兜底 query 才用测评名：此时没有任何 measure 可依据，只能整体问一句。
    return queries or [f"{test} 结果 解读"]


def measure_is_high(measure) -> bool:
    """强度是否值得专门提一句。

    一律看数值 ``strength``；``factor_level`` 只用于 prompt 措辞。
    """
    return measure.strength is not None and abs(measure.strength) >= 0.5


def retrieval_namespaces(req: EnhanceRequest) -> list[str]:
    """检索范围由 Go 侧元数据给定（§7.2），Python 不自己判断测评类型。

    契约里 ``domains`` 至少 1 项且已通过 ``TestInfo`` 校验（未知域、重复、
    policy 与 domains 不一致都会被拒），这里直接返回即可。
    """
    return list(req.test.domains)
