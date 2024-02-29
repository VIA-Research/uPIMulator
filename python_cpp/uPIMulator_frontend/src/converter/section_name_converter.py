from abi.section.section_name import SectionName


class SectionNameConverter:
    def __init__(self):
        pass

    @staticmethod
    def convert_to_section_name(section_name: str) -> SectionName:
        if section_name == "atomic":
            return SectionName.ATOMIC
        elif section_name == "bss":
            return SectionName.BSS
        elif section_name == "data":
            return SectionName.DATA
        elif section_name == "debug_abbrev":
            return SectionName.DEBUG_ABBREV
        elif section_name == "debug_frame":
            return SectionName.DEBUG_FRAME
        elif section_name == "debug_info":
            return SectionName.DEBUG_INFO
        elif section_name == "debug_line":
            return SectionName.DEBUG_LINE
        elif section_name == "debug_loc":
            return SectionName.DEBUG_LOC
        elif section_name == "debug_ranges":
            return SectionName.DEBUG_RANGES
        elif section_name == "debug_str":
            return SectionName.DEBUG_STR
        elif section_name == "dpu_host":
            return SectionName.DPU_HOST
        elif section_name == "mram":
            return SectionName.MRAM
        elif section_name == "rodata":
            return SectionName.RODATA
        elif section_name == "stack_sizes":
            return SectionName.STACK_SIZES
        elif section_name == "text":
            return SectionName.TEXT
        else:
            raise ValueError

    @staticmethod
    def convert_to_str(section_name: SectionName) -> str:
        if section_name == SectionName.ATOMIC:
            return "atomic"
        elif section_name == SectionName.BSS:
            return "bss"
        elif section_name == SectionName.DATA:
            return "data"
        elif section_name == SectionName.DEBUG_ABBREV:
            return "debug_abbrev"
        elif section_name == SectionName.DEBUG_FRAME:
            return "debug_frame"
        elif section_name == SectionName.DEBUG_INFO:
            return "debug_info"
        elif section_name == SectionName.DEBUG_LINE:
            return "debug_line"
        elif section_name == SectionName.DEBUG_LOC:
            return "debug_loc"
        elif section_name == SectionName.DEBUG_RANGES:
            return "debug_ranges"
        elif section_name == SectionName.DEBUG_STR:
            return "debug_ranges"
        elif section_name == SectionName.DPU_HOST:
            return "dpu_host"
        elif section_name == SectionName.MRAM:
            return "mram"
        elif section_name == SectionName.RODATA:
            return "rodata"
        elif section_name == SectionName.STACK_SIZES:
            return "stack_sizes"
        elif section_name == SectionName.TEXT:
            return "text"
        else:
            raise ValueError
