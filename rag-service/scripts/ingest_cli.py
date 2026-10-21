#!/usr/bin/env python
"""命令行入库（§4 :207）。

推荐直接用 ``python -m app.ingest.pipeline split|ingest``（§7.4）；
本脚本是等价的一步式封装：切块 + 向量化 + 写 Chroma。

用法::

    /Users/iron/rag/bin/python scripts/ingest_cli.py                 # 全量（切块 + 入库）
    /Users/iron/rag/bin/python scripts/ingest_cli.py --ns typology   # 只刷一个构念域
    /Users/iron/rag/bin/python scripts/ingest_cli.py --dry-run       # 只切块，不写库

``--ns`` 的取值是**构念域**（ADR-0008）：``symptom`` / ``trait`` / ``typology`` /
``social`` / ``wellbeing`` / ``career`` / ``ability`` / ``common``，**不是测评名**。
"""

from __future__ import annotations

import argparse
import sys
from pathlib import Path

SERVICE_ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(SERVICE_ROOT))


def main() -> int:
    parser = argparse.ArgumentParser(description="知识库入库（P0-2）")
    parser.add_argument("--ns", action="append", help="只处理指定 namespace，可重复")
    parser.add_argument("--dry-run", action="store_true", help="只切块，不写库")
    args = parser.parse_args()

    from app.core.logging import configure_logging
    from app.core.config import get_settings
    from app.ingest.pipeline import ingest_all

    s = get_settings()
    configure_logging(s.log_level)

    namespaces = args.ns or s.kb_namespaces
    print(f"namespace: {namespaces}")
    print(f"raw:       {s.data_raw_dir}")
    print(f"chroma:    {s.chroma_path}")
    if args.dry_run:
        print("模式:      dry-run（只切块，不写 Chroma）")

    stats = ingest_all(namespaces=namespaces, rebuild=not args.dry_run)
    total = sum(stats.values())
    for ns, n in sorted(stats.items()):
        print(f"  {ns:10} {n:5d} 块")
    print(f"  {'合计':10} {total:5d} 块")
    return 0 if total else 1


if __name__ == "__main__":
    raise SystemExit(main())
