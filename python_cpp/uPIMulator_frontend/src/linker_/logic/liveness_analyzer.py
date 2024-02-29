from typing import Union

from abi.binary.executable import Executable
from abi.binary.relocatable import Relocatable
from converter.symbol_converter import SymbolConverter
from parser_.grammar.assemblyListener import assemblyListener
from parser_.grammar.assemblyParser import assemblyParser


class LivenessAnalyzer(assemblyListener):
    def __init__(self, binary: Union[Relocatable, Executable]):
        self._binary: Union[Relocatable, Executable] = binary

    def exitGlobal_directive(self, ctx: assemblyParser.Global_directiveContext) -> None:
        symbol_name = str(ctx.Identifier())

        # TODO(bongjoon.hyun@gmail.com): __sys_used_mram_end will be defined in the linker script
        if symbol_name != "__sys_used_mram_end":
            self._binary.liveness().add_global_symbol(symbol_name)

    def exitSet_directive(self, ctx: assemblyParser.Set_directiveContext) -> None:
        self._binary.liveness().add_def(str(ctx.Identifier(i=0)))
        self._binary.liveness().add_use(str(ctx.Identifier(i=1)))

    def exitSize_directive(self, ctx: assemblyParser.Size_directiveContext) -> None:
        self._binary.liveness().add_use(str(ctx.Identifier()))

    def exitType_directive(self, ctx: assemblyParser.Type_directiveContext):
        self._binary.liveness().add_symbol(
            str(ctx.Identifier()), SymbolConverter.convert_to_symbol(ctx.symbol_type().getText()),
        )

    def exitPrimary_expression(self, ctx: assemblyParser.Primary_expressionContext) -> None:
        identifier = ctx.Identifier()
        if identifier:
            self._binary.liveness().add_use(str(identifier))

    def exitLabel(self, ctx: assemblyParser.LabelContext) -> None:
        label_name = str(ctx.Identifier())
        assert label_name not in self._binary.liveness().defs()

        # TODO(bongjoon.hyun@gmail.com): __sys_used_mram_end will be defined in the linker script
        if label_name != "__sys_used_mram_end":
            self._binary.liveness().checkout_def(label_name)
