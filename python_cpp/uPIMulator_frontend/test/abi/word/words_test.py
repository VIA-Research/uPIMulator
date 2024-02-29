from typing import Set, Union

import pytest

from abi.word.data_address_word import DataAddressWord
from abi.word.data_word import DataWord
from abi.word.double_data_word import DoubleDataWord
from abi.word.instruction_address_word import InstructionAddressWord
from abi.word.instruction_word import InstructionWord
from abi.word.representation import Representation
from initializer.int_initializer import IntInitializer


@pytest.fixture
def words() -> Set[Union[DataAddressWord, DataWord, InstructionAddressWord, InstructionWord, DoubleDataWord,]]:
    return {
        DataAddressWord(),
        DataWord(),
        InstructionAddressWord(),
        InstructionWord(),
        DoubleDataWord(),
    }


def test_bit(words: Set[Union[DataAddressWord, DataWord, InstructionAddressWord, InstructionWord]]):
    for word in words:
        for i in range(word.width()):
            assert not word.bit(i)

            word.set_bit(i)
            assert word.bit(i)

            word.clear_bit(i)
            assert not word.bit(i)


def test_bit_slice(words: Set[Union[DataAddressWord, DataWord, InstructionAddressWord, InstructionWord]]):
    for word in words:
        for _ in range(100):
            representation = IntInitializer.value_by_list(list(Representation))

            slice_width = IntInitializer.value_by_range(1, word.width())
            begin = IntInitializer.value_by_range(0, word.width() - slice_width + 1)
            end = IntInitializer.value_by_range(begin + slice_width, word.width() + 1)
            value = IntInitializer.value_by_width(representation, slice_width)

            word.set_bit_slice(begin, end, value)
            assert value == word.bit_slice(representation, begin, end)


def test_value(words: Set[Union[DataAddressWord, DataWord, InstructionAddressWord, InstructionWord]]):
    for word in words:
        for _ in range(100):
            representation = IntInitializer.value_by_list(list(Representation))
            value = IntInitializer.value_by_width(representation, word.width())
            word.set_value(value)
            assert value == word.value(representation)


def test_zero(words: Set[Union[DataAddressWord, DataWord, InstructionAddressWord, InstructionWord]]):
    for word in words:
        word.set_value(0)

        for i in range(word.width()):
            assert not word.bit(i)


def test_lneg(words: Set[Union[DataAddressWord, DataWord, InstructionAddressWord, InstructionWord]]):
    for word in words:
        word.set_value(-1)

        for i in range(word.width()):
            assert word.bit(i)


def test_bytes(words: Set[Union[DataAddressWord, DataWord, InstructionAddressWord, InstructionWord]]):
    for word in words:
        for _ in range(100):
            representation = IntInitializer.value_by_list(list(Representation))
            value = IntInitializer.value_by_width(representation, word.width())

            word.set_value(value)
            bytes_ = word.to_bytes()

            word.set_value(0)
            word.from_bytes(bytes_)

            assert value == word.value(representation)
