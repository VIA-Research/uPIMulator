import math
from typing import List, Tuple, Union

from abi.isa.instruction.condition import Condition
from abi.isa.instruction.endian import Endian
from abi.isa.instruction.instruction import Instruction
from abi.isa.instruction.op_code import OpCode
from abi.isa.instruction.suffix import Suffix
from abi.isa.register.gp_register import GPRegister
from abi.isa.register.pair_register import PairRegister
from abi.isa.register.sp_register import SPRegister
from abi.word.instruction_word import InstructionWord
from abi.word.representation import Representation
from encoder.byte import Byte
from util.config_loader import ConfigLoader


class InstructionEncoder:
    def __init__(self):
        pass

    @staticmethod
    def encode(instruction: Instruction) -> List[Byte]:
        instruction_word = InstructionWord()
        InstructionEncoder._encode_opcode(instruction, instruction_word)
        InstructionEncoder._encode_suffix(instruction, instruction_word)

        suffix = instruction.suffix()
        if suffix == Suffix.RICI:
            InstructionEncoder._encode_rici(instruction, instruction_word)
        elif suffix == Suffix.RRI:
            InstructionEncoder._encode_rri(instruction, instruction_word)
        elif suffix == Suffix.RRIC:
            InstructionEncoder._encode_rric(instruction, instruction_word)
        elif suffix == Suffix.RRICI:
            InstructionEncoder._encode_rrici(instruction, instruction_word)
        elif suffix == Suffix.RRIF:
            InstructionEncoder._encode_rrif(instruction, instruction_word)
        elif suffix == Suffix.RRR:
            InstructionEncoder._encode_rrr(instruction, instruction_word)
        elif suffix == Suffix.RRRC:
            InstructionEncoder._encode_rrrc(instruction, instruction_word)
        elif suffix == Suffix.RRRCI:
            InstructionEncoder._encode_rrrci(instruction, instruction_word)
        elif suffix == Suffix.ZRI:
            InstructionEncoder._encode_zri(instruction, instruction_word)
        elif suffix == Suffix.ZRIC:
            InstructionEncoder._encode_zric(instruction, instruction_word)
        elif suffix == Suffix.ZRICI:
            InstructionEncoder._encode_zrici(instruction, instruction_word)
        elif suffix == Suffix.ZRIF:
            InstructionEncoder._encode_zrif(instruction, instruction_word)
        elif suffix == Suffix.ZRR:
            InstructionEncoder._encode_zrr(instruction, instruction_word)
        elif suffix == Suffix.ZRRC:
            InstructionEncoder._encode_zrrc(instruction, instruction_word)
        elif suffix == Suffix.ZRRCI:
            InstructionEncoder._encode_zrrci(instruction, instruction_word)
        elif suffix == Suffix.S_RRI or suffix == Suffix.U_RRI:
            InstructionEncoder._encode_s_rri(instruction, instruction_word)
        elif suffix == Suffix.S_RRIC or suffix == Suffix.U_RRIC:
            InstructionEncoder._encode_s_rric(instruction, instruction_word)
        elif suffix == Suffix.S_RRICI or suffix == Suffix.U_RRICI:
            InstructionEncoder._encode_s_rrici(instruction, instruction_word)
        elif suffix == Suffix.S_RRIF or suffix == Suffix.U_RRIF:
            InstructionEncoder._encode_s_rrif(instruction, instruction_word)
        elif suffix == Suffix.S_RRR or suffix == Suffix.U_RRR:
            InstructionEncoder._encode_s_rrr(instruction, instruction_word)
        elif suffix == Suffix.S_RRRC or suffix == Suffix.U_RRRC:
            InstructionEncoder._encode_s_rrrc(instruction, instruction_word)
        elif suffix == Suffix.S_RRRCI or suffix == Suffix.U_RRRCI:
            InstructionEncoder._encode_s_rrrci(instruction, instruction_word)
        elif suffix == Suffix.RR:
            InstructionEncoder._encode_rr(instruction, instruction_word)
        elif suffix == Suffix.RRC:
            InstructionEncoder._encode_rrc(instruction, instruction_word)
        elif suffix == Suffix.RRCI:
            InstructionEncoder._encode_rrci(instruction, instruction_word)
        elif suffix == Suffix.ZR:
            InstructionEncoder._encode_zr(instruction, instruction_word)
        elif suffix == Suffix.ZRC:
            InstructionEncoder._encode_zrc(instruction, instruction_word)
        elif suffix == Suffix.ZRCI:
            InstructionEncoder._encode_zrci(instruction, instruction_word)
        elif suffix == Suffix.S_RR or suffix == Suffix.U_RR:
            InstructionEncoder._encode_s_rr(instruction, instruction_word)
        elif suffix == Suffix.S_RRC or suffix == Suffix.U_RRC:
            InstructionEncoder._encode_s_rrc(instruction, instruction_word)
        elif suffix == Suffix.S_RRCI or suffix == Suffix.U_RRCI:
            InstructionEncoder._encode_s_rrci(instruction, instruction_word)
        elif suffix == Suffix.DRDICI:
            InstructionEncoder._encode_drdici(instruction, instruction_word)
        elif suffix == Suffix.RRRI:
            InstructionEncoder._encode_rrri(instruction, instruction_word)
        elif suffix == Suffix.RRRICI:
            InstructionEncoder._encode_rrrici(instruction, instruction_word)
        elif suffix == Suffix.ZRRI:
            InstructionEncoder._encode_zrri(instruction, instruction_word)
        elif suffix == Suffix.ZRRICI:
            InstructionEncoder._encode_zrrici(instruction, instruction_word)
        elif suffix == Suffix.S_RRRI or suffix == Suffix.U_RRRI:
            InstructionEncoder._encode_s_rrri(instruction, instruction_word)
        elif suffix == Suffix.S_RRRICI or suffix == Suffix.U_RRRICI:
            InstructionEncoder._encode_s_rrrici(instruction, instruction_word)
        elif suffix == Suffix.RIR:
            InstructionEncoder._encode_rir(instruction, instruction_word)
        elif suffix == Suffix.RIRC:
            InstructionEncoder._encode_rirc(instruction, instruction_word)
        elif suffix == Suffix.RIRCI:
            InstructionEncoder._encode_rirci(instruction, instruction_word)
        elif suffix == Suffix.ZIR:
            InstructionEncoder._encode_zir(instruction, instruction_word)
        elif suffix == Suffix.ZIRC:
            InstructionEncoder._encode_zirc(instruction, instruction_word)
        elif suffix == Suffix.ZIRCI:
            InstructionEncoder._encode_zirci(instruction, instruction_word)
        elif suffix == Suffix.S_RIRC or suffix == Suffix.U_RIRC:
            InstructionEncoder._encode_s_rirc(instruction, instruction_word)
        elif suffix == Suffix.S_RIRCI or suffix == Suffix.U_RIRCI:
            InstructionEncoder._encode_s_rirci(instruction, instruction_word)
        elif suffix == Suffix.R:
            InstructionEncoder._encode_r(instruction, instruction_word)
        elif suffix == Suffix.RCI:
            InstructionEncoder._encode_rci(instruction, instruction_word)
        elif suffix == Suffix.Z:
            InstructionEncoder._encode_z(instruction, instruction_word)
        elif suffix == Suffix.ZCI:
            InstructionEncoder._encode_zci(instruction, instruction_word)
        elif suffix == Suffix.S_R or suffix == Suffix.U_R:
            InstructionEncoder._encode_s_r(instruction, instruction_word)
        elif suffix == Suffix.S_RCI or suffix == Suffix.U_RCI:
            InstructionEncoder._encode_s_rci(instruction, instruction_word)
        elif suffix == Suffix.CI:
            InstructionEncoder._encode_ci(instruction, instruction_word)
        elif suffix == Suffix.I:
            InstructionEncoder._encode_i(instruction, instruction_word)
        elif suffix == Suffix.DDCI:
            InstructionEncoder._encode_ddci(instruction, instruction_word)
        elif suffix == Suffix.ERRI:
            InstructionEncoder._encode_erri(instruction, instruction_word)
        elif suffix == Suffix.S_ERRI or suffix == Suffix.U_ERRI:
            InstructionEncoder._encode_s_erri(instruction, instruction_word)
        elif suffix == Suffix.EDRI:
            InstructionEncoder._encode_edri(instruction, instruction_word)
        elif suffix == Suffix.ERII:
            InstructionEncoder._encode_erii(instruction, instruction_word)
        elif suffix == Suffix.ERIR:
            InstructionEncoder._encode_erir(instruction, instruction_word)
        elif suffix == Suffix.ERID:
            InstructionEncoder._encode_erid(instruction, instruction_word)
        elif suffix == Suffix.DMA_RRI:
            InstructionEncoder._encode_dma_rri(instruction, instruction_word)
        else:
            raise ValueError

        return instruction_word.to_bytes()

    @staticmethod
    def decode(bytes_: List[Byte]) -> Instruction:
        instruction_word = InstructionWord()
        instruction_word.from_bytes(bytes_)

        op_code = InstructionEncoder._decode_op_code(instruction_word)
        suffix = InstructionEncoder._decode_suffix(instruction_word)

        if suffix == Suffix.RICI:
            return InstructionEncoder._decode_rici(op_code, suffix, instruction_word)
        elif suffix == Suffix.RRI:
            return InstructionEncoder._decode_rri(op_code, suffix, instruction_word)
        elif suffix == Suffix.RRIC:
            return InstructionEncoder._decode_rric(op_code, suffix, instruction_word)
        elif suffix == Suffix.RRICI:
            return InstructionEncoder._decode_rrici(op_code, suffix, instruction_word)
        elif suffix == Suffix.RRIF:
            return InstructionEncoder._decode_rrif(op_code, suffix, instruction_word)
        elif suffix == Suffix.RRR:
            return InstructionEncoder._decode_rrr(op_code, suffix, instruction_word)
        elif suffix == Suffix.RRRC:
            return InstructionEncoder._decode_rrrc(op_code, suffix, instruction_word)
        elif suffix == Suffix.RRRCI:
            return InstructionEncoder._decode_rrrci(op_code, suffix, instruction_word)
        elif suffix == Suffix.ZRI:
            return InstructionEncoder._decode_zri(op_code, suffix, instruction_word)
        elif suffix == Suffix.ZRIC:
            return InstructionEncoder._decode_zric(op_code, suffix, instruction_word)
        elif suffix == Suffix.ZRICI:
            return InstructionEncoder._decode_zrici(op_code, suffix, instruction_word)
        elif suffix == Suffix.ZRIF:
            return InstructionEncoder._decode_zrif(op_code, suffix, instruction_word)
        elif suffix == Suffix.ZRR:
            return InstructionEncoder._decode_zrr(op_code, suffix, instruction_word)
        elif suffix == Suffix.ZRRC:
            return InstructionEncoder._decode_zrrc(op_code, suffix, instruction_word)
        elif suffix == Suffix.ZRRCI:
            return InstructionEncoder._decode_zrrci(op_code, suffix, instruction_word)
        elif suffix == Suffix.S_RRI or suffix == Suffix.U_RRI:
            return InstructionEncoder._decode_s_rri(op_code, suffix, instruction_word)
        elif suffix == Suffix.S_RRIC or suffix == Suffix.U_RRIC:
            return InstructionEncoder._decode_s_rric(op_code, suffix, instruction_word)
        elif suffix == Suffix.S_RRICI or suffix == Suffix.U_RRICI:
            return InstructionEncoder._decode_s_rrici(op_code, suffix, instruction_word)
        elif suffix == Suffix.S_RRIF or suffix == Suffix.U_RRIF:
            return InstructionEncoder._decode_s_rrif(op_code, suffix, instruction_word)
        elif suffix == Suffix.S_RRR or suffix == Suffix.U_RRR:
            return InstructionEncoder._decode_s_rrr(op_code, suffix, instruction_word)
        elif suffix == Suffix.S_RRRC or suffix == Suffix.U_RRRC:
            return InstructionEncoder._decode_s_rrrc(op_code, suffix, instruction_word)
        elif suffix == Suffix.S_RRRCI or suffix == Suffix.U_RRRCI:
            return InstructionEncoder._decode_s_rrrci(op_code, suffix, instruction_word)
        elif suffix == Suffix.RR:
            return InstructionEncoder._decode_rr(op_code, suffix, instruction_word)
        elif suffix == Suffix.RRC:
            return InstructionEncoder._decode_rrc(op_code, suffix, instruction_word)
        elif suffix == Suffix.RRCI:
            return InstructionEncoder._decode_rrci(op_code, suffix, instruction_word)
        elif suffix == Suffix.ZR:
            return InstructionEncoder._decode_zr(op_code, suffix, instruction_word)
        elif suffix == Suffix.ZRC:
            return InstructionEncoder._decode_zrc(op_code, suffix, instruction_word)
        elif suffix == Suffix.ZRCI:
            return InstructionEncoder._decode_zrci(op_code, suffix, instruction_word)
        elif suffix == Suffix.S_RR or suffix == Suffix.U_RR:
            return InstructionEncoder._decode_s_rr(op_code, suffix, instruction_word)
        elif suffix == Suffix.S_RRC or suffix == Suffix.U_RRC:
            return InstructionEncoder._decode_s_rrc(op_code, suffix, instruction_word)
        elif suffix == Suffix.S_RRCI or suffix == Suffix.U_RRCI:
            return InstructionEncoder._decode_s_rrci(op_code, suffix, instruction_word)
        elif suffix == Suffix.DRDICI:
            return InstructionEncoder._decode_drdici(op_code, suffix, instruction_word)
        elif suffix == Suffix.RRRI:
            return InstructionEncoder._decode_rrri(op_code, suffix, instruction_word)
        elif suffix == Suffix.RRRICI:
            return InstructionEncoder._decode_rrrici(op_code, suffix, instruction_word)
        elif suffix == Suffix.ZRRI:
            return InstructionEncoder._decode_zrri(op_code, suffix, instruction_word)
        elif suffix == Suffix.ZRRICI:
            return InstructionEncoder._decode_zrrici(op_code, suffix, instruction_word)
        elif suffix == Suffix.S_RRRI or suffix == Suffix.U_RRRI:
            return InstructionEncoder._decode_s_rrri(op_code, suffix, instruction_word)
        elif suffix == Suffix.S_RRRICI or suffix == Suffix.U_RRRICI:
            return InstructionEncoder._decode_s_rrrici(op_code, suffix, instruction_word)
        elif suffix == Suffix.RIR:
            return InstructionEncoder._decode_rir(op_code, suffix, instruction_word)
        elif suffix == Suffix.RIRC:
            return InstructionEncoder._decode_rirc(op_code, suffix, instruction_word)
        elif suffix == Suffix.RIRCI:
            return InstructionEncoder._decode_rirci(op_code, suffix, instruction_word)
        elif suffix == Suffix.ZIR:
            return InstructionEncoder._decode_zir(op_code, suffix, instruction_word)
        elif suffix == Suffix.ZIRC:
            return InstructionEncoder._decode_zirc(op_code, suffix, instruction_word)
        elif suffix == Suffix.ZIRCI:
            return InstructionEncoder._decode_zirci(op_code, suffix, instruction_word)
        elif suffix == Suffix.S_RIRC or suffix == Suffix.U_RIRC:
            return InstructionEncoder._decode_s_rirc(op_code, suffix, instruction_word)
        elif suffix == Suffix.S_RIRCI or suffix == Suffix.U_RIRCI:
            return InstructionEncoder._decode_s_rirci(op_code, suffix, instruction_word)
        elif suffix == Suffix.R:
            return InstructionEncoder._decode_r(op_code, suffix, instruction_word)
        elif suffix == Suffix.RCI:
            return InstructionEncoder._decode_rci(op_code, suffix, instruction_word)
        elif suffix == Suffix.Z:
            return InstructionEncoder._decode_z(op_code, suffix, instruction_word)
        elif suffix == Suffix.ZCI:
            return InstructionEncoder._decode_zci(op_code, suffix, instruction_word)
        elif suffix == Suffix.S_R or suffix == Suffix.U_R:
            return InstructionEncoder._decode_s_r(op_code, suffix, instruction_word)
        elif suffix == Suffix.S_RCI or suffix == Suffix.U_RCI:
            return InstructionEncoder._decode_s_rci(op_code, suffix, instruction_word)
        elif suffix == Suffix.CI:
            return InstructionEncoder._decode_ci(op_code, suffix, instruction_word)
        elif suffix == Suffix.I:
            return InstructionEncoder._decode_i(op_code, suffix, instruction_word)
        elif suffix == Suffix.DDCI:
            return InstructionEncoder._decode_ddci(op_code, suffix, instruction_word)
        elif suffix == Suffix.ERRI:
            return InstructionEncoder._decode_erri(op_code, suffix, instruction_word)
        elif suffix == Suffix.S_ERRI or suffix == Suffix.U_ERRI:
            return InstructionEncoder._decode_s_erri(op_code, suffix, instruction_word)
        elif suffix == Suffix.EDRI:
            return InstructionEncoder._decode_edri(op_code, suffix, instruction_word)
        elif suffix == Suffix.ERII:
            return InstructionEncoder._decode_erii(op_code, suffix, instruction_word)
        elif suffix == Suffix.ERIR:
            return InstructionEncoder._decode_erir(op_code, suffix, instruction_word)
        elif suffix == Suffix.ERID:
            return InstructionEncoder._decode_erid(op_code, suffix, instruction_word)
        elif suffix == Suffix.DMA_RRI:
            return InstructionEncoder._decode_dma_rri(op_code, suffix, instruction_word)
        else:
            raise ValueError

    @staticmethod
    def _op_code_slice() -> Tuple[int, int]:
        return 0, math.ceil(math.log2(len(OpCode)))

    @staticmethod
    def _suffix_slice() -> Tuple[int, int]:
        _, begin = InstructionEncoder._op_code_slice()
        end = begin + math.ceil(math.log2(len(Suffix)))
        return begin, end

    @staticmethod
    def _register_width() -> int:
        return math.ceil(math.log2(ConfigLoader.num_gp_registers() + len(SPRegister)))

    @staticmethod
    def _condition_width() -> int:
        return math.ceil(math.log2(len(Condition)))

    @staticmethod
    def _pc_width() -> int:
        return ConfigLoader.iram_address_width()

    @staticmethod
    def _endian_width() -> int:
        return math.ceil(math.log2(len(Endian)))

    @staticmethod
    def _encode_opcode(instruction: Instruction, instruction_word: InstructionWord) -> None:
        instruction_word.set_bit_slice(*InstructionEncoder._op_code_slice(), instruction.op_code().value)

    @staticmethod
    def _decode_op_code(instruction_word: InstructionWord) -> OpCode:
        begin, end = InstructionEncoder._op_code_slice()
        return OpCode(instruction_word.bit_slice(Representation.UNSIGNED, begin, end))

    @staticmethod
    def _encode_suffix(instruction: Instruction, instruction_word: InstructionWord) -> None:
        instruction_word.set_bit_slice(*InstructionEncoder._suffix_slice(), instruction.suffix().value)

    @staticmethod
    def _decode_suffix(instruction_word: InstructionWord) -> Suffix:
        begin, end = InstructionEncoder._suffix_slice()
        return Suffix(instruction_word.bit_slice(Representation.UNSIGNED, begin, end))

    @staticmethod
    def _encode_register(
        instruction_word: InstructionWord, begin: int, end: int, register: Union[GPRegister, SPRegister, PairRegister],
    ) -> None:
        if isinstance(register, (GPRegister, PairRegister)):
            instruction_word.set_bit_slice(begin, end, register.index())
        elif isinstance(register, SPRegister):
            instruction_word.set_bit_slice(begin, end, ConfigLoader.num_gp_registers() + register.value)
        else:
            raise ValueError

    @staticmethod
    def _decode_gp_register(instruction_word: InstructionWord, begin: int, end: int) -> GPRegister:
        gp_register = instruction_word.bit_slice(Representation.UNSIGNED, begin, end)
        if gp_register < ConfigLoader.num_gp_registers():
            return GPRegister(gp_register)
        else:
            raise ValueError

    @staticmethod
    def _decode_source_register(instruction_word: InstructionWord, begin: int, end: int,) -> Instruction.SourceRegister:
        source_register = instruction_word.bit_slice(Representation.UNSIGNED, begin, end)
        if source_register < ConfigLoader.num_gp_registers():
            return GPRegister(source_register)
        elif source_register < ConfigLoader.num_gp_registers() + len(SPRegister):
            return SPRegister(source_register - ConfigLoader.num_gp_registers())
        else:
            raise ValueError

    @staticmethod
    def _decode_pair_register(instruction_word: InstructionWord, begin: int, end: int) -> PairRegister:
        pair_register = instruction_word.bit_slice(Representation.UNSIGNED, begin, end)
        if pair_register < ConfigLoader.num_gp_registers():
            return PairRegister(pair_register)
        else:
            raise ValueError

    @staticmethod
    def _encode_imm(instruction_word: InstructionWord, begin: int, end: int, value: int) -> None:
        instruction_word.set_bit_slice(begin, end, value)

    @staticmethod
    def _decode_imm(instruction_word: InstructionWord, begin: int, end: int, representation: Representation,) -> int:
        return instruction_word.bit_slice(representation, begin, end)

    @staticmethod
    def _encode_off(instruction_word: InstructionWord, begin: int, end: int, value: int) -> None:
        InstructionEncoder._encode_imm(instruction_word, begin, end, value)

    @staticmethod
    def _decode_off(instruction_word: InstructionWord, begin: int, end: int, representation: Representation,) -> int:
        return InstructionEncoder._decode_imm(instruction_word, begin, end, representation)

    @staticmethod
    def _encode_condition(instruction_word: InstructionWord, begin: int, end: int, condition: Condition) -> None:
        instruction_word.set_bit_slice(begin, end, condition.value)

    @staticmethod
    def _decode_condition(instruction_word: InstructionWord, begin: int, end: int) -> Condition:
        return Condition(instruction_word.bit_slice(Representation.UNSIGNED, begin, end))

    @staticmethod
    def _encode_pc(instruction_word: InstructionWord, begin: int, end: int, pc: int):
        InstructionEncoder._encode_imm(instruction_word, begin, end, pc)

    @staticmethod
    def _decode_pc(instruction_word: InstructionWord, begin: int, end: int) -> int:
        return InstructionEncoder._decode_imm(instruction_word, begin, end, Representation.UNSIGNED)

    @staticmethod
    def _encode_endian(instruction_word: InstructionWord, begin: int, end: int, endian: Endian) -> None:
        instruction_word.set_bit_slice(begin, end, endian.value)

    @staticmethod
    def _decode_endian(instruction_word: InstructionWord, begin: int, end: int) -> Endian:
        return Endian(instruction_word.bit_slice(Representation.UNSIGNED, begin, end))

    @staticmethod
    def _encode_rici(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RICIOpCodes
        assert instruction.suffix() == Suffix.RICI

        _, ra_begin = InstructionEncoder._suffix_slice()
        ra_end = ra_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        imm = instruction.imm()
        imm_begin, imm_end = ra_end, ra_end + imm.width()
        InstructionEncoder._encode_imm(instruction_word, imm_begin, imm_end, imm.value())

        condition = instruction.condition()
        condition_begin, condition_end = (
            imm_end,
            imm_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

        pc = instruction.pc()
        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        InstructionEncoder._encode_pc(instruction_word, pc_begin, pc_end, pc.value())

    @staticmethod
    def _decode_rici(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RICIOpCodes
        assert suffix == Suffix.RICI

        _, ra_begin = InstructionEncoder._suffix_slice()
        ra_end = ra_begin + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        imm_begin, imm_end = ra_end, ra_end + 16
        imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.SIGNED)

        condition_begin, condition_end = (
            imm_end,
            imm_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        pc = InstructionEncoder._decode_pc(instruction_word, pc_begin, pc_end)

        return Instruction(op_code, suffix, ra=ra, imm=imm, condition=condition, pc=pc)

    @staticmethod
    def _encode_rri(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRIOpCodes
        assert instruction.suffix() == Suffix.RRI

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rc_begin, rc_end, instruction.rc())

        ra_begin, ra_end = rc_end, rc_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        imm = instruction.imm()
        imm_begin, imm_end = ra_end, ra_end + imm.width()
        InstructionEncoder._encode_imm(instruction_word, imm_begin, imm_end, imm.value())

    @staticmethod
    def _decode_rri(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRIOpCodes
        assert suffix == Suffix.RRI

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        rc = InstructionEncoder._decode_gp_register(instruction_word, rc_begin, rc_end)

        ra_begin, ra_end = rc_end, rc_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        if op_code in Instruction.AddRRIOpCodes:
            imm_begin, imm_end = ra_end, ra_end + 32
            imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.UNSIGNED)
        elif op_code in Instruction.AsrRRIOpCodes:
            imm_begin, imm_end = ra_end, ra_end + 5
            imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.UNSIGNED)
        elif op_code in Instruction.CallRRIOpCodes:
            imm_begin, imm_end = ra_end, ra_end + 24
            imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.SIGNED)
        else:
            raise ValueError

        return Instruction(op_code, suffix, rc=rc, ra=ra, imm=imm)

    @staticmethod
    def _encode_rric(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRICOpCodes
        assert instruction.suffix() == Suffix.RRIC

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rc_begin, rc_end, instruction.rc())

        ra_begin, ra_end = rc_end, rc_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        imm = instruction.imm()
        imm_begin, imm_end = ra_end, ra_end + imm.width()
        InstructionEncoder._encode_imm(instruction_word, imm_begin, imm_end, imm.value())

        condition = instruction.condition()
        condition_begin, condition_end = (
            imm_end,
            imm_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

    @staticmethod
    def _decode_rric(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRICOpCodes
        assert suffix == Suffix.RRIC

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        rc = InstructionEncoder._decode_gp_register(instruction_word, rc_begin, rc_end)

        ra_begin, ra_end = rc_end, rc_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        if op_code in Instruction.AddRRICOpCodes or op_code in Instruction.SubRRICOpCodes:
            imm_begin, imm_end = ra_end, ra_end + 24
            imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.SIGNED)
        elif op_code in Instruction.AsrRRICOpCodes:
            imm_begin, imm_end = ra_end, ra_end + 5
            imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.UNSIGNED)
        else:
            raise ValueError

        condition_begin, condition_end = (
            imm_end,
            imm_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        return Instruction(op_code, suffix, rc=rc, ra=ra, imm=imm, condition=condition)

    @staticmethod
    def _encode_rrici(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRICIOpCodes
        assert instruction.suffix() == Suffix.RRICI

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rc_begin, rc_end, instruction.rc())

        ra_begin, ra_end = rc_end, rc_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        imm = instruction.imm()
        imm_begin, imm_end = ra_end, ra_end + imm.width()
        InstructionEncoder._encode_imm(instruction_word, imm_begin, imm_end, imm.value())

        condition = instruction.condition()
        condition_begin, condition_end = (
            imm_end,
            imm_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

        pc = instruction.pc()
        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        InstructionEncoder._encode_pc(instruction_word, pc_begin, pc_end, pc.value())

    @staticmethod
    def _decode_rrici(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRICOpCodes
        assert suffix == Suffix.RRICI

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        rc = InstructionEncoder._decode_gp_register(instruction_word, rc_begin, rc_end)

        ra_begin, ra_end = rc_end, rc_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        if (
            op_code in Instruction.AddRRICIOpCodes
            or op_code in Instruction.AndRRICIOpCodes
            or op_code in Instruction.SubRRICIOpCodes
        ):
            imm_begin, imm_end = ra_end, ra_end + 8
            imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.SIGNED)
        elif op_code in Instruction.AsrRRICIOpCodes:
            imm_begin, imm_end = ra_end, ra_end + 5
            imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.UNSIGNED)
        else:
            raise ValueError

        condition_begin, condition_end = (
            imm_end,
            imm_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        pc = InstructionEncoder._decode_pc(instruction_word, pc_begin, pc_end)

        return Instruction(op_code, suffix, rc=rc, ra=ra, imm=imm, condition=condition, pc=pc)

    @staticmethod
    def _encode_rrif(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRIFOpCodes
        assert instruction.suffix() == Suffix.RRIF

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rc_begin, rc_end, instruction.rc())

        ra_begin, ra_end = rc_end, rc_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        imm = instruction.imm()
        imm_begin, imm_end = ra_end, ra_end + imm.width()
        InstructionEncoder._encode_imm(instruction_word, imm_begin, imm_end, imm.value())

        condition = instruction.condition()
        condition_begin, condition_end = (
            imm_end,
            imm_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

    @staticmethod
    def _decode_rrif(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRIFOpCodes
        assert suffix == Suffix.RRIF

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        rc = InstructionEncoder._decode_gp_register(instruction_word, rc_begin, rc_end)

        ra_begin, ra_end = rc_end, rc_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        imm_begin, imm_end = ra_end, ra_end + 24
        imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.SIGNED)

        condition_begin, condition_end = (
            imm_end,
            imm_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        return Instruction(op_code, suffix, rc=rc, ra=ra, imm=imm, condition=condition)

    @staticmethod
    def _encode_rrr(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRROpCodes
        assert instruction.suffix() == Suffix.RRR

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rc_begin, rc_end, instruction.rc())

        ra_begin, ra_end = rc_end, rc_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        rb_begin, rb_end = ra_end, ra_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rb_begin, rb_end, instruction.rb())

    @staticmethod
    def _decode_rrr(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRROpCodes
        assert suffix == Suffix.RRR

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        rc = InstructionEncoder._decode_gp_register(instruction_word, rc_begin, rc_end)

        ra_begin, ra_end = rc_end, rc_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        rb_begin, rb_end = ra_end, ra_end + InstructionEncoder._register_width()
        rb = InstructionEncoder._decode_source_register(instruction_word, rb_begin, rb_end)

        return Instruction(op_code, suffix, rc=rc, ra=ra, rb=rb)

    @staticmethod
    def _encode_rrrc(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRRCOpCodes
        assert instruction.suffix() == Suffix.RRRC

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rc_begin, rc_end, instruction.rc())

        ra_begin, ra_end = rc_end, rc_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        rb_begin, rb_end = ra_end, ra_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rb_begin, rb_end, instruction.rb())

        condition = instruction.condition()
        condition_begin, condition_end = (
            rb_end,
            rb_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

    @staticmethod
    def _decode_rrrc(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRRCOpCodes
        assert suffix == Suffix.RRRC

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        rc = InstructionEncoder._decode_gp_register(instruction_word, rc_begin, rc_end)

        ra_begin, ra_end = rc_end, rc_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        rb_begin, rb_end = ra_end, ra_end + InstructionEncoder._register_width()
        rb = InstructionEncoder._decode_source_register(instruction_word, rb_begin, rb_end)

        condition_begin, condition_end = (
            rb_end,
            rb_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        return Instruction(op_code, suffix, rc=rc, ra=ra, rb=rb, condition=condition)

    @staticmethod
    def _encode_rrrci(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRRCIOpCodes
        assert instruction.suffix() == Suffix.RRRCI

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rc_begin, rc_end, instruction.rc())

        ra_begin, ra_end = rc_end, rc_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        rb_begin, rb_end = ra_end, ra_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rb_begin, rb_end, instruction.rb())

        condition = instruction.condition()
        condition_begin, condition_end = (
            rb_end,
            rb_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

        pc = instruction.pc()
        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        InstructionEncoder._encode_pc(instruction_word, pc_begin, pc_end, pc.value())

    @staticmethod
    def _decode_rrrci(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRRCIOpCodes
        assert suffix == Suffix.RRRCI

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        rc = InstructionEncoder._decode_gp_register(instruction_word, rc_begin, rc_end)

        ra_begin, ra_end = rc_end, rc_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        rb_begin, rb_end = ra_end, ra_end + InstructionEncoder._register_width()
        rb = InstructionEncoder._decode_source_register(instruction_word, rb_begin, rb_end)

        condition_begin, condition_end = (
            rb_end,
            rb_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        pc = InstructionEncoder._decode_pc(instruction_word, pc_begin, pc_end)

        return Instruction(op_code, suffix, rc=rc, ra=ra, rb=rb, condition=condition, pc=pc)

    @staticmethod
    def _encode_zri(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRIOpCodes
        assert instruction.suffix() == Suffix.ZRI

        _, ra_begin = InstructionEncoder._suffix_slice()
        ra_end = ra_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        imm = instruction.imm()
        imm_begin, imm_end = ra_end, ra_end + imm.width()
        InstructionEncoder._encode_imm(instruction_word, imm_begin, imm_end, imm.value())

    @staticmethod
    def _decode_zri(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRIOpCodes
        assert suffix == Suffix.ZRI

        _, ra_begin = InstructionEncoder._suffix_slice()
        ra_end = ra_begin + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        if op_code in Instruction.AddRRIOpCodes:
            imm_begin, imm_end = ra_end, ra_end + 32
            imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.UNSIGNED)
        elif op_code in Instruction.AsrRRIOpCodes:
            imm_begin, imm_end = ra_end, ra_end + 5
            imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.UNSIGNED)
        elif op_code in Instruction.CallRRIOpCodes:
            imm_begin, imm_end = ra_end, ra_end + 28
            imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.SIGNED)
        else:
            raise ValueError

        return Instruction(op_code, suffix, ra=ra, imm=imm)

    @staticmethod
    def _encode_zric(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRICOpCodes
        assert instruction.suffix() == Suffix.ZRIC

        _, ra_begin = InstructionEncoder._suffix_slice()
        ra_end = ra_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        imm = instruction.imm()
        imm_begin, imm_end = ra_end, ra_end + imm.width()
        InstructionEncoder._encode_imm(instruction_word, imm_begin, imm_end, imm.value())

        condition = instruction.condition()
        condition_begin, condition_end = (
            imm_end,
            imm_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

    @staticmethod
    def _decode_zric(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRICOpCodes
        assert suffix == Suffix.ZRIC

        _, ra_begin = InstructionEncoder._suffix_slice()
        ra_end = ra_begin + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        if op_code in Instruction.AddRRICOpCodes or op_code in Instruction.SubRRICOpCodes:
            imm_begin, imm_end = ra_end, ra_end + 27
            imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.SIGNED)
        elif op_code in Instruction.AsrRRICOpCodes:
            imm_begin, imm_end = ra_end, ra_end + 5
            imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.UNSIGNED)
        else:
            raise ValueError

        condition_begin, condition_end = (
            imm_end,
            imm_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        return Instruction(op_code, suffix, ra=ra, imm=imm, condition=condition)

    @staticmethod
    def _encode_zrici(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRICIOpCodes
        assert instruction.suffix() == Suffix.ZRICI

        _, ra_begin = InstructionEncoder._suffix_slice()
        ra_end = ra_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        imm = instruction.imm()
        imm_begin, imm_end = ra_end, ra_end + imm.width()
        InstructionEncoder._encode_imm(instruction_word, imm_begin, imm_end, imm.value())

        condition = instruction.condition()
        condition_begin, condition_end = (
            imm_end,
            imm_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

        pc = instruction.pc()
        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        InstructionEncoder._encode_pc(instruction_word, pc_begin, pc_end, pc.value())

    @staticmethod
    def _decode_zrici(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRICOpCodes
        assert suffix == Suffix.ZRICI

        _, ra_begin = InstructionEncoder._suffix_slice()
        ra_end = ra_begin + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        if (
            op_code in Instruction.AddRRICIOpCodes
            or op_code in Instruction.AndRRICIOpCodes
            or op_code in Instruction.SubRRICIOpCodes
        ):
            imm_begin, imm_end = ra_end, ra_end + 11
            imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.SIGNED)
        elif op_code in Instruction.AsrRRICIOpCodes:
            imm_begin, imm_end = ra_end, ra_end + 5
            imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.UNSIGNED)
        else:
            raise ValueError

        condition_begin, condition_end = (
            imm_end,
            imm_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        pc = InstructionEncoder._decode_pc(instruction_word, pc_begin, pc_end)

        return Instruction(op_code, suffix, ra=ra, imm=imm, condition=condition, pc=pc)

    @staticmethod
    def _encode_zrif(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRIFOpCodes
        assert instruction.suffix() == Suffix.ZRIF

        _, ra_begin = InstructionEncoder._suffix_slice()
        ra_end = ra_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        imm = instruction.imm()
        imm_begin, imm_end = ra_end, ra_end + imm.width()
        InstructionEncoder._encode_imm(instruction_word, imm_begin, imm_end, imm.value())

        condition = instruction.condition()
        condition_begin, condition_end = (
            imm_end,
            imm_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

    @staticmethod
    def _decode_zrif(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRIFOpCodes
        assert suffix == Suffix.ZRIF

        _, ra_begin = InstructionEncoder._suffix_slice()
        ra_end = ra_begin + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        imm_begin, imm_end = ra_end, ra_end + 27
        imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.SIGNED)

        condition_begin, condition_end = (
            imm_end,
            imm_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        return Instruction(op_code, suffix, ra=ra, imm=imm, condition=condition)

    @staticmethod
    def _encode_zrr(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRROpCodes
        assert instruction.suffix() == Suffix.ZRR

        _, ra_begin = InstructionEncoder._suffix_slice()
        ra_end = ra_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        rb_begin, rb_end = ra_end, ra_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rb_begin, rb_end, instruction.rb())

    @staticmethod
    def _decode_zrr(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRROpCodes
        assert suffix == Suffix.ZRR

        _, ra_begin = InstructionEncoder._suffix_slice()
        ra_end = ra_begin + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        rb_begin, rb_end = ra_end, ra_end + InstructionEncoder._register_width()
        rb = InstructionEncoder._decode_source_register(instruction_word, rb_begin, rb_end)

        return Instruction(op_code, suffix, ra=ra, rb=rb)

    @staticmethod
    def _encode_zrrc(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRRCOpCodes
        assert instruction.suffix() == Suffix.ZRRC

        _, ra_begin = InstructionEncoder._suffix_slice()
        ra_end = ra_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        rb_begin, rb_end = ra_end, ra_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rb_begin, rb_end, instruction.rb())

        condition = instruction.condition()
        condition_begin, condition_end = (
            rb_end,
            rb_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

    @staticmethod
    def _decode_zrrc(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRRCOpCodes
        assert suffix == Suffix.ZRRC

        _, ra_begin = InstructionEncoder._suffix_slice()
        ra_end = ra_begin + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        rb_begin, rb_end = ra_end, ra_end + InstructionEncoder._register_width()
        rb = InstructionEncoder._decode_source_register(instruction_word, rb_begin, rb_end)

        condition_begin, condition_end = (
            rb_end,
            rb_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        return Instruction(op_code, suffix, ra=ra, rb=rb, condition=condition)

    @staticmethod
    def _encode_zrrci(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRRCIOpCodes
        assert instruction.suffix() == Suffix.ZRRCI

        _, ra_begin = InstructionEncoder._suffix_slice()
        ra_end = ra_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        rb_begin, rb_end = ra_end, ra_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rb_begin, rb_end, instruction.rb())

        condition = instruction.condition()
        condition_begin, condition_end = (
            rb_end,
            rb_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

        pc = instruction.pc()
        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        InstructionEncoder._encode_pc(instruction_word, pc_begin, pc_end, pc.value())

    @staticmethod
    def _decode_zrrci(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRRCIOpCodes
        assert suffix == Suffix.ZRRCI

        _, ra_begin = InstructionEncoder._suffix_slice()
        ra_end = ra_begin + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        rb_begin, rb_end = ra_end, ra_end + InstructionEncoder._register_width()
        rb = InstructionEncoder._decode_source_register(instruction_word, rb_begin, rb_end)

        condition_begin, condition_end = (
            rb_end,
            rb_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        pc = InstructionEncoder._decode_pc(instruction_word, pc_begin, pc_end)

        return Instruction(op_code, suffix, ra=ra, rb=rb, condition=condition, pc=pc)

    @staticmethod
    def _encode_s_rri(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRIOpCodes
        assert instruction.suffix() == Suffix.S_RRI or instruction.suffix() == Suffix.U_RRI

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, dc_begin, dc_end, instruction.dc())

        ra_begin, ra_end = dc_end, dc_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        imm = instruction.imm()
        imm_begin, imm_end = ra_end, ra_end + imm.width()
        InstructionEncoder._encode_imm(instruction_word, imm_begin, imm_end, imm.value())

    @staticmethod
    def _decode_s_rri(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRIOpCodes
        assert suffix == Suffix.S_RRI or suffix == Suffix.U_RRI

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        dc = InstructionEncoder._decode_pair_register(instruction_word, dc_begin, dc_end)

        ra_begin, ra_end = dc_end, dc_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        if op_code in Instruction.AddRRIOpCodes:
            imm_begin, imm_end = ra_end, ra_end + 32
            imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.UNSIGNED)
        elif op_code in Instruction.AsrRRIOpCodes:
            imm_begin, imm_end = ra_end, ra_end + 5
            imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.UNSIGNED)
        elif op_code in Instruction.CallRRIOpCodes:
            imm_begin, imm_end = ra_end, ra_end + 24
            imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.SIGNED)
        else:
            raise ValueError

        return Instruction(op_code, suffix, dc=dc, ra=ra, imm=imm)

    @staticmethod
    def _encode_s_rric(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRICOpCodes
        assert instruction.suffix() == Suffix.S_RRIC or instruction.suffix() == Suffix.U_RRIC

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, dc_begin, dc_end, instruction.dc())

        ra_begin, ra_end = dc_end, dc_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        imm = instruction.imm()
        imm_begin, imm_end = ra_end, ra_end + imm.width()
        InstructionEncoder._encode_imm(instruction_word, imm_begin, imm_end, imm.value())

        condition = instruction.condition()
        condition_begin, condition_end = (
            imm_end,
            imm_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

    @staticmethod
    def _decode_s_rric(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRICOpCodes
        assert suffix == Suffix.S_RRIC or suffix == Suffix.U_RRIC

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        dc = InstructionEncoder._decode_pair_register(instruction_word, dc_begin, dc_end)

        ra_begin, ra_end = dc_end, dc_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        if op_code in Instruction.AddRRICOpCodes or op_code in Instruction.SubRRICOpCodes:
            imm_begin, imm_end = ra_end, ra_end + 24
            imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.SIGNED)
        elif op_code in Instruction.AsrRRICOpCodes:
            imm_begin, imm_end = ra_end, ra_end + 5
            imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.UNSIGNED)
        else:
            raise ValueError

        condition_begin, condition_end = (
            imm_end,
            imm_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        return Instruction(op_code, suffix, dc=dc, ra=ra, imm=imm, condition=condition)

    @staticmethod
    def _encode_s_rrici(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRICIOpCodes
        assert instruction.suffix() == Suffix.S_RRICI or instruction.suffix() == Suffix.U_RRICI

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, dc_begin, dc_end, instruction.dc())

        ra_begin, ra_end = dc_end, dc_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        imm = instruction.imm()
        imm_begin, imm_end = ra_end, ra_end + imm.width()
        InstructionEncoder._encode_imm(instruction_word, imm_begin, imm_end, imm.value())

        condition = instruction.condition()
        condition_begin, condition_end = (
            imm_end,
            imm_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

        pc = instruction.pc()
        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        InstructionEncoder._encode_pc(instruction_word, pc_begin, pc_end, pc.value())

    @staticmethod
    def _decode_s_rrici(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRICOpCodes
        assert suffix == Suffix.S_RRICI or suffix == Suffix.U_RRICI

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        dc = InstructionEncoder._decode_pair_register(instruction_word, dc_begin, dc_end)

        ra_begin, ra_end = dc_end, dc_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        if (
            op_code in Instruction.AddRRICIOpCodes
            or op_code in Instruction.AndRRICIOpCodes
            or op_code in Instruction.SubRRICIOpCodes
        ):
            imm_begin, imm_end = ra_end, ra_end + 8
            imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.SIGNED)
        elif op_code in Instruction.AsrRRICIOpCodes:
            imm_begin, imm_end = ra_end, ra_end + 5
            imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.UNSIGNED)
        else:
            raise ValueError

        condition_begin, condition_end = (
            imm_end,
            imm_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        pc = InstructionEncoder._decode_pc(instruction_word, pc_begin, pc_end)

        return Instruction(op_code, suffix, dc=dc, ra=ra, imm=imm, condition=condition, pc=pc)

    @staticmethod
    def _encode_s_rrif(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRIFOpCodes
        assert instruction.suffix() == Suffix.S_RRIF or instruction.suffix() == Suffix.U_RRIF

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, dc_begin, dc_end, instruction.dc())

        ra_begin, ra_end = dc_end, dc_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        imm = instruction.imm()
        imm_begin, imm_end = ra_end, ra_end + imm.width()
        InstructionEncoder._encode_imm(instruction_word, imm_begin, imm_end, imm.value())

        condition = instruction.condition()
        condition_begin, condition_end = (
            imm_end,
            imm_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

    @staticmethod
    def _decode_s_rrif(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRIFOpCodes
        assert suffix == Suffix.S_RRIF or suffix == Suffix.U_RRIF

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        dc = InstructionEncoder._decode_pair_register(instruction_word, dc_begin, dc_end)

        ra_begin, ra_end = dc_end, dc_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        imm_begin, imm_end = ra_end, ra_end + 24
        imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.SIGNED)

        condition_begin, condition_end = (
            imm_end,
            imm_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        return Instruction(op_code, suffix, dc=dc, ra=ra, imm=imm, condition=condition)

    @staticmethod
    def _encode_s_rrr(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRROpCodes
        assert instruction.suffix() == Suffix.S_RRR or instruction.suffix() == Suffix.U_RRR

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, dc_begin, dc_end, instruction.dc())

        ra_begin, ra_end = dc_end, dc_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        rb_begin, rb_end = ra_end, ra_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rb_begin, rb_end, instruction.rb())

    @staticmethod
    def _decode_s_rrr(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRROpCodes
        assert suffix == Suffix.S_RRR or suffix == Suffix.U_RRR

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        dc = InstructionEncoder._decode_pair_register(instruction_word, dc_begin, dc_end)

        ra_begin, ra_end = dc_end, dc_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        rb_begin, rb_end = ra_end, ra_end + InstructionEncoder._register_width()
        rb = InstructionEncoder._decode_source_register(instruction_word, rb_begin, rb_end)

        return Instruction(op_code, suffix, dc=dc, ra=ra, rb=rb)

    @staticmethod
    def _encode_s_rrrc(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRRCOpCodes
        assert instruction.suffix() == Suffix.S_RRRC or instruction.suffix() == Suffix.U_RRRC

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, dc_begin, dc_end, instruction.dc())

        ra_begin, ra_end = dc_end, dc_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        rb_begin, rb_end = ra_end, ra_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rb_begin, rb_end, instruction.rb())

        condition = instruction.condition()
        condition_begin, condition_end = (
            rb_end,
            rb_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

    @staticmethod
    def _decode_s_rrrc(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRRCOpCodes
        assert suffix == Suffix.S_RRRC or suffix == Suffix.U_RRRC

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        dc = InstructionEncoder._decode_pair_register(instruction_word, dc_begin, dc_end)

        ra_begin, ra_end = dc_end, dc_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        rb_begin, rb_end = ra_end, ra_end + InstructionEncoder._register_width()
        rb = InstructionEncoder._decode_source_register(instruction_word, rb_begin, rb_end)

        condition_begin, condition_end = (
            rb_end,
            rb_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        return Instruction(op_code, suffix, dc=dc, ra=ra, rb=rb, condition=condition)

    @staticmethod
    def _encode_s_rrrci(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRRCIOpCodes
        assert instruction.suffix() == Suffix.S_RRRCI or instruction.suffix() == Suffix.U_RRRCI

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, dc_begin, dc_end, instruction.dc())

        ra_begin, ra_end = dc_end, dc_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        rb_begin, rb_end = ra_end, ra_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rb_begin, rb_end, instruction.rb())

        condition = instruction.condition()
        condition_begin, condition_end = (
            rb_end,
            rb_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

        pc = instruction.pc()
        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        InstructionEncoder._encode_pc(instruction_word, pc_begin, pc_end, pc.value())

    @staticmethod
    def _decode_s_rrrci(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRRCIOpCodes
        assert suffix == Suffix.S_RRRCI or suffix == Suffix.U_RRRCI

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        dc = InstructionEncoder._decode_pair_register(instruction_word, dc_begin, dc_end)

        ra_begin, ra_end = dc_end, dc_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        rb_begin, rb_end = ra_end, ra_end + InstructionEncoder._register_width()
        rb = InstructionEncoder._decode_source_register(instruction_word, rb_begin, rb_end)

        condition_begin, condition_end = (
            rb_end,
            rb_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        pc = InstructionEncoder._decode_pc(instruction_word, pc_begin, pc_end)

        return Instruction(op_code, suffix, dc=dc, ra=ra, rb=rb, condition=condition, pc=pc)

    @staticmethod
    def _encode_rr(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RROpCodes
        assert instruction.suffix() == Suffix.RR

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rc_begin, rc_end, instruction.rc())

        ra_begin, ra_end = rc_end, rc_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

    @staticmethod
    def _decode_rr(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RROpCodes
        assert suffix == Suffix.RR

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        rc = InstructionEncoder._decode_gp_register(instruction_word, rc_begin, rc_end)

        ra_begin, ra_end = rc_end, rc_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        return Instruction(op_code, suffix, rc=rc, ra=ra)

    @staticmethod
    def _encode_rrc(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRCOpCodes
        assert instruction.suffix() == Suffix.RRC

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rc_begin, rc_end, instruction.rc())

        ra_begin, ra_end = rc_end, rc_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        condition = instruction.condition()
        condition_begin, condition_end = (
            ra_end,
            ra_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

    @staticmethod
    def _decode_rrc(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRCOpCodes
        assert suffix == Suffix.RRC

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        rc = InstructionEncoder._decode_gp_register(instruction_word, rc_begin, rc_end)

        ra_begin, ra_end = rc_end, rc_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        condition_begin, condition_end = (
            ra_end,
            ra_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        return Instruction(op_code, suffix, rc=rc, ra=ra, condition=condition)

    @staticmethod
    def _encode_rrci(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRCIOpCodes
        assert instruction.suffix() == Suffix.RRCI

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rc_begin, rc_end, instruction.rc())

        ra_begin, ra_end = rc_end, rc_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        condition = instruction.condition()
        condition_begin, condition_end = (
            ra_end,
            ra_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

        pc = instruction.pc()
        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        InstructionEncoder._encode_pc(instruction_word, pc_begin, pc_end, pc.value())

    @staticmethod
    def _decode_rrci(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRCIOpCodes
        assert suffix == Suffix.RRCI

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        rc = InstructionEncoder._decode_gp_register(instruction_word, rc_begin, rc_end)

        ra_begin, ra_end = rc_end, rc_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        condition_begin, condition_end = (
            ra_end,
            ra_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        pc = InstructionEncoder._decode_pc(instruction_word, pc_begin, pc_end)

        return Instruction(op_code, suffix, rc=rc, ra=ra, condition=condition, pc=pc)

    @staticmethod
    def _encode_zr(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RROpCodes
        assert instruction.suffix() == Suffix.ZR

        _, ra_begin = InstructionEncoder._suffix_slice()
        ra_end = ra_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

    @staticmethod
    def _decode_zr(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RROpCodes
        assert suffix == Suffix.ZR

        _, ra_begin = InstructionEncoder._suffix_slice()
        ra_end = ra_begin + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        return Instruction(op_code, suffix, ra=ra)

    @staticmethod
    def _encode_zrc(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRCOpCodes
        assert instruction.suffix() == Suffix.ZRC

        _, ra_begin = InstructionEncoder._suffix_slice()
        ra_end = ra_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        condition = instruction.condition()
        condition_begin, condition_end = (
            ra_end,
            ra_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

    @staticmethod
    def _decode_zrc(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRCOpCodes
        assert suffix == Suffix.ZRC

        _, ra_begin = InstructionEncoder._suffix_slice()
        ra_end = ra_begin + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        condition_begin, condition_end = (
            ra_end,
            ra_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        return Instruction(op_code, suffix, ra=ra, condition=condition)

    @staticmethod
    def _encode_zrci(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRCIOpCodes
        assert instruction.suffix() == Suffix.ZRCI

        _, ra_begin = InstructionEncoder._suffix_slice()
        ra_end = ra_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        condition = instruction.condition()
        condition_begin, condition_end = (
            ra_end,
            ra_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

        pc = instruction.pc()
        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        InstructionEncoder._encode_pc(instruction_word, pc_begin, pc_end, pc.value())

    @staticmethod
    def _decode_zrci(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRCIOpCodes
        assert suffix == Suffix.ZRCI

        _, ra_begin = InstructionEncoder._suffix_slice()
        ra_end = ra_begin + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        condition_begin, condition_end = (
            ra_end,
            ra_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        pc = InstructionEncoder._decode_pc(instruction_word, pc_begin, pc_end)

        return Instruction(op_code, suffix, ra=ra, condition=condition, pc=pc)

    @staticmethod
    def _encode_s_rr(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RROpCodes
        assert instruction.suffix() == Suffix.S_RR or instruction.suffix() == Suffix.U_RR

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, dc_begin, dc_end, instruction.dc())

        ra_begin, ra_end = dc_end, dc_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

    @staticmethod
    def _decode_s_rr(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RROpCodes
        assert suffix == Suffix.S_RR or suffix == Suffix.U_RR

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        dc = InstructionEncoder._decode_pair_register(instruction_word, dc_begin, dc_end)

        ra_begin, ra_end = dc_end, dc_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        return Instruction(op_code, suffix, dc=dc, ra=ra)

    @staticmethod
    def _encode_s_rrc(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRCOpCodes
        assert instruction.suffix() == Suffix.S_RRC or instruction.suffix() == Suffix.U_RRC

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, dc_begin, dc_end, instruction.dc())

        ra_begin, ra_end = dc_end, dc_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        condition = instruction.condition()
        condition_begin, condition_end = (
            ra_end,
            ra_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

    @staticmethod
    def _decode_s_rrc(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRCOpCodes
        assert suffix == Suffix.S_RRC or suffix == Suffix.U_RRC

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        dc = InstructionEncoder._decode_pair_register(instruction_word, dc_begin, dc_end)

        ra_begin, ra_end = dc_end, dc_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        condition_begin, condition_end = (
            ra_end,
            ra_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        return Instruction(op_code, suffix, dc=dc, ra=ra, condition=condition)

    @staticmethod
    def _encode_s_rrci(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRCIOpCodes
        assert instruction.suffix() == Suffix.S_RRCI or instruction.suffix() == Suffix.U_RRCI

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, dc_begin, dc_end, instruction.dc())

        ra_begin, ra_end = dc_end, dc_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        condition = instruction.condition()
        condition_begin, condition_end = (
            ra_end,
            ra_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

        pc = instruction.pc()
        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        InstructionEncoder._encode_pc(instruction_word, pc_begin, pc_end, pc.value())

    @staticmethod
    def _decode_s_rrci(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRCIOpCodes
        assert suffix == Suffix.S_RRCI or suffix == Suffix.U_RRCI

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        dc = InstructionEncoder._decode_pair_register(instruction_word, dc_begin, dc_end)

        ra_begin, ra_end = dc_end, dc_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        condition_begin, condition_end = (
            ra_end,
            ra_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        pc = InstructionEncoder._decode_pc(instruction_word, pc_begin, pc_end)

        return Instruction(op_code, suffix, dc=dc, ra=ra, condition=condition, pc=pc)

    @staticmethod
    def _encode_drdici(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.DRDICIOpCodes
        assert instruction.suffix() == Suffix.DRDICI

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, dc_begin, dc_end, instruction.dc())

        ra_begin, ra_end = dc_end, dc_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        db_begin, db_end = ra_end, ra_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, db_begin, db_end, instruction.db())

        imm = instruction.imm()
        imm_begin, imm_end = db_end, db_end + imm.width()
        InstructionEncoder._encode_imm(instruction_word, imm_begin, imm_end, imm.value())

        condition = instruction.condition()
        condition_begin, condition_end = (
            imm_end,
            imm_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

        pc = instruction.pc()
        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        InstructionEncoder._encode_pc(instruction_word, pc_begin, pc_end, pc.value())

    @staticmethod
    def _decode_drdici(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.DRDICIOpCodes
        assert suffix == Suffix.DRDICI

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        dc = InstructionEncoder._decode_pair_register(instruction_word, dc_begin, dc_end)

        ra_begin, ra_end = dc_end, dc_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        db_begin, db_end = ra_end, ra_end + InstructionEncoder._register_width()
        db = InstructionEncoder._decode_pair_register(instruction_word, db_begin, db_end)

        imm_begin, imm_end = db_end, db_end + 5
        imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.UNSIGNED)

        condition_begin, condition_end = (
            imm_end,
            imm_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        pc = InstructionEncoder._decode_pc(instruction_word, pc_begin, pc_end)

        return Instruction(op_code, suffix, dc=dc, ra=ra, db=db, imm=imm, condition=condition, pc=pc)

    @staticmethod
    def _encode_rrri(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRRIOpCodes
        assert instruction.suffix() == Suffix.RRRI

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rc_begin, rc_end, instruction.rc())

        ra_begin, ra_end = rc_end, rc_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        rb_begin, rb_end = ra_end, ra_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rb_begin, rb_end, instruction.rb())

        imm = instruction.imm()
        imm_begin, imm_end = rb_end, rb_end + imm.width()
        InstructionEncoder._encode_imm(instruction_word, imm_begin, imm_end, imm.value())

    @staticmethod
    def _decode_rrri(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRRIOpCodes
        assert suffix == Suffix.RRRI

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        rc = InstructionEncoder._decode_gp_register(instruction_word, rc_begin, rc_end)

        ra_begin, ra_end = rc_end, rc_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        rb_begin, rb_end = ra_end, ra_end + InstructionEncoder._register_width()
        rb = InstructionEncoder._decode_source_register(instruction_word, rb_begin, rb_end)

        imm_begin, imm_end = rb_end, rb_end + 5
        imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.UNSIGNED)

        return Instruction(op_code, suffix, rc=rc, ra=ra, rb=rb, imm=imm)

    @staticmethod
    def _encode_rrrici(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRRICIOpCodes
        assert instruction.suffix() == Suffix.RRRICI

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rc_begin, rc_end, instruction.rc())

        ra_begin, ra_end = rc_end, rc_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        rb_begin, rb_end = ra_end, ra_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rb_begin, rb_end, instruction.rb())

        imm = instruction.imm()
        imm_begin, imm_end = rb_end, rb_end + imm.width()
        InstructionEncoder._encode_imm(instruction_word, imm_begin, imm_end, imm.value())

        condition = instruction.condition()
        condition_begin, condition_end = (
            imm_end,
            imm_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

        pc = instruction.pc()
        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        InstructionEncoder._encode_pc(instruction_word, pc_begin, pc_end, pc.value())

    @staticmethod
    def _decode_rrrici(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRRICIOpCodes
        assert suffix == Suffix.RRRICI

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        rc = InstructionEncoder._decode_gp_register(instruction_word, rc_begin, rc_end)

        ra_begin, ra_end = rc_end, rc_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        rb_begin, rb_end = ra_end, ra_end + InstructionEncoder._register_width()
        rb = InstructionEncoder._decode_source_register(instruction_word, rb_begin, rb_end)

        imm_begin, imm_end = rb_end, rb_end + 5
        imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.UNSIGNED)

        condition_begin, condition_end = (
            imm_end,
            imm_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        pc = InstructionEncoder._decode_pc(instruction_word, pc_begin, pc_end)

        return Instruction(op_code, suffix, rc=rc, ra=ra, rb=rb, imm=imm, condition=condition, pc=pc)

    @staticmethod
    def _encode_zrri(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRRIOpCodes
        assert instruction.suffix() == Suffix.ZRRI

        _, ra_begin = InstructionEncoder._suffix_slice()
        ra_end = ra_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        rb_begin, rb_end = ra_end, ra_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rb_begin, rb_end, instruction.rb())

        imm = instruction.imm()
        imm_begin, imm_end = rb_end, rb_end + imm.width()
        InstructionEncoder._encode_imm(instruction_word, imm_begin, imm_end, imm.value())

    @staticmethod
    def _decode_zrri(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRRIOpCodes
        assert suffix == Suffix.ZRRI

        _, ra_begin = InstructionEncoder._suffix_slice()
        ra_end = ra_begin + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        rb_begin, rb_end = ra_end, ra_end + InstructionEncoder._register_width()
        rb = InstructionEncoder._decode_source_register(instruction_word, rb_begin, rb_end)

        imm_begin, imm_end = rb_end, rb_end + 5
        imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.UNSIGNED)

        return Instruction(op_code, suffix, ra=ra, rb=rb, imm=imm)

    @staticmethod
    def _encode_zrrici(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRRICIOpCodes
        assert instruction.suffix() == Suffix.ZRRICI

        _, ra_begin = InstructionEncoder._suffix_slice()
        ra_end = ra_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        rb_begin, rb_end = ra_end, ra_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rb_begin, rb_end, instruction.rb())

        imm = instruction.imm()
        imm_begin, imm_end = rb_end, rb_end + imm.width()
        InstructionEncoder._encode_imm(instruction_word, imm_begin, imm_end, imm.value())

        condition = instruction.condition()
        condition_begin, condition_end = (
            imm_end,
            imm_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

        pc = instruction.pc()
        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        InstructionEncoder._encode_pc(instruction_word, pc_begin, pc_end, pc.value())

    @staticmethod
    def _decode_zrrici(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRRICIOpCodes
        assert suffix == Suffix.ZRRICI

        _, ra_begin = InstructionEncoder._suffix_slice()
        ra_end = ra_begin + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        rb_begin, rb_end = ra_end, ra_end + InstructionEncoder._register_width()
        rb = InstructionEncoder._decode_source_register(instruction_word, rb_begin, rb_end)

        imm_begin, imm_end = rb_end, rb_end + 5
        imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.UNSIGNED)

        condition_begin, condition_end = (
            imm_end,
            imm_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        pc = InstructionEncoder._decode_pc(instruction_word, pc_begin, pc_end)

        return Instruction(op_code, suffix, ra=ra, rb=rb, imm=imm, condition=condition, pc=pc)

    @staticmethod
    def _encode_s_rrri(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRRIOpCodes
        assert instruction.suffix() == Suffix.S_RRRI or instruction.suffix() == Suffix.U_RRRI

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, dc_begin, dc_end, instruction.dc())

        ra_begin, ra_end = dc_end, dc_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        rb_begin, rb_end = ra_end, ra_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rb_begin, rb_end, instruction.rb())

        imm = instruction.imm()
        imm_begin, imm_end = rb_end, rb_end + imm.width()
        InstructionEncoder._encode_imm(instruction_word, imm_begin, imm_end, imm.value())

    @staticmethod
    def _decode_s_rrri(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRRIOpCodes
        assert suffix == Suffix.S_RRRI or suffix == Suffix.U_RRRI

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        dc = InstructionEncoder._decode_pair_register(instruction_word, dc_begin, dc_end)

        ra_begin, ra_end = dc_end, dc_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        rb_begin, rb_end = ra_end, ra_end + InstructionEncoder._register_width()
        rb = InstructionEncoder._decode_source_register(instruction_word, rb_begin, rb_end)

        imm_begin, imm_end = rb_end, rb_end + 5
        imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.UNSIGNED)

        return Instruction(op_code, suffix, dc=dc, ra=ra, rb=rb, imm=imm)

    @staticmethod
    def _encode_s_rrrici(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RRRICIOpCodes
        assert instruction.suffix() == Suffix.S_RRRICI or instruction.suffix() == Suffix.U_RRRICI

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, dc_begin, dc_end, instruction.dc())

        ra_begin, ra_end = dc_end, dc_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        rb_begin, rb_end = ra_end, ra_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rb_begin, rb_end, instruction.rb())

        imm = instruction.imm()
        imm_begin, imm_end = rb_end, rb_end + imm.width()
        InstructionEncoder._encode_imm(instruction_word, imm_begin, imm_end, imm.value())

        condition = instruction.condition()
        condition_begin, condition_end = (
            imm_end,
            imm_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

        pc = instruction.pc()
        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        InstructionEncoder._encode_pc(instruction_word, pc_begin, pc_end, pc.value())

    @staticmethod
    def _decode_s_rrrici(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RRRICIOpCodes
        assert suffix == Suffix.S_RRRICI or suffix == Suffix.U_RRRICI

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        dc = InstructionEncoder._decode_pair_register(instruction_word, dc_begin, dc_end)

        ra_begin, ra_end = dc_end, dc_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        rb_begin, rb_end = ra_end, ra_end + InstructionEncoder._register_width()
        rb = InstructionEncoder._decode_source_register(instruction_word, rb_begin, rb_end)

        imm_begin, imm_end = rb_end, rb_end + 5
        imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.UNSIGNED)

        condition_begin, condition_end = (
            imm_end,
            imm_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        pc = InstructionEncoder._decode_pc(instruction_word, pc_begin, pc_end)

        return Instruction(op_code, suffix, dc=dc, ra=ra, rb=rb, imm=imm, condition=condition, pc=pc)

    @staticmethod
    def _encode_rir(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RIROpCodes
        assert instruction.suffix() == Suffix.RIR

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rc_begin, rc_end, instruction.rc())

        imm = instruction.imm()
        imm_begin, imm_end = rc_end, rc_end + imm.width()
        InstructionEncoder._encode_imm(instruction_word, imm_begin, imm_end, imm.value())

        ra_begin, ra_end = imm_end, imm_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

    @staticmethod
    def _decode_rir(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RIROpCodes
        assert suffix == Suffix.RIR

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        rc = InstructionEncoder._decode_gp_register(instruction_word, rc_begin, rc_end)

        imm_begin, imm_end = rc_end, rc_end + 32
        imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.UNSIGNED)

        ra_begin, ra_end = imm_end, imm_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        return Instruction(op_code, suffix, rc=rc, imm=imm, ra=ra)

    @staticmethod
    def _encode_rirc(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RIRCOpCodes
        assert instruction.suffix() == Suffix.RIRC

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rc_begin, rc_end, instruction.rc())

        imm = instruction.imm()
        imm_begin, imm_end = rc_end, rc_end + imm.width()
        InstructionEncoder._encode_imm(instruction_word, imm_begin, imm_end, imm.value())

        ra_begin, ra_end = imm_end, imm_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        condition = instruction.condition()
        condition_begin, condition_end = (
            ra_end,
            ra_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

    @staticmethod
    def _decode_rirc(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RIRCOpCodes
        assert suffix == Suffix.RIRC

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        rc = InstructionEncoder._decode_gp_register(instruction_word, rc_begin, rc_end)

        imm_begin, imm_end = rc_end, rc_end + 24
        imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.SIGNED)

        ra_begin, ra_end = imm_end, imm_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        condition_begin, condition_end = (
            ra_end,
            ra_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        return Instruction(op_code, suffix, rc=rc, imm=imm, ra=ra, condition=condition)

    @staticmethod
    def _encode_rirci(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RIRCIOpCodes
        assert instruction.suffix() == Suffix.RIRCI

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rc_begin, rc_end, instruction.rc())

        imm = instruction.imm()
        imm_begin, imm_end = rc_end, rc_end + imm.width()
        InstructionEncoder._encode_imm(instruction_word, imm_begin, imm_end, imm.value())

        ra_begin, ra_end = imm_end, imm_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        condition = instruction.condition()
        condition_begin, condition_end = (
            ra_end,
            ra_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

        pc = instruction.pc()
        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        InstructionEncoder._encode_pc(instruction_word, pc_begin, pc_end, pc.value())

    @staticmethod
    def _decode_rirci(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RIRCIOpCodes
        assert suffix == Suffix.RIRCI

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        rc = InstructionEncoder._decode_gp_register(instruction_word, rc_begin, rc_end)

        imm_begin, imm_end = rc_end, rc_end + 8
        imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.SIGNED)

        ra_begin, ra_end = imm_end, imm_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        condition_begin, condition_end = (
            ra_end,
            ra_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        pc = InstructionEncoder._decode_pc(instruction_word, pc_begin, pc_end)

        return Instruction(op_code, suffix, rc=rc, imm=imm, ra=ra, condition=condition, pc=pc)

    @staticmethod
    def _encode_zir(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RIROpCodes
        assert instruction.suffix() == Suffix.ZIR

        imm = instruction.imm()
        _, imm_begin = InstructionEncoder._suffix_slice()
        imm_end = imm_begin + imm.width()
        InstructionEncoder._encode_imm(instruction_word, imm_begin, imm_end, imm.value())

        ra_begin, ra_end = imm_end, imm_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

    @staticmethod
    def _decode_zir(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RIROpCodes
        assert suffix == Suffix.ZIR

        _, imm_begin = InstructionEncoder._suffix_slice()
        imm_end = imm_begin + 32
        imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.UNSIGNED)

        ra_begin, ra_end = imm_end, imm_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        return Instruction(op_code, suffix, imm=imm, ra=ra)

    @staticmethod
    def _encode_zirc(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RIRCOpCodes
        assert instruction.suffix() == Suffix.ZIRC

        imm = instruction.imm()
        _, imm_begin = InstructionEncoder._suffix_slice()
        imm_end = imm_begin + imm.width()
        InstructionEncoder._encode_imm(instruction_word, imm_begin, imm_end, imm.value())

        ra_begin, ra_end = imm_end, imm_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        condition = instruction.condition()
        condition_begin, condition_end = (
            ra_end,
            ra_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

    @staticmethod
    def _decode_zirc(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RIRCOpCodes
        assert suffix == Suffix.ZIRC

        _, imm_begin = InstructionEncoder._suffix_slice()
        imm_end = imm_begin + 27
        imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.SIGNED)

        ra_begin, ra_end = imm_end, imm_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        condition_begin, condition_end = (
            ra_end,
            ra_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        return Instruction(op_code, suffix, imm=imm, ra=ra, condition=condition)

    @staticmethod
    def _encode_zirci(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RIRCIOpCodes
        assert instruction.suffix() == Suffix.ZIRCI

        imm = instruction.imm()
        _, imm_begin = InstructionEncoder._suffix_slice()
        imm_end = imm_begin + imm.width()
        InstructionEncoder._encode_imm(instruction_word, imm_begin, imm_end, imm.value())

        ra_begin, ra_end = imm_end, imm_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        condition = instruction.condition()
        condition_begin, condition_end = (
            ra_end,
            ra_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

        pc = instruction.pc()
        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        InstructionEncoder._encode_pc(instruction_word, pc_begin, pc_end, pc.value())

    @staticmethod
    def _decode_zirci(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RIRCIOpCodes
        assert suffix == Suffix.ZIRCI

        _, imm_begin = InstructionEncoder._suffix_slice()
        imm_end = imm_begin + 11
        imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.SIGNED)

        ra_begin, ra_end = imm_end, imm_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        condition_begin, condition_end = (
            ra_end,
            ra_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        pc = InstructionEncoder._decode_pc(instruction_word, pc_begin, pc_end)

        return Instruction(op_code, suffix, imm=imm, ra=ra, condition=condition, pc=pc)

    @staticmethod
    def _encode_s_rirc(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RIRCOpCodes
        assert instruction.suffix() == Suffix.S_RIRC or instruction.suffix() == Suffix.U_RIRC

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, dc_begin, dc_end, instruction.dc())

        imm = instruction.imm()
        imm_begin, imm_end = dc_end, dc_end + imm.width()
        InstructionEncoder._encode_imm(instruction_word, imm_begin, imm_end, imm.value())

        ra_begin, ra_end = imm_end, imm_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        condition = instruction.condition()
        condition_begin, condition_end = (
            ra_end,
            ra_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

    @staticmethod
    def _decode_s_rirc(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RIRCOpCodes
        assert suffix == Suffix.S_RIRC or suffix == Suffix.U_RIRC

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        dc = InstructionEncoder._decode_pair_register(instruction_word, dc_begin, dc_end)

        imm_begin, imm_end = dc_end, dc_end + 24
        imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.SIGNED)

        ra_begin, ra_end = imm_end, imm_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        condition_begin, condition_end = (
            ra_end,
            ra_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        return Instruction(op_code, suffix, dc=dc, imm=imm, ra=ra, condition=condition)

    @staticmethod
    def _encode_s_rirci(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RIRCIOpCodes
        assert instruction.suffix() == Suffix.S_RIRCI or instruction.suffix() == Suffix.U_RIRCI

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, dc_begin, dc_end, instruction.dc())

        imm = instruction.imm()
        imm_begin, imm_end = dc_end, dc_end + imm.width()
        InstructionEncoder._encode_imm(instruction_word, imm_begin, imm_end, imm.value())

        ra_begin, ra_end = imm_end, imm_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        condition = instruction.condition()
        condition_begin, condition_end = (
            ra_end,
            ra_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

        pc = instruction.pc()
        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        InstructionEncoder._encode_pc(instruction_word, pc_begin, pc_end, pc.value())

    @staticmethod
    def _decode_s_rirci(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RIRCIOpCodes
        assert suffix == Suffix.S_RIRCI or suffix == Suffix.U_RIRCI

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        dc = InstructionEncoder._decode_pair_register(instruction_word, dc_begin, dc_end)

        imm_begin, imm_end = dc_end, dc_end + 8
        imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.SIGNED)

        ra_begin, ra_end = imm_end, imm_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        condition_begin, condition_end = (
            ra_end,
            ra_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        pc = InstructionEncoder._decode_pc(instruction_word, pc_begin, pc_end)

        return Instruction(op_code, suffix, dc=dc, imm=imm, ra=ra, condition=condition, pc=pc)

    @staticmethod
    def _encode_r(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.ROpCodes
        assert instruction.suffix() == Suffix.R

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rc_begin, rc_end, instruction.rc())

    @staticmethod
    def _decode_r(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.ROpCodes
        assert suffix == Suffix.R

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        rc = InstructionEncoder._decode_gp_register(instruction_word, rc_begin, rc_end)

        return Instruction(op_code, suffix, rc=rc)

    @staticmethod
    def _encode_rci(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RCIOpCodes
        assert instruction.suffix() == Suffix.RCI

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rc_begin, rc_end, instruction.rc())

        condition = instruction.condition()
        condition_begin, condition_end = (
            rc_end,
            rc_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

        pc = instruction.pc()
        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        InstructionEncoder._encode_pc(instruction_word, pc_begin, pc_end, pc.value())

    @staticmethod
    def _decode_rci(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RCIOpCodes
        assert suffix == Suffix.RCI

        _, rc_begin = InstructionEncoder._suffix_slice()
        rc_end = rc_begin + InstructionEncoder._register_width()
        rc = InstructionEncoder._decode_gp_register(instruction_word, rc_begin, rc_end)

        condition_begin, condition_end = (
            rc_end,
            rc_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        pc = InstructionEncoder._decode_pc(instruction_word, pc_begin, pc_end)

        return Instruction(op_code, suffix, rc=rc, condition=condition, pc=pc)

    @staticmethod
    def _encode_z(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.ROpCodes or instruction.op_code() == OpCode.NOP
        assert instruction.suffix() == Suffix.Z

    @staticmethod
    def _decode_z(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.ROpCodes or op_code == OpCode.NOP
        assert suffix == Suffix.Z

        return Instruction(op_code, suffix)

    @staticmethod
    def _encode_zci(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RCIOpCodes
        assert instruction.suffix() == Suffix.ZCI

        condition = instruction.condition()
        _, condition_begin = InstructionEncoder._suffix_slice()
        condition_end = condition_begin + InstructionEncoder._condition_width()
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

        pc = instruction.pc()
        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        InstructionEncoder._encode_pc(instruction_word, pc_begin, pc_end, pc.value())

    @staticmethod
    def _decode_zci(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RCIOpCodes
        assert suffix == Suffix.ZCI

        _, condition_begin = InstructionEncoder._suffix_slice()
        condition_end = condition_begin + InstructionEncoder._condition_width()
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        pc = InstructionEncoder._decode_pc(instruction_word, pc_begin, pc_end)

        return Instruction(op_code, suffix, condition=condition, pc=pc)

    @staticmethod
    def _encode_s_r(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.ROpCodes
        assert instruction.suffix() == Suffix.S_R or instruction.suffix() == Suffix.U_R

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, dc_begin, dc_end, instruction.dc())

    @staticmethod
    def _decode_s_r(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.ROpCodes
        assert suffix == Suffix.S_R or suffix == Suffix.U_R

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        dc = InstructionEncoder._decode_pair_register(instruction_word, dc_begin, dc_end)

        return Instruction(op_code, suffix, dc=dc)

    @staticmethod
    def _encode_s_rci(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.RCIOpCodes
        assert instruction.suffix() == Suffix.S_RCI or instruction.suffix() == Suffix.U_RCI

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, dc_begin, dc_end, instruction.dc())

        condition = instruction.condition()
        condition_begin, condition_end = (
            dc_end,
            dc_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

        pc = instruction.pc()
        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        InstructionEncoder._encode_pc(instruction_word, pc_begin, pc_end, pc.value())

    @staticmethod
    def _decode_s_rci(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.RCIOpCodes
        assert suffix == Suffix.S_RCI or suffix == Suffix.U_RCI

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        dc = InstructionEncoder._decode_pair_register(instruction_word, dc_begin, dc_end)

        condition_begin, condition_end = (
            dc_end,
            dc_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        pc = InstructionEncoder._decode_pc(instruction_word, pc_begin, pc_end)

        return Instruction(op_code, suffix, dc=dc, condition=condition, pc=pc)

    @staticmethod
    def _encode_ci(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.CIOpCodes
        assert instruction.suffix() == Suffix.CI

        condition = instruction.condition()
        _, condition_begin = InstructionEncoder._suffix_slice()
        condition_end = condition_begin + InstructionEncoder._condition_width()
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

        pc = instruction.pc()
        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        InstructionEncoder._encode_pc(instruction_word, pc_begin, pc_end, pc.value())

    @staticmethod
    def _decode_ci(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.CIOpCodes
        assert suffix == Suffix.CI

        _, condition_begin = InstructionEncoder._suffix_slice()
        condition_end = condition_begin + InstructionEncoder._condition_width()
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        pc = InstructionEncoder._decode_pc(instruction_word, pc_begin, pc_end)

        return Instruction(op_code, suffix, condition=condition, pc=pc)

    @staticmethod
    def _encode_i(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.IOpCodes
        assert instruction.suffix() == Suffix.I

        imm = instruction.imm()
        _, imm_begin = InstructionEncoder._suffix_slice()
        imm_end = imm_begin + imm.width()
        InstructionEncoder._encode_imm(instruction_word, imm_begin, imm_end, imm.value())

    @staticmethod
    def _decode_i(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.IOpCodes
        assert suffix == Suffix.I

        _, imm_begin = InstructionEncoder._suffix_slice()
        imm_end = imm_begin + 24
        imm = InstructionEncoder._decode_imm(instruction_word, imm_begin, imm_end, Representation.SIGNED)

        return Instruction(op_code, suffix, imm=imm)

    @staticmethod
    def _encode_ddci(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.DDCIOpCodes
        assert instruction.suffix() == Suffix.DDCI

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, dc_begin, dc_end, instruction.dc())

        db_begin, db_end = dc_end, dc_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, db_begin, db_end, instruction.db())

        condition = instruction.condition()
        condition_begin, condition_end = (
            db_end,
            db_end + InstructionEncoder._condition_width(),
        )
        InstructionEncoder._encode_condition(instruction_word, condition_begin, condition_end, condition)

        pc = instruction.pc()
        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        InstructionEncoder._encode_pc(instruction_word, pc_begin, pc_end, pc.value())

    @staticmethod
    def _decode_ddci(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.DDCIOpCodes
        assert suffix == Suffix.DDCI

        _, dc_begin = InstructionEncoder._suffix_slice()
        dc_end = dc_begin + InstructionEncoder._register_width()
        dc = InstructionEncoder._decode_pair_register(instruction_word, dc_begin, dc_end)

        db_begin, db_end = dc_end, dc_end + InstructionEncoder._register_width()
        db = InstructionEncoder._decode_pair_register(instruction_word, db_begin, db_end)

        condition_begin, condition_end = (
            db_end,
            db_end + InstructionEncoder._condition_width(),
        )
        condition = InstructionEncoder._decode_condition(instruction_word, condition_begin, condition_end)

        pc_begin, pc_end = condition_end, condition_end + InstructionEncoder._pc_width()
        pc = InstructionEncoder._decode_pc(instruction_word, pc_begin, pc_end)

        return Instruction(op_code, suffix, dc=dc, db=db, condition=condition, pc=pc)

    @staticmethod
    def _encode_erri(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.ERRIOpCodes
        assert instruction.suffix() == Suffix.ERRI

        _, endian_begin = InstructionEncoder._suffix_slice()
        endian_end = endian_begin + InstructionEncoder._endian_width()
        InstructionEncoder._encode_endian(instruction_word, endian_begin, endian_end, instruction.endian())

        rc_begin, rc_end = endian_end, endian_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rc_begin, rc_end, instruction.rc())

        ra_begin, ra_end = rc_end, rc_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        off = instruction.off()
        off_begin, off_end = ra_end, ra_end + off.width()
        InstructionEncoder._encode_off(instruction_word, off_begin, off_end, off.value())

    @staticmethod
    def _decode_erri(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.ERRIOpCodes
        assert suffix == Suffix.ERRI

        _, endian_begin = InstructionEncoder._suffix_slice()
        endian_end = endian_begin + InstructionEncoder._endian_width()
        endian = InstructionEncoder._decode_endian(instruction_word, endian_begin, endian_end)

        rc_begin, rc_end = endian_end, endian_end + InstructionEncoder._register_width()
        rc = InstructionEncoder._decode_gp_register(instruction_word, rc_begin, rc_end)

        ra_begin, ra_end = rc_end, rc_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        off_begin, off_end = ra_end, ra_end + 24
        off = InstructionEncoder._decode_off(instruction_word, off_begin, off_end, Representation.SIGNED)

        return Instruction(op_code, suffix, endian=endian, rc=rc, ra=ra, off=off)

    @staticmethod
    def _encode_s_erri(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.ERRIOpCodes
        assert instruction.suffix() == Suffix.S_ERRI or instruction.suffix() == Suffix.U_ERRI

        _, endian_begin = InstructionEncoder._suffix_slice()
        endian_end = endian_begin + InstructionEncoder._endian_width()
        InstructionEncoder._encode_endian(instruction_word, endian_begin, endian_end, instruction.endian())

        dc_begin, dc_end = endian_end, endian_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, dc_begin, dc_end, instruction.dc())

        ra_begin, ra_end = dc_end, dc_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        off = instruction.off()
        off_begin, off_end = ra_end, ra_end + off.width()
        InstructionEncoder._encode_off(instruction_word, off_begin, off_end, off.value())

    @staticmethod
    def _decode_s_erri(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.ERRIOpCodes
        assert suffix == Suffix.S_ERRI or suffix == Suffix.U_ERRI

        _, endian_begin = InstructionEncoder._suffix_slice()
        endian_end = endian_begin + InstructionEncoder._endian_width()
        endian = InstructionEncoder._decode_endian(instruction_word, endian_begin, endian_end)

        dc_begin, dc_end = endian_end, endian_end + InstructionEncoder._register_width()
        dc = InstructionEncoder._decode_pair_register(instruction_word, dc_begin, dc_end)

        ra_begin, ra_end = dc_end, dc_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        off_begin, off_end = ra_end, ra_end + 24
        off = InstructionEncoder._decode_off(instruction_word, off_begin, off_end, Representation.SIGNED)

        return Instruction(op_code, suffix, endian=endian, dc=dc, ra=ra, off=off)

    @staticmethod
    def _encode_edri(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.EDRIOpCodes
        assert instruction.suffix() == Suffix.EDRI

        _, endian_begin = InstructionEncoder._suffix_slice()
        endian_end = endian_begin + InstructionEncoder._endian_width()
        InstructionEncoder._encode_endian(instruction_word, endian_begin, endian_end, instruction.endian())

        dc_begin, dc_end = endian_end, endian_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, dc_begin, dc_end, instruction.dc())

        ra_begin, ra_end = dc_end, dc_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        off = instruction.off()
        off_begin, off_end = ra_end, ra_end + off.width()
        InstructionEncoder._encode_off(instruction_word, off_begin, off_end, off.value())

    @staticmethod
    def _decode_edri(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.EDRIOpCodes
        assert suffix == Suffix.EDRI

        _, endian_begin = InstructionEncoder._suffix_slice()
        endian_end = endian_begin + InstructionEncoder._endian_width()
        endian = InstructionEncoder._decode_endian(instruction_word, endian_begin, endian_end)

        dc_begin, dc_end = endian_end, endian_end + InstructionEncoder._register_width()
        dc = InstructionEncoder._decode_pair_register(instruction_word, dc_begin, dc_end)

        ra_begin, ra_end = dc_end, dc_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        off_begin, off_end = ra_end, ra_end + 24
        off = InstructionEncoder._decode_off(instruction_word, off_begin, off_end, Representation.SIGNED)

        return Instruction(op_code, suffix, endian=endian, dc=dc, ra=ra, off=off)

    @staticmethod
    def _encode_erii(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.ERIIOpCodes
        assert instruction.suffix() == Suffix.ERII

        _, endian_begin = InstructionEncoder._suffix_slice()
        endian_end = endian_begin + InstructionEncoder._endian_width()
        InstructionEncoder._encode_endian(instruction_word, endian_begin, endian_end, instruction.endian())

        ra_begin, ra_end = endian_end, endian_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        off = instruction.off()
        off_begin, off_end = ra_end, ra_end + off.width()
        InstructionEncoder._encode_off(instruction_word, off_begin, off_end, off.value())

        imm = instruction.imm()
        imm_begin, imm_end = off_end, off_end + imm.width()
        InstructionEncoder._encode_off(instruction_word, imm_begin, imm_end, imm.value())

    @staticmethod
    def _decode_erii(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.ERIIOpCodes
        assert suffix == Suffix.ERII

        _, endian_begin = InstructionEncoder._suffix_slice()
        endian_end = endian_begin + InstructionEncoder._endian_width()
        endian = InstructionEncoder._decode_endian(instruction_word, endian_begin, endian_end)

        ra_begin, ra_end = endian_end, endian_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        # NOTE(bongjoon.hyun@gmail.com): original width is 12
        off_begin, off_end = ra_end, ra_end + 24
        off = InstructionEncoder._decode_off(instruction_word, off_begin, off_end, Representation.SIGNED)

        # NOTE(bongjoon.hyun@gmail.com): original width is 8
        imm_begin, imm_end = off_end, off_end + 16
        imm = InstructionEncoder._decode_off(instruction_word, imm_begin, imm_end, Representation.SIGNED)

        return Instruction(op_code, suffix, endian=endian, ra=ra, off=off, imm=imm)

    @staticmethod
    def _encode_erir(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.ERIROpCodes
        assert instruction.suffix() == Suffix.ERIR

        _, endian_begin = InstructionEncoder._suffix_slice()
        endian_end = endian_begin + InstructionEncoder._endian_width()
        InstructionEncoder._encode_endian(instruction_word, endian_begin, endian_end, instruction.endian())

        ra_begin, ra_end = endian_end, endian_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        off = instruction.off()
        off_begin, off_end = ra_end, ra_end + off.width()
        InstructionEncoder._encode_off(instruction_word, off_begin, off_end, off.value())

        rb_begin, rb_end = off_end, off_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rb_begin, rb_end, instruction.rb())

    @staticmethod
    def _decode_erir(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.ERIROpCodes
        assert suffix == Suffix.ERIR

        _, endian_begin = InstructionEncoder._suffix_slice()
        endian_end = endian_begin + InstructionEncoder._endian_width()
        endian = InstructionEncoder._decode_endian(instruction_word, endian_begin, endian_end)

        ra_begin, ra_end = endian_end, endian_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        off_begin, off_end = ra_end, ra_end + 24
        off = InstructionEncoder._decode_off(instruction_word, off_begin, off_end, Representation.SIGNED)

        rb_begin, rb_end = off_end, off_end + InstructionEncoder._register_width()
        rb = InstructionEncoder._decode_source_register(instruction_word, rb_begin, rb_end)

        return Instruction(op_code, suffix, endian=endian, ra=ra, off=off, rb=rb)

    @staticmethod
    def _encode_erid(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.ERIDOpCodes
        assert instruction.suffix() == Suffix.ERID

        _, endian_begin = InstructionEncoder._suffix_slice()
        endian_end = endian_begin + InstructionEncoder._endian_width()
        InstructionEncoder._encode_endian(instruction_word, endian_begin, endian_end, instruction.endian())

        ra_begin, ra_end = endian_end, endian_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        off = instruction.off()
        off_begin, off_end = ra_end, ra_end + off.width()
        InstructionEncoder._encode_off(instruction_word, off_begin, off_end, off.value())

        db_begin, db_end = off_end, off_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, db_begin, db_end, instruction.db())

    @staticmethod
    def _decode_erid(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.ERIDOpCodes
        assert suffix == Suffix.ERID

        _, endian_begin = InstructionEncoder._suffix_slice()
        endian_end = endian_begin + InstructionEncoder._endian_width()
        endian = InstructionEncoder._decode_endian(instruction_word, endian_begin, endian_end)

        ra_begin, ra_end = endian_end, endian_end + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        off_begin, off_end = ra_end, ra_end + 24
        off = InstructionEncoder._decode_off(instruction_word, off_begin, off_end, Representation.SIGNED)

        db_begin, db_end = off_end, off_end + InstructionEncoder._register_width()
        db = InstructionEncoder._decode_pair_register(instruction_word, db_begin, db_end)

        return Instruction(op_code, suffix, endian=endian, ra=ra, off=off, db=db)

    @staticmethod
    def _encode_dma_rri(instruction: Instruction, instruction_word: InstructionWord) -> None:
        assert instruction.op_code() in Instruction.DMARRIOpCodes
        assert instruction.suffix() == Suffix.DMA_RRI

        _, ra_begin = InstructionEncoder._suffix_slice()
        ra_end = ra_begin + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, ra_begin, ra_end, instruction.ra())

        rb_begin, rb_end = ra_end, ra_end + InstructionEncoder._register_width()
        InstructionEncoder._encode_register(instruction_word, rb_begin, rb_end, instruction.rb())

        imm = instruction.imm()
        imm_begin, imm_end = rb_end, rb_end + imm.width()
        InstructionEncoder._encode_off(instruction_word, imm_begin, imm_end, imm.value())

    @staticmethod
    def _decode_dma_rri(op_code: OpCode, suffix: Suffix, instruction_word: InstructionWord) -> Instruction:
        assert op_code in Instruction.DMARRIOpCodes
        assert suffix == Suffix.DMA_RRI

        _, ra_begin = InstructionEncoder._suffix_slice()
        ra_end = ra_begin + InstructionEncoder._register_width()
        ra = InstructionEncoder._decode_source_register(instruction_word, ra_begin, ra_end)

        rb_begin, rb_end = ra_end, ra_end + InstructionEncoder._register_width()
        rb = InstructionEncoder._decode_source_register(instruction_word, rb_begin, rb_end)

        imm_begin, imm_end = rb_end, rb_end + 8
        imm = InstructionEncoder._decode_off(instruction_word, imm_begin, imm_end, Representation.UNSIGNED)

        return Instruction(op_code, suffix, ra=ra, rb=rb, imm=imm)
