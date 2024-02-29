from abi.isa.instruction.instruction import Instruction
from abi.isa.instruction.op_code import OpCode
from abi.isa.instruction.suffix import Suffix
from converter.condition_converter import ConditionConverter
from converter.endian_converter import EndianConverter
from converter.op_code_converter import OpCodeConverter
from converter.register_converter import RegisterConverter
from converter.suffix_converter import SuffixConverter


class InstructionConverter:
    def __init__(self):
        pass

    @staticmethod
    def convert_to_string(instruction: Instruction) -> str:
        op_code = OpCodeConverter.convert_to_string(instruction.op_code())
        suffix = SuffixConverter.convert_to_string(instruction.suffix())

        if instruction.suffix() == Suffix.RICI:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_rici_to_string(instruction)}"
        elif instruction.suffix() == Suffix.RRI:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_rri_to_string(instruction)}"
        elif instruction.suffix() == Suffix.RRIC:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_rric_to_string(instruction)}"
        elif instruction.suffix() == Suffix.RRICI:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_rrici_to_string(instruction)}"
        elif instruction.suffix() == Suffix.RRIF:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_rrif_to_string(instruction)}"
        elif instruction.suffix() == Suffix.RRR:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_rrr_to_string(instruction)}"
        elif instruction.suffix() == Suffix.RRRC:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_rrrc_to_string(instruction)}"
        elif instruction.suffix() == Suffix.RRRCI:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_rrrci_to_string(instruction)}"
        elif instruction.suffix() == Suffix.ZRI:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_zri_to_string(instruction)}"
        elif instruction.suffix() == Suffix.ZRIC:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_zric_to_string(instruction)}"
        elif instruction.suffix() == Suffix.ZRICI:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_zrici_to_string(instruction)}"
        elif instruction.suffix() == Suffix.ZRIF:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_zrif_to_string(instruction)}"
        elif instruction.suffix() == Suffix.ZRR:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_zrr_to_string(instruction)}"
        elif instruction.suffix() == Suffix.ZRRC:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_zrrc_to_string(instruction)}"
        elif instruction.suffix() == Suffix.ZRRCI:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_zrrci_to_string(instruction)}"
        elif instruction.suffix() == Suffix.S_RRI or instruction.suffix() == Suffix.U_RRI:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_s_rri_to_string(instruction)}"
        elif instruction.suffix() == Suffix.S_RRIC or instruction.suffix() == Suffix.U_RRIC:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_s_rric_to_string(instruction)}"
        elif instruction.suffix() == Suffix.S_RRICI or instruction.suffix() == Suffix.U_RRICI:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_s_rrici_to_string(instruction)}"
        elif instruction.suffix() == Suffix.S_RRIF or instruction.suffix() == Suffix.U_RRIF:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_s_rrif_to_string(instruction)}"
        elif instruction.suffix() == Suffix.S_RRR or instruction.suffix() == Suffix.U_RRR:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_s_rrr_to_string(instruction)}"
        elif instruction.suffix() == Suffix.S_RRRC or instruction.suffix() == Suffix.U_RRRC:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_s_rrrc_to_string(instruction)}"
        elif instruction.suffix() == Suffix.S_RRRCI or instruction.suffix() == Suffix.U_RRRCI:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_s_rrrci_to_string(instruction)}"
        elif instruction.suffix() == Suffix.RR:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_rr_to_string(instruction)}"
        elif instruction.suffix() == Suffix.RRC:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_rrc_to_string(instruction)}"
        elif instruction.suffix() == Suffix.RRCI:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_rrci_to_string(instruction)}"
        elif instruction.suffix() == Suffix.ZR:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_zr_to_string(instruction)}"
        elif instruction.suffix() == Suffix.ZRC:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_zrc_to_string(instruction)}"
        elif instruction.suffix() == Suffix.ZRCI:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_zrci_to_string(instruction)}"
        elif instruction.suffix() == Suffix.S_RR or instruction.suffix() == Suffix.U_RR:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_s_rr_to_string(instruction)}"
        elif instruction.suffix() == Suffix.S_RRC or instruction.suffix() == Suffix.U_RRC:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_s_rrc_to_string(instruction)}"
        elif instruction.suffix() == Suffix.S_RRCI or instruction.suffix() == Suffix.U_RRCI:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_s_rrci_to_string(instruction)}"
        elif instruction.suffix() == Suffix.DRDICI:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_drdici_to_string(instruction)}"
        elif instruction.suffix() == Suffix.RRRI:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_rrri_to_string(instruction)}"
        elif instruction.suffix() == Suffix.RRRICI:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_rrrici_to_string(instruction)}"
        elif instruction.suffix() == Suffix.ZRRI:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_zrri_to_string(instruction)}"
        elif instruction.suffix() == Suffix.ZRRICI:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_zrrici_to_string(instruction)}"
        elif instruction.suffix() == Suffix.S_RRRI or instruction.suffix() == Suffix.U_RRRI:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_s_rrri_to_string(instruction)}"
        elif instruction.suffix() == Suffix.S_RRRICI or instruction.suffix() == Suffix.U_RRRICI:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_s_rrrici_to_string(instruction)}"
        elif instruction.suffix() == Suffix.RIR:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_rir_to_string(instruction)}"
        elif instruction.suffix() == Suffix.RIRC:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_rirc_to_string(instruction)}"
        elif instruction.suffix() == Suffix.RIRCI:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_rirci_to_string(instruction)}"
        elif instruction.suffix() == Suffix.ZIR:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_zir_to_string(instruction)}"
        elif instruction.suffix() == Suffix.ZIRC:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_zirc_to_string(instruction)}"
        elif instruction.suffix() == Suffix.ZIRCI:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_zirci_to_string(instruction)}"
        elif instruction.suffix() == Suffix.S_RIRC or instruction.suffix() == Suffix.U_RIRC:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_s_rirc_to_string(instruction)}"
        elif instruction.suffix() == Suffix.S_RIRCI or instruction.suffix() == Suffix.U_RIRCI:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_s_rirci_to_string(instruction)}"
        elif instruction.suffix() == Suffix.R:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_r_to_string(instruction)}"
        elif instruction.suffix() == Suffix.RCI:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_rci_to_string(instruction)}"
        elif instruction.suffix() == Suffix.Z:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_z_to_string(instruction)}"
        elif instruction.suffix() == Suffix.ZCI:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_zci_to_string(instruction)}"
        elif instruction.suffix() == Suffix.S_R or instruction.suffix() == Suffix.U_R:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_s_r_to_string(instruction)}"
        elif instruction.suffix() == Suffix.S_RCI or instruction.suffix() == Suffix.U_RCI:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_s_rci_to_string(instruction)}"
        elif instruction.suffix() == Suffix.CI:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_ci_to_string(instruction)}"
        elif instruction.suffix() == Suffix.I:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_i_to_string(instruction)}"
        elif instruction.suffix() == Suffix.DDCI:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_ddci_to_string(instruction)}"
        elif instruction.suffix() == Suffix.ERRI:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_erri_to_string(instruction)}"
        elif instruction.suffix() == Suffix.S_ERRI or instruction.suffix() == Suffix.U_ERRI:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_s_erri_to_string(instruction)}"
        elif instruction.suffix() == Suffix.EDRI:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_edri_to_string(instruction)}"
        elif instruction.suffix() == Suffix.ERII:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_erii_to_string(instruction)}"
        elif instruction.suffix() == Suffix.ERIR:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_erir_to_string(instruction)}"
        elif instruction.suffix() == Suffix.ERID:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_erid_to_string(instruction)}"
        elif instruction.suffix() == Suffix.DMA_RRI:
            return f"{op_code}, {suffix}, {InstructionConverter._convert_dma_rri_to_string(instruction)}"
        else:
            raise ValueError

    @staticmethod
    def _convert_rici_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RICIOpCodes
        assert instruction.suffix() == Suffix.RICI

        ra = RegisterConverter.convert_to_string(instruction.ra())
        imm = str(instruction.imm().value())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        pc = str(instruction.pc().value())
        return f"{ra}, {imm}, {condition}, {pc}"

    @staticmethod
    def _convert_rri_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRIOpCodes
        assert instruction.suffix() == Suffix.RRI

        rc = RegisterConverter.convert_to_string(instruction.rc())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        imm = str(instruction.imm().value())
        return f"{rc}, {ra}, {imm}"

    @staticmethod
    def _convert_rric_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRICOpCodes
        assert instruction.suffix() == Suffix.RRIC

        rc = RegisterConverter.convert_to_string(instruction.rc())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        imm = str(instruction.imm().value())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        return f"{rc}, {ra}, {imm}, {condition}"

    @staticmethod
    def _convert_rrici_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRICIOpCodes
        assert instruction.suffix() == Suffix.RRICI

        rc = RegisterConverter.convert_to_string(instruction.rc())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        imm = str(instruction.imm().value())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        pc = str(instruction.pc().value())
        return f"{rc}, {ra}, {imm}, {condition}, {pc}"

    @staticmethod
    def _convert_rrif_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRIFOpCodes
        assert instruction.suffix() == Suffix.RRIF

        rc = RegisterConverter.convert_to_string(instruction.rc())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        imm = str(instruction.imm().value())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        return f"{rc}, {ra}, {imm}, {condition}"

    @staticmethod
    def _convert_rrr_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRROpCodes
        assert instruction.suffix() == Suffix.RRR

        rc = RegisterConverter.convert_to_string(instruction.rc())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        rb = RegisterConverter.convert_to_string(instruction.rb())
        return f"{rc}, {ra}, {rb}"

    @staticmethod
    def _convert_rrrc_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRRCOpCodes
        assert instruction.suffix() == Suffix.RRRC

        rc = RegisterConverter.convert_to_string(instruction.rc())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        rb = RegisterConverter.convert_to_string(instruction.rb())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        return f"{rc}, {ra}, {rb}, {condition}"

    @staticmethod
    def _convert_rrrci_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRRCIOpCodes
        assert instruction.suffix() == Suffix.RRRCI

        rc = RegisterConverter.convert_to_string(instruction.rc())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        rb = RegisterConverter.convert_to_string(instruction.rb())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        pc = str(instruction.pc().value())
        return f"{rc}, {ra}, {rb}, {condition}, {pc}"

    @staticmethod
    def _convert_zri_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRIOpCodes
        assert instruction.suffix() == Suffix.ZRI

        ra = RegisterConverter.convert_to_string(instruction.ra())
        imm = str(instruction.imm().value())
        return f"{ra}, {imm}"

    @staticmethod
    def _convert_zric_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRICOpCodes
        assert instruction.suffix() == Suffix.ZRIC

        ra = RegisterConverter.convert_to_string(instruction.ra())
        imm = str(instruction.imm().value())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        return f"{ra}, {imm}, {condition}"

    @staticmethod
    def _convert_zrici_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRICIOpCodes
        assert instruction.suffix() == Suffix.ZRICI

        ra = RegisterConverter.convert_to_string(instruction.ra())
        imm = str(instruction.imm().value())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        pc = str(instruction.pc().value())
        return f"{ra}, {imm}, {condition}, {pc}"

    @staticmethod
    def _convert_zrif_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRIFOpCodes
        assert instruction.suffix() == Suffix.ZRIF

        ra = RegisterConverter.convert_to_string(instruction.ra())
        imm = str(instruction.imm().value())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        return f"{ra}, {imm}, {condition}"

    @staticmethod
    def _convert_zrr_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRROpCodes
        assert instruction.suffix() == Suffix.ZRR

        ra = RegisterConverter.convert_to_string(instruction.ra())
        rb = RegisterConverter.convert_to_string(instruction.rb())
        return f"{ra}, {rb}"

    @staticmethod
    def _convert_zrrc_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRRCOpCodes
        assert instruction.suffix() == Suffix.ZRRC

        ra = RegisterConverter.convert_to_string(instruction.ra())
        rb = RegisterConverter.convert_to_string(instruction.rb())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        return f"{ra}, {rb}, {condition}"

    @staticmethod
    def _convert_zrrci_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRRCIOpCodes
        assert instruction.suffix() == Suffix.ZRRCI

        ra = RegisterConverter.convert_to_string(instruction.ra())
        rb = RegisterConverter.convert_to_string(instruction.rb())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        pc = str(instruction.pc().value())
        return f"{ra}, {rb}, {condition}, {pc}"

    @staticmethod
    def _convert_s_rri_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRIOpCodes
        assert instruction.suffix() == Suffix.S_RRI or instruction.suffix() == Suffix.U_RRI

        dc = RegisterConverter.convert_to_string(instruction.dc())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        imm = str(instruction.imm().value())
        return f"{dc}, {ra}, {imm}"

    @staticmethod
    def _convert_s_rric_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRICOpCodes
        assert instruction.suffix() == Suffix.S_RRIC or instruction.suffix() == Suffix.U_RRIC

        dc = RegisterConverter.convert_to_string(instruction.dc())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        imm = str(instruction.imm().value())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        return f"{dc}, {ra}, {imm}, {condition}"

    @staticmethod
    def _convert_s_rrici_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRICIOpCodes
        assert instruction.suffix() == Suffix.S_RRICI or instruction.suffix() == Suffix.U_RRICI

        dc = RegisterConverter.convert_to_string(instruction.dc())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        imm = str(instruction.imm().value())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        pc = str(instruction.pc().value())
        return f"{dc}, {ra}, {imm}, {condition}, {pc}"

    @staticmethod
    def _convert_s_rrif_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRIFOpCodes
        assert instruction.suffix() == Suffix.S_RRIF or instruction.suffix() == Suffix.U_RRIF

        dc = RegisterConverter.convert_to_string(instruction.dc())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        imm = str(instruction.imm().value())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        return f"{dc}, {ra}, {imm}, {condition}"

    @staticmethod
    def _convert_s_rrr_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRROpCodes
        assert instruction.suffix() == Suffix.S_RRR or instruction.suffix() == Suffix.U_RRR

        dc = RegisterConverter.convert_to_string(instruction.dc())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        rb = RegisterConverter.convert_to_string(instruction.rb())
        return f"{dc}, {ra}, {rb}"

    @staticmethod
    def _convert_s_rrrc_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRRCOpCodes
        assert instruction.suffix() == Suffix.S_RRRC or instruction.suffix() == Suffix.U_RRRC

        dc = RegisterConverter.convert_to_string(instruction.dc())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        rb = RegisterConverter.convert_to_string(instruction.rb())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        return f"{dc}, {ra}, {rb}, {condition}"

    @staticmethod
    def _convert_s_rrrci_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRRCIOpCodes
        assert instruction.suffix() == Suffix.S_RRRCI or instruction.suffix() == Suffix.U_RRRCI

        dc = RegisterConverter.convert_to_string(instruction.dc())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        rb = RegisterConverter.convert_to_string(instruction.rb())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        pc = str(instruction.pc().value())
        return f"{dc}, {ra}, {rb}, {condition}, {pc}"

    @staticmethod
    def _convert_rr_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RROpCodes
        assert instruction.suffix() == Suffix.RR

        rc = RegisterConverter.convert_to_string(instruction.rc())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        return f"{rc}, {ra}"

    @staticmethod
    def _convert_rrc_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRCOpCodes
        assert instruction.suffix() == Suffix.RRC

        rc = RegisterConverter.convert_to_string(instruction.rc())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        return f"{rc}, {ra}, {condition}"

    @staticmethod
    def _convert_rrci_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRCIOpCodes
        assert instruction.suffix() == Suffix.RRCI

        rc = RegisterConverter.convert_to_string(instruction.rc())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        pc = str(instruction.pc().value())
        return f"{rc}, {ra}, {condition}, {pc}"

    @staticmethod
    def _convert_zr_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RROpCodes
        assert instruction.suffix() == Suffix.ZR

        ra = RegisterConverter.convert_to_string(instruction.ra())
        return f"{ra}"

    @staticmethod
    def _convert_zrc_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRCOpCodes
        assert instruction.suffix() == Suffix.ZRC

        ra = RegisterConverter.convert_to_string(instruction.ra())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        return f"{ra}, {condition}"

    @staticmethod
    def _convert_zrci_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRCIOpCodes
        assert instruction.suffix() == Suffix.ZRCI

        ra = RegisterConverter.convert_to_string(instruction.ra())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        pc = str(instruction.pc().value())
        return f"{ra}, {condition}, {pc}"

    @staticmethod
    def _convert_s_rr_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RROpCodes
        assert instruction.suffix() == Suffix.S_RR or instruction.suffix() == Suffix.U_RR

        dc = RegisterConverter.convert_to_string(instruction.dc())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        return f"{dc}, {ra}"

    @staticmethod
    def _convert_s_rrc_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRCOpCodes
        assert instruction.suffix() == Suffix.S_RRC or instruction.suffix() == Suffix.U_RRC

        dc = RegisterConverter.convert_to_string(instruction.dc())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        return f"{dc}, {ra}, {condition}"

    @staticmethod
    def _convert_s_rrci_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRCIOpCodes
        assert instruction.suffix() == Suffix.S_RRCI or instruction.suffix() == Suffix.U_RRCI

        dc = RegisterConverter.convert_to_string(instruction.dc())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        pc = str(instruction.pc().value())
        return f"{dc}, {ra}, {condition}, {pc}"

    @staticmethod
    def _convert_drdici_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.DRDICIOpCodes
        assert instruction.suffix() == Suffix.DRDICI

        dc = RegisterConverter.convert_to_string(instruction.dc())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        db = RegisterConverter.convert_to_string(instruction.db())
        imm = str(instruction.imm().value())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        pc = str(instruction.pc().value())
        return f"{dc}, {ra}, {db}, {imm}, {condition}, {pc}"

    @staticmethod
    def _convert_rrri_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRRIOpCodes
        assert instruction.suffix() == Suffix.RRRI

        rc = RegisterConverter.convert_to_string(instruction.rc())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        rb = RegisterConverter.convert_to_string(instruction.rb())
        imm = str(instruction.imm().value())
        return f"{rc}, {ra}, {rb}, {imm}"

    @staticmethod
    def _convert_rrrici_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRRICIOpCodes
        assert instruction.suffix() == Suffix.RRRICI

        rc = RegisterConverter.convert_to_string(instruction.rc())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        rb = RegisterConverter.convert_to_string(instruction.rb())
        imm = str(instruction.imm().value())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        pc = str(instruction.pc().value())
        return f"{rc}, {ra}, {rb}, {imm}, {condition}, {pc}"

    @staticmethod
    def _convert_zrri_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRRIOpCodes
        assert instruction.suffix() == Suffix.ZRRI

        ra = RegisterConverter.convert_to_string(instruction.ra())
        rb = RegisterConverter.convert_to_string(instruction.rb())
        imm = str(instruction.imm().value())
        return f"{ra}, {rb}, {imm}"

    @staticmethod
    def _convert_zrrici_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRRICIOpCodes
        assert instruction.suffix() == Suffix.ZRRICI

        ra = RegisterConverter.convert_to_string(instruction.ra())
        rb = RegisterConverter.convert_to_string(instruction.rb())
        imm = str(instruction.imm().value())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        pc = str(instruction.pc().value())
        return f"{ra}, {rb}, {imm}, {condition}, {pc}"

    @staticmethod
    def _convert_s_rrri_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRRIOpCodes
        assert instruction.suffix() == Suffix.S_RRRI or instruction.suffix() == Suffix.U_RRRI

        dc = RegisterConverter.convert_to_string(instruction.dc())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        rb = RegisterConverter.convert_to_string(instruction.rb())
        imm = str(instruction.imm().value())
        return f"{dc}, {ra}, {rb}, {imm}"

    @staticmethod
    def _convert_s_rrrici_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RRRICIOpCodes
        assert instruction.suffix() == Suffix.S_RRRICI or instruction.suffix() == Suffix.U_RRRICI

        dc = RegisterConverter.convert_to_string(instruction.dc())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        rb = RegisterConverter.convert_to_string(instruction.rb())
        imm = str(instruction.imm().value())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        pc = str(instruction.pc().value())
        return f"{dc}, {ra}, {rb}, {imm}, {condition}, {pc}"

    @staticmethod
    def _convert_rir_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RIROpCodes
        assert instruction.suffix() == Suffix.RIR

        rc = RegisterConverter.convert_to_string(instruction.rc())
        imm = str(instruction.imm().value())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        return f"{rc}, {imm}, {ra}"

    @staticmethod
    def _convert_rirc_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RIRCOpCodes
        assert instruction.suffix() == Suffix.RIRC

        rc = RegisterConverter.convert_to_string(instruction.rc())
        imm = str(instruction.imm().value())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        return f"{rc}, {imm}, {ra}, {condition}"

    @staticmethod
    def _convert_rirci_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RIRCIOpCodes
        assert instruction.suffix() == Suffix.RIRCI

        rc = RegisterConverter.convert_to_string(instruction.rc())
        imm = str(instruction.imm().value())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        pc = str(instruction.pc().value())
        return f"{rc}, {imm}, {ra}, {condition}, {pc}"

    @staticmethod
    def _convert_zir_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RIROpCodes
        assert instruction.suffix() == Suffix.ZIR

        imm = str(instruction.imm().value())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        return f"{imm}, {ra}"

    @staticmethod
    def _convert_zirc_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RIRCOpCodes
        assert instruction.suffix() == Suffix.ZIRC

        imm = str(instruction.imm().value())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        return f"{imm}, {ra}, {condition}"

    @staticmethod
    def _convert_zirci_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RIRCIOpCodes
        assert instruction.suffix() == Suffix.ZIRCI

        imm = str(instruction.imm().value())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        pc = str(instruction.pc().value())
        return f"{imm}, {ra}, {condition}, {pc}"

    @staticmethod
    def _convert_s_rirc_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RIRCOpCodes
        assert instruction.suffix() == Suffix.S_RIRC or instruction.suffix() == Suffix.U_RIRC

        dc = RegisterConverter.convert_to_string(instruction.dc())
        imm = str(instruction.imm().value())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        return f"{dc}, {imm}, {ra}, {condition}"

    @staticmethod
    def _convert_s_rirci_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RIRCIOpCodes
        assert instruction.suffix() == Suffix.S_RIRCI or instruction.suffix() == Suffix.U_RIRCI

        dc = RegisterConverter.convert_to_string(instruction.dc())
        imm = str(instruction.imm().value())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        pc = str(instruction.pc().value())
        return f"{dc}, {imm}, {ra}, {condition}, {pc}"

    @staticmethod
    def _convert_r_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.ROpCodes
        assert instruction.suffix() == Suffix.R

        rc = RegisterConverter.convert_to_string(instruction.rc())
        return f"{rc}"

    @staticmethod
    def _convert_rci_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RCIOpCodes
        assert instruction.suffix() == Suffix.RCI

        rc = RegisterConverter.convert_to_string(instruction.rc())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        pc = str(instruction.pc().value())
        return f"{rc}, {condition}, {pc}"

    @staticmethod
    def _convert_z_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.ROpCodes or instruction.op_code() == OpCode.NOP
        assert instruction.suffix() == Suffix.Z

        return ""

    @staticmethod
    def _convert_zci_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RCIOpCodes
        assert instruction.suffix() == Suffix.ZCI

        condition = ConditionConverter.convert_to_string(instruction.condition())
        pc = str(instruction.pc().value())
        return f"{condition}, {pc}"

    @staticmethod
    def _convert_s_r_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.ROpCodes
        assert instruction.suffix() == Suffix.S_R or instruction.suffix() == Suffix.U_R

        dc = RegisterConverter.convert_to_string(instruction.dc())
        return f"{dc}"

    @staticmethod
    def _convert_s_rci_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.RCIOpCodes
        assert instruction.suffix() == Suffix.S_RCI or instruction.suffix() == Suffix.U_RCI

        dc = RegisterConverter.convert_to_string(instruction.dc())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        pc = str(instruction.pc().value())
        return f"{dc}, {condition}, {pc}"

    @staticmethod
    def _convert_ci_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.CIOpCodes
        assert instruction.suffix() == Suffix.CI

        condition = ConditionConverter.convert_to_string(instruction.condition())
        pc = str(instruction.pc().value())
        return f"{condition}, {pc}"

    @staticmethod
    def _convert_i_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.IOpCodes
        assert instruction.suffix() == Suffix.I

        imm = str(instruction.imm().value())
        return f"{imm}"

    @staticmethod
    def _convert_ddci_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.DDCIOpCodes
        assert instruction.suffix() == Suffix.DDCI

        dc = RegisterConverter.convert_to_string(instruction.dc())
        db = RegisterConverter.convert_to_string(instruction.db())
        condition = ConditionConverter.convert_to_string(instruction.condition())
        pc = str(instruction.pc().value())
        return f"{dc}, {db}, {condition}, {pc}"

    @staticmethod
    def _convert_erri_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.ERRIOpCodes
        assert instruction.suffix() == Suffix.ERRI

        endian = EndianConverter.convert_to_string(instruction.endian())
        rc = RegisterConverter.convert_to_string(instruction.rc())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        off = str(instruction.off().value())
        return f"{endian}, {rc}, {ra}, {off}"

    @staticmethod
    def _convert_s_erri_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.ERRIOpCodes
        assert instruction.suffix() == Suffix.S_ERRI or instruction.suffix() == Suffix.U_ERRI

        endian = EndianConverter.convert_to_string(instruction.endian())
        dc = RegisterConverter.convert_to_string(instruction.dc())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        off = str(instruction.off().value())
        return f"{endian}, {dc}, {ra}, {off}"

    @staticmethod
    def _convert_edri_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.EDRIOpCodes
        assert instruction.suffix() == Suffix.EDRI

        endian = EndianConverter.convert_to_string(instruction.endian())
        dc = RegisterConverter.convert_to_string(instruction.dc())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        off = str(instruction.off().value())
        return f"{endian}, {dc}, {ra}, {off}"

    @staticmethod
    def _convert_erii_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.ERIIOpCodes
        assert instruction.suffix() == Suffix.ERII

        endian = EndianConverter.convert_to_string(instruction.endian())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        off = str(instruction.off().value())
        imm = str(instruction.imm().value())
        return f"{endian}, {ra}, {off}, {imm}"

    @staticmethod
    def _convert_erir_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.ERIROpCodes
        assert instruction.suffix() == Suffix.ERIR

        endian = EndianConverter.convert_to_string(instruction.endian())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        off = str(instruction.off().value())
        rb = RegisterConverter.convert_to_string(instruction.rb())
        return f"{endian}, {ra}, {off}, {rb}"

    @staticmethod
    def _convert_erid_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.ERIDOpCodes
        assert instruction.suffix() == Suffix.ERID

        endian = EndianConverter.convert_to_string(instruction.endian())
        ra = RegisterConverter.convert_to_string(instruction.ra())
        off = str(instruction.off().value())
        db = RegisterConverter.convert_to_string(instruction.db())
        return f"{endian}, {ra}, {off}, {db}"

    @staticmethod
    def _convert_dma_rri_to_string(instruction: Instruction) -> str:
        assert instruction.op_code() in Instruction.DMARRIOpCodes
        assert instruction.suffix() == Suffix.DMA_RRI

        ra = RegisterConverter.convert_to_string(instruction.ra())
        rb = RegisterConverter.convert_to_string(instruction.rb())
        imm = str(instruction.imm().value())
        return f"{ra}, {rb}, {imm}"
