#include "converter/op_code_converter.h"

#include <stdexcept>

namespace upmem_sim::converter {

std::string OpCodeConverter::to_string(abi::instruction::OpCode op_code) {
  if (op_code == abi::instruction::ACQUIRE) {
    return "acquire";
  } else if (op_code == abi::instruction::RELEASE) {
    return "release";
  } else if (op_code == abi::instruction::ADD) {
    return "add";
  } else if (op_code == abi::instruction::ADDC) {
    return "addc";
  } else if (op_code == abi::instruction::AND) {
    return "and";
  } else if (op_code == abi::instruction::ANDN) {
    return "andn";
  } else if (op_code == abi::instruction::ASR) {
    return "asr";
  } else if (op_code == abi::instruction::CAO) {
    return "cao";
  } else if (op_code == abi::instruction::CLO) {
    return "clo";
  } else if (op_code == abi::instruction::CLS) {
    return "cls";
  } else if (op_code == abi::instruction::CLZ) {
    return "clz";
  } else if (op_code == abi::instruction::CMPB4) {
    return "cmpb4";
  } else if (op_code == abi::instruction::DIV_STEP) {
    return "div_step";
  } else if (op_code == abi::instruction::EXTSB) {
    return "extsb";
  } else if (op_code == abi::instruction::EXTSH) {
    return "extsh";
  } else if (op_code == abi::instruction::EXTUB) {
    return "extub";
  } else if (op_code == abi::instruction::EXTUH) {
    return "extuh";
  } else if (op_code == abi::instruction::LSL) {
    return "lsl";
  } else if (op_code == abi::instruction::LSL_ADD) {
    return "lsl_add";
  } else if (op_code == abi::instruction::LSL_SUB) {
    return "lsl_sub";
  } else if (op_code == abi::instruction::LSL1) {
    return "lsl1";
  } else if (op_code == abi::instruction::LSL1X) {
    return "lsl1x";
  } else if (op_code == abi::instruction::LSLX) {
    return "lslx";
  } else if (op_code == abi::instruction::LSR) {
    return "lsr";
  } else if (op_code == abi::instruction::LSR_ADD) {
    return "lsr_add";
  } else if (op_code == abi::instruction::LSR1) {
    return "lsr1";
  } else if (op_code == abi::instruction::LSR1X) {
    return "lsr1x";
  } else if (op_code == abi::instruction::LSRX) {
    return "lsrx";
  } else if (op_code == abi::instruction::MUL_SH_SH) {
    return "mul_sh_sh";
  } else if (op_code == abi::instruction::MUL_SH_SL) {
    return "mul_sh_sl";
  } else if (op_code == abi::instruction::MUL_SH_UH) {
    return "mul_sh_uh";
  } else if (op_code == abi::instruction::MUL_SH_UL) {
    return "mul_sh_ul";
  } else if (op_code == abi::instruction::MUL_SL_SH) {
    return "mul_sl_sh";
  } else if (op_code == abi::instruction::MUL_SL_SL) {
    return "mul_sl_sl";
  } else if (op_code == abi::instruction::MUL_SL_UH) {
    return "mul_sl_uh";
  } else if (op_code == abi::instruction::MUL_SL_UL) {
    return "mul_sl_ul";
  } else if (op_code == abi::instruction::MUL_STEP) {
    return "mul_step";
  } else if (op_code == abi::instruction::MUL_UH_UH) {
    return "mul_uh_uh";
  } else if (op_code == abi::instruction::MUL_UH_UL) {
    return "mul_uh_ul";
  } else if (op_code == abi::instruction::MUL_UL_UH) {
    return "mul_ul_uh";
  } else if (op_code == abi::instruction::MUL_UL_UL) {
    return "mul_ul_ul";
  } else if (op_code == abi::instruction::NAND) {
    return "nand";
  } else if (op_code == abi::instruction::NOR) {
    return "nor";
  } else if (op_code == abi::instruction::NXOR) {
    return "nxor";
  } else if (op_code == abi::instruction::OR) {
    return "or";
  } else if (op_code == abi::instruction::ORN) {
    return "orn";
  } else if (op_code == abi::instruction::ROL) {
    return "rol";
  } else if (op_code == abi::instruction::ROL_ADD) {
    return "rol_add";
  } else if (op_code == abi::instruction::ROR) {
    return "ror";
  } else if (op_code == abi::instruction::RSUB) {
    return "rsub";
  } else if (op_code == abi::instruction::RSUBC) {
    return "rsubc";
  } else if (op_code == abi::instruction::SUB) {
    return "sub";
  } else if (op_code == abi::instruction::SUBC) {
    return "subc";
  } else if (op_code == abi::instruction::XOR) {
    return "xor";
  } else if (op_code == abi::instruction::BOOT) {
    return "boot";
  } else if (op_code == abi::instruction::RESUME) {
    return "resume";
  } else if (op_code == abi::instruction::STOP) {
    return "stop";
  } else if (op_code == abi::instruction::CALL) {
    return "call";
  } else if (op_code == abi::instruction::FAULT) {
    return "fault";
  } else if (op_code == abi::instruction::NOP) {
    return "nop";
  } else if (op_code == abi::instruction::SATS) {
    return "sats";
  } else if (op_code == abi::instruction::MOVD) {
    return "movd";
  } else if (op_code == abi::instruction::SWAPD) {
    return "swapd";
  } else if (op_code == abi::instruction::HASH) {
    return "hash";
  } else if (op_code == abi::instruction::TIME) {
    return "time";
  } else if (op_code == abi::instruction::TIME_CFG) {
    return "time_cfg";
  } else if (op_code == abi::instruction::LBS) {
    return "lbs";
  } else if (op_code == abi::instruction::LBU) {
    return "lbu";
  } else if (op_code == abi::instruction::LD) {
    return "ld";
  } else if (op_code == abi::instruction::LHS) {
    return "lhs";
  } else if (op_code == abi::instruction::LHU) {
    return "lhu";
  } else if (op_code == abi::instruction::LW) {
    return "lw";
  } else if (op_code == abi::instruction::SB) {
    return "sb";
  } else if (op_code == abi::instruction::SB_ID) {
    return "sb_id";
  } else if (op_code == abi::instruction::SD) {
    return "sd";
  } else if (op_code == abi::instruction::SD_ID) {
    return "sd_id";
  } else if (op_code == abi::instruction::SH) {
    return "sh";
  } else if (op_code == abi::instruction::SH_ID) {
    return "sh_id";
  } else if (op_code == abi::instruction::SW) {
    return "sw";
  } else if (op_code == abi::instruction::SW_ID) {
    return "sw_id";
  } else if (op_code == abi::instruction::LDMA) {
    return "ldma";
  } else if (op_code == abi::instruction::LDMAI) {
    return "ldmai";
  } else if (op_code == abi::instruction::SDMA) {
    return "sdma";
  } else {
    throw std::invalid_argument("");
  }
}

}  // namespace upmem_sim::converter
