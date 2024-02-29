from abi.binary.executable import Executable
from abi.directive.ascii_directive import AsciiDirective
from abi.directive.asciz_directive import AscizDirective
from abi.directive.byte_directive import ByteDirective
from abi.directive.long_directive import LongDirective
from abi.directive.quad_directive import QuadDirective
from abi.directive.short_directive import ShortDirective
from abi.directive.zero_directive import ZeroDirective
from abi.section.section_flag import SectionFlag
from abi.section.section_name import SectionName
from abi.section.section_type import SectionType
from converter.section_flag_converter import SectionFlagConverter
from converter.section_name_converter import SectionNameConverter
from converter.section_type_converter import SectionTypeConverter
from encoder.ascii_encoder import AsciiEncoder
from encoder.byte import Byte
from parser_.grammar.assemblyListener import assemblyListener
from parser_.grammar.assemblyParser import assemblyParser
from util.config_loader import ConfigLoader


class LabelAssigner(assemblyListener):
    def __init__(self, executable: Executable):
        self._executable: Executable = executable

    def exitAscii_directive(self, ctx: assemblyParser.Ascii_directiveContext) -> None:
        string_literal = str(ctx.StringLiteral())[1:-1].replace("\\b", "■")

        i = 0
        characters = ""
        while i < len(string_literal):
            if string_literal[i] == "\\":
                octal = Byte(int(string_literal[i + 1 : i + 4], base=8))
                characters += AsciiEncoder().decode(octal)
                i += 3
            else:
                characters += string_literal[i]
                i += 1
        ascii_directive = AsciiDirective(characters)
        cur_label = self._executable.cur_label()
        cur_label.set_size(cur_label.size() + ascii_directive.size())

    def exitAsciz_directive(self, ctx: assemblyParser.Asciz_directiveContext) -> None:
        characters = str(ctx.StringLiteral())[1:-1]

        asciz_directive = AscizDirective(characters)
        cur_label = self._executable.cur_label()
        cur_label.set_size(cur_label.size() + asciz_directive.size())

    def exitByte_directive(self, ctx: assemblyParser.Byte_directiveContext) -> None:
        cur_label = self._executable.cur_label()
        cur_label.set_size(cur_label.size() + ByteDirective.size())

    def exitLong_directive(self, ctx: assemblyParser.Long_directiveContext) -> None:
        cur_label = self._executable.cur_label()
        cur_label.set_size(cur_label.size() + LongDirective.size())

    def exitQuad_directive(self, ctx: assemblyParser.Quad_directiveContext) -> None:
        cur_label = self._executable.cur_label()
        cur_label.set_size(cur_label.size() + QuadDirective.size())

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

    def exitShort_directive(self, ctx: assemblyParser.Short_directiveContext) -> None:
        cur_label = self._executable.cur_label()
        cur_label.set_size(cur_label.size() + ShortDirective.size())

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

    def exitZero_directive(self, ctx: assemblyParser.Zero_directiveContext) -> None:
        size, value = int(ctx.number(i=0).getText()), ctx.number(i=1)
        if value is None:
            value = 0
        else:
            value = int(value.getText())

        zero_directive = ZeroDirective(size, value)

        cur_label = self._executable.cur_label()
        cur_label.set_size(cur_label.size() + zero_directive.size())

    def exitInstruction(self, ctx: assemblyParser.InstructionContext):
        instruction_size = ConfigLoader.iram_data_width() // 8
        cur_label = self._executable.cur_label()
        cur_label.set_size(cur_label.size() + instruction_size)

    def exitLabel(self, ctx: assemblyParser.LabelContext):
        label_name = str(ctx.Identifier())

        # TODO(bongjoon.hyun@gmail.com): __sys_used_mram_end will be defined in the linker script
        if label_name != "__sys_used_mram_end":
            self._executable.append_label(str(ctx.Identifier()))
