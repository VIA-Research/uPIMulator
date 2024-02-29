from abi.isa.register.sp_register import SPRegister as SoftSPRegister
from abi.word.data_word import DataWord
from abi.word.representation import Representation
from util.config_loader import ConfigLoader


class SPRegister:
    def __init__(self, id_: int):
        assert 0 <= id_ < ConfigLoader.max_num_tasklets()

        self._zero = DataWord()
        self._zero.set_value(0)

        self._one = DataWord()
        self._one.set_value(1)

        self._lneg = DataWord()
        self._lneg.set_value(-1)

        self._mneg = DataWord()
        self._mneg.set_bit(self._mneg.width() - 1)

        self._id = DataWord()
        self._id.set_value(id_)

        self._id2 = DataWord()
        self._id2.set_value(2 * id_)

        self._id4 = DataWord()
        self._id4.set_value(4 * id_)

        self._id8 = DataWord()
        self._id8.set_value(8 * id_)

    def read(self, soft_sp_register: SoftSPRegister, representation: Representation) -> int:
        if soft_sp_register == SoftSPRegister.ZERO:
            return self._zero.value(representation)
        elif soft_sp_register == SoftSPRegister.ONE:
            return self._one.value(representation)
        elif soft_sp_register == SoftSPRegister.LNEG:
            return self._lneg.value(representation)
        elif soft_sp_register == SoftSPRegister.MNEG:
            return self._mneg.value(representation)
        elif soft_sp_register == SoftSPRegister.ID:
            return self._id.value(representation)
        elif soft_sp_register == SoftSPRegister.ID2:
            return self._id2.value(representation)
        elif soft_sp_register == SoftSPRegister.ID4:
            return self._id4.value(representation)
        elif soft_sp_register == SoftSPRegister.ID8:
            return self._id8.value(representation)
        else:
            raise ValueError
