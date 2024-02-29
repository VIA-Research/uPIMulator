from abi.word._base_word import BaseWord
from util.config_loader import ConfigLoader


class InstructionAddressWord(BaseWord):
    def __init__(self):
        super().__init__(ConfigLoader.iram_address_width())
