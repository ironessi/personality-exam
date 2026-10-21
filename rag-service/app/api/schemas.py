"""`POST /api/v1/enhance` 的请求/响应契约。

与 doc/RAG增强系统开发文档.md §5.2（:274-338）逐字段对应。契约归属见
doc/adr/0001-增强请求契约归属-go侧.md：

  * **Go 侧**负责把 42 个异构测评结果规范化成本文件的形状；
    ``strength ∈ [-1,1]``、``factor_level``、``flags`` 全部在 Go 侧算好。
  * **Python 侧不得读取任何原始测评字段**（``raw_result`` / ``user_answers``
    之类根本不在契约里）。
  * 遇到不认识的 ``schema_version`` **必须显式报错**，不做猜测性兼容。

为什么要这么严：MBTI 的 E/S 满分 21、T/J 满分 14，SCL90 的因子分又是 1~5 的
均分——原始分之间不可比。把「解释原始分」的职责留在 Go 侧，Python 只面对
已经归一化的数值。
"""

from __future__ import annotations

from typing import Literal

from pydantic import BaseModel, ConfigDict, Field, model_validator

SUPPORTED_SCHEMA_VERSION = 1

#: 知识库分区键 = **构念域**（ADR-0008），**不是测评名**。
#: 同一份值三处共用：Chroma where 过滤、BM25 索引选择、data/raw 子目录名。
KNOWN_DOMAINS: tuple[str, ...] = (
    "symptom",  # 症状严重度
    "trait",  # 特质人格
    "typology",  # 类型学人格
    "social",  # 人际与社会
    "wellbeing",  # 福祉与积极心理
    "career",  # 职业兴趣
    "ability",  # 能力与智力
    "common",  # 通用层：任何测评都可检索
)

#: 增强策略（ADR-0008）。这不是性能开关，是**安全措施**：
#: ``common_only`` 用于受出版方版权保护、且必须由受训主试实施的测验
#: （斯坦福-比奈、罗夏）；``disabled`` 用于娱乐测评（疯狂动物城、SBTI）。
#: 对它们生成个性化深度解读等于伪造权威。
EnhancePolicy = Literal["full", "common_only", "disabled"]


class TestInfo(BaseModel):
    """测评身份**与检索范围**。

    ``code``（MBTI=13/SCL90=1/DISC=28）与 ``key``（``mbti``）都只是身份，
    **都不是分区键**——分区键是 ``domains``。

    检索范围由 Go 侧按构念域算好后放进 ``domains``（ADR-0008），
    Python **不判断测评类型**：它只照 ``domains`` 过滤 Chroma 与 BM25。
    这样新增一个测评时，本文件与检索层都不用改。
    """

    model_config = ConfigDict(extra="forbid")

    code: int = Field(..., description="consts.go:8-55 里的数字编码，共 42 个")
    key: str = Field(..., description="机读名，如 mbti。⚠️ 不是 namespace")
    name: str = Field(..., description="展示名，如 16型人格测试（MBTI）")
    domains: list[str] = Field(
        ...,
        min_length=1,
        description="检索范围，取值见 KNOWN_DOMAINS；通常是「本域 + common」",
    )
    policy: EnhancePolicy = Field(
        ..., description="full=正常增强 / common_only=只检索通用层 / disabled=不接入"
    )

    @model_validator(mode="after")
    def _check_retrieval_scope(self) -> "TestInfo":
        unknown = [d for d in self.domains if d not in KNOWN_DOMAINS]
        if unknown:
            raise ValueError(
                f"未知构念域 {unknown}；合法取值见 KNOWN_DOMAINS（ADR-0008）。"
                "不要按测评名建 namespace。"
            )
        if len(set(self.domains)) != len(self.domains):
            raise ValueError(f"domains 有重复项：{self.domains}")
        if self.policy == "common_only" and self.domains != ["common"]:
            raise ValueError(
                f"policy=common_only 的检索范围只能是 ['common']，收到 {self.domains}"
            )
        if self.policy == "full" and "common" not in self.domains:
            raise ValueError(
                f"policy=full 的检索范围必须包含 common，收到 {self.domains}"
            )
        return self


class ResultInfo(BaseModel):
    """整体结果。``total`` 的量纲随量表而异，仅用于叙述，不参与算强度。"""

    model_config = ConfigDict(extra="forbid")

    level: str
    total: float | None = None
    total_max: float | None = None


