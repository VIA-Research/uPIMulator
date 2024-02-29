from abi.isa.instruction.op_code import OpCode


class OpCodeConverter:
    def __init__(self):
        pass

    @staticmethod
    def convert_to_op_code(op_code: str) -> OpCode:
        if op_code == "acquire":
            return OpCode.ACQUIRE
        elif op_code == "release":
            return OpCode.RELEASE
        elif op_code == "add":
            return OpCode.ADD
        elif op_code == "addc":
            return OpCode.ADDC
        elif op_code == "and":
            return OpCode.AND
        elif op_code == "andn":
            return OpCode.ANDN
        elif op_code == "asr":
            return OpCode.ASR
        elif op_code == "cao":
            return OpCode.CAO
        elif op_code == "clo":
            return OpCode.CLO
        elif op_code == "cls":
            return OpCode.CLS
        elif op_code == "clz":
            return OpCode.CLZ
        elif op_code == "cmpb4":
            return OpCode.CMPB4
        elif op_code == "div_step":
            return OpCode.DIV_STEP
        elif op_code == "extsb":
            return OpCode.EXTSB
        elif op_code == "extsh":
            return OpCode.EXTSH
        elif op_code == "extub":
            return OpCode.EXTUB
        elif op_code == "extuh":
            return OpCode.EXTUH
        elif op_code == "lsl":
            return OpCode.LSL
        elif op_code == "lsl_add":
            return OpCode.LSL_ADD
        elif op_code == "lsl_sub":
            return OpCode.LSL_SUB
        elif op_code == "lsl1":
            return OpCode.LSL1
        elif op_code == "lsl1x":
            return OpCode.LSL1X
        elif op_code == "lslx":
            return OpCode.LSLX
        elif op_code == "lsr":
            return OpCode.LSR
        elif op_code == "lsr_add":
            return OpCode.LSR_ADD
        elif op_code == "lsr1":
            return OpCode.LSR1
        elif op_code == "lsr1x":
            return OpCode.LSR1X
        elif op_code == "lsrx":
            return OpCode.LSRX
        elif op_code == "mul_sh_sh":
            return OpCode.MUL_SH_SH
        elif op_code == "mul_sh_sl":
            return OpCode.MUL_SH_SL
        elif op_code == "mul_sh_uh":
            return OpCode.MUL_SH_UH
        elif op_code == "mul_sh_ul":
            return OpCode.MUL_SH_UL
        elif op_code == "mul_sl_sh":
            return OpCode.MUL_SL_SH
        elif op_code == "mul_sl_sl":
            return OpCode.MUL_SL_SL
        elif op_code == "mul_sl_uh":
            return OpCode.MUL_SL_UH
        elif op_code == "mul_sl_ul":
            return OpCode.MUL_SL_UL
        elif op_code == "mul_step":
            return OpCode.MUL_STEP
        elif op_code == "mul_uh_uh":
            return OpCode.MUL_UH_UH
        elif op_code == "mul_uh_ul":
            return OpCode.MUL_UH_UL
        elif op_code == "mul_ul_uh":
            return OpCode.MUL_UL_UH
        elif op_code == "mul_ul_ul":
            return OpCode.MUL_UL_UL
        elif op_code == "nand":
            return OpCode.NAND
        elif op_code == "nor":
            return OpCode.NOR
        elif op_code == "nxor":
            return OpCode.NXOR
        elif op_code == "or":
            return OpCode.OR
        elif op_code == "orn":
            return OpCode.ORN
        elif op_code == "rol":
            return OpCode.ROL
        elif op_code == "rol_add":
            return OpCode.ROL_ADD
        elif op_code == "ror":
            return OpCode.ROR
        elif op_code == "rsub":
            return OpCode.RSUB
        elif op_code == "rsubc":
            return OpCode.RSUBC
        elif op_code == "sub":
            return OpCode.SUB
        elif op_code == "subc":
            return OpCode.SUBC
        elif op_code == "xor":
            return OpCode.XOR
        elif op_code == "boot":
            return OpCode.BOOT
        elif op_code == "resume":
            return OpCode.RESUME
        elif op_code == "stop":
            return OpCode.STOP
        elif op_code == "call":
            return OpCode.CALL
        elif op_code == "fault":
            return OpCode.FAULT
        elif op_code == "nop":
            return OpCode.NOP
        elif op_code == "sats":
            return OpCode.SATS
        elif op_code == "movd":
            return OpCode.MOVD
        elif op_code == "swapd":
            return OpCode.SWAPD
        elif op_code == "hash":
            return OpCode.HASH
        elif op_code == "time":
            return OpCode.TIME
        elif op_code == "time_cfg":
            return OpCode.TIME_CFG
        elif op_code == "lbs":
            return OpCode.LBS
        elif op_code == "lbu":
            return OpCode.LBU
        elif op_code == "ld":
            return OpCode.LD
        elif op_code == "lhs":
            return OpCode.LHS
        elif op_code == "lhu":
            return OpCode.LHU
        elif op_code == "lw":
            return OpCode.LW
        elif op_code == "sb":
            return OpCode.SB
        elif op_code == "sb_id":
            return OpCode.SB_ID
        elif op_code == "sd":
            return OpCode.SD
        elif op_code == "sd_id":
            return OpCode.SD_ID
        elif op_code == "sh":
            return OpCode.SH
        elif op_code == "sh_id":
            return OpCode.SH_ID
        elif op_code == "sw":
            return OpCode.SW
        elif op_code == "sw_id":
            return OpCode.SW_ID
        elif op_code == "ldma":
            return OpCode.LDMA
        elif op_code == "ldmai":
            return OpCode.LDMAI
        elif op_code == "sdma":
            return OpCode.SDMA
        else:
            raise ValueError

    @staticmethod
    def convert_to_string(op_code: OpCode) -> str:
        if op_code == OpCode.ACQUIRE:
            return "acquire"
        elif op_code == OpCode.RELEASE:
            return "release"
        elif op_code == OpCode.ADD:
            return "add"
        elif op_code == OpCode.ADDC:
            return "addc"
        elif op_code == OpCode.AND:
            return "and"
        elif op_code == OpCode.ANDN:
            return "andn"
        elif op_code == OpCode.ASR:
            return "asr"
        elif op_code == OpCode.CAO:
            return "cao"
        elif op_code == OpCode.CLO:
            return "clo"
        elif op_code == OpCode.CLS:
            return "cls"
        elif op_code == OpCode.CLZ:
            return "clz"
        elif op_code == OpCode.CMPB4:
            return "cmpb4"
        elif op_code == OpCode.DIV_STEP:
            return "div_step"
        elif op_code == OpCode.EXTSB:
            return "extsb"
        elif op_code == OpCode.EXTSH:
            return "extsh"
        elif op_code == OpCode.EXTUB:
            return "extub"
        elif op_code == OpCode.EXTUH:
            return "extuh"
        elif op_code == OpCode.LSL:
            return "lsl"
        elif op_code == OpCode.LSL_ADD:
            return "lsl_add"
        elif op_code == OpCode.LSL_SUB:
            return "lsl_sub"
        elif op_code == OpCode.LSL1:
            return "lsl1"
        elif op_code == OpCode.LSL1X:
            return "lsl1x"
        elif op_code == OpCode.LSLX:
            return "lslx"
        elif op_code == OpCode.LSR:
            return "lsr"
        elif op_code == OpCode.LSR_ADD:
            return "lsr_add"
        elif op_code == OpCode.LSR1:
            return "lsr1"
        elif op_code == OpCode.LSR1X:
            return "lsr1x"
        elif op_code == OpCode.LSRX:
            return "lsrx"
        elif op_code == OpCode.MUL_SH_SH:
            return "mul_sh_sh"
        elif op_code == OpCode.MUL_SH_SL:
            return "mul_sh_sl"
        elif op_code == OpCode.MUL_SH_UH:
            return "mul_sh_uh"
        elif op_code == OpCode.MUL_SH_UL:
            return "mul_sh_ul"
        elif op_code == OpCode.MUL_SL_SH:
            return "mul_sl_sh"
        elif op_code == OpCode.MUL_SL_SL:
            return "mul_sl_sl"
        elif op_code == OpCode.MUL_SL_UH:
            return "mul_sl_uh"
        elif op_code == OpCode.MUL_SL_UL:
            return "mul_sl_ul"
        elif op_code == OpCode.MUL_STEP:
            return "mul_step"
        elif op_code == OpCode.MUL_UH_UH:
            return "mul_uh_uh"
        elif op_code == OpCode.MUL_UH_UL:
            return "mul_uh_ul"
        elif op_code == OpCode.MUL_UL_UH:
            return "mul_ul_uh"
        elif op_code == OpCode.MUL_UL_UL:
            return "mul_ul_ul"
        elif op_code == OpCode.NAND:
            return "nand"
        elif op_code == OpCode.NOR:
            return "nor"
        elif op_code == OpCode.NXOR:
            return "nxor"
        elif op_code == OpCode.OR:
            return "or"
        elif op_code == OpCode.ORN:
            return "orn"
        elif op_code == OpCode.ROL:
            return "rol"
        elif op_code == OpCode.ROL_ADD:
            return "rol_add"
        elif op_code == OpCode.ROR:
            return "ror"
        elif op_code == OpCode.RSUB:
            return "rsub"
        elif op_code == OpCode.RSUBC:
            return "rsubc"
        elif op_code == OpCode.SUB:
            return "sub"
        elif op_code == OpCode.SUBC:
            return "subc"
        elif op_code == OpCode.XOR:
            return "xor"
        elif op_code == OpCode.BOOT:
            return "boot"
        elif op_code == OpCode.RESUME:
            return "resume"
        elif op_code == OpCode.STOP:
            return "stop"
        elif op_code == OpCode.CALL:
            return "call"
        elif op_code == OpCode.FAULT:
            return "fault"
        elif op_code == OpCode.NOP:
            return "nop"
        elif op_code == OpCode.SATS:
            return "sats"
        elif op_code == OpCode.MOVD:
            return "movd"
        elif op_code == OpCode.SWAPD:
            return "swapd"
        elif op_code == OpCode.HASH:
            return "hash"
        elif op_code == OpCode.TIME:
            return "time"
        elif op_code == OpCode.TIME_CFG:
            return "time_cfg"
        elif op_code == OpCode.LBS:
            return "lbs"
        elif op_code == OpCode.LBU:
            return "lbu"
        elif op_code == OpCode.LD:
            return "ld"
        elif op_code == OpCode.LHS:
            return "lhs"
        elif op_code == OpCode.LHU:
            return "lhu"
        elif op_code == OpCode.LW:
            return "lw"
        elif op_code == OpCode.SB:
            return "sb"
        elif op_code == OpCode.SB_ID:
            return "sb_id"
        elif op_code == OpCode.SD:
            return "sd"
        elif op_code == OpCode.SD_ID:
            return "sd_id"
        elif op_code == OpCode.SH:
            return "sh"
        elif op_code == OpCode.SH_ID:
            return "sh_id"
        elif op_code == OpCode.SW:
            return "sw"
        elif op_code == OpCode.SW_ID:
            return "sw_id"
        elif op_code == OpCode.LDMA:
            return "ldma"
        elif op_code == OpCode.LDMAI:
            return "ldmai"
        elif op_code == OpCode.SDMA:
            return "sdma"
        else:
            raise ValueError
