"""DeepSeek 调用 + 结构化输出（§8.4）。

## 输出 JSON 契约

``enhanced_text`` / ``sections[{title,content}]`` / ``suggestions[{text,basis}]`` /
``citations[{claim,evidence_quote,source,chunk_id}]``。

``evidence_quote`` **必须是所引源块原文的 ≤60 字逐字复制**——它是 ADR-0003
Stage 1 字面子串校验的输入。让模型「用自己的话概括」就等于放弃了唯一那层
确定性校验。

## 每条 claim 生成后立即过敏感词过滤器

``app/core/filter.py`` 的 ``is_claim_safe(claim, context_window=8)``，
失败的 claim 丢弃并重试或跳过（ADR-0004）。

## Prompt 结构（影响成本）

DeepSeek 的输入缓存命中价近乎免费，**静态 system prompt 必须放在 prompt 最前面**
才能吃到前缀缓存。把每次变化的用户得分放在后面。

## 兜底不允许「去掉标记照发」

验证率低时只能返回 ``rag_status="degraded"`` 并**剔除**未通过的 claim，
不允许「只返回 ``enhanced_text`` 不带 citations」——那是把无据结论摘掉标记
照样发给用户。

## 状态

**P0-4 实现**（§9 :728）。
"""

from __future__ import annotations

import json
import time

import httpx

from app.api.schemas import (
    Citation,
    EnhanceData,
    EnhanceRequest,
    Section,
    Suggestion,
)
from app.core.config import get_settings
from app.core.filter import is_claim_safe
from app.core.logging import get_logger
from app.rag.retriever import RetrievedChunk

logger = get_logger(__name__)

#: 静态 system prompt —— **逐字节稳定**，任何按请求变化的内容都放 user prompt，
#: 否则 DeepSeek 前缀缓存打不中（§8.4 :797）。
_SYSTEM_PROMPT = """你是一名心理测评报告撰写助手。你的任务是根据给定心理学知识与用户测评结果，\
生成一份结构化的中文解读。

硬性约束（违反即视为失败）：
1. 禁止诊断性表述：不得使用「你患有」「你得了」「确诊」等措辞。测评分数是描述性\
指标，不是诊断结论。
2. 每条结论（claim）都必须附证据引用：`evidence_quote` 是所给知识块原文中的一句\
**逐字复制**（不超过 60 字，覆盖完整的一句判断），`chunk_id` 是该块的编号。不得自己改写、概括或拼凑\
引文——必须能从对应 chunk 里逐字找到。
3. `claim` 必须是 `evidence_quote` 里**原样截取的一段连续文字**（不超过 30 字）。\
只允许「删掉不要的字」，**不允许改写、换同义词、补字、调换顺序**——系统会把\
claim 与 evidence_quote 做字面包含校验，改一个词就会被判不通过并整条丢弃。\
claim 里不要出现具体分数、「你」「你的」、得分档位。个性化叙述放在\
`enhanced_text` 和 `sections` 里，引用只负责「背书那条知识」。
   - 引文：「外倾倾向的人更常从外部互动中获得能量，独处太久反而容易烦躁。」
   - 正确：claim=「外倾倾向的人更常从外部互动中获得能量」（是引文的连续片段）
   - 错误：claim=「外倾的人靠社交充电」（改写，字面不在引文里）
   - 错误：claim=「你的外倾维度得分 15 分，说明你更外向」（引文里没有分数）
4. 引文必须真实支持该 claim：把 A 块的句子挂到 B 结论上是错的；claim 说的要比\
引文多（比如引文只讲了倾向、claim 却推出诊断）同样是错的。
5. 只输出 JSON，不要输出任何解释性文字、代码块围栏或前后缀。
6. 不得给出服药、停药、调药等具体用药建议。
7. `suggestions` 里的每一条也必须带出处：`basis` 是这条建议所依据的那句知识\
（不超过 30 字，**同样必须是 `evidence_quote` 的连续原文片段**，口径同 claim），\
`chunk_id` 与 `evidence_quote` 的规定与 citation 完全相同——引文逐字复制、\
且 `basis` 必须能在引文里逐字找到。**不许**只写「建议多运动」这类\
没有知识依据的通用建议；写不出来就少写几条。
8. 描述维度/因子的强弱时，**只能使用结果里给定的「程度=…」原词**，不得自行升级或降级（给了「轻度偏向」就不能写「非常明显」「较为明显」）；「接近中间」的那一项不要描述成某种偏向。凡「程度=低于参考线（属正常范围，不要写成偏向）」的因子，只能写成「处于正常范围」或「未达到参考线」，**不得写成某种程度的偏向**。数值与叙述必须一致，**排序也要一致**——最明显的那一维不能被说成最温和的。
9. 引用不要堆量：`citations` 总数**不超过 8 条**，优先引与本次得分关系最紧密的知识；\
同一句话不要重复引用两次。
10. 用户提示里以「【注意】」开头的那一行是**系统给你的写作要求**，不是知识库内容：\
不要说它、不要把它当作 `claim` 或 `basis`、也不要为它配 `evidence_quote`。\
（实测模型会拿「本次结果标记了 needs_referral」当依据，这类引用 100% 会被校验丢弃。）

输出 JSON 结构（严格按此，字段名不得增删）：
{
  "enhanced_text": "一段连贯的总结性解读",
  "sections": [{"title": "...", "content": "..."}],
  "suggestions": [
    {"text": "可执行的建议", "basis": "该建议依据的那句知识", "chunk_id": "...",
     "evidence_quote": "..."}
  ],
  "citations": [
    {"claim": "...", "evidence_quote": "...", "chunk_id": "...", "source": "...",
     "section": 0}
  ]
}
```

每条 citation 都必须带 `section`：它是这条引用所支持的 `sections` 下标（0 起），\
填错会让覆盖率统计失真，**不得随手填 0**。
"""


