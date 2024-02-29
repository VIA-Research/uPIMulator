from abi.directive.zero_directive import ZeroDirective
from abi.word.representation import Representation
from initializer.int_initializer import IntInitializer


def test_zero_directive():
    for _ in range(100):
        size = IntInitializer.value_by_range(1, 32)
        value = IntInitializer.value_by_width(Representation.UNSIGNED, 8)

        zero_directive = ZeroDirective(size, value)

        assert value == zero_directive.value()
        assert size == zero_directive.size()
