from enum import Enum, auto


class SectionName(Enum):
    ATOMIC = 1
    BSS = auto()
    DATA = auto()
    DEBUG_ABBREV = auto()
    DEBUG_FRAME = auto()
    DEBUG_INFO = auto()
    DEBUG_LINE = auto()
    DEBUG_LOC = auto()
    DEBUG_RANGES = auto()
    DEBUG_STR = auto()
    DPU_HOST = auto()
    MRAM = auto()
    RODATA = auto()
    STACK_SIZES = auto()
    TEXT = auto()
