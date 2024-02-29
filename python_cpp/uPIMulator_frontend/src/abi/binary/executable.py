import os
from typing import List, Optional, Set, Union

from abi.binary.liveness import Liveness
from abi.binary.relocatable import Relocatable
from abi.directive.ascii_directive import AsciiDirective
from abi.directive.asciz_directive import AscizDirective
from abi.directive.byte_directive import ByteDirective
from abi.directive.long_directive import LongDirective
from abi.directive.quad_directive import QuadDirective
from abi.directive.short_directive import ShortDirective
from abi.directive.zero_directive import ZeroDirective
from abi.isa.instruction.instruction import Instruction
from abi.label.label import Label
from abi.section.section import Section
from abi.section.section_flag import SectionFlag
from abi.section.section_name import SectionName
from abi.section.section_type import SectionType
from parser_.parser import Parser


class Executable:
    def __init__(self, filepath: str, relocatable: Set[Relocatable]):
        self._filepath: str = filepath
        self._relocatables: Set[Relocatable] = relocatable
        self._lines: str = self._init_lines()
        self._liveness: Liveness = Liveness()
        self._sections: Set[Section] = set()
        self._cur_section: Optional[Section] = None

    def filepath(self) -> str:
        return self._filepath

    def lines(self) -> str:
        return self._lines

    def liveness(self) -> Liveness:
        return self._liveness

    def sections(self, section_name: SectionName) -> Set[Section]:
        sections: Set[Section] = set()
        for section in self._sections:
            if section.section_name() == section_name:
                sections.add(section)
        return sections

    def section(self, section_name: SectionName, name: str) -> Optional[Section]:
        for section in self._sections:
            if section_name == section.section_name() and name == section.name():
                return section
        return None

    def cur_section(self) -> Section:
        assert self._cur_section is not None
        return self._cur_section

    def checkout_section(
        self, section_name: SectionName, name: str, section_flags: Set[SectionFlag], section_type: SectionType,
    ) -> None:
        if self.section(section_name, name) is None:
            self._sections.add(Section(section_name, name, section_flags, section_type))

        section = self.section(section_name, name)
        assert section is not None
        self._cur_section = section
        section.checkout_hidden_label()

    def label(self, label_name: str) -> Optional[Label]:
        for section in self._sections:
            label = section.label(label_name)
            if label is not None:
                return label
        return None

    def labels(self) -> List[Label]:
        labels: List[Label] = []
        for section in self._sections:
            labels += section.labels()
        return labels

    def cur_label(self) -> Label:
        assert self._cur_section is not None
        return self._cur_section.cur_label()

    def append_label(self, label_name: str) -> None:
        assert self._cur_section is not None
        self._cur_section.append_label(label_name)

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
        assert self._cur_section is not None
        self._cur_section.append_assembler_instruction(assembler_instruction)

    def _init_lines(self) -> str:
        lines = ""
        for relocatable in self._relocatables:
            prefix = f"{relocatable.filepath().split(os.path.sep)[-1].split('.')[0]}"

            with open(relocatable.filepath(), encoding="ISO-8859-1") as file:
                relocatable_lines = Parser.preprocess("".join(file.readlines()))

                if prefix != "main":
                    for relocatable_local_symbols in relocatable.liveness().local_symbols():
                        relocatable_lines = relocatable_lines.replace(
                            f"{relocatable_local_symbols},", f"{prefix}.{relocatable_local_symbols},",
                        )
                        relocatable_lines = relocatable_lines.replace(
                            f"{relocatable_local_symbols} ", f"{prefix}.{relocatable_local_symbols} ",
                        )
                        relocatable_lines = relocatable_lines.replace(
                            f"{relocatable_local_symbols}\t", f"{prefix}.{relocatable_local_symbols}\t",
                        )
                        relocatable_lines = relocatable_lines.replace(
                            f"{relocatable_local_symbols}\n", f"{prefix}.{relocatable_local_symbols}\n",
                        )
                        relocatable_lines = relocatable_lines.replace(
                            f"{relocatable_local_symbols}:", f"{prefix}.{relocatable_local_symbols}:",
                        )

                lines += relocatable_lines
        return lines
