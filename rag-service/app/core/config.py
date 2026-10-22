"""应用配置。

所有字段都能被环境变量或 rag-service/.env 覆盖，字段名大写即为环境变量名
（``embedding_model_path`` ← ``EMBEDDING_MODEL_PATH``）。
"""

from __future__ import annotations

from functools import lru_cache
from pathlib import Path
from typing import Literal

from pydantic import Field
from pydantic_settings import BaseSettings, SettingsConfigDict

# rag-service/ —— app/core/config.py 往上数三层
SERVICE_ROOT: Path = Path(__file__).resolve().parents[2]

# BAAI/bge-small-zh-v1.5 在开发机的 HF 快照目录（已实测存在，512 维）
DEFAULT_EMBEDDING_MODEL_PATH = (
    "/Users/iron/.cache/huggingface/hub/"
    "models--BAAI--bge-small-zh-v1.5/snapshots/"
    "7999e1d3359715c523056ef9478215996d62a620"
)


class Settings(BaseSettings):
    """进程级配置。字段顺序按「服务 → 模型 → 向量库 → LLM → 检索」分组。"""

    model_config = SettingsConfigDict(
        env_file=SERVICE_ROOT / ".env",
        env_file_encoding="utf-8",
        case_sensitive=False,
        extra="ignore",
    )

    # ── 服务 ────────────────────────────────────────────────
    host: str = "127.0.0.1"
    port: int = 8001
    log_level: str = "INFO"
    # 留空 = 不校验调用方身份。仅当只监听 127.0.0.1 时才可留空。
    service_api_key: str = ""

    # ── 本地向量模型 ────────────────────────────────────────
    embedding_model_path: str = DEFAULT_EMBEDDING_MODEL_PATH
    # 断言用：模型真实维度不符则启动失败，避免静默产生错维向量
    embedding_dim: int = 512

    # ── 引用验证 Stage 2（NLI 语义蕴含，ADR-0003） ──────────
    # 优先用本地目录：容器/离线环境不该在请求链路上拉 400MB 权重。
    # 相对路径按 rag-service/ 解析；留空则回退到 HF 模型名（联网下载）。
    # ⚠️ 权重缺失时 Stage 2 会静默退化为 Stage 1（见 routes 的 warning），
    # 生产必须显式检查 /health 或启动日志。
    nli_model_path: str = ""
    nli_model_name: str = "IDEA-CCNL/Erlangshen-Roberta-110M-NLI"

    # ── Chroma ──────────────────────────────────────────────
    chroma_path: Path = SERVICE_ROOT / "data" / "chroma"
    collection_name: str = "psychology_knowledge"

    # ── DeepSeek ────────────────────────────────────────────
    llm_base_url: str = "https://api.deepseek.com/v1"
    llm_api_key: str = ""
    llm_model: str = "deepseek-flash"
    # 只作用于后台 worker 的调用；提交接口的 3 秒期限在 Go 侧（ADR-0007）
    llm_timeout_seconds: float = 30.0
    # 实测：该模型带思考时，2048 会被 reasoning_tokens 吃光而 content 为空。
    # 关思考（默认）后一次完整解读约用 600–900 tokens，4096 留足余量；
    # 万一有人改回 thinking=auto，这个预算也不至于立刻饿死。
    llm_max_tokens: int = 4096
    llm_temperature: float = 0.3
    # ⚠️ 默认关闭思考。DeepSeek V4 系列默认带 reasoning，而本任务是
    # 「在给定材料内受约束地提取与组织」，不是开放推理；思考会吃光输出
    # token 预算（导致 content 为空、JSON 解析失败）并挤占 30s 超时。
    # 取值：disabled（发 thinking={"type":"disabled"}）/ auto（用模型默认）。
    # 实测 auto 下 `finish_reason=length` 且 content=""，已由 synthesizer 显式报错兜底。
    llm_thinking: Literal["disabled", "auto"] = "disabled"

    # ── 检索 ────────────────────────────────────────────────
    # 知识库分区键 = **构念域**（ADR-0008），不是测评名。
    # 同一份值三处共用：Chroma where 过滤 / BM25 索引选择 / data/raw 子目录名。
    # 按测评建 namespace 会得到 42 个几乎同质的索引；构念域是封闭集合（7 个 + common），
    # 新增测评只需登记它落在哪些域，不新增目录、不改检索代码。
    kb_namespaces: list[str] = Field(
        default_factory=lambda: [
            "symptom",  # 症状严重度
            "trait",  # 特质人格
            "typology",  # 类型学人格
            "social",  # 人际与社会
            "wellbeing",  # 福祉与积极心理
            "career",  # 职业兴趣
            "ability",  # 能力与智力
            "common",  # 通用层：任何测评都可检索
        ]
    )
    top_k: int = 5
    bm25_top_k: int = 5
    # RRF 融合常数。**默认 3，不是教科书的 60。**
    #
    # Cormack 2009 的 k=60 是给上千条候选设计的：那种规模下 rank 的贡献是
    # 1/60 ~ 1/1000，几乎只用来打破平局。而本服务每条通道只取 top_k*3
    # （生产 top_k=5 → 15 条）候选，用 60 会让 rank 完全失去区分度，于是
    # 「两个通道都排第 9/10」压过「某个通道排第 1」。
    #
    # 实测（30 条 fixture，bm25=llamaindex，hybrid 通道）：
    #   k=60 → MRR 0.716 · k=10 → 0.722 · k=8 → 0.726 · k=5 → 0.729
    #   k=4  → 0.736 · k=3 → 0.769 · k=2/1 → 0.769（平台期）
    # 关键词命中率全程 87% 不变，越域命中全程 0。
    #
    # 具体故障样例：查询「MBTI 外倾 内倾 维度 高分 低分 含义」，
    # 不相关的《心理咨询如何起效》在 vector 排 #10、BM25 排 #9，
    # k=60 时融合分 1/69 + 1/68 = 0.0292 高居第一，而真正该出的
    # 《MBTI维度强度解读》BM25 排 #1 只有 1/60 = 0.0167。
    # k=3 时翻转为 0.25 vs 0.148。
    #
    # 调这个值的唯一判据是 scripts/eval_retrieval.py 的 hybrid MRR；
    # 不要凭直觉调大——那会把「两通道都平庸」重新抬回第一。
    rrf_k: int = Field(default=3, ge=1, le=100)
    # BM25 通道用哪套实现（ADR-0009）：
    #   "llamaindex" —— 官方 BM25Retriever（底层 bm25s），分词靠正则换成汉字级
    #   "jieba"      —— jieba 词级分词 + rank_bm25，BM25 这层是自写的
    #
    # 默认 llamaindex：30 条 fixture 实测 top-5 命中 85% / MRR 0.772，
    # 优于 jieba 的 75% / 0.700。
    #
    # ⚠️ 已知弱点，别当作没有：换成 llamaindex 后 **top-1 落在主构念域的占比从 80% 降到 70%**。
    # 汉字级分词让 common 域的通用段落更容易拿高分（字符覆盖广），于是命中率上去、
    # 命中内容的专属性下来。抽查反例：`16型人格 INTJ 的认知功能与类型动力学 解释`
    # 在 ["typology","common"] 下 top-5 全落在 common。fixture 扩到 100 条以上要复测；
    # 若报告出现「引用很贴题但很空泛」，先怀疑这里。
    #
    # LlamaIndex 的 BM25Retriever 没有 tokenizer 参数，默认正则 r"(?u)\b\w\w+\b"
    # 会把整句中文切成一个 token，故必须显式换 pattern（见 rag/bm25.py 的 CJK_TOKEN_PATTERN）。
    bm25_backend: Literal["jieba", "llamaindex"] = "llamaindex"

    # ── 启动预热 ────────────────────────────────────────────
    # 必须保持 True：首次 encode 有 3~4s 开销（见 rag/embedding.py 顶部说明）
    warmup_on_startup: bool = True

    # 按 token 长度暖机的区间（含端点）。每个长度的**首次** encode 要 200~2400ms，
    # 30 条 fixture 顺序跑会把检索 P95 从 35ms 抬到 837ms（详见 embedding.py 的实测记录）。
    # 代价是启动多花十几秒。设成 lo > hi 可关闭（不推荐，只在明确知道请求长度分布时用）。
    warmup_length_lo: int = Field(default=6, ge=0, le=512)
    warmup_length_hi: int = Field(default=30, ge=0, le=512)

    # ── 派生 ────────────────────────────────────────────────
    @property
    def data_raw_dir(self) -> Path:
        return SERVICE_ROOT / "data" / "raw"

    @property
    def data_processed_dir(self) -> Path:
        return SERVICE_ROOT / "data" / "processed"

    @property
    def llm_enabled(self) -> bool:
        """没配 Key 就别去调 LLM——/health 的 llm_reachable 会如实报 false。"""
        return bool(self.llm_api_key.strip())

    def ensure_dirs(self) -> None:
        """P0-1 验收项「Chroma 持久化目录可写」的前置。"""
        self.chroma_path.mkdir(parents=True, exist_ok=True)
        self.data_processed_dir.mkdir(parents=True, exist_ok=True)
        for ns in self.kb_namespaces:
            (self.data_raw_dir / ns).mkdir(parents=True, exist_ok=True)


@lru_cache(maxsize=1)
def get_settings() -> Settings:
    """进程内单例。改 .env 需重启（或清缓存）才生效。"""
    return Settings()
