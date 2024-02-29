from abi.directive.quad_directive import QuadDirective
from abi.word.representation import Representation
from initializer.int_initializer import IntInitializer


def test_quad_directive():
    for _ in range(100):
        value = IntInitializer.value_by_width(Representation.UNSIGNED, 64)

        quad_directive = QuadDirective(value)

        assert value == quad_directive.value()
        assert 8 == quad_directive.size()
