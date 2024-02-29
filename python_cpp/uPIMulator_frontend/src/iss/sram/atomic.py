from typing import List

from abi.word.data_address_word import DataAddressWord
from abi.word.representation import Representation
from iss.sram.lock import Lock
from util.config_loader import ConfigLoader


class Atomic:
    def __init__(self):
        self._address: DataAddressWord = DataAddressWord()
        self._address.set_value(ConfigLoader.atomic_offset())

        self._locks: List[Lock] = [Lock() for _ in range(ConfigLoader.atomic_size())]

    def address(self) -> int:
        return self._address.value(Representation.UNSIGNED)

    def can_acquire(self, address: int) -> bool:
        index = self._index(address)
        return self._locks[index].can_acquire()

    def acquire(self, address: int, id_: int) -> None:
        index = self._index(address)
        self._locks[index].acquire(id_)

    def can_release(self, address: int, id_: int) -> bool:
        index = self._index(address)
        return self._locks[index].can_release(id_)

    def release(self, address: int, id_: int) -> None:
        index = self._index(address)
        self._locks[index].release(id_)

    def _index(self, address: int) -> int:
        index = address - self.address()
        assert 0 <= index < len(self._locks)
        return index
