import pytest

from abi.isa.exception import Exception_
from abi.isa.flag import Flag
from abi.isa.instruction.condition import Condition
from abi.isa.register.gp_register import GPRegister
from abi.isa.register.pair_register import PairRegister
from abi.isa.register.sp_register import SPRegister
from abi.word.data_word import DataWord
from abi.word.double_data_word import DoubleDataWord
from abi.word.instruction_address_word import InstructionAddressWord
from abi.word.representation import Representation
from initializer.int_initializer import IntInitializer
from iss.register.register_file import RegisterFile
from util.config_loader import ConfigLoader


@pytest.fixture
def id_() -> int:
    return IntInitializer.value_by_range(0, ConfigLoader.max_num_tasklets())


@pytest.fixture
def instruction_address_word_value() -> int:
    return IntInitializer.value_by_width(Representation.UNSIGNED, InstructionAddressWord().width())


@pytest.fixture
def data_word_value() -> int:
    return IntInitializer.value_by_width(Representation.SIGNED, DataWord().width())


@pytest.fixture
def double_data_word_value() -> int:
    return IntInitializer.value_by_width(Representation.SIGNED, DoubleDataWord().width())


def test_gp_register(id_: int, data_word_value: int):
    register_file = RegisterFile(id_)
    for i in range(ConfigLoader.num_gp_registers()):
        gp_register = GPRegister(i)

        register_file.write(gp_register, data_word_value)
        assert data_word_value == register_file.read(gp_register, Representation.SIGNED)

        register_file.cycle()


def test_sp_register(id_: int):
    register_file = RegisterFile(id_)

    assert register_file.read(SPRegister.ZERO, Representation.SIGNED) == 0
    assert register_file.read(SPRegister.ONE, Representation.SIGNED) == 1
    assert register_file.read(SPRegister.ID, Representation.SIGNED) == id_
    assert register_file.read(SPRegister.LNEG, Representation.SIGNED) == -1
    assert register_file.read(SPRegister.MNEG, Representation.SIGNED) == -(2 ** (DataWord().width() - 1))
    assert register_file.read(SPRegister.ID2, Representation.SIGNED) == (2 * id_)
    assert register_file.read(SPRegister.ID4, Representation.SIGNED) == (4 * id_)
    assert register_file.read(SPRegister.ID8, Representation.SIGNED) == (8 * id_)


def test_pair_register(id_: int, double_data_word_value: int):
    register_file = RegisterFile(id_)

    for i in range(ConfigLoader.num_gp_registers()):
        if i % 2 == 0:
            pair_register = PairRegister(i)

            register_file.write(pair_register, double_data_word_value)
            assert double_data_word_value == register_file.read(pair_register, Representation.SIGNED)

            register_file.cycle()


def test_pc_register(id_: int, instruction_address_word_value: int):
    register_file = RegisterFile(id_)

    register_file.write_pc(instruction_address_word_value)
    assert instruction_address_word_value == register_file.read_pc()

    register_file.increment_pc()
    assert (instruction_address_word_value + InstructionAddressWord().size()) == register_file.read_pc()


def test_condition_register(id_: int):
    register_file = RegisterFile(id_)

    for condition in Condition:
        if condition == Condition.TRUE:
            assert register_file.condition(condition)
        elif condition == Condition.FALSE:
            assert not register_file.condition(condition)
        elif (
            condition == Condition.Z
            or condition == Condition.NZ
            or condition == Condition.C
            or condition == Condition.NC
        ):
            pass
        else:
            assert not register_file.condition(condition)

            register_file.set_condition(condition)
            assert register_file.condition(condition)

            register_file.clear_condition(condition)
            assert not register_file.condition(condition)


def test_exception_register(id_: int):
    register_file = RegisterFile(id_)

    for exception in Exception_:
        assert not register_file.exception(exception)

        register_file.set_exception(exception)
        assert register_file.exception(exception)

        register_file.clear_exception(exception)
        assert not register_file.exception(exception)


def test_flag_register(id_: int):
    register_file = RegisterFile(id_)

    for flag in Flag:
        assert not register_file.flag(flag)

        register_file.set_flag(flag)
        assert register_file.flag(flag)

        register_file.clear_flag(flag)
        assert not register_file.flag(flag)


def test_clear_all(id_: int):
    register_file = RegisterFile(id_)

    for condition in Condition:
        if (
            condition == Condition.TRUE
            or condition == Condition.FALSE
            or condition == Condition.Z
            or condition == Condition.NZ
            or condition == Condition.C
            or condition == Condition.NC
        ):
            pass
        else:
            assert not register_file.condition(condition)

    for exception in Exception_:
        assert not register_file.exception(exception)

    for flag in Flag:
        assert not register_file.flag(flag)
