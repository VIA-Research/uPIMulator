from abi.directive.asciz_directive import AscizDirective
from initializer.int_initializer import IntInitializer
from initializer.str_initializer import StrInitializer


def test_ascii_directive():
    for _ in range(100):
        characters_width = IntInitializer.value_by_range(1, 64)
        characters = StrInitializer.identifier(characters_width)

        asciz_directive = AscizDirective(characters)

        assert characters == asciz_directive.characters()
        assert len(characters) + 1 == asciz_directive.size()
