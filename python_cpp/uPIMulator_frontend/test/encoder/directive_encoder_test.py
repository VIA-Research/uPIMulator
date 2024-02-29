from abi.directive.byte_directive import ByteDirective
from abi.directive.long_directive import LongDirective
from abi.directive.quad_directive import QuadDirective
from abi.directive.short_directive import ShortDirective
from encoder.ascii_encoder import AsciiEncoder
from encoder.directive_encoder import DirectiveEncoder
from initializer.directive_initializer import DirectiveInitializer


def test_ascii_directive():
    ascii_directive = DirectiveInitializer.ascii_directive()

    bytes_ = DirectiveEncoder.encode(ascii_directive)
    characters = "".join([AsciiEncoder.decode(byte) for byte in bytes_])

    assert ascii_directive.characters() == characters


def test_asciz_directive():
    asciz_directive = DirectiveInitializer.asciz_directive()

    bytes_ = DirectiveEncoder.encode(asciz_directive)
    characters = "".join([AsciiEncoder.decode(byte) for byte in bytes_])

    assert asciz_directive.characters() == characters[:-1]


def test_byte_directive():
    byte_directive = DirectiveInitializer.byte_directive()

    bytes_ = DirectiveEncoder.encode(byte_directive)

    for i in range(ByteDirective.size()):
        byte_directive.value() == bytes_[i].value()


def test_short_directive():
    short_directive = DirectiveInitializer.short_directive()

    bytes_ = DirectiveEncoder.encode(short_directive)

    for i in range(ShortDirective.size()):
        short_directive.value() == bytes_[i].value()


def test_long_directive():
    long_directive = DirectiveInitializer.long_directive()

    bytes_ = DirectiveEncoder.encode(long_directive)

    for i in range(LongDirective.size()):
        long_directive.value() == bytes_[i].value()


def test_quad_directive():
    quad_directive = DirectiveInitializer.quad_directive()

    bytes_ = DirectiveEncoder.encode(quad_directive)

    for i in range(QuadDirective.size()):
        quad_directive.value() == bytes_[i].value()
