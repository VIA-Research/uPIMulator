from queue import Queue
from typing import List

from abi.word.data_address_word import DataAddressWord
from abi.word.data_word import DataWord
from abi.word.representation import Representation
from iss.dram.mram_command import MRAMCommand
from util.config_loader import ConfigLoader


class MRAM:
    def __init__(self):
        self._address: DataAddressWord = DataAddressWord()
        self._address.set_value(ConfigLoader.mram_offset())

        self._data_words: List[DataWord] = [DataWord() for _ in range(ConfigLoader.mram_size() // DataWord().size())]
        self._mram_command_queue: Queue[MRAMCommand] = Queue(-1)

        assert ConfigLoader.min_access_granularity() % DataWord().size() == 0
        self._num_min_access_data_words = ConfigLoader.min_access_granularity() // DataWord().size()

    def address(self) -> int:
        return self._address.value(Representation.UNSIGNED)

    def can_push(self) -> bool:
        return True

    def push(self, mram_command: MRAMCommand) -> None:
        self._mram_command_queue.put(mram_command)

    def can_pop(self) -> bool:
        return not self._mram_command_queue.empty()

    def pop(self) -> MRAMCommand:
        mram_command = self._mram_command_queue.get()
        if mram_command.operation() == MRAMCommand.Operation.READ:
            if mram_command.address() % DataWord().size() == 0:
                self._aligned_read(mram_command)
            else:
                self._unaligned_read(mram_command)
        elif mram_command.operation() == MRAMCommand.Operation.WRITE:
            if mram_command.address() % DataWord().size() == 0:
                self._aligned_write(mram_command)
            else:
                self._unaligned_write(mram_command)
        else:
            raise ValueError
        return mram_command

    def _aligned_read(self, mram_command: MRAMCommand) -> None:
        assert mram_command.address() % DataWord().size() == 0

        begin = self._index(mram_command.begin_address())
        end = self._index(mram_command.end_address())
        mram_command.set_data_words(self._data_words[begin:end])

    def _aligned_write(self, mram_command: MRAMCommand) -> None:
        assert mram_command.address() % DataWord().size() == 0

        begin = self._index(mram_command.address())
        end = self._index(mram_command.address() + mram_command.size())
        self._data_words[begin:end] = mram_command.data_words()

    def _unaligned_read(self, mram_command: MRAMCommand) -> None:
        data_words: List[DataWord] = [DataWord() for _ in range(mram_command.size() // DataWord().size())]
        for address in range(mram_command.begin_address(), mram_command.end_address()):
            mram_index = self._index(address)
            mram_offset = self._offset(address)

            byte = self._data_words[mram_index].bit_slice(
                Representation.UNSIGNED, 8 * mram_offset, 8 * (mram_offset + 1)
            )

            mram_command_index = (address - mram_command.begin_address()) // DataWord().size()
            mram_command_offset = (address - mram_command.begin_address()) % DataWord().size()

            data_words[mram_command_index].set_bit_slice(8 * mram_command_offset, 8 * (mram_command_offset + 1), byte)

        mram_command.set_data_words(data_words)

    def _unaligned_write(self, mram_command: MRAMCommand) -> None:
        for address in range(mram_command.begin_address(), mram_command.end_address()):
            mram_command_index = (address - mram_command.begin_address()) // DataWord().size()
            mram_command_offset = (address - mram_command.begin_address()) % DataWord().size()

            byte = mram_command.data_words()[mram_command_index].bit_slice(
                Representation.UNSIGNED, 8 * mram_command_offset, 8 * (mram_command_offset + 1),
            )

            mram_index = self._index(address)
            mram_offset = self._offset(address)

            self._data_words[mram_index].set_bit_slice(8 * mram_offset, 8 * (mram_offset + 1), byte)

    def cycle(self) -> None:
        pass

    def _index(self, address: int) -> int:
        index = (address - self.address()) // DataWord().size()
        assert 0 <= index < index + self._num_min_access_data_words <= len(self._data_words)
        return index

    def _offset(self, address: int) -> int:
        offset = (address - self.address()) % DataWord().size()
        return offset
