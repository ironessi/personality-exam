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

    def test_negative_factor_strength_is_normal_not_a_bias(self):
        """SCL-90 因子是单极的：负 strength = 低于参考线，不是「偏向低分」。

        实测依据：scl-12 十项全 1 分，生成结果写成「各因子均处于正常范围，
        程度为「明显偏向」」——|strength|=0.5 被套进双极的量表，一句话自相矛盾。
        """
        assert strength_band(-0.5, "factor") == "低于参考线（属正常范围，不要写成偏向）"
        assert strength_band(-1.0, "factor") == "低于参考线（属正常范围，不要写成偏向）"

    def test_positive_factor_strength_keeps_bands(self):
        assert strength_band(0.55, "factor") == "明显偏向"
        assert strength_band(0.1, "factor") == "接近中间（不要强调任何一端的偏向）"

    def test_negative_dimension_strength_still_uses_magnitude(self):
        """维度是双极的：负值只表示落在另一极，强度照样看绝对值。"""
        assert strength_band(-0.67, "dimension") == "明显偏向"


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


class TestFitQuote:
    """引文裁剪：超长引文曾让整次增强 500（suggestion 侧没有 try）。"""

    def test_short_quote_untouched(self):
        from app.rag.synthesizer import _fit_quote

        assert _fit_quote("短句。") == "短句。"

    def test_long_quote_cut_at_sentence_end(self):
        from app.rag.synthesizer import _fit_quote

        quote = "第一句讲清楚一件事。" + "第二句很长很长" * 12
        got = _fit_quote(quote)
        assert got == "第一句讲清楚一件事。"

    def test_long_quote_without_boundary_is_dropped(self):
        """连逗号顿号都没有时宁可丢掉，也不硬切出半句冒充原文。"""
        from app.rag.synthesizer import _fit_quote

        assert _fit_quote("没有任何标点" * 20) is None

    def test_long_quote_without_sentence_end_falls_back_to_comma(self):
        """整段没有句号的长复句在逗号处收尾。

        实测依据：20 条验收用例里被丢弃的引文全部是这种「A 会…，擅长…，因此…」
        的长复句，语料本身就爱用逗号连写。硬丢等于白丢一条 citation；在逗号处
        截断仍是**逐字原文**，Stage 1 与 Stage 1.5 照常能验。
        """
        from app.rag.synthesizer import _fit_quote

        quote = "F 倾向明显的人会把人、关系影响、价值取向与和谐程度放在前面考虑，" + "擅长体察语气变化" * 10
        got = _fit_quote(quote)
        assert got is not None
        assert len(got) <= 64
        assert got.endswith("，")
        assert quote.startswith(got)

    def test_soft_boundary_must_keep_meaningful_length(self):
        """次级边界切得太靠前就放弃：剩几个字的「引文」没有背书价值。"""
        from app.rag.synthesizer import _fit_quote

        assert _fit_quote("短，" + "没有句号的超长后半段" * 10) is None

    def test_empty_quote(self):
        from app.rag.synthesizer import _fit_quote

        assert _fit_quote(None) is None
        assert _fit_quote("   ") is None

    def test_overlong_suggestion_does_not_break_generation(self):
        """一条坏建议不该拖垮整次生成：它自己失去引文（随后被验证器丢弃），
        其余建议必须原样保留。"""
        from app.rag.synthesizer import _to_enhance_data

        payload = {
            "enhanced_text": "总结",
            "sections": [{"title": "节", "content": "内容"}],
            "suggestions": [
                {"text": "好建议", "basis": "依据一", "chunk_id": "c1", "evidence_quote": "原文一句。"},
                {"text": "坏建议", "basis": "依据二", "chunk_id": "c1",
                 "evidence_quote": "没有任何标点" * 20},
            ],
            "citations": [],
        }
        data = _to_enhance_data(_req([0.5]), payload, "m", 1, 1, 1)
        got = {s.text: s.evidence_quote for s in data.suggestions}
        assert got["好建议"] == "原文一句。"
        assert got["坏建议"] is None, "无句末标点的超长引文应被丢掉（随后由验证器剔除该建议）"

    def test_overlong_citation_is_dropped_not_fatal(self):
        from app.rag.synthesizer import _to_enhance_data

        payload = {
            "enhanced_text": "总结",
            "sections": [{"title": "节", "content": "内容"}],
            "suggestions": [],
            "citations": [
                {"claim": "结论", "source": "s.md", "chunk_id": "c1",
                 "evidence_quote": "没有任何标点" * 20, "section": 0},
            ],
        }
        data = _to_enhance_data(_req([0.5]), payload, "m", 1, 1, 1)
        assert data.citations == [] or data.citations[0].evidence_quote is None
