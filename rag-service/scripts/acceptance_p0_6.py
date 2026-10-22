"""P0-6 联调验收：20 条真实测评用例，走**完整链路**（Go 提交接口 → 队列 → worker → 回填）。

为什么不用 `e2e_enhance.py` 那套直连 Python 的脚本：P0-6 的验收对象是"用户提交一份
测评，最后在详情接口里拿到增强解读"这件事本身。直连 `/enhance` 会跳过 Go 侧的
`BuildRequest`、3 秒同步期限、队列、worker 回填——而那一段恰好是缺陷最多的地方
（参考 `doc/RAG增强系统开发文档.md` §9 P0-5 验收记录里那两个真 bug）。

用法（需要 Go 服务器在 8000、rag-service 在 8001、`rag.enable: true`）::

    /Users/iron/rag/bin/python scripts/acceptance_p0_6.py
    /Users/iron/rag/bin/python scripts/acceptance_p0_6.py --limit 4      # 先跑几条看看

产物写到 `--out`（默认 `.tmp/acceptance/`）：

* `cases.json`    —— 20 条用例的作答与预期
* `results.jsonl` —— 每条用例的原始响应（含耗时、rag 内容）
* `review.csv`    —— 人工评审表（每条用例的每一节一行，留给评审人打分）
* `summary.json`  —— 汇总指标

**它不判断"解读写得好不好"**：那需要人。脚本只负责把机器能判的判掉
（有没有 rag、引用是否 verified、claim 是否引文片段、chunk 是否越域、耗时），
并把需要人看的内容排成一张表。
"""

from __future__ import annotations

import argparse
import csv
import json
import re
import statistics
import sys
import time
import unicodedata
from dataclasses import dataclass
from pathlib import Path
from typing import Any

import httpx
import jwt

BASE_URL = "http://127.0.0.1:8000"
JWT_KEY = "db03d23b03ec405793b38f10592a2f34"  # internal/consts/consts.go:4
USER_ID = 9001
USERNAME = "rag_acceptance"

# 每个测评允许召回的构念域（Go 侧 `consts.RetrievalDomains` 的等价物，
# 这里写死是为了让脚本能独立断言"越域"，而不是信任被测系统自己说的）
ALLOWED_DOMAINS = {
    13: {"typology", "common"},
    1: {"symptom", "common"},
}

# MBTI 题号（相对）→ 维度；每题最高 7 分（见 internal/logic/evaluation/MBTI.go）
MBTI_LAYOUT = {
    "E": [1, 5, 9],
    "S": [2, 6, 10],
    "T": [3, 7],
    "J": [4, 8],
}

# SCL-90 十因子题号（1 基，见 internal/logic/evaluation/SCL90.go::SCL90Factors）
SCL90_FACTORS = {
    "somatization": [1, 4, 12, 27, 40, 42, 48, 49, 52, 53, 56, 58],
    "obsessive": [3, 9, 10, 28, 38, 45, 46, 51, 55, 65],
    "interpersonal": [6, 21, 34, 36, 37, 41, 61, 69, 73],
    "depression": [5, 14, 15, 20, 22, 26, 29, 30, 31, 32, 54, 71, 79],
    "anxiety": [2, 17, 23, 33, 39, 57, 72, 78, 80, 86],
    "hostility": [11, 24, 63, 67, 74, 81],
    "phobic": [13, 25, 47, 50, 70, 75, 82],
    "paranoid": [8, 18, 43, 68, 76, 83],
    "psychoticism": [7, 16, 35, 62, 77, 84, 85, 87, 88, 90],
    "other": [19, 44, 59, 60, 64, 66, 89],
}

SCL90_CRISIS_ITEM = 58  # 「有轻生念头」

QUESTION_IDS = {13: list(range(376, 386)), 1: list(range(1, 96))}


@dataclass
class Case:
    case_id: str
    test_id: int
    test_name: str
    profile: str
    answers: list[dict[str, Any]]


