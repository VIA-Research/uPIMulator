import pytest

from abi.word.data_word import DataWord
from abi.word.representation import Representation
from initializer.int_initializer import IntInitializer
from iss.sram.wram import WRAM
from util.config_loader import ConfigLoader


@pytest.fixture
def wram() -> WRAM:
    return WRAM()


def test_wram(wram: WRAM):
    for _ in range(100):
        address = IntInitializer.value_by_range(
            ConfigLoader.wram_offset(), ConfigLoader.wram_offset() + ConfigLoader.wram_size() - DataWord().size(),
        )

        if address % DataWord().size() == 0:
            data_word = DataWord()
            data_word.set_value(IntInitializer.value_by_width(Representation.UNSIGNED, data_word.width()))

            wram.write(address, data_word)
            data_word.value(Representation.UNSIGNED) == wram.read(address).value(Representation.UNSIGNED)
