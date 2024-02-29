from typing import List

from encoder.byte import Byte


class Bin:
    def __init__(self, bytes_: List[Byte]):
        self._bytes: List[Byte] = bytes_

    def dump(self, filepath: str):
        with open(filepath, "w") as file:
            lines = ""
            for byte in self._bytes:
                lines += f"{byte.value()}\n"
            file.writelines(lines)
