from enum import Enum, auto


class Condition(Enum):
    # condition is always true, whatever the operation result is
    TRUE = 0

    # condition is always false, whatever the operation result is
    FALSE = auto()

    # operation result equals zero
    Z = auto()

    # operation result does not equal zero
    NZ = auto()

    # operation result is even
    E = auto()

    # operation result is odd
    # trunk-ignore(flake8/E741)
    O = auto()

    # operation result is greater or equal to zero
    PL = auto()

    # operation result is strictly lower than zero
    MI = auto()

    # operation result has overflow set
    OV = auto()

    # operation result does not have overflow set
    NOV = auto()

    # operation result has carry set
    C = auto()

    # operation result does not have carry set
    NC = auto()

    # source register operand is equal to zero
    SZ = auto()

    # source register operand is not equal to zero
    SNZ = auto()

    # source register operand is positive or null
    SPL = auto()

    # source register operand is strictly negative
    SMI = auto()

    # source register operand is odd
    SO = auto()

    # source register operand is even
    SE = auto()

    # NC5 NC6 NC7 NC8 NC9 NC10 NC11 NC12 NC13 NC14 operation result set the
    # carry flag number 5, 6, 7, 8, 9, 10, 11, 12, 13 or 14 respectively.
    # these conditions may come in handy to quickly detect buffer overflows.
    NC5 = auto()
    NC6 = auto()
    NC7 = auto()
    NC8 = auto()
    NC9 = auto()
    NC10 = auto()
    NC11 = auto()
    NC12 = auto()
    NC13 = auto()
    NC14 = auto()

    # operation is a bit count and the result is the maximum count value
    MAX = auto()

    # operation is a bit count and the result is not the maximum count value
    NMAX = auto()

    # second operand is a register with a value having bit number 5 equal to one
    SH32 = auto()

    # second operand is a register with a value having bit number 5 equal to zero
    NSH32 = auto()

    # first operand is equal to the second operand
    EQ = auto()

    # first operand is different from the second operand
    NEQ = auto()

    # first operand is respectively lower than, lower or equal to, greater than,
    # greater or equal to the second operand when performing an unsigned comparison
    LTU = auto()
    LEU = auto()
    GTU = auto()
    GEU = auto()

    # first operand of is respectively lower than, lower or equal to, greater
    # than, greater or equal to the second operand when performing a signed
    # comparison
    LTS = auto()
    LES = auto()
    GTS = auto()
    GES = auto()

    # operation result is null and ZeroFlag is set
    XZ = auto()

    # operation result is not null or ZeroFlag is not set
    XNZ = auto()

    # either operation result holds carry flag or ZeroFlag is set
    XLEU = auto()

    # operation result holds carry flag and ZeroFlag is not set
    XGTU = auto()

    # ZeroFLag is set and either operation result is negative or overflows
    XLES = auto()

    # ZeroFlag is not set and either operation result is positive or null or overflows
    XGTS = auto()

    # operation is an 8x8 multiplication and the result is a less than 256
    SMALL = auto()

    # operation is an 8x8 multiplication and the result is greater than 255
    LARGE = auto()
