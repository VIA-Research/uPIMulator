from typing import Optional, Union

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
from abi.isa.instruction.op_code import OpCode
from abi.isa.instruction.suffix import Suffix
from abi.isa.register.gp_register import GPRegister
from abi.isa.register.pair_register import PairRegister
from abi.isa.register.sp_register import SPRegister
from abi.word.immediate import Immediate
from abi.word.representation import Representation
from util.config_loader import ConfigLoader


class Instruction:
    DestinationRegister = Union[GPRegister, SPRegister]
    SourceRegister = Union[GPRegister, SPRegister]

    AcquireRICIOpCodes = {OpCode.ACQUIRE}
    ReleaseRICIOpCodes = {OpCode.RELEASE}
    BootRICIOpCodes = {OpCode.BOOT, OpCode.RESUME}
    RICIOpCodes = {*AcquireRICIOpCodes, *ReleaseRICIOpCodes, *BootRICIOpCodes}

    AddRRIOpCodes = {OpCode.ADD, OpCode.ADDC, OpCode.AND, OpCode.OR, OpCode.XOR}
    AsrRRIOpCodes = {
        OpCode.ASR,
        OpCode.LSL,
        OpCode.LSL1,
        OpCode.LSL1X,
        OpCode.LSLX,
        OpCode.LSR,
        OpCode.LSR1,
        OpCode.LSR1X,
        OpCode.LSRX,
        OpCode.ROL,
        OpCode.ROR,
    }
    CallRRIOpCodes = {OpCode.CALL}
    RRIOpCodes = {*AddRRIOpCodes, *AsrRRIOpCodes, *CallRRIOpCodes}

    AddRRICOpCodes = {
        OpCode.ADD,
        OpCode.ADDC,
        OpCode.AND,
        OpCode.ANDN,
        OpCode.NAND,
        OpCode.NOR,
        OpCode.NXOR,
        OpCode.OR,
        OpCode.ORN,
        OpCode.XOR,
        OpCode.HASH,
    }
    AsrRRICOpCodes = {
        OpCode.ASR,
        OpCode.LSL,
        OpCode.LSL1,
        OpCode.LSL1X,
        OpCode.LSLX,
        OpCode.LSR,
        OpCode.LSR1,
        OpCode.LSR1X,
        OpCode.LSRX,
        OpCode.ROL,
        OpCode.ROR,
    }
    SubRRICOpCodes = {OpCode.SUB, OpCode.SUBC}
    RRICOpCodes = {*AddRRICOpCodes, *AsrRRICOpCodes, *SubRRICOpCodes}

    AddRRICIOpCodes = {OpCode.ADD, OpCode.ADDC}
    AndRRICIOpCodes = {
        OpCode.AND,
        OpCode.ANDN,
        OpCode.NAND,
        OpCode.NOR,
        OpCode.NXOR,
        OpCode.OR,
        OpCode.ORN,
        OpCode.XOR,
        OpCode.HASH,
    }
    AsrRRICIOpCodes = {
        OpCode.ASR,
        OpCode.LSL,
        OpCode.LSL1,
        OpCode.LSL1X,
        OpCode.LSLX,
        OpCode.LSR,
        OpCode.LSR1,
        OpCode.LSR1X,
        OpCode.LSRX,
        OpCode.ROL,
        OpCode.ROR,
    }
    SubRRICIOpCodes = {OpCode.SUB, OpCode.SUBC}
    RRICIOpCodes = {
        *AddRRICIOpCodes,
        *AndRRICIOpCodes,
        *AsrRRICIOpCodes,
        *SubRRICIOpCodes,
    }

    RRIFOpCodes = {
        OpCode.ADD,
        OpCode.ADDC,
        OpCode.AND,
        OpCode.ANDN,
        OpCode.NAND,
        OpCode.NOR,
        OpCode.NXOR,
        OpCode.OR,
        OpCode.ORN,
        OpCode.SUB,
        OpCode.SUBC,
        OpCode.XOR,
        OpCode.HASH,
    }

    RRROpCodes = {
        OpCode.ADD,
        OpCode.ADDC,
        OpCode.AND,
        OpCode.ANDN,
        OpCode.ASR,
        OpCode.CMPB4,
        OpCode.LSL,
        OpCode.LSL1,
        OpCode.LSL1X,
        OpCode.LSLX,
        OpCode.LSR,
        OpCode.LSR1,
        OpCode.LSR1X,
        OpCode.LSRX,
        OpCode.MUL_SH_SH,
        OpCode.MUL_SH_SL,
        OpCode.MUL_SH_UH,
        OpCode.MUL_SH_UL,
        OpCode.MUL_SL_SH,
        OpCode.MUL_SL_SL,
        OpCode.MUL_SL_UH,
        OpCode.MUL_SL_UL,
        OpCode.MUL_UH_UH,
        OpCode.MUL_UH_UL,
        OpCode.MUL_UL_UH,
        OpCode.MUL_UL_UL,
        OpCode.NAND,
        OpCode.NOR,
        OpCode.NXOR,
        OpCode.OR,
        OpCode.ORN,
        OpCode.ROL,
        OpCode.ROR,
        OpCode.RSUB,
        OpCode.RSUBC,
        OpCode.SUB,
        OpCode.SUBC,
        OpCode.XOR,
        OpCode.HASH,
        OpCode.CALL,
    }

    AddRRRCOpCodes = {
        OpCode.ADD,
        OpCode.ADDC,
        OpCode.AND,
        OpCode.ANDN,
        OpCode.ASR,
        OpCode.CMPB4,
        OpCode.LSL,
        OpCode.LSL1,
        OpCode.LSL1X,
        OpCode.LSLX,
        OpCode.LSR,
        OpCode.LSR1,
        OpCode.LSR1X,
        OpCode.LSRX,
        OpCode.MUL_SH_SH,
        OpCode.MUL_SH_SL,
        OpCode.MUL_SH_UH,
        OpCode.MUL_SH_UL,
        OpCode.MUL_SL_SH,
        OpCode.MUL_SL_SL,
        OpCode.MUL_SL_UH,
        OpCode.MUL_SL_UL,
        OpCode.MUL_UH_UH,
        OpCode.MUL_UH_UL,
        OpCode.MUL_UL_UH,
        OpCode.MUL_UL_UL,
        OpCode.NAND,
        OpCode.NOR,
        OpCode.NXOR,
        OpCode.ROL,
        OpCode.ROR,
        OpCode.OR,
        OpCode.ORN,
        OpCode.XOR,
        OpCode.HASH,
        OpCode.CALL,
    }
    RsubRRRCOpCodes = {OpCode.RSUB, OpCode.RSUBC}
    SubRRRCOpCodes = {OpCode.SUB, OpCode.SUBC}
    RRRCOpCodes = {*AddRRRCOpCodes, *RsubRRRCOpCodes, *SubRRRCOpCodes}

    AddRRRCIOpCodes = {OpCode.ADD, OpCode.ADDC}
    AndRRRCIOpCodes = {
        OpCode.AND,
        OpCode.ANDN,
        OpCode.NAND,
        OpCode.NOR,
        OpCode.NXOR,
        OpCode.OR,
        OpCode.ORN,
        OpCode.XOR,
        OpCode.HASH,
    }
    AsrRRRCIOpCodes = {
        OpCode.ASR,
        OpCode.CMPB4,
        OpCode.LSL,
        OpCode.LSL1,
        OpCode.LSL1X,
        OpCode.LSLX,
        OpCode.LSR,
        OpCode.LSR1,
        OpCode.LSR1X,
        OpCode.LSRX,
        OpCode.ROL,
        OpCode.ROR,
    }
    MulRRRCIOpCodes = {
        OpCode.MUL_SH_SH,
        OpCode.MUL_SH_SL,
        OpCode.MUL_SH_UH,
        OpCode.MUL_SH_UL,
        OpCode.MUL_SL_SH,
        OpCode.MUL_SL_SL,
        OpCode.MUL_SL_UH,
        OpCode.MUL_SL_UL,
        OpCode.MUL_UH_UH,
        OpCode.MUL_UH_UL,
        OpCode.MUL_UL_UH,
        OpCode.MUL_UL_UL,
    }
    RsubRRRCIOpCodes = {OpCode.RSUB, OpCode.RSUBC, OpCode.SUB, OpCode.SUBC}
    RRRCIOpCodes = {
        *AddRRRCIOpCodes,
        *AndRRRCIOpCodes,
        *AsrRRRCIOpCodes,
        *RsubRRRCIOpCodes,
    }

    RROpCodes = {
        OpCode.CAO,
        OpCode.CLO,
        OpCode.CLS,
        OpCode.CLZ,
        OpCode.EXTSB,
        OpCode.EXTSH,
        OpCode.EXTUB,
        OpCode.EXTUH,
        OpCode.SATS,
        OpCode.TIME_CFG,
    }

    RRCOpCodes = {
        OpCode.CAO,
        OpCode.CLO,
        OpCode.CLS,
        OpCode.CLZ,
        OpCode.EXTSB,
        OpCode.EXTSH,
        OpCode.EXTUB,
        OpCode.EXTUH,
        OpCode.SATS,
    }

    CaoRRCIOpCodes = {OpCode.CAO, OpCode.CLO, OpCode.CLS, OpCode.CLZ}
    ExtsbRRCIOpCodes = {
        OpCode.EXTSB,
        OpCode.EXTSH,
        OpCode.EXTUB,
        OpCode.EXTUH,
        OpCode.SATS,
    }
    TimeCfgRRCIOpCodes = {OpCode.TIME_CFG}
    RRCIOpCodes = {*CaoRRCIOpCodes, *ExtsbRRCIOpCodes, *TimeCfgRRCIOpCodes}

    DivStepDRDICIOpCodes = {OpCode.DIV_STEP}
    MulStepDRDICIOpCodes = {OpCode.MUL_STEP}
    DRDICIOpCodes = {*DivStepDRDICIOpCodes, *MulStepDRDICIOpCodes}

    RRRIOpCodes = {OpCode.LSL_ADD, OpCode.LSL_SUB, OpCode.LSR_ADD, OpCode.ROL_ADD}
    RRRICIOpCodes = {OpCode.LSL_ADD, OpCode.LSL_SUB, OpCode.LSR_ADD, OpCode.ROL_ADD}

    RIROpCodes = {OpCode.SUB, OpCode.SUBC}
    RIRCOpCodes = {OpCode.SUB, OpCode.SUBC}
    RIRCIOpCodes = {OpCode.SUB, OpCode.SUBC}

    ROpCodes = {OpCode.TIME}
    RCIOpCodes = {OpCode.TIME}

    CIOpCodes = {OpCode.STOP}
    IOpCodes = {OpCode.FAULT}
    DDCIOpCodes = {OpCode.MOVD, OpCode.SWAPD}

    ERRIOpCodes = {OpCode.LBS, OpCode.LBU, OpCode.LHS, OpCode.LHU, OpCode.LW}
    EDRIOpCodes = {OpCode.LD}

    ERIIOpCodes = {
        OpCode.SB,
        OpCode.SB_ID,
        OpCode.SD,
        OpCode.SD_ID,
        OpCode.SH,
        OpCode.SH_ID,
        OpCode.SW,
        OpCode.SW_ID,
        OpCode.SD,
        OpCode.SD_ID,
    }
    ERIROpCodes = {OpCode.SB, OpCode.SH, OpCode.SW}
    ERIDOpCodes = {OpCode.SD}

    DMARRIOpCodes = {OpCode.LDMA, OpCode.LDMAI, OpCode.SDMA}

    def __init__(
        self,
        op_code: OpCode,
        suffix: Suffix,
        rc: Optional[DestinationRegister] = None,
        ra: Optional[SourceRegister] = None,
        rb: Optional[SourceRegister] = None,
        dc: Optional[PairRegister] = None,
        db: Optional[PairRegister] = None,
        condition: Optional[Condition] = None,
        imm: Optional[int] = None,
        off: Optional[int] = None,
        pc: Optional[int] = None,
        endian: Optional[Endian] = None,
    ):
        self._op_code: OpCode = op_code
        self._suffix: Suffix = suffix
        self._rc: Union[GPRegister, SPRegister, None] = None
        self._ra: Union[GPRegister, SPRegister, None] = None
        self._rb: Union[GPRegister, SPRegister, None] = None
        self._dc: Optional[PairRegister] = None
        self._db: Optional[PairRegister] = None
        self._condition: Optional[Condition] = None
        self._imm: Optional[Immediate] = None
        self._pc: Optional[Immediate] = None
        self._off: Optional[Immediate] = None
        self._endian: Optional[Endian] = None

        self._init_instruction(rc, ra, rb, dc, db, condition, imm, off, pc, endian)

    def op_code(self) -> OpCode:
        return self._op_code

    def suffix(self) -> Suffix:
        return self._suffix

    def rc(self) -> Union[GPRegister, SPRegister]:
        assert self._rc is not None
        return self._rc

    def ra(self) -> Union[GPRegister, SPRegister]:
        assert self._ra is not None
        return self._ra

    def rb(self) -> Union[GPRegister, SPRegister]:
        assert self._rb is not None
        return self._rb

    def dc(self) -> PairRegister:
        assert self._dc is not None
        return self._dc

    def db(self) -> PairRegister:
        assert self._db is not None
        return self._db

    def condition(self) -> Condition:
        assert self._condition is not None
        return self._condition

    def imm(self) -> Immediate:
        assert self._imm is not None
        return self._imm

    def off(self) -> Immediate:
        assert self._off is not None
        return self._off

    def pc(self) -> Immediate:
        assert self._pc is not None
        return self._pc

    def endian(self) -> Endian:
        assert self._endian is not None
        return self._endian

    def _init_instruction(
        self,
        rc: Optional[DestinationRegister],
        ra: Optional[SourceRegister],
        rb: Optional[SourceRegister],
        dc: Optional[PairRegister],
        db: Optional[PairRegister],
        condition: Optional[Condition],
        imm: Optional[int],
        off: Optional[int],
        pc: Optional[int],
        endian: Optional[Endian],
    ) -> None:
        if self._suffix == Suffix.RICI:
            assert ra is not None
            assert imm is not None
            assert condition is not None
            assert pc is not None

            self._init_rici(ra, imm, condition, pc)
        elif self._suffix == Suffix.RRI:
            assert isinstance(rc, GPRegister)
            assert ra is not None
            assert imm is not None

            self._init_rri(rc, ra, imm)
        elif self._suffix == Suffix.RRIC:
            assert isinstance(rc, GPRegister)
            assert ra is not None
            assert imm is not None
            assert condition is not None

            self._init_rric(rc, ra, imm, condition)
        elif self._suffix == Suffix.RRICI:
            assert isinstance(rc, GPRegister)
            assert ra is not None
            assert imm is not None
            assert condition is not None
            assert pc is not None

            self._init_rrici(rc, ra, imm, condition, pc)
        elif self._suffix == Suffix.RRIF:
            assert isinstance(rc, GPRegister)
            assert ra is not None
            assert imm is not None
            assert condition is not None

            self._init_rrif(rc, ra, imm, condition)
        elif self._suffix == Suffix.RRR:
            assert isinstance(rc, GPRegister)
            assert ra is not None
            assert rb is not None

            self._init_rrr(rc, ra, rb)
        elif self._suffix == Suffix.RRRC:
            assert isinstance(rc, GPRegister)
            assert ra is not None
            assert rb is not None
            assert condition is not None

            self._init_rrrc(rc, ra, rb, condition)
        elif self._suffix == Suffix.RRRCI:
            assert isinstance(rc, GPRegister)
            assert ra is not None
            assert rb is not None
            assert condition is not None
            assert pc is not None

            self._init_rrrci(rc, ra, rb, condition, pc)
        elif self._suffix == Suffix.ZRI:
            assert ra is not None
            assert imm is not None

            self._init_zri(ra, imm)
        elif self._suffix == Suffix.ZRIC:
            assert ra is not None
            assert imm is not None
            assert condition is not None

            self._init_zric(ra, imm, condition)
        elif self._suffix == Suffix.ZRICI:
            assert ra is not None
            assert imm is not None
            assert condition is not None
            assert pc is not None

            self._init_zrici(ra, imm, condition, pc)
        elif self._suffix == Suffix.ZRIF:
            assert ra is not None
            assert imm is not None
            assert condition is not None

            self._init_zrif(ra, imm, condition)
        elif self._suffix == Suffix.ZRR:
            assert ra is not None
            assert rb is not None

            self._init_zrr(ra, rb)
        elif self._suffix == Suffix.ZRRC:
            assert ra is not None
            assert rb is not None
            assert condition is not None

            self._init_zrrc(ra, rb, condition)
        elif self._suffix == Suffix.ZRRCI:
            assert ra is not None
            assert rb is not None
            assert condition is not None
            assert pc is not None

            self._init_zrrci(ra, rb, condition, pc)
        elif self._suffix == Suffix.S_RRI or self._suffix == Suffix.U_RRI:
            assert dc is not None
            assert ra is not None
            assert imm is not None

            self._init_s_rri(dc, ra, imm)
        elif self._suffix == Suffix.S_RRIC or self._suffix == Suffix.U_RRIC:
            assert dc is not None
            assert ra is not None
            assert imm is not None
            assert condition is not None

            self._init_s_rric(dc, ra, imm, condition)
        elif self._suffix == Suffix.S_RRICI or self._suffix == Suffix.U_RRICI:
            assert dc is not None
            assert ra is not None
            assert imm is not None
            assert condition is not None
            assert pc is not None

            self._init_s_rrici(dc, ra, imm, condition, pc)
        elif self._suffix == Suffix.S_RRIF or self._suffix == Suffix.U_RRIF:
            assert dc is not None
            assert ra is not None
            assert imm is not None
            assert condition is not None

            self._init_s_rrif(dc, ra, imm, condition)
        elif self._suffix == Suffix.S_RRR or self._suffix == Suffix.U_RRR:
            assert dc is not None
            assert ra is not None
            assert rb is not None

            self._init_s_rrr(dc, ra, rb)
        elif self._suffix == Suffix.S_RRRC or self._suffix == Suffix.U_RRRC:
            assert dc is not None
            assert ra is not None
            assert rb is not None
            assert condition is not None

            self._init_s_rrrc(dc, ra, rb, condition)
        elif self._suffix == Suffix.S_RRRCI or self._suffix == Suffix.U_RRRCI:
            assert dc is not None
            assert ra is not None
            assert rb is not None
            assert condition is not None
            assert pc is not None

            self._init_s_rrrci(dc, ra, rb, condition, pc)
        elif self._suffix == Suffix.RR:
            assert isinstance(rc, GPRegister)
            assert ra is not None

            self._init_rr(rc, ra)
        elif self._suffix == Suffix.RRC:
            assert isinstance(rc, GPRegister)
            assert ra is not None
            assert condition is not None

            self._init_rrc(rc, ra, condition)
        elif self._suffix == Suffix.RRCI:
            assert isinstance(rc, GPRegister)
            assert ra is not None
            assert condition is not None
            assert pc is not None

            self._init_rrci(rc, ra, condition, pc)
        elif self._suffix == Suffix.ZR:
            assert ra is not None

            self._init_zr(ra)
        elif self._suffix == Suffix.ZRC:
            assert ra is not None
            assert condition is not None

            self._init_zrc(ra, condition)
        elif self._suffix == Suffix.ZRCI:
            assert ra is not None
            assert condition is not None
            assert pc is not None

            self._init_zrci(ra, condition, pc)
        elif self._suffix == Suffix.S_RR or self._suffix == Suffix.U_RR:
            assert dc is not None
            assert ra is not None

            self._init_s_rr(dc, ra)
        elif self._suffix == Suffix.S_RRC or self._suffix == Suffix.U_RRC:
            assert dc is not None
            assert ra is not None
            assert condition is not None

            self._init_s_rrc(dc, ra, condition)
        elif self._suffix == Suffix.S_RRCI or self._suffix == Suffix.U_RRCI:
            assert dc is not None
            assert ra is not None
            assert condition is not None
            assert pc is not None

            self._init_s_rrci(dc, ra, condition, pc)
        elif self._suffix == Suffix.DRDICI:
            assert dc is not None
            assert ra is not None
            assert db is not None
            assert imm is not None
            assert condition is not None
            assert pc is not None

            self._init_drdici(dc, ra, db, imm, condition, pc)
        elif self._suffix == Suffix.RRRI:
            assert isinstance(rc, GPRegister)
            assert ra is not None
            assert rb is not None
            assert imm is not None

            self._init_rrri(rc, ra, rb, imm)
        elif self._suffix == Suffix.RRRICI:
            assert isinstance(rc, GPRegister)
            assert ra is not None
            assert rb is not None
            assert imm is not None
            assert condition is not None
            assert pc is not None

            self._init_rrrici(rc, ra, rb, imm, condition, pc)
        elif self._suffix == Suffix.ZRRI:
            assert ra is not None
            assert rb is not None
            assert imm is not None

            self._init_zrri(ra, rb, imm)
        elif self._suffix == Suffix.ZRRICI:
            assert ra is not None
            assert rb is not None
            assert imm is not None
            assert condition is not None
            assert pc is not None

            self._init_zrrici(ra, rb, imm, condition, pc)
        elif self._suffix == Suffix.S_RRRI or self._suffix == Suffix.U_RRRI:
            assert dc is not None
            assert ra is not None
            assert rb is not None
            assert imm is not None

            self._init_s_rrri(dc, ra, rb, imm)
        elif self._suffix == Suffix.S_RRRICI or self._suffix == Suffix.U_RRRICI:
            assert dc is not None
            assert ra is not None
            assert rb is not None
            assert imm is not None
            assert condition is not None
            assert pc is not None

            self._init_s_rrrici(dc, ra, rb, imm, condition, pc)
        elif self._suffix == Suffix.RIR:
            assert isinstance(rc, GPRegister)
            assert imm is not None
            assert ra is not None

            self._init_rir(rc, imm, ra)
        elif self._suffix == Suffix.RIRC:
            assert isinstance(rc, GPRegister)
            assert imm is not None
            assert ra is not None
            assert condition is not None

            self._init_rirc(rc, imm, ra, condition)
        elif self._suffix == Suffix.RIRCI:
            assert isinstance(rc, GPRegister)
            assert imm is not None
            assert ra is not None
            assert condition is not None
            assert pc is not None

            self._init_rirci(rc, imm, ra, condition, pc)
        elif self._suffix == Suffix.ZIR:
            assert imm is not None
            assert ra is not None

            self._init_zir(imm, ra)
        elif self._suffix == Suffix.ZIRC:
            assert imm is not None
            assert ra is not None
            assert condition is not None

            self._init_zirc(imm, ra, condition)
        elif self._suffix == Suffix.ZIRCI:
            assert imm is not None
            assert ra is not None
            assert condition is not None
            assert pc is not None

            self._init_zirci(imm, ra, condition, pc)
        elif self._suffix == Suffix.S_RIRC or self._suffix == Suffix.U_RIRC:
            assert dc is not None
            assert imm is not None
            assert ra is not None
            assert condition is not None

            self._init_s_rirc(dc, imm, ra, condition)
        elif self._suffix == Suffix.S_RIRCI or self._suffix == Suffix.U_RIRCI:
            assert dc is not None
            assert imm is not None
            assert ra is not None
            assert condition is not None
            assert pc is not None

            self._init_s_rirci(dc, imm, ra, condition, pc)
        elif self._suffix == Suffix.R:
            assert isinstance(rc, GPRegister)

            self._init_r(rc)
        elif self._suffix == Suffix.RCI:
            assert isinstance(rc, GPRegister)
            assert condition is not None
            assert pc is not None

            self._init_rci(rc, condition, pc)
        elif self._suffix == Suffix.Z:
            self._init_z()
        elif self._suffix == Suffix.ZCI:
            assert condition is not None
            assert pc is not None

            self._init_zci(condition, pc)
        elif self._suffix == Suffix.S_R or self._suffix == Suffix.U_R:
            assert dc is not None

            self._init_s_r(dc)
        elif self._suffix == Suffix.S_RCI or self._suffix == Suffix.U_RCI:
            assert dc is not None
            assert condition is not None
            assert pc is not None

            self._init_s_rci(dc, condition, pc)
        elif self._suffix == Suffix.CI:
            assert condition is not None
            assert pc is not None

            self._init_ci(condition, pc)
        elif self._suffix == Suffix.I:
            assert imm is not None

            self._init_i(imm)
        elif self._suffix == Suffix.DDCI:
            assert dc is not None
            assert db is not None
            assert condition is not None
            assert pc is not None

            self._init_ddci(dc, db, condition, pc)
        elif self._suffix == Suffix.ERRI:
            assert endian is not None
            assert isinstance(rc, GPRegister)
            assert ra is not None
            assert off is not None

            self._init_erri(endian, rc, ra, off)
        elif self._suffix == Suffix.S_ERRI or self._suffix == Suffix.U_ERRI:
            assert endian is not None
            assert dc is not None
            assert ra is not None
            assert off is not None

            self._init_s_erri(endian, dc, ra, off)
        elif self._suffix == Suffix.EDRI:
            assert endian is not None
            assert dc is not None
            assert ra is not None
            assert off is not None

            self._init_edri(endian, dc, ra, off)
        elif self._suffix == Suffix.ERII:
            assert endian is not None
            assert ra is not None
            assert off is not None
            assert imm is not None

            self._init_erii(endian, ra, off, imm)
        elif self._suffix == Suffix.ERIR:
            assert endian is not None
            assert ra is not None
            assert off is not None
            assert rb is not None

            self._init_erir(endian, ra, off, rb)
        elif self._suffix == Suffix.ERID:
            assert endian is not None
            assert ra is not None
            assert off is not None
            assert db is not None

            self._init_erid(endian, ra, off, db)
        elif self._suffix == Suffix.DMA_RRI:
            assert ra is not None
            assert rb is not None
            assert imm is not None

            self._init_dma_rri(ra, rb, imm)
        else:
            raise ValueError

    def _init_rici(self, ra: SourceRegister, imm: int, condition: Condition, pc: int) -> None:
        assert self._op_code in Instruction.RICIOpCodes
        assert self._suffix == Suffix.RICI

        self._ra = ra
        self._imm = Immediate(Representation.SIGNED, 16, imm)

        if self._op_code in Instruction.AcquireRICIOpCodes:
            self._condition = AcquireCC(condition).condition()
        elif self._op_code in Instruction.ReleaseRICIOpCodes:
            self._condition = ReleaseCC(condition).condition()
        elif self._op_code in Instruction.BootRICIOpCodes:
            self._condition = BootCC(condition).condition()
        else:
            raise ValueError

        self._pc = Immediate(Representation.UNSIGNED, ConfigLoader.iram_address_width(), pc)

    def _init_rri(self, rc: GPRegister, ra: SourceRegister, imm: int) -> None:
        assert self._op_code in Instruction.RRIOpCodes
        assert self._suffix == Suffix.RRI

        self._rc = rc
        self._ra = ra

        if self._op_code in Instruction.AddRRIOpCodes:
            self._imm = Immediate(Representation.UNSIGNED, 32, imm)
        elif self._op_code in Instruction.AsrRRIOpCodes:
            self._imm = Immediate(Representation.UNSIGNED, 5, imm)
        elif self._op_code in Instruction.CallRRIOpCodes:
            self._imm = Immediate(Representation.SIGNED, 24, imm)
        else:
            raise ValueError

    def _init_rric(self, rc: GPRegister, ra: SourceRegister, imm: int, condition: Condition,) -> None:
        assert self._op_code in Instruction.RRICOpCodes
        assert self._suffix == Suffix.RRIC

        self._rc = rc
        self._ra = ra

        if self._op_code in Instruction.AddRRICOpCodes or self._op_code in Instruction.SubRRICOpCodes:
            self._imm = Immediate(Representation.SIGNED, 24, imm)
        elif self._op_code in Instruction.AsrRRICOpCodes:
            self._imm = Immediate(Representation.UNSIGNED, 5, imm)
        else:
            raise ValueError

        if self._op_code in Instruction.AddRRICOpCodes or self._op_code in Instruction.AsrRRICOpCodes:
            self._condition = LogSetCC(condition).condition()
        elif self._op_code in Instruction.SubRRICOpCodes:
            self._condition = ExtSubSetCC(condition).condition()
        else:
            raise ValueError

    def _init_rrici(self, rc: GPRegister, ra: SourceRegister, imm: int, condition: Condition, pc: int,) -> None:
        assert self._op_code in Instruction.RRICIOpCodes
        assert self._suffix == Suffix.RRICI

        self._rc = rc
        self._ra = ra

        if (
            self._op_code in Instruction.AddRRICIOpCodes
            or self._op_code in Instruction.AndRRICIOpCodes
            or self._op_code in Instruction.SubRRICIOpCodes
        ):
            self._imm = Immediate(Representation.SIGNED, 8, imm)
        elif self._op_code in Instruction.AsrRRICIOpCodes:
            self._imm = Immediate(Representation.UNSIGNED, 5, imm)
        else:
            raise ValueError

        if self._op_code in Instruction.AddRRICIOpCodes:
            self._condition = AddNZCC(condition).condition()
        elif self._op_code in Instruction.AndRRICIOpCodes:
            self._condition = LogNZCC(condition).condition()
        elif self._op_code in Instruction.AsrRRICIOpCodes:
            self._condition = ImmShiftNZCC(condition).condition()
        elif self._op_code in Instruction.SubRRICIOpCodes:
            self._condition = SubNZCC(condition).condition()
        else:
            raise ValueError

        self._pc = Immediate(Representation.UNSIGNED, ConfigLoader.iram_address_width(), pc)

    def _init_rrif(self, rc: GPRegister, ra: SourceRegister, imm: int, condition: Condition,) -> None:
        assert self._op_code in Instruction.RRIFOpCodes
        assert self._suffix == Suffix.RRIF

        self._rc = rc
        self._ra = ra
        self._imm = Immediate(Representation.SIGNED, 24, imm)
        self._condition = FalseCC(condition).condition()

    def _init_rrr(self, rc: GPRegister, ra: SourceRegister, rb: SourceRegister) -> None:
        assert self._op_code in Instruction.RRROpCodes
        assert self._suffix == Suffix.RRR

        self._rc = rc
        self._ra = ra
        self._rb = rb

    def _init_rrrc(self, rc: GPRegister, ra: SourceRegister, rb: SourceRegister, condition: Condition,) -> None:
        assert self._op_code in Instruction.RRRCOpCodes
        assert self._suffix == Suffix.RRRC

        self._rc = rc
        self._ra = ra
        self._rb = rb

        if self._op_code in Instruction.AddRRRCOpCodes:
            self._condition = LogSetCC(condition).condition()
        elif self._op_code in Instruction.RsubRRRCOpCodes:
            self._condition = SubSetCC(condition).condition()
        elif self._op_code in Instruction.SubRRICIOpCodes:
            self._condition = ExtSubSetCC(condition).condition()
        else:
            raise ValueError

    def _init_rrrci(
        self, rc: GPRegister, ra: SourceRegister, rb: SourceRegister, condition: Condition, pc: int,
    ) -> None:
        assert self._op_code in Instruction.RRRCIOpCodes
        assert self._suffix == Suffix.RRRCI

        self._rc = rc
        self._ra = ra
        self._rb = rb

        if self._op_code in Instruction.AddRRRCIOpCodes:
            self._condition = AddNZCC(condition).condition()
        elif self._op_code in Instruction.AndRRRCIOpCodes:
            self._condition = LogNZCC(condition).condition()
        elif self._op_code in Instruction.AsrRRRCIOpCodes:
            self._condition = ShiftNZCC(condition).condition()
        elif self._op_code in Instruction.MulRRRCIOpCodes:
            self._condition = MulNZCC(condition).condition()
        elif self._op_code in Instruction.RsubRRRCIOpCodes:
            self._condition = SubNZCC(condition).condition()
        else:
            raise ValueError

        self._pc = Immediate(Representation.UNSIGNED, ConfigLoader.iram_address_width(), pc)

    def _init_zri(self, ra: SourceRegister, imm: int):
        assert self._op_code in Instruction.RRIOpCodes
        assert self._suffix == Suffix.ZRI

        self._ra = ra

        if self._op_code in Instruction.AddRRIOpCodes:
            self._imm = Immediate(Representation.UNSIGNED, 32, imm)
        elif self._op_code in Instruction.AsrRRIOpCodes:
            self._imm = Immediate(Representation.UNSIGNED, 5, imm)
        elif self._op_code in Instruction.CallRRIOpCodes:
            self._imm = Immediate(Representation.SIGNED, 28, imm)
        else:
            raise ValueError

    def _init_zric(self, ra: SourceRegister, imm: int, condition: Condition) -> None:
        assert self._op_code in Instruction.RRICOpCodes
        assert self._suffix == Suffix.ZRIC

        self._ra = ra

        if self._op_code in Instruction.AddRRICOpCodes or self._op_code in Instruction.SubRRICOpCodes:
            self._imm = Immediate(Representation.SIGNED, 27, imm)
        elif self._op_code in Instruction.AsrRRICOpCodes:
            self._imm = Immediate(Representation.UNSIGNED, 5, imm)
        else:
            raise ValueError

        if self._op_code in Instruction.AddRRICOpCodes or self._op_code in Instruction.AsrRRICOpCodes:
            self._condition = LogSetCC(condition).condition()
        elif self._op_code in Instruction.SubRRICOpCodes:
            self._condition = ExtSubSetCC(condition).condition()
        else:
            raise ValueError

    def _init_zrici(self, ra: SourceRegister, imm: int, condition: Condition, pc: int) -> None:
        assert self._op_code in Instruction.RRICIOpCodes
        assert self._suffix == Suffix.ZRICI

        self._ra = ra

        if (
            self._op_code in Instruction.AddRRICIOpCodes
            or self._op_code in Instruction.AndRRICIOpCodes
            or self._op_code in Instruction.SubRRICIOpCodes
        ):
            self._imm = Immediate(Representation.SIGNED, 11, imm)
        elif self._op_code in Instruction.AsrRRICIOpCodes:
            self._imm = Immediate(Representation.UNSIGNED, 5, imm)
        else:
            raise ValueError

        if self._op_code in Instruction.AddRRICIOpCodes:
            self._condition = AddNZCC(condition).condition()
        elif self._op_code in Instruction.AndRRICIOpCodes:
            self._condition = LogNZCC(condition).condition()
        elif self._op_code in Instruction.AsrRRICIOpCodes:
            self._condition = ImmShiftNZCC(condition).condition()
        elif self._op_code in Instruction.SubRRICIOpCodes:
            self._condition = SubNZCC(condition).condition()
        else:
            raise ValueError

        self._pc = Immediate(Representation.UNSIGNED, ConfigLoader.iram_address_width(), pc)

    def _init_zrif(self, ra: SourceRegister, imm: int, condition: Condition) -> None:
        assert self._op_code in Instruction.RRIFOpCodes
        assert self._suffix == Suffix.ZRIF

        self._ra = ra
        self._imm = Immediate(Representation.SIGNED, 27, imm)
        self._condition = FalseCC(condition).condition()

    def _init_zrr(self, ra: SourceRegister, rb: SourceRegister) -> None:
        assert self._op_code in Instruction.RRROpCodes
        assert self._suffix == Suffix.ZRR

        self._ra = ra
        self._rb = rb

    def _init_zrrc(self, ra: SourceRegister, rb: SourceRegister, condition: Condition) -> None:
        assert self._op_code in Instruction.RRRCOpCodes
        assert self._suffix == Suffix.ZRRC

        self._ra = ra
        self._rb = rb

        if self._op_code in Instruction.AddRRRCOpCodes:
            self._condition = LogSetCC(condition).condition()
        elif self._op_code in Instruction.RsubRRRCOpCodes:
            self._condition = SubSetCC(condition).condition()
        elif self._op_code in Instruction.SubRRICIOpCodes:
            self._condition = ExtSubSetCC(condition).condition()
        else:
            raise ValueError

    def _init_zrrci(self, ra: SourceRegister, rb: SourceRegister, condition: Condition, pc: int) -> None:
        assert self._op_code in Instruction.RRRCIOpCodes
        assert self._suffix == Suffix.ZRRCI

        self._ra = ra
        self._rb = rb

        if self._op_code in Instruction.AddRRRCIOpCodes:
            self._condition = AddNZCC(condition).condition()
        elif self._op_code in Instruction.AndRRRCIOpCodes:
            self._condition = LogNZCC(condition).condition()
        elif self._op_code in Instruction.AsrRRRCIOpCodes:
            self._condition = ShiftNZCC(condition).condition()
        elif self._op_code in Instruction.MulRRRCIOpCodes:
            self._condition = MulNZCC(condition).condition()
        elif self._op_code in Instruction.RsubRRRCIOpCodes:
            self._condition = SubNZCC(condition).condition()
        else:
            raise ValueError

        self._pc = Immediate(Representation.UNSIGNED, ConfigLoader.iram_address_width(), pc)

    def _init_s_rri(self, dc: PairRegister, ra: SourceRegister, imm: int):
        assert self._op_code in Instruction.RRIOpCodes
        assert self._suffix == Suffix.S_RRI or self._suffix == Suffix.U_RRI

        self._dc = dc
        self._ra = ra

        if self._op_code in Instruction.AddRRIOpCodes:
            self._imm = Immediate(Representation.UNSIGNED, 32, imm)
        elif self._op_code in Instruction.AsrRRIOpCodes:
            self._imm = Immediate(Representation.UNSIGNED, 5, imm)
        elif self._op_code in Instruction.CallRRIOpCodes:
            self._imm = Immediate(Representation.SIGNED, 24, imm)
        else:
            raise ValueError

    def _init_s_rric(self, dc: PairRegister, ra: SourceRegister, imm: int, condition: Condition) -> None:
        assert self._op_code in Instruction.RRICOpCodes
        assert self._suffix == Suffix.S_RRIC or self._suffix == Suffix.U_RRIC

        self._dc = dc
        self._ra = ra

        if self._op_code in Instruction.AddRRICOpCodes or self._op_code in Instruction.SubRRICOpCodes:
            self._imm = Immediate(Representation.SIGNED, 24, imm)
        elif self._op_code in Instruction.AsrRRICOpCodes:
            self._imm = Immediate(Representation.UNSIGNED, 5, imm)
        else:
            raise ValueError

        if self._op_code in Instruction.AddRRICOpCodes or self._op_code in Instruction.AsrRRICOpCodes:
            self._condition = LogSetCC(condition).condition()
        elif self._op_code in Instruction.SubRRICOpCodes:
            self._condition = ExtSubSetCC(condition).condition()
        else:
            raise ValueError

    def _init_s_rrici(self, dc: PairRegister, ra: SourceRegister, imm: int, condition: Condition, pc: int,) -> None:
        assert self._op_code in Instruction.RRICIOpCodes
        assert self._suffix == Suffix.S_RRICI or self._suffix == Suffix.U_RRICI

        self._dc = dc
        self._ra = ra

        if (
            self._op_code in Instruction.AddRRICIOpCodes
            or self._op_code in Instruction.AndRRICIOpCodes
            or self._op_code in Instruction.SubRRICIOpCodes
        ):
            self._imm = Immediate(Representation.SIGNED, 8, imm)
        elif self._op_code in Instruction.AsrRRICIOpCodes:
            self._imm = Immediate(Representation.UNSIGNED, 5, imm)
        else:
            raise ValueError

        if self._op_code in Instruction.AddRRICIOpCodes:
            self._condition = AddNZCC(condition).condition()
        elif self._op_code in Instruction.AndRRICIOpCodes:
            self._condition = LogNZCC(condition).condition()
        elif self._op_code in Instruction.AsrRRICIOpCodes:
            self._condition = ImmShiftNZCC(condition).condition()
        elif self._op_code in Instruction.SubRRICIOpCodes:
            self._condition = SubNZCC(condition).condition()
        else:
            raise ValueError

        self._pc = Immediate(Representation.UNSIGNED, ConfigLoader.iram_address_width(), pc)

    def _init_s_rrif(self, dc: PairRegister, ra: SourceRegister, imm: int, condition: Condition) -> None:
        assert self._op_code in Instruction.RRIFOpCodes
        assert self._suffix == Suffix.S_RRIF or self._suffix == Suffix.U_RRIF

        self._dc = dc
        self._ra = ra
        self._imm = Immediate(Representation.SIGNED, 24, imm)
        self._condition = FalseCC(condition).condition()

    def _init_s_rrr(self, dc: PairRegister, ra: SourceRegister, rb: SourceRegister) -> None:
        assert self._op_code in Instruction.RRROpCodes
        assert self._suffix == Suffix.S_RRR or self._suffix == Suffix.U_RRR

        self._dc = dc
        self._ra = ra
        self._rb = rb

    def _init_s_rrrc(self, dc: PairRegister, ra: SourceRegister, rb: SourceRegister, condition: Condition,) -> None:
        assert self._op_code in Instruction.RRRCOpCodes
        assert self._suffix == Suffix.S_RRRC or self._suffix == Suffix.U_RRRC

        self._dc = dc
        self._ra = ra
        self._rb = rb

        if self._op_code in Instruction.AddRRRCOpCodes:
            self._condition = LogSetCC(condition).condition()
        elif self._op_code in Instruction.RsubRRRCOpCodes:
            self._condition = SubSetCC(condition).condition()
        elif self._op_code in Instruction.SubRRICIOpCodes:
            self._condition = ExtSubSetCC(condition).condition()
        else:
            raise ValueError

    def _init_s_rrrci(
        self, dc: PairRegister, ra: SourceRegister, rb: SourceRegister, condition: Condition, pc: int,
    ) -> None:
        assert self._op_code in Instruction.RRRCIOpCodes
        assert self._suffix == Suffix.S_RRRCI or self._suffix == Suffix.U_RRRCI

        self._dc = dc
        self._ra = ra
        self._rb = rb

        if self._op_code in Instruction.AddRRRCIOpCodes:
            self._condition = AddNZCC(condition).condition()
        elif self._op_code in Instruction.AndRRRCIOpCodes:
            self._condition = LogNZCC(condition).condition()
        elif self._op_code in Instruction.AsrRRRCIOpCodes:
            self._condition = ShiftNZCC(condition).condition()
        elif self._op_code in Instruction.MulRRRCIOpCodes:
            self._condition = MulNZCC(condition).condition()
        elif self._op_code in Instruction.RsubRRRCIOpCodes:
            self._condition = SubNZCC(condition).condition()
        else:
            raise ValueError

        self._pc = Immediate(Representation.UNSIGNED, ConfigLoader.iram_address_width(), pc)

    def _init_rr(self, rc: GPRegister, ra: SourceRegister) -> None:
        assert self._op_code in Instruction.RROpCodes
        assert self._suffix == Suffix.RR

        self._rc = rc
        self._ra = ra

    def _init_rrc(self, rc: GPRegister, ra: SourceRegister, condition: Condition) -> None:
        assert self._op_code in Instruction.RRCOpCodes
        assert self._suffix == Suffix.RRC

        self._rc = rc
        self._ra = ra
        self._condition = LogSetCC(condition).condition()

    def _init_rrci(self, rc: GPRegister, ra: SourceRegister, condition: Condition, pc: int) -> None:
        assert self._op_code in Instruction.RRCIOpCodes
        assert self._suffix == Suffix.RRCI

        self._rc = rc
        self._ra = ra

        if self._op_code in Instruction.CaoRRCIOpCodes:
            self._condition = CountNZCC(condition).condition()
        elif self._op_code in Instruction.ExtsbRRCIOpCodes:
            self._condition = LogNZCC(condition).condition()
        elif self._op_code in Instruction.TimeCfgRRCIOpCodes:
            self._condition = TrueCC(condition).condition()
        else:
            raise ValueError

        self._pc = Immediate(Representation.UNSIGNED, ConfigLoader.iram_address_width(), pc)

    def _init_zr(self, ra: SourceRegister) -> None:
        assert self._op_code in Instruction.RROpCodes
        assert self._suffix == Suffix.ZR

        self._ra = ra

    def _init_zrc(self, ra: SourceRegister, condition: Condition) -> None:
        assert self._op_code in Instruction.RRCOpCodes
        assert self._suffix == Suffix.ZRC

        self._ra = ra
        self._condition = LogSetCC(condition).condition()

    def _init_zrci(self, ra: SourceRegister, condition: Condition, pc: int) -> None:
        assert self._op_code in Instruction.RRCIOpCodes
        assert self._suffix == Suffix.ZRCI

        self._ra = ra

        if self._op_code in Instruction.CaoRRCIOpCodes:
            self._condition = CountNZCC(condition).condition()
        elif self._op_code in Instruction.ExtsbRRCIOpCodes:
            self._condition = LogNZCC(condition).condition()
        elif self._op_code in Instruction.TimeCfgRRCIOpCodes:
            self._condition = TrueCC(condition).condition()
        else:
            raise ValueError

        self._pc = Immediate(Representation.UNSIGNED, ConfigLoader.iram_address_width(), pc)

    def _init_s_rr(self, dc: PairRegister, ra: SourceRegister) -> None:
        assert self._op_code in Instruction.RROpCodes
        assert self._suffix == Suffix.S_RR or self._suffix == Suffix.U_RR

        self._dc = dc
        self._ra = ra

    def _init_s_rrc(self, dc: PairRegister, ra: SourceRegister, condition: Condition) -> None:
        assert self._op_code in Instruction.RRCOpCodes
        assert self._suffix == Suffix.S_RRC or self._suffix == Suffix.U_RRC

        self._dc = dc
        self._ra = ra
        self._condition = LogSetCC(condition).condition()

    def _init_s_rrci(self, dc: PairRegister, ra: SourceRegister, condition: Condition, pc: int) -> None:
        assert self._op_code in Instruction.RRCIOpCodes
        assert self._suffix == Suffix.S_RRCI or self._suffix == Suffix.U_RRCI

        self._dc = dc
        self._ra = ra

        if self._op_code in Instruction.CaoRRCIOpCodes:
            self._condition = CountNZCC(condition).condition()
        elif self._op_code in Instruction.ExtsbRRCIOpCodes:
            self._condition = LogNZCC(condition).condition()
        elif self._op_code in Instruction.TimeCfgRRCIOpCodes:
            self._condition = TrueCC(condition).condition()
        else:
            raise ValueError

        self._pc = Immediate(Representation.UNSIGNED, ConfigLoader.iram_address_width(), pc)

    def _init_drdici(
        self, dc: PairRegister, ra: SourceRegister, db: PairRegister, imm: int, condition: Condition, pc: int,
    ) -> None:
        assert self._op_code in Instruction.DRDICIOpCodes
        assert self._suffix == Suffix.DRDICI

        self._dc = dc
        self._ra = ra
        self._db = db
        self._imm = Immediate(Representation.UNSIGNED, 5, imm)

        if self._op_code in Instruction.DivStepDRDICIOpCodes:
            self._condition = DivCC(condition).condition()
        elif self._op_code in Instruction.MulStepDRDICIOpCodes:
            self._condition = BootCC(condition).condition()
        else:
            raise ValueError

        self._pc = Immediate(Representation.UNSIGNED, ConfigLoader.iram_address_width(), pc)

    def _init_rrri(self, rc: GPRegister, ra: SourceRegister, rb: SourceRegister, imm: int) -> None:
        assert self._op_code in Instruction.RRRIOpCodes
        assert self._suffix == Suffix.RRRI

        self._rc = rc
        self._ra = ra
        self._rb = rb
        self._imm = Immediate(Representation.UNSIGNED, 5, imm)

    def _init_rrrici(
        self, rc: GPRegister, ra: SourceRegister, rb: SourceRegister, imm: int, condition: Condition, pc: int,
    ) -> None:
        assert self._op_code in Instruction.RRRICIOpCodes
        assert self._suffix == Suffix.RRRICI

        self._rc = rc
        self._ra = ra
        self._rb = rb
        self._imm = Immediate(Representation.UNSIGNED, 5, imm)
        self._condition = DivNZCC(condition).condition()
        self._pc = Immediate(Representation.UNSIGNED, ConfigLoader.iram_address_width(), pc)

    def _init_zrri(self, ra: SourceRegister, rb: SourceRegister, imm: int) -> None:
        assert self._op_code in Instruction.RRRIOpCodes
        assert self._suffix == Suffix.ZRRI

        self._ra = ra
        self._rb = rb
        self._imm = Immediate(Representation.UNSIGNED, 5, imm)

    def _init_zrrici(self, ra: SourceRegister, rb: SourceRegister, imm: int, condition: Condition, pc: int,) -> None:
        assert self._op_code in Instruction.RRRICIOpCodes
        assert self._suffix == Suffix.ZRRICI

        self._ra = ra
        self._rb = rb
        self._imm = Immediate(Representation.UNSIGNED, 5, imm)
        self._condition = DivNZCC(condition).condition()
        self._pc = Immediate(Representation.UNSIGNED, ConfigLoader.iram_address_width(), pc)

    def _init_s_rrri(self, dc: PairRegister, ra: SourceRegister, rb: SourceRegister, imm: int) -> None:
        assert self._op_code in Instruction.RRRIOpCodes
        assert self._suffix == Suffix.S_RRRI or self._suffix == Suffix.U_RRRI

        self._dc = dc
        self._ra = ra
        self._rb = rb
        self._imm = Immediate(Representation.UNSIGNED, 5, imm)

    def _init_s_rrrici(
        self, dc: PairRegister, ra: SourceRegister, rb: SourceRegister, imm: int, condition: Condition, pc: int,
    ) -> None:
        assert self._op_code in Instruction.RRRICIOpCodes
        assert self._suffix == Suffix.S_RRRICI or self._suffix == Suffix.U_RRRICI

        self._dc = dc
        self._ra = ra
        self._rb = rb
        self._imm = Immediate(Representation.UNSIGNED, 5, imm)
        self._condition = DivNZCC(condition).condition()
        self._pc = Immediate(Representation.UNSIGNED, ConfigLoader.iram_address_width(), pc)

    def _init_rir(self, rc: GPRegister, imm: int, ra: SourceRegister) -> None:
        assert self._op_code in Instruction.RIROpCodes
        assert self._suffix == Suffix.RIR

        self._rc = rc
        self._imm = Immediate(Representation.UNSIGNED, 32, imm)
        self._ra = ra

    def _init_rirc(self, rc: GPRegister, imm: int, ra: SourceRegister, condition: Condition):
        assert self._op_code in Instruction.RIRCOpCodes
        assert self._suffix == Suffix.RIRC

        self._rc = rc
        self._imm = Immediate(Representation.SIGNED, 24, imm)
        self._ra = ra
        self._condition = SubSetCC(condition).condition()

    def _init_rirci(self, rc: GPRegister, imm: int, ra: SourceRegister, condition: Condition, pc: int,) -> None:
        assert self._op_code in Instruction.RIRCIOpCodes
        assert self._suffix == Suffix.RIRCI

        self._rc = rc
        self._imm = Immediate(Representation.SIGNED, 8, imm)
        self._ra = ra
        self._condition = SubNZCC(condition).condition()
        self._pc = Immediate(Representation.UNSIGNED, ConfigLoader.iram_address_width(), pc)

    def _init_zir(self, imm: int, ra: SourceRegister) -> None:
        assert self._op_code in Instruction.RIROpCodes
        assert self._suffix == Suffix.ZIR

        self._imm = Immediate(Representation.UNSIGNED, 32, imm)
        self._ra = ra

    def _init_zirc(self, imm: int, ra: SourceRegister, condition: Condition) -> None:
        assert self._op_code in Instruction.RIRCOpCodes
        assert self._suffix == Suffix.ZIRC

        self._imm = Immediate(Representation.SIGNED, 27, imm)
        self._ra = ra
        self._condition = SubSetCC(condition).condition()

    def _init_zirci(self, imm: int, ra: SourceRegister, condition: Condition, pc: int) -> None:
        assert self._op_code in Instruction.RIRCIOpCodes
        assert self._suffix == Suffix.ZIRCI

        self._imm = Immediate(Representation.SIGNED, 11, imm)
        self._ra = ra
        self._condition = SubNZCC(condition).condition()
        self._pc = Immediate(Representation.UNSIGNED, ConfigLoader.iram_address_width(), pc)

    def _init_s_rirc(self, dc: PairRegister, imm: int, ra: SourceRegister, condition: Condition) -> None:
        assert self._op_code in Instruction.RIRCOpCodes
        assert self._suffix == Suffix.S_RIRC or self._suffix == Suffix.U_RIRC

        self._dc = dc
        self._imm = Immediate(Representation.SIGNED, 24, imm)
        self._ra = ra
        self._condition = SubSetCC(condition).condition()

    def _init_s_rirci(self, dc: PairRegister, imm: int, ra: SourceRegister, condition: Condition, pc: int,) -> None:
        assert self._op_code in Instruction.RIRCIOpCodes
        assert self._suffix == Suffix.S_RIRCI or self._suffix == Suffix.U_RIRCI

        self._dc = dc
        self._imm = Immediate(Representation.SIGNED, 8, imm)
        self._ra = ra
        self._condition = SubNZCC(condition).condition()
        self._pc = Immediate(Representation.UNSIGNED, ConfigLoader.iram_address_width(), pc)

    def _init_r(self, rc: GPRegister) -> None:
        assert self._op_code in Instruction.ROpCodes
        assert self._suffix == Suffix.R

        self._rc = rc

    def _init_rci(self, rc: GPRegister, condition: Condition, pc: int) -> None:
        assert self._op_code in Instruction.RCIOpCodes
        assert self._suffix == Suffix.RCI

        self._rc = rc
        self._condition = TrueCC(condition).condition()
        self._pc = Immediate(Representation.UNSIGNED, ConfigLoader.iram_address_width(), pc)

    def _init_z(self) -> None:
        assert self._op_code in Instruction.ROpCodes or self._op_code == OpCode.NOP
        assert self._suffix == Suffix.Z

    def _init_zci(self, condition: Condition, pc: int) -> None:
        assert self._op_code in Instruction.RCIOpCodes
        assert self._suffix == Suffix.ZCI

        self._condition = TrueCC(condition).condition()
        self._pc = Immediate(Representation.UNSIGNED, ConfigLoader.iram_address_width(), pc)

    def _init_s_r(self, dc: PairRegister) -> None:
        assert self._op_code in Instruction.ROpCodes
        assert self._suffix == Suffix.S_R or self._suffix == Suffix.U_R

        self._dc = dc

    def _init_s_rci(self, dc: PairRegister, condition: Condition, pc: int) -> None:
        assert self._op_code in Instruction.RCIOpCodes
        assert self._suffix == Suffix.S_RCI or self._suffix == Suffix.U_RCI

        self._dc = dc
        self._condition = TrueCC(condition).condition()
        self._pc = Immediate(Representation.UNSIGNED, ConfigLoader.iram_address_width(), pc)

    def _init_ci(self, condition: Condition, pc: int) -> None:
        assert self._op_code in Instruction.CIOpCodes
        assert self._suffix == Suffix.CI

        self._condition = BootCC(condition).condition()
        self._pc = Immediate(Representation.UNSIGNED, ConfigLoader.iram_address_width(), pc)

    def _init_i(self, imm: int) -> None:
        assert self._op_code in Instruction.IOpCodes
        assert self._suffix == Suffix.I

        self._imm = Immediate(Representation.SIGNED, 24, imm)

    def _init_ddci(self, dc: PairRegister, db: PairRegister, condition: Condition, pc: int) -> None:
        assert self._op_code in Instruction.DDCIOpCodes
        assert self._suffix == Suffix.DDCI

        self._dc = dc
        self._db = db
        self._condition = TrueFalseCC(condition).condition()
        self._pc = Immediate(Representation.UNSIGNED, ConfigLoader.iram_address_width(), pc)

    def _init_erri(self, endian: Endian, rc: GPRegister, ra: SourceRegister, off: int) -> None:
        assert self._op_code in Instruction.ERRIOpCodes
        assert self._suffix == Suffix.ERRI

        self._endian = endian
        self._rc = rc
        self._ra = ra
        self._off = Immediate(Representation.SIGNED, 24, off)

    def _init_s_erri(self, endian: Endian, dc: PairRegister, ra: SourceRegister, off: int) -> None:
        assert self._op_code in Instruction.ERRIOpCodes
        assert self._suffix == Suffix.S_ERRI or self._suffix == Suffix.U_ERRI

        self._endian = endian
        self._dc = dc
        self._ra = ra
        self._off = Immediate(Representation.SIGNED, 24, off)

    def _init_edri(self, endian: Endian, dc: PairRegister, ra: SourceRegister, off: int) -> None:
        assert self._op_code in Instruction.EDRIOpCodes
        assert self._suffix == Suffix.EDRI

        self._endian = endian
        self._dc = dc
        self._ra = ra
        self._off = Immediate(Representation.SIGNED, 24, off)

    def _init_erii(self, endian: Endian, ra: SourceRegister, off: int, imm: int) -> None:
        assert self._op_code in Instruction.ERIIOpCodes
        assert self._suffix == Suffix.ERII

        self._endian = endian
        self._ra = ra
        # NOTE(bongjoon.hyun@gmail.com): original width is 12
        self._off = Immediate(Representation.SIGNED, 24, off)
        # NOTE(bongjoon.hyun@gmail.com): original width is 8
        self._imm = Immediate(Representation.SIGNED, 16, imm)

    def _init_erir(self, endian: Endian, ra: SourceRegister, off: int, rb: SourceRegister):
        assert self._op_code in Instruction.ERIROpCodes
        assert self._suffix == Suffix.ERIR

        self._endian = endian
        self._ra = ra
        self._off = Immediate(Representation.SIGNED, 24, off)
        self._rb = rb

    def _init_erid(self, endian: Endian, ra: SourceRegister, off: int, db: PairRegister):
        assert self._op_code in Instruction.ERIDOpCodes
        assert self._suffix == Suffix.ERID

        self._endian = endian
        self._ra = ra
        self._off = Immediate(Representation.SIGNED, 24, off)
        self._db = db

    def _init_dma_rri(self, ra: SourceRegister, rb: SourceRegister, imm: int):
        assert self._op_code in Instruction.DMARRIOpCodes
        assert self._suffix == Suffix.DMA_RRI

        self._ra = ra
        self._rb = rb
        self._imm = Immediate(Representation.UNSIGNED, 8, imm)
