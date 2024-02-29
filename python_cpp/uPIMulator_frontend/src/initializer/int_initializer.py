import random
from typing import Any, List

from abi.word.representation import Representation


class IntInitializer:
    def __init__(self):
        pass

    @staticmethod
    def value_by_range(min_value: int, max_value: int) -> int:
        assert min_value < max_value

        return random.randrange(min_value, max_value)

    @staticmethod
    def value_by_width(representation: Representation, width: int):
        assert width > 0

        if representation == Representation.UNSIGNED:
            return IntInitializer.value_by_range(0, 2 ** width)
        elif representation == Representation.SIGNED:
            return IntInitializer.value_by_range(-(2 ** (width - 1)), 2 ** (width - 1))
        else:
            raise ValueError

    @staticmethod
    def value_by_list(list_: List[Any]) -> Any:
        return random.choice(list_)
