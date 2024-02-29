from typing import Set

from abi.isa.cc._base_cc import BaseCC
from abi.isa.instruction.condition import Condition


class TrueFalseCC(BaseCC):
    def __init__(self, condition: Condition):
        super().__init__(TrueFalseCC.conditions(), condition)

    @staticmethod
    def conditions() -> Set[Condition]:
        return {Condition.TRUE, Condition.FALSE}
