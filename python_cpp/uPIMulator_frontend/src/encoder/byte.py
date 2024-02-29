class Byte:
    def __init__(self, value: int):
        assert 0 <= value < 2 ** 8

        self._value: int = value

    def value(self):
        return self._value
