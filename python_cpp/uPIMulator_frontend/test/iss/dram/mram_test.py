import pytest

from abi.word.data_word import DataWord
from abi.word.representation import Representation
from initializer.int_initializer import IntInitializer
from iss.dram.mram import MRAM
from iss.dram.mram_command import MRAMCommand
from util.config_loader import ConfigLoader


@pytest.fixture
def mram() -> MRAM:
    return MRAM()


def test_mram(mram: MRAM):
    for _ in range(100):
        address = IntInitializer.value_by_range(
            ConfigLoader.mram_offset(),
            ConfigLoader.mram_offset() + ConfigLoader.mram_size() - ConfigLoader.min_access_granularity(),
        )

        if address % ConfigLoader.min_access_granularity() == 0:
            write_mram_command = MRAMCommand(
                MRAMCommand.Operation.WRITE, address, ConfigLoader.min_access_granularity(),
            )
            read_mram_command = MRAMCommand(MRAMCommand.Operation.READ, address, ConfigLoader.min_access_granularity(),)

            data_words = [DataWord() for _ in range(ConfigLoader.min_access_granularity() // DataWord().size())]
            for data_word in data_words:
                data_word.set_value(IntInitializer.value_by_width(Representation.UNSIGNED, DataWord().width()))
            write_mram_command.set_data_words(data_words)

            assert not mram.can_pop()

            assert mram.can_push()
            mram.push(write_mram_command)

            assert mram.can_push()
            mram.push(read_mram_command)

            assert mram.can_pop()
            assert write_mram_command == mram.pop()

            assert mram.can_pop()
            assert read_mram_command == mram.pop()

            for data_word, read_mram_command_data_word in zip(data_words, read_mram_command.data_words()):
                assert data_word.value(Representation.UNSIGNED) == read_mram_command_data_word.value(
                    Representation.UNSIGNED
                )

            assert not mram.can_pop()
