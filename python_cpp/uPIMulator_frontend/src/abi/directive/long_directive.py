from abi.word.immediate import Immediate
from abi.word.representation import Representation


class LongDirective:
    def __init__(self, value: int):
        self._value: Immediate = Immediate(Representation.UNSIGNED, 8 * LongDirective.size(), value)

    def value(self) -> int:
        return self._value.value()

    @staticmethod
    def size() -> int:
        return 4
