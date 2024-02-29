from abi.isa.exception import Exception_


class ExceptionRegister:
    def __init__(self):
        self._bits = [False for _ in range(len(Exception_))]

    def exception(self, exception: Exception_) -> bool:
        return self._bits[exception.value]

    def set_exception(self, exception: Exception_) -> None:
        self._bits[exception.value] = True

    def clear_exception(self, exception: Exception_) -> None:
        self._bits[exception.value] = False
