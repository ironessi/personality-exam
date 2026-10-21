#!/usr/bin/env python3
"""批量质量评测：把 support_precision / coverage 从「单次日志」变成可复现的数字。

``scripts/e2e_enhance.py`` 只跑一遍，回答「这条链路通不通」；本脚本回答
「多跑几份档案，质量指标稳定在什么区间」。两者都要有：单次通过可能是运气，
而人工核验 50 条引用的成本又太高（ADR-0003 的遗留项）。

指标口径严格按 ADR-0003（与 ``app/core/verifier.py`` 一致）：

* ``support_precision`` = 通过两级验证的 citation 数 / 生成的 citation 数
* ``coverage``           = 含 ≥1 条通过 citation 的 section 数 / section 总数
* 建议通过率            = 通过的 suggestion 数 / 生成的 suggestion 数

每个指标都给 **Wilson 95% 置信区间**：小样本下 3/3 和 30/30 的可信度差很远，
只报点估计会让人误以为「100% 就是稳了」。

用法::

    /Users/iron/rag/bin/python scripts/eval_quality.py              # 默认每份档案跑 2 轮
    /Users/iron/rag/bin/python scripts/eval_quality.py --repeat 3
    /Users/iron/rag/bin/python scripts/eval_quality.py --dry-run    # 只跑检索

注意：每次真实生成约消耗 2~3k prompt token + 1~2k completion token，
`--repeat` 直接乘上去。**不要**在调 prompt 的过程中反复跑大 repeat。
"""

from __future__ import annotations

import argparse
import asyncio
import copy
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent.parent))

from app.api.schemas import EnhanceRequest  # noqa: E402
from app.core.verifier import (  # noqa: E402
    stage2_status,
    support_coverage,
    support_precision,
    verify_citations,
    verify_suggestions,
)
from app.rag.query_builder import build_queries, retrieval_namespaces  # noqa: E402
from app.rag.retriever import retrieve  # noqa: E402
from app.rag.synthesizer import synthesize  # noqa: E402
from scripts.e2e_enhance import CASES  # noqa: E402


def _mbti_variant() -> dict:
    """另一份 MBTI 档案：四维全部贴着中间，用来看「全中间」会不会被写成极端。"""
    case = copy.deepcopy(CASES["mbti"])
    for m, (score, strength, pole) in zip(
        case["measures"],
        [(12, 0.0, "E"), (12, 0.0, "S"), (8.0, 0.0, "T"), (8, 0.0, "J")],
    ):
        m["score"] = score
        m["strength"] = strength
        m["pole"] = pole
    case["labels"] = {"personality_type": "ESTJ"}
    return case


def _scl90_mild_variant() -> dict:
    """总分不高、只有个别因子轻度过线的档案（最容易产出「正确的废话」）。"""
    case = copy.deepcopy(CASES["scl90"])
    case["result"] = {"level": "轻度异常", "total": 118, "total_max": 450}
    case["measures"] = case["measures"][:2]
    case["measures"][0]["score"] = 2.05
    case["measures"][0]["strength"] = 0.025
    case["measures"][1]["score"] = 1.8
    case["measures"][1]["strength"] = -0.1
    case["flags"] = []
    return case


PROFILES: dict[str, dict] = {
    "mbti-typical": CASES["mbti"],
    "mbti-all-middle": _mbti_variant(),
    "scl90-referral": CASES["scl90"],
    "scl90-mild": _scl90_mild_variant(),
}


async def run_once(raw: dict, dry_run: bool) -> dict:
    req = EnhanceRequest.model_validate(raw)
    queries = build_queries(req)
    chunks = retrieve(queries, namespaces=retrieval_namespaces(req))
    if not chunks:
        return {"error": "检索为空"}

    if dry_run:
        allowed = set(retrieval_namespaces(req))
        crossed = [c for c in chunks if c.kb_namespace not in allowed]
        return {
            "chunks": len(chunks),
            "crossed": len(crossed),
            "dry_run": True,
        }

    data = await synthesize(req, chunks)
    chunks_by_id = {c.chunk_id: c.text for c in chunks}

    passed, details = verify_citations(data.citations, chunks_by_id)
    sugg_passed, sugg_details = verify_suggestions(data.suggestions, chunks_by_id)

    return {
        "sections": len(data.sections),
        "citations_total": len(details),
        "citations_passed": len(passed),
        "suggestions_total": len(sugg_details),
        "suggestions_passed": len(sugg_passed),
        "text_len": len(data.enhanced_text),
        "prompt_tokens": data.prompt_tokens,
        "completion_tokens": data.completion_tokens,
        "latency_ms": data.latency_ms,
    }


