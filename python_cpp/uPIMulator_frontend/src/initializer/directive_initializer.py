from abi.directive.ascii_directive import AsciiDirective
from abi.directive.asciz_directive import AscizDirective
from abi.directive.byte_directive import ByteDirective
from abi.directive.long_directive import LongDirective
from abi.directive.quad_directive import QuadDirective
from abi.directive.short_directive import ShortDirective
from initializer.int_initializer import IntInitializer
from initializer.str_initializer import StrInitializer


class DirectiveInitializer:
    def __init__(self):
        pass

    @staticmethod
    def ascii_directive():
        width = IntInitializer.value_by_range(1, 128)
        characters = StrInitializer.identifier(width)
        return AsciiDirective(characters)

    @staticmethod
    def asciz_directive():
        width = IntInitializer.value_by_range(1, 128)
        characters = StrInitializer.identifier(width)
        return AscizDirective(characters)

    @staticmethod
    def byte_directive():
        value = IntInitializer.value_by_range(0, 2 ** (8 * ByteDirective.size()))
        return ByteDirective(value)

    @staticmethod
    def short_directive():
        value = IntInitializer.value_by_range(0, 2 ** (8 * ShortDirective.size()))
        return ShortDirective(value)

    @staticmethod
    def long_directive():
        value = IntInitializer.value_by_range(0, 2 ** (8 * LongDirective.size()))
        return LongDirective(value)

    @staticmethod
    def quad_directive():
        value = IntInitializer.value_by_range(0, 2 ** (8 * QuadDirective.size()))
        return QuadDirective(value)
