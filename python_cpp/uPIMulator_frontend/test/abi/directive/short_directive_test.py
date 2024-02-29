from abi.directive.short_directive import ShortDirective
from abi.word.representation import Representation
from initializer.int_initializer import IntInitializer


def test_short_directive():
    for _ in range(100):
        value = IntInitializer.value_by_width(Representation.UNSIGNED, 16)

        short_directive = ShortDirective(value)

        assert value == short_directive.value()
        assert 2 == short_directive.size()
