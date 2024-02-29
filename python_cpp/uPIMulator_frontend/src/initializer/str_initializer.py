import random


class StrInitializer:
    def __init__(self):
        pass

    @staticmethod
    def identifier(width: int) -> str:
        assert width > 0

        return random.choice(StrInitializer._non_digit_characters()) + "".join(
            [random.choice(StrInitializer._characters()) for _ in range(width - 1)]
        )

    @staticmethod
    def _non_digit_characters() -> str:
        return StrInitializer._alphabets() + StrInitializer._digits()

    @staticmethod
    def _characters() -> str:
        return StrInitializer._non_digit_characters() + StrInitializer._digits()

    @staticmethod
    def _alphabets() -> str:
        return "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

    @staticmethod
    def _digits() -> str:
        return "0123456789"

    @staticmethod
    def _special_symbols() -> str:
        return "._"
