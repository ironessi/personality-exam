"""HTTP 路由。

P0-1 阶段只有 ``/health`` 是真实实现；``/enhance`` 与 ``/admin/reload`` 在
P0-4 落地。``/enhance`` 串起 query_builder → retrieve → synthesize → verify。
"""

from __future__ import annotations

import time

from fastapi import APIRouter, Header, HTTPException, Request, status

from app.api.schemas import EnhanceData, EnhanceRequest, EnhanceResponse, HealthResponse
from app.core.config import get_settings
from app.core.logging import get_logger

logger = get_logger(__name__)

router = APIRouter(prefix="/api/v1", tags=["rag"])

# /health 里 llm_reachable 的缓存：别让健康检查把 DeepSeek 打爆
_LLM_PROBE_TTL_S = 60.0
_llm_probe_cache: dict[str, float] = {"at": 0.0, "ok": 0.0}

#: Stage 2 不可用的 WARNING 只报一次（每进程）。
_stage2_warned: dict[str, bool] = {"done": False}


def _verify_caller(provided: str | None) -> None:
    """内网留空即不校验；一旦配了 service_api_key 就必须带对。"""
    expected = get_settings().service_api_key.strip()
    if not expected:
        return
    if provided != expected:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="X-API-Key 缺失或不正确",
        )


def _llm_reachable() -> bool:
    """探一次 DeepSeek 的 /models。结果缓存 60 s，失败不抛错。"""
    s = get_settings()
    if not s.llm_enabled:
        return False

    now = time.monotonic()
    if now - _llm_probe_cache["at"] < _LLM_PROBE_TTL_S:
        return bool(_llm_probe_cache["ok"])

    ok = False
    try:
        import httpx

        resp = httpx.get(
            f"{s.llm_base_url.rstrip('/')}/models",
            headers={"Authorization": f"Bearer {s.llm_api_key}"},
            timeout=5.0,
        )
        ok = resp.status_code == 200
        if not ok:
            logger.warning("LLM 探活返回 %s：%s", resp.status_code, resp.text[:200])
    except Exception as exc:  # noqa: BLE001
        logger.warning("LLM 探活失败：%s", exc)

    _llm_probe_cache.update(at=now, ok=1.0 if ok else 0.0)
    return ok


@router.get("/health", response_model=HealthResponse)
def health(request: Request) -> HealthResponse:
    """P0-1 自检端点。

    响应形状固定为 ``{status, embedding_loaded, collection_count, llm_reachable}``
    （§5.2 :343）。注意 ``status`` 是 ``"ok"``/``"degraded"``，与业务字段
    ``rag_status``（``pending``/``ready``/``degraded``/``failed``）**不是一回事**。
    """
    embedding_loaded = bool(getattr(request.app.state, "embedding_loaded", False))
    chroma_ok = bool(getattr(request.app.state, "chroma_ok", False))

    if chroma_ok:
        try:
            from app.rag.index import collection_count

            count = collection_count()
        except Exception:  # noqa: BLE001
            count = 0
    else:
        count = 0

    healthy = embedding_loaded and chroma_ok
    return HealthResponse(
        status="ok" if healthy else "degraded",
        embedding_loaded=embedding_loaded,
        collection_count=count,
        llm_reachable=_llm_reachable(),
        stage2_ready=bool(getattr(request.app.state, "stage2_ready", False)),
    )


