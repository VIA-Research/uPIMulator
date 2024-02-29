from abi.binary.liveness import Liveness
from parser_.parser import Parser


class Relocatable:
    def __init__(self, filepath: str):
        self._filepath: str = filepath
        self._lines: str = self._init_lines()
        self._liveness: Liveness = Liveness()

    def filepath(self) -> str:
        return self._filepath

    def lines(self) -> str:
        return self._lines

    def liveness(self) -> Liveness:
        return self._liveness

    def _init_lines(self) -> str:
        with open(self._filepath, encoding="ISO-8859-1") as file:
            return Parser.preprocess("".join(file.readlines()))
