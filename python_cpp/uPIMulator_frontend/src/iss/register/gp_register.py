from abi.isa.register.gp_register import GPRegister as SoftGPRegister
from abi.word.data_word import DataWord
from abi.word.representation import Representation


class GPRegister:
    def __init__(self, index: int):
        self._soft_gp_regsiter: SoftGPRegister = SoftGPRegister(index)
        self._word: DataWord = DataWord()

    def index(self) -> int:
        return self._soft_gp_regsiter.index()

    def read(self, representation: Representation) -> int:
        return self._word.value(representation)

    def write(self, value: int) -> None:
        self._word.set_value(value)
