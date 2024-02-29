from abi.directive.byte_directive import ByteDirective
from abi.word.representation import Representation
from initializer.int_initializer import IntInitializer


def test_byte_directive():
    for _ in range(100):
        value = IntInitializer.value_by_width(Representation.UNSIGNED, 8)

        byte_directive = ByteDirective(value)

        assert value == byte_directive.value()
        assert 1 == byte_directive.size()
