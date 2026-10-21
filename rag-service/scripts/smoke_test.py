#!/usr/bin/env python
"""P0-1 环境自检（§9 :736-760 / 附录 C :958-983）。

四项验收：

| # | 验收项 | 判定 |
|---|---|---|
| 1 | 本地 BGE 可加载 | ``embedding_loaded: true`` |
| 2 | 加载 < 3 s、编码 3 条 < 5 s | 见下方「关于计时口径」 |
| 3 | DeepSeek 连通 | ``llm_reachable: true`` |
| 4 | Chroma 持久化目录可写 | 目录可写且集合可开 |

## 关于计时口径（重要）

实测发现「模型加载 < 3 s」这个验收项**取决于把什么算进计时**：

| 口径 | 实测 |
|---|---|
| ``import torch`` + ``import sentence_transformers`` | **5.04 s** |
| 纯 ``SentenceTransformer(path)`` | 0.21 s 冷 / 0.08 s 热 |
| 加载后**第一次** ``encode(3 条)`` | **4.27 s** |
| 第二次 / 第三次 | 0.029 s / 0.008 s |

所以本脚本**分开报告**三个数字：``import_s`` / ``load_s`` / ``first_encode_s``，
并按附录 C 的原意以 ``load_s`` 判定第 2 项（``load_s < 3``）。
``first_encode_s`` 单独报告，因为它才是决定 ADR-0007 那个 3 秒同步期限能否
达成的量——而它靠 ``warmup()`` 在启动期烧掉，不该由用户承担。

用法::

    /Users/iron/rag/bin/python scripts/smoke_test.py
    /Users/iron/rag/bin/python scripts/smoke_test.py --skip-llm   # 无 Key 时
"""

from __future__ import annotations

import argparse
import json
import sys
import time
from pathlib import Path

# 允许直接 `python scripts/smoke_test.py`
SERVICE_ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(SERVICE_ROOT))

RESULTS: dict[str, object] = {}


def _hr(title: str) -> None:
    print(f"\n{'─' * 68}\n  {title}\n{'─' * 68}")


def check_embedding(skip_load: bool = False) -> bool:
    """① 模型可加载 ② 加载与编码耗时 ③ 维度与归一化。"""
    _hr("1-3. 本地向量模型（BAAI/bge-small-zh-v1.5）")

    from app.core.config import get_settings

    s = get_settings()
    model_path = Path(s.embedding_model_path)

    if not model_path.exists():
        print(f"✗ 模型路径不存在：{model_path}")
        RESULTS["embedding_loaded"] = False
        return False
    print(f"  模型路径  {model_path}")
    print(f"  快照目录存在 ✓")

    t0 = time.perf_counter()
    from sentence_transformers import SentenceTransformer  # noqa: F401

    import_s = time.perf_counter() - t0
    print(f"  import 耗时        {import_s:.3f}s   （torch + sentence_transformers）")

    if skip_load:
        print("  --skip-load：跳过实际加载")
        return True

    t0 = time.perf_counter()
    model = SentenceTransformer(s.embedding_model_path)
    load_s = time.perf_counter() - t0
    print(f"  模型加载耗时       {load_s:.3f}s   ← 验收判定用这个（要求 < 3s）")

    # 维度：优先新名，避免 sentence-transformers 5.x 的 FutureWarning
    get_dim = getattr(model, "get_embedding_dimension", None) or getattr(
        model, "get_sentence_embedding_dimension"
    )
    dim = int(get_dim())
    dim_ok = dim == s.embedding_dim
    print(f"  向量维度           {dim}（配置期望 {s.embedding_dim}）{'✓' if dim_ok else '✗'}")

    samples = [
        "外倾者倾向于从社交互动中获取心理能量",
        "躯体化因子偏高可能表现为头痛与肌肉酸痛",
        "本结果仅供参考，不构成医学诊断",
    ]

    t0 = time.perf_counter()
    emb = model.encode(samples, normalize_embeddings=True)
    first_encode_s = time.perf_counter() - t0
    print(
        f"  首次 encode 3 条   {first_encode_s:.3f}s   ← 启动期必须由 warmup 烧掉"
    )

    t0 = time.perf_counter()
    model.encode(samples, normalize_embeddings=True)
    warm_encode_s = time.perf_counter() - t0
    print(f"  稳态 encode 3 条   {warm_encode_s:.3f}s   ← 验收判定用这个（要求 < 5s）")

    # 归一化：cosine 空间的前提
    try:
        import numpy as np

        norms = np.linalg.norm(np.asarray(emb), axis=1)
        norm_ok = bool(np.allclose(norms, 1.0, atol=1e-4))
        print(
            f"  L2 范数            {[round(float(n), 6) for n in norms]} "
            f"{'✓' if norm_ok else '✗ 未归一化，cosine 检索会失真'}"
        )
    except ImportError:
        norm_ok = True
        print("  （numpy 不可用，跳过归一化断言）")

    RESULTS.update(
        embedding_loaded=True,
        embedding_dim=dim,
        import_s=round(import_s, 3),
        load_s=round(load_s, 3),
        first_encode_s=round(first_encode_s, 3),
        warm_encode_s=round(warm_encode_s, 3),
    )

    ok = dim_ok and norm_ok and load_s < 3.0 and warm_encode_s < 5.0
    print(f"  → 加载 {load_s:.3f}s / 稳态编码 {warm_encode_s:.3f}s  {'✓' if ok else '✗'}")
    return ok


