from abi.word.instruction_address_word import InstructionAddressWord
from abi.word.representation import Representation


class PCRegister:
    def __init__(self):
        self._word: InstructionAddressWord = InstructionAddressWord()

    def read(self) -> int:
        return self._word.value(Representation.UNSIGNED)

    def write(self, value: int) -> None:
        self._word.set_value(value)

    def increment(self) -> None:
        self.write(self.read() + (InstructionAddressWord().size()))
