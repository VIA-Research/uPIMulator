from abi.isa.flag import Flag


class FlagRegister:
    def __init__(self):
        self._bits = [False for _ in range(len(Flag))]

    def flag(self, flag: Flag) -> bool:
        return self._bits[flag.value]

    def set_flag(self, flag: Flag) -> None:
        self._bits[flag.value] = True

    def clear_flag(self, flag: Flag) -> None:
        self._bits[flag.value] = False
