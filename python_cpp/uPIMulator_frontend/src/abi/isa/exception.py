from enum import Enum, auto


class Exception_(Enum):
    MEMORY_FAULT = 0
    DMA_FAULT = auto()
    HEAP_FULL = auto()
    DIVISION_BY_ZERO = auto()
    ASSERT = auto()
    HALT = auto()
    PRINT_OVERFLOW = auto()
    ALREADY_PROFILING = auto()
    NOT_PROFILING = auto()