def check_chroma() -> bool:
    """④ Chroma 持久化目录可写。"""
    _hr("4. Chroma 持久化目录")

    from app.core.config import get_settings

    s = get_settings()
    s.chroma_path.mkdir(parents=True, exist_ok=True)
    print(f"  目录  {s.chroma_path}")

    probe = s.chroma_path / ".write_probe"
    try:
        probe.write_text("ok", encoding="utf-8")
        probe.unlink()
        print("  写入探测           ✓")
    except OSError as exc:
        print(f"  写入探测           ✗ {exc}")
        RESULTS["chroma_writable"] = False
        return False

    try:
        import chromadb
        from chromadb.config import Settings as ChromaSettings

        client = chromadb.PersistentClient(
            path=str(s.chroma_path),
            settings=ChromaSettings(anonymized_telemetry=False),
        )
        col = client.get_or_create_collection(
            name=s.collection_name,
            metadata={"hnsw:space": "cosine"},
        )
        count = col.count()
        print(f"  集合               {s.collection_name}（{count} 块，hnsw:space=cosine）✓")
        RESULTS.update(chroma_writable=True, collection_count=count)
        return True
    except Exception as exc:  # noqa: BLE001
        print(f"  集合打开失败       ✗ {exc}")
        RESULTS["chroma_writable"] = False
        return False


def check_llm(skip: bool = False) -> bool:
    """⑤ DeepSeek 连通。"""
    _hr("5. DeepSeek 连通")

    from app.core.config import get_settings

    s = get_settings()

    if skip:
        print("  --skip-llm：跳过")
        RESULTS["llm_reachable"] = None
        return True

    if not s.llm_enabled:
        print("  ✗ 未配置 LLM_API_KEY（在 rag-service/.env 里填上后重跑）")
        RESULTS["llm_reachable"] = False
        return False

    print(f"  base_url  {s.llm_base_url}")
    print(f"  model     {s.llm_model}")

    try:
        import httpx

        t0 = time.perf_counter()
        resp = httpx.post(
            f"{s.llm_base_url.rstrip('/')}/chat/completions",
            headers={
                "Authorization": f"Bearer {s.llm_api_key}",
                "Content-Type": "application/json",
            },
            json={
                "model": s.llm_model,
                "messages": [{"role": "user", "content": "ping"}],
                "max_tokens": 1,
            },
            timeout=15.0,
        )
        elapsed = time.perf_counter() - t0
        ok = resp.status_code == 200
        print(f"  HTTP {resp.status_code}  {elapsed:.3f}s  {'✓' if ok else '✗'}")
        if not ok:
            print(f"  响应：{resp.text[:300]}")
        RESULTS.update(llm_reachable=ok, llm_latency_s=round(elapsed, 3))
        return ok
    except Exception as exc:  # noqa: BLE001
        print(f"  ✗ {exc}")
        RESULTS["llm_reachable"] = False
        return False


