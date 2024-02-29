import math
from typing import List

from abi.word.data_word import DataWord
from abi.word.instruction_word import InstructionWord
from abi.word.representation import Representation
from encoder.byte import Byte
from iss.dram.mram import MRAM
from iss.dram.mram_command import MRAMCommand
from iss.sram.atomic import Atomic
from iss.sram.iram import IRAM
from iss.sram.wram import WRAM
from util.config_loader import ConfigLoader


class DMA:
    def __init__(self, atomic: Atomic, iram: IRAM, wram: WRAM, mram: MRAM):
        self._atomic: Atomic = atomic
        self._iram: IRAM = iram
        self._wram: WRAM = wram
        self._mram: MRAM = mram

    def host_dma_transfer_to_atomic(self, address: int, bytes_: List[Byte]) -> None:
        for byte in bytes_:
            assert byte.value() == 0

    def host_dma_transfer_to_iram(self, address: int, bytes_: List[Byte]) -> None:
        assert len(bytes_) % InstructionWord().size() == 0
        num_instruction_words = math.ceil(len(bytes_) // InstructionWord().size())
        for i in range(num_instruction_words):
            instruction_word = InstructionWord()

            begin = instruction_word.size() * i
            end = instruction_word.size() * (i + 1)

            instruction_word.from_bytes(bytes_[begin:end])
            self._iram.write(address + begin, instruction_word)

    def host_dma_transfer_to_wram(self, address: int, bytes_: List[Byte]) -> None:
        for i in range(len(bytes_)):
            cur_address = address + i
            base_address = (cur_address // DataWord().size()) * DataWord().size()
            offset = cur_address % DataWord().size()

            data_word = self._wram.read(base_address)
            data_word.set_bit_slice(8 * offset, 8 * (offset + 1), bytes_[i].value())
            self._wram.write(base_address, data_word)

    def host_dma_transfer_from_wram(self, address: int, size: int) -> List[Byte]:
        bytes_: List[Byte] = []
        for i in range(size):
            cur_address = address + i
            base_address = (cur_address // DataWord().size()) * DataWord().size()
            offset = cur_address % DataWord().size()

            data_word = self._wram.read(base_address)
            bytes_.append(Byte(data_word.bit_slice(Representation.UNSIGNED, 8 * offset, 8 * (offset + 1))))
        return bytes_

    def host_dma_transfer_to_mram(self, address: int, bytes_: List[Byte]) -> None:
        num_data_words = math.ceil(len(bytes_) / DataWord().size())
        data_words: List[DataWord] = []
        for i in range(num_data_words):
            data_word = DataWord()

            begin = data_word.size() * i
            end = min(data_word.size() * (i + 1), len(bytes_))

            if end - begin < data_word.size():
                data_word.from_bytes(bytes_[begin:end] + [Byte(0) for _ in range(data_word.size() - end + begin)])
            else:
                data_word.from_bytes(bytes_[begin:end])
            data_words.append(data_word)

        while (len(data_words) * DataWord().size()) % ConfigLoader.min_access_granularity() != 0:
            data_words.append(DataWord())

        mram_command = MRAMCommand(MRAMCommand.Operation.WRITE, address, len(data_words) * DataWord().size())
        mram_command.set_data_words(data_words)

        assert self._mram.can_push()
        self._mram.push(mram_command)
        assert self._mram.can_pop()
        assert mram_command == self._mram.pop()

    def dpu_dma_transfer_from_mram_to_wram(self, src_address: int, dst_address: int, size: int) -> None:
        assert size % ConfigLoader.min_access_granularity() == 0
        assert size % DataWord().size() == 0

        mram_command = MRAMCommand(MRAMCommand.Operation.READ, src_address, size)

        assert self._mram.can_push()
        self._mram.push(mram_command)
        assert self._mram.can_pop()
        assert mram_command == self._mram.pop()

        for i, data_word in enumerate(mram_command.data_words()):
            self._wram.write(dst_address + i * data_word.size(), data_word)

    def dpu_dma_transfer_from_wram_to_mram(self, src_address: int, dst_address: int, size: int) -> None:
        assert size % ConfigLoader.min_access_granularity() == 0
        assert size % DataWord().size() == 0

        num_data_words = size // DataWord().size()
        data_words: List[DataWord] = []
        for i in range(num_data_words):
            data_word = self._wram.read(src_address + i * DataWord().size())
            data_words.append(data_word)

        mram_command = MRAMCommand(MRAMCommand.Operation.WRITE, dst_address, size)
        mram_command.set_data_words(data_words)

        assert self._mram.can_push()
        self._mram.push(mram_command)
        assert self._mram.can_pop()
        assert mram_command == self._mram.pop()
