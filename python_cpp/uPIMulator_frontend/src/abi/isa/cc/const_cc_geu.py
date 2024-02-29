from typing import Set

from abi.isa.cc._base_cc import BaseCC
from abi.isa.instruction.condition import Condition


class ConstCCGEU(BaseCC):
    def __init__(self, condition: Condition):
        super().__init__(ConstCCGEU.conditions(), condition)

    @staticmethod
    def conditions() -> Set[Condition]:
        return {Condition.GEU}
