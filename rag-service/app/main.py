"""FastAPI 入口。

启动顺序（lifespan）都围绕一件事：**在开始接客之前把首次前向的开销烧掉**。
未预热的服务在冷启动后第一个请求上必然撑爆 ADR-0007 的 3 秒同步期限。
"""

from __future__ import annotations

import time
from contextlib import asynccontextmanager

from fastapi import FastAPI

from app import __version__
from app.api.routes import router
from app.core.config import get_settings
from app.core.logging import configure_logging, get_logger

logger = get_logger(__name__)


@asynccontextmanager
async def lifespan(app: FastAPI):
    s = get_settings()
    configure_logging(s.log_level)

    t_boot = time.perf_counter()
    logger.info("rag-service v%s 启动中 …", __version__)

    # 1. 目录（P0-1 验收项之一：Chroma 持久化目录可写）
    s.ensure_dirs()

    # 2. Chroma —— 失败也让进程起来，/health 会如实报 degraded
    try:
        from app.rag.index import collection_count, health_ok

        app.state.chroma_ok = health_ok()
        logger.info("Chroma 就绪 | 集合块数=%d", collection_count())
    except Exception as exc:  # noqa: BLE001
        app.state.chroma_ok = False
        logger.error("Chroma 初始化失败（服务仍启动，/health 会报 degraded）：%s", exc)

    # 3. 向量模型 + 预热 —— 这一步故意同步阻塞，且必须做完
    if s.warmup_on_startup:
        try:
            from app.rag.embedding import get_embedding_model

            model = get_embedding_model()
            first_encode_s = model.warmup()
            # 3b. 按 token 长度暖机 —— 这是**另一个**一次性开销，`warmup()` 烧不掉它。
            #     实测每个新长度的首次前向要 200~2400ms，30 条 fixture 顺序跑会把
            #     检索 P95 从 19ms 抬到 837ms（中位数仍是 30ms，所以只有 P95 能看出问题）。
            #     详见 `app/rag/embedding.py::warmup_lengths` 的实测记录。
            lengths_s, lengths_n = model.warmup_lengths(s.warmup_length_lo, s.warmup_length_hi)
            app.state.warmup_lengths_s = round(lengths_s, 3)
            app.state.warmup_lengths_n = lengths_n
            app.state.embedding_loaded = model.loaded
            app.state.warmup_first_encode_s = round(first_encode_s, 3)
        except Exception as exc:  # noqa: BLE001
            app.state.embedding_loaded = False
            app.state.warmup_first_encode_s = None
            logger.error("向量模型预热失败（服务仍启动，/enhance 不可用）：%s", exc)
    else:
        app.state.embedding_loaded = False
        app.state.warmup_first_encode_s = None
        app.state.warmup_lengths_s = None
        app.state.warmup_lengths_n = 0
        logger.warning("WARMUP_ON_STARTUP=false —— 第一个请求会承担 3~4s 首次前向开销")

    # 4. 检索栈预热 —— 三笔独立的一次性开销，少预热任何一笔都会落在第一个真实请求上
    #    · LlamaIndex 的 embedding 适配器：首次调用比底层单例多付约 1.2s
    #    · VectorStoreIndex 装配 + 首次 Chroma 查询：冷启动合计约 5.5s
    #    · BM25 索引构建（jieba 分词整个语料）：约 0.4s/域
    app.state.retrieval_warmup_s = None
    if s.warmup_on_startup and app.state.chroma_ok:
        try:
            from app.rag.index import collection_count
            from app.rag.retriever import retrieve

            if collection_count() > 0:
                t_retrieval = time.perf_counter()
                retrieve(
                    ["预热查询"],
                    namespaces=list(s.kb_namespaces),
                    top_k=3,
                    mode="hybrid",
                )
                app.state.retrieval_warmup_s = round(time.perf_counter() - t_retrieval, 3)
                logger.info("检索栈预热完成 | 耗时=%.3fs", app.state.retrieval_warmup_s)
            else:
                logger.warning("知识库为空，跳过检索预热（先跑 scripts/ingest_cli.py）")
        except Exception as exc:  # noqa: BLE001
            logger.error("检索栈预热失败（服务仍启动，检索会承担冷启动开销）：%s", exc)

    # 5. Stage 2（NLI 语义蕴含）预热 —— 必须显式做，且必须把结果记进日志。
    #
    #    verifier 的 ``_NliModel.get()`` 是懒加载：不预热就会把 3–5s 的首次
    #    前向摊到第一个真实请求上（而请求预算只有 30s）。更糟的是**权重缺失时
    #    它只会静默退化成单级验证**，所以这里无论成败都要打一条日志，
    #    并把结论存进 app.state 供 /health 读取。
    app.state.stage2_ready = False
    app.state.stage2_error = None
    if s.warmup_on_startup:
        try:
            from app.core.verifier import stage2_status

            t_nli = time.perf_counter()
            ready, err = stage2_status()
            app.state.stage2_ready = ready
            app.state.stage2_error = err
            if ready:
                logger.info(
                    "Stage 2（NLI 语义蕴含）就绪 | 预热耗时=%.3fs", time.perf_counter() - t_nli
                )
            else:
                logger.warning(
                    "Stage 2（NLI 语义蕴含）不可用，将只跑字面级验证：%s。"
                    "「抄对句子挂错结论」不会被拦下——生产必须装好 NLI 权重。",
                    err,
                )
        except Exception as exc:  # noqa: BLE001
            app.state.stage2_error = str(exc)
            logger.error("Stage 2 预热异常（服务仍启动，验证退化为单级）：%s", exc)

    app.state.boot_seconds = round(time.perf_counter() - t_boot, 3)
    logger.info(
        "启动完成 | 总耗时=%.3fs | embedding_loaded=%s | llm_configured=%s",
        app.state.boot_seconds,
        app.state.embedding_loaded,
        s.llm_enabled,
    )

    yield

    logger.info("rag-service 退出")


settings = get_settings()

app = FastAPI(
    title="心理测评 RAG 增强服务",
    version=__version__,
    description=(
        "为心理测评结果生成深度解读。契约见 doc/RAG增强系统开发文档.md §5.2。\n\n"
        "注意：这是**唯一**的 Go↔Python 通道（同步 HTTP）。Redis 队列夹在 Go 的"
        "提交接口与 Go 的 worker 之间，不是 Go 与 Python 之间的通道。"
    ),
    lifespan=lifespan,
)

app.include_router(router)


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(
        "app.main:app",
        host=settings.host,
        port=settings.port,
        reload=False,
    )
