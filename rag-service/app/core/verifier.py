"""两级引用验证（ADR-0003）。

## 不要用相似度阈值验引用

``char_overlap_ratio(claim, chunk) >= 0.35``（最长公共子串占比）这类写法是错的：
它只能证明「**抄了原文**」，不能证明「**原文支持这个结论**」。

- 抄对句子挂错结论 → 被判 ``verified=true``（漏放）
- 正确转述 → 被判死剔除（误杀）

## 正确的两级

**Stage 1 机械验证**（零成本、确定性、可离线）
校验 ``evidence_quote`` 是否为所引 chunk 原文的**字面子串**——规范化空白与
全/半角标点后比较。生成侧被强制为每条 claim 附 ``evidence_quote``（≤60 字逐字
复制）+ ``chunk_id``。这层能干掉编造引文、拼凑引文与 chunk_id 幻觉。

**Stage 2 语义蕴含**
本地 NLI 交叉编码器判断 ``evidence_quote ⊨ claim``。用
``IDEA-CCNL/Erlangshen-Roberta-110M-NLI``（中文 110M）——
知识库语料与用户输入都是中文，英文主训的多语模型没有优势却更占内存与 CPU；
``MoritzLaurer/mDeBERTa-v3-base-mnli-xnli``（多语 279M）留作备选。
判定阈值先在 **50 条人工标注样本**上对齐。

## 兜底不允许「去掉标记照发」

验证率低时只能返回 ``rag_status="degraded"`` 并**剔除**未通过的 claim，
不允许「只返回 ``enhanced_text`` 不带 citations」——那是把无据结论摘掉标记
照样发给用户。

## 质量红线：两个指标必须成对上报

| 指标 | 定义 | 红线 |
|---|---|---|
| ``support_precision`` | 通过验证的 claim / 输出 claim 总数 | ≥ 0.95，且报 Wilson 95% CI |
| ``coverage`` | 有 claim 覆盖的 measure 数 / 总 measure 数 | ≥ 0.8 |

只报 ``support_precision`` 时，「每篇只输出一条最保险结论」即可刷到 100%。
"""

from __future__ import annotations

import math
import re
import unicodedata
from pathlib import Path
from dataclasses import dataclass

#: Stage 2 语义蕴含的 NLI 模型（ADR-0003 已定）。
#: 首次使用时按需加载，不在 import 期加载（torch/transformers 太重，会拖慢启动）。
_NLI_MODEL_NAME = "IDEA-CCNL/Erlangshen-Roberta-110M-NLI"
#: 语义蕴含判定阈值。ADR-0003 明确：先在 50 条人工标注样本上对齐后再定最终值；
#: 这里给一个保守默认（entailment 概率 ≥ 0.5 判蕴含），标注完成后用实测值覆盖。
_NLI_ENTAILMENT_THRESHOLD = 0.5

#: 半角 → 全角标点映射（只覆盖中文文本里最常见的几对，避免过度规范化）。
_HALF_TO_FULL = {
    ",": "，",
    ";": "；",
    ":": "：",
    "!": "！",
    "?": "？",
    "(": "（",
    ")": "）",
    ".": "。",
}

#: 零宽字符与不可见字符，折叠空白时一并去掉。
_ZERO_WIDTH = {
    "\u200b",  # zero width space
    "\u200c",  # zero width non-joiner
    "\u200d",  # zero width joiner
    "\ufeff",  # BOM / zero width no-break space
    "\u00a0",  # no-break space
}

_WHITESPACE_RE = re.compile(r"\s+")


@dataclass(frozen=True)
class VerificationResult:
    claim: str
    chunk_id: str
    passed_stage1: bool
    passed_stage2: bool | None  # None = Stage 2 未启用或未跑
    reason: str

    @property
    def verified(self) -> bool:
        return self.passed_stage1 and self.passed_stage2 is not False


def normalize_for_substring(text: str) -> str:
    """Stage 1 的比较基准：折叠空白、统一全/半角标点、去零宽字符。

    两侧都必须过这个函数后再做 ``in`` 判断，否则「原文里的换行/全角逗号」
    会造成假阴性。
    """
    if not text:
        return ""

    # 1. 去零宽/不可见字符
    out = "".join(ch for ch in text if ch not in _ZERO_WIDTH)

    # 2. NFKC 归一化（把全角 ASCII 字母数字、组合字符规范化，但保留中文标点语义）
    out = unicodedata.normalize("NFKC", out)

    # 3. 半角标点 → 全角，统一标点形态
    out = "".join(_HALF_TO_FULL.get(ch, ch) for ch in out)

    # 4. 折叠所有连续空白为单个空格，再整体去空格（字面子串比较不关心空格）
    out = _WHITESPACE_RE.sub("", out)

    return out