def _chunk_lines(chunks: list[RetrievedChunk]) -> str:
    """把检索块序列化成给模型看的文本，块之间用明显分隔。

    每块带上 ``chunk_id`` 与来源，模型产出 ``evidence_quote``/``chunk_id`` 时
    才有据可依；分隔符用不会出现在正文里的记号，便于模型识别边界。
    """
    lines: list[str] = []
    for i, c in enumerate(chunks, start=1):
        lines.append(
            f"[{i}] chunk_id={c.chunk_id} source={c.source} "
            f"section={c.section_title}\n{c.text}"
        )
    return "\n\n".join(lines)


_STRENGTH_BANDS: tuple[tuple[float, str], ...] = (
    (0.20, "接近中间（不要强调任何一端的偏向）"),
    (0.50, "轻度偏向"),
    (0.80, "明显偏向"),
    (float("inf"), "非常明显"),
)


_FACTOR_BELOW_CUTOFF = "低于参考线（属正常范围，不要写成偏向）"


def strength_band(strength: float | None, kind: str = "dimension") -> str:
    """把 ``strength`` 折成**模型必须照抄**的程度词。

    为什么不让模型自己判断强弱：实测拿 E=+0.33 / S=-0.67 的档案去问，模型写出
    「外倾和思考倾向较为明显，直觉和判断倾向相对温和」——把最弱的一维说成最明显，
    把最强的一维说成最温和，正好反了。strength 是数值，程度词必须由代码定，
    否则每次生成都可能得到一套自相矛盾的叙述。

    ``kind`` 为什么必须传：两种测量项的 strength 方向语义不同。

    * ``dimension``（MBTI 的四维）是**双极**的：负值只表示落在另一极，
      ``|strength|`` 越大就是"离中间越远"，所以用绝对值分档，方向由 ``pole`` 说明。
    * ``factor``（SCL-90 的十因子）是**单极**的：负值意味着低于临床参考线，
      那是"没有这方面的问题"，不是"轻度偏向低分"。20 条验收用例实测抓到过这个
      错误——十项全 1 分的档案被写成「位于正常范围，程度为「明显偏向」」，
      一句话里自相矛盾，而它读起来还很像专业表述。
    """
    if strength is None:
        return "无"
    if kind == "factor" and strength < 0:
        return _FACTOR_BELOW_CUTOFF
    mag = abs(strength)
    for upper, label in _STRENGTH_BANDS:
        if mag < upper:
            return label
    return "非常明显"


def _measures_lines(req: EnhanceRequest) -> str:
    """把测评结果序列化给模型。strength/factor_level 已由 Go 侧算好，直接转述。"""
    lines: list[str] = [f"测评：{req.test.name}", f"总体：{req.result.level}"]
    for m in req.measures:
        bits = [f"{m.label}({m.kind}) score={m.score}"]
        if m.strength is not None:
            bits.append(f"strength={m.strength:+.2f}")
            bits.append(f"程度={strength_band(m.strength, m.kind)}")
        if m.factor_level:
            bits.append(f"factor_level={m.factor_level}")
        if m.threshold is not None and m.pole:
            bits.append(f"threshold={m.threshold} pole={m.pole}")
        lines.append("；".join(bits))
    if req.needs_referral:
        lines.append("【注意】本次结果标记了 needs_referral，请给出寻求专业帮助的建议。")
    return "\n".join(lines)


