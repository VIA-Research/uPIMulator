from typing import List, Optional, Union

from abi.directive.ascii_directive import AsciiDirective
from abi.directive.asciz_directive import AscizDirective
from abi.directive.byte_directive import ByteDirective
from abi.directive.long_directive import LongDirective
from abi.directive.quad_directive import QuadDirective
from abi.directive.short_directive import ShortDirective
from abi.directive.zero_directive import ZeroDirective
from abi.isa.instruction.instruction import Instruction
from abi.word.data_address_word import DataAddressWord
from abi.word.representation import Representation
from encoder.byte import Byte
from encoder.directive_encoder import DirectiveEncoder
from encoder.instruction_encoder import InstructionEncoder
from util.config_loader import ConfigLoader


class Label:
    def __init__(self, name: str):
        self._name = name
        self._address: Optional[DataAddressWord] = None
        self._size: int = 0
        self._assembler_instructions: List[
            Union[
                AsciiDirective,
                AscizDirective,
                ByteDirective,
                LongDirective,
                QuadDirective,
                ShortDirective,
                ZeroDirective,
                Instruction,
            ]
        ] = []

    def name(self) -> str:
        return self._name

    def address(self) -> Optional[int]:
        if self._address is not None:
            return self._address.value(Representation.UNSIGNED)
        else:
            return None

    def begin_address(self) -> int:
        address = self.address()
        assert address is not None
        return address

    def end_address(self) -> int:
        return self.begin_address() + self.size()

    def set_address(self, address: int) -> None:
        assert self._address is None
        self._address = DataAddressWord()

        assert address >= 0
        self._address.set_value(address)

    def size(self) -> int:
        return max(self._size, self._assembler_instructions_size())

    def set_size(self, size: int) -> None:
        assert size >= self._assembler_instructions_size()
        self._size = size

    def append_assembler_instruction(
        self,
        assembler_instruction: Union[
            AsciiDirective,
            AscizDirective,
            ByteDirective,
            LongDirective,
            QuadDirective,
            ShortDirective,
            ZeroDirective,
            Instruction,
        ],
    ) -> None:
        self._assembler_instructions.append(assembler_instruction)

    def to_bytes(self) -> List[Byte]:
        bytes_: List[Byte] = []
        for assembler_instruction in self._assembler_instructions:
            if isinstance(
                assembler_instruction,
                (
                    AsciiDirective,
                    AscizDirective,
                    ByteDirective,
                    LongDirective,
                    QuadDirective,
                    ShortDirective,
                    ZeroDirective,
                ),
            ):
                bytes_ += DirectiveEncoder.encode(assembler_instruction)
            elif isinstance(assembler_instruction, Instruction):
                bytes_ += InstructionEncoder.encode(assembler_instruction)
            else:
                raise ValueError
        return bytes_

    def _assembler_instructions_size(self) -> int:
        size = 0
        for assembler_instruction in self._assembler_instructions:
            if isinstance(
                assembler_instruction,
                (
                    AsciiDirective,
                    AscizDirective,
                    ByteDirective,
                    LongDirective,
                    QuadDirective,
                    ShortDirective,
                    ZeroDirective,
                ),
            ):
                size += assembler_instruction.size()
            elif isinstance(assembler_instruction, Instruction):
                size += ConfigLoader.iram_data_width() // 8
            else:
                raise ValueError
        return size
