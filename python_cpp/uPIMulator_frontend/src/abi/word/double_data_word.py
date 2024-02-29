from abi.word._base_word import BaseWord
from util.config_loader import ConfigLoader


class DoubleDataWord(BaseWord):
    def __init__(self):
        assert ConfigLoader.atomic_data_width() == ConfigLoader.wram_data_width() == ConfigLoader.mram_data_width()

        super().__init__(2 * ConfigLoader().atomic_data_width())