def build_prompt(req: EnhanceRequest, chunks: list[RetrievedChunk]) -> tuple[str, str]:
    """返回 ``(system_prompt, user_prompt)``。

    ``system_prompt`` 必须是**逐字节稳定**的常量：任何按请求拼接的内容都要放进
    ``user_prompt``，否则前缀缓存永远打不中。

    system 中必须包含：角色设定、禁止诊断性表述、输出 JSON schema 说明、
    ``evidence_quote`` 的逐字复制要求。
    """
    user = (
        "以下是本次测评结果：\n\n"
        f"{_measures_lines(req)}\n\n"
        "以下是检索到的知识块（引用时必须逐字复制其中的句子）：\n\n"
        f"{_chunk_lines(chunks)}\n\n"
        "请据此生成解读，严格按 system 指定的 JSON 结构输出。"
    )
    return _SYSTEM_PROMPT, user


def _parse_llm_json(raw: str) -> dict:
    """从 LLM 输出里抠出 JSON 对象。

    模型偶尔会包一层 ```json ...``` 围栏，或前后带废话；先按裸解析，失败则
    尝试去掉围栏，再失败则取第一个 ``{`` 到最后一个 ``}`` 之间。
    """
    text = raw.strip()
    # 去代码围栏。这里必须取 split 的**中间段**：```json\n{...}\n``` 切开是
    # ['', 'json\n{...}\n', '']，原先写 [-1] 会拿到空串，于是围栏一出现就必挂。
    if text.startswith("```"):
        parts = text.split("```")
        text = (parts[1] if len(parts) > 1 else "").strip()
        if text.startswith("json"):
            text = text[4:].strip()
    try:
        return json.loads(text)
    except json.JSONDecodeError:
        start = text.find("{")
        end = text.rfind("}")
        if start != -1 and end != -1 and end > start:
            return json.loads(text[start : end + 1])
        raise


_MAX_QUOTE_CHARS = 64  # 与 Citation.evidence_quote / Suggestion.evidence_quote 的 max_length 一致
_QUOTE_BOUNDARIES = "。！？；…"
# 次级边界：句末标点一个都没有时的退路。
# 实测 20 条验收用例里，被丢弃的引文**全部**是"超过 64 字且整段没有句号"的
# 长复句（语料里常见「A 倾向的人会…，擅长…，因此…」），一整条 citation 就这么没了。
# 在逗号/顿号处收尾虽然读起来像半句，但它是**逐字原文**，Stage 1 与 Stage 1.5
# 都照常能验——丢掉它才是真的损失。
_QUOTE_SOFT_BOUNDARIES = "，、：,"
_QUOTE_MIN_SOFT_CHARS = 20  # 次级边界至少要保住 20 字，否则句子太碎，不如不要


def _fit_quote(raw: object) -> str | None:
    """把模型给的引文裁到契约允许的长度，保留完整句子。

    为什么必须裁：``evidence_quote`` 的 ``max_length=64`` 是硬约束，超长会抛
    ``ValidationError``。citation 那边异常被吞掉（只是少一条），suggestion 那边
    原先没有 try ——**一条 70 字的引文能让整次增强 500**。实测踩过：
    ``scl90-mild`` 档案直接抛 ``string_too_long``。

    剪裁策略：优先在句末标点处切；整段没有句末标点时退到逗号/顿号（且至少保住
    ``_QUOTE_MIN_SOFT_CHARS`` 字）；连这个都没有才返回 None，调用方按「没有引文」
    处理。**任何情况下都不做硬切**——硬切出来的半句既过不了验证，又容易被误认为
    原文的完整表述。
    """
    quote = str(raw or "").strip()
    if not quote:
        return None
    if len(quote) <= _MAX_QUOTE_CHARS:
        return quote

    head = quote[:_MAX_QUOTE_CHARS]
    cut = max(head.rfind(ch) for ch in _QUOTE_BOUNDARIES)
    if cut > 0:
        return head[: cut + 1]

    soft = max(head.rfind(ch) for ch in _QUOTE_SOFT_BOUNDARIES)
    if soft + 1 >= _QUOTE_MIN_SOFT_CHARS:
        logger.info("引文无句末标点，改在逗号处收尾：%s…", quote[:40])
        return head[: soft + 1]

    logger.warning("引文超过 %d 字且找不到任何可切分标点，丢弃：%s…", _MAX_QUOTE_CHARS, quote[:40])
    return None


