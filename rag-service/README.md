# rag-service — 心理测评 RAG 增强服务

为 GoFrame 心理测评平台提供**深度解读**能力：把测评得分翻译成有文献依据的个性化
解读与建议，且每条结论都可溯源。

架构、契约与验收标准见 [`../doc/RAG增强系统开发文档.md`](../doc/RAG增强系统开发文档.md)，
关键决策见 [`../doc/adr/`](../doc/adr/)，术语见 [`../doc/GLOSSARY.md`](../doc/GLOSSARY.md)。

---

## 当前进度

| 阶段 | 内容 | 状态 |
|---|---|---|
| **P0-1** | 环境与骨架 | 🔨 进行中（目录、配置、`/health`、启动预热、自检脚本已就绪） |
| P0-2 | 知识库入库 | ⬜ 未开始 —— `app/ingest/*` |
| P0-3 | 混合检索 | ⬜ 未开始 —— `app/rag/{query_builder,retriever}.py` |
| P0-4 | 生成 + 两级引用验证 | ⬜ 未开始 —— `app/rag/synthesizer.py`、`app/core/{verifier,filter}.py` |
| P0-5 | Go 侧集成 | ⬜ 未开始 —— `internal/logic/rag/*` |
| P0-6 | 泛化度量 | ⬜ 未开始 —— ADR-0006 |

`/enhance` 与 `/admin/reload` 当前显式返回 **501**。这是刻意的：返回假数据会被下游
当成真结果写进用户的 `result_data.rag`。

---

## 环境

用**已有的 `/Users/iron/rag` 虚拟环境**（uv 建的，Python 3.12.7，已含 torch / chromadb /
sentence-transformers 等）。

```bash
# 补齐缺失依赖
/Users/iron/rag/bin/python -m pip install -r requirements.txt
# 或从 uv 重装
# uv pip install --python /Users/iron/rag/bin/python -r requirements.txt
```

配置：

```bash
cp .env.example .env
# 然后编辑 .env，至少填上 LLM_API_KEY
```

> ⚠️ `.env` 已被 `.gitignore` 忽略。API Key 绝不能入库。

---

## 自检（P0-1 验收）

```bash
/Users/iron/rag/bin/python scripts/smoke_test.py
# 暂时没有 Key 时
/Users/iron/rag/bin/python scripts/smoke_test.py --skip-llm
```

脚本按 §9 的四项验收逐条报结果，并把计时**拆成三段**报告：

| 阶段 | 实测 | 说明 |
|---|---|---|
| `import torch` + `sentence_transformers` | ~5.0 s | **不计入**「加载 < 3s」判定 |
| `SentenceTransformer(path)` | 0.21 s 冷 / 0.08 s 热 | 验收判定用这个 |
| 加载后**首次** `encode(3 条)` | **~4.3 s** | 启动期由 `warmup()` 烧掉 |
| 稳态 `encode(3 条)` | ~0.01 s | 这才是 §11「Embedding < 0.5s」所指 |

---

## 运行

```bash
/Users/iron/rag/bin/python -m uvicorn app.main:app --host 127.0.0.1 --port 8001
```

启动时会**同步**加载向量模型并跑 2 轮 dummy encode。这是必须的：ADR-0007 的同步
快路径只有 3 秒，而首次前向要 3~4 秒——不预热的话，每次服务重启都会让第一个用户
白等 3 秒再被转成异步。日志里会打印 `首次 encode N 条=x.xxxs（本该由第一个用户承担）`。

```bash
curl -s localhost:8001/api/v1/health | python -m json.tool
# {"status":"ok","embedding_loaded":true,"collection_count":0,"llm_reachable":true}
```

---

## 测试

```bash
cd rag-service
/Users/iron/rag/bin/python -m pytest tests/ -v
```

骨架阶段只有 `/health` 的形状与契约校验用例在跑；其余用例已写好并标了
`skip(reason="P0-x 实现后启用")`，实现对应模块后去掉 skip 即可。

---

## 目录

```
app/
├── main.py              FastAPI 入口 + lifespan 预热
├── api/
│   ├── routes.py        /health、/enhance(501)、/admin/reload(501)
│   └── schemas.py       §5.2 规范契约（ADR-0001）
├── core/
│   ├── config.py        Settings（pydantic-settings）
│   ├── logging.py       统一日志
│   ├── verifier.py      两级引用验证（ADR-0003）
│   └── filter.py        敏感词过滤（ADR-0004）
├── rag/
│   ├── embedding.py     本地 BGE 单例 + 预热 + LlamaIndex BaseEmbedding 适配器
│   ├── index.py         Chroma 连接、集合与 ChromaVectorStore 装配
│   ├── bm25.py          BM25 双后端（jieba 词级 / LlamaIndex 汉字级），每域独立
│   ├── retriever.py     混合检索（LlamaIndex VectorIndexRetriever + RRF，ADR-0005/0009）
│   ├── query_builder.py 结果 → 检索 query（§8.1 的三处陷阱）
│   └── synthesizer.py   DeepSeek + 结构化输出
└── ingest/              loader / splitter / pipeline（P0-2）
data/
├── raw/<构念域>/      语料唯一真源（**入库**）
│                      symptom/trait/typology/social/wellbeing/career/ability/common
├── processed/                      切块产物（可重建，不入库）
└── chroma/                         向量库（可重建，不入库）
```

---

## 三件容易搞错的事

1. **`kb_namespace` 是构念域，三处必须一致**：Chroma `where` 过滤、BM25 索引选择、
   `data/raw/` 子目录名。取值只有 8 个（`symptom`/`trait`/`typology`/`social`/
   `wellbeing`/`career`/`ability`/`common`），**不是测评名**——按测评建 namespace 会
   得到 42 个几乎同质的索引，而 `chunk_id` 前缀就是 namespace，两套命名混用时
   `upsert` 会相互覆盖且**不报错**（ADR-0008）。BM25 **没有过滤能力**，必须按
   namespace 各建独立索引；混用一个索引会静默串味（类型学的 query 召回症状严重度的
   块，不报错，只是引用「有点对但不对题」）。

2. **Redis 不是 Go↔Python 的通道**。`rag:queue` 夹在 Go 的提交接口与 Go 的 worker
   之间，解耦的是 Go **内部**的「提交」与「调 LLM」。Go 与 Python 之间只有一条同步
   HTTP：`POST /api/v1/enhance`。

3. **`status` ≠ `rag_status`**。前者是 `/health` 的服务健康（`ok`/`degraded`），后者是
   测评结果里的增强状态（`pending`/`ready`/`degraded`/`failed`）。
