"""HTTP 接口测试（P0-1 起即可跑）。

``/health`` 的响应形状必须固定为 §5.2 :343 的四字段。注意 ``status``
（``"ok"``/``"degraded"``）与业务字段 ``rag_status``
（``pending``/``ready``/``degraded``/``failed``）**不是一回事**——后者是
测评结果里的增强状态，前者是服务健康。
"""

from __future__ import annotations

import pytest
from fastapi.testclient import TestClient
from pydantic import ValidationError

from app.main import app


@pytest.fixture(scope="module")
def client():
    # `with` 会真的跑 lifespan（加载模型 + 预热），所以首个用例约 10s、且依赖
    # data/chroma 与本地模型权重存在。若想跑得快，可去掉 `with`——此时
    # app.state 上不会有 embedding_loaded，/health 会如实返回 degraded，
    # 形状断言仍然通过。
    with TestClient(app) as c:
        yield c


def test_health_shape(client):
    resp = client.get("/api/v1/health")
    assert resp.status_code == 200

    body = resp.json()
    # stage2_ready 是 2026-10-21 新增的**可选**字段（默认 False）：Stage 2 缺权重时
    # 服务仍能跑，但两级验证实际退化，运维必须看得见，又不该让 status 变 degraded
    # 把 Go 侧误判成"整个增强不可用"。
    assert set(body) == {
        "status",
        "embedding_loaded",
        "collection_count",
        "llm_reachable",
        "stage2_ready",
    }
    assert body["status"] in {"ok", "degraded"}
    assert isinstance(body["embedding_loaded"], bool)
    assert isinstance(body["collection_count"], int)
    assert isinstance(body["llm_reachable"], bool)
    assert isinstance(body["stage2_ready"], bool)


def test_enhance_rejects_unknown_schema_version(client):
    """未知 schema_version 必须被拒，不得静默兼容（ADR-0001）。"""
    payload = {
        "schema_version": 99,
        "test": {
            "code": 13,
            "key": "mbti",
            "name": "MBTI",
            "domains": ["typology", "common"],
            "policy": "full",
        },
        "result": {"level": "完成测试"},
        "measures": [
            {
                "key": "EI",
                "label": "外倾/内倾",
                "kind": "dimension",
                "score": 15,
                "threshold": 12,
                "pole": "E",
            }
        ],
    }
    assert client.post("/api/v1/enhance", json=payload).status_code == 422


def test_enhance_rejects_strength_out_of_range(client):
    payload = {
        "schema_version": 1,
        "test": {
            "code": 13,
            "key": "mbti",
            "name": "MBTI",
            "domains": ["typology", "common"],
            "policy": "full",
        },
        "result": {"level": "完成测试"},
        "measures": [
            {
                "key": "EI",
                "label": "外倾/内倾",
                "kind": "dimension",
                "score": 15,
                "threshold": 12,
                "pole": "E",
                "strength": 1.5,
            }
        ],
    }
    assert client.post("/api/v1/enhance", json=payload).status_code == 422


def test_enhance_returns_structured_json_for_whitelisted_tests(client, monkeypatch):
    """§9 —— 对 §6.1 白名单内的测评均返回结构化 JSON。

    用 monkeypatch 替换 ``synthesize``，避免真实 LLM 依赖，同时仍覆盖
    query_builder → retrieve → verify 的完整路由串联。
    """
    from app.api.schemas import EnhanceData

    async def fake_synthesize(req):
        return EnhanceData(
            enhanced_text="测试解读",
            sections=[{"title": "概述", "content": "内容"}],
            suggestions=[{"text": "建议", "basis": "test"}],
            citations=[],
            model="deepseek-flash",
            prompt_tokens=1,
            completion_tokens=1,
            latency_ms=1,
            rag_status="ready",
        )

    import app.api.routes as routes

    monkeypatch.setattr(routes, "_enhance_async", fake_synthesize)

    payload = {
        "schema_version": 1,
        "test": {
            "code": 13,
            "key": "mbti",
            "name": "16型人格测试（MBTI）",
            "domains": ["typology", "common"],
            "policy": "full",
        },
        "result": {"level": "完成测试"},
        "measures": [
            {
                "key": "EI",
                "label": "外倾/内倾",
                "kind": "dimension",
                "score": 15,
                "threshold": 12,
                "pole": "E",
            }
        ],
    }
    resp = client.post("/api/v1/enhance", json=payload)
    assert resp.status_code == 200
    body = resp.json()
    assert body["success"] is True
    assert body["data"]["rag_status"] == "ready"
    assert body["data"]["enhanced_text"] == "测试解读"


def test_enhance_rejects_disabled_policy(client):
    """policy=disabled 的测评不该走增强，路由层防御性拦 422。"""
    payload = {
        "schema_version": 1,
        "test": {
            "code": 19,
            "key": "zootopia",
            "name": "疯狂动物城人格",
            "domains": ["common"],
            "policy": "disabled",
        },
        "result": {"level": "完成测试"},
        "measures": [
            {
                "key": "X",
                "label": "某维度",
                "kind": "dimension",
                "score": 1,
                "threshold": 1,
                "pole": "A",
            }
        ],
    }
    resp = client.post("/api/v1/enhance", json=payload)
    assert resp.status_code == 422


