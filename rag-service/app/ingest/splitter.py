"""语义切块 + 元数据（§7.3）。

每个 chunk 必须携带四个元数据字段（§7.3 :599）：

| 字段 | 说明 |
|---|---|
| ``source`` | 展示用文件名 |
| ``kb_namespace`` | **构念域**（ADR-0008），过滤检索与索引选择的唯一权威值 |
| ``section_title`` | 章节定位，供引用回显 |
| ``chunk_id`` | 形如 ``typology_00037``，Stage 1 字面校验的锚点 |

⚠️ ``kb_namespace`` 的取值是**构念域**（``symptom`` / ``trait`` / ``typology`` /
``social`` / ``wellbeing`` / ``career`` / ``ability`` / ``common``），**不是测评名**。
把 ``mbti`` 当 namespace 会让同一份知识在多个测评目录下各存一份，而 ``chunk_id``
前缀就是 namespace——两套命名混用时 ``upsert`` 会相互覆盖且**不报错**（ADR-0008）。

切块参数取 §7.3 的中值区间：块 300~400 字符、重叠 50 字符。
**不要把块调大**：块越大越容易在一个块里混进两个主题，检索命中的块会
把无关内容一起带进 prompt，稀释引用并拉高 token 成本。

``chunk_id``：``f"{kb_namespace}_{seq:05d}"``，``seq`` 由调用方按
「namespace → 文档路径 → 文档内块序」的顺序连续分配（见 ``iter_raw_docs`` 与
``pipeline.split_corpus``），因此对**同一份语料**是确定性可重现的——
重复入库得到同一批 id，Chroma ``upsert`` 自然幂等，不产生重复块。

已知代价：语料**增删**会让排序靠后的文档整体重编号，历史记录里已写死的
``chunk_id`` 可能指向新块。P0 语料固定，可接受；语料开始频繁变动时应改成
把 id 绑定到内容哈希（届时 §5.2 契约里的 ``chunk_id`` 语义要同步说明）。
"""

from __future__ import annotations

import re
from dataclasses import dataclass, field
from pathlib import Path

# §7.3 切块参数
DEFAULT_CHUNK_SIZE = 400
DEFAULT_MIN_CHUNK = 300
DEFAULT_OVERLAP = 50

_HEADING_RE = re.compile(r"^(#{1,6})\s+(.+?)\s*$")
# 同一套规则的多行版本：一次扫出全文所有标题，用来判定分节级别
_HEADING_ML_RE = re.compile(r"^(#{1,6})\s+(.+?)\s*$", re.MULTILINE)
_DOC_TITLE_RE = re.compile(r"^#\s+(.+?)\s*$", re.MULTILINE)
_LIST_ITEM_RE = re.compile(r"^\s*(?:[-*+]|\d+[.)])\s+")
_SENTENCE_END = "。！？；!?;"
# 中文之间的软换行要直接接起来，不留空格
_CJK = r"\u4e00-\u9fff\u3400-\u4dbf"
_CJK_SPACE_RE = re.compile(rf"(?<=[{_CJK}]) (?=[{_CJK}])")


@dataclass
class Chunk:
    chunk_id: str
    text: str
    source: str
    kb_namespace: str
    section_title: str = ""
    metadata: dict = field(default_factory=dict)

    def to_metadata(self) -> dict:
        """Chroma / BM25 共用的元数据：§7.3 要求的四个字段，一个都不少。"""
        return {
            "source": self.source,
            "kb_namespace": self.kb_namespace,
            "section_title": self.section_title,
            "chunk_id": self.chunk_id,
        }

    def to_dict(self) -> dict:
        """``data/processed/*.jsonl`` 的一行（可审计）。"""
        return {"chunk_id": self.chunk_id, "text": self.text, **self.to_metadata()}


def _doc_title(doc) -> str:
    """文档标题：取第一个 ``#`` 标题；纯文本文件退回文件名。

    它只在**小节标题为空**时兜底，好过让引用回显出空字符串。
    """
    m = _DOC_TITLE_RE.search(doc.text)
    if m:
        return m.group(1).strip()
    return Path(doc.source).stem


