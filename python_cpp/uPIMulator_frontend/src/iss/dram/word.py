from typing import List

from abi.word.data_word import DataWord
from abi.word.representation import Representation
from util.config_loader import ConfigLoader


class Word:
    def __init__(self, data_words: List[DataWord]):
        assert ConfigLoader.min_access_granularity() % DataWord().size() == 0

        self._data_words: List[DataWord] = [
            DataWord() for _ in range(ConfigLoader.min_access_granularity() // DataWord().size())
        ]

        for src_data_word, dst_data_word in zip(data_words, self._data_words):
            dst_data_word.set_value(src_data_word.value(Representation.UNSIGNED))

    def data_words(self) -> List[DataWord]:
        dst_data_words: List[DataWord] = [DataWord() for _ in range(len(self._data_words))]
        for src_data_word, dst_data_word in zip(self._data_words, dst_data_words):
            dst_data_word.set_value(src_data_word.value(Representation.UNSIGNED))
        return dst_data_words