def test_admin_reload_requires_key_when_configured(client, monkeypatch):
    """key 未配时放行（调 reload_indexes），配了就必须校验 X-API-Key。"""
    from app.core.config import get_settings

    # 未配 key → 放行（走 reload_indexes，可能因 Chroma 空而返回计数）
    monkeypatch.setattr(get_settings(), "service_api_key", "")
    resp = client.post("/api/v1/admin/reload")
    # 未配 key 时不应 401；成功或业务 500 都算「鉴权已放行」
    assert resp.status_code != 401

    # 配了 key → 缺 X-API-Key 必须 401
    monkeypatch.setattr(get_settings(), "service_api_key", "secret-key-123")
    resp = client.post("/api/v1/admin/reload")
    assert resp.status_code == 401


# ── 构念域与增强策略（ADR-0008）──────────────────────────────────


def _request(**test_overrides):
    """构造一个合法请求；``test_overrides`` 用来注入非法字段做负例。"""
    from app.api.schemas import EnhanceRequest

    test = {
        "code": 13,
        "key": "mbti",
        "name": "16型人格测试（MBTI）",
        "domains": ["typology", "common"],
        "policy": "full",
    }
    test.update(test_overrides)
    return EnhanceRequest(
        schema_version=1,
        test=test,
        result={"level": "完成测试"},
        measures=[
            {
                "key": "EI",
                "label": "外倾/内倾",
                "kind": "dimension",
                "score": 15,
                "threshold": 12,
                "pole": "E",
            }
        ],
    )


def test_domains_reject_unknown_namespace():
    """按测评名建 namespace 是 ADR-0008 明确取代掉的旧做法，必须报错而不是静默放过。"""
    with pytest.raises(ValidationError):
        _request(domains=["mbti", "common"])


def test_domains_reject_duplicates():
    with pytest.raises(ValidationError):
        _request(domains=["common", "common"])


def test_policy_common_only_must_retrieve_common_alone():
    with pytest.raises(ValidationError):
        _request(domains=["typology", "common"], policy="common_only")


def test_policy_full_must_include_common():
    with pytest.raises(ValidationError):
        _request(domains=["typology"], policy="full")


def test_common_only_accepts_common_alone():
    req = _request(domains=["common"], policy="common_only")
    assert req.test.policy == "common_only"
    assert req.test.domains == ["common"]


def test_test_info_no_longer_has_kb_namespace():
    """旧字段必须真的消失，否则两套分区键会并存（GLOSSARY 第 4 条歧义）。"""
    from app.api.schemas import TestInfo

    assert "kb_namespace" not in TestInfo.model_fields


def test_enhance_async_drops_suggestions_without_evidence(monkeypatch):
    """建议没出处时必须被丢弃，而不是「留着但标 verified=False」。

    原实现只验 citations，suggestions 是整份报告里唯一没有出处的部分——前端
    也不展示 verified，留着等于照样把无据建议发给用户。
    """
    import asyncio

    from app.api import routes
    from app.api.schemas import (
        Citation,
        EnhanceData,
        EnhanceRequest,
        Section,
        Suggestion,
    )
    from app.rag.retriever import RetrievedChunk

    quote = "偏直觉的人更容易被模式、关联、可能性吸引。"
    chunk = RetrievedChunk(
        chunk_id="typology_00001",
        text=quote + "抽象概念、未来图景、隐喻对他们更好用。",
        source="16型人格类型详解.md",
        kb_namespace="typology",
        section_title="直觉倾向",
    )

    async def fake_synthesize(req, chunks):
        return EnhanceData(
            enhanced_text="总结",
            sections=[Section(title="直觉", content="内容")],
            suggestions=[
                Suggestion(
                    text="在需要核对细节时建立检查清单。",
                    basis="直觉倾向的人更容易被模式与可能性吸引",
                    chunk_id="typology_00001",
                    evidence_quote=quote,
                ),
                Suggestion(
                    text="每天冥想二十分钟。",
                    basis="这句话是模型编的",
                    chunk_id="typology_00001",
                    evidence_quote="原文里根本没有这一句。",
                ),
            ],
            citations=[
                Citation(
                    claim="直觉倾向的人更容易被模式与可能性吸引",
                    source="16型人格类型详解.md",
                    chunk_id="typology_00001",
                    evidence_quote=quote,
                    section=0,
                )
            ],
        )

    monkeypatch.setattr("app.rag.retriever.retrieve", lambda *a, **k: [chunk])
    monkeypatch.setattr("app.rag.synthesizer.synthesize", fake_synthesize)
    # Stage 2 的判定在本用例里无关（我们验的是「有没有走同一套验证」），
    # 且真实 NLI 不认同义改写，会让通过与否随模型版本漂移。
    monkeypatch.setattr("app.core.verifier._nli_entails", lambda **k: True)

    req = EnhanceRequest.model_validate(
        {
            "schema_version": 1,
            "test": {
                "code": 13,
                "key": "mbti",
                "name": "16型人格测试（MBTI）",
                "domains": ["typology", "common"],
                "policy": "full",
            },
            "result": {"level": "完成测试"},
            "measures": [
                {
                    "key": "SN",
                    "label": "感觉/直觉",
                    "kind": "dimension",
                    "score": 6,
                    "threshold": 12,
                    "pole": "N",
                }
            ],
        }
    )

    data = asyncio.run(routes._enhance_async(req))

    assert [s.text for s in data.suggestions] == ["在需要核对细节时建立检查清单。"]
    assert data.suggestions[0].verified is True
    assert data.citations[0].verified is True
    assert data.rag_status == "ready"