def _heading_levels(text: str) -> tuple[int | None, int | None]:
    """返回 ``(文档标题级别, 分节级别)``，两者都可能为 ``None``。

    规则：首个标题若出现在最浅级别上、且该级别只有它一个，它就只是**文档标题**
    （不进正文，也不参与分节）；**分节级别**取其余标题里最浅的一级。没有任何
    标题的纯文本文档两级都是 ``None``，整篇作为一节。

    为什么必须按文档自己的顶层结构分节：只要把任意级别的标题都当分节线，一份用
    ``###`` 组织小节的语料就会被切成一堆碎块——实测一份 25,213 字符的文档被切成
    **135 块，平均 175 字符、最短 83 字符**。碎块装不下一段完整论证，检索命中后
    进 prompt 的上下文太薄，引用回验也更容易失败。
    """
    levels = [len(m.group(1)) for m in _HEADING_ML_RE.finditer(text)]
    if not levels:
        return None, None
    shallowest = min(levels)
    if levels[0] == shallowest and levels.count(shallowest) == 1:
        return shallowest, (min(levels[1:]) if len(levels) > 1 else None)
    return None, shallowest


def _parse_sections(text: str) -> list[tuple[str, str]]:
    """按 markdown 标题把正文切成 ``(section_title, body)``。

    只按**分节级别**（见 ``_heading_levels``）断开：更深层的小节标题留在正文里，
    当作段落的引导行。标题行**不进** body——它已经作为 ``section_title`` 单独
    携带，再留在正文只会让每个块开头都黏着一行标题。
    """
    title_level, level = _heading_levels(text)
    sections: list[tuple[str, str]] = []
    title = ""
    buf: list[str] = []

    def flush() -> None:
        body = "\n".join(buf).strip()
        buf.clear()
        if body:
            sections.append((title, body))

    for line in text.splitlines():
        m = _HEADING_RE.match(line)
        if m:
            lv = len(m.group(1))
            if level is not None and lv == level:
                flush()
                title = m.group(2).strip()
                continue
            if title_level is not None and lv <= title_level:
                continue  # 文档标题行本身不进正文
        buf.append(line)
    flush()
    return sections


def _join_lines(lines: list[str]) -> str:
    return _CJK_SPACE_RE.sub("", " ".join(lines))


def _blocks(body: str) -> list[str]:
    """按空行切块；列表项与深层小节标题各自成块。

    分隔符不在这里补，统一由 ``_units`` 负责——它是唯一知道「这一段有没有以
    句末标点收尾」的地方。
    """
    out: list[str] = []
    for raw in re.split(r"\n\s*\n", body):
        lines = [ln.strip() for ln in raw.splitlines() if ln.strip()]
        if not lines:
            continue
        if any(_LIST_ITEM_RE.match(ln) for ln in lines):
            out.extend(lines)
            continue
        # 深层小节标题自成一块：语料里的 `### X` 后面**直接**跟正文（不留空行），
        # 交给 _join_lines 会被当成软换行续行，粘出「### 核心特征INTJ 型的人…」
        run: list[str] = []
        for ln in lines:
            if _HEADING_RE.match(ln):
                if run:
                    out.append(_join_lines(run))
                    run = []
                out.append(ln)
            else:
                run.append(ln)
        if run:
            out.append(_join_lines(run))
    return out


def split_sentences(text: str) -> list[str]:
    """按中英文句末标点切句，标点归前一句。"""
    out: list[str] = []
    buf: list[str] = []
    for ch in text:
        buf.append(ch)
        if ch in _SENTENCE_END:
            s = "".join(buf).strip()
            if s:
                out.append(s)
            buf = []
    tail = "".join(buf).strip()
    if tail:
        out.append(tail)
    return out


def _units(body: str) -> list[str]:
    """正文 → 句子序列（列表项也各自成句）。

    ``_pack`` 是**无分隔**拼接（``cur += unit``），所以对没有以句末标点收尾的块
    （列表项、裸标题、表格行）要在其末句补一个换行；不补就会拼出
    「…获取能量- 第二条：…」这种读不通、还会污染字面回验的粘句。
    """
    out: list[str] = []
    for block in _blocks(body):
        sentences = split_sentences(block)
        if sentences and block[-1] not in _SENTENCE_END:
            sentences[-1] += "\n"
        out.extend(sentences)
    return out


def _hard_split(sentence: str, chunk_size: int, overlap: int) -> list[str]:
    """单句长过一整块时的兜底硬切（步长 = chunk_size - overlap）。

    正常语料不该走到这里；走到了说明有超长表格或没加标点的段落，
    这里保证不丢字、不死循环。
    """
    step = max(1, chunk_size - overlap)
    return [sentence[i : i + chunk_size] for i in range(0, len(sentence), step)]


