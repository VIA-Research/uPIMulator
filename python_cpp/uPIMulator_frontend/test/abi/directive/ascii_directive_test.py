from abi.directive.ascii_directive import AsciiDirective
from initializer.int_initializer import IntInitializer
from initializer.str_initializer import StrInitializer


def test_ascii_directive():
    for _ in range(100):
        characters_width = IntInitializer.value_by_range(1, 64)
        characters = StrInitializer.identifier(characters_width)

        ascii_directive = AsciiDirective(characters)

        assert characters == ascii_directive.characters()
        assert len(characters) == ascii_directive.size()
