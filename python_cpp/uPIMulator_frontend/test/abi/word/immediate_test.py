from abi.word.immediate import Immediate
from abi.word.representation import Representation
from initializer.int_initializer import IntInitializer


def test_value():
    for _ in range(100):
        representation = IntInitializer.value_by_list(list(Representation))
        width = IntInitializer.value_by_range(1, 64)
        value = IntInitializer.value_by_width(representation, width)
        immediate = Immediate(representation, width, value)

        assert value == immediate.value()
