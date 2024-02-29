from abi.directive.long_directive import LongDirective
from abi.word.representation import Representation
from initializer.int_initializer import IntInitializer


def test_long_directive():
    for _ in range(100):
        value = IntInitializer.value_by_width(Representation.UNSIGNED, 32)

        long_directive = LongDirective(value)

        assert value == long_directive.value()
        assert 4 == long_directive.size()
