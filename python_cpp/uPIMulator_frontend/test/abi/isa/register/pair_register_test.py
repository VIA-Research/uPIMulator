from abi.isa.register.pair_register import PairRegister
from initializer.int_initializer import IntInitializer
from util.config_loader import ConfigLoader


def test_pair_register():
    index = (IntInitializer.value_by_range(0, ConfigLoader.num_gp_registers()) // 2) * 2
    pair_register = PairRegister(index)

    assert index == pair_register.index()