def stage1_substring_ok(evidence_quote: str, chunk_text: str) -> bool:
    """``evidence_quote`` 是否为 ``chunk_text`` 的规范化字面子串。"""
    q = normalize_for_substring(evidence_quote)
    c = normalize_for_substring(chunk_text)
    if not q:
        return False  # 空引文不能算通过
    return q in c


# ── Stage 2：NLI 语义蕴含（懒加载）──────────────────────────


def _resolve_model_ref() -> str:
    """决定 NLI 权重从哪加载：本地目录优先，其次 HuggingFace 模型名。

    相对路径按 ``rag-service/`` 解析（与 ``.env`` 同一基准），这样
    ``NLI_MODEL_PATH=.models/Erlangshen`` 在任意工作目录下都成立。
    """
    try:
        from app.core.config import SERVICE_ROOT, get_settings

        s = get_settings()
        raw = (s.nli_model_path or "").strip()
        if not raw:
            return s.nli_model_name or _NLI_MODEL_NAME
        p = Path(raw)
        if not p.is_absolute():
            p = SERVICE_ROOT / p
        if p.exists():
            return str(p)
        # 配了路径却不存在：不要静默回退到联网下载（容器里那会卡住请求）
        raise FileNotFoundError(f"nli_model_path 不存在：{p}")
    except FileNotFoundError:
        raise
    except Exception:  # noqa: BLE001 —— 配置读不到就退回内置模型名
        return _NLI_MODEL_NAME


class _NliModel:
    """进程内单例的 NLI 交叉编码器。

    懒加载：``transformers`` / ``torch`` 的 import 与权重下载都发生在首次
    ``entails`` 调用，不在模块 import 期。加载失败不抛穿——返回 ``None`` 让
    调用方把 Stage 2 视为「未启用」（``passed_stage2=None``），整条链路退化为
    Stage 1 only，而不是让 /enhance 崩溃。
    """

    _instance = None

    @classmethod
    def get(cls) -> "_NliModel | None":
        """返回已加载好的单例；首次调用会**真的去加载**权重。

        ⚠️ 这里必须调 ``_ensure_loaded()``：``__init__`` 只做初始化、不加载，
        若 ``get()`` 只看 ``_available`` 就会永远返回 ``None``——于是
        ``entails()`` 一次都不会被调用，Stage 2 静默永久失效，而且**日志里
        一句解释都没有**（踩过：模型明明装好了，Stage2 仍然 0/N）。
        """
        if cls._instance is None:
            cls._instance = _NliModel()
        inst = cls._instance
        return inst if inst._ensure_loaded() else None

    @classmethod
    def get_error(cls) -> str | None:
        """加载失败原因（未尝试加载时为 ``None``）。"""
        return getattr(cls._instance, "_load_error", None)

    def __init__(self) -> None:
        self._available = False
        self._pipeline = None
        self._load_error: str | None = None

    def _ensure_loaded(self) -> bool:
        if self._available:
            return True
        if self._load_error is not None:
            return False
        try:
            from transformers import pipeline

            self._pipeline = pipeline(
                "text-classification",
                model=_resolve_model_ref(),
            )
            self._available = True
        except Exception as exc:  # noqa: BLE001 —— 权重缺失/网络失败都不该拖垮服务
            self._load_error = str(exc)
        return self._available

    def entails(self, premise: str, hypothesis: str) -> bool | None:
        """``premise ⊨ hypothesis`` 是否成立。

        返回 ``None`` 表示模型不可用（Stage 2 未启用）。
        """
        if not self._ensure_loaded():
            return None
        try:
            out = self._pipeline(
                {"text": premise, "text_pair": hypothesis},
                top_k=1,
            )
            # ⚠️ 返回形状随版本不同：``top_k=1`` 时新版是**扁平列表**
            # ``[{"label": "ENTAILMENT", "score": 0.94}]``，老版是嵌套
            # ``[[{...}]]``。按 ``out[0][0]`` 取值在新版上会 KeyError/TypeError，
            # 而异常被下面吞掉 → ``entails`` 永远返回 None → Stage 2 静默失效。
            item = out
            while isinstance(item, list) and item and isinstance(item[0], list):
                item = item[0]
            if isinstance(item, list):
                if not item:
                    return None
                top = item[0]
            else:
                top = item
            # 该模型 id2label 是大写：CONTRADICTION / NEUTRAL / ENTAILMENT
            if str(top.get("label", "")).lower() != "entailment":
                return False
            return float(top.get("score", 0.0)) >= _NLI_ENTAILMENT_THRESHOLD
        except Exception:  # noqa: BLE001
            return None


def _nli_entails(premise: str, hypothesis: str) -> bool | None:
    model = _NliModel.get()
    if model is None:
        return None
    return model.entails(premise, hypothesis)


