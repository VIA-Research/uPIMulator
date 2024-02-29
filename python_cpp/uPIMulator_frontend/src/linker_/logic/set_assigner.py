from abi.binary.executable import Executable
from abi.section.section_flag import SectionFlag
from abi.section.section_name import SectionName
from abi.section.section_type import SectionType
from converter.section_flag_converter import SectionFlagConverter
from converter.section_name_converter import SectionNameConverter
from converter.section_type_converter import SectionTypeConverter
from parser_.grammar.assemblyListener import assemblyListener
from parser_.grammar.assemblyParser import assemblyParser


class SetAssigner(assemblyListener):
    def __init__(self, executable: Executable):
        self._executable: Executable = executable

    def exitSection_directive(self, ctx: assemblyParser.Section_directiveContext) -> None:
        section_name = SectionNameConverter.convert_to_section_name(str(ctx.section_name().children[0])[1:])

        name = ctx.Identifier()
        if name is not None:
            name = str(name)
        else:
            name = ""

        section_flags = SectionFlagConverter.convert_to_section_flags(str(ctx.StringLiteral())[1:-1])
        section_type = SectionTypeConverter.convert_to_section_type(str(ctx.section_types().children[0]))

        self._executable.checkout_section(section_name, name, section_flags, section_type)

    def exitSet_directive(self, ctx: assemblyParser.Set_directiveContext) -> None:
        src_label_name = str(ctx.Identifier(i=0))
        dst_label_name = str(ctx.Identifier(i=1))

        src_label = self._executable.cur_section().label(src_label_name)
        assert src_label is not None

        self._executable.append_label(dst_label_name)
        dst_label = self._executable.cur_section().label(dst_label_name)
        assert dst_label is not None

        dst_label.set_address(src_label.address())

    def exitSize_directive(self, ctx: assemblyParser.Size_directiveContext) -> None:
        # TODO(bongjoon.hyun@gmail.com): implement this
        pass

    def exitStack_sizes_directive(self, ctx: assemblyParser.Stack_sizes_directiveContext):
        section_name = SectionName.STACK_SIZES

        name = ""
        section_flags = SectionFlagConverter.convert_to_section_flags(str(ctx.StringLiteral())[1:-1])
        section_type = SectionTypeConverter.convert_to_section_type(str(ctx.section_types().children[0]))

        self._executable.checkout_section(section_name, name, section_flags, section_type)

    def exitText_directive(self, ctx: assemblyParser.Text_directiveContext) -> None:
        self._executable.checkout_section(
            SectionName.TEXT, "", {SectionFlag.ALLOC, SectionFlag.EXECINSTR}, SectionType.PROG_BITS,
        )
