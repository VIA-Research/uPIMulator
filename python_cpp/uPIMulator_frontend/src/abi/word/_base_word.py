import math
from typing import List

from abi.word.representation import Representation
from encoder.byte import Byte


class BaseWord:
    def __init__(self, width: int):
        assert width > 0

        self._bits: List[bool] = [False for _ in range(width)]

    def width(self) -> int:
        return len(self._bits)

    def size(self) -> int:
        return self.width() // 8

    def sign_bit(self) -> bool:
        return self._bits[-1]

    def bit(self, i: int) -> bool:
        return self._bits[i]

    def set_bit(self, i: int) -> None:
        self._bits[i] = True

    def clear_bit(self, i: int) -> None:
        self._bits[i] = False

    def bit_slice(self, representation: Representation, begin: int, end: int) -> int:
        assert 0 <= begin < end <= self.width()

        slice_width = end - begin

        value = 0
        for i in range(slice_width):
            if self.bit(i + begin):
                if representation == Representation.SIGNED and i == slice_width - 1:
                    value -= 2 ** i
                else:
                    value += 2 ** i
        return value

    def set_bit_slice(self, begin: int, end: int, value: int) -> None:
        assert 0 <= begin < end <= self.width()

        if value >= 0:
            self._set_positive_bit_slice(begin, end, value)
        else:
            self._set_negative_bit_slice(begin, end, value)

    def value(self, representation: Representation) -> int:
        return self.bit_slice(representation, 0, self.width())

    def set_value(self, value: int) -> None:
        self.set_bit_slice(0, self.width(), value)

    def to_bytes(self) -> List[Byte]:
        num_bytes = math.ceil(self.width() / 8.0)

        bytes_: List[Byte] = []
        for i in range(num_bytes):
            begin = 8 * i
            end = min(8 * (i + 1), self.width())

            bytes_.append(Byte(self.bit_slice(Representation.UNSIGNED, begin, end)))
        return bytes_

    def from_bytes(self, bytes_: List[Byte]) -> None:
        for i, byte in enumerate(bytes_):
            begin = 8 * i
            end = min(8 * (i + 1), self.width())

            self.set_bit_slice(begin, end, byte.value())

    def _set_positive_bit_slice(self, begin: int, end: int, value: int) -> None:
        assert value >= 0

        slice_width = end - begin
        for i in range(slice_width):
            if value % 2:
                self.set_bit(i + begin)
            else:
                self.clear_bit(i + begin)

            value //= 2

        assert value == 0

    def _set_negative_bit_slice(self, begin: int, end: int, value: int) -> None:
        assert value < 0

        slice_width = end - begin

        self.set_bit(end - 1)
        value += 2 ** (slice_width - 1)

        if begin + 1 < end:
            self._set_positive_bit_slice(begin, end - 1, value)
