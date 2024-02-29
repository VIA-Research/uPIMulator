from typing import Set

from abi.section.section_flag import SectionFlag


class SectionFlagConverter:
    def __init__(self):
        pass

    @staticmethod
    def convert_to_section_flags(section_flags: str) -> Set[SectionFlag]:
        def convert(section_flag: str) -> SectionFlag:
            if section_flag == "a":
                return SectionFlag.ALLOC
            elif section_flag == "w":
                return SectionFlag.WRITE
            elif section_flag == "x":
                return SectionFlag.EXECINSTR
            elif section_flag == "o":
                return SectionFlag.LINK_ORDER
            elif section_flag == "M":
                return SectionFlag.MERGE
            elif section_flag == "S":
                return SectionFlag.STRINGS
            else:
                raise ValueError

        return {convert(section_flag) for section_flag in section_flags}
