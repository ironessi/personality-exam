"""生成侧测试：prompt 的确定性与「程度词」口径。

这里测的都是**不该交给模型判断**的东西：强弱程度、敏感词、JSON 兜底解析。
"""

from __future__ import annotations

import pytest

from app.api.schemas import EnhanceRequest, Measure, TestInfo
from app.rag.synthesizer import (
    _measures_lines,
    _parse_llm_json,
    build_prompt,
    strength_band,
)


class TestStrengthBand:
    """程度词必须由代码定：实测模型会把最弱的一维说成最明显的一维。"""

    @pytest.mark.parametrize(
        ("strength", "expected"),
        [
            (0.0, "接近中间（不要强调任何一端的偏向）"),
            (0.19, "接近中间（不要强调任何一端的偏向）"),
            (-0.19, "接近中间（不要强调任何一端的偏向）"),
            (0.20, "轻度偏向"),
            (0.33, "轻度偏向"),
            (-0.49, "轻度偏向"),
            (0.50, "明显偏向"),
            (-0.67, "明显偏向"),
            (0.80, "非常明显"),
            (-1.0, "非常明显"),
        ],
    )
    def test_band_boundaries(self, strength, expected):
        assert strength_band(strength) == expected

    def test_none_strength_is_not_described(self):
        """没有 strength（如 DISC 那种只有名次的结果）时不能编一个程度。"""
        assert strength_band(None) == "无"


def _req(strengths: list[float | None]) -> EnhanceRequest:
    return EnhanceRequest.model_validate(
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
                    "key": f"K{i}",
                    "label": f"维度{i}",
                    "kind": "dimension",
                    "score": 10,
                    "threshold": 12,
                    "pole": "E",
                    **({"strength": s} if s is not None else {}),
                }
                for i, s in enumerate(strengths)
            ],
        }
    )


def test_measures_lines_carry_band_words():
    """prompt 里必须同时出现数值与程度词，模型才有东西可照抄。"""
    lines = _measures_lines(_req([0.3333, -0.6667]))

    assert "strength=+0.33" in lines
    assert "程度=轻度偏向" in lines
    assert "strength=-0.67" in lines
    assert "程度=明显偏向" in lines


def test_measures_lines_omit_band_without_strength():
    lines = _measures_lines(_req([None]))

    assert "strength=" not in lines
    assert "程度=" not in lines


def test_system_prompt_is_byte_stable():
    """system prompt 必须与请求无关，否则 DeepSeek 前缀缓存永远打不中。"""
    a, _ = build_prompt(_req([0.1]), [])
    b, _ = build_prompt(_req([-0.9, 0.5]), [])
    assert a == b


class TestParseLlmJson:
    def test_plain_json(self):
        assert _parse_llm_json('{"a": 1}') == {"a": 1}

    def test_code_fence(self):
        assert _parse_llm_json('```json\n{"a": 1}\n```') == {"a": 1}

    def test_prose_around_object(self):
        assert _parse_llm_json('好的，结果如下：{"a": 1} 希望有帮助') == {"a": 1}
