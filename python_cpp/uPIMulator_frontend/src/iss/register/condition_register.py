from abi.isa.instruction.condition import Condition


class ConditionRegister:
    def __init__(self):
        self._bits = [False for _ in range(len(Condition))]

    def condition(self, condition: Condition) -> bool:
        if condition == Condition.TRUE:
            return True
        elif condition == Condition.FALSE:
            return False
        else:
            return self._bits[condition.value]

    def set_condition(self, condition: Condition) -> None:
        assert condition != Condition.TRUE and condition != Condition.FALSE

        self._bits[condition.value] = True

    def clear_condition(self, condition: Condition) -> None:
        assert condition != Condition.TRUE and condition != Condition.FALSE

        self._bits[condition.value] = False
