from abi.word.immediate import Immediate
from abi.word.representation import Representation


class ZeroDirective:
    def __init__(self, size: int, value: int):
        self._size: int = size
        self._value: Immediate = Immediate(Representation.UNSIGNED, 8, value)

    def size(self) -> int:
        return self._size

    def value(self) -> int:
        return self._value.value()
