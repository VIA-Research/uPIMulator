from abi.isa.register.gp_register import GPRegister
from initializer.int_initializer import IntInitializer
from util.config_loader import ConfigLoader


def test_gp_register():
    index = IntInitializer.value_by_range(0, ConfigLoader.num_gp_registers())
    gp_register = GPRegister(index)

    assert index == gp_register.index()
