"""切块测试（P0-2）。

验收（§9 :774-777）：

- 每个 chunk 均含 ``source`` / ``kb_namespace`` / ``section_title`` / ``chunk_id``
  四个元数据字段；
- ``chunk_id`` 幂等——同一语料重复切块得到同一批 id，重复入库不产生重复块；
- 无整段版权原文（这一条靠人抽查，代码层能守的是「不得把标题行混进正文」
  与「不得丢正文」）。
"""

from __future__ import annotations

from pathlib import Path

import pytest

from app.ingest.loader import RawDoc
from app.ingest.splitter import (
    DEFAULT_CHUNK_SIZE,
    DEFAULT_MIN_CHUNK,
    DEFAULT_OVERLAP,
    split_document,
    split_sentences,
)

FOUR_FIELDS = {"source", "kb_namespace", "section_title", "chunk_id"}


def sent(i: int) -> str:
    """第 i 句：编号唯一，便于逐句核对有没有丢正文。"""
    return f"第{i:03d}句，这是一句长度适中的中文，用来把测试语料填到目标块大小。"


def body(n: int, per_para: int = 4) -> str:
    """把 n 句拼成若干自然段（空行分隔）。"""
    paras = []
    for start in range(0, n, per_para):
        paras.append("".join(sent(i) for i in range(start, min(start + per_para, n))))
    return "\n\n".join(paras)


def one_section(n: int = 40) -> str:
    return "## 第一节\n\n" + body(n)


def make_doc(text: str, ns: str = "typology", source: str = "测试语料.md") -> RawDoc:
    return RawDoc(path=Path(source), text=text, kb_namespace=ns, source=source)


# ── 元数据契约（§9 :775）─────────────────────────────────────────


def test_every_chunk_carries_the_four_fields():
    chunks = split_document(make_doc(one_section(40)))
    assert chunks, "至少应切出一块"
    for c in chunks:
        md = c.to_metadata()
        assert set(md) == FOUR_FIELDS
        assert all(str(v).strip() for v in md.values()), f"{c.chunk_id} 有字段为空"


def test_chunk_id_format_and_namespace_prefix():
    chunks = split_document(make_doc(one_section(40), ns="symptom"))
    assert [c.chunk_id for c in chunks] == [f"symptom_{i:05d}" for i in range(len(chunks))]


def test_start_seq_offsets_ids_without_collision():
    first = split_document(make_doc(one_section(20)))
    second = split_document(make_doc(one_section(20)), start_seq=len(first))
    assert second[0].chunk_id == f"typology_{len(first):05d}"
    assert len({c.chunk_id for c in first + second}) == len(first) + len(second)


def test_section_title_comes_from_heading():
    text = f"## 甲节\n\n{body(3)}\n\n## 乙节\n\n{body(3)}"
    assert {c.section_title for c in split_document(make_doc(text))} == {"甲节", "乙节"}


def test_heading_line_does_not_leak_into_chunk_text():
    text = f"## 甲节\n\n{body(10)}"
    assert all("## 甲节" not in c.text for c in split_document(make_doc(text)))


def test_h1_title_becomes_fallback_section_title():
    text = f"# 我的文档\n\n{body(10)}"
    assert {c.section_title for c in split_document(make_doc(text))} == {"我的文档"}


def test_plain_text_falls_back_to_filename():
    chunks = split_document(make_doc(body(10), source="无标题语料.md"))
    assert {c.section_title for c in chunks} == {"无标题语料"}


# ── chunk_id 幂等（§9 :777）──────────────────────────────────────


def test_deterministic_across_repeated_calls():
    doc = make_doc(one_section(60))
    a = [(c.chunk_id, c.text, c.section_title) for c in split_document(doc)]
    b = [(c.chunk_id, c.text, c.section_title) for c in split_document(doc)]
    assert a == b


# ── 块尺寸与重叠（§7.3）───────────────────────────────────────────


def test_no_chunk_exceeds_chunk_size():
    for c in split_document(make_doc(one_section(80))):
        assert len(c.text) <= DEFAULT_CHUNK_SIZE, f"{c.chunk_id} 有 {len(c.text)} 字符"