def _clean_citation(c: dict, n_sections: int) -> Citation | None:
    """把模型产出的一条 citation 转成 ``Citation``，过敏感词过滤器。

    任一环节不合格返回 ``None``（丢弃），不抛穿：单条坏引文不该拖垮整次生成。

    ``section`` 越界（模型把 0 起下标写成 1 起，或编了个不存在的节）一律置
    ``None`` 而不是夹到合法范围——**猜一个节号会让 ``coverage`` 虚高**，
    而覆盖率正是用来拦住"每节只写一句保守结论"的指标（ADR-0003）。
    """
    try:
        claim = (c.get("claim") or "").strip()
        if not claim:
            return None
        if not is_claim_safe(claim):
            logger.warning("敏感词过滤器拦截 claim：%s", claim[:80])
            return None
        raw_section = c.get("section")
        section: int | None = None
        if isinstance(raw_section, int) and not isinstance(raw_section, bool):
            if 0 <= raw_section < n_sections:
                section = raw_section
            else:
                logger.debug("citation 的 section 越界，置 None：%r（共 %d 节）", raw_section, n_sections)
        cit = Citation(
            claim=claim,
            source=str(c.get("source") or ""),
            chunk_id=str(c.get("chunk_id") or ""),
            evidence_quote=_fit_quote(c.get("evidence_quote")),
            section=section,
        )
        return cit
    except Exception as exc:  # noqa: BLE001 —— 单条损坏，丢弃即可
        logger.warning("丢弃一条坏 citation：%s", exc)
        return None


def _clean_suggestion(s: dict) -> Suggestion | None:
    """把一条建议转成 ``Suggestion``；不合格返回 None。

    必须逐条 try：建议是列表里的独立元素，一条超长引文/一个坏字段不该让整次
    增强失败。此前这里是列表推导式里直接构造，``max_length`` 一超就 500。
    """
    try:
        text = str(s.get("text") or "").strip()
        basis = str(s.get("basis") or "").strip()
        if not text or not basis:
            return None
        # 建议文本与依据都要过敏感词过滤：原先只过滤 citations，
        # 「建议停用抗抑郁药」这类句子可以原样出现在建议里。
        if not is_claim_safe(text) or not is_claim_safe(basis):
            logger.warning("敏感词过滤器拦截 suggestion：%s", text[:80])
            return None
        return Suggestion(
            text=text,
            basis=basis,
            chunk_id=str(s.get("chunk_id") or ""),
            evidence_quote=_fit_quote(s.get("evidence_quote")),
        )
    except Exception as exc:  # noqa: BLE001 —— 单条损坏，丢弃即可
        logger.warning("丢弃一条坏 suggestion：%s", exc)
        return None


def _to_enhance_data(
    req: EnhanceRequest,
    payload: dict,
    model: str,
    prompt_tokens: int,
    completion_tokens: int,
    latency_ms: int,
) -> EnhanceData:
    sections = [
        Section(title=str(s.get("title") or ""), content=str(s.get("content") or ""))
        for s in (payload.get("sections") or [])
        if isinstance(s, dict)
    ]
    # 建议也要过敏感词过滤器：原先只过滤 citations，于是「建议停用抗抑郁药」这类
    # 句子可以原样出现在 suggestions 里发给用户（is_claim_safe 的给药规则不做
    # 否定豁免，所以「不建议服用」同样会被拦下——那是刻意的）。
    suggestions = [
        sugg
        for s in (payload.get("suggestions") or [])
        if isinstance(s, dict) and (sugg := _clean_suggestion(s)) is not None
    ]
    citations = [
        cit
        for c in (payload.get("citations") or [])
        if isinstance(c, dict) and (cit := _clean_citation(c, len(sections))) is not None
    ]
    return EnhanceData(
        enhanced_text=str(payload.get("enhanced_text") or ""),
        sections=sections,
        suggestions=suggestions,
        citations=citations,
        model=model,
        prompt_tokens=prompt_tokens,
        completion_tokens=completion_tokens,
        latency_ms=latency_ms,
        rag_status="ready",
    )


