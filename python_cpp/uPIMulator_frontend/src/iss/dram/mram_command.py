from enum import Enum, auto
from typing import List

from abi.word.data_address_word import DataAddressWord
from abi.word.data_word import DataWord
from abi.word.representation import Representation
from util.config_loader import ConfigLoader


class MRAMCommand:
    class Operation(Enum):
        READ = 0
        WRITE = auto()

    def __init__(self, operation: Operation, address: int, size: int):
        assert size % ConfigLoader.min_access_granularity() == 0
        assert size % DataWord().size() == 0

        self._operation: MRAMCommand.Operation = operation

        self._address: DataAddressWord = DataAddressWord()
        self._address.set_value(address)

        self._size = size

        self._data_words: List[DataWord] = [DataWord() for _ in range(size // DataWord().size())]

    def operation(self) -> Operation:
        return self._operation

    def address(self) -> int:
        return self._address.value(Representation.UNSIGNED)

    def size(self) -> int:
        return self._size

    def begin_address(self) -> int:
        return self.address()

    def end_address(self) -> int:
        return self.address() + self.size()

    def data_words(self) -> List[DataWord]:
        data_words: List[DataWord] = [DataWord() for _ in range(len(self._data_words))]
        for i, data_word in enumerate(self._data_words):
            data_words[i].set_value(data_word.value(Representation.UNSIGNED))
        return data_words

    def data_word(self, address: int) -> DataWord:
        index = self._index(address)

        data_word = DataWord()
        data_word.set_value(self._data_words[index].value(Representation.UNSIGNED))
        return data_word

    def set_data_words(self, data_words: List[DataWord]) -> None:
        for i, data_word in enumerate(data_words):
            self._data_words[i].set_value(data_word.value(Representation.UNSIGNED))

    def set_data_word(self, address: int, data_word: DataWord) -> None:
        index = self._index(address)
        self._data_words[index].set_value(data_word.value(Representation.UNSIGNED))

    def _index(self, address: int) -> int:
        assert (address - self.address()) % DataWord().size() == 0
        return (address - self.address()) // DataWord().size()