def test_only_the_last_chunk_of_a_section_may_be_short():
    text = f"## 甲节\n\n{body(40)}\n\n## 乙节\n\n{body(40)}"
    by_section: dict[str, list] = {}
    for c in split_document(make_doc(text)):
        by_section.setdefault(c.section_title, []).append(c)
    for title, group in by_section.items():
        for c in group[:-1]:
            assert len(c.text) >= DEFAULT_MIN_CHUNK, f"{title}/{c.chunk_id} 只有 {len(c.text)} 字符"


def test_consecutive_chunks_in_a_section_overlap():
    chunks = split_document(make_doc(one_section(60)))
    assert len(chunks) >= 3
    overlapped = sum(
        1
        for a, b in zip(chunks, chunks[1:])
        if a.section_title == b.section_title and b.text[:DEFAULT_OVERLAP] in a.text
    )
    assert overlapped >= 1, "同节相邻块之间应有重叠，否则跨块信息会断"


def test_overlap_is_not_the_whole_previous_chunk():
    """种子不得把上一块整段复制过来（那会产出近似重复的块）。"""
    for c in split_document(make_doc(one_section(60))):
        assert len(c.text) > DEFAULT_OVERLAP


# ── 内容完整性 ──────────────────────────────────────────────────


def test_no_body_sentence_is_lost():
    text = one_section(60)
    joined = "".join(c.text for c in split_document(make_doc(text)))
    for i in range(60):
        assert sent(i) in joined, f"第 {i} 句丢了"


def test_list_items_are_not_glued_into_one_sentence():
    text = "## 列表\n\n- 第一条：外向者从人际互动中获取能量\n- 第二条：内向者从独处中恢复精力\n"
    chunks = split_document(make_doc(text))
    joined = "".join(c.text for c in chunks)
    assert "第一条" in joined and "第二条" in joined
    # 两个条目可以同块（400 字符预算装得下），但不得被拼成一句：
    # 列表项通常不以句号收尾，不留分隔就会得到「…能量第二条：…」
    assert "能量第二条" not in joined, "列表项被粘成了一句"


def test_deeper_headings_do_not_fragment_sections():
    """### 小节不该成为分节线。

    实测教训：一份 25,213 字符的语料用 ### 组织 128 个小节，被判成 146 个章节后
    切成 135 块、平均 175 字符、最短 83 字符——碎块装不下一段完整论证。
    """
    text = (
        f"# 文档标题\n\n## 甲型\n\n{body(20)}\n\n### 优势\n{body(20)}"
        f"\n\n### 盲点\n{body(20)}\n\n## 乙型\n\n{body(20)}"
    )
    chunks = split_document(make_doc(text))
    assert {c.section_title for c in chunks} == {"甲型", "乙型"}
    joined = "".join(c.text for c in chunks)
    assert "### 优势" in joined and "### 盲点" in joined, "深层标题应留在正文里"


def test_deeper_heading_is_not_glued_to_the_following_sentence():
    text = "## 甲节\n\n### 优势\n优势通常表现为清晰的判断。\n"
    joined = "".join(c.text for c in split_document(make_doc(text)))
    assert "### 优势优势" not in joined, "深层标题与正文粘成了一句"


def test_doc_title_line_is_not_leaked_into_any_chunk():
    text = f"# 文档标题\n\n## 甲节\n\n{body(20)}"
    chunks = split_document(make_doc(text))
    assert all("文档标题" not in c.text for c in chunks)


def test_oversized_sentence_without_punctuation_is_hard_split():
    text = "## 长句\n\n" + ("啊" * 1000)
    chunks = split_document(make_doc(text))
    assert len(chunks) >= 3
    assert sum(len(c.text) for c in chunks) >= 1000
    assert all(len(c.text) <= DEFAULT_CHUNK_SIZE for c in chunks)


# ── 参数校验：宁可报错，也不要悄悄产出畸形块 ──────────────────────


@pytest.mark.parametrize(
    "kwargs",
    [
        {"chunk_size": 0},
        {"chunk_size": -1},
        {"overlap": -1},
        {"overlap": DEFAULT_CHUNK_SIZE},
        {"min_chunk": 0},
        {"min_chunk": DEFAULT_CHUNK_SIZE + 1},
    ],
)
def test_invalid_parameters_raise(kwargs):
    with pytest.raises(ValueError):
        split_document(make_doc(one_section(10)), **kwargs)


def test_split_sentences_keeps_terminator_with_previous_sentence():
    assert split_sentences("甲。乙！丙？") == ["甲。", "乙！", "丙？"]
