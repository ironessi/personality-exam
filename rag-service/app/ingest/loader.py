"""原始文献加载：txt / md（§4 :195）。

命名约定：``data/raw/{namespace}/`` 下的文件，其 ``kb_namespace`` 取**父目录名**。
这是 ``kb_namespace`` 三处共用（Chroma where / BM25 索引选择 / 原始目录）的由来，
所以文件不能放在 ``data/raw/`` 根下——那会得到一个空的 namespace。

``SOURCES.md`` 是来源与许可登记表，不是语料；``_`` / ``.`` 开头的文件同理，
两者都会被 ``iter_raw_docs`` 跳过。

pdf 需要额外依赖（``pypdf``），P0-2 的语料全是自撰 markdown，暂不声明；
``SUPPORTED_SUFFIXES`` 之外的格式会被明确跳过并告警，而不是静默吞掉。
"""

from __future__ import annotations

from dataclasses import dataclass
from pathlib import Path

from app.core.logging import get_logger

logger = get_logger(__name__)


@dataclass
class RawDoc:
    path: Path
    text: str
    kb_namespace: str
    source: str  # 展示用文件名，如 16型人格类型详解.md


SUPPORTED_SUFFIXES: tuple[str, ...] = (".txt", ".md", ".markdown")

# 不是语料的文件名（小写比较）：来源登记表、说明文件
SKIP_NAMES: frozenset[str] = frozenset({"sources.md", "readme.md", "changelog.md"})
# 下划线/点开头的文件按约定是草稿或隐藏文件
SKIP_PREFIXES: tuple[str, ...] = ("_", ".")


def _read_text(path: Path) -> str:
    """按 UTF-8 读；失败再试 GBK（中文语料最常见的两种编码）。

    不用 ``errors="ignore"`` 兜底：那会**静默吃掉**字符，切块结果看不出问题，
    却会让引用回验（Stage 1 字面子串校验）永远匹配不上原文。
    """
    raw = path.read_bytes()
    for enc in ("utf-8", "utf-8-sig", "gbk"):
        try:
            return raw.decode(enc)
        except UnicodeDecodeError:
            continue
    logger.warning("无法用 utf-8/gbk 解码，改用 utf-8 errors=replace：%s", path)
    return raw.decode("utf-8", errors="replace")


def should_skip(path: Path) -> bool:
    """``SOURCES.md``、``README.md``、``_``/``.`` 开头的文件都不是语料。"""
    name = path.name.lower()
    return name in SKIP_NAMES or name.startswith(SKIP_PREFIXES)


def load_file(path: Path, kb_namespace: str) -> RawDoc:
    """加载单个文件。

    Raises:
        ValueError: 后缀不在 ``SUPPORTED_SUFFIXES`` 里。
    """
    if path.suffix.lower() not in SUPPORTED_SUFFIXES:
        raise ValueError(
            f"不支持的语料格式 {path.suffix!r}：{path}。"
            f"支持 {SUPPORTED_SUFFIXES}；pdf 请先转成 markdown。"
        )
    return RawDoc(
        path=path,
        text=_read_text(path),
        kb_namespace=kb_namespace,
        source=path.name,
    )


def iter_raw_docs(raw_root: Path) -> list[RawDoc]:
    """遍历 ``data/raw/{ns}/**``，按目录名定 namespace。

    返回顺序**确定性**：``sorted(rglob)`` 即按「namespace → 相对路径」排序。
    切块的 ``chunk_id`` 序号建立在这个顺序上，顺序一抖动，同一份语料每次
    重建就会得到不同的 id，历史引用随之悬空。

    Returns:
        加载成功的 ``RawDoc`` 列表；``data/raw`` 不存在时返回空列表（不是异常）。
    """
    raw_root = Path(raw_root)
    if not raw_root.is_dir():
        logger.warning("raw 目录不存在：%s", raw_root)
        return []

    docs: list[RawDoc] = []
    for path in sorted(raw_root.rglob("*")):
        if not path.is_file():
            continue
        if should_skip(path):
            continue

        rel = path.relative_to(raw_root)
        # namespace = 第一层目录名。根下的文件没有 namespace，
        # 跳过而不是猜一个——猜错会污染过滤检索。
        if len(rel.parts) < 2:
            logger.warning("跳过 data/raw 根下的文件（无 namespace）：%s", rel)
            continue
        if path.suffix.lower() not in SUPPORTED_SUFFIXES:
            logger.warning("跳过不支持的格式：%s", rel)
            continue

        docs.append(load_file(path, rel.parts[0]))

    logger.info("加载 %d 份原始文档 ← %s", len(docs), raw_root)
    return docs
