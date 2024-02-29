import math
from typing import List

import pytest

from abi.word.data_word import DataWord
from abi.word.instruction_word import InstructionWord
from abi.word.representation import Representation
from encoder.byte import Byte
from initializer.int_initializer import IntInitializer
from iss.dpu.dma import DMA
from iss.dram.mram import MRAM
from iss.dram.mram_command import MRAMCommand
from iss.sram.atomic import Atomic
from iss.sram.iram import IRAM
from iss.sram.wram import WRAM
from util.config_loader import ConfigLoader


@pytest.fixture
def atomic() -> Atomic:
    return Atomic()


@pytest.fixture
def iram() -> IRAM:
    return IRAM()


@pytest.fixture
def wram() -> WRAM:
    return WRAM()


@pytest.fixture
def mram() -> MRAM:
    return MRAM()


@pytest.fixture
def atomic_address() -> int:
    return ConfigLoader.atomic_offset()


@pytest.fixture
def atomic_bytes() -> List[Byte]:
    return [Byte(0) for _ in range(100)]


@pytest.fixture
def iram_address() -> int:
    return ConfigLoader.iram_offset()


@pytest.fixture
def iram_bytes() -> List[Byte]:
    bytes_: List[Byte] = []
    for _ in range(100):
        instruction_word = InstructionWord()
        instruction_word.set_value(IntInitializer.value_by_width(Representation.UNSIGNED, instruction_word.width()))
        bytes_ += instruction_word.to_bytes()
    return bytes_


@pytest.fixture
def wram_address() -> int:
    return ConfigLoader.wram_offset()


@pytest.fixture
def wram_bytes() -> List[Byte]:
    return [Byte(IntInitializer.value_by_range(0, 2 ** 8)) for _ in range(100)]


@pytest.fixture
def mram_address() -> int:
    return ConfigLoader.mram_offset()


@pytest.fixture
def unaligned_mram_bytes() -> List[Byte]:
    return [Byte(IntInitializer.value_by_range(0, 2 ** 8)) for _ in range(100)]


@pytest.fixture
def aligned_mram_bytes() -> List[Byte]:
    return [Byte(IntInitializer.value_by_range(0, 2 ** 8)) for _ in range(100 * ConfigLoader.min_access_granularity())]


def test_host_dma_transfer_to_atomic(
    atomic: Atomic, iram: IRAM, wram: WRAM, mram: MRAM, atomic_address: int, atomic_bytes: List[Byte],
):
    dma = DMA(atomic, iram, wram, mram)
    dma.host_dma_transfer_to_atomic(atomic_address, atomic_bytes)

    for address in range(atomic_address, atomic_address + len(atomic_bytes)):
        assert atomic.can_acquire(address)


def test_host_dma_transfer_to_iram(
    atomic: Atomic, iram: IRAM, wram: WRAM, mram: MRAM, iram_address: int, iram_bytes: List[Byte],
):
    dma = DMA(atomic, iram, wram, mram)
    dma.host_dma_transfer_to_iram(iram_address, iram_bytes)

    num_instruction_word = len(iram_bytes) // InstructionWord().size()
    bytes_: List[Byte] = []
    for i in range(num_instruction_word):
        instruction_word = iram.read(iram_address + i * InstructionWord().size())
        bytes_ += instruction_word.to_bytes()

    for iram_byte, byte in zip(iram_bytes, bytes_):
        assert iram_byte.value() == byte.value()


def test_host_dma_transfer_to_wram(
    atomic: Atomic, iram: IRAM, wram: WRAM, mram: MRAM, wram_address: int, wram_bytes: List[Byte],
):
    dma = DMA(atomic, iram, wram, mram)
    dma.host_dma_transfer_to_wram(wram_address, wram_bytes)

    num_data_words = math.ceil(len(wram_bytes) / DataWord().size())
    bytes_: List[Byte] = []
    for i in range(num_data_words):
        data_word = wram.read(wram_address + i * DataWord().size())
        bytes_ += data_word.to_bytes()

    for i in range(len(wram_bytes)):
        assert wram_bytes[i].value() == bytes_[i].value()


def test_host_dma_transfer_to_mram(
    atomic: Atomic, iram: IRAM, wram: WRAM, mram: MRAM, mram_address: int, unaligned_mram_bytes: List[Byte],
):
    dma = DMA(atomic, iram, wram, mram)
    dma.host_dma_transfer_to_mram(mram_address, unaligned_mram_bytes)

    mram_command_size = (
        math.ceil(len(unaligned_mram_bytes) / ConfigLoader.min_access_granularity())
        * ConfigLoader.min_access_granularity()
    )
    mram_command = MRAMCommand(MRAMCommand.Operation.READ, mram_address, mram_command_size)

    assert mram.can_push()
    mram.push(mram_command)
    assert mram.can_pop()
    mram_command_bytes: List[Byte] = []
    for data_word in mram.pop().data_words():
        mram_command_bytes += data_word.to_bytes()

    for i in range(len(unaligned_mram_bytes)):
        assert unaligned_mram_bytes[i].value() == mram_command_bytes[i].value()


def test_dpu_dma_transfer_from_mram_to_wram(
    atomic: Atomic,
    iram: IRAM,
    wram: WRAM,
    mram: MRAM,
    wram_address: int,
    mram_address: int,
    aligned_mram_bytes: List[Byte],
):
    dma = DMA(atomic, iram, wram, mram)
    dma.host_dma_transfer_to_mram(mram_address, aligned_mram_bytes)
    dma.dpu_dma_transfer_from_mram_to_wram(mram_address, wram_address, len(aligned_mram_bytes))

    num_data_words = len(aligned_mram_bytes) // DataWord().size()
    bytes_: List[Byte] = []
    for i in range(num_data_words):
        data_word = wram.read(wram_address + i * DataWord().size())
        bytes_ += data_word.to_bytes()

    for i in range(len(aligned_mram_bytes)):
        assert aligned_mram_bytes[i].value() == bytes_[i].value()


def test_dpu_dma_transfer_from_wram_to_mram(
    atomic: Atomic,
    iram: IRAM,
    wram: WRAM,
    mram: MRAM,
    wram_address: int,
    mram_address: int,
    aligned_mram_bytes: List[Byte],
):
    dma = DMA(atomic, iram, wram, mram)
    dma.host_dma_transfer_to_wram(wram_address, aligned_mram_bytes)
    dma.dpu_dma_transfer_from_wram_to_mram(wram_address, mram_address, len(aligned_mram_bytes))

    mram_command_size = (
        len(aligned_mram_bytes) // ConfigLoader.min_access_granularity()
    ) * ConfigLoader.min_access_granularity()

    mram_command = MRAMCommand(MRAMCommand.Operation.READ, mram_address, mram_command_size)

    assert mram.can_push()
    mram.push(mram_command)
    assert mram.can_pop()
    mram_command_bytes: List[Byte] = []
    for data_word in mram.pop().data_words():
        mram_command_bytes += data_word.to_bytes()

    for i in range(len(aligned_mram_bytes)):
        assert aligned_mram_bytes[i].value() == mram_command_bytes[i].value()
