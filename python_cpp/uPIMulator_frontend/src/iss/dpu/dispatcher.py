from abi.word.data_word import DataWord
from abi.word.double_data_word import DoubleDataWord
from abi.word.representation import Representation
from iss.sram.wram import WRAM


class Dispatcher:
    def __init__(self, wram: WRAM):
        self._wram: WRAM = wram

    def lbs(self, address: int) -> int:
        base_address = (address // DataWord().size()) * DataWord().size()
        offset = address % DataWord().size()

        data_word = DataWord()
        data_word.set_value(self._wram.read(base_address).value(Representation.SIGNED))
        return data_word.bit_slice(Representation.SIGNED, 8 * offset, 8 * (offset + 1))

    def lbu(self, address: int) -> int:
        base_address = (address // DataWord().size()) * DataWord().size()
        offset = address % DataWord().size()

        data_word = DataWord()
        data_word.set_value(self._wram.read(base_address).value(Representation.UNSIGNED))
        return data_word.bit_slice(Representation.UNSIGNED, 8 * offset, 8 * (offset + 1))

    def lhs(self, address: int) -> int:
        data_word = DataWord()
        data_word.set_bit_slice(0, 8, self.lbs(address))
        data_word.set_bit_slice(8, 16, self.lbs(address + 1))
        return data_word.bit_slice(Representation.SIGNED, 0, 16)

    def lhu(self, address: int) -> int:
        data_word = DataWord()
        data_word.set_bit_slice(0, 8, self.lbu(address))
        data_word.set_bit_slice(8, 16, self.lbu(address + 1))
        return data_word.bit_slice(Representation.UNSIGNED, 0, 16)

    def lw(self, address: int) -> int:
        data_word = DataWord()
        data_word.set_bit_slice(0, 8, self.lbu(address))
        data_word.set_bit_slice(8, 16, self.lbu(address + 1))
        data_word.set_bit_slice(16, 24, self.lbu(address + 2))
        data_word.set_bit_slice(24, 32, self.lbu(address + 3))
        return data_word.value(Representation.UNSIGNED)

    def ld(self, address: int) -> int:
        double_data_word = DoubleDataWord()
        double_data_word.set_bit_slice(0, DataWord().width(), self.lw(address))
        double_data_word.set_bit_slice(
            DataWord().width(), 2 * DataWord().width(), self.lw(address + DataWord().size()),
        )
        return double_data_word.value(Representation.UNSIGNED)

    def sb(self, address: int, value: int) -> None:
        base_address = (address // DataWord().size()) * DataWord().size()
        offset = address % DataWord().size()

        data_word = self._wram.read(base_address)
        data_word.set_bit_slice(8 * offset, 8 * (offset + 1), value)

        self._wram.write(base_address, data_word)

    def sh(self, address: int, value: int) -> None:
        value_word = DataWord()
        value_word.set_value(value)

        self.sb(address, value_word.bit_slice(Representation.UNSIGNED, 0, 8))
        self.sb(address + 1, value_word.bit_slice(Representation.UNSIGNED, 8, 16))

    def sw(self, address: int, value: int) -> None:
        value_word = DataWord()
        value_word.set_value(value)

        self.sb(address, value_word.bit_slice(Representation.UNSIGNED, 0, 8))
        self.sb(address + 1, value_word.bit_slice(Representation.UNSIGNED, 8, 16))
        self.sb(address + 2, value_word.bit_slice(Representation.UNSIGNED, 16, 24))
        self.sb(address + 3, value_word.bit_slice(Representation.UNSIGNED, 24, 32))

    def sd(self, address: int, value: int) -> None:
        value_word = DoubleDataWord()
        value_word.set_value(value)

        self.sw(
            address, value_word.bit_slice(Representation.UNSIGNED, 0, DataWord().width()),
        )
        self.sw(
            address + DataWord().size(),
            value_word.bit_slice(Representation.UNSIGNED, DataWord().width(), 2 * DataWord().width()),
        )
