from typing import Set

from abi.isa.cc._base_cc import BaseCC
from abi.isa.instruction.condition import Condition


class ShiftNZCC(BaseCC):
    def __init__(self, condition: Condition):
        super().__init__(
            ShiftNZCC.conditions(), condition,
        )

    @staticmethod
    def conditions() -> Set[Condition]:
        return {
            Condition.Z,
            Condition.NZ,
            Condition.XZ,
            Condition.XNZ,
            Condition.E,
            Condition.O,
            Condition.PL,
            Condition.MI,
            Condition.SZ,
            Condition.SNZ,
            Condition.SE,
            Condition.SO,
            Condition.SPL,
            Condition.SMI,
            Condition.SH32,
            Condition.NSH32,
            Condition.TRUE,
        }
