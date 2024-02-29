from abi.word._base_word import BaseWord
from util.config_loader import ConfigLoader


class DataAddressWord(BaseWord):
    def __init__(self):
        assert (
            ConfigLoader.atomic_address_width()
            == ConfigLoader.iram_address_width()
            == ConfigLoader.wram_address_width()
            == ConfigLoader.mram_address_width()
        )

        super().__init__(ConfigLoader.mram_address_width())