def _mbti_case(case_id: str, profile: str, sums: dict[str, int]) -> Case:
    """按四个维度的目标总分构造作答。"""
    scores = {qid: 1 for qid in QUESTION_IDS[13]}
    for dim, targets in sums.items():
        qids = MBTI_LAYOUT[dim]
        per, extra = divmod(targets, len(qids))
        for i, qid in enumerate(qids):
            # 每题 1~7 分；把余数摊到前几题上，避免超界
            scores[qid] = min(7, per + (1 if i < extra else 0))
    answers = [
        {"question_id": qid, "answer": str(score), "score": score}
        for qid, score in sorted(scores.items())
    ]
    return Case(case_id, 13, "16型人格测试（MBTI）", profile, answers)


def _scl90_case(case_id: str, profile: str, boosts: dict[str, int]) -> Case:
    """全题默认 1 分，按因子把该因子的题项抬到指定分值。"""
    scores = {qid: 1 for qid in QUESTION_IDS[1]}
    for factor, value in boosts.items():
        if factor == "crisis":
            scores[SCL90_CRISIS_ITEM] = value
            continue
        for qid in SCL90_FACTORS[factor]:
            scores[qid] = value
    answers = [
        {"question_id": qid, "answer": str(score), "score": score}
        for qid, score in sorted(scores.items())
    ]
    return Case(case_id, 1, "SCL-90症状自评量表", profile, answers)


def build_cases() -> list[Case]:
    cases = [
        _mbti_case("mbti-01", "ESTJ 强偏好", {"E": 18, "S": 18, "T": 12, "J": 12}),
        _mbti_case("mbti-02", "ENFP", {"E": 17, "S": 5, "T": 5, "J": 4}),
        _mbti_case("mbti-03", "ISTJ", {"E": 5, "S": 18, "T": 13, "J": 11}),
        _mbti_case("mbti-04", "ESTP", {"E": 18, "S": 17, "T": 13, "J": 3}),
        _mbti_case("mbti-05", "四维全贴阈值", {"E": 12, "S": 13, "T": 8, "J": 8}),
        _mbti_case("mbti-06", "ENTJ（与历史用例同档，验缓存）", {"E": 15, "S": 6, "T": 11, "J": 9}),
        _mbti_case("mbti-07", "INFP 极低分", {"E": 4, "S": 4, "T": 3, "J": 3}),
        _mbti_case("mbti-08", "四维均低于阈值", {"E": 10, "S": 10, "T": 6, "J": 6}),
        _mbti_case("mbti-09", "四维均略高于阈值", {"E": 14, "S": 14, "T": 9, "J": 9}),
        _mbti_case("mbti-10", "反向极端（I/N/F/P 全高）", {"E": 3, "S": 3, "T": 2, "J": 2}),
        _scl90_case("scl-01", "仅躯体化偏高", {"somatization": 5}),
        _scl90_case("scl-02", "抑郁偏高", {"depression": 4}),
        _scl90_case("scl-03", "焦虑偏高", {"anxiety": 5}),
        _scl90_case("scl-04", "敌对偏高", {"hostility": 5}),
        _scl90_case("scl-05", "恐怖偏高", {"phobic": 5}),
        _scl90_case("scl-06", "偏执偏高", {"paranoid": 5}),
        _scl90_case("scl-07", "精神病性偏高", {"psychoticism": 4}),
        _scl90_case("scl-08", "强迫偏高", {"obsessive": 5}),
        _scl90_case("scl-09", "人际敏感偏高", {"interpersonal": 5}),
        _scl90_case("scl-10", "全题 3 分（整体偏高）", {f: 3 for f in SCL90_FACTORS}),
        _scl90_case("scl-11", "全题 2 分（整体轻度）", {f: 2 for f in SCL90_FACTORS}),
        _scl90_case("scl-12", "危机信号（题 58 答 3）", {"crisis": 3}),
        # 以下 3 条是把评审样本凑到 ≥100 节用的"组合档案"：
        # 单因子偏高只考"能不能讲清一个构念"，多因子同时偏高才考"会不会把几件事混成一件"。
        _scl90_case("scl-13", "抑郁+焦虑双高", {"depression": 4, "anxiety": 4}),
        _scl90_case("scl-14", "全题 4 分（整体明显偏高）", {f: 4 for f in SCL90_FACTORS}),
        _scl90_case("scl-15", "强迫+人际+敌对三高", {"obsessive": 5, "interpersonal": 4, "hostility": 4}),
    ]
    return cases


