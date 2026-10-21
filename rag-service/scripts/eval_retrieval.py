#!/usr/bin/env python
"""检索质量与输出泛化度量（§8.2 验收 + ADR-0006）。

## 检索质量（P0-3，§9 :783）

跨构念域过滤必须**两条通道单独验**：``typology`` query 不召回 ``symptom`` 块。

- 向量通道：``where={"kb_namespace": {"$in": [...]}}``
- BM25 通道：按 namespace 选**独立索引**（ADR-0005）

只验向量通道会漏掉 BM25 串味——而 BM25 串味不报错，只是引用「有点对但不对题」。
所以下面把三种模式分别跑一遍，**越域命中数必须是 0**。

## 降级态必须单独报数（ADR-0008 遗留第 1 条）

``trait`` / ``wellbeing`` / ``ability`` 三域语料为 0，落在这些域的测评只检索
``common``。它们的命中率天然低于 ``full`` 态，**不能与 ``full`` 混在一起报一个总
数字**——否则既有内容的三域看起来像坏了，没有内容的三域看起来像还没坏。
fixture 里用 ``"degraded": true`` 标记，本脚本分开统计。

## 泛化度量（ADR-0006 / 主文档 §13）

只有「可用率 ≥ 90%」这类主观指标时，无法证明增强真的解决了泛化问题
（同一个模板套不同得分档案也能「可用」）。要求：

- 20 个分数档案 × 3 个测评
- 任意两档案输出的 **3-gram Jaccard 中位数 < 0.2**（否则说明在套模板）
- 每篇至少包含 **5 个信息要素**
- 与**旧模板基线**做对照

属 **P0-6**，尚未实现（依赖 ``/enhance`` 真正产出报告文本）。
"""

from __future__ import annotations

import argparse
import json
import sys
from pathlib import Path

SERVICE_ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(SERVICE_ROOT))

FIXTURE = SERVICE_ROOT / "tests" / "fixtures" / "queries.jsonl"
MODES = ("vector", "bm25", "hybrid")


def load_fixture(path: Path = FIXTURE) -> list[dict]:
    if not path.exists():
        raise SystemExit(f"fixture 不存在：{path}")
    rows: list[dict] = []
    for lineno, line in enumerate(path.read_text(encoding="utf-8").splitlines(), start=1):
        if not line.strip():
            continue
        try:
            rec = json.loads(line)
        except json.JSONDecodeError as exc:
            raise SystemExit(f"{path}:{lineno} 不是合法 JSON：{exc}") from exc
        for key in ("id", "query", "domains", "expect_any_of"):
            if key not in rec:
                raise SystemExit(f"{path}:{lineno} 缺字段 {key!r}")
        if not rec["domains"]:
            raise SystemExit(f"{path}:{lineno} domains 为空")
        if "common" not in rec["domains"]:
            # 每个测评的检索范围都是「本域 + common」（ADR-0008），漏掉 common
            # 意味着 fixture 写的不是 Go 侧真会发出来的 domains，评测结论无意义。
            raise SystemExit(f"{path}:{lineno} domains 必须包含 'common'（ADR-0008）")
        rows.append(rec)
    if not rows:
        raise SystemExit(f"{path} 是空的")
    return rows


