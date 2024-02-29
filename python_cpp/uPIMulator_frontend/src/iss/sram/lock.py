from typing import Optional


class Lock:
    def __init__(self):
        self._id: Optional[int] = None

    def can_acquire(self):
        return self._id is None

    def acquire(self, id_: int):
        assert self.can_acquire()
        self._id = id_

    def can_release(self, id_: int):
        return self._id is None or self._id == id_

    def release(self, id_: int):
        assert self.can_release(id_)
        self._id = None
