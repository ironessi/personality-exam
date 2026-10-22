"""`EmbeddingModel.warmup_lengths` 的行为测试。

这里**不加载真实模型**：真机验证（P95 从 837ms 降到 19ms）在
`scripts/eval_retrieval.py` 与启动日志里做，单测只保证「按长度逐个编码」
这件事本身没写错——尤其是别把 29 个长度写成只编码一次。
"""

from __future__ import annotations

from typing import Any

from app.rag.embedding import EmbeddingModel


class _FakeTensor:
    def __init__(self, length: int) -> None:
        self.shape = (1, length)


class _FakeTokenizer:
    """把「字数」当 token 数，够用来驱动 warmup_lengths 的加字循环。"""

    def __call__(self, texts: list[str], **_kwargs: Any) -> dict[str, _FakeTensor]:
        return {"input_ids": _FakeTensor(len(texts[0]))}


class _FakeSentenceTransformer:
    def __init__(self) -> None:
        self.tokenizer = _FakeTokenizer()
        self.encoded: list[str] = []

    def encode(self, texts: list[str], **_kwargs: Any) -> list[list[float]]:
        self.encoded.extend(texts)
        return [[0.0] for _ in texts]


def _model_with_fake_transformer() -> tuple[EmbeddingModel, _FakeSentenceTransformer]:
    model = EmbeddingModel(model_path="/nonexistent", expected_dim=512)
    fake = _FakeSentenceTransformer()
    model._model = fake  # noqa: SLF001 —— 单测替身，绕开真实权重加载
    model._dim = 512  # noqa: SLF001
    return model, fake


def test_warmup_lengths_encodes_once_per_length() -> None:
    model, fake = _model_with_fake_transformer()

    elapsed, covered = model.warmup_lengths(lo=3, hi=7)

    assert covered == 5
    assert elapsed >= 0.0
    assert len(fake.encoded) == 5, "长度区间内每个长度都要各编码一次"


def test_warmup_lengths_grows_text_until_token_count_reaches_target() -> None:
    model, fake = _model_with_fake_transformer()

    model.warmup_lengths(lo=4, hi=4)

    assert fake.encoded == ["字" * 4], "字级切分下，4 个 token 需要 4 个字"


def test_warmup_lengths_accepts_single_length_range() -> None:
    model, fake = _model_with_fake_transformer()

    _, covered = model.warmup_lengths(lo=1, hi=1)

    assert covered == 1
    assert fake.encoded == ["字"]


def test_warmup_lengths_clamps_zero_to_one() -> None:
    """token 长度为 0 没有意义；传 0 应从 1 开始，而不是空转一轮。"""
    model, fake = _model_with_fake_transformer()

    _, covered = model.warmup_lengths(lo=0, hi=2)

    assert covered == 2
    assert len(fake.encoded) == 2
