from abi.word.immediate import Immediate
from abi.word.representation import Representation


class ByteDirective:
    def __init__(self, value: int):
        self._value: Immediate = Immediate(Representation.UNSIGNED, 8 * ByteDirective.size(), value)

    def value(self) -> int:
        return self._value.value()

    @staticmethod
    def size() -> int:
        return 1
