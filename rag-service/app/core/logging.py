"""统一日志。

两个目标：
1. 单行可读，带毫秒时间戳与 logger 名；
2. 生产可用：``LOG_JSON=true`` 时输出 JSON Lines，便于将来采集。

第三方的 uvicorn / chromadb / sentence_transformers 日志在 INFO 下噪音很大，
统一压到 WARNING，只留自己的 INFO。
"""

from __future__ import annotations

import json
import logging
import sys
from datetime import datetime, timezone

_CONFIGURED = False

# 这些 logger 在 INFO 级别会刷屏，压掉
_NOISY_LOGGERS = (
    "httpx",
    "httpcore",
    "urllib3",
    "chromadb",
    "sentence_transformers",
    "transformers",
    "torch",
    "jieba",
    "watchfiles",
)


class JsonFormatter(logging.Formatter):
    def format(self, record: logging.LogRecord) -> str:
        payload = {
            "ts": datetime.fromtimestamp(record.created, tz=timezone.utc).isoformat(
                timespec="milliseconds"
            ),
            "level": record.levelname,
            "logger": record.name,
            "msg": record.getMessage(),
        }
        if record.exc_info:
            payload["exc"] = self.formatException(record.exc_info)
        return json.dumps(payload, ensure_ascii=False)


class TextFormatter(logging.Formatter):
    def __init__(self) -> None:
        super().__init__(
            fmt="%(asctime)s.%(msecs)03d %(levelname)-7s %(name)-28s %(message)s",
            datefmt="%H:%M:%S",
        )


def configure_logging(level: str = "INFO", *, json_output: bool = False) -> None:
    """幂等：重复调用不会叠加 handler（否则每次 reload 日志翻倍）。"""
    global _CONFIGURED

    handler = logging.StreamHandler(sys.stdout)
    handler.setFormatter(JsonFormatter() if json_output else TextFormatter())

    root = logging.getLogger()
    if not _CONFIGURED:
        root.handlers.clear()
        root.addHandler(handler)
    root.setLevel(getattr(logging, level.upper(), logging.INFO))

    for name in _NOISY_LOGGERS:
        logging.getLogger(name).setLevel(logging.WARNING)

    _CONFIGURED = True


def get_logger(name: str) -> logging.Logger:
    return logging.getLogger(name)
