from __future__ import annotations

from util.config_loader import ConfigLoader


class GPRegister:
    def __init__(self, index: int):
        assert 0 <= index < ConfigLoader.num_gp_registers()

        self._index: int = index

    def __eq__(self, other: GPRegister) -> bool:
        return self._index == other.index()

    def index(self) -> int:
        return self._index
