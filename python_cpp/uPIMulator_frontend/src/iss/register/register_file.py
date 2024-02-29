from typing import Union

from abi.isa.exception import Exception_
from abi.isa.flag import Flag
from abi.isa.instruction.condition import Condition
from abi.isa.register.gp_register import GPRegister as SoftGPRegister
from abi.isa.register.pair_register import PairRegister as SoftPairRegister
from abi.isa.register.sp_register import SPRegister as SoftSPRegister
from abi.word.double_data_word import DoubleDataWord
from abi.word.representation import Representation
from iss.register.condition_register import ConditionRegister as HardConditionRegister
from iss.register.exception_register import ExceptionRegister as HardExceptionRegister
from iss.register.flag_register import FlagRegister as HardFlagRegister
from iss.register.gp_register import GPRegister as HardGPRegister
from iss.register.pc_register import PCRegister as HardPCRegister
from iss.register.sp_register import SPRegister as HardSPRegister
from util.config_loader import ConfigLoader


class RegisterFile:
    Register = Union[SoftGPRegister, SoftSPRegister, SoftPairRegister]

    def __init__(self, id_: int):
        self._gp_registers = [HardGPRegister(i) for i in range(ConfigLoader.num_gp_registers())]
        self._sp_register = HardSPRegister(id_)
        self._pc_register = HardPCRegister()
        self._condition_register = HardConditionRegister()
        self._exception_register = HardExceptionRegister()
        self._flag_register = HardFlagRegister()

    def read(self, register: Register, representation: Representation) -> int:
        if isinstance(register, SoftGPRegister):
            return self._read_gp_register(register, representation)
        elif isinstance(register, SoftSPRegister):
            return self._read_sp_register(register, representation)
        elif isinstance(register, SoftPairRegister):
            return self._read_pair_register(register, representation)
        else:
            raise ValueError

    def write(self, register: Register, value: int) -> None:
        if isinstance(register, SoftGPRegister):
            self._write_gp_register(register, value)
        elif isinstance(register, SoftSPRegister):
            raise ValueError
        elif isinstance(register, SoftPairRegister):
            self._write_pair_register(register, value)
        else:
            raise ValueError

    def read_pc(self) -> int:
        return self._pc_register.read()

    def write_pc(self, value: int) -> None:
        self._pc_register.write(value)

    def increment_pc(self) -> None:
        self._pc_register.increment()

    def condition(self, condition: Condition) -> bool:
        return self._condition_register.condition(condition)

    def set_condition(self, condition: Condition) -> None:
        self._condition_register.set_condition(condition)

    def clear_condition(self, condition: Condition) -> None:
        self._condition_register.clear_condition(condition)

    def exception(self, exception: Exception_) -> bool:
        return self._exception_register.exception(exception)

    def set_exception(self, exception: Exception_) -> None:
        self._exception_register.set_exception(exception)

    def clear_exception(self, exception: Exception_) -> None:
        self._exception_register.clear_exception(exception)

    def flag(self, flag: Flag) -> bool:
        return self._flag_register.flag(flag)

    def set_flag(self, flag: Flag) -> None:
        self._flag_register.set_flag(flag)

    def clear_flag(self, flag: Flag) -> None:
        self._flag_register.clear_flag(flag)

    def clear_conditions(self) -> None:
        for condition in Condition:
            if condition == Condition.TRUE or condition == Condition.FALSE:
                pass
            else:
                self.clear_condition(condition)

    def cycle(self) -> None:
        pass

    def _read_gp_register(self, soft_gp_register: SoftGPRegister, representation: Representation) -> int:
        return self._gp_registers[soft_gp_register.index()].read(representation)

    def _read_sp_register(self, soft_sp_register: SoftSPRegister, representation: Representation):
        return self._sp_register.read(soft_sp_register, representation)

    def _read_pair_register(self, soft_pair_register: SoftPairRegister, representation: Representation) -> int:
        even_register = self._gp_registers[soft_pair_register.even_register().index()]
        odd_register = self._gp_registers[soft_pair_register.odd_register().index()]

        double_data_word = DoubleDataWord()
        double_data_word.set_bit_slice(0, double_data_word.width() // 2, odd_register.read(Representation.UNSIGNED))
        double_data_word.set_bit_slice(
            double_data_word.width() // 2, double_data_word.width(), even_register.read(representation),
        )

        return double_data_word.value(representation)

    def _write_gp_register(self, soft_gp_register: SoftGPRegister, value: int) -> None:
        self._gp_registers[soft_gp_register.index()].write(value)

    def _write_pair_register(self, soft_pair_register: SoftPairRegister, value: int) -> None:
        even_register = self._gp_registers[soft_pair_register.even_register().index()]
        odd_register = self._gp_registers[soft_pair_register.odd_register().index()]

        double_data_word = DoubleDataWord()
        double_data_word.set_value(value)

        odd_register.write(double_data_word.bit_slice(Representation.UNSIGNED, 0, double_data_word.width() // 2))
        even_register.write(
            double_data_word.bit_slice(
                Representation.UNSIGNED, double_data_word.width() // 2, double_data_word.width(),
            )
        )
