from typing import Set

from abi.isa.cc._base_cc import BaseCC
from abi.isa.instruction.condition import Condition


class CountNZCC(BaseCC):
    def __init__(self, condition: Condition):
        super().__init__(
            CountNZCC.conditions(), condition,
        )

    @staticmethod
    def conditions() -> Set[Condition]:
        return {
            Condition.Z,
            Condition.NZ,
            Condition.XZ,
            Condition.XNZ,
            Condition.SZ,
            Condition.SNZ,
            Condition.SPL,
            Condition.SMI,
            Condition.MAX,
            Condition.NMAX,
            Condition.TRUE,
        }
