class AscizDirective:
    def __init__(self, characters: str):
        self._characters: str = characters

    def characters(self) -> str:
        return self._characters

    def size(self) -> int:
        return len(self._characters) + 1
