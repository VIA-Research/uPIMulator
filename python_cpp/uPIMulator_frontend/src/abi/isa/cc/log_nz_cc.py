from typing import Set

from abi.isa.cc._base_cc import BaseCC
from abi.isa.instruction.condition import Condition


class LogNZCC(BaseCC):
    def __init__(self, condition: Condition):
        super().__init__(
            LogNZCC.conditions(), condition,
        )

    @staticmethod
    def conditions() -> Set[Condition]:
        return {
            Condition.Z,
            Condition.NZ,
            Condition.XZ,
            Condition.XNZ,
            Condition.PL,
            Condition.MI,
            Condition.SZ,
            Condition.SNZ,
            Condition.SPL,
            Condition.SMI,
            Condition.TRUE,
        }
