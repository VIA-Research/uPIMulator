from typing import List, Optional, Set, Union

from abi.directive.ascii_directive import AsciiDirective
from abi.directive.asciz_directive import AscizDirective
from abi.directive.byte_directive import ByteDirective
from abi.directive.long_directive import LongDirective
from abi.directive.quad_directive import QuadDirective
from abi.directive.short_directive import ShortDirective
from abi.directive.zero_directive import ZeroDirective
from abi.isa.instruction.instruction import Instruction
from abi.label.label import Label
from abi.section.section_flag import SectionFlag
from abi.section.section_name import SectionName
from abi.section.section_type import SectionType
from abi.word.instruction_word import InstructionWord
from converter.section_name_converter import SectionNameConverter
from encoder.byte import Byte


class Section:
    def __init__(
        self, section_name: SectionName, name: str, section_flags: Set[SectionFlag], section_type: SectionType,
    ):
        self._section_name: SectionName = section_name
        self._name = name
        self._section_flags: Set[SectionFlag] = section_flags
        self._section_type: SectionType = section_type

        self._hidden_label = Label(self._hidden_label_name())
        self._labels: List[Label] = [self._hidden_label]
        self._cur_label: Label = self._hidden_label

    def section_name(self) -> SectionName:
        return self._section_name

    def name(self) -> str:
        return self._name

    def address(self) -> Optional[int]:
        return self._labels[0].address()

    def set_address(self, address: int) -> None:
        assert self.address() is None
        assert address >= 0

        offset = 0
        for label in self._labels:
            label.set_address(address + offset)
            offset += label.size()

            if self._section_name == SectionName.TEXT:
                label_address = label.address()
                label_size = label.size()

                assert label_address is not None
                assert label_address % InstructionWord().size() == 0
                assert label_size % InstructionWord().size() == 0

    def size(self) -> int:
        return sum([label.size() for label in self._labels])

    def labels(self) -> List[Label]:
        return self._labels

    def label(self, label_name: str) -> Optional[Label]:
        for label in self._labels:
            if label_name == label.name():
                return label
        return None

    def checkout_hidden_label(self) -> None:
        self._cur_label = self._hidden_label

    def cur_label(self) -> Label:
        return self._cur_label

    def append_label(self, label_name: str):
        if self.label(label_name) is None:
            self._labels.append(Label(label_name))

        label = self.label(label_name)
        assert label is not None
        self._cur_label = label

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
        if self._section_name == SectionName.TEXT:
            assert isinstance(assembler_instruction, Instruction)
        else:
            assert isinstance(
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
            )

        self._cur_label.append_assembler_instruction(assembler_instruction)

    def to_bytes(self) -> List[Byte]:
        bytes_: List[Byte] = []
        for label in self._labels:
            bytes_ += label.to_bytes()
        return bytes_

    def _hidden_label_name(self) -> str:
        return f"{SectionNameConverter.convert_to_str(self._section_name)}.{self._name}"
