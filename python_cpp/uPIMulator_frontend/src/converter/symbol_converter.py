from abi.label.symbol import Symbol


class SymbolConverter:
    def __init__(self):
        pass

    @staticmethod
    def convert_to_symbol(symbol: str) -> Symbol:
        if symbol == "@function":
            return Symbol.FUNCTION
        elif symbol == "@object":
            return Symbol.OBJECT
        else:
            raise ValueError
