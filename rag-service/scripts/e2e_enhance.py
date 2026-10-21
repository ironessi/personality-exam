#!/usr/bin/env python3
"""端到端验证：真实 DeepSeek 打通 ``/enhance`` 全链路（P0-4 收尾）。

为什么单独有这个脚本：``pytest`` 里的 74 个用例**全部 mock 掉了** ``synthesize``，
所以「真实 LLM 会不会按 JSON 结构返回」「evidence_quote 是否真的逐字复制」
「Stage 2 NLI 是否真的能分辨错挂结论」这三件事从未被验证过。
在接 Go 之前必须先把这一环打通，否则出问题时分不清是 Go 层还是 Python 层。

用法::

    /Users/iron/rag/bin/python scripts/e2e_enhance.py
    /Users/iron/rag/bin/python scripts/e2e_enhance.py --case scl90
    /Users/iron/rag/bin/python scripts/e2e_enhance.py --dry-run   # 不调 LLM，只看检索

退出码：0 全绿 / 1 有结构性失败。
"""

from __future__ import annotations

import argparse
import asyncio
import json
import sys
import time
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent.parent))

from app.api.schemas import EnhanceRequest  # noqa: E402

# ── 两个真实档案 ────────────────────────────────────────────
#
# strength 的计算口径（Go 侧 P0-5 要照着实现）：
#   dimension: (score - threshold) / (max - threshold)，再 clamp 到 [-1,1]
#   factor:    (score - clinical_cutoff) / clinical_cutoff，再 clamp 到 [-1,1]
#
# 这里手算出数值，是为了**不依赖 Go 侧**也能验证 Python 链路。

CASES: dict[str, dict] = {
    "mbti": {
        "schema_version": 1,
        "test": {
            "code": 13,
            "key": "mbti",
            "name": "16型人格测试（MBTI）",
            "domains": ["typology", "common"],
            "policy": "full",
        },
        "result": {"level": "完成测试", "total": None},
        "measures": [
            # E=15/max21/threshold12 → (15-12)/(21-12)=0.333
            {
                "key": "EI",
                "label": "外倾 / 内倾",
                "kind": "dimension",
                "score": 15,
                "max": 21,
                "threshold": 12,
                "pole": "E",
                "strength": 0.333,
            },
            # S=6 → (6-12)/(21-12)=-0.667 强内倾直觉
            {
                "key": "SN",
                "label": "实感 / 直觉",
                "kind": "dimension",
                "score": 6,
                "max": 21,
                "threshold": 12,
                "pole": "N",
                "strength": -0.667,
            },
            # T=11/max14/threshold8 → (11-8)/(14-8)=0.5
            {
                "key": "TF",
                "label": "思考 / 情感",
                "kind": "dimension",
                "score": 11,
                "max": 14,
                "threshold": 8,
                "pole": "T",
                "strength": 0.5,
            },
            # J=9 → (9-8)/(14-8)=0.167
            {
                "key": "JP",
                "label": "判断 / 知觉",
                "kind": "dimension",
                "score": 9,
                "max": 14,
                "threshold": 8,
                "pole": "J",
                "strength": 0.167,
            },
        ],
        "labels": {"personality_type": "INTJ"},
        "flags": [],
    },
    "scl90": {
        "schema_version": 1,
        "test": {
            "code": 1,
            "key": "scl90",
            "name": "SCL-90 症状自评量表",
            "domains": ["symptom", "common"],
            "policy": "full",
        },
        "result": {"level": "明显异常", "total": 178, "total_max": 450},
        "measures": [
            # (2.4-2.0)/2.0 = 0.2
            {
                "key": "somatization",
                "label": "躯体化",
                "kind": "factor",
                "score": 2.4,
                "max": 5,
                "clinical_cutoff": 2.0,
                "factor_level": "偏高",
                "strength": 0.2,
            },
            # (3.1-2.0)/2.0 = 0.55 → 显著偏高
            {
                "key": "depression",
                "label": "抑郁",
                "kind": "factor",
                "score": 3.1,
                "max": 5,
                "clinical_cutoff": 2.0,
                "factor_level": "偏高",
                "strength": 0.55,
            },
            # (1.4-2.0)/2.0 = -0.3 正常
            {
                "key": "phobic_anxiety",
                "label": "恐怖焦虑",
                "kind": "factor",
                "score": 1.4,
                "max": 5,
                "clinical_cutoff": 2.0,
                "factor_level": "正常",
                "strength": -0.3,
            },
        ],
        "labels": {},
        "flags": ["needs_referral"],
    },
}


def _hr(title: str) -> None:
    print(f"\n{'─' * 68}\n{title}\n{'─' * 68}")


