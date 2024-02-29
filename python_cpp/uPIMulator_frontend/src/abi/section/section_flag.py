from enum import Enum, auto


class SectionFlag(Enum):
    # the section is allocatable
    ALLOC = 1

    # the section is writable
    WRITE = auto()

    # the section is executable
    EXECINSTR = auto()

    # the section has a link-order restriction
    LINK_ORDER = auto()

    # the section can be merged
    MERGE = auto()

    # the section contains null-terminated string
    STRINGS = auto()
