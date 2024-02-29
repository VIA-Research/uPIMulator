from enum import Enum, auto


class Suffix(Enum):
    RICI = 0

    RRI = auto()
    RRIC = auto()
    RRICI = auto()
    RRIF = auto()
    RRR = auto()
    RRRC = auto()
    RRRCI = auto()

    ZRI = auto()
    ZRIC = auto()
    ZRICI = auto()
    ZRIF = auto()
    ZRR = auto()
    ZRRC = auto()
    ZRRCI = auto()

    S_RRI = auto()
    S_RRIC = auto()
    S_RRICI = auto()
    S_RRIF = auto()
    S_RRR = auto()
    S_RRRC = auto()
    S_RRRCI = auto()

    U_RRI = auto()
    U_RRIC = auto()
    U_RRICI = auto()
    U_RRIF = auto()
    U_RRR = auto()
    U_RRRC = auto()
    U_RRRCI = auto()

    RR = auto()
    RRC = auto()
    RRCI = auto()

    ZR = auto()
    ZRC = auto()
    ZRCI = auto()

    S_RR = auto()
    S_RRC = auto()
    S_RRCI = auto()

    U_RR = auto()
    U_RRC = auto()
    U_RRCI = auto()

    DRDICI = auto()

    RRRI = auto()
    RRRICI = auto()

    ZRRI = auto()
    ZRRICI = auto()

    S_RRRI = auto()
    S_RRRICI = auto()

    U_RRRI = auto()
    U_RRRICI = auto()

    RIR = auto()
    RIRC = auto()
    RIRCI = auto()

    ZIR = auto()
    ZIRC = auto()
    ZIRCI = auto()

    S_RIRC = auto()
    S_RIRCI = auto()

    U_RIRC = auto()
    U_RIRCI = auto()

    R = auto()
    RCI = auto()

    Z = auto()
    ZCI = auto()

    S_R = auto()
    S_RCI = auto()

    U_R = auto()
    U_RCI = auto()

    CI = auto()
    # trunk-ignore(flake8/E741)
    I = auto()

    DDCI = auto()

    ERRI = auto()

    S_ERRI = auto()
    U_ERRI = auto()

    EDRI = auto()

    ERII = auto()
    ERIR = auto()
    ERID = auto()

    DMA_RRI = auto()