def _fmt_ci(point: float, low: float, high: float, n: int) -> str:
    if n == 0:
        return "  —  "
    return f"{point:.2f} [{low:.2f},{high:.2f}]"


async def main_async(args) -> int:
    ok, err = stage2_status()
    print(f"Stage 2（NLI）: {'就绪' if ok else '不可用 → ' + str(err)}")
    print(f"档案 {len(PROFILES)} 份 × {args.repeat} 轮\n")

    rows: list[dict] = []
    for name, raw in PROFILES.items():
        for i in range(args.repeat):
            try:
                r = await run_once(raw, args.dry_run)
            except Exception as exc:  # noqa: BLE001
                print(f"✗ {name} 第 {i + 1} 轮抛异常：{exc}")
                continue
            r["profile"] = name
            r["round"] = i + 1
            rows.append(r)
            if args.dry_run:
                print(f"  {name} #{i + 1}: chunks={r['chunks']} 越域={r['crossed']}")
            else:
                print(
                    f"  {name} #{i + 1}: sections={r['sections']} "
                    f"citations={r['citations_passed']}/{r['citations_total']} "
                    f"建议={r['suggestions_passed']}/{r['suggestions_total']} "
                    f"{r['text_len']}字 {r['latency_ms']}ms"
                )

    if args.dry_run:
        bad = [r for r in rows if r.get("crossed")]
        print(f"\n越域命中 {len(bad)} 次")
        return 0 if not bad else 1

    total = sum(r["citations_total"] for r in rows)
    passed = sum(r["citations_passed"] for r in rows)
    precision, low, high = support_precision(passed, total)

    # coverage 按「每次运行的 section 覆盖」取平均：整篇覆盖到 0 次也算数据点，
    # 不能只统计"有引用的那几次"，否则覆盖率会系统性偏高。
    covs = []
    for r in rows:
        if r["sections"]:
            # 这里只能拿到总数，用 passed/n 近似 —— 精确 coverage 需要 section 下标，
            # 而那要重跑一遍验证；本脚本只做量级监控，精确值看服务端日志 `coverage=`。
            covs.append(min(r["citations_passed"] / r["sections"], 1.0))
    cov = sum(covs) / len(covs) if covs else 0.0

    sugg_total = sum(r["suggestions_total"] for r in rows)
    sugg_passed = sum(r["suggestions_passed"] for r in rows)
    sp, sl, sh = support_precision(sugg_passed, sugg_total)

    print("\n" + "─" * 60)
    print(f"引用：{passed}/{total}  support_precision = {_fmt_ci(precision, low, high, total)}")
    print(f"建议：{sugg_passed}/{sugg_total}  通过率 = {_fmt_ci(sp, sl, sh, sugg_total)}")
    print(f"section 覆盖（近似）：{cov:.2f}")
    print("注意：精确 coverage 以服务端日志的 `coverage=` 为准；本脚本只是量级监控。")

    if not ok:
        print("\n⚠ Stage 2 不可用：上面的 precision 只反映字面校验，不能当验收数字。")
        return 1
    if total == 0:
        print("\n✗ 一条引用都没有生成——prompt 结构或模型输出有问题。")
        return 1
    return 0


def main() -> int:
    p = argparse.ArgumentParser(description="批量质量评测（precision / coverage）")
    p.add_argument("--repeat", type=int, default=2, help="每份档案跑几轮（默认 2）")
    p.add_argument("--dry-run", action="store_true", help="只跑检索，不调 LLM")
    return asyncio.run(main_async(p.parse_args()))


if __name__ == "__main__":
    raise SystemExit(main())
