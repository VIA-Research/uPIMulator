import os
from typing import List, Set, Union

import antlr4

from abi.binary.executable import Executable
from abi.binary.relocatable import Relocatable
from assembler.assembler import Assembler
from linker_.linker_script import LinkerScript
from linker_.logic.instruction_assigner import InstructionAssigner
from linker_.logic.label_assigner import LabelAssigner
from linker_.logic.liveness_analyzer import LivenessAnalyzer
from linker_.logic.set_assigner import SetAssigner
from parser_.parser import Parser
from util.path_collector import PathCollector


class Linker:
    def __init__(self, num_tasklets: int):
        self._libs: List[Relocatable] = self._init_libs()

    def link(self, filepath: str, data_prep_param: list, num_dpus: int):
        num_tasklets = Linker._num_tasklets(filepath)
        linker_script = LinkerScript(num_tasklets)

        relocatable = Relocatable(filepath)
        self._analyze_liveness(relocatable)
        relocatables = self._resolve_symbols(relocatable, linker_script)

        executable = Executable(filepath, relocatables)
        Linker._write_bin(executable, num_dpus)

        self._analyze_liveness(executable)
        self._assign_labels(executable)
        linker_script.assign_address(executable)
        self._assign_set(executable)
        self._assign_instruction(executable, linker_script)

        Assembler.assemble(executable, linker_script, data_prep_param, num_dpus)

    def _analyze_liveness(self, binary: Union[Relocatable, Executable]):
        document = Parser.parse_lines(binary.lines())
        liveness_analyzer = LivenessAnalyzer(binary)
        parse_tree_walker = antlr4.ParseTreeWalker()
        parse_tree_walker.walk(liveness_analyzer, document)

    def _resolve_symbols(self, relocatable: Relocatable, linker_script: LinkerScript) -> Set[Relocatable]:
        relocatables: Set[Relocatable] = {relocatable, self._crt0_relocatable()}

        unresolved_symbols = self._unresolved_symbols(relocatables, linker_script)
        while unresolved_symbols:
            has_resolved = False
            for unresolved_symbol in unresolved_symbols:
                for lib in self._libs:
                    if unresolved_symbol in lib.liveness().defs():
                        relocatables.add(lib)
                        has_resolved = True
            assert has_resolved

            unresolved_symbols = self._unresolved_symbols(relocatables, linker_script)

        return relocatables

    def _defs(self, relocatables: Set[Relocatable]) -> Set[str]:
        defs: Set[str] = set()
        for relocatable in relocatables:
            defs.update(relocatable.liveness().defs())
        return defs

    def _unresolved_symbols(self, relocatables: Set[Relocatable], linker_script: LinkerScript) -> Set[str]:
        defs = self._defs(relocatables)
        unresolved_symbols: Set[str] = set()
        for relocatable in relocatables:
            for unresolved_symbol in relocatable.liveness().unresolved_symbols():
                if (
                    unresolved_symbol not in defs
                    and unresolved_symbol not in linker_script.symbol_names()
                    and unresolved_symbol not in linker_script.constant_names()
                ):
                    unresolved_symbols.add(unresolved_symbol)
        return unresolved_symbols

    def _assign_labels(self, executable: Executable) -> None:
        document = Parser.parse_lines(executable.lines())
        label_assigner = LabelAssigner(executable)
        parse_tree_walker = antlr4.ParseTreeWalker()
        parse_tree_walker.walk(label_assigner, document)

    def _assign_set(self, executable: Executable) -> None:
        document = Parser.parse_lines(executable.lines())
        set_assigner = SetAssigner(executable)
        parse_tree_walker = antlr4.ParseTreeWalker()
        parse_tree_walker.walk(set_assigner, document)

    def _assign_instruction(self, executable: Executable, linker_script: LinkerScript) -> None:
        document = Parser.parse_lines(executable.lines())
        instruction_assigner = InstructionAssigner(executable, linker_script)
        parse_tree_walker = antlr4.ParseTreeWalker()
        parse_tree_walker.walk(instruction_assigner, document)

    def _init_libs(self) -> List[Relocatable]:
        relocatables: List[Relocatable] = []
        for lib_name in Linker.libs_names():
            asm_lib_path = os.path.join(PathCollector.asm_path_in_local(), lib_name)
            for root_path, _, filenames in os.walk(asm_lib_path):
                for filename in filenames:
                    if filename.split(".")[-1] == "S":
                        filepath = os.path.join(root_path, filename)
                        relocatable = Relocatable(filepath)
                        self._analyze_liveness(relocatable)
                        relocatables.append(relocatable)
        return relocatables

    def _crt0_relocatable(self) -> Relocatable:
        for lib in self._libs:
            crt0_filepath = os.path.join(PathCollector.asm_path_in_local(), "misc", "crt0.S")
            if lib.filepath() == crt0_filepath:
                return lib
        raise ValueError

    @staticmethod
    def _write_bin(executable: Executable, num_dpus: int) -> None:
        filepath = os.path.join(PathCollector.bin_path_in_local(), f"{num_dpus}_dpus", *executable.filepath().split(os.path.sep)[-2:])

        dirname = os.path.dirname(filepath)
        if not os.path.exists(dirname):
            os.makedirs(dirname)

        with open(filepath, "w") as file:
            file.writelines(executable.lines())

    @staticmethod
    def _num_tasklets(filepath: str) -> int:
        return int(filepath.split(os.path.sep)[-2].split(".")[-1])

    @staticmethod
    def libs_names() -> List[str]:
        return ["misc", "stdlib", "syslib"]
