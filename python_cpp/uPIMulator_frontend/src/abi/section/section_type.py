from enum import Enum, auto


class SectionType(Enum):
    # section contains either initialized data and instructions or instructions only
    PROG_BITS = 1

    # section contains only zero-initialized data
    NO_BITS = auto()
