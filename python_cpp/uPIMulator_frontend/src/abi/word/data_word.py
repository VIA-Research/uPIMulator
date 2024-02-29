from abi.word._base_word import BaseWord
from util.config_loader import ConfigLoader


class DataWord(BaseWord):
    def __init__(self):
        assert ConfigLoader.atomic_data_width() == ConfigLoader.wram_data_width() == ConfigLoader.mram_data_width()

        super().__init__(ConfigLoader().atomic_data_width())
