from abi.binary.executable import Executable
from abi.directive.ascii_directive import AsciiDirective
from abi.directive.asciz_directive import AscizDirective
from abi.directive.byte_directive import ByteDirective
from abi.directive.long_directive import LongDirective
from abi.directive.quad_directive import QuadDirective
from abi.directive.short_directive import ShortDirective
from abi.directive.zero_directive import ZeroDirective
from abi.isa.instruction.condition import Condition
from abi.isa.instruction.endian import Endian
from abi.isa.instruction.instruction import Instruction
from abi.isa.instruction.op_code import OpCode
from abi.isa.instruction.suffix import Suffix
from abi.isa.register.sp_register import SPRegister
from abi.section.section_flag import SectionFlag
from abi.section.section_name import SectionName
from abi.section.section_type import SectionType
from abi.word.instruction_word import InstructionWord
from converter.condition_converter import ConditionConverter
from converter.endian_converter import EndianConverter
from converter.op_code_converter import OpCodeConverter
from converter.register_converter import RegisterConverter
from converter.section_flag_converter import SectionFlagConverter
from converter.section_name_converter import SectionNameConverter
from converter.section_type_converter import SectionTypeConverter
from encoder.ascii_encoder import AsciiEncoder
from encoder.byte import Byte
from linker_.linker_script import LinkerScript
from parser_.grammar.assemblyListener import assemblyListener
from parser_.grammar.assemblyParser import assemblyParser


