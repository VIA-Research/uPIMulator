from abi.isa.instruction.condition import Condition


class ConditionConverter:
    def __init__(self):
        pass

    @staticmethod
    def convert_to_condition(condition: str) -> Condition:
        if condition == "true":
            return Condition.TRUE
        elif condition == "false":
            return Condition.FALSE
        elif condition == "z":
            return Condition.Z
        elif condition == "nz":
            return Condition.NZ
        elif condition == "e":
            return Condition.E
        elif condition == "o":
            return Condition.O
        elif condition == "pl":
            return Condition.PL
        elif condition == "mi":
            return Condition.MI
        elif condition == "ov":
            return Condition.OV
        elif condition == "nov":
            return Condition.NOV
        elif condition == "c":
            return Condition.C
        elif condition == "nc":
            return Condition.NC
        elif condition == "sz":
            return Condition.SZ
        elif condition == "snz":
            return Condition.SNZ
        elif condition == "spl":
            return Condition.SPL
        elif condition == "smi":
            return Condition.SMI
        elif condition == "so":
            return Condition.SO
        elif condition == "se":
            return Condition.SE
        elif condition == "nc5":
            return Condition.NC5
        elif condition == "nc6":
            return Condition.NC6
        elif condition == "nc7":
            return Condition.NC7
        elif condition == "nc8":
            return Condition.NC8
        elif condition == "nc9":
            return Condition.NC9
        elif condition == "nc10":
            return Condition.NC10
        elif condition == "nc11":
            return Condition.NC11
        elif condition == "nc12":
            return Condition.NC12
        elif condition == "nc13":
            return Condition.NC13
        elif condition == "nc14":
            return Condition.NC14
        elif condition == "max":
            return Condition.MAX
        elif condition == "nmax":
            return Condition.NMAX
        elif condition == "sh32":
            return Condition.SH32
        elif condition == "nsh32":
            return Condition.NSH32
        elif condition == "eq":
            return Condition.EQ
        elif condition == "neq":
            return Condition.NEQ
        elif condition == "ltu":
            return Condition.LTU
        elif condition == "leu":
            return Condition.LEU
        elif condition == "gtu":
            return Condition.GTU
        elif condition == "geu":
            return Condition.GEU
        elif condition == "lts":
            return Condition.LTS
        elif condition == "les":
            return Condition.LES
        elif condition == "gts":
            return Condition.GTS
        elif condition == "ges":
            return Condition.GES
        elif condition == "xz":
            return Condition.XZ
        elif condition == "xnz":
            return Condition.XNZ
        elif condition == "xleu":
            return Condition.XLEU
        elif condition == "xgtu":
            return Condition.XGTU
        elif condition == "xles":
            return Condition.XLES
        elif condition == "xgts":
            return Condition.XGTS
        elif condition == "small":
            return Condition.SMALL
        elif condition == "large":
            return Condition.LARGE
        else:
            raise ValueError

    @staticmethod
    def convert_to_string(condition: Condition) -> str:
        if condition == Condition.TRUE:
            return "true"
        elif condition == Condition.FALSE:
            return "false"
        elif condition == Condition.Z:
            return "z"
        elif condition == Condition.NZ:
            return "nz"
        elif condition == Condition.E:
            return "e"
        elif condition == Condition.O:
            return "o"
        elif condition == Condition.PL:
            return "pl"
        elif condition == Condition.MI:
            return "mi"
        elif condition == Condition.OV:
            return "ov"
        elif condition == Condition.NOV:
            return "nov"
        elif condition == Condition.C:
            return "c"
        elif condition == Condition.NC:
            return "nc"
        elif condition == Condition.SZ:
            return "sz"
        elif condition == Condition.SNZ:
            return "snz"
        elif condition == Condition.SPL:
            return "spl"
        elif condition == Condition.SMI:
            return "smi"
        elif condition == Condition.SO:
            return "so"
        elif condition == Condition.SE:
            return "se"
        elif condition == Condition.NC5:
            return "nc5"
        elif condition == Condition.NC6:
            return "nc6"
        elif condition == Condition.NC7:
            return "nc7"
        elif condition == Condition.NC8:
            return "nc8"
        elif condition == Condition.NC9:
            return "nc9"
        elif condition == Condition.NC10:
            return "nc10"
        elif condition == Condition.NC11:
            return "nc11"
        elif condition == Condition.NC12:
            return "nc12"
        elif condition == Condition.NC13:
            return "nc13"
        elif condition == Condition.NC14:
            return "nc14"
        elif condition == Condition.MAX:
            return "max"
        elif condition == Condition.NMAX:
            return "nmax"
        elif condition == Condition.SH32:
            return "sh32"
        elif condition == Condition.NSH32:
            return "nsh32"
        elif condition == Condition.EQ:
            return "eq"
        elif condition == Condition.NEQ:
            return "neq"
        elif condition == Condition.LTU:
            return "ltu"
        elif condition == Condition.LEU:
            return "leu"
        elif condition == Condition.GTU:
            return "gtu"
        elif condition == Condition.GEU:
            return "geu"
        elif condition == Condition.LTS:
            return "lts"
        elif condition == Condition.LES:
            return "les"
        elif condition == Condition.GTS:
            return "gts"
        elif condition == Condition.GES:
            return "ges"
        elif condition == Condition.XZ:
            return "xz"
        elif condition == Condition.XNZ:
            return "xnz"
        elif condition == Condition.XLEU:
            return "xleu"
        elif condition == Condition.XGTU:
            return "xgtu"
        elif condition == Condition.XLES:
            return "xles"
        elif condition == Condition.XGTS:
            return "xgts"
        elif condition == Condition.SMALL:
            return "small"
        elif condition == Condition.LARGE:
            return "large"
        else:
            raise ValueError