def run_mode(rows: list[dict], mode: str, top_k: int, backend: str | None = None) -> dict:
    from app.rag.retriever import retrieve

    detail: list[dict] = []
    violations: list[str] = []

    for rec in rows:
        domains = list(rec["domains"])
        try:
            hits = retrieve(
                [rec["query"]],
                namespaces=domains,
                top_k=top_k,
                mode=mode,
                bm25_backend=backend,
            )
        except Exception as exc:  # noqa: BLE001 —— 评测脚本要把失败原因打出来
            detail.append(
                {
                    "id": rec["id"],
                    "degraded": bool(rec.get("degraded")),
                    "returned": 0,
                    "hit": False,
                    "error": f"{type(exc).__name__}: {exc}",
                    "chunk_ids": [],
                }
            )
            continue

        bad = [h.chunk_id for h in hits if h.kb_namespace not in domains]
        if bad:
            violations.append(f"{rec['id']} → {bad}")

        # 逐块判断，才能算出「第一个命中的排第几」。
        # 只看「top-5 里有没有关键词」会把「命中在最后一名」和「命中在第一名」
        # 算成一样好——而这两者对报告质量的影响完全不同。
        rank = 0
        matched: list[str] = []
        for pos, h in enumerate(hits, start=1):
            blob = f"{h.section_title}\n{h.text}"
            kws = [kw for kw in rec["expect_any_of"] if kw in blob]
            if kws:
                rank = pos
                matched = kws
                break
        detail.append(
            {
                "id": rec["id"],
                "degraded": bool(rec.get("degraded")),
                "returned": len(hits),
                "hit": bool(matched),
                "rank": rank,
                "matched": matched,
                "chunk_ids": [h.chunk_id for h in hits],
                "namespaces": [h.kb_namespace for h in hits],
            }
        )

    def _rate(items: list[dict]) -> tuple[int, int]:
        return sum(1 for i in items if i["hit"]), len(items)

    def _mrr(items: list[dict]) -> float:
        """Mean Reciprocal Rank：命中在第 1 名得 1.0，第 5 名得 0.2，未命中得 0。"""
        if not items:
            return 0.0
        return sum(1.0 / i["rank"] for i in items if i.get("rank")) / len(items)

    full = [i for i in detail if not i["degraded"]]
    degraded = [i for i in detail if i["degraded"]]
    fh, fn = _rate(full)
    dh, dn = _rate(degraded)

    return {
        "mode": mode,
        "violations": violations,
        "detail": detail,
        "full_hit": fh,
        "full_total": fn,
        "full_mrr": round(_mrr(full), 4),
        "degraded_hit": dh,
        "degraded_total": dn,
        "degraded_mrr": round(_mrr(degraded), 4),
        "empty": sum(1 for i in detail if i["returned"] == 0),
    }


def report_retrieval(
    rows: list[dict],
    top_k: int,
    as_json: bool,
    backend: str | None = None,
) -> int:
    results = [run_mode(rows, mode, top_k, backend) for mode in MODES]

    from app.core.config import get_settings

    effective = backend or get_settings().bm25_backend

    if as_json:
        print(
            json.dumps(
                {"fixture": len(rows), "top_k": top_k, "bm25_backend": effective, "modes": results},
                ensure_ascii=False,
                indent=2,
            )
        )
        return 1 if any(r["violations"] for r in results) else 0

    from app.core.config import get_settings

    effective = backend or get_settings().bm25_backend
    print(f"fixture：{len(rows)} 条 query，top_k={top_k}，bm25 后端={effective}")
    ok = True
    for r in results:
        print(f"\n=== 通道 {r['mode']} ===")
        print(
            f"  越域命中：{len(r['violations'])} 处"
            + ("  ✅" if not r["violations"] else "  ❌ " + "; ".join(r["violations"][:5]))
        )
        print(f"  返回空结果：{r['empty']} 条")
        if r["full_total"]:
            print(
                f"  full 态关键词命中：{r['full_hit']}/{r['full_total']} = "
                f"{r['full_hit'] / r['full_total']:.0%}  MRR={r['full_mrr']:.3f}"
            )
        if r["degraded_total"]:
            print(
                f"  降级态关键词命中：{r['degraded_hit']}/{r['degraded_total']} = "
                f"{r['degraded_hit'] / r['degraded_total']:.0%}  MRR={r['degraded_mrr']:.3f}"
                f"  ← fixture 里标了 degraded 的条目；语料补齐后应为 0 条"
            )
        if r["violations"]:
            ok = False
            print("  逐条：")
            for line in r["violations"]:
                print(f"    - {line}")

    print("\n=== 判定 ===")
    print("越域过滤（三条通道全部为 0）：" + ("✅ 通过" if ok else "❌ 不通过"))
    return 0 if ok else 1


def main() -> int:
    parser = argparse.ArgumentParser(description="检索质量与泛化度量")
    parser.add_argument("--mode", choices=["retrieval", "diversity"], default="retrieval")
    parser.add_argument("--fixture", type=Path, default=FIXTURE)
    parser.add_argument("--top-k", type=int, default=None)
    parser.add_argument("--json", action="store_true")
    parser.add_argument(
        "--bm25-backend",
        choices=["jieba", "llamaindex"],
        default=None,
        help="BM25 通道实现（ADR-0009 实测用）；默认取 settings.bm25_backend",
    )
    args = parser.parse_args()

    if args.mode == "diversity":
        print("❌ 尚未实现：泛化度量属 P0-6 / ADR-0006，依赖 /enhance 产出报告文本", file=sys.stderr)
        return 2

    from app.core.config import get_settings

    top_k = args.top_k or get_settings().top_k
    return report_retrieval(load_fixture(args.fixture), top_k, args.json, args.bm25_backend)


if __name__ == "__main__":
    raise SystemExit(main())