def check_schemas() -> bool:
    """附加：契约模型自检（不需要网络与模型）。"""
    _hr("6. 契约模型（附加）")

    from app.api.schemas import EnhanceRequest, Measure

    sample = {
        "schema_version": 1,
        "test": {
            "code": 13,
            "key": "mbti",
            "name": "16型人格测试（MBTI）",
            "domains": ["typology", "common"],
            "policy": "full",
        },
        "result": {"level": "完成测试", "total": 41, "total_max": 70},
        "measures": [
            {
                "key": "EI",
                "label": "外倾(E)/内倾(I)",
                "kind": "dimension",
                "score": 15,
                "max": 21,
                "threshold": 12,
                "pole": "E",
                "strength": 0.43,
            },
            {
                "key": "somatization",
                "label": "躯体化",
                "kind": "factor",
                "score": 2.1,
                "max": 5,
                "clinical_cutoff": 2.0,
                "factor_level": "偏高",
                "strength": 0.03,
            },
        ],
        "labels": {"type": "ENFP"},
        "flags": ["needs_referral"],
        "options": {
            "top_k": 5,
            "max_suggestions": 5,
            "language": "zh-CN",
            "include_citations": True,
        },
    }

    checks: list[tuple[str, bool]] = []

    try:
        req = EnhanceRequest.model_validate(sample)
        checks.append(("§5.2 样例请求通过校验", len(req.measures) == 2))
        checks.append(("needs_referral 可读", req.needs_referral is True))
    except Exception as exc:  # noqa: BLE001
        print(f"  ✗ 样例请求校验失败：{exc}")
        RESULTS["schema_ok"] = False
        return False

    # 未知 schema_version 必须显式报错，不得静默兼容
    bad = {**sample, "schema_version": 99}
    try:
        EnhanceRequest.model_validate(bad)
        checks.append(("未知 schema_version 被拒绝", False))
    except Exception:  # noqa: BLE001
        checks.append(("未知 schema_version 被拒绝", True))

    # strength 越界必须拒绝
    try:
        Measure(
            key="EI", label="x", kind="dimension",
            score=1, threshold=12, pole="E", strength=1.5,
        )
        checks.append(("strength 越界被拒绝", False))
    except Exception:  # noqa: BLE001
        checks.append(("strength 越界被拒绝", True))

    # dimension 缺 threshold 必须拒绝（Go 侧义务）
    try:
        Measure(key="EI", label="x", kind="dimension", score=1, pole="E")
        checks.append(("dimension 缺 threshold 被拒绝", False))
    except Exception:  # noqa: BLE001
        checks.append(("dimension 缺 threshold 被拒绝", True))

    for name, ok in checks:
        print(f"  {'✓' if ok else '✗'} {name}")

    RESULTS["schema_ok"] = all(ok for _, ok in checks)
    return RESULTS["schema_ok"]  # type: ignore[return-value]


def main() -> int:
    parser = argparse.ArgumentParser(description="rag-service P0-1 环境自检")
    parser.add_argument("--skip-llm", action="store_true", help="不探测 DeepSeek（无 Key 时）")
    parser.add_argument("--skip-load", action="store_true", help="不加载模型，只查路径")
    parser.add_argument("--json", action="store_true", help="只输出 JSON 结果")
    args = parser.parse_args()

    print("=" * 68)
    print("  rag-service 环境自检（P0-1）")
    print("=" * 68)

    verdicts = {
        "embedding": check_embedding(skip_load=args.skip_load),
        "chroma": check_chroma(),
        "schemas": check_schemas(),
        "llm": check_llm(skip=args.skip_llm),
    }

    _hr("汇总")
    for key, val in RESULTS.items():
        print(f"  {key:22} {val}")
    print()
    for name, ok in verdicts.items():
        print(f"  {'✓' if ok else '✗'} {name}")

    all_ok = all(verdicts.values())
    print(f"\n  {'✅ 全部通过' if all_ok else '❌ 存在失败项'}")

    if args.json:
        print(json.dumps({"verdicts": verdicts, "results": RESULTS}, ensure_ascii=False, indent=2))

    return 0 if all_ok else 1


if __name__ == "__main__":
    raise SystemExit(main())
