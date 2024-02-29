from typing import Set

from abi.isa.cc._base_cc import BaseCC
from abi.isa.instruction.condition import Condition


class SubNZCC(BaseCC):
    def __init__(self, condition: Condition):
        super().__init__(
            SubNZCC.conditions(), condition,
        )

    @staticmethod
    def conditions() -> Set[Condition]:
        return {
            Condition.C,
            Condition.NC,
            Condition.Z,
            Condition.NZ,
            Condition.XZ,
            Condition.XNZ,
            Condition.OV,
            Condition.NOV,
            Condition.MI,
            Condition.PL,
            Condition.EQ,
            Condition.NEQ,
            Condition.SPL,
            Condition.SMI,
            Condition.GES,
            Condition.GEU,
            Condition.GTS,
            Condition.GTU,
            Condition.LES,
            Condition.LEU,
            Condition.LTS,
            Condition.LTU,
            Condition.XGTS,
            Condition.XGTU,
            Condition.XLES,
            Condition.XLEU,
            Condition.TRUE,
        }
