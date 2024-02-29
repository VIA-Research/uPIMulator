from typing import Set

from abi.isa.instruction.condition import Condition


class BaseCC:
    def __init__(self, conditions: Set[Condition], condition: Condition):
        assert condition in conditions
        self._condition: Condition = condition

    def condition(self):
        return self._condition
