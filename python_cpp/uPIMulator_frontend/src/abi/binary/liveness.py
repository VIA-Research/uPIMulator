from typing import Dict, Optional, Set

from abi.label.symbol import Symbol


class Liveness:
    def __init__(self):
        self._def_uses: Dict[str, Set[str]] = {}
        self._global_symbols: Set[str] = set()
        self._symbols: Dict[str, Symbol] = {}
        self._unresolved_symbols: Set[str] = set()

        self._cur_def: Optional[str] = None

    def defs(self) -> Set[str]:
        return set(self._def_uses.keys())

    def checkout_def(self, label_name: str) -> None:
        self.add_def(label_name)
        self._cur_def = label_name

    def add_def(self, label_name: str) -> None:
        if label_name not in self._def_uses:
            self._def_uses[label_name] = set()

            if label_name in self._unresolved_symbols:
                self._unresolved_symbols.remove(label_name)

    def uses(self, def_: str) -> Set[str]:
        if def_ in self._def_uses:
            return self._def_uses[def_].copy()
        else:
            return set()

    def add_use(self, label_name: str) -> None:
        assert self._cur_def is not None
        self._def_uses[self._cur_def].add(label_name)

        if label_name not in self._def_uses:
            self._unresolved_symbols.add(label_name)

    def global_symbols(self) -> Set[str]:
        return self._global_symbols.copy()

    def local_symbols(self) -> Set[str]:
        local_symbols = self.defs()
        for global_symbol in self.global_symbols():
            local_symbols.remove(global_symbol)
        return local_symbols

    def symbol(self, symbol_name: str) -> Optional[Symbol]:
        return self._symbols.get(symbol_name)

    def add_global_symbol(self, symbol_name: str) -> None:
        self._global_symbols.add(symbol_name)

    def add_symbol(self, symbol_name: str, symbol: Symbol) -> None:
        assert symbol_name not in self._symbols
        self._symbols[symbol_name] = symbol

    def unresolved_symbols(self) -> Set[str]:
        return self._unresolved_symbols.copy()