def stage2_status() -> tuple[bool, str | None]:
    """Stage 2（NLI 语义蕴含）是否真的可用，以及不可用的原因。

    为什么需要这个函数：``entails`` 在模型加载失败时返回 ``None``，而
    ``verified`` 的定义是「Stage 1 过且 Stage 2 不是 False」——也就是说
    **模型缺失时两级验证会静默退化成一字面级**，``verified=True`` 照常写回。
    这在生产里是危险的：没人会注意到防幻觉的第二道闸门根本没合上。

    实测触发场景（2026-10）：NLI 权重未下载且 DSH 沙箱只允许写工作区，
    HF 报 ``PermissionError ... Check cache directory permissions``，
    于是本进程内 Stage 2 全程 ``None``。

    故调用方（``routes._enhance_async``）应在首次发现不可用时**大声记一条
    WARNING**，让运维知道当前跑的是单级验证。
    """
    model = _NliModel.get()
    if model is None:
        # ``_load_error`` 是**实例**属性，不在未实例化时存在——必须从单例上读，
        # 否则这里会抛 AttributeError，把「Stage 2 不可用」这条告警本身弄丢。
        return False, _NliModel.get_error() or "NLI 模型未能加载"
    return True, None


# ── 主入口 ──────────────────────────────────────────────────


def verify_citations(
    citations: list,
    chunks_by_id: dict,
    *,
    run_stage2: bool = True,
) -> tuple[list, list[VerificationResult]]:
    """跑两级验证。

    Args:
        citations: ``Citation`` 列表（Pydantic 模型）。
        chunks_by_id: ``chunk_id -> chunk_text`` 映射，用于 Stage 1 字面校验。
        run_stage2: 是否跑 NLI 语义蕴含。评测脚本可能只想要 Stage 1。

    Returns:
        ``(通过的 citations, 每条 claim 的判定明细)``。
        调用方据明细算 ``support_precision`` / ``coverage`` 并**成对**上报。
    """
    passed: list = []
    details: list[VerificationResult] = []

    for c in citations:
        chunk_text = chunks_by_id.get(c.chunk_id, "")
        if not chunk_text:
            details.append(
                VerificationResult(
                    claim=c.claim,
                    chunk_id=c.chunk_id,
                    passed_stage1=False,
                    passed_stage2=None,
                    reason="chunk_id 不存在（幻觉或已删）",
                )
            )
            continue

        ok1 = stage1_substring_ok(c.evidence_quote or "", chunk_text)
        if not ok1:
            details.append(
                VerificationResult(
                    claim=c.claim,
                    chunk_id=c.chunk_id,
                    passed_stage1=False,
                    passed_stage2=None,
                    reason="evidence_quote 不是所引块原文的字面子串",
                )
            )
            continue

        ok2: bool | None = None
        if run_stage2:
            ok2 = _nli_entails(premise=c.evidence_quote or "", hypothesis=c.claim)

        verified = ok2 is not False  # Stage 1 已过；Stage 2 未启用(None)视为通过
        if verified:
            c.verified = True
            passed.append(c)
        details.append(
            VerificationResult(
                claim=c.claim,
                chunk_id=c.chunk_id,
                passed_stage1=True,
                passed_stage2=ok2,
                reason=(
                    "通过" if verified
                    else "evidence_quote 不支持该 claim（语义蕴含不成立）"
                ),
            )
        )

    return passed, details


def support_precision(passed: int, total: int) -> tuple[float, float, float]:
    """返回 ``(precision, wilson_low, wilson_high)``。

    样本量小时单点估计没有意义——50 条里通过 48 条，Wilson 95% CI 下界约 0.86，
    远低于 0.95 的红线。
    """
    if total <= 0:
        return 0.0, 0.0, 0.0

    p = passed / total
    n = total
    z = 1.959963984540054  # 95% 置信区间对应的 z 值

    denom = 1 + z * z / n
    centre = (p + z * z / (2 * n)) / denom
    half = (z * math.sqrt((p * (1 - p) + z * z / (4 * n)) / n)) / denom
    return p, max(0.0, centre - half), min(1.0, centre + half)


def support_coverage(verified_citations, total_sections: int) -> float:
    """覆盖率：**含 ≥1 条通过验证 claim 的 section 占比**（ADR-0003）。

    为什么必须和 ``support_precision`` 成对上报：只报精确率时，一个"每节只写
    一句最保险的话"的系统能轻松拿到 100%，而报告同时变得毫无信息量。覆盖率
    把这个退化解法打回原形。

    没有 ``section`` 标记的引用**不计入覆盖**（``_clean_citation`` 会把越界值
    置 None）——宁可少算，不可把"模型漏填"当成"覆盖到了第 0 节"。
    """
    if total_sections <= 0:
        return 0.0
    covered = {
        c.section
        for c in verified_citations
        if getattr(c, "section", None) is not None and 0 <= c.section < total_sections
    }
    return len(covered) / total_sections
