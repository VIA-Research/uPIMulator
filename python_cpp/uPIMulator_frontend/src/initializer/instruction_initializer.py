from typing import Set

from abi.isa.cc.acquire_cc import AcquireCC
from abi.isa.cc.add_nz_cc import AddNZCC
from abi.isa.cc.boot_cc import BootCC
from abi.isa.cc.count_nz_cc import CountNZCC
from abi.isa.cc.div_cc import DivCC
from abi.isa.cc.div_nz_cc import DivNZCC
from abi.isa.cc.ext_sub_set_cc import ExtSubSetCC
from abi.isa.cc.false_cc import FalseCC
from abi.isa.cc.imm_shift_nz_cc import ImmShiftNZCC
from abi.isa.cc.log_nz_cc import LogNZCC
from abi.isa.cc.log_set_cc import LogSetCC
from abi.isa.cc.mul_nz_cc import MulNZCC
from abi.isa.cc.release_cc import ReleaseCC
from abi.isa.cc.shift_nz_cc import ShiftNZCC
from abi.isa.cc.sub_nz_cc import SubNZCC
from abi.isa.cc.sub_set_cc import SubSetCC
from abi.isa.cc.true_cc import TrueCC
from abi.isa.cc.true_false_cc import TrueFalseCC
from abi.isa.instruction.condition import Condition
from abi.isa.instruction.endian import Endian
from abi.isa.instruction.instruction import Instruction
from abi.isa.instruction.op_code import OpCode
from abi.isa.instruction.suffix import Suffix
from abi.isa.register.gp_register import GPRegister
from abi.isa.register.pair_register import PairRegister
from abi.isa.register.sp_register import SPRegister
from abi.word.representation import Representation
from initializer.int_initializer import IntInitializer
from util.config_loader import ConfigLoader


