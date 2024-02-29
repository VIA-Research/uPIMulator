from typing import List, Union

from abi.directive.ascii_directive import AsciiDirective
from abi.directive.asciz_directive import AscizDirective
from abi.directive.byte_directive import ByteDirective
from abi.directive.long_directive import LongDirective
from abi.directive.quad_directive import QuadDirective
from abi.directive.short_directive import ShortDirective
from abi.directive.zero_directive import ZeroDirective
from abi.word.immediate import Immediate
from abi.word.representation import Representation
from encoder.ascii_encoder import AsciiEncoder
from encoder.byte import Byte


class DirectiveEncoder:
    def __init__(self):
        pass

    @staticmethod
    def encode(
        directive: Union[
            AsciiDirective, AscizDirective, ByteDirective, LongDirective, QuadDirective, ShortDirective, ZeroDirective,
        ]
    ) -> List[Byte]:
        if isinstance(directive, AsciiDirective):
            return DirectiveEncoder._encode_ascii(directive)
        elif isinstance(directive, AscizDirective):
            return DirectiveEncoder._encode_asciz(directive)
        elif isinstance(directive, (ByteDirective, LongDirective, QuadDirective, ShortDirective)):
            return DirectiveEncoder._encode_byte(directive)
        elif isinstance(directive, ZeroDirective):
            return DirectiveEncoder._encode_zero(directive)
        else:
            raise ValueError

    @staticmethod
    def decode(
        bytes_: List[Byte],
    ) -> Union[
        AsciiDirective, AscizDirective, ByteDirective, LongDirective, QuadDirective, ShortDirective, ZeroDirective,
    ]:
        raise AttributeError

    @staticmethod
    def _encode_ascii(directive: AsciiDirective):
        return [AsciiEncoder.encode(character) for character in directive.characters()]

    @staticmethod
    def _encode_asciz(directive: AscizDirective):
        bytes_ = [AsciiEncoder.encode(character) for character in directive.characters()]
        bytes_.append(AsciiEncoder.encode(AsciiEncoder.unknown()))
        return bytes_

    @staticmethod
    def _encode_byte(
        directive: Union[ByteDirective, LongDirective, QuadDirective, ShortDirective, ZeroDirective,]
    ):
        return Immediate(Representation.UNSIGNED, 8 * directive.size(), directive.value()).to_bytes()

    @staticmethod
    def _encode_zero(directive: ZeroDirective):
        assert isinstance(directive, ZeroDirective)

        byte = Immediate(Representation.UNSIGNED, 8, directive.value()).to_bytes()[0]
        return [byte for _ in range(directive.size())]
