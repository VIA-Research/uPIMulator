from typing import Set

from abi.isa.cc._base_cc import BaseCC
from abi.isa.instruction.condition import Condition


class ConstCCGE0(BaseCC):
    def __init__(self, condition: Condition):
        super().__init__(ConstCCGE0.conditions(), condition)

    @staticmethod
    def conditions() -> Set[Condition]:
        return {Condition.PL}
