from enum import Enum, auto

from iss.register.register_file import RegisterFile


class Thread:
    class State(Enum):
        EMBRYO = 0
        RUNNABLE = auto()
        SLEEP = auto()
        ZOMBIE = auto()

    def __init__(self, id_: int):
        self._id: int = id_
        self._state: Thread.State = Thread.State.EMBRYO
        self._register_file: RegisterFile = RegisterFile(id_)
        self._issue_cycles: int = 0

    def id_(self) -> int:
        return self._id

    def state(self) -> State:
        return self._state

    def set_thread_state(self, thread_state: State) -> None:
        self._state = thread_state

    def register_file(self) -> RegisterFile:
        return self._register_file

    def issue_cycles(self) -> int:
        return self._issue_cycles

    def increment_issue_cycles(self) -> None:
        self._issue_cycles += 1

    def reset_issue_cycles(self) -> None:
        self._issue_cycles = 0
