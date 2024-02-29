from typing import List

from abi.word.data_address_word import DataAddressWord
from abi.word.data_word import DataWord
from abi.word.representation import Representation
from util.config_loader import ConfigLoader


class WRAM:
    def __init__(self):
        self._address: DataAddressWord = DataAddressWord()
        self._address.set_value(ConfigLoader.wram_offset())

        self._data_words: List[DataWord] = [DataWord() for _ in range(ConfigLoader.wram_size() // DataWord().size())]

    def address(self) -> int:
        return self._address.value(Representation.UNSIGNED)

    def read(self, address: int) -> DataWord:
        index = self._index(address)
        data_word = DataWord()
        data_word.set_value(self._data_words[index].value(Representation.UNSIGNED))
        return data_word

    def write(self, address: int, data_word: DataWord) -> None:
        index = self._index(address)
        self._data_words[index].set_value(data_word.value(Representation.UNSIGNED))

    def _index(self, address: int) -> int:
        assert address % DataWord().size() == 0
        index = (address - self.address()) // DataWord().size()
        assert 0 <= index < len(self._data_words)

        return index