async def _enhance_async(req: EnhanceRequest) -> EnhanceData:
    """async 包装：synthesize 是 async，verify 之后是同步。"""
    from app.core.verifier import (
        stage2_status,
        support_coverage,
        support_precision,
        verify_citations,
        verify_suggestions,
    )
    from app.rag.query_builder import build_queries, retrieval_namespaces
    from app.rag.retriever import retrieve
    from app.rag.synthesizer import synthesize

    queries = build_queries(req)
    namespaces = retrieval_namespaces(req)
    chunks = retrieve(queries, namespaces=namespaces)

    if not chunks:
        return EnhanceData(
            enhanced_text="",
            sections=[],
            suggestions=[],
            citations=[],
            rag_status="failed",
        )

    data = await synthesize(req, chunks)

    # Stage 2 不可用时两级验证会静默退化成一字面级（见 verifier.stage2_status）。
    # 只在首次发现时大声报警一次，避免每条日志都刷屏，也避免没人注意到。
    if not _stage2_warned["done"]:
        ready, err = stage2_status()
        if not ready:
            _stage2_warned["done"] = True
            logger.warning(
                "Stage 2（NLI 语义蕴含）不可用，当前仅跑字面级验证：%s。"
                "这意味着「抄对句子挂错结论」不会被拦下——生产环境必须装好 NLI 权重。",
                err,
            )

    chunks_by_id = {c.chunk_id: c.text for c in chunks}
    passed, details = verify_citations(data.citations, chunks_by_id)

    total = len(details)
    n_passed = len(passed)
    precision, low, high = support_precision(n_passed, total)

    data.citations = passed

    # 建议走同一套两级验证，把 basis 当 claim（ADR-0003 补充）。没通过的直接丢弃：
    # 留着它们等于把「看起来有据」的建议发给用户，而前端并不会展示 verified。
    sugg_passed, sugg_details = verify_suggestions(data.suggestions, chunks_by_id)
    data.suggestions = sugg_passed
    data.rag_status = "ready"

    # ADR-0003 要求精确率与覆盖率**成对上报**：单看精确率时，「每节只写一句
    # 最保险的话」能刷到 100%，而报告同时退化成正确的废话。
    coverage = support_coverage(passed, len(data.sections))

    logger.info(
        "verify: passed=%d/%d precision=%.3f wilson=[%.3f,%.3f] coverage=%.3f "
        "suggestions=%d/%d",
        n_passed,
        total,
        precision,
        low,
        high,
        coverage,
        len(sugg_passed),
        len(sugg_details),
    )

    if total > 0 and n_passed == 0:
        data.rag_status = "degraded"

    return data


@router.post("/enhance", response_model=EnhanceResponse)
async def enhance(
    req: EnhanceRequest, x_api_key: str | None = Header(default=None)
) -> EnhanceResponse:
    """生成增强解读（P0-4）。

    串起 query_builder → retrieve → synthesize → verify。任何一步失败都降级为
    ``success=False`` 并带 error，而不是返回假数据——假数据会被下游当成真结果
    写进用户的 ``result_data.rag``（ADR-0002 的失败降级）。
    """
    _verify_caller(x_api_key)

    # policy=disabled 的测评本就不该调用本接口；防御性拦一道。
    if req.test.policy == "disabled":
        raise HTTPException(
            status_code=status.HTTP_422_UNPROCESSABLE_ENTITY,
            detail="该测评的增强策略是 disabled，不生成个性化解读（ADR-0008）",
        )

    try:
        data = await _enhance_async(req)
        return EnhanceResponse(success=True, data=data)
    except Exception as exc:  # noqa: BLE001 —— 失败降级，不把 500 抛给 Go
        logger.exception("增强生成失败：%s", exc)
        return EnhanceResponse(
            success=False,
            data=None,
            error="增强生成失败，已降级（不影响原始测评结果）",
        )


@router.post("/admin/reload")
def admin_reload(x_api_key: str | None = Header(default=None)) -> dict[str, int | str]:
    """重载知识库（P0-2 实现）。

    §6.1 的 ``apiKey`` 默认为空，所以**鉴权靠 service_api_key**；暴露端口时
    必须设置它，否则任何人都能触发重建索引。
    """
    _verify_caller(x_api_key)
    try:
        from app.ingest.pipeline import reload_indexes

        return reload_indexes()
    except Exception as exc:  # noqa: BLE001
        logger.exception("知识库重载失败：%s", exc)
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"知识库重载失败：{exc}",
        )
