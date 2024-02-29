from abi.isa.instruction.endian import Endian


class EndianConverter:
    def __init__(self):
        pass

    @staticmethod
    def convert_to_endian(endian: str) -> Endian:
        if endian == "!little":
            return Endian.LITTLE
        elif endian == "!big":
            return Endian.BIG
        else:
            raise ValueError

    @staticmethod
    def convert_to_string(endian: Endian) -> str:
        if endian == Endian.LITTLE:
            return "!little"
        elif endian == Endian.BIG:
            return "!big"
        else:
            raise ValueError