def make_token() -> str:
    now = int(time.time())
    return jwt.encode(
        {"Id": USER_ID, "Username": USERNAME, "exp": now + 86400, "iat": now},
        JWT_KEY,
        algorithm="HS256",
    )


def normalize(text: str) -> str:
    text = unicodedata.normalize("NFKC", text)
    return re.sub(r"\s+", "", text)


def submit(client: httpx.Client, token: str, case: Case) -> dict[str, Any]:
    t0 = time.perf_counter()
    resp = client.post(
        f"{BASE_URL}/api/user/tests/submit",
        headers={"Authorization": f"Bearer {token}"},
        json={
            "test_id": case.test_id,
            "test_name": case.test_name,
            "answers": case.answers,
        },
        timeout=120.0,
    )
    elapsed_ms = (time.perf_counter() - t0) * 1000
    resp.raise_for_status()
    body = resp.json()
    return {
        "http_status": resp.status_code,
        "submit_ms": round(elapsed_ms, 1),
        "record_id": body["data"]["record_id"],
        "inline_rag": "rag" in (body["data"].get("result") or {}),
    }


def poll(client: httpx.Client, token: str, record_id: int, deadline_s: float = 180.0) -> dict[str, Any]:
    """等 worker 回填。返回最终 `result`。"""
    t0 = time.perf_counter()
    while time.perf_counter() - t0 < deadline_s:
        resp = client.get(
            f"{BASE_URL}/api/user/tests/{record_id}",
            headers={"Authorization": f"Bearer {token}"},
            timeout=30.0,
        )
        resp.raise_for_status()
        result = resp.json()["data"]["result"]
        status = result.get("rag_status")
        if status in {"ready", "degraded", "failed"}:
            return {"wait_s": round(time.perf_counter() - t0, 1), "result": result}
        time.sleep(1.0)
    return {"wait_s": round(time.perf_counter() - t0, 1), "result": result, "timeout": True}


def audit(result: dict[str, Any], test_id: int) -> dict[str, Any]:
    rag = result.get("rag") or {}
    citations = rag.get("citations") or []
    suggestions = rag.get("suggestions") or []
    allowed = ALLOWED_DOMAINS[test_id]

    out_of_domain = []
    extractive_bad = []
    for c in citations:
        prefix = str(c.get("chunk_id", "")).split("_")[0]
        if prefix and prefix not in allowed:
            out_of_domain.append(c.get("chunk_id"))
        quote, claim = c.get("evidence_quote") or "", c.get("claim") or ""
        if quote and normalize(claim) not in normalize(quote):
            extractive_bad.append(claim)

    bad_suggestion_basis = []
    for s in suggestions:
        quote, basis = s.get("evidence_quote") or "", s.get("basis") or ""
        if quote and normalize(basis) not in normalize(quote):
            bad_suggestion_basis.append(basis)

    return {
        "rag_status": result.get("rag_status"),
        "sections": len(rag.get("sections") or []),
        "citations": len(citations),
        "citations_verified": sum(1 for c in citations if c.get("verified")),
        "suggestions": len(suggestions),
        "suggestions_verified": sum(1 for s in suggestions if s.get("verified")),
        "out_of_domain": out_of_domain,
        "non_extractive_claims": extractive_bad,
        "non_extractive_basis": bad_suggestion_basis,
        "domains": sorted({str(c.get("chunk_id", "")).split("_")[0] for c in citations if c.get("chunk_id")}),
        "enhanced_chars": len(rag.get("enhanced_text") or ""),
        "has_disclaimer": bool(rag.get("disclaimer")),
    }


