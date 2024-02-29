from __future__ import annotations

from abi.isa.register.gp_register import GPRegister
from util.config_loader import ConfigLoader


class PairRegister:
    def __init__(self, index: int):
        assert 0 <= index < ConfigLoader.num_gp_registers()
        assert index % 2 == 0

        self._index: int = index

    def __eq__(self, other: PairRegister) -> bool:
        return self._index == other.index()

    def index(self) -> int:
        return self._index

    def even_register(self) -> GPRegister:
        return GPRegister(self._index)

    def odd_register(self) -> GPRegister:
        return GPRegister(self._index + 1)
