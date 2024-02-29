import math
from typing import List

from abi.word.data_word import BaseWord
from abi.word.representation import Representation
from encoder.byte import Byte


class Immediate:
    def __init__(self, representation: Representation, width: int, value: int):
        assert width > 0

        self._representation: Representation = representation
        self._width: int = width

        self._word: BaseWord = BaseWord(width)
        self._word.set_bit_slice(0, width, value)

    def representation(self) -> Representation:
        return self._representation

    def width(self) -> int:
        return self._width

    def bit(self, i: int) -> bool:
        return self._word.bit(i)

    def bit_slice(self, representation: Representation, begin: int, end: int) -> int:
        return self._word.bit_slice(representation, begin, end)

    def value(self) -> int:
        return self._word.bit_slice(self._representation, 0, self._width)

    def to_bytes(self) -> List[Byte]:
        num_bytes = math.ceil(self.width() / 8.0)
        return self._word.to_bytes()[:num_bytes]