async def synthesize(req: EnhanceRequest, chunks: list[RetrievedChunk]) -> EnhanceData:
    """调 DeepSeek 并解析成 ``EnhanceData``。

    返回前必须：
      1. 逐条 claim 过 ``is_claim_safe``；
      2. 写清 ``model`` / ``prompt_tokens`` / ``completion_tokens`` / ``latency_ms``
         （仅供服务端统计，Go 侧不得透传给前端）。

    Raises:
        RuntimeError: LLM 不可用、网络失败、或输出无法解析成合法 JSON。
    """
    s = get_settings()
    if not s.llm_enabled:
        raise RuntimeError("未配置 llm_api_key，无法调用 DeepSeek")

    system, user = build_prompt(req, chunks)

    body = {
        "model": s.llm_model,
        "messages": [
            {"role": "system", "content": system},
            {"role": "user", "content": user},
        ],
        "max_tokens": s.llm_max_tokens,
        "temperature": s.llm_temperature,
        "response_format": {"type": "json_object"},
    }
    # ⚠️ 必须显式关思考。DeepSeek V4 系列默认带 reasoning，实测在完整 prompt
    # （system 609 + user 2497 字）下，2048 的 max_tokens **全部被
    # reasoning_tokens 吃光**，`finish_reason="length"` 而 `content=""`，
    # 于是 JSON 解析直接炸在 "Expecting value: line 1 column 1"。
    # 该现象只在真实调用下暴露（pytest 全 mock 掉了 synthesize）。
    # 关掉思考的理由：本任务是「在给定材料内受约束地提取与组织」，不是开放推理；
    # 且思考会吃掉 30s 超时预算与 token 成本。
    if s.llm_thinking == "disabled":
        body["thinking"] = {"type": "disabled"}

    start = time.monotonic()
    try:
        async with httpx.AsyncClient(timeout=s.llm_timeout_seconds) as client:
            resp = await client.post(
                f"{s.llm_base_url.rstrip('/')}/chat/completions",
                headers={"Authorization": f"Bearer {s.llm_api_key}"},
                json=body,
            )
        resp.raise_for_status()
    except httpx.HTTPStatusError as exc:
        raise RuntimeError(f"DeepSeek 返回 {exc.response.status_code}") from exc
    except httpx.HTTPError as exc:
        raise RuntimeError(f"调用 DeepSeek 失败：{exc}") from exc

    latency_ms = int((time.monotonic() - start) * 1000)
    data = resp.json()

    try:
        choice = data["choices"][0]
        message = choice["message"]
        content = message["content"]
        usage = data.get("usage") or {}
        prompt_tokens = int(usage.get("prompt_tokens") or 0)
        completion_tokens = int(usage.get("completion_tokens") or 0)
        finish_reason = choice.get("finish_reason")
        reasoning_tokens = int(
            (usage.get("completion_tokens_details") or {}).get("reasoning_tokens") or 0
        )
    except (KeyError, IndexError, TypeError) as exc:
        raise RuntimeError(f"DeepSeek 响应缺少预期字段：{exc}") from exc

    # 空 content 是最容易被误诊成「JSON 格式错」的故障：真因通常是 token 预算
    # 被 reasoning 吃光（finish_reason="length"）。这里给出可操作的报错，
    # 不要再让它掉进 json.loads 的 "Expecting value"。
    if not (content or "").strip():
        raise RuntimeError(
            "DeepSeek 返回空 content"
            f"（finish_reason={finish_reason}, reasoning_tokens={reasoning_tokens}, "
            f"completion_tokens={completion_tokens}, max_tokens={s.llm_max_tokens}）。"
            "若 reasoning_tokens 接近 max_tokens，说明思考吃光了输出预算"
            "——检查 llm_thinking 是否为 disabled，或调大 llm_max_tokens。"
        )

    try:
        payload = _parse_llm_json(content)
    except json.JSONDecodeError as exc:
        # JSON 坏了绝大多数情况是**输出被截断**（finish_reason="length"），
        # 少数是模型擅自加了尾逗号。两种情况都需要原始响应尾部才能判定，
        # 所以把 finish_reason 与末尾 200 字一起报出来，别只留一句
        # "Expecting ',' delimiter" 让人猜。
        raise RuntimeError(
            f"DeepSeek 返回的 JSON 无法解析：{exc}"
            f"（finish_reason={finish_reason}, completion_tokens={completion_tokens}, "
            f"max_tokens={s.llm_max_tokens}）。"
            f"响应末尾：…{content[-200:]!r}"
        ) from exc
    if not isinstance(payload, dict):
        raise RuntimeError("LLM 输出不是 JSON 对象")

    return _to_enhance_data(
        req,
        payload,
        model=s.llm_model,
        prompt_tokens=prompt_tokens,
        completion_tokens=completion_tokens,
        latency_ms=latency_ms,
    )
