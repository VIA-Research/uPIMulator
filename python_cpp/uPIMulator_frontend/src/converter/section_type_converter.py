from abi.section.section_type import SectionType


class SectionTypeConverter:
    def __init__(self):
        pass

    @staticmethod
    def convert_to_section_type(section_type: str) -> SectionType:
        if section_type == "@progbits":
            return SectionType.PROG_BITS
        elif section_type == "@nobits":
            return SectionType.NO_BITS
        else:
            raise ValueError
