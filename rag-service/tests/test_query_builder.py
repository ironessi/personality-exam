"""`app.rag.query_builder` 的口径测试。

这里锁的不是「query 长什么样」这种易碎的字面量，而是三条**有实测依据**的口径：

1. per-measure query **不注入测评名**（``16型人格测试（MBTI）`` 会把
   ``common`` 里的「测评工具总览」hub 块抬到 top-1，实测见模块 docstring）；
2. ``label`` 里的 ``" / "`` 归一成 ``" 与 "``；
3. 只有**没有 measure 时的兜底 query** 才允许出现测评名。
"""

from __future__ import annotations

import pytest

from app.api.schemas import EnhanceRequest
from app.rag.query_builder import (
    build_queries,
    measure_is_high,
    retrieval_namespaces,
)

TEST_NAME = "16型人格测试（MBTI）"


def _req(measures: list[dict], *, flags: list[str] | None = None) -> EnhanceRequest:
    return EnhanceRequest(
        schema_version=1,
        test={
            "code": 13,
            "key": "mbti",
            "name": TEST_NAME,
            "domains": ["typology", "common"],
            "policy": "full",
        },
        result={"level": "完成测试", "total": None, "total_max": None},
        measures=measures,
        flags=flags or [],
    )


def _dim(key: str, label: str, score: float = 15.0, pole: str = "E") -> dict:
    return {
        "key": key,
        "label": label,
        "kind": "dimension",
        "score": score,
        "max": 21.0,
        "threshold": 12.0,
        "pole": pole,
        "strength": 0.333,
    }


def _factor(key: str, label: str, score: float = 3.1) -> dict:
    return {
        "key": key,
        "label": label,
        "kind": "factor",
        "score": score,
        "clinical_cutoff": 2.0,
        "factor_level": "偏高",
        "strength": 0.55,
    }


def test_dimension_query_does_not_inject_test_name() -> None:
    """回归：测评名会把 common 的 hub 块抬到 top-1（bge-small-zh 分不开）。"""
    queries = build_queries(_req([_dim("e", "外倾 / 内倾")]))
    assert len(queries) == 1
    assert TEST_NAME not in queries[0]
    # 测评名里的「MBTI」「人格」「测试」任一被注入都算回退
    assert "MBTI" not in queries[0]
    assert "人格" not in queries[0]


def test_factor_query_does_not_inject_test_name() -> None:
    queries = build_queries(_req([_factor("som", "躯体化")]))
    assert TEST_NAME not in queries[0]
    assert "因子" in queries[0]


def test_slash_label_is_normalized_to_conjunction() -> None:
    queries = build_queries(_req([_dim("e", "外倾 / 内倾")]))
    assert "外倾 与 内倾" in queries[0]
    assert "/" not in queries[0]
    assert "／" not in queries[0]


def test_full_width_slash_and_packed_slash_are_normalized() -> None:
    packed = build_queries(_req([_dim("e", "外倾/内倾")]))[0]
    fullwidth = build_queries(_req([_dim("e", "外倾／内倾")]))[0]
    spaced = build_queries(_req([_dim("e", "外倾 / 内倾")]))[0]
    assert packed == fullwidth == spaced


def test_one_query_per_measure_covering_every_label() -> None:
    """coverage 红线要求每个 measure 都被检索到，不能只挑高的写。"""
    req = _req(
        [
            _dim("e", "外倾 / 内倾"),
            _dim("s", "实感 / 直觉", score=6.0, pole="S"),
            _dim("t", "思考 / 情感", score=11.0, pole="T"),
            _dim("j", "判断 / 知觉", score=9.0, pole="J"),
        ]
    )
    queries = build_queries(req)
    assert len(queries) == 4
    for label in ("外倾", "实感", "思考", "判断"):
        assert any(label in q for q in queries), label


def test_needs_referral_appends_one_referral_query() -> None:
    req = _req([_factor("som", "躯体化")], flags=["needs_referral"])
    queries = build_queries(req)
    assert len(queries) == 2
    assert queries[-1] == "何时需要寻求专业心理帮助"


def test_no_referral_flag_means_no_referral_query() -> None:
    queries = build_queries(_req([_factor("som", "躯体化")]))
    assert all("专业心理帮助" not in q for q in queries)


def test_fallback_query_may_use_test_name() -> None:
    """没有 measure 时只能整体问一句——这是唯一允许出现测评名的地方。

    ``measuers`` 在 schema 层有 ``min_length=1``，所以绕开校验器直接构造一个
    只有 ``test`` 的替身，专测兜底分支。
    """

    class _Bare:
        class test:  # noqa: N801 - 模拟 TestInfo
            name = TEST_NAME

        measures: list = []
        needs_referral = False

    assert build_queries(_Bare()) == [f"{TEST_NAME} 结果 解读"]


@pytest.mark.parametrize(
    ("strength", "expected"),
    [
        (0.5, True),
        (-0.5, True),
        (0.49, False),
        (-0.49, False),
        (0.0, False),
        (None, False),
    ],
)
def test_measure_is_high_reads_numeric_strength_only(
    strength: float | None, expected: bool
) -> None:
    class _M:
        pass

    m = _M()
    m.strength = strength
    m.factor_level = "偏高"  # 字符串不参与判断
    assert measure_is_high(m) is expected


def test_retrieval_namespaces_are_taken_verbatim_from_contract() -> None:
    """Python 不判断测评类型——domains 是 Go 算好的。"""
    req = _req([_dim("e", "外倾 / 内倾")])
    assert retrieval_namespaces(req) == ["typology", "common"]
