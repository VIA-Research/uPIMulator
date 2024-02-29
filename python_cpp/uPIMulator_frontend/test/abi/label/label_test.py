import pytest

from abi.label.label import Label
from abi.word.data_address_word import DataAddressWord
from abi.word.instruction_address_word import InstructionAddressWord
from abi.word.representation import Representation
from initializer.int_initializer import IntInitializer
from initializer.str_initializer import StrInitializer


@pytest.fixture
def label() -> Label:
    return Label("foo")


def test_name():
    for _ in range(100):
        name_width = IntInitializer.value_by_range(1, 64)
        name = StrInitializer.identifier(name_width)

        label = Label(name)

        assert name == label.name()


def test_address(label: Label):
    address_width = min(InstructionAddressWord().width(), DataAddressWord().width())
    address = IntInitializer.value_by_width(Representation.UNSIGNED, address_width)

    label.set_address(address)

    assert address == label.address()


def test_size(label: Label):
    address_width = min(InstructionAddressWord().width(), DataAddressWord().width())
    max_size = 2 ** address_width - 1
    size = IntInitializer.value_by_range(0, max_size + 1)

    label.set_size(size)

    assert size == label.size()


def test_begin_address(label: Label):
    address_width = min(InstructionAddressWord().width(), DataAddressWord().width())
    address = IntInitializer.value_by_width(Representation.UNSIGNED, address_width)

    label.set_address(address)

    assert address == label.begin_address()


def test_end_address(label: Label):
    address_width = min(InstructionAddressWord().width(), DataAddressWord().width())
    address = IntInitializer.value_by_width(Representation.UNSIGNED, address_width)
    max_size = 2 ** address_width - 1
    size = IntInitializer.value_by_range(0, max_size + 1)

    label.set_address(address)
    label.set_size(size)

    assert address + size == label.end_address()
