#include "converter/op_code_converter.h"

#include <stdexcept>

namespace upmem_profiler::converter {

abi::instruction::OpCode OpCodeConverter::to_op_code(std::string op_code) {
  if (op_code == "acquire") {
    return abi::instruction::ACQUIRE;
  } else if (op_code == "release") {
    return abi::instruction::RELEASE;
  } else if (op_code == "add") {
    return abi::instruction::ADD;
  } else if (op_code == "addc") {
    return abi::instruction::ADDC;
  } else if (op_code == "and") {
    return abi::instruction::AND;
  } else if (op_code == "andn") {
    return abi::instruction::ANDN;
  } else if (op_code == "asr") {
    return abi::instruction::ASR;
  } else if (op_code == "cao") {
    return abi::instruction::CAO;
  } else if (op_code == "clo") {
    return abi::instruction::CLO;
  } else if (op_code == "cls") {
    return abi::instruction::CLS;
  } else if (op_code == "clz") {
    return abi::instruction::CLZ;
  } else if (op_code == "cmpb4") {
    return abi::instruction::CMPB4;
  } else if (op_code == "div_step") {
    return abi::instruction::DIV_STEP;
  } else if (op_code == "extsb") {
    return abi::instruction::EXTSB;
  } else if (op_code == "extsh") {
    return abi::instruction::EXTSH;
  } else if (op_code == "extub") {
    return abi::instruction::EXTUB;
  } else if (op_code == "extuh") {
    return abi::instruction::EXTUH;
  } else if (op_code == "lsl") {
    return abi::instruction::LSL;
  } else if (op_code == "lsl_add") {
    return abi::instruction::LSL_ADD;
  } else if (op_code == "lsl_sub") {
    return abi::instruction::LSL_SUB;
  } else if (op_code == "lsl1") {
    return abi::instruction::LSL1;
  } else if (op_code == "lsl1x") {
    return abi::instruction::LSL1X;
  } else if (op_code == "lslx") {
    return abi::instruction::LSLX;
  } else if (op_code == "lsr") {
    return abi::instruction::LSR;
  } else if (op_code == "lsr_add") {
    return abi::instruction::LSR_ADD;
  } else if (op_code == "lsr1") {
    return abi::instruction::LSR1;
  } else if (op_code == "lsr1x") {
    return abi::instruction::LSR1X;
  } else if (op_code == "lsrx") {
    return abi::instruction::LSRX;
  } else if (op_code == "mul_sh_sh") {
    return abi::instruction::MUL_SH_SH;
  } else if (op_code == "mul_sh_sl") {
    return abi::instruction::MUL_SH_SL;
  } else if (op_code == "mul_sh_uh") {
    return abi::instruction::MUL_SH_UH;
  } else if (op_code == "mul_sh_ul") {
    return abi::instruction::MUL_SH_UL;
  } else if (op_code == "mul_sl_sh") {
    return abi::instruction::MUL_SL_SH;
  } else if (op_code == "mul_sl_sl") {
    return abi::instruction::MUL_SL_SL;
  } else if (op_code == "mul_sl_uh") {
    return abi::instruction::MUL_SL_UH;
  } else if (op_code == "mul_sl_ul") {
    return abi::instruction::MUL_SL_UL;
  } else if (op_code == "mul_step") {
    return abi::instruction::MUL_STEP;
  } else if (op_code == "mul_uh_uh") {
    return abi::instruction::MUL_UH_UH;
  } else if (op_code == "mul_uh_ul") {
    return abi::instruction::MUL_UH_UL;
  } else if (op_code == "mul_ul_uh") {
    return abi::instruction::MUL_UL_UH;
  } else if (op_code == "mul_ul_ul") {
    return abi::instruction::MUL_UL_UL;
  } else if (op_code == "nand") {
    return abi::instruction::NAND;
  } else if (op_code == "nor") {
    return abi::instruction::NOR;
  } else if (op_code == "nxor") {
    return abi::instruction::NXOR;
  } else if (op_code == "or") {
    return abi::instruction::OR;
  } else if (op_code == "orn") {
    return abi::instruction::ORN;
  } else if (op_code == "rol") {
    return abi::instruction::ROL;
  } else if (op_code == "rol_add") {
    return abi::instruction::ROL_ADD;
  } else if (op_code == "ror") {
    return abi::instruction::ROR;
  } else if (op_code == "rsub") {
    return abi::instruction::RSUB;
  } else if (op_code == "rsubc") {
    return abi::instruction::RSUBC;
  } else if (op_code == "sub") {
    return abi::instruction::SUB;
  } else if (op_code == "subc") {
    return abi::instruction::SUBC;
  } else if (op_code == "xor") {
    return abi::instruction::XOR;
  } else if (op_code == "boot") {
    return abi::instruction::BOOT;
  } else if (op_code == "resume") {
    return abi::instruction::RESUME;
  } else if (op_code == "stop") {
    return abi::instruction::STOP;
  } else if (op_code == "call") {
    return abi::instruction::CALL;
  } else if (op_code == "fault") {
    return abi::instruction::FAULT;
  } else if (op_code == "nop") {
    return abi::instruction::NOP;
  } else if (op_code == "sats") {
    return abi::instruction::SATS;
  } else if (op_code == "movd") {
    return abi::instruction::MOVD;
  } else if (op_code == "swapd") {
    return abi::instruction::SWAPD;
  } else if (op_code == "hash") {
    return abi::instruction::HASH;
  } else if (op_code == "time") {
    return abi::instruction::TIME;
  } else if (op_code == "time_cfg") {
    return abi::instruction::TIME_CFG;
  } else if (op_code == "lbs") {
    return abi::instruction::LBS;
  } else if (op_code == "lbu") {
    return abi::instruction::LBU;
  } else if (op_code == "ld") {
    return abi::instruction::LD;
  } else if (op_code == "lhs") {
    return abi::instruction::LHS;
  } else if (op_code == "lhu") {
    return abi::instruction::LHU;
  } else if (op_code == "lw") {
    return abi::instruction::LW;
  } else if (op_code == "sb") {
    return abi::instruction::SB;
  } else if (op_code == "sb_id") {
    return abi::instruction::SB_ID;
  } else if (op_code == "sd") {
    return abi::instruction::SD;
  } else if (op_code == "sd_id") {
    return abi::instruction::SD_ID;
  } else if (op_code == "sh") {
    return abi::instruction::SH;
  } else if (op_code == "SH_ID") {
    return abi::instruction::SH_ID;
  } else if (op_code == "sw") {
    return abi::instruction::SW;
  } else if (op_code == "sw_id") {
    return abi::instruction::SW_ID;
  } else if (op_code == "ldma") {
    return abi::instruction::LDMA;
  } else if (op_code == "ldmai") {
    return abi::instruction::LDMAI;
  } else if (op_code == "sdma") {
    return abi::instruction::SDMA;
  } else {
    throw std::invalid_argument("");
  }
}

} // namespace upmem_profiler::converter
