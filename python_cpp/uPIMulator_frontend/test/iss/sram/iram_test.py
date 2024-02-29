import pytest

from abi.word.instruction_word import InstructionWord
from abi.word.representation import Representation
from initializer.int_initializer import IntInitializer
from iss.sram.iram import IRAM
from util.config_loader import ConfigLoader


@pytest.fixture
def iram() -> IRAM:
    return IRAM()


def test_iram(iram: IRAM):
    for _ in range(100):
        address = IntInitializer.value_by_range(
            ConfigLoader.iram_offset(),
            ConfigLoader.iram_offset() + ConfigLoader.iram_size() - InstructionWord().size(),
        )

        if address % InstructionWord().size() == 0:
            instruction_word = InstructionWord()
            instruction_word.set_value(IntInitializer.value_by_width(Representation.UNSIGNED, instruction_word.width()))

            iram.write(address, instruction_word)
            instruction_word.value(Representation.UNSIGNED) == iram.read(address).value(Representation.UNSIGNED)