def _overlap_seed(text: str, overlap: int) -> str:
    """从块尾取**整句**作为下一块的开头（重叠 50 字符，§7.3）。

    优先取整句：字符硬切会把下一块的开头截成半句话，那半句既检索不到，
    又会污染引用回验的字面匹配。但**末句本身长过上界**时（中文里 200 字的
    长句并不罕见），只能硬切它的尾部——否则重叠会整个落空，而落空是
    **静默**的：块照样切出来了，只是跨块信息断了，没人会收到报错。
    """
    if overlap <= 0:
        return ""
    cap = overlap * 3  # 种子最多占 150 字符，别让它吃掉小半块
    seed = ""
    for s in reversed(split_sentences(text)):
        if not seed and len(s) > cap:
            seed = s[-overlap:]
            break
        if seed and len(seed) + len(s) > cap:
            break
        seed = s + seed
        if len(seed) >= overlap:
            break

    # 整块本身就短于 overlap 时，别把整块复制成下一块的开头（那会产出近似重复的块）
    return "" if seed.strip() == text.strip() else seed


def _pack(
    units: list[str],
    *,
    chunk_size: int,
    min_chunk: int,
    overlap: int,
) -> list[str]:
    """贪心打包：块内文本**绝不超过** ``chunk_size``，填到 ``min_chunk`` 以上才收口。

    重叠必须在这里一起算，不能事后再拼：种子占掉的长度要从本块的预算里扣，
    否则「块 ≤ 400」与「重叠 50」互相矛盾，必有一方静默失效。
    """
    chunks: list[str] = []
    seed = ""
    idx = 0
    n = len(units)

    while idx < n:
        # 种子放得下才带上；放不下就丢种子，绝不产出「只有重叠」的块
        cur = seed if len(seed) + len(units[idx]) <= chunk_size else ""
        added = 0

        while idx < n and len(cur) + len(units[idx]) <= chunk_size:
            cur += units[idx]
            idx += 1
            added += 1
            if len(cur) >= min_chunk:
                nxt = units[idx] if idx < n else None
                if nxt is None or len(cur) + len(nxt) > chunk_size:
                    break

        if added == 0:
            # 防御：单句就超长时 units 已在 split_document 里硬切过，正常到不了这里
            break

        chunks.append(cur)
        seed = _overlap_seed(cur, overlap)

    return chunks


def split_document(
    doc,
    *,
    chunk_size: int = DEFAULT_CHUNK_SIZE,
    min_chunk: int = DEFAULT_MIN_CHUNK,
    overlap: int = DEFAULT_OVERLAP,
    start_seq: int = 0,
) -> list[Chunk]:
    """把一个 ``RawDoc`` 切成带元数据的 ``Chunk`` 列表。

    Args:
        start_seq: 本 namespace 内已分配的块数，由调用方连续维护，
            保证 ``chunk_id`` 在 namespace 内唯一且可重现。

    Raises:
        ValueError: 参数越界（宁可报错，也不要悄悄产出畸形块）。
    """
    if chunk_size <= 0:
        raise ValueError(f"chunk_size 必须为正：{chunk_size}")
    if not 0 <= overlap < chunk_size:
        raise ValueError(f"overlap 必须在 [0, chunk_size) 内：{overlap} vs {chunk_size}")
    if not 0 < min_chunk <= chunk_size:
        raise ValueError(f"min_chunk 必须在 (0, chunk_size] 内：{min_chunk} vs {chunk_size}")

    doc_title = _doc_title(doc)
    chunks: list[Chunk] = []
    seq = start_seq

    for title, body in _parse_sections(doc.text):
        section_title = title or doc_title
        units: list[str] = []
        for unit in _units(body):
            if len(unit) <= chunk_size:
                units.append(unit)
            else:
                units.extend(_hard_split(unit, chunk_size, overlap))

        for text in _pack(
            units, chunk_size=chunk_size, min_chunk=min_chunk, overlap=overlap
        ):
            text = text.strip()
            if not text:
                continue
            chunks.append(
                Chunk(
                    chunk_id=f"{doc.kb_namespace}_{seq:05d}",
                    text=text,
                    source=doc.source,
                    kb_namespace=doc.kb_namespace,
                    section_title=section_title,
                )
            )
            seq += 1

    return chunks
