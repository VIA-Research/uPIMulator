from typing import List

from abi.word.instruction_address_word import InstructionAddressWord
from abi.word.instruction_word import InstructionWord
from abi.word.representation import Representation
from util.config_loader import ConfigLoader


class IRAM:
    def __init__(self):
        self._address: InstructionAddressWord = InstructionAddressWord()
        self._address.set_value(ConfigLoader.iram_offset())

        self._instruction_words: List[InstructionWord] = [
            InstructionWord() for _ in range(ConfigLoader.iram_size() // InstructionWord().size())
        ]

    def address(self) -> int:
        return self._address.value(Representation.UNSIGNED)

    def read(self, address: int) -> InstructionWord:
        index = self._index(address)
        instruction_word = InstructionWord()
        instruction_word.set_value(self._instruction_words[index].value(Representation.UNSIGNED))
        return instruction_word

    def write(self, address: int, instruction_word: InstructionWord) -> None:
        index = self._index(address)
        self._instruction_words[index].set_value(instruction_word.value(Representation.UNSIGNED))

    def _index(self, address: int) -> int:
        assert address % InstructionWord().size() == 0
        index = (address - self.address()) // InstructionWord().size()
        assert 0 <= index < len(self._instruction_words)

        return index