class Measure(BaseModel):
    """统一承载 dimension 与 factor，用 ``kind`` 判别（§5.2 :283-289）。

    不要拆回 ``dimensions`` / ``factors`` 两个数组：那正是 ADR-0001 取代掉的
    「按测评分支理解」的旧形状，会让每个量表各写一套模板。

    ``strength`` 是**唯一**允许控制流程的量：``[-1,1]``，Go 侧算好。
    ``factor_level``（因子分 ≥ 2.0 记「偏高」）只喂 prompt，字符串比较既脆又难测。
    """

    model_config = ConfigDict(extra="forbid")

    key: str
    label: str
    kind: Literal["dimension", "factor"]

    score: float
    max: float | None = None

    # 维度专用
    threshold: float | None = None
    pole: str | None = None

    # 因子专用
    clinical_cutoff: float | None = None
    factor_level: str | None = None

    # 两者共用：控制流程只看它
    strength: float | None = Field(
        default=None, ge=-1.0, le=1.0, description="归一化强度，[-1,1]；Go 侧计算"
    )

    @model_validator(mode="after")
    def _check_kind_specific_fields(self) -> "Measure":
        if self.kind == "dimension":
            if self.threshold is None:
                raise ValueError(f"dimension {self.key!r} 缺少 threshold（Go 侧必须给）")
            if self.pole is None:
                raise ValueError(f"dimension {self.key!r} 缺少 pole（如 'E' / 'I'）")
        else:  # factor
            if self.clinical_cutoff is None:
                raise ValueError(
                    f"factor {self.key!r} 缺少 clinical_cutoff（如 2.0），"
                    "否则无从判断「偏高」"
                )
        return self


class EnhanceOptions(BaseModel):
    model_config = ConfigDict(extra="forbid")

    top_k: int = Field(default=5, ge=1, le=50)
    max_suggestions: int = Field(default=5, ge=1, le=20)
    language: str = "zh-CN"
    include_citations: bool = True


class EnhanceRequest(BaseModel):
    model_config = ConfigDict(extra="forbid")

    schema_version: int
    test: TestInfo
    result: ResultInfo
    measures: list[Measure] = Field(..., min_length=1)
    labels: dict[str, str] = Field(default_factory=dict)
    flags: list[str] = Field(default_factory=list, description="如 needs_referral")
    options: EnhanceOptions = Field(default_factory=EnhanceOptions)

    @model_validator(mode="after")
    def _check_schema_version(self) -> "EnhanceRequest":
        if self.schema_version != SUPPORTED_SCHEMA_VERSION:
            raise ValueError(
                f"不支持的 schema_version={self.schema_version}，"
                f"本服务只认 {SUPPORTED_SCHEMA_VERSION}。"
                "契约变更时必须显式升级版本号，不要静默兼容。"
            )
        return self

    @property
    def needs_referral(self) -> bool:
        return "needs_referral" in self.flags


# ── 响应 ────────────────────────────────────────────────────


class Section(BaseModel):
    model_config = ConfigDict(extra="forbid")

    title: str
    content: str


class Suggestion(BaseModel):
    model_config = ConfigDict(extra="forbid")

    text: str
    basis: str = Field(..., description="如 MBTI_16型人格详解#p12")


class Citation(BaseModel):
    """``evidence_quote`` 必须是所引源块原文的 ≤60 字**逐字复制**。

    这是 ADR-0003 Stage 1 的机械校验依据：规范化空白与全/半角标点后做字面子串
    比较。用「最长公共子串占比 ≥ 阈值」验引用是错的——它只能证明「抄了原文」，
    不能证明「原文支持这个结论」。
    """

    model_config = ConfigDict(extra="forbid")

    claim: str
    source: str
    chunk_id: str
    score: float | None = None
    evidence_quote: str | None = Field(default=None, max_length=64)
    verified: bool = False
    #: 这条引用背书的是 ``sections`` 里的第几节（0 起）。只用于算 ``coverage``
    #: （ADR-0003 要求 precision / coverage 成对上报），不参与展示。
    #: 模型漏填时为 ``None``——**不要当成第 0 节**，否则 coverage 会被系统性高估。
    section: int | None = Field(default=None, ge=0)


RagStatus = Literal["pending", "ready", "degraded", "failed"]


class EnhanceData(BaseModel):
    """⚠️ 只有这里的**业务字段**允许落库。

    ``model`` / ``prompt_tokens`` / ``completion_tokens`` / ``latency_ms`` 仅供
    服务端日志与成本统计；Go 侧写 ``resultMap["rag"]`` 时必须显式构造白名单结构体，
    **不得整体透传**（§5.2 :330）。
    """

    model_config = ConfigDict(extra="forbid")

    enhanced_text: str
    sections: list[Section] = Field(default_factory=list)
    suggestions: list[Suggestion] = Field(default_factory=list)
    citations: list[Citation] = Field(default_factory=list)

    model: str | None = None
    prompt_tokens: int | None = None
    completion_tokens: int | None = None
    latency_ms: int | None = None

    rag_status: RagStatus = "ready"


class EnhanceResponse(BaseModel):
    model_config = ConfigDict(extra="forbid")

    success: bool
    data: EnhanceData | None = None
    error: str | None = None
    disclaimer: str = (
        "本结果基于心理学公开文献生成，仅供参考，不构成医学诊断。"
        "如有持续困扰，请咨询专业心理医生。"
    )


class HealthResponse(BaseModel):
    model_config = ConfigDict(extra="forbid")

    status: Literal["ok", "degraded"]
    embedding_loaded: bool
    collection_count: int
    llm_reachable: bool
    #: Stage 2（NLI 语义蕴含）是否已加载。**不参与 status 判定**：
    #: 缺 NLI 权重时服务仍能跑（退化为只验字面），但这是「两级验证」承诺的
    #: 实际降级，运维必须看得见——所以单列一个字段，而不是悄悄降 status
    #: 让 Go 侧以为整个增强都挂了。
    stage2_ready: bool = False
