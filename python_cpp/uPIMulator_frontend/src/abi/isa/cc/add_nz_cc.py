from typing import Set

from abi.isa.cc._base_cc import BaseCC
from abi.isa.instruction.condition import Condition


class AddNZCC(BaseCC):
    def __init__(self, condition: Condition):
        super().__init__(
            AddNZCC.conditions(), condition,
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
            Condition.PL,
            Condition.MI,
            Condition.SZ,
            Condition.SNZ,
            Condition.SPL,
            Condition.SMI,
            Condition.NC5,
            Condition.NC6,
            Condition.NC7,
            Condition.NC8,
            Condition.NC9,
            Condition.NC10,
            Condition.NC11,
            Condition.NC12,
            Condition.NC13,
            Condition.NC14,
            Condition.TRUE,
        }