class InstructionAssigner(assemblyListener):
    def __init__(self, executable: Executable, linker_script: LinkerScript):
        self._executable = executable
        self._linker_script = linker_script

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

        self._executable.append_assembler_instruction(ascii_directive)

    def exitAsciz_directive(self, ctx: assemblyParser.Asciz_directiveContext) -> None:
        characters = str(ctx.StringLiteral())[1:-1]
        asciz_directive = AscizDirective(characters)
        self._executable.append_assembler_instruction(asciz_directive)

    def exitByte_directive(self, ctx: assemblyParser.Byte_directiveContext) -> None:
        value = self._evaluate_program_counter(ctx.program_counter())
        byte_directive = ByteDirective(value)
        self._executable.append_assembler_instruction(byte_directive)

    def exitLong_directive(self, ctx: assemblyParser.Long_directiveContext) -> None:
        value = self._evaluate_program_counter(ctx.program_counter())
        long_directive = LongDirective(value)
        self._executable.append_assembler_instruction(long_directive)

    def exitQuad_directive(self, ctx: assemblyParser.Quad_directiveContext) -> None:
        value = self._evaluate_program_counter(ctx.program_counter())
        quad_directive = QuadDirective(value)
        self._executable.append_assembler_instruction(quad_directive)

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
        value = self._evaluate_program_counter(ctx.program_counter())
        short_directive = ShortDirective(value)
        self._executable.append_assembler_instruction(short_directive)

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
        self._executable.append_assembler_instruction(zero_directive)

    def exitRici_instruction(self, ctx: assemblyParser.Rici_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rici_op_code().children[0])[1:])
        suffix = Suffix.RICI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = self._evaluate_program_counter(ctx.program_counter(i=0))
        condition = ConditionConverter.convert_to_condition(str(ctx.condition().children[0]))
        pc = self._evaluate_program_counter(ctx.program_counter(i=1))

        instruction = Instruction(op_code, suffix, ra=ra, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitRri_instruction(self, ctx: assemblyParser.Rri_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        suffix = Suffix.RRI
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = self._evaluate_number(ctx.number())

        instruction = Instruction(op_code, suffix, rc=rc, ra=ra, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert rc == instruction.rc()
        assert ra == instruction.ra()
        assert (imm % (2 ** instruction.imm().width())) == (
            instruction.imm().value() % (2 ** instruction.imm().width())
        )

        self._executable.append_assembler_instruction(instruction)

    def exitRric_instruction(self, ctx: assemblyParser.Rric_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = self._evaluate_number(ctx.number())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())

        if condition != Condition.FALSE:
            suffix = Suffix.RRIC
        else:
            suffix = Suffix.RRIF

        instruction = Instruction(op_code, suffix, rc=rc, ra=ra, imm=imm, condition=condition)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert rc == instruction.rc()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()

        self._executable.append_assembler_instruction(instruction)

    def exitRrici_instruction(self, ctx: assemblyParser.Rrici_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        suffix = Suffix.RRICI
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = self._evaluate_number(ctx.number())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, rc=rc, ra=ra, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert rc == instruction.rc()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitRrr_instruction(self, ctx: assemblyParser.Rrr_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        suffix = Suffix.RRR
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())

        instruction = Instruction(op_code, suffix, rc=rc, ra=ra, rb=rb)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert rc == instruction.rc()
        assert ra == instruction.ra()
        assert rb == instruction.rb()

        self._executable.append_assembler_instruction(instruction)

    def exitRrrc_instruction(self, ctx: assemblyParser.Rrrc_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        suffix = Suffix.RRRC
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())

        instruction = Instruction(op_code, suffix, rc=rc, ra=ra, rb=rb, condition=condition)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert rc == instruction.rc()
        assert ra == instruction.ra()
        assert rb == instruction.rb()
        assert condition == instruction.condition()

        self._executable.append_assembler_instruction(instruction)

    def exitRrrci_instruction(self, ctx: assemblyParser.Rrrci_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        suffix = Suffix.RRRCI
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, rc=rc, ra=ra, rb=rb, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert rc == instruction.rc()
        assert ra == instruction.ra()
        assert rb == instruction.rb()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitZri_instruction(self, ctx: assemblyParser.Zri_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        suffix = Suffix.ZRI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = self._evaluate_number(ctx.number())

        instruction = Instruction(op_code, suffix, ra=ra, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()

        self._executable.append_assembler_instruction(instruction)

    def exitZric_instruction(self, ctx: assemblyParser.Zric_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = self._evaluate_program_counter(ctx.program_counter())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())

        if condition != Condition.FALSE:
            suffix = Suffix.ZRIC
        else:
            suffix = Suffix.ZRIF

        instruction = Instruction(op_code, suffix, ra=ra, imm=imm, condition=condition)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert condition == instruction.condition()

        self._executable.append_assembler_instruction(instruction)

    def exitZrici_instruction(self, ctx: assemblyParser.Zrici_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        suffix = Suffix.ZRICI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = self._evaluate_number(ctx.number())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, ra=ra, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitZrr_instruction(self, ctx: assemblyParser.Zrr_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        suffix = Suffix.ZRR
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())

        instruction = Instruction(op_code, suffix, ra=ra, rb=rb)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert rb == instruction.rb()

        self._executable.append_assembler_instruction(instruction)

    def exitZrrc_instruction(self, ctx: assemblyParser.Zrrc_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        suffix = Suffix.ZRRC
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())

        instruction = Instruction(op_code, suffix, ra=ra, rb=rb, condition=condition)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert rb == instruction.rb()
        assert condition == instruction.condition()

        self._executable.append_assembler_instruction(instruction)

    def exitZrrci_instruction(self, ctx: assemblyParser.Zrrci_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        suffix = Suffix.ZRRCI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, ra=ra, rb=rb, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert rb == instruction.rb()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitS_rri_instruction(self, ctx: assemblyParser.S_rri_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        suffix = Suffix.S_RRI
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = self._evaluate_number(ctx.number())

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()

        self._executable.append_assembler_instruction(instruction)

    def exitS_rric_instruction(self, ctx: assemblyParser.S_rric_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = self._evaluate_number(ctx.number())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())

        if condition != Condition.FALSE:
            suffix = Suffix.S_RRIC
        else:
            suffix = Suffix.S_RRIF

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra, imm=imm, condition=condition)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()

        self._executable.append_assembler_instruction(instruction)

    def exitS_rrici_instruction(self, ctx: assemblyParser.S_rrici_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        suffix = Suffix.S_RRICI
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = self._evaluate_number(ctx.number())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitS_rrr_instruction(self, ctx: assemblyParser.S_rrr_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        suffix = Suffix.S_RRR
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra, rb=rb)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert rb == instruction.rb()

        self._executable.append_assembler_instruction(instruction)

    def exitS_rrrc_instruction(self, ctx: assemblyParser.S_rrrc_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        suffix = Suffix.S_RRRC
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra, rb=rb, condition=condition)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert rb == instruction.rb()
        assert condition == instruction.condition()

        self._executable.append_assembler_instruction(instruction)

    def exitS_rrrci_instruction(self, ctx: assemblyParser.S_rrrci_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        suffix = Suffix.S_RRRCI
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra, rb=rb, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert rb == instruction.rb()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitU_rri_instruction(self, ctx: assemblyParser.U_rri_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        suffix = Suffix.U_RRI
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = self._evaluate_number(ctx.number())

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()

        self._executable.append_assembler_instruction(instruction)

    def exitU_rric_instruction(self, ctx: assemblyParser.U_rric_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = self._evaluate_number(ctx.number())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())

        if condition != Condition.FALSE:
            suffix = Suffix.U_RRIC
        else:
            suffix = Suffix.U_RRIF

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra, imm=imm, condition=condition)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()

        self._executable.append_assembler_instruction(instruction)

    def exitU_rrici_instruction(self, ctx: assemblyParser.U_rrici_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        suffix = Suffix.U_RRICI
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = self._evaluate_number(ctx.number())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitU_rrr_instruction(self, ctx: assemblyParser.U_rrr_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        suffix = Suffix.U_RRR
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra, rb=rb)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert rb == instruction.rb()

        self._executable.append_assembler_instruction(instruction)

    def exitU_rrrc_instruction(self, ctx: assemblyParser.U_rrrc_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        suffix = Suffix.U_RRRC
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra, rb=rb, condition=condition)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert rb == instruction.rb()
        assert condition == instruction.condition()

        self._executable.append_assembler_instruction(instruction)

    def exitU_rrrci_instruction(self, ctx: assemblyParser.U_rrrci_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        suffix = Suffix.U_RRRCI
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra, rb=rb, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert rb == instruction.rb()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitRr_instruction(self, ctx: assemblyParser.Rr_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rr_op_code().children[0])[1:])
        suffix = Suffix.RR
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())

        instruction = Instruction(op_code, suffix, rc=rc, ra=ra)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert rc == instruction.rc()
        assert ra == instruction.ra()

        self._executable.append_assembler_instruction(instruction)

    def exitRrc_instruction(self, ctx: assemblyParser.Rrc_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rr_op_code().children[0])[1:])
        suffix = Suffix.RRC
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())

        instruction = Instruction(op_code, suffix, rc=rc, ra=ra, condition=condition)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert rc == instruction.rc()
        assert ra == instruction.ra()
        assert condition == instruction.condition()

        self._executable.append_assembler_instruction(instruction)

    def exitRrci_instruction(self, ctx: assemblyParser.Rrci_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rr_op_code().children[0])[1:])
        suffix = Suffix.RRCI
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, rc=rc, ra=ra, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert rc == instruction.rc()
        assert ra == instruction.ra()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitZr_instruction(self, ctx: assemblyParser.Zr_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rr_op_code().children[0])[1:])
        suffix = Suffix.ZR
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())

        instruction = Instruction(op_code, suffix, ra=ra)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()

        self._executable.append_assembler_instruction(instruction)

    def exitZrc_instruction(self, ctx: assemblyParser.Zrc_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rr_op_code().children[0])[1:])
        suffix = Suffix.ZRC
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())

        instruction = Instruction(op_code, suffix, ra=ra, condition=condition)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert condition == instruction.condition()

        self._executable.append_assembler_instruction(instruction)

    def exitZrci_instruction(self, ctx: assemblyParser.Zrci_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rr_op_code().children[0])[1:])
        suffix = Suffix.ZRCI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, ra=ra, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitS_rr_instruction(self, ctx: assemblyParser.S_rr_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rr_op_code().children[0])[1:])
        suffix = Suffix.S_RR
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()

        self._executable.append_assembler_instruction(instruction)

    def exitS_rrc_instruction(self, ctx: assemblyParser.S_rrc_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rr_op_code().children[0])[1:])
        suffix = Suffix.S_RRC
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra, condition=condition)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert condition == instruction.condition()

        self._executable.append_assembler_instruction(instruction)

    def exitS_rrci_instruction(self, ctx: assemblyParser.S_rrci_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rr_op_code().children[0])[1:])
        suffix = Suffix.S_RRCI
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitU_rr_instruction(self, ctx: assemblyParser.U_rr_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rr_op_code().children[0])[1:])
        suffix = Suffix.U_RR
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()

        self._executable.append_assembler_instruction(instruction)

    def exitU_rrc_instruction(self, ctx: assemblyParser.U_rrc_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rr_op_code().children[0])[1:])
        suffix = Suffix.U_RRC
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra, condition=condition)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert condition == instruction.condition()

        self._executable.append_assembler_instruction(instruction)

    def exitU_rrci_instruction(self, ctx: assemblyParser.U_rrci_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rr_op_code().children[0])[1:])
        suffix = Suffix.U_RRCI
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitDrdici_instruction(self, ctx: assemblyParser.Drdici_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.drdici_op_code().children[0])[1:])
        suffix = Suffix.DRDICI
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister(i=0)))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        db = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister(i=1)))
        imm = self._evaluate_number(ctx.number())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra, db=db, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert db == instruction.db()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitRrri_instruction(self, ctx: assemblyParser.Rrri_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rrri_op_code().children[0])[1:])
        suffix = Suffix.RRRI
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())
        imm = self._evaluate_number(ctx.number())

        instruction = Instruction(op_code, suffix, rc=rc, ra=ra, rb=rb, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert rc == instruction.rc()
        assert ra == instruction.ra()
        assert rb == instruction.rb()
        assert imm == instruction.imm().value()

        self._executable.append_assembler_instruction(instruction)

    def exitRrrici_instruction(self, ctx: assemblyParser.Rrrici_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rrri_op_code().children[0])[1:])
        suffix = Suffix.RRRICI
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())
        imm = self._evaluate_number(ctx.number())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, rc=rc, ra=ra, rb=rb, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert rc == instruction.rc()
        assert ra == instruction.ra()
        assert rb == instruction.rb()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitZrri_instruction(self, ctx: assemblyParser.Zrri_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rrri_op_code().children[0])[1:])
        suffix = Suffix.ZRRI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())
        imm = self._evaluate_number(ctx.number())

        instruction = Instruction(op_code, suffix, ra=ra, rb=rb, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert rb == instruction.rb()
        assert imm == instruction.imm().value()

        self._executable.append_assembler_instruction(instruction)

    def exitZrrici_instruction(self, ctx: assemblyParser.Zrrici_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rrri_op_code().children[0])[1:])
        suffix = Suffix.ZRRICI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())
        imm = self._evaluate_number(ctx.number())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, ra=ra, rb=rb, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert rb == instruction.rb()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitS_rrri_instruction(self, ctx: assemblyParser.S_rrri_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rrri_op_code().children[0])[1:])
        suffix = Suffix.S_RRRI
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())
        imm = self._evaluate_number(ctx.number())

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra, rb=rb, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert rb == instruction.rb()
        assert imm == instruction.imm().value()

        self._executable.append_assembler_instruction(instruction)

    def exitS_rrrici_instruction(self, ctx: assemblyParser.S_rrrici_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rrri_op_code().children[0])[1:])
        suffix = Suffix.S_RRRICI
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())
        imm = self._evaluate_number(ctx.number())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra, rb=rb, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert rb == instruction.rb()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitU_rrri_instruction(self, ctx: assemblyParser.U_rrri_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rrri_op_code().children[0])[1:])
        suffix = Suffix.U_RRRI
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())
        imm = self._evaluate_number(ctx.number())

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra, rb=rb, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert rb == instruction.rb()
        assert imm == instruction.imm().value()

        self._executable.append_assembler_instruction(instruction)

    def exitU_rrrici_instruction(self, ctx: assemblyParser.U_rrrici_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rrri_op_code().children[0])[1:])
        suffix = Suffix.U_RRRICI
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())
        imm = self._evaluate_number(ctx.number())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra, rb=rb, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert rb == instruction.rb()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitRir_instruction(self, ctx: assemblyParser.Rir_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        suffix = Suffix.RIR
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        imm = self._evaluate_number(ctx.number())
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())

        instruction = Instruction(op_code, suffix, rc=rc, imm=imm, ra=ra)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert rc == instruction.rc()
        assert imm == instruction.imm().value()
        assert ra == instruction.ra()

        self._executable.append_assembler_instruction(instruction)

    def exitRirc_instruction(self, ctx: assemblyParser.Rirc_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        imm = self._evaluate_number(ctx.number())
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())

        if condition != Condition.FALSE:
            suffix = Suffix.RIRC
        else:
            suffix = Suffix.RRIF

        instruction = Instruction(op_code, suffix, rc=rc, imm=imm, ra=ra, condition=condition)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert rc == instruction.rc()
        assert imm == instruction.imm().value()
        assert ra == instruction.ra()
        assert condition == instruction.condition()

        self._executable.append_assembler_instruction(instruction)

    def exitRirci_instruction(self, ctx: assemblyParser.Rirci_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        suffix = Suffix.RIRCI
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        imm = self._evaluate_number(ctx.number())
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, rc=rc, imm=imm, ra=ra, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert rc == instruction.rc()
        assert imm == instruction.imm().value()
        assert ra == instruction.ra()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitZir_instruction(self, ctx: assemblyParser.Zir_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        suffix = Suffix.ZIR
        imm = self._evaluate_number(ctx.number())
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())

        instruction = Instruction(op_code, suffix, imm=imm, ra=ra)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert imm == instruction.imm().value()
        assert ra == instruction.ra()

        self._executable.append_assembler_instruction(instruction)

    def exitZirc_instruction(self, ctx: assemblyParser.Zirc_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        imm = self._evaluate_number(ctx.number())
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())

        if condition != Condition.FALSE:
            suffix = Suffix.ZIRC
        else:
            suffix = Suffix.ZRIF

        instruction = Instruction(op_code, suffix, imm=imm, ra=ra, condition=condition)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert imm == instruction.imm().value()
        assert ra == instruction.ra()
        assert condition == instruction.condition()

        self._executable.append_assembler_instruction(instruction)

    def exitZirci_instruction(self, ctx: assemblyParser.Zirci_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        suffix = Suffix.ZIRCI
        imm = self._evaluate_number(ctx.number())
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, imm=imm, ra=ra, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert imm == instruction.imm().value()
        assert ra == instruction.ra()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitS_rirc_instruction(self, ctx: assemblyParser.S_rirc_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        suffix = Suffix.S_RIRC
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        imm = self._evaluate_number(ctx.number())
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())

        instruction = Instruction(op_code, suffix, dc=dc, imm=imm, ra=ra, condition=condition)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert imm == instruction.imm().value()
        assert ra == instruction.ra()
        assert condition == instruction.condition()

        self._executable.append_assembler_instruction(instruction)

    def exitS_rirci_instruction(self, ctx: assemblyParser.S_rirci_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        suffix = Suffix.S_RIRCI
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        imm = self._evaluate_number(ctx.number())
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, dc=dc, imm=imm, ra=ra, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert imm == instruction.imm().value()
        assert ra == instruction.ra()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitU_rirc_instruction(self, ctx: assemblyParser.U_rirc_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        suffix = Suffix.U_RIRC
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        imm = self._evaluate_number(ctx.number())
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())

        instruction = Instruction(op_code, suffix, dc=dc, imm=imm, ra=ra, condition=condition)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert imm == instruction.imm().value()
        assert ra == instruction.ra()
        assert condition == instruction.condition()

        self._executable.append_assembler_instruction(instruction)

    def exitU_rirci_instruction(self, ctx: assemblyParser.U_rirci_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.rri_op_code().children[0])[1:])
        suffix = Suffix.U_RIRCI
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        imm = self._evaluate_number(ctx.number())
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, dc=dc, imm=imm, ra=ra, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert imm == instruction.imm().value()
        assert ra == instruction.ra()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitR_instruction(self, ctx: assemblyParser.R_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.r_op_code().children[0])[1:])
        suffix = Suffix.R
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))

        instruction = Instruction(op_code, suffix, rc=rc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert rc == instruction.rc()

        self._executable.append_assembler_instruction(instruction)

    def exitRci_instruction(self, ctx: assemblyParser.Rci_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.r_op_code().children[0])[1:])
        suffix = Suffix.RCI
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, rc=rc, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert rc == instruction.rc()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitZ_instruction(self, ctx: assemblyParser.Z_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.r_op_code().children[0])[1:])
        suffix = Suffix.Z

        instruction = Instruction(op_code, suffix)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()

        self._executable.append_assembler_instruction(instruction)

    def exitZci_instruction(self, ctx: assemblyParser.Zci_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.r_op_code().children[0])[1:])
        suffix = Suffix.ZCI
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitS_r_instruction(self, ctx: assemblyParser.S_r_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.r_op_code().children[0])[1:])
        suffix = Suffix.S_R
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))

        instruction = Instruction(op_code, suffix, dc=dc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()

        self._executable.append_assembler_instruction(instruction)

    def exitS_rci_instruction(self, ctx: assemblyParser.S_rci_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.r_op_code().children[0])[1:])
        suffix = Suffix.S_RCI
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, dc=dc, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitU_r_instruction(self, ctx: assemblyParser.U_r_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.r_op_code().children[0])[1:])
        suffix = Suffix.U_R
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))

        instruction = Instruction(op_code, suffix, dc=dc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()

        self._executable.append_assembler_instruction(instruction)

    def exitU_rci_instruction(self, ctx: assemblyParser.U_rci_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.r_op_code().children[0])[1:])
        suffix = Suffix.U_RCI
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, dc=dc, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitCi_instruction(self, ctx: assemblyParser.Ci_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.ci_op_code().children[0])[1:])
        suffix = Suffix.CI
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitI_instruction(self, ctx: assemblyParser.I_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.i_op_code().children[0])[1:])
        suffix = Suffix.I
        imm = self._evaluate_number(ctx.number())

        instruction = Instruction(op_code, suffix, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert imm == instruction.imm().value()

        self._executable.append_assembler_instruction(instruction)

    def exitDdci_instruction(self, ctx: assemblyParser.Ddci_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.ddci_op_code().children[0])[1:])
        suffix = Suffix.DDCI
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister(i=0)))
        db = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister(i=1)))
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, dc=dc, db=db, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert db == instruction.db()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitErri_instruction(self, ctx: assemblyParser.Erri_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.load_op_code().children[0])[1:])
        suffix = Suffix.ERRI
        endian = EndianConverter.convert_to_endian(ctx.endian().getText())
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        off = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, endian=endian, rc=rc, ra=ra, off=off)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert rc == instruction.rc()
        assert ra == instruction.ra()
        assert off == instruction.off().value()

        self._executable.append_assembler_instruction(instruction)

    def exitEdri_instruction(self, ctx: assemblyParser.Edri_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.load_op_code().children[0])[1:])
        suffix = Suffix.EDRI
        endian = EndianConverter.convert_to_endian(ctx.endian().getText())
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        off = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, endian=endian, dc=dc, ra=ra, off=off)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert off == instruction.off().value()

        self._executable.append_assembler_instruction(instruction)

    def exitS_erri_instruction(self, ctx: assemblyParser.S_erri_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.load_op_code().children[0])[1:])
        suffix = Suffix.S_ERRI
        endian = EndianConverter.convert_to_endian(ctx.endian().getText())
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        off = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, endian=endian, dc=dc, ra=ra, off=off)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert off == instruction.off().value()

        self._executable.append_assembler_instruction(instruction)

    def exitU_erri_instruction(self, ctx: assemblyParser.U_erri_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.load_op_code().children[0])[1:])
        suffix = Suffix.U_ERRI
        endian = EndianConverter.convert_to_endian(ctx.endian().getText())
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        off = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, endian=endian, dc=dc, ra=ra, off=off)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert off == instruction.off().value()

        self._executable.append_assembler_instruction(instruction)

    def exitErii_instruction(self, ctx: assemblyParser.Erii_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.store_op_code().children[0])[1:])
        suffix = Suffix.ERII
        endian = EndianConverter.convert_to_endian(ctx.endian().getText())
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        off = self._evaluate_number(ctx.number(i=0))
        imm = self._evaluate_number(ctx.number(i=1))

        instruction = Instruction(op_code, suffix, endian=endian, ra=ra, off=off, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert off == instruction.off().value()

        self._executable.append_assembler_instruction(instruction)

    def exitErir_instruction(self, ctx: assemblyParser.Erir_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.store_op_code().children[0])[1:])
        suffix = Suffix.ERIR
        endian = EndianConverter.convert_to_endian(ctx.endian().getText())
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        off = self._evaluate_program_counter(ctx.program_counter())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())

        instruction = Instruction(op_code, suffix, endian=endian, ra=ra, off=off, rb=rb)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert ra == instruction.ra()
        assert off == instruction.off().value()
        assert rb == instruction.rb()

        self._executable.append_assembler_instruction(instruction)

    def exitErid_instruction(self, ctx: assemblyParser.Erid_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.store_op_code().children[0])[1:])
        suffix = Suffix.ERID
        endian = EndianConverter.convert_to_endian(ctx.endian().getText())
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        off = self._evaluate_program_counter(ctx.program_counter())
        db = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))

        instruction = Instruction(op_code, suffix, endian=endian, ra=ra, off=off, db=db)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert ra == instruction.ra()
        assert off == instruction.off().value()
        assert db == instruction.db()

        self._executable.append_assembler_instruction(instruction)

    def exitDma_rri_instruction(self, ctx: assemblyParser.Dma_rri_instructionContext) -> None:
        op_code = OpCodeConverter.convert_to_op_code(str(ctx.dma_op_code().children[0])[1:])
        suffix = Suffix.DMA_RRI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())
        imm = self._evaluate_number(ctx.number())

        instruction = Instruction(op_code, suffix, ra=ra, rb=rb, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert rb == instruction.rb()
        assert imm == instruction.imm().value()

        self._executable.append_assembler_instruction(instruction)

    def exitAndn_rrif_instruction(self, ctx: assemblyParser.Andn_rrif_instructionContext) -> None:
        op_code = OpCode.ANDN
        suffix = Suffix.RRIF
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        imm = self._evaluate_number(ctx.number())
        condition = Condition.FALSE

        instruction = Instruction(op_code, suffix, rc=rc, ra=ra, imm=imm, condition=condition)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert rc == instruction.rc()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()

        self._executable.append_assembler_instruction(instruction)

    def exitNand_rrif_instruction(self, ctx: assemblyParser.Nand_rrif_instructionContext) -> None:
        op_code = OpCode.NAND
        suffix = Suffix.RRIF
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        imm = self._evaluate_number(ctx.number())
        condition = Condition.FALSE

        instruction = Instruction(op_code, suffix, rc=rc, ra=ra, imm=imm, condition=condition)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert rc == instruction.rc()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()

        self._executable.append_assembler_instruction(instruction)

    def exitNor_rrif_instruction(self, ctx: assemblyParser.Nor_rrif_instructionContext) -> None:
        op_code = OpCode.NOR
        suffix = Suffix.RRIF
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        imm = self._evaluate_number(ctx.number())
        condition = Condition.FALSE

        instruction = Instruction(op_code, suffix, rc=rc, ra=ra, imm=imm, condition=condition)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert rc == instruction.rc()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()

        self._executable.append_assembler_instruction(instruction)

    def exitNxor_rrif_instruction(self, ctx: assemblyParser.Nxor_rrif_instructionContext) -> None:
        op_code = OpCode.NXOR
        suffix = Suffix.RRIF
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        imm = self._evaluate_number(ctx.number())
        condition = Condition.FALSE

        instruction = Instruction(op_code, suffix, rc=rc, ra=ra, imm=imm, condition=condition)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert rc == instruction.rc()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()

        self._executable.append_assembler_instruction(instruction)

    def exitOrn_rrif_instruction(self, ctx: assemblyParser.Orn_rrif_instructionContext) -> None:
        op_code = OpCode.ORN
        suffix = Suffix.RRIF
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        imm = self._evaluate_number(ctx.number())
        condition = Condition.FALSE

        instruction = Instruction(op_code, suffix, rc=rc, ra=ra, imm=imm, condition=condition)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert rc == instruction.rc()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()

        self._executable.append_assembler_instruction(instruction)

    def exitHash_rrif_instruction(self, ctx: assemblyParser.Hash_rrif_instructionContext) -> None:
        op_code = OpCode.HASH
        suffix = Suffix.RRIF
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        imm = self._evaluate_number(ctx.number())
        condition = Condition.FALSE

        instruction = Instruction(op_code, suffix, rc=rc, ra=ra, imm=imm, condition=condition)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert rc == instruction.rc()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()

        self._executable.append_assembler_instruction(instruction)

    def exitMove_ri_instruction(self, ctx: assemblyParser.Move_ri_instructionContext) -> None:
        op_code = OpCode.OR
        suffix = Suffix.RRI
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        ra = SPRegister.ZERO
        imm = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, rc=rc, ra=ra, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert rc == instruction.rc()
        assert ra == instruction.ra()
        assert (imm % (2 ** instruction.imm().width())) == (
            instruction.imm().value() % (2 ** instruction.imm().width())
        )

        self._executable.append_assembler_instruction(instruction)

    def exitMove_rici_instruction(self, ctx: assemblyParser.Move_rici_instructionContext) -> None:
        op_code = OpCode.OR
        suffix = Suffix.RRICI
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        ra = SPRegister.ZERO
        imm = self._evaluate_number(ctx.number())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, rc=rc, ra=ra, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert rc == instruction.rc()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitMove_rr_instruction(self, ctx: assemblyParser.Move_rr_instructionContext) -> None:
        op_code = OpCode.OR
        suffix = Suffix.RRIF
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = 0
        condition = Condition.FALSE

        instruction = Instruction(op_code, suffix, rc=rc, ra=ra, imm=imm, condition=condition)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert rc == instruction.rc()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()

        self._executable.append_assembler_instruction(instruction)

    def exitMove_rrci_instruction(self, ctx: assemblyParser.Move_rrci_instructionContext) -> None:
        op_code = OpCode.OR
        suffix = Suffix.RRICI
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = 0
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, rc=rc, ra=ra, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert rc == instruction.rc()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitMove_s_ri_instruction(self, ctx: assemblyParser.Move_s_ri_instructionContext) -> None:
        op_code = OpCode.AND
        suffix = Suffix.S_RRI
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = SPRegister.LNEG
        imm = self._evaluate_number(ctx.number())

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert (imm % (2 ** instruction.imm().width())) == (
            instruction.imm().value() % (2 ** instruction.imm().width())
        )

        self._executable.append_assembler_instruction(instruction)

    def exitMove_s_rici_instruction(self, ctx: assemblyParser.Move_s_rici_instructionContext) -> None:
        op_code = OpCode.OR
        suffix = Suffix.S_RRICI
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = SPRegister.ZERO
        imm = self._evaluate_number(ctx.number())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitMove_s_rr_instruction(self, ctx: assemblyParser.Move_s_rr_instructionContext) -> None:
        op_code = OpCode.OR
        suffix = Suffix.S_RRIF
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = 0
        condition = Condition.FALSE

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra, imm=imm, condition=condition)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()

        self._executable.append_assembler_instruction(instruction)

    def exitMove_s_rrci_instruction(self, ctx: assemblyParser.Move_s_rrci_instructionContext) -> None:
        op_code = OpCode.OR
        suffix = Suffix.S_RRICI
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = 0
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitMove_u_ri_instruction(self, ctx: assemblyParser.Move_u_ri_instructionContext) -> None:
        op_code = OpCode.AND
        suffix = Suffix.U_RRI
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = SPRegister.LNEG
        imm = self._evaluate_number(ctx.number())

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert (imm % (2 ** instruction.imm().width())) == (
            instruction.imm().value() % (2 ** instruction.imm().width())
        )

        self._executable.append_assembler_instruction(instruction)

    def exitMove_u_rici_instruction(self, ctx: assemblyParser.Move_u_rici_instructionContext) -> None:
        op_code = OpCode.OR
        suffix = Suffix.U_RRICI
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = SPRegister.ZERO
        imm = self._evaluate_number(ctx.number())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitMove_u_rr_instruction(self, ctx: assemblyParser.Move_u_rr_instructionContext) -> None:
        op_code = OpCode.OR
        suffix = Suffix.U_RRIF
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = 0
        condition = Condition.FALSE

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra, imm=imm, condition=condition)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()

        self._executable.append_assembler_instruction(instruction)

    def exitMove_u_rrci_instruction(self, ctx: assemblyParser.Move_u_rrci_instructionContext) -> None:
        op_code = OpCode.OR
        suffix = Suffix.U_RRICI
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = 0
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitNeg_rr_instruction(self, ctx: assemblyParser.Neg_rr_instructionContext) -> None:
        op_code = OpCode.SUB
        suffix = Suffix.RIR
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        imm = 0
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())

        instruction = Instruction(op_code, suffix, rc=rc, imm=imm, ra=ra)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert rc == instruction.rc()
        assert imm == instruction.imm().value()
        assert ra == instruction.ra()

        self._executable.append_assembler_instruction(instruction)

    def exitNeg_rrci_instruction(self, ctx: assemblyParser.Neg_rrci_instructionContext) -> None:
        op_code = OpCode.SUB
        suffix = Suffix.RIRCI
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        imm = 0
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, rc=rc, imm=imm, ra=ra, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert rc == instruction.rc()
        assert imm == instruction.imm().value()
        assert ra == instruction.ra()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitNot_rr_instruction(self, ctx: assemblyParser.Not_rr_instructionContext) -> None:
        op_code = OpCode.XOR
        suffix = Suffix.RRI
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = -1

        instruction = Instruction(op_code, suffix, rc=rc, ra=ra, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert rc == instruction.rc()
        assert ra == instruction.ra()
        assert (imm % (2 ** instruction.imm().width())) == (
            instruction.imm().value() % (2 ** instruction.imm().width())
        )

        self._executable.append_assembler_instruction(instruction)

    def exitNot_rrci_instruction(self, ctx: assemblyParser.Not_rrci_instructionContext) -> None:
        op_code = OpCode.XOR
        suffix = Suffix.RRICI
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = -1
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, rc=rc, ra=ra, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert rc == instruction.rc()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitNot_zrci_instruction(self, ctx: assemblyParser.Not_zrci_instructionContext) -> None:
        op_code = OpCode.XOR
        suffix = Suffix.ZRICI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = -1
        condition = ConditionConverter.convert_to_condition(ctx.condition().getText())
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, ra=ra, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitJeq_rii_instruction(self, ctx: assemblyParser.Jeq_rii_instructionContext) -> None:
        op_code = OpCode.SUB
        suffix = Suffix.ZRICI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = self._evaluate_program_counter(ctx.program_counter(i=0))
        condition = Condition.Z
        pc = self._evaluate_program_counter(ctx.program_counter(i=1))

        instruction = Instruction(op_code, suffix, ra=ra, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitJeq_rri_instruction(self, ctx: assemblyParser.Jeq_rri_instructionContext) -> None:
        op_code = OpCode.SUB
        suffix = Suffix.ZRRCI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())
        condition = Condition.Z
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, ra=ra, rb=rb, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert rb == instruction.rb()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitJneq_rii_instruction(self, ctx: assemblyParser.Jneq_rii_instructionContext) -> None:
        op_code = OpCode.SUB
        suffix = Suffix.ZRICI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = self._evaluate_number(ctx.number())
        condition = Condition.NZ
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, ra=ra, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitJneq_rri_instruction(self, ctx: assemblyParser.Jneq_rri_instructionContext) -> None:
        op_code = OpCode.SUB
        suffix = Suffix.ZRRCI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())
        condition = Condition.NZ
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, ra=ra, rb=rb, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert rb == instruction.rb()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitJz_ri_instruction(self, ctx: assemblyParser.Jz_ri_instructionContext) -> None:
        op_code = OpCode.SUB
        suffix = Suffix.ZRICI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = 0
        condition = Condition.Z
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, ra=ra, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitJnz_ri_instruction(self, ctx: assemblyParser.Jnz_ri_instructionContext) -> None:
        op_code = OpCode.SUB
        suffix = Suffix.ZRICI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = 0
        condition = Condition.NZ
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, ra=ra, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitJltu_rii_instruction(self, ctx: assemblyParser.Jltu_rii_instructionContext) -> None:
        op_code = OpCode.SUB
        suffix = Suffix.ZRICI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = self._evaluate_number(ctx.number())
        condition = Condition.LTU
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, ra=ra, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitJltu_rri_instruction(self, ctx: assemblyParser.Jltu_rri_instructionContext) -> None:
        op_code = OpCode.SUB
        suffix = Suffix.ZRRCI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())
        condition = Condition.LTU
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, ra=ra, rb=rb, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert rb == instruction.rb()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitJgtu_rii_instruction(self, ctx: assemblyParser.Jgtu_rii_instructionContext) -> None:
        op_code = OpCode.SUB
        suffix = Suffix.ZRICI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = self._evaluate_number(ctx.number())
        condition = Condition.GTU
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, ra=ra, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitJgtu_rri_instruction(self, ctx: assemblyParser.Jgtu_rri_instructionContext) -> None:
        op_code = OpCode.SUB
        suffix = Suffix.ZRRCI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())
        condition = Condition.GTU
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, ra=ra, rb=rb, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert rb == instruction.rb()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitJleu_rii_instruction(self, ctx: assemblyParser.Jleu_rii_instructionContext) -> None:
        op_code = OpCode.SUB
        suffix = Suffix.ZRICI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = self._evaluate_number(ctx.number())
        condition = Condition.LEU
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, ra=ra, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitJleu_rri_instruction(self, ctx: assemblyParser.Jleu_rri_instructionContext) -> None:
        op_code = OpCode.SUB
        suffix = Suffix.ZRRCI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())
        condition = Condition.LEU
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, ra=ra, rb=rb, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert rb == instruction.rb()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitJgeu_rii_instruction(self, ctx: assemblyParser.Jgeu_rii_instructionContext) -> None:
        op_code = OpCode.SUB
        suffix = Suffix.ZRICI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = self._evaluate_number(ctx.number())
        condition = Condition.GEU
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, ra=ra, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitJgeu_rri_instruction(self, ctx: assemblyParser.Jgeu_rri_instructionContext) -> None:
        op_code = OpCode.SUB
        suffix = Suffix.ZRRCI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())
        condition = Condition.GEU
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, ra=ra, rb=rb, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert rb == instruction.rb()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitJlts_rii_instruction(self, ctx: assemblyParser.Jlts_rii_instructionContext) -> None:
        op_code = OpCode.SUB
        suffix = Suffix.ZRICI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = self._evaluate_number(ctx.number())
        condition = Condition.LTS
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, ra=ra, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitJlts_rri_instruction(self, ctx: assemblyParser.Jlts_rri_instructionContext) -> None:
        op_code = OpCode.SUB
        suffix = Suffix.ZRRCI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())
        condition = Condition.LTS
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, ra=ra, rb=rb, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert rb == instruction.rb()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitJgts_rii_instruction(self, ctx: assemblyParser.Jgts_rii_instructionContext) -> None:
        op_code = OpCode.SUB
        suffix = Suffix.ZRICI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = self._evaluate_number(ctx.number())
        condition = Condition.GTS
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, ra=ra, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitJgts_rri_instruction(self, ctx: assemblyParser.Jgts_rri_instructionContext) -> None:
        op_code = OpCode.SUB
        suffix = Suffix.ZRRCI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())
        condition = Condition.GTS
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, ra=ra, rb=rb, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert rb == instruction.rb()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitJles_rii_instruction(self, ctx: assemblyParser.Jles_rii_instructionContext) -> None:
        op_code = OpCode.SUB
        suffix = Suffix.ZRICI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = self._evaluate_number(ctx.number())
        condition = Condition.LES
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, ra=ra, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitJles_rri_instruction(self, ctx: assemblyParser.Jles_rri_instructionContext) -> None:
        op_code = OpCode.SUB
        suffix = Suffix.ZRRCI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())
        condition = Condition.LES
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, ra=ra, rb=rb, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert rb == instruction.rb()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitJges_rii_instruction(self, ctx: assemblyParser.Jges_rii_instructionContext) -> None:
        op_code = OpCode.SUB
        suffix = Suffix.ZRICI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = self._evaluate_number(ctx.number())
        condition = Condition.GES
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, ra=ra, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitJges_rri_instruction(self, ctx: assemblyParser.Jges_rri_instructionContext) -> None:
        op_code = OpCode.SUB
        suffix = Suffix.ZRRCI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())
        condition = Condition.GES
        pc = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, ra=ra, rb=rb, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert rb == instruction.rb()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitJump_ri_instruction(self, ctx: assemblyParser.Jump_ri_instructionContext) -> None:
        op_code = OpCode.CALL
        suffix = Suffix.ZRI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, ra=ra, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()

        self._executable.append_assembler_instruction(instruction)

    def exitJump_i_instruction(self, ctx: assemblyParser.Jump_i_instructionContext) -> None:
        op_code = OpCode.CALL
        suffix = Suffix.ZRI
        ra = SPRegister.ZERO
        imm = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, ra=ra, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()

        self._executable.append_assembler_instruction(instruction)

    def exitJump_r_instruction(self, ctx: assemblyParser.Jump_r_instructionContext) -> None:
        op_code = OpCode.CALL
        suffix = Suffix.ZRI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = 0

        instruction = Instruction(op_code, suffix, ra=ra, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()

        self._executable.append_assembler_instruction(instruction)

    def exitDiv_step_drdici_instruction(self, ctx: assemblyParser.Div_step_drdici_instructionContext) -> None:
        op_code = OpCode.DIV_STEP
        suffix = Suffix.DRDICI
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister(i=0)))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        db = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister(i=1)))
        imm = self._evaluate_number(ctx.number())
        condition = Condition.FALSE
        pc = 0

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra, db=db, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert db == instruction.db()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitMul_step_drdici_instruction(self, ctx: assemblyParser.Mul_step_drdici_instructionContext) -> None:
        op_code = OpCode.MUL_STEP
        suffix = Suffix.DRDICI
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister(i=0)))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        db = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister(i=1)))
        imm = self._evaluate_number(ctx.number())
        condition = Condition.FALSE
        pc = 0

        instruction = Instruction(op_code, suffix, dc=dc, ra=ra, db=db, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert db == instruction.db()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitBoot_rici_instruction(self, ctx: assemblyParser.Boot_rici_instructionContext) -> None:
        op_code = OpCode.BOOT
        suffix = Suffix.RICI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = self._evaluate_number(ctx.number())
        condition = Condition.FALSE
        pc = 0

        instruction = Instruction(op_code, suffix, ra=ra, imm=imm, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitResume_rici_instruction(self, ctx: assemblyParser.Resume_rici_instructionContext) -> None:
        op_code = OpCode.RESUME
        suffix = Suffix.RICI
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = self._evaluate_number(ctx.number())
        condition = Condition.FALSE
        pc = 0

        instruction = Instruction(op_code, suffix, ra=ra, imm=imm, condiiton=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitStop_ci_instruction(self, ctx: assemblyParser.Stop_ci_instructionContext) -> None:
        op_code = OpCode.STOP
        suffix = Suffix.CI
        condition = Condition.FALSE
        pc = 0

        instruction = Instruction(op_code, suffix, condiiton=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitCall_ri_instruction(self, ctx: assemblyParser.Call_ri_instructionContext) -> None:
        op_code = OpCode.CALL
        suffix = Suffix.RRI
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        ra = SPRegister.ZERO
        imm = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, rc=rc, ra=ra, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert rc == instruction.rc()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()

        self._executable.append_assembler_instruction(instruction)

    def exitCall_rr_instruction(self, ctx: assemblyParser.Call_rr_instructionContext) -> None:
        op_code = OpCode.CALL
        suffix = Suffix.RRI
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        imm = 0

        instruction = Instruction(op_code, suffix, rc=rc, ra=ra, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert rc == instruction.rc()
        assert ra == instruction.ra()
        assert imm == instruction.imm().value()

        self._executable.append_assembler_instruction(instruction)

    def exitBkp_instruction(self, ctx: assemblyParser.Bkp_instructionContext) -> None:
        op_code = OpCode.FAULT
        suffix = Suffix.I
        imm = 0

        instruction = Instruction(op_code, suffix, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert imm == instruction.imm().value()

        self._executable.append_assembler_instruction(instruction)

    def exitMovd_ddci_instruction(self, ctx: assemblyParser.Movd_ddci_instructionContext) -> None:
        op_code = OpCode.MOVD
        suffix = Suffix.DDCI
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister(i=0)))
        db = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister(i=1)))
        condition = Condition.FALSE
        pc = 0

        instruction = Instruction(op_code, suffix, dc=dc, db=db, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert db == instruction.db()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitSwapd_ddci_instruction(self, ctx: assemblyParser.Swapd_ddci_instructionContext) -> None:
        op_code = OpCode.SWAPD
        suffix = Suffix.DDCI
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister(i=0)))
        db = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister(i=1)))
        condition = Condition.FALSE
        pc = 0

        instruction = Instruction(op_code, suffix, dc=dc, db=db, condition=condition, pc=pc)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert dc == instruction.dc()
        assert db == instruction.db()
        assert condition == instruction.condition()
        assert pc == instruction.pc().value()

        self._executable.append_assembler_instruction(instruction)

    def exitTime_cfg_zr_instruction(self, ctx: assemblyParser.Time_cfg_zr_instructionContext) -> None:
        op_code = OpCode.TIME_CFG
        suffix = Suffix.ZR
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())

        instruction = Instruction(op_code, suffix, ra=ra)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert ra == instruction.ra()

        self._executable.append_assembler_instruction(instruction)

    def exitLbs_erri_instruction(self, ctx: assemblyParser.Lbs_erri_instructionContext) -> None:
        op_code = OpCode.LBS
        suffix = Suffix.ERRI
        endian = Endian.LITTLE
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        off = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, endian=endian, rc=rc, ra=ra, off=off)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert rc == instruction.rc()
        assert ra == instruction.ra()
        assert off == instruction.off().value()

        self._executable.append_assembler_instruction(instruction)

    def exitLbs_s_erri_instruction(self, ctx: assemblyParser.Lbs_s_erri_instructionContext) -> None:
        op_code = OpCode.LBS
        suffix = Suffix.S_ERRI
        endian = Endian.LITTLE
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        off = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, endian=endian, dc=dc, ra=ra, off=off)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert off == instruction.off().value()

        self._executable.append_assembler_instruction(instruction)

    def exitLbu_erri_instruction(self, ctx: assemblyParser.Lbu_erri_instructionContext) -> None:
        op_code = OpCode.LBU
        suffix = Suffix.ERRI
        endian = Endian.LITTLE
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        off = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, endian=endian, rc=rc, ra=ra, off=off)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert rc == instruction.rc()
        assert ra == instruction.ra()
        assert off == instruction.off().value()

        self._executable.append_assembler_instruction(instruction)

    def exitLbu_u_erri_instruction(self, ctx: assemblyParser.Lbu_u_erri_instructionContext) -> None:
        op_code = OpCode.LBS
        suffix = Suffix.U_ERRI
        endian = Endian.LITTLE
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        off = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, endian=endian, dc=dc, ra=ra, off=off)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert off == instruction.off().value()

        self._executable.append_assembler_instruction(instruction)

    def exitLd_edri_instruction(self, ctx: assemblyParser.Ld_edri_instructionContext) -> None:
        op_code = OpCode.LD
        suffix = Suffix.EDRI
        endian = Endian.LITTLE
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        off = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, endian=endian, dc=dc, ra=ra, off=off)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert off == instruction.off().value()

        self._executable.append_assembler_instruction(instruction)

    def exitLhs_erri_instruction(self, ctx: assemblyParser.Lhs_erri_instructionContext) -> None:
        op_code = OpCode.LHS
        suffix = Suffix.ERRI
        endian = Endian.LITTLE
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        off = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, endian=endian, rc=rc, ra=ra, off=off)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert rc == instruction.rc()
        assert ra == instruction.ra()
        assert off == instruction.off().value()

        self._executable.append_assembler_instruction(instruction)

    def exitLhs_s_erri_instruction(self, ctx: assemblyParser.Lhs_s_erri_instructionContext) -> None:
        op_code = OpCode.LHS
        suffix = Suffix.S_ERRI
        endian = Endian.LITTLE
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        off = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, endian=endian, dc=dc, ra=ra, off=off)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert off == instruction.off().value()

        self._executable.append_assembler_instruction(instruction)

    def exitLhu_erri_instruction(self, ctx: assemblyParser.Lhu_erri_instructionContext) -> None:
        op_code = OpCode.LHU
        suffix = Suffix.ERRI
        endian = Endian.LITTLE
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        off = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, endian=endian, rc=rc, ra=ra, off=off)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert rc == instruction.rc()
        assert ra == instruction.ra()
        assert off == instruction.off().value()

        self._executable.append_assembler_instruction(instruction)

    def exitLhu_u_erri_instruction(self, ctx: assemblyParser.Lhu_u_erri_instructionContext) -> None:
        op_code = OpCode.LHU
        suffix = Suffix.U_ERRI
        endian = Endian.LITTLE
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        off = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, endian=endian, dc=dc, ra=ra, off=off)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert off == instruction.off().value()

        self._executable.append_assembler_instruction(instruction)

    def exitLw_erri_instruction(self, ctx: assemblyParser.Lw_erri_instructionContext) -> None:
        op_code = OpCode.LW
        suffix = Suffix.ERRI
        endian = Endian.LITTLE
        rc = RegisterConverter.convert_to_gp_register(str(ctx.GPRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        off = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, endian=endian, rc=rc, ra=ra, off=off)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert rc == instruction.rc()
        assert ra == instruction.ra()
        assert off == instruction.off().value()

        self._executable.append_assembler_instruction(instruction)

    def exitLw_s_erri_instruction(self, ctx: assemblyParser.Lw_s_erri_instructionContext) -> None:
        op_code = OpCode.LW
        suffix = Suffix.S_ERRI
        endian = Endian.LITTLE
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        off = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, endian=endian, dc=dc, ra=ra, off=off)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert off == instruction.off().value()

        self._executable.append_assembler_instruction(instruction)

    def exitLw_u_erri_instruction(self, ctx: assemblyParser.Lw_u_erri_instructionContext) -> None:
        op_code = OpCode.LW
        suffix = Suffix.U_ERRI
        endian = Endian.LITTLE
        dc = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        off = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, endian=endian, dc=dc, ra=ra, off=off)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert dc == instruction.dc()
        assert ra == instruction.ra()
        assert off == instruction.off().value()

        self._executable.append_assembler_instruction(instruction)

    def exitSb_erii_instruction(self, ctx: assemblyParser.Sb_erii_instructionContext) -> None:
        op_code = OpCode.SB
        suffix = Suffix.ERII
        endian = Endian.LITTLE
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        off = self._evaluate_number(ctx.number())
        imm = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, endian=endian, ra=ra, off=off, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert ra == instruction.ra()
        assert off == instruction.off().value()
        assert (imm % (2 ** instruction.imm().width())) == (
            instruction.imm().value() % (2 ** instruction.imm().width())
        )

        self._executable.append_assembler_instruction(instruction)

    def exitSb_erir_instruction(self, ctx: assemblyParser.Sb_erir_instructionContext) -> None:
        op_code = OpCode.SB
        suffix = Suffix.ERIR
        endian = Endian.LITTLE
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        off = self._evaluate_program_counter(ctx.program_counter())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())

        instruction = Instruction(op_code, suffix, endian=endian, ra=ra, off=off, rb=rb)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert ra == instruction.ra()
        assert off == instruction.off().value()
        assert rb == instruction.rb()

        self._executable.append_assembler_instruction(instruction)

    def exitSb_id_rii_instruction(self, ctx: assemblyParser.Sb_id_rii_instructionContext) -> None:
        op_code = OpCode.SB_ID
        suffix = Suffix.ERII
        endian = Endian.LITTLE
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        off = self._evaluate_number(ctx.number(i=0))
        imm = self._evaluate_number(ctx.number(i=1))

        instruction = Instruction(op_code, suffix, endian=endian, ra=ra, off=off, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert ra == instruction.ra()
        assert off == instruction.off().value()
        assert imm == instruction.imm().value()

        self._executable.append_assembler_instruction(instruction)

    def exitSb_id_ri_instruction(self, ctx: assemblyParser.Sb_id_ri_instructionContext) -> None:
        op_code = OpCode.SB_ID
        suffix = Suffix.ERII
        endian = Endian.LITTLE
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        off = self._evaluate_number(ctx.number())
        imm = 0

        instruction = Instruction(op_code, suffix, endian=endian, ra=ra, off=off, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert ra == instruction.ra()
        assert off == instruction.off().value()
        assert imm == instruction.imm().value()

        self._executable.append_assembler_instruction(instruction)

    def exitSd_erii_instruction(self, ctx: assemblyParser.Sd_erii_instructionContext) -> None:
        op_code = OpCode.SD
        suffix = Suffix.ERII
        endian = Endian.LITTLE
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        off = self._evaluate_program_counter(ctx.program_counter(i=0))
        imm = self._evaluate_program_counter(ctx.program_counter(i=1))

        instruction = Instruction(op_code, suffix, endian=endian, ra=ra, off=off, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert ra == instruction.ra()
        assert off == instruction.off().value()
        assert imm == instruction.imm().value()

        self._executable.append_assembler_instruction(instruction)

    def exitSd_erid_instruction(self, ctx: assemblyParser.Sd_erid_instructionContext) -> None:
        op_code = OpCode.SD
        suffix = Suffix.ERID
        endian = Endian.LITTLE
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        off = self._evaluate_program_counter(ctx.program_counter())
        db = RegisterConverter.convert_to_pair_register(str(ctx.PairRegister()))

        instruction = Instruction(op_code, suffix, endian=endian, ra=ra, off=off, db=db)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert ra == instruction.ra()
        assert off == instruction.off().value()
        assert db == instruction.db()

        self._executable.append_assembler_instruction(instruction)

    def exitSd_id_rii_instruction(self, ctx: assemblyParser.Sd_id_rii_instructionContext) -> None:
        op_code = OpCode.SD_ID
        suffix = Suffix.ERII
        endian = Endian.LITTLE
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        off = self._evaluate_number(ctx.number(i=0))
        imm = self._evaluate_number(ctx.number(i=1))

        instruction = Instruction(op_code, suffix, endian=endian, ra=ra, off=off, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert ra == instruction.ra()
        assert off == instruction.off().value()
        assert imm == instruction.imm().value()

        self._executable.append_assembler_instruction(instruction)

    def exitSd_id_ri_instruction(self, ctx: assemblyParser.Sd_id_ri_instructionContext) -> None:
        op_code = OpCode.SD_ID
        suffix = Suffix.ERII
        endian = Endian.LITTLE
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        off = self._evaluate_number(ctx.number())
        imm = 0

        instruction = Instruction(op_code, suffix, endian=endian, ra=ra, off=off, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert ra == instruction.ra()
        assert off == instruction.off().value()
        assert imm == instruction.imm().value()

        self._executable.append_assembler_instruction(instruction)

    def exitSh_erii_instruction(self, ctx: assemblyParser.Sh_erii_instructionContext) -> None:
        op_code = OpCode.SH
        suffix = Suffix.ERII
        endian = Endian.LITTLE
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        off = self._evaluate_number(ctx.number())
        imm = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, endian=endian, ra=ra, off=off, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert ra == instruction.ra()
        assert off == instruction.off().value()
        assert imm == instruction.imm().value()

        self._executable.append_assembler_instruction(instruction)

    def exitSh_erir_instruction(self, ctx: assemblyParser.Sh_erir_instructionContext) -> None:
        op_code = OpCode.SH
        suffix = Suffix.ERIR
        endian = Endian.LITTLE
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        off = self._evaluate_program_counter(ctx.program_counter())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())

        instruction = Instruction(op_code, suffix, endian=endian, ra=ra, off=off, rb=rb)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert ra == instruction.ra()
        assert off == instruction.off().value()
        assert rb == instruction.rb()

        self._executable.append_assembler_instruction(instruction)

    def exitSh_id_rii_instruction(self, ctx: assemblyParser.Sh_id_rii_instructionContext) -> None:
        op_code = OpCode.SH_ID
        suffix = Suffix.ERII
        endian = Endian.LITTLE
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        off = self._evaluate_number(ctx.number(i=0))
        imm = self._evaluate_number(ctx.number(i=1))

        instruction = Instruction(op_code, suffix, endian=endian, ra=ra, off=off, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert ra == instruction.ra()
        assert off == instruction.off().value()
        assert imm == instruction.imm().value()

        self._executable.append_assembler_instruction(instruction)

    def exitSh_id_ri_instruction(self, ctx: assemblyParser.Sh_id_ri_instructionContext) -> None:
        op_code = OpCode.SH_ID
        suffix = Suffix.ERII
        endian = Endian.LITTLE
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        off = self._evaluate_number(ctx.number())
        imm = 0

        instruction = Instruction(op_code, suffix, endian=endian, ra=ra, off=off, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert ra == instruction.ra()
        assert off == instruction.off().value()
        assert imm == instruction.imm().value()

        self._executable.append_assembler_instruction(instruction)

    def exitSw_erii_instruction(self, ctx: assemblyParser.Sw_erii_instructionContext) -> None:
        op_code = OpCode.SW
        suffix = Suffix.ERII
        endian = Endian.LITTLE
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        off = self._evaluate_number(ctx.number())
        imm = self._evaluate_program_counter(ctx.program_counter())

        instruction = Instruction(op_code, suffix, endian=endian, ra=ra, off=off, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert ra == instruction.ra()
        assert off == instruction.off().value()
        assert imm == instruction.imm().value()

        self._executable.append_assembler_instruction(instruction)

    def exitSw_erir_instruction(self, ctx: assemblyParser.Sw_erir_instructionContext) -> None:
        op_code = OpCode.SW
        suffix = Suffix.ERIR
        endian = Endian.LITTLE
        ra = RegisterConverter.convert_to_source_register(ctx.src_register(i=0).getText())
        off = self._evaluate_program_counter(ctx.program_counter())
        rb = RegisterConverter.convert_to_source_register(ctx.src_register(i=1).getText())

        instruction = Instruction(op_code, suffix, endian=endian, ra=ra, off=off, rb=rb)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert ra == instruction.ra()
        assert off == instruction.off().value()
        assert rb == instruction.rb()

        self._executable.append_assembler_instruction(instruction)

    def exitSw_id_rii_instruction(self, ctx: assemblyParser.Sw_id_rii_instructionContext) -> None:
        op_code = OpCode.SW_ID
        suffix = Suffix.ERII
        endian = Endian.LITTLE
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        off = self._evaluate_number(ctx.number(i=0))
        imm = self._evaluate_number(ctx.number(i=1))

        instruction = Instruction(op_code, suffix, endian=endian, ra=ra, off=off, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert ra == instruction.ra()
        assert off == instruction.off().value()
        assert imm == instruction.imm().value()

        self._executable.append_assembler_instruction(instruction)

    def exitSw_id_ri_instruction(self, ctx: assemblyParser.Sw_id_ri_instructionContext) -> None:
        op_code = OpCode.SW_ID
        suffix = Suffix.ERII
        endian = Endian.LITTLE
        ra = RegisterConverter.convert_to_source_register(ctx.src_register().getText())
        off = self._evaluate_number(ctx.number())
        imm = 0

        instruction = Instruction(op_code, suffix, endian=endian, ra=ra, off=off, imm=imm)

        assert op_code == instruction.op_code()
        assert suffix == instruction.suffix()
        assert endian == instruction.endian()
        assert ra == instruction.ra()
        assert off == instruction.off().value()
        assert imm == instruction.imm().value()

        self._executable.append_assembler_instruction(instruction)

    def exitLabel(self, ctx: assemblyParser.LabelContext) -> None:
        label_name = str(ctx.Identifier())

        # TODO(bongjoon.hyun@gmail.com): __sys_used_mram_end will be defined in the linker script
        if label_name != "__sys_used_mram_end":
            self._executable.append_label(str(ctx.Identifier()))

    def _evaluate_program_counter(self, ctx: assemblyParser.Program_counterContext) -> int:
        child = ctx.children[0]
        if isinstance(child, assemblyParser.Primary_expressionContext):
            return self._evaluate_primary_expression(child)
        elif isinstance(child, assemblyParser.Add_expressionContext):
            return self._evaluate_add_expression(child)
        elif isinstance(child, assemblyParser.Sub_expressionContext):
            return self._evaluate_sub_expression(child)
        else:
            raise ValueError

    def _evaluate_primary_expression(self, ctx: assemblyParser.Primary_expressionContext) -> int:
        identifier = ctx.Identifier()
        if identifier is not None:
            return self._evaluate_identifier(str(identifier))

        child = ctx.children[0]
        if isinstance(child, assemblyParser.NumberContext):
            return InstructionAssigner._evaluate_number(child)
        elif isinstance(child, assemblyParser.Section_nameContext):
            return self._evaluate_section_name(child)
        else:
            raise ValueError

    def _evaluate_add_expression(self, ctx: assemblyParser.Add_expressionContext) -> int:
        primary_expression1 = self._evaluate_primary_expression(ctx.primary_expression(i=0))
        primary_expression2 = self._evaluate_primary_expression(ctx.primary_expression(i=1))

        if ctx.primary_expression(i=0).getText() != "NR_TASKLETS" and isinstance(
            ctx.primary_expression(i=1).children[0], assemblyParser.NumberContext
        ):
            return primary_expression1 + InstructionWord().size() * primary_expression2
        else:
            return primary_expression1 + primary_expression2

    def _evaluate_sub_expression(self, ctx: assemblyParser.Sub_expressionContext) -> int:
        primary_expression1 = self._evaluate_primary_expression(ctx.primary_expression(i=0))
        primary_expression2 = self._evaluate_primary_expression(ctx.primary_expression(i=1))

        if ctx.primary_expression(i=0).getText() != "NR_TASKLETS" and isinstance(
            ctx.primary_expression(i=1).children[0], assemblyParser.NumberContext
        ):
            return primary_expression1 - InstructionWord().size() * primary_expression2
        else:
            return primary_expression1 - primary_expression2

    @staticmethod
    def _evaluate_number(ctx: assemblyParser.NumberContext) -> int:
        if ctx.getText()[:2] == "0x":
            return int(ctx.getText()[2:], 16)
        else:
            return int(ctx.getText())

    def _evaluate_identifier(self, identifier: str) -> int:
        label = self._executable.label(identifier)
        if label is not None:
            assert label.address() is not None
            return label.address()

        if identifier in self._linker_script.symbol_names():
            assert self._linker_script.symbol(identifier).address() is not None
            return self._linker_script.symbol(identifier).address()

        if identifier in self._linker_script.constant_names():
            assert self._linker_script.constant(identifier) != 0
            return self._linker_script.constant(identifier)

        raise ValueError

    def _evaluate_section_name(self, ctx: assemblyParser.Section_nameContext) -> int:
        section_name = SectionNameConverter.convert_to_section_name(str(ctx.children[0])[1:])

        return self._executable.section(section_name, "").address()