def _strength_audit(payload: dict) -> list[str]:
    """校验手算 strength 与公式一致，防止脚本自己写错。"""
    problems: list[str] = []
    for m in payload["measures"]:
        if m["kind"] == "dimension":
            denom = m["max"] - m["threshold"]
            expect = (m["score"] - m["threshold"]) / denom if denom else 0.0
        else:
            cutoff = m["clinical_cutoff"]
            expect = (m["score"] - cutoff) / cutoff if cutoff else 0.0
        expect = max(-1.0, min(1.0, expect))
        if abs(expect - m["strength"]) > 0.01:
            problems.append(
                f"{m['key']}: 手算 {m['strength']:+.3f} ≠ 公式 {expect:+.3f}"
            )
    return problems


def _label_ngrams(label: str, n: int = 2) -> list[str]:
    """把 label 拆成 ``n``-gram 词表，用于宽松的「对题」判断。

    ``恐怖焦虑`` → ``["恐怖", "怖焦", "焦虑"]``。取 ``n=2`` 是因为中文双字词
    已经足够区分主题，再长就会漏掉「label 比语料多一个字」的情况。
    """
    grams: list[str] = []
    for part in label.split():
        if len(part) < n:
            continue
        grams.extend(part[i : i + n] for i in range(len(part) - n + 1))
    return grams


def _audit_per_measure(req: EnhanceRequest, queries: list[str], namespaces: list[str]) -> bool:
    """逐 measure 单独检索，检查「本域内容是否真的进榜」。

    两个断言：

    1. **每个 measure 的 top-5 里至少有一块来自本构念域**（不是 ``common``）。
       全落 ``common`` 说明检索被「测评工具总览」型 hub 块吃掉了。
    2. ``policy="full"`` 时 **top-1 不得是 ``common``**——``common`` 是兜底层，
       不该盖过本域内容。

    ``common_only`` / 没有非通用域的档案跳过这两条。
    """
    from app.rag.query_builder import _normalize_label
    from app.rag.retriever import retrieve

    own_domains = [d for d in namespaces if d != "common"]
    print("\n[2b] 逐 measure 检索审计")
    if not own_domains:
        print("     ⚠ 只有 common 域（common_only 策略），跳过本域检查")
        return True

    ok = True
    for m, q in zip(req.measures, queries):
        hits = retrieve([q], namespaces=namespaces, top_k=3)
        in_own = [c for c in hits if c.kb_namespace != "common"]
        label = _normalize_label(m.label)
        if not in_own:
            print(f"     ✗ {label}：top-3 全在 common，本域零命中")
            for c in hits:
                print(f"         [{c.kb_namespace}] {c.source} | {c.section_title}")
            ok = False
            continue
        top = hits[0]
        if top.kb_namespace == "common":
            print(
                f"     ✗ {label}：top-1 落在 common —— {top.source} | {top.section_title}"
            )
            print("         common 是兜底层，不该盖过本域内容（见 query_builder 的实测记录）")
            ok = False
            continue
        # 标签词是否出现在命中块里。Corpus 用词常比 label 短（label「恐怖焦虑」
        # 在语料里写作「恐怖」），所以取标签的 2-gram 做匹配；**不匹配只告警**，
        # 因为 label 与语料的措辞差异是正常的，不能当硬失败。
        tokens = _label_ngrams(label)
        matched = any(
            t in (c.text or "") or t in (c.section_title or "")
            for c in hits[:3]
            for t in tokens
        )
        flag = "✓" if matched else "⚠"
        print(
            f"     {flag} {label} → [{top.kb_namespace}] {top.source}"
            f" | {top.section_title}（本域 {len(in_own)}/3）"
        )
        if not matched:
            print(f"         标签词 {tokens} 未出现在 top-3 正文/标题里（仅告警，不算失败）")

    if ok:
        print("     ✓ 每个 measure 都检索到本域且对题的内容")
    return ok