def write_review_csv(path: Path, rows: list[dict[str, Any]]) -> None:
    fieldnames = [
        "编号",
        "测评",
        "档案",
        "章节标题",
        "内容摘要",
        "引用条数",
        "已验证",
        "分域",
        "可用性评分(3=可用 2=勉强 1=不可用)",
        "问题类型(事实错/跑题/空话/口径过硬/其他)",
        "备注",
    ]
    with path.open("w", encoding="utf-8-sig", newline="") as fh:
        writer = csv.DictWriter(fh, fieldnames=fieldnames)
        writer.writeheader()
        writer.writerows(rows)


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--out", default=".tmp/acceptance")
    parser.add_argument("--limit", type=int, default=0, help="只跑前 N 条（调试用）")
    args = parser.parse_args()

    out_dir = Path(args.out)
    out_dir.mkdir(parents=True, exist_ok=True)

    cases = build_cases()
    if args.limit:
        cases = cases[: args.limit]
    (out_dir / "cases.json").write_text(
        json.dumps([c.__dict__ for c in cases], ensure_ascii=False, indent=2), encoding="utf-8"
    )

    token = make_token()
    records: list[dict[str, Any]] = []
    review_rows: list[dict[str, Any]] = []

    with httpx.Client() as client:
        for i, case in enumerate(cases, 1):
            print(f"[{i}/{len(cases)}] {case.case_id} {case.profile} …", flush=True)
            sub = submit(client, token, case)
            waited = poll(client, token, sub["record_id"])
            audited = audit(waited["result"], case.test_id)
            record = {
                "case_id": case.case_id,
                "profile": case.profile,
                "test_id": case.test_id,
                **sub,
                **waited,
                **audited,
            }
            records.append(record)

            rag = waited["result"].get("rag") or {}
            for idx, section in enumerate(rag.get("sections") or []):
                cited = [c for c in (rag.get("citations") or []) if c.get("section") == idx]
                review_rows.append(
                    {
                        "编号": f"{case.case_id}-s{idx + 1}",
                        "测评": case.test_name,
                        "档案": case.profile,
                        "章节标题": (section.get("title") or "")[:40],
                        "内容摘要": re.sub(r"\s+", " ", section.get("content") or "")[:160],
                        "引用条数": len(cited),
                        "已验证": sum(1 for c in cited if c.get("verified")),
                        "分域": ",".join(sorted({str(c.get("chunk_id", "")).split("_")[0] for c in cited})),
                        "可用性评分(3=可用 2=勉强 1=不可用)": "",
                        "问题类型(事实错/跑题/空话/口径过硬/其他)": "",
                        "备注": "",
                    }
                )
            print(
                f"    record={sub['record_id']} submit={sub['submit_ms']}ms wait={waited['wait_s']}s "
                f"status={audited['rag_status']} 节={audited['sections']} "
                f"引用={audited['citations']}({audited['citations_verified']} verified) "
                f"建议={audited['suggestions']} 越域={len(audited['out_of_domain'])}",
                flush=True,
            )

    with (out_dir / "results.jsonl").open("w", encoding="utf-8") as fh:
        for r in records:
            fh.write(json.dumps(r, ensure_ascii=False) + "\n")

    write_review_csv(out_dir / "review.csv", review_rows)

    ok = [r for r in records if r["rag_status"] == "ready"]
    cites = sum(r["citations"] for r in records)
    verified = sum(r["citations_verified"] for r in records)
    sugg = sum(r["suggestions"] for r in records)
    sugg_ok = sum(r["suggestions_verified"] for r in records)
    summary = {
        "cases": len(records),
        "ready": len(ok),
        "degraded": sum(1 for r in records if r["rag_status"] == "degraded"),
        "failed": sum(1 for r in records if r["rag_status"] == "failed"),
        "citations": cites,
        "citations_verified": verified,
        "citation_pass_rate": round(verified / cites, 3) if cites else 0.0,
        "suggestions": sugg,
        "suggestions_verified": sugg_ok,
        "suggestion_pass_rate": round(sugg_ok / sugg, 3) if sugg else 0.0,
        "out_of_domain_total": sum(len(r["out_of_domain"]) for r in records),
        "non_extractive_total": sum(len(r["non_extractive_claims"]) for r in records),
        "sections_total": sum(r["sections"] for r in records),
        "review_rows": len(review_rows),
        "submit_ms_median": round(statistics.median(r["submit_ms"] for r in records), 1),
        "sync_hit": sum(1 for r in records if r["inline_rag"]),
        "wait_s_median": round(statistics.median(r["wait_s"] for r in records), 1),
        "enhanced_chars_median": round(statistics.median(r["enhanced_chars"] for r in records), 1),
    }
    (out_dir / "summary.json").write_text(json.dumps(summary, ensure_ascii=False, indent=2), encoding="utf-8")
    print(json.dumps(summary, ensure_ascii=False, indent=2))
    return 0


if __name__ == "__main__":
    sys.exit(main())
