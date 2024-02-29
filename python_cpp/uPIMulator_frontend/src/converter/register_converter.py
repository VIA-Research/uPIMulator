from typing import Union

from abi.isa.register.gp_register import GPRegister
from abi.isa.register.pair_register import PairRegister
from abi.isa.register.sp_register import SPRegister


class RegisterConverter:
    def __init__(self):
        pass

    @staticmethod
    def convert_to_gp_register(gp_register: str) -> GPRegister:
        assert gp_register[0] == "r"
        index = int(gp_register[1:])
        return GPRegister(index)

    @staticmethod
    def convert_to_sp_register(sp_register: str) -> SPRegister:
        if sp_register == "zero":
            return SPRegister.ZERO
        elif sp_register == "one":
            return SPRegister.ONE
        elif sp_register == "lneg":
            return SPRegister.LNEG
        elif sp_register == "mneg":
            return SPRegister.MNEG
        elif sp_register == "id":
            return SPRegister.ID
        elif sp_register == "id2":
            return SPRegister.ID2
        elif sp_register == "id4":
            return SPRegister.ID4
        elif sp_register == "id8":
            return SPRegister.ID8
        else:
            raise ValueError

    @staticmethod
    def convert_to_pair_register(pair_register: str) -> PairRegister:
        assert pair_register[0] == "d"
        index = int(pair_register[1:])
        return PairRegister(index)

    @staticmethod
    def convert_to_zero_register(zero_register: str) -> SPRegister:
        assert zero_register == "zero"
        return SPRegister.ZERO

    @staticmethod
    def convert_to_source_register(source_register: str,) -> Union[GPRegister, SPRegister]:
        if source_register[0] == "r":
            return RegisterConverter.convert_to_gp_register(source_register)
        else:
            return RegisterConverter.convert_to_sp_register(source_register)

    @staticmethod
    def convert_to_string(register: Union[GPRegister, SPRegister, PairRegister]) -> str:
        if isinstance(register, GPRegister):
            return f"r{register.index()}"
        elif isinstance(register, PairRegister):
            return f"d{register.index()}"
        elif register == SPRegister.ZERO:
            return "zero"
        elif register == SPRegister.ONE:
            return "one"
        elif register == SPRegister.LNEG:
            return "lneg"
        elif register == SPRegister.MNEG:
            return "mneg"
        elif register == SPRegister.ID:
            return "id"
        elif register == SPRegister.ID2:
            return "id2"
        elif register == SPRegister.ID4:
            return "id4"
        elif register == SPRegister.ID8:
            return "id8"
        else:
            raise ValueError
