"""敏感词过滤器（ADR-0004）。

## 契约

导出 ``is_claim_safe(claim: str, context_window: int = 8) -> bool``。
集成点是 ``app/rag/synthesizer.py`` 生成每条 claim 之后立即调用，失败的 claim
丢弃并重试或跳过。**不作用于检索原文**——只审生成出来的 claim。

## 为什么不能用「单词黑名单、命中即拦」

``确诊|你患有|你得了|建议服用|建议停药|抗抑郁药|精神分裂|你必须`` 这类写法有
三处真实误杀（§10.2 :827）：

| 词 | 会被误杀的合规句 |
|---|---|
| ``确诊`` | 「SCL-90 **不能用于确诊**」 |
| ``精神分裂`` | 因子 9「精神病性」的科普解释 |
| ``抗抑郁药`` | 「请勿自行调整**抗抑郁药**，需遵医嘱」 |

正确做法是**断言句式黑名单 + 否定语境豁免**：拦的是「把结论安到用户头上」的
句式（你患有 X / 你的 X 说明你有 Y），而不是词本身；句子里出现否定、劝诫、
免责语境时豁免。
"""

from __future__ import annotations

import re

# 否定/劝诫/免责语境标记。命中断言句式后，向左右各看 context_window 个字符，
# 出现任一标记即豁免——「SCL-90 不能用于确诊」里「确诊」前有「不能」。
NEGATION_MARKERS: tuple[str, ...] = (
    "不能",
    "不可",
    "不构成",
    "不足以",
    "无法",
    "难以",
    "不代表",
    "并非",
    "不是",
    "请勿",
    "不要",
    "禁止",
    "避免",
    "仅供参考",
    "建议咨询",
    "需遵医嘱",
    "请就医",
)

# 断言句式：把诊断结论安在用户头上。这是唯一要拦的东西——
# 不是「出现某词」，而是「用断言句式把结论说成用户的事实」。
ASSERTION_PATTERNS: tuple[str, ...] = (
    "你患有",
    "你得了",
    "你已经确诊",
    "你有明显的",
    "你属于重度",
    "你必须",
    "建议服用",
    "建议停药",
    "自行调整用药",
)

# 「建议服用/停用/加量/减量」这类直接给药指令，句式比上面的「你患有」更危险，
# 且否定语境下也不该出现（「不建议服用」仍是给药建议）。单独一个集合，命中即拦、
# 不做否定豁免——因为它本身就是动作指令，而不是对用户状态的断言。
MEDICATION_DIRECTIVES: tuple[str, ...] = (
    "建议服用",
    "建议停用",
    "建议停药",
    "建议加量",
    "建议减量",
    "自行调整用药",
    "自行停药",
    "自行加药",
    "自行减药",
)

# 「抗抑郁药」这类具体药名单独出现不拦（科普语境合法），但被上述
# 指令句式包裹时必须拦。为简洁，把「给药动作」统一进指令集合，
# 这里不再单列药名。

# 预编译
_ASSERTION_RE = re.compile("|".join(re.escape(p) for p in ASSERTION_PATTERNS))
_MEDICATION_RE = re.compile("|".join(re.escape(p) for p in MEDICATION_DIRECTIVES))


def is_claim_safe(claim: str, context_window: int = 8) -> bool:
    """claim 是否可以输出。

    ``context_window``：命中断言句式时，向左右各看这么多字符判断是否处于否定
    /劝诫语境；处于否定语境则豁免（例如「SCL-90 不能用于确诊」）。

    规则（ADR-0004）：
    1. 直接给药/调药指令（``MEDICATION_DIRECTIVES``）——命中即拦，**不做**否定豁免。
       理由：这类句式本身就是要用户去服药/停药/改药，无论语气是否否定都不该由
       一个测评报告产出；「不建议服用」依然是给药建议。
    2. 断言句式（``ASSERTION_PATTERNS``）——命中后回看上下文，出现否定/劝诫/
       免责标记则放行；否则拦截。这是「把诊断结论安到用户头上」的判定。
    """
    if not claim or not claim.strip():
        return True  # 空 claim 不该出现，但也不构成危险；交由上层丢弃

    # 1. 给药指令：命中即拦，不看语境
    if _MEDICATION_RE.search(claim):
        return False

    # 2. 断言句式：命中后看否定语境
    for m in _ASSERTION_RE.finditer(claim):
        start, end = m.span()
        # 向左看 context_window 个字符
        left = claim[max(0, start - context_window):start]
        # 向右看 context_window 个字符（覆盖「不能用于确诊」这类后置否定）
        right = claim[end:end + context_window]
        window = left + right
        if not any(marker in window for marker in NEGATION_MARKERS):
            return False

    return True
