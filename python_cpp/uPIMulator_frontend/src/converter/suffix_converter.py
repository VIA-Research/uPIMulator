from abi.isa.instruction.suffix import Suffix


class SuffixConverter:
    def __init__(self):
        pass

    @staticmethod
    def convert_to_string(suffix: Suffix):
        if suffix == Suffix.RICI:
            return "rici"
        elif suffix == Suffix.RRI:
            return "rri"
        elif suffix == Suffix.RRIC:
            return "rric"
        elif suffix == Suffix.RRICI:
            return "rrici"
        elif suffix == Suffix.RRIF:
            return "rrif"
        elif suffix == Suffix.RRR:
            return "rrr"
        elif suffix == Suffix.RRRC:
            return "rrrc"
        elif suffix == Suffix.RRRCI:
            return "rrrci"
        elif suffix == Suffix.ZRI:
            return "zri"
        elif suffix == Suffix.ZRIC:
            return "zric"
        elif suffix == Suffix.ZRICI:
            return "zrici"
        elif suffix == Suffix.ZRIF:
            return "zrif"
        elif suffix == Suffix.ZRR:
            return "zrr"
        elif suffix == Suffix.ZRRC:
            return "zrrc"
        elif suffix == Suffix.ZRRCI:
            return "zrrci"
        elif suffix == Suffix.S_RRI:
            return "s_rri"
        elif suffix == Suffix.S_RRIC:
            return "s_rric"
        elif suffix == Suffix.S_RRICI:
            return "s_rrici"
        elif suffix == Suffix.S_RRIF:
            return "s_rrif"
        elif suffix == Suffix.S_RRR:
            return "s_rrr"
        elif suffix == Suffix.S_RRRC:
            return "s_rrrc"
        elif suffix == Suffix.S_RRRCI:
            return "s_rrrci"
        elif suffix == Suffix.U_RRI:
            return "u_rri"
        elif suffix == Suffix.U_RRIC:
            return "u_rric"
        elif suffix == Suffix.U_RRICI:
            return "u_rrici"
        elif suffix == Suffix.U_RRIF:
            return "u_rrif"
        elif suffix == Suffix.U_RRR:
            return "u_rrr"
        elif suffix == Suffix.U_RRRC:
            return "u_rrrc"
        elif suffix == Suffix.U_RRRCI:
            return "u_rrrci"
        elif suffix == Suffix.RR:
            return "rr"
        elif suffix == Suffix.RRC:
            return "rrc"
        elif suffix == Suffix.RRCI:
            return "rrci"
        elif suffix == Suffix.ZR:
            return "zr"
        elif suffix == Suffix.ZRC:
            return "zrc"
        elif suffix == Suffix.ZRCI:
            return "zrci"
        elif suffix == Suffix.S_RR:
            return "s_rr"
        elif suffix == Suffix.S_RRC:
            return "s_rrc"
        elif suffix == Suffix.S_RRCI:
            return "s_rrci"
        elif suffix == Suffix.U_RR:
            return "u_rr"
        elif suffix == Suffix.U_RRC:
            return "u_rrc"
        elif suffix == Suffix.U_RRCI:
            return "u_rrci"
        elif suffix == Suffix.DRDICI:
            return "drdici"
        elif suffix == Suffix.RRRI:
            return "rrri"
        elif suffix == Suffix.RRRICI:
            return "rrrici"
        elif suffix == Suffix.ZRRI:
            return "zrri"
        elif suffix == Suffix.ZRRICI:
            return "zrrici"
        elif suffix == Suffix.S_RRRI:
            return "s_rrri"
        elif suffix == Suffix.S_RRRICI:
            return "s_rrrici"
        elif suffix == Suffix.U_RRRI:
            return "u_rrri"
        elif suffix == Suffix.U_RRRICI:
            return "u_rrrici"
        elif suffix == Suffix.RIR:
            return "rir"
        elif suffix == Suffix.RIRC:
            return "rirc"
        elif suffix == Suffix.RIRCI:
            return "rirci"
        elif suffix == Suffix.ZIR:
            return "zir"
        elif suffix == Suffix.ZIRC:
            return "zirc"
        elif suffix == Suffix.ZIRCI:
            return "zirci"
        elif suffix == Suffix.S_RIRC:
            return "s_rirc"
        elif suffix == Suffix.S_RIRCI:
            return "s_rirci"
        elif suffix == Suffix.U_RIRC:
            return "u_rirc"
        elif suffix == Suffix.U_RIRCI:
            return "u_rirci"
        elif suffix == Suffix.R:
            return "r"
        elif suffix == Suffix.RCI:
            return "rci"
        elif suffix == Suffix.Z:
            return "z"
        elif suffix == Suffix.ZCI:
            return "zci"
        elif suffix == Suffix.S_R:
            return "s_r"
        elif suffix == Suffix.S_RCI:
            return "s_rci"
        elif suffix == Suffix.U_R:
            return "u_r"
        elif suffix == Suffix.U_RCI:
            return "u_rci"
        elif suffix == Suffix.CI:
            return "ci"
        elif suffix == Suffix.I:
            return "i"
        elif suffix == Suffix.DDCI:
            return "ddci"
        elif suffix == Suffix.ERRI:
            return "erri"
        elif suffix == Suffix.S_ERRI:
            return "s_erri"
        elif suffix == Suffix.U_ERRI:
            return "u_erri"
        elif suffix == Suffix.EDRI:
            return "edri"
        elif suffix == Suffix.ERII:
            return "erii"
        elif suffix == Suffix.ERIR:
            return "erir"
        elif suffix == Suffix.ERID:
            return "erid"
        elif suffix == Suffix.DMA_RRI:
            return "dma_rri"
        else:
            raise ValueError