async def run_case(name: str, dry_run: bool) -> bool:
    from app.core.config import get_settings
    from app.core.verifier import verify_citations
    from app.rag.query_builder import build_queries, retrieval_namespaces
    from app.rag.retriever import retrieve

    payload = CASES[name]
    _hr(f"用例 {name.upper()}  ·  policy={payload['test']['policy']}")

    bad = _strength_audit(payload)
    if bad:
        print("✗ strength 手算与公式不一致：")
        for b in bad:
            print(f"    {b}")
        return False
    print("✓ strength 公式自洽")

    req = EnhanceRequest.model_validate(payload)
    print(f"✓ 契约校验通过  code={req.test.code}  domains={req.test.domains}")
    print(f"  needs_referral={req.needs_referral}")

    # ── 1. query builder ──
    queries = build_queries(req)
    namespaces = retrieval_namespaces(req)
    print(f"\n[1] query 共 {len(queries)} 条，检索范围 {namespaces}")
    for q in queries:
        print(f"    · {q}")

    # ── 2. 检索 ──
    t0 = time.perf_counter()
    chunks = retrieve(queries, namespaces=namespaces)
    t_ret = time.perf_counter() - t0
    print(f"\n[2] 检索 {len(chunks)} 块，{t_ret * 1000:.0f}ms")
    if not chunks:
        print("✗ 检索为空——知识库没建好或 domains 错")
        return False
    for c in chunks[:5]:
        print(f"    [{c.kb_namespace}] {c.chunk_id}  {c.source}  fused={c.fused_score:.4f}")

    # 越域检查
    allowed = set(namespaces)
    crossed = [c for c in chunks if c.kb_namespace not in allowed]
    if crossed:
        print(f"✗ 越域命中 {len(crossed)} 块：{[c.kb_namespace for c in crossed]}")
        return False
    print(f"✓ 无越域命中（允许 {sorted(allowed)}）")

    # ── 2b. 逐 measure 检索审计 ──
    #
    # 这是**踩过坑才加的**：曾经 top-5 全是 [common] 的「测评工具总览」块，
    # 每条引文都通过了字面校验、整条链路「通过」，但解读全是泛泛而谈。
    # 越域检查抓不到它（common 本来就在允许范围内），必须单独看「本域是否进榜」。
    audit_ok = _audit_per_measure(req, queries, namespaces)
    if not audit_ok:
        return False

    if dry_run:
        print("\n--dry-run：跳过 LLM")
        return True

    # ── 3. 生成 ──
    from app.rag.synthesizer import synthesize

    settings = get_settings()
    if not settings.llm_enabled:
        print("\n✗ 未配置 llm_api_key，无法真实调用")
        return False
    print(f"\n[3] 调用 DeepSeek  model={settings.llm_model} …")
    t0 = time.perf_counter()
    data = await synthesize(req, chunks)
    t_gen = time.perf_counter() - t0
    print(f"✓ 生成完成 {t_gen:.1f}s  tokens in/out={data.prompt_tokens}/{data.completion_tokens}")
    print(f"  enhanced_text {len(data.enhanced_text)} 字")
    print(f"  sections={len(data.sections)}  suggestions={len(data.suggestions)}"
          f"  citations={len(data.citations)}")
    if not data.enhanced_text.strip():
        print("✗ enhanced_text 为空——LLM 没按 JSON 结构返回")
        return False

    # ── 4. 两级引用验证 ──
    chunks_by_id = {c.chunk_id: c.text for c in chunks}
    passed, details = verify_citations(data.citations, chunks_by_id)
    n1 = sum(1 for d in details if d.passed_stage1)
    n2 = sum(1 for d in details if d.passed_stage2)
    print(f"\n[4] 引用验证  共 {len(details)} 条  Stage1 过 {n1}  Stage2 过 {n2}")
    for d in details:
        mark = "✓" if d.verified else "✗"
        print(f"    {mark} [{d.chunk_id}] {d.claim[:34]}")
        if not d.verified:
            print(f"        reason={d.reason}")

    if data.citations:
        # 关键断言：LLM 给的引文必须能在原文里找到，否则 prompt 约束失效
        if n1 == 0:
            print("✗ 没有任何一条引文通过字面校验——evidence_quote 约束失效")
            return False
        print(f"✓ 字面命中率 {n1}/{len(details)} = {n1 / len(details):.0%}")
    else:
        print("⚠ LLM 未给出任何 citation（结构合法但溯源缺失）")

    print(f"\n{name.upper()} 通过。")
    return True


async def main_async(args) -> int:
    names = ["mbti", "scl90"] if args.case == "all" else [args.case]
    results = {}
    for n in names:
        try:
            results[n] = await run_case(n, args.dry_run)
        except Exception as exc:  # noqa: BLE001
            import traceback

            traceback.print_exc()
            print(f"\n✗ {n} 抛异常：{exc}")
            results[n] = False

    _hr("汇总")
    for n, ok in results.items():
        print(f"  {'✓' if ok else '✗'} {n}")
    return 0 if all(results.values()) else 1


def main() -> int:
    p = argparse.ArgumentParser(description="真实 DeepSeek 端到端验证 /enhance")
    p.add_argument("--case", choices=["mbti", "scl90", "all"], default="all")
    p.add_argument("--dry-run", action="store_true", help="只跑检索，不调 LLM")
    return asyncio.run(main_async(p.parse_args()))


if __name__ == "__main__":
    raise SystemExit(main())