class InstructionInitializer:
    def __init__(self):
        pass

    @staticmethod
    def instruction(op_code: OpCode, suffix: Suffix) -> Instruction:
        if suffix == Suffix.RICI:
            return InstructionInitializer._rici(op_code, suffix)
        elif suffix == Suffix.RRI:
            return InstructionInitializer._rri(op_code, suffix)
        elif suffix == Suffix.RRIC:
            return InstructionInitializer._rric(op_code, suffix)
        elif suffix == Suffix.RRICI:
            return InstructionInitializer._rrici(op_code, suffix)
        elif suffix == Suffix.RRIF:
            return InstructionInitializer._rrif(op_code, suffix)
        elif suffix == Suffix.RRR:
            return InstructionInitializer._rrr(op_code, suffix)
        elif suffix == Suffix.RRRC:
            return InstructionInitializer._rrrc(op_code, suffix)
        elif suffix == Suffix.RRRCI:
            return InstructionInitializer._rrrci(op_code, suffix)
        elif suffix == Suffix.ZRI:
            return InstructionInitializer._zri(op_code, suffix)
        elif suffix == Suffix.ZRIC:
            return InstructionInitializer._zric(op_code, suffix)
        elif suffix == Suffix.ZRICI:
            return InstructionInitializer._zrici(op_code, suffix)
        elif suffix == Suffix.ZRIF:
            return InstructionInitializer._zrif(op_code, suffix)
        elif suffix == Suffix.ZRR:
            return InstructionInitializer._zrr(op_code, suffix)
        elif suffix == Suffix.ZRRC:
            return InstructionInitializer._zrrc(op_code, suffix)
        elif suffix == Suffix.ZRRCI:
            return InstructionInitializer._zrrci(op_code, suffix)
        elif suffix == Suffix.S_RRI or suffix == Suffix.U_RRI:
            return InstructionInitializer._s_rri(op_code, suffix)
        elif suffix == Suffix.S_RRIC or suffix == Suffix.U_RRIC:
            return InstructionInitializer._s_rric(op_code, suffix)
        elif suffix == Suffix.S_RRICI or suffix == Suffix.U_RRICI:
            return InstructionInitializer._s_rrici(op_code, suffix)
        elif suffix == Suffix.S_RRIF or suffix == Suffix.U_RRIF:
            return InstructionInitializer._s_rrif(op_code, suffix)
        elif suffix == Suffix.S_RRR or suffix == Suffix.U_RRR:
            return InstructionInitializer._s_rrr(op_code, suffix)
        elif suffix == Suffix.S_RRRC or suffix == Suffix.U_RRRC:
            return InstructionInitializer._s_rrrc(op_code, suffix)
        elif suffix == Suffix.S_RRRCI or suffix == Suffix.U_RRRCI:
            return InstructionInitializer._s_rrrci(op_code, suffix)
        elif suffix == Suffix.RR:
            return InstructionInitializer._rr(op_code, suffix)
        elif suffix == Suffix.RRC:
            return InstructionInitializer._rrc(op_code, suffix)
        elif suffix == Suffix.RRCI:
            return InstructionInitializer._rrci(op_code, suffix)
        elif suffix == Suffix.ZR:
            return InstructionInitializer._zr(op_code, suffix)
        elif suffix == Suffix.ZRC:
            return InstructionInitializer._zrc(op_code, suffix)
        elif suffix == Suffix.ZRCI:
            return InstructionInitializer._zrci(op_code, suffix)
        elif suffix == Suffix.S_RR or suffix == Suffix.U_RR:
            return InstructionInitializer._s_rr(op_code, suffix)
        elif suffix == Suffix.S_RRC or suffix == Suffix.U_RRC:
            return InstructionInitializer._s_rrc(op_code, suffix)
        elif suffix == Suffix.S_RRCI or suffix == Suffix.U_RRCI:
            return InstructionInitializer._s_rrci(op_code, suffix)
        elif suffix == Suffix.DRDICI:
            return InstructionInitializer._drdici(op_code, suffix)
        elif suffix == Suffix.RRRI:
            return InstructionInitializer._rrri(op_code, suffix)
        elif suffix == Suffix.RRRICI:
            return InstructionInitializer._rrrici(op_code, suffix)
        elif suffix == Suffix.ZRRI:
            return InstructionInitializer._zrri(op_code, suffix)
        elif suffix == Suffix.ZRRICI:
            return InstructionInitializer._zrrici(op_code, suffix)
        elif suffix == Suffix.S_RRRI or suffix == Suffix.U_RRRI:
            return InstructionInitializer._s_rrri(op_code, suffix)
        elif suffix == Suffix.S_RRRICI or suffix == Suffix.U_RRRICI:
            return InstructionInitializer._s_rrrici(op_code, suffix)
        elif suffix == Suffix.RIR:
            return InstructionInitializer._rir(op_code, suffix)
        elif suffix == Suffix.RIRC:
            return InstructionInitializer._rirc(op_code, suffix)
        elif suffix == Suffix.RIRCI:
            return InstructionInitializer._rirci(op_code, suffix)
        elif suffix == Suffix.ZIR:
            return InstructionInitializer._zir(op_code, suffix)
        elif suffix == Suffix.ZIRC:
            return InstructionInitializer._zirc(op_code, suffix)
        elif suffix == Suffix.ZIRCI:
            return InstructionInitializer._zirci(op_code, suffix)
        elif suffix == Suffix.S_RIRC or suffix == Suffix.U_RIRC:
            return InstructionInitializer._s_rirc(op_code, suffix)
        elif suffix == Suffix.S_RIRCI or suffix == Suffix.U_RIRCI:
            return InstructionInitializer._s_rirci(op_code, suffix)
        elif suffix == Suffix.R:
            return InstructionInitializer._r(op_code, suffix)
        elif suffix == Suffix.RCI:
            return InstructionInitializer._rci(op_code, suffix)
        elif suffix == Suffix.Z:
            return InstructionInitializer._z(op_code, suffix)
        elif suffix == Suffix.ZCI:
            return InstructionInitializer._zci(op_code, suffix)
        elif suffix == Suffix.S_R or suffix == Suffix.U_R:
            return InstructionInitializer._s_r(op_code, suffix)
        elif suffix == Suffix.S_RCI or suffix == Suffix.U_RCI:
            return InstructionInitializer._s_rci(op_code, suffix)
        elif suffix == Suffix.CI:
            return InstructionInitializer._ci(op_code, suffix)
        elif suffix == Suffix.I:
            return InstructionInitializer._i(op_code, suffix)
        elif suffix == Suffix.DDCI:
            return InstructionInitializer._ddci(op_code, suffix)
        elif suffix == Suffix.ERRI:
            return InstructionInitializer._erri(op_code, suffix)
        elif suffix == Suffix.S_ERRI or suffix == Suffix.U_ERRI:
            return InstructionInitializer._s_erri(op_code, suffix)
        elif suffix == Suffix.EDRI:
            return InstructionInitializer._edri(op_code, suffix)
        elif suffix == Suffix.ERII:
            return InstructionInitializer._erii(op_code, suffix)
        elif suffix == Suffix.ERIR:
            return InstructionInitializer._erir(op_code, suffix)
        elif suffix == Suffix.ERID:
            return InstructionInitializer._erid(op_code, suffix)
        elif suffix == Suffix.DMA_RRI:
            return InstructionInitializer._dma_rri(op_code, suffix)
        else:
            raise ValueError

    @staticmethod
    def _gp_register() -> Instruction.DestinationRegister:
        return GPRegister(IntInitializer.value_by_range(0, ConfigLoader.num_gp_registers()))

    @staticmethod
    def _source_register() -> Instruction.SourceRegister:
        num_gp_registers = ConfigLoader.num_gp_registers()
        num_sp_registers = len(SPRegister)
        num_source_registers = num_gp_registers + num_sp_registers

        register_index = IntInitializer.value_by_range(0, num_source_registers)
        if register_index < ConfigLoader.num_gp_registers():
            return GPRegister(register_index)
        else:
            return SPRegister(register_index - num_gp_registers)

    @staticmethod
    def _pair_register() -> PairRegister:
        register_index = (IntInitializer.value_by_range(0, ConfigLoader.num_gp_registers()) // 2) * 2
        return PairRegister(register_index)

    @staticmethod
    def _imm(representation: Representation, width: int) -> int:
        return IntInitializer.value_by_width(representation, width)

    @staticmethod
    def _condition(conditions: Set[Condition]) -> Condition:
        return IntInitializer.value_by_list(list(conditions))

    @staticmethod
    def _pc() -> int:
        return InstructionInitializer._imm(Representation.UNSIGNED, ConfigLoader.iram_address_width())

    @staticmethod
    def _endian() -> Endian:
        return IntInitializer.value_by_list(list(Endian))

    @staticmethod
    def _rici(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RICIOpCodes
        assert suffix == Suffix.RICI

        ra = InstructionInitializer._source_register()
        imm = InstructionInitializer._imm(Representation.SIGNED, 16)

        if op_code in Instruction.AcquireRICIOpCodes:
            condition = InstructionInitializer._condition(AcquireCC.conditions())
        elif op_code in Instruction.ReleaseRICIOpCodes:
            condition = InstructionInitializer._condition(ReleaseCC.conditions())
        elif op_code in Instruction.BootRICIOpCodes:
            condition = InstructionInitializer._condition(BootCC.conditions())
        else:
            raise ValueError

        pc = InstructionInitializer._pc()

        return Instruction(op_code, suffix, ra=ra, imm=imm, condition=condition, pc=pc)

    @staticmethod
    def _rri(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRIOpCodes
        assert suffix == Suffix.RRI

        rc = InstructionInitializer._gp_register()
        ra = InstructionInitializer._source_register()

        if op_code in Instruction.AddRRIOpCodes:
            imm = InstructionInitializer._imm(Representation.UNSIGNED, 32)
        elif op_code in Instruction.AsrRRIOpCodes:
            imm = InstructionInitializer._imm(Representation.UNSIGNED, 5)
        elif op_code in Instruction.CallRRIOpCodes:
            imm = InstructionInitializer._imm(Representation.SIGNED, 24)
        else:
            raise ValueError

        return Instruction(op_code, suffix, rc=rc, ra=ra, imm=imm)

    @staticmethod
    def _rric(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRICOpCodes
        assert suffix == Suffix.RRIC

        rc = InstructionInitializer._gp_register()
        ra = InstructionInitializer._source_register()

        if op_code in Instruction.AddRRICOpCodes or op_code in Instruction.SubRRICOpCodes:
            imm = InstructionInitializer._imm(Representation.SIGNED, 24)
        elif op_code in Instruction.AsrRRICOpCodes:
            imm = InstructionInitializer._imm(Representation.UNSIGNED, 5)
        else:
            raise ValueError

        if op_code in Instruction.AddRRICOpCodes or op_code in Instruction.AsrRRICOpCodes:
            condition = InstructionInitializer._condition(LogSetCC.conditions())
        elif op_code in Instruction.SubRRICOpCodes:
            condition = InstructionInitializer._condition(ExtSubSetCC.conditions())
        else:
            raise ValueError

        return Instruction(op_code, suffix, rc=rc, ra=ra, imm=imm, condition=condition)

    @staticmethod
    def _rrici(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRICIOpCodes
        assert suffix == Suffix.RRICI

        rc = InstructionInitializer._gp_register()
        ra = InstructionInitializer._source_register()

        if (
            op_code in Instruction.AddRRICIOpCodes
            or op_code in Instruction.AndRRICIOpCodes
            or op_code in Instruction.SubRRICIOpCodes
        ):
            imm = InstructionInitializer._imm(Representation.SIGNED, 8)
        elif op_code in Instruction.AsrRRICIOpCodes:
            imm = InstructionInitializer._imm(Representation.UNSIGNED, 5)
        else:
            raise ValueError

        if op_code in Instruction.AddRRICIOpCodes:
            condition = InstructionInitializer._condition(AddNZCC.conditions())
        elif op_code in Instruction.AndRRICIOpCodes:
            condition = InstructionInitializer._condition(LogNZCC.conditions())
        elif op_code in Instruction.AsrRRICIOpCodes:
            condition = InstructionInitializer._condition(ImmShiftNZCC.conditions())
        elif op_code in Instruction.SubRRICIOpCodes:
            condition = InstructionInitializer._condition(SubNZCC.conditions())
        else:
            raise ValueError

        pc = InstructionInitializer._pc()

        return Instruction(op_code, suffix, rc=rc, ra=ra, imm=imm, condition=condition, pc=pc)

    @staticmethod
    def _rrif(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRIFOpCodes
        assert suffix == Suffix.RRIF

        rc = InstructionInitializer._gp_register()
        ra = InstructionInitializer._source_register()
        imm = InstructionInitializer._imm(Representation.SIGNED, 24)
        condition = InstructionInitializer._condition(FalseCC.conditions())

        return Instruction(op_code, suffix, rc=rc, ra=ra, imm=imm, condition=condition)

    @staticmethod
    def _rrr(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRROpCodes
        assert suffix == Suffix.RRR

        rc = InstructionInitializer._gp_register()
        ra = InstructionInitializer._source_register()
        rb = InstructionInitializer._source_register()

        return Instruction(op_code, suffix, rc=rc, ra=ra, rb=rb)

    @staticmethod
    def _rrrc(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRRCOpCodes
        assert suffix == Suffix.RRRC

        rc = InstructionInitializer._gp_register()
        ra = InstructionInitializer._source_register()
        rb = InstructionInitializer._source_register()

        if op_code in Instruction.AddRRRCOpCodes:
            condition = InstructionInitializer._condition(LogSetCC.conditions())
        elif op_code in Instruction.RsubRRRCOpCodes:
            condition = InstructionInitializer._condition(SubSetCC.conditions())
        elif op_code in Instruction.SubRRRCOpCodes:
            condition = InstructionInitializer._condition(ExtSubSetCC.conditions())
        else:
            raise ValueError

        return Instruction(op_code, suffix, rc=rc, ra=ra, rb=rb, condition=condition)

    @staticmethod
    def _rrrci(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRRCIOpCodes
        assert suffix == Suffix.RRRCI

        rc = InstructionInitializer._gp_register()
        ra = InstructionInitializer._source_register()
        rb = InstructionInitializer._source_register()

        if op_code in Instruction.AddRRRCIOpCodes:
            condition = InstructionInitializer._condition(AddNZCC.conditions())
        elif op_code in Instruction.AndRRRCIOpCodes:
            condition = InstructionInitializer._condition(LogNZCC.conditions())
        elif op_code in Instruction.AsrRRRCIOpCodes:
            condition = InstructionInitializer._condition(ShiftNZCC.conditions())
        elif op_code in Instruction.MulRRRCIOpCodes:
            condition = InstructionInitializer._condition(MulNZCC.conditions())
        elif op_code in Instruction.RsubRRRCIOpCodes:
            condition = InstructionInitializer._condition(SubNZCC.conditions())
        else:
            raise ValueError

        pc = InstructionInitializer._pc()

        return Instruction(op_code, suffix, rc=rc, ra=ra, rb=rb, condition=condition, pc=pc)

    @staticmethod
    def _zri(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRIOpCodes
        assert suffix == Suffix.ZRI

        ra = InstructionInitializer._source_register()

        if op_code in Instruction.AddRRIOpCodes:
            imm = InstructionInitializer._imm(Representation.UNSIGNED, 32)
        elif op_code in Instruction.AsrRRIOpCodes:
            imm = InstructionInitializer._imm(Representation.UNSIGNED, 5)
        elif op_code in Instruction.CallRRIOpCodes:
            imm = InstructionInitializer._imm(Representation.SIGNED, 28)
        else:
            raise ValueError

        return Instruction(op_code, suffix, ra=ra, imm=imm)

    @staticmethod
    def _zric(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRICOpCodes
        assert suffix == Suffix.ZRIC

        ra = InstructionInitializer._source_register()

        if op_code in Instruction.AddRRICOpCodes or op_code in Instruction.SubRRICOpCodes:
            imm = InstructionInitializer._imm(Representation.SIGNED, 27)
        elif op_code in Instruction.AsrRRICOpCodes:
            imm = InstructionInitializer._imm(Representation.UNSIGNED, 5)
        else:
            raise ValueError

        if op_code in Instruction.AddRRICOpCodes or op_code in Instruction.AsrRRICOpCodes:
            condition = InstructionInitializer._condition(LogSetCC.conditions())
        elif op_code in Instruction.SubRRICOpCodes:
            condition = InstructionInitializer._condition(ExtSubSetCC.conditions())
        else:
            raise ValueError

        return Instruction(op_code, suffix, ra=ra, imm=imm, condition=condition)

    @staticmethod
    def _zrici(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRICIOpCodes
        assert suffix == Suffix.ZRICI

        ra = InstructionInitializer._source_register()

        if (
            op_code in Instruction.AddRRICIOpCodes
            or op_code in Instruction.AndRRICIOpCodes
            or op_code in Instruction.SubRRICIOpCodes
        ):
            imm = InstructionInitializer._imm(Representation.SIGNED, 11)
        elif op_code in Instruction.AsrRRICIOpCodes:
            imm = InstructionInitializer._imm(Representation.UNSIGNED, 5)
        else:
            raise ValueError

        if op_code in Instruction.AddRRICIOpCodes:
            condition = InstructionInitializer._condition(AddNZCC.conditions())
        elif op_code in Instruction.AndRRICIOpCodes:
            condition = InstructionInitializer._condition(LogNZCC.conditions())
        elif op_code in Instruction.AsrRRICIOpCodes:
            condition = InstructionInitializer._condition(ImmShiftNZCC.conditions())
        elif op_code in Instruction.SubRRICIOpCodes:
            condition = InstructionInitializer._condition(SubNZCC.conditions())
        else:
            raise ValueError

        pc = InstructionInitializer._pc()

        return Instruction(op_code, suffix, ra=ra, imm=imm, condition=condition, pc=pc)

    @staticmethod
    def _zrif(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRIFOpCodes
        assert suffix == Suffix.ZRIF

        ra = InstructionInitializer._source_register()
        imm = InstructionInitializer._imm(Representation.SIGNED, 27)
        condition = InstructionInitializer._condition(FalseCC.conditions())

        return Instruction(op_code, suffix, ra=ra, imm=imm, condition=condition)

    @staticmethod
    def _zrr(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRROpCodes
        assert suffix == Suffix.ZRR

        ra = InstructionInitializer._source_register()
        rb = InstructionInitializer._source_register()

        return Instruction(op_code, suffix, ra=ra, rb=rb)

    @staticmethod
    def _zrrc(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRRCOpCodes
        assert suffix == Suffix.ZRRC

        ra = InstructionInitializer._source_register()
        rb = InstructionInitializer._source_register()

        if op_code in Instruction.AddRRRCOpCodes:
            condition = InstructionInitializer._condition(LogSetCC.conditions())
        elif op_code in Instruction.RsubRRRCOpCodes:
            condition = InstructionInitializer._condition(SubSetCC.conditions())
        elif op_code in Instruction.SubRRRCOpCodes:
            condition = InstructionInitializer._condition(ExtSubSetCC.conditions())
        else:
            raise ValueError

        return Instruction(op_code, suffix, ra=ra, rb=rb, condition=condition)

    @staticmethod
    def _zrrci(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRRCIOpCodes
        assert suffix == Suffix.ZRRCI

        ra = InstructionInitializer._source_register()
        rb = InstructionInitializer._source_register()

        if op_code in Instruction.AddRRRCIOpCodes:
            condition = InstructionInitializer._condition(AddNZCC.conditions())
        elif op_code in Instruction.AndRRRCIOpCodes:
            condition = InstructionInitializer._condition(LogNZCC.conditions())
        elif op_code in Instruction.AsrRRRCIOpCodes:
            condition = InstructionInitializer._condition(ShiftNZCC.conditions())
        elif op_code in Instruction.MulRRRCIOpCodes:
            condition = InstructionInitializer._condition(MulNZCC.conditions())
        elif op_code in Instruction.RsubRRRCIOpCodes:
            condition = InstructionInitializer._condition(SubNZCC.conditions())
        else:
            raise ValueError

        pc = InstructionInitializer._pc()

        return Instruction(op_code, suffix, ra=ra, rb=rb, condition=condition, pc=pc)

    @staticmethod
    def _s_rri(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRIOpCodes
        assert suffix == Suffix.S_RRI or suffix == Suffix.U_RRI

        dc = InstructionInitializer._pair_register()
        ra = InstructionInitializer._source_register()

        if op_code in Instruction.AddRRIOpCodes:
            imm = InstructionInitializer._imm(Representation.UNSIGNED, 32)
        elif op_code in Instruction.AsrRRIOpCodes:
            imm = InstructionInitializer._imm(Representation.UNSIGNED, 5)
        elif op_code in Instruction.CallRRIOpCodes:
            imm = InstructionInitializer._imm(Representation.SIGNED, 24)
        else:
            raise ValueError

        return Instruction(op_code, suffix, dc=dc, ra=ra, imm=imm)

    @staticmethod
    def _s_rric(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRICOpCodes
        assert suffix == Suffix.S_RRIC or suffix == Suffix.U_RRIC

        dc = InstructionInitializer._pair_register()
        ra = InstructionInitializer._source_register()

        if op_code in Instruction.AddRRICOpCodes or op_code in Instruction.SubRRICOpCodes:
            imm = InstructionInitializer._imm(Representation.SIGNED, 24)
        elif op_code in Instruction.AsrRRICOpCodes:
            imm = InstructionInitializer._imm(Representation.UNSIGNED, 5)
        else:
            raise ValueError

        if op_code in Instruction.AddRRICOpCodes or op_code in Instruction.AsrRRICOpCodes:
            condition = InstructionInitializer._condition(LogSetCC.conditions())
        elif op_code in Instruction.SubRRICOpCodes:
            condition = InstructionInitializer._condition(ExtSubSetCC.conditions())
        else:
            raise ValueError

        return Instruction(op_code, suffix, dc=dc, ra=ra, imm=imm, condition=condition)

    @staticmethod
    def _s_rrici(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRICIOpCodes
        assert suffix == Suffix.S_RRICI or suffix == Suffix.U_RRICI

        dc = InstructionInitializer._pair_register()
        ra = InstructionInitializer._source_register()

        if (
            op_code in Instruction.AddRRICIOpCodes
            or op_code in Instruction.AndRRICIOpCodes
            or op_code in Instruction.SubRRICIOpCodes
        ):
            imm = InstructionInitializer._imm(Representation.SIGNED, 8)
        elif op_code in Instruction.AsrRRICIOpCodes:
            imm = InstructionInitializer._imm(Representation.UNSIGNED, 5)
        else:
            raise ValueError

        if op_code in Instruction.AddRRICIOpCodes:
            condition = InstructionInitializer._condition(AddNZCC.conditions())
        elif op_code in Instruction.AndRRICIOpCodes:
            condition = InstructionInitializer._condition(LogNZCC.conditions())
        elif op_code in Instruction.AsrRRICIOpCodes:
            condition = InstructionInitializer._condition(ImmShiftNZCC.conditions())
        elif op_code in Instruction.SubRRICIOpCodes:
            condition = InstructionInitializer._condition(SubNZCC.conditions())
        else:
            raise ValueError

        pc = InstructionInitializer._pc()

        return Instruction(op_code, suffix, dc=dc, ra=ra, imm=imm, condition=condition, pc=pc)

    @staticmethod
    def _s_rrif(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRIFOpCodes
        assert suffix == Suffix.S_RRIF or suffix == Suffix.U_RRIF

        dc = InstructionInitializer._pair_register()
        ra = InstructionInitializer._source_register()
        imm = InstructionInitializer._imm(Representation.SIGNED, 24)
        condition = InstructionInitializer._condition(FalseCC.conditions())

        return Instruction(op_code, suffix, dc=dc, ra=ra, imm=imm, condition=condition)

    @staticmethod
    def _s_rrr(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRROpCodes
        assert suffix == Suffix.S_RRR or suffix == Suffix.U_RRR

        dc = InstructionInitializer._pair_register()
        ra = InstructionInitializer._source_register()
        rb = InstructionInitializer._source_register()

        return Instruction(op_code, suffix, dc=dc, ra=ra, rb=rb)

    @staticmethod
    def _s_rrrc(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRRCOpCodes
        assert suffix == Suffix.S_RRRC or suffix == Suffix.U_RRRC

        dc = InstructionInitializer._pair_register()
        ra = InstructionInitializer._source_register()
        rb = InstructionInitializer._source_register()

        if op_code in Instruction.AddRRRCOpCodes:
            condition = InstructionInitializer._condition(LogSetCC.conditions())
        elif op_code in Instruction.RsubRRRCOpCodes:
            condition = InstructionInitializer._condition(SubSetCC.conditions())
        elif op_code in Instruction.SubRRRCOpCodes:
            condition = InstructionInitializer._condition(ExtSubSetCC.conditions())
        else:
            raise ValueError

        return Instruction(op_code, suffix, dc=dc, ra=ra, rb=rb, condition=condition)

    @staticmethod
    def _s_rrrci(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRRCIOpCodes
        assert suffix == Suffix.S_RRRCI or suffix == Suffix.U_RRRCI

        dc = InstructionInitializer._pair_register()
        ra = InstructionInitializer._source_register()
        rb = InstructionInitializer._source_register()

        if op_code in Instruction.AddRRRCIOpCodes:
            condition = InstructionInitializer._condition(AddNZCC.conditions())
        elif op_code in Instruction.AndRRRCIOpCodes:
            condition = InstructionInitializer._condition(LogNZCC.conditions())
        elif op_code in Instruction.AsrRRRCIOpCodes:
            condition = InstructionInitializer._condition(ShiftNZCC.conditions())
        elif op_code in Instruction.MulRRRCIOpCodes:
            condition = InstructionInitializer._condition(MulNZCC.conditions())
        elif op_code in Instruction.RsubRRRCIOpCodes:
            condition = InstructionInitializer._condition(SubNZCC.conditions())
        else:
            raise ValueError

        pc = InstructionInitializer._pc()

        return Instruction(op_code, suffix, dc=dc, ra=ra, rb=rb, condition=condition, pc=pc)

    @staticmethod
    def _rr(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RROpCodes
        assert suffix == Suffix.RR

        rc = InstructionInitializer._gp_register()
        ra = InstructionInitializer._source_register()

        return Instruction(op_code, suffix, rc=rc, ra=ra)

    @staticmethod
    def _rrc(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRCOpCodes
        assert suffix == Suffix.RRC

        rc = InstructionInitializer._gp_register()
        ra = InstructionInitializer._source_register()
        condition = InstructionInitializer._condition(LogSetCC.conditions())

        return Instruction(op_code, suffix, rc=rc, ra=ra, condition=condition)

    @staticmethod
    def _rrci(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRCIOpCodes
        assert suffix == Suffix.RRCI

        rc = InstructionInitializer._gp_register()
        ra = InstructionInitializer._source_register()

        if op_code in Instruction.CaoRRCIOpCodes:
            condition = InstructionInitializer._condition(CountNZCC.conditions())
        elif op_code in Instruction.ExtsbRRCIOpCodes:
            condition = InstructionInitializer._condition(LogNZCC.conditions())
        elif op_code in Instruction.TimeCfgRRCIOpCodes:
            condition = InstructionInitializer._condition(TrueCC.conditions())
        else:
            raise ValueError

        pc = InstructionInitializer._pc()

        return Instruction(op_code, suffix, rc=rc, ra=ra, condition=condition, pc=pc)

    @staticmethod
    def _zr(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RROpCodes
        assert suffix == Suffix.ZR

        ra = InstructionInitializer._source_register()

        return Instruction(op_code, suffix, ra=ra)

    @staticmethod
    def _zrc(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRCOpCodes
        assert suffix == Suffix.ZRC

        ra = InstructionInitializer._source_register()
        condition = InstructionInitializer._condition(LogSetCC.conditions())

        return Instruction(op_code, suffix, ra=ra, condition=condition)

    @staticmethod
    def _zrci(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRCIOpCodes
        assert suffix == Suffix.ZRCI

        ra = InstructionInitializer._source_register()

        if op_code in Instruction.CaoRRCIOpCodes:
            condition = InstructionInitializer._condition(CountNZCC.conditions())
        elif op_code in Instruction.ExtsbRRCIOpCodes:
            condition = InstructionInitializer._condition(LogNZCC.conditions())
        elif op_code in Instruction.TimeCfgRRCIOpCodes:
            condition = InstructionInitializer._condition(TrueCC.conditions())
        else:
            raise ValueError

        pc = InstructionInitializer._pc()

        return Instruction(op_code, suffix, ra=ra, condition=condition, pc=pc)

    @staticmethod
    def _s_rr(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RROpCodes
        assert suffix == Suffix.S_RR or suffix == Suffix.U_RR

        dc = InstructionInitializer._pair_register()
        ra = InstructionInitializer._source_register()

        return Instruction(op_code, suffix, dc=dc, ra=ra)

    @staticmethod
    def _s_rrc(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRCOpCodes
        assert suffix == Suffix.S_RRC or suffix == Suffix.U_RRC

        dc = InstructionInitializer._pair_register()
        ra = InstructionInitializer._source_register()
        condition = InstructionInitializer._condition(LogSetCC.conditions())

        return Instruction(op_code, suffix, dc=dc, ra=ra, condition=condition)

    @staticmethod
    def _s_rrci(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRCIOpCodes
        assert suffix == Suffix.S_RRCI or suffix == Suffix.U_RRCI

        dc = InstructionInitializer._pair_register()
        ra = InstructionInitializer._source_register()

        if op_code in Instruction.CaoRRCIOpCodes:
            condition = InstructionInitializer._condition(CountNZCC.conditions())
        elif op_code in Instruction.ExtsbRRCIOpCodes:
            condition = InstructionInitializer._condition(LogNZCC.conditions())
        elif op_code in Instruction.TimeCfgRRCIOpCodes:
            condition = InstructionInitializer._condition(TrueCC.conditions())
        else:
            raise ValueError

        pc = InstructionInitializer._pc()

        return Instruction(op_code, suffix, dc=dc, ra=ra, condition=condition, pc=pc)

    @staticmethod
    def _drdici(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.DRDICIOpCodes
        assert suffix == Suffix.DRDICI

        dc = InstructionInitializer._pair_register()
        ra = InstructionInitializer._source_register()
        db = InstructionInitializer._pair_register()
        imm = InstructionInitializer._imm(Representation.UNSIGNED, 5)

        if op_code in Instruction.DivStepDRDICIOpCodes:
            condition = InstructionInitializer._condition(DivCC.conditions())
        elif op_code in Instruction.MulStepDRDICIOpCodes:
            condition = InstructionInitializer._condition(BootCC.conditions())
        else:
            raise ValueError

        pc = InstructionInitializer._pc()

        return Instruction(op_code, suffix, dc=dc, ra=ra, db=db, imm=imm, condition=condition, pc=pc)

    @staticmethod
    def _rrri(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRRIOpCodes
        assert suffix == Suffix.RRRI

        rc = InstructionInitializer._gp_register()
        ra = InstructionInitializer._source_register()
        rb = InstructionInitializer._source_register()
        imm = InstructionInitializer._imm(Representation.UNSIGNED, 5)

        return Instruction(op_code, suffix, rc=rc, ra=ra, rb=rb, imm=imm)

    @staticmethod
    def _rrrici(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRRICIOpCodes
        assert suffix == Suffix.RRRICI

        rc = InstructionInitializer._gp_register()
        ra = InstructionInitializer._source_register()
        rb = InstructionInitializer._source_register()
        imm = InstructionInitializer._imm(Representation.UNSIGNED, 5)
        condition = InstructionInitializer._condition(DivNZCC.conditions())
        pc = InstructionInitializer._pc()

        return Instruction(op_code, suffix, rc=rc, ra=ra, rb=rb, imm=imm, condition=condition, pc=pc)

    @staticmethod
    def _zrri(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRRIOpCodes
        assert suffix == Suffix.ZRRI

        ra = InstructionInitializer._source_register()
        rb = InstructionInitializer._source_register()
        imm = InstructionInitializer._imm(Representation.UNSIGNED, 5)

        return Instruction(op_code, suffix, ra=ra, rb=rb, imm=imm)

    @staticmethod
    def _zrrici(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRRICIOpCodes
        assert suffix == Suffix.ZRRICI

        ra = InstructionInitializer._source_register()
        rb = InstructionInitializer._source_register()
        imm = InstructionInitializer._imm(Representation.UNSIGNED, 5)
        condition = InstructionInitializer._condition(DivNZCC.conditions())
        pc = InstructionInitializer._pc()

        return Instruction(op_code, suffix, ra=ra, rb=rb, imm=imm, condition=condition, pc=pc)

    @staticmethod
    def _s_rrri(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRRIOpCodes
        assert suffix == Suffix.S_RRRI or suffix == Suffix.U_RRRI

        dc = InstructionInitializer._pair_register()
        ra = InstructionInitializer._source_register()
        rb = InstructionInitializer._source_register()
        imm = InstructionInitializer._imm(Representation.UNSIGNED, 5)

        return Instruction(op_code, suffix, dc=dc, ra=ra, rb=rb, imm=imm)

    @staticmethod
    def _s_rrrici(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RRRICIOpCodes
        assert suffix == Suffix.S_RRRICI or suffix == Suffix.U_RRRICI

        dc = InstructionInitializer._pair_register()
        ra = InstructionInitializer._source_register()
        rb = InstructionInitializer._source_register()
        imm = InstructionInitializer._imm(Representation.UNSIGNED, 5)
        condition = InstructionInitializer._condition(DivNZCC.conditions())
        pc = InstructionInitializer._pc()

        return Instruction(op_code, suffix, dc=dc, ra=ra, rb=rb, imm=imm, condition=condition, pc=pc)

    @staticmethod
    def _rir(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RIROpCodes
        assert suffix == Suffix.RIR

        rc = InstructionInitializer._gp_register()
        imm = InstructionInitializer._imm(Representation.UNSIGNED, 32)
        ra = InstructionInitializer._source_register()

        return Instruction(op_code, suffix, rc=rc, imm=imm, ra=ra)

    @staticmethod
    def _rirc(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RIRCOpCodes
        assert suffix == Suffix.RIRC

        rc = InstructionInitializer._gp_register()
        imm = InstructionInitializer._imm(Representation.SIGNED, 24)
        ra = InstructionInitializer._source_register()
        condition = InstructionInitializer._condition(SubSetCC.conditions())

        return Instruction(op_code, suffix, rc=rc, imm=imm, ra=ra, condition=condition)

    @staticmethod
    def _rirci(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RIRCIOpCodes
        assert suffix == Suffix.RIRCI

        rc = InstructionInitializer._gp_register()
        imm = InstructionInitializer._imm(Representation.SIGNED, 8)
        ra = InstructionInitializer._source_register()
        condition = InstructionInitializer._condition(SubNZCC.conditions())
        pc = InstructionInitializer._pc()

        return Instruction(op_code, suffix, rc=rc, imm=imm, ra=ra, condition=condition, pc=pc)

    @staticmethod
    def _zir(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RIROpCodes
        assert suffix == Suffix.ZIR

        imm = InstructionInitializer._imm(Representation.UNSIGNED, 32)
        ra = InstructionInitializer._source_register()

        return Instruction(op_code, suffix, imm=imm, ra=ra)

    @staticmethod
    def _zirc(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RIRCOpCodes
        assert suffix == Suffix.ZIRC

        imm = InstructionInitializer._imm(Representation.SIGNED, 27)
        ra = InstructionInitializer._source_register()
        condition = InstructionInitializer._condition(SubSetCC.conditions())

        return Instruction(op_code, suffix, imm=imm, ra=ra, condition=condition)

    @staticmethod
    def _zirci(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RIRCIOpCodes
        assert suffix == Suffix.ZIRCI

        imm = InstructionInitializer._imm(Representation.SIGNED, 11)
        ra = InstructionInitializer._source_register()
        condition = InstructionInitializer._condition(SubNZCC.conditions())
        pc = InstructionInitializer._pc()

        return Instruction(op_code, suffix, imm=imm, ra=ra, condition=condition, pc=pc)

    @staticmethod
    def _s_rirc(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RIRCOpCodes
        assert suffix == Suffix.S_RIRC or suffix == Suffix.U_RIRC

        dc = InstructionInitializer._pair_register()
        imm = InstructionInitializer._imm(Representation.SIGNED, 24)
        ra = InstructionInitializer._source_register()
        condition = InstructionInitializer._condition(SubSetCC.conditions())

        return Instruction(op_code, suffix, dc=dc, imm=imm, ra=ra, condition=condition)

    @staticmethod
    def _s_rirci(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RIRCIOpCodes
        assert suffix == Suffix.S_RIRCI or suffix == Suffix.U_RIRCI

        dc = InstructionInitializer._pair_register()
        imm = InstructionInitializer._imm(Representation.SIGNED, 8)
        ra = InstructionInitializer._source_register()
        condition = InstructionInitializer._condition(SubNZCC.conditions())
        pc = InstructionInitializer._pc()

        return Instruction(op_code, suffix, dc=dc, imm=imm, ra=ra, condition=condition, pc=pc)

    @staticmethod
    def _r(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.ROpCodes
        assert suffix == Suffix.R

        rc = InstructionInitializer._gp_register()

        return Instruction(op_code, suffix, rc=rc)

    @staticmethod
    def _rci(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RCIOpCodes
        assert suffix == Suffix.RCI

        rc = InstructionInitializer._gp_register()
        condition = InstructionInitializer._condition(TrueCC.conditions())
        pc = InstructionInitializer._pc()

        return Instruction(op_code, suffix, rc=rc, condition=condition, pc=pc)

    @staticmethod
    def _z(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.ROpCodes or op_code == OpCode.NOP
        assert suffix == Suffix.Z

        return Instruction(op_code, suffix)

    @staticmethod
    def _zci(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RCIOpCodes
        assert suffix == Suffix.ZCI

        condition = InstructionInitializer._condition(TrueCC.conditions())
        pc = InstructionInitializer._pc()

        return Instruction(op_code, suffix, condition=condition, pc=pc)

    @staticmethod
    def _s_r(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.ROpCodes
        assert suffix == Suffix.S_R or suffix == Suffix.U_R

        dc = InstructionInitializer._pair_register()

        return Instruction(op_code, suffix, dc=dc)

    @staticmethod
    def _s_rci(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.RCIOpCodes
        assert suffix == Suffix.S_RCI or suffix == Suffix.U_RCI

        dc = InstructionInitializer._pair_register()
        condition = InstructionInitializer._condition(TrueCC.conditions())
        pc = InstructionInitializer._pc()

        return Instruction(op_code, suffix, dc=dc, condition=condition, pc=pc)

    @staticmethod
    def _ci(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.CIOpCodes
        assert suffix == Suffix.CI

        condition = InstructionInitializer._condition(BootCC.conditions())
        pc = InstructionInitializer._pc()

        return Instruction(op_code, suffix, condition=condition, pc=pc)

    @staticmethod
    def _i(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.IOpCodes
        assert suffix == Suffix.I

        imm = InstructionInitializer._imm(Representation.SIGNED, 24)

        return Instruction(op_code, suffix, imm=imm)

    @staticmethod
    def _ddci(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.DDCIOpCodes
        assert suffix == Suffix.DDCI

        dc = InstructionInitializer._pair_register()
        db = InstructionInitializer._pair_register()
        condition = InstructionInitializer._condition(TrueFalseCC.conditions())
        pc = InstructionInitializer._pc()

        return Instruction(op_code, suffix, dc=dc, db=db, condition=condition, pc=pc)

    @staticmethod
    def _erri(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.ERRIOpCodes
        assert suffix == Suffix.ERRI

        endian = InstructionInitializer._endian()
        rc = InstructionInitializer._gp_register()
        ra = InstructionInitializer._source_register()
        off = InstructionInitializer._imm(Representation.SIGNED, 24)

        return Instruction(op_code, suffix, endian=endian, rc=rc, ra=ra, off=off)

    @staticmethod
    def _s_erri(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.ERRIOpCodes
        assert suffix == Suffix.S_ERRI or suffix == Suffix.U_ERRI

        endian = InstructionInitializer._endian()
        dc = InstructionInitializer._pair_register()
        ra = InstructionInitializer._source_register()
        off = InstructionInitializer._imm(Representation.SIGNED, 24)

        return Instruction(op_code, suffix, endian=endian, dc=dc, ra=ra, off=off)

    @staticmethod
    def _edri(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.EDRIOpCodes
        assert suffix == Suffix.EDRI

        endian = InstructionInitializer._endian()
        dc = InstructionInitializer._pair_register()
        ra = InstructionInitializer._source_register()
        off = InstructionInitializer._imm(Representation.SIGNED, 24)

        return Instruction(op_code, suffix, endian=endian, dc=dc, ra=ra, off=off)

    @staticmethod
    def _erii(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.ERIIOpCodes
        assert suffix == Suffix.ERII

        endian = InstructionInitializer._endian()
        ra = InstructionInitializer._source_register()
        # NOTE(bongjoon.hyun@gmail.com): original width is 12
        off = InstructionInitializer._imm(Representation.SIGNED, 24)
        # NOTE(bongjoon.hyun@gmail.com): original width is 8
        imm = InstructionInitializer._imm(Representation.SIGNED, 16)

        return Instruction(op_code, suffix, endian=endian, ra=ra, off=off, imm=imm)

    @staticmethod
    def _erir(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.ERIROpCodes
        assert suffix == Suffix.ERIR

        endian = InstructionInitializer._endian()
        ra = InstructionInitializer._source_register()
        off = InstructionInitializer._imm(Representation.SIGNED, 24)
        rb = InstructionInitializer._source_register()

        return Instruction(op_code, suffix, endian=endian, ra=ra, off=off, rb=rb)

    @staticmethod
    def _erid(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.ERIDOpCodes
        assert suffix == Suffix.ERID

        endian = InstructionInitializer._endian()
        ra = InstructionInitializer._source_register()
        off = InstructionInitializer._imm(Representation.SIGNED, 24)
        db = InstructionInitializer._pair_register()

        return Instruction(op_code, suffix, endian=endian, ra=ra, off=off, db=db)

    @staticmethod
    def _dma_rri(op_code: OpCode, suffix: Suffix) -> Instruction:
        assert op_code in Instruction.DMARRIOpCodes
        assert suffix == Suffix.DMA_RRI

        ra = InstructionInitializer._source_register()
        rb = InstructionInitializer._source_register()
        imm = InstructionInitializer._imm(Representation.UNSIGNED, 8)

        return Instruction(op_code, suffix, ra=ra, rb=rb, imm=imm)
