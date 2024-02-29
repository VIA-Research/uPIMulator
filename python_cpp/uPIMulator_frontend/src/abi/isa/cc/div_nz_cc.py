from typing import Set

from abi.isa.cc._base_cc import BaseCC
from abi.isa.instruction.condition import Condition


class DivNZCC(BaseCC):
    def __init__(self, condition: Condition):
        super().__init__(
            DivNZCC.conditions(), condition,
        )

    @staticmethod
    def conditions() -> Set[Condition]:
        return {
            Condition.SZ,
            Condition.SNZ,
            Condition.SPL,
            Condition.SMI,
            Condition.TRUE,
        }
