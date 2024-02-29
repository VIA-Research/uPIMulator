#include "converter/instruction_converter.h"

#include <sstream>

#include "converter/condition_converter.h"
#include "converter/endian_converter.h"
#include "converter/op_code_converter.h"
#include "converter/reg_converter.h"
#include "converter/suffix_converter.h"

namespace upmem_sim::converter {

std::string InstructionConverter::to_string(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;

  ss << "[" << instruction->thread()->id() << "] ";
  ss << OpCodeConverter::to_string(instruction->op_code()) << ", ";
  ss << SuffixConverter::to_string(instruction->suffix()) << ", ";

  abi::instruction::Suffix suffix = instruction->suffix();
  if (suffix == abi::instruction::RICI) {
    ss << to_string_rici(instruction);
  } else if (suffix == abi::instruction::RRI) {
    ss << to_string_rri(instruction);
  } else if (suffix == abi::instruction::RRIC) {
    ss << to_string_rric(instruction);
  } else if (suffix == abi::instruction::RRICI) {
    ss << to_string_rrici(instruction);
  } else if (suffix == abi::instruction::RRIF) {
    ss << to_string_rrif(instruction);
  } else if (suffix == abi::instruction::RRR) {
    ss << to_string_rrr(instruction);
  } else if (suffix == abi::instruction::RRRC) {
    ss << to_string_rrrc(instruction);
  } else if (suffix == abi::instruction::RRRCI) {
    ss << to_string_rrrci(instruction);
  } else if (suffix == abi::instruction::ZRI) {
    ss << to_string_zri(instruction);
  } else if (suffix == abi::instruction::ZRIC) {
    ss << to_string_zric(instruction);
  } else if (suffix == abi::instruction::ZRICI) {
    ss << to_string_zrici(instruction);
  } else if (suffix == abi::instruction::ZRIF) {
    ss << to_string_zrif(instruction);
  } else if (suffix == abi::instruction::ZRR) {
    ss << to_string_zrr(instruction);
  } else if (suffix == abi::instruction::ZRRC) {
    ss << to_string_zrrc(instruction);
  } else if (suffix == abi::instruction::ZRRCI) {
    ss << to_string_zrrci(instruction);
  } else if (suffix == abi::instruction::S_RRI or
             suffix == abi::instruction::U_RRI) {
    ss << to_string_s_rri(instruction);
  } else if (suffix == abi::instruction::S_RRIC or
             suffix == abi::instruction::U_RRIC) {
    ss << to_string_s_rric(instruction);
  } else if (suffix == abi::instruction::S_RRICI or
             suffix == abi::instruction::U_RRICI) {
    ss << to_string_s_rrici(instruction);
  } else if (suffix == abi::instruction::S_RRIF or
             suffix == abi::instruction::U_RRIF) {
    ss << to_string_s_rrif(instruction);
  } else if (suffix == abi::instruction::S_RRR or
             suffix == abi::instruction::U_RRR) {
    ss << to_string_s_rrr(instruction);
  } else if (suffix == abi::instruction::S_RRRC or
             suffix == abi::instruction::U_RRRC) {
    ss << to_string_s_rrrc(instruction);
  } else if (suffix == abi::instruction::S_RRRCI or
             suffix == abi::instruction::U_RRRCI) {
    ss << to_string_s_rrrci(instruction);
  } else if (suffix == abi::instruction::RR) {
    ss << to_string_rr(instruction);
  } else if (suffix == abi::instruction::RRC) {
    ss << to_string_rrc(instruction);
  } else if (suffix == abi::instruction::RRCI) {
    ss << to_string_rrci(instruction);
  } else if (suffix == abi::instruction::ZR) {
    ss << to_string_zr(instruction);
  } else if (suffix == abi::instruction::ZRC) {
    ss << to_string_zrc(instruction);
  } else if (suffix == abi::instruction::ZRCI) {
    ss << to_string_zrci(instruction);
  } else if (suffix == abi::instruction::S_RR or
             suffix == abi::instruction::U_RR) {
    ss << to_string_s_rr(instruction);
  } else if (suffix == abi::instruction::S_RRC or
             suffix == abi::instruction::U_RRC) {
    ss << to_string_s_rrc(instruction);
  } else if (suffix == abi::instruction::S_RRCI or
             suffix == abi::instruction::U_RRCI) {
    ss << to_string_s_rrci(instruction);
  } else if (suffix == abi::instruction::DRDICI) {
    ss << to_string_drdici(instruction);
  } else if (suffix == abi::instruction::RRRI) {
    ss << to_string_rrri(instruction);
  } else if (suffix == abi::instruction::RRRICI) {
    ss << to_string_rrrici(instruction);
  } else if (suffix == abi::instruction::ZRRI) {
    ss << to_string_zrri(instruction);
  } else if (suffix == abi::instruction::ZRRICI) {
    ss << to_string_zrrici(instruction);
  } else if (suffix == abi::instruction::S_RRRI or
             suffix == abi::instruction::U_RRRI) {
    ss << to_string_s_rrri(instruction);
  } else if (suffix == abi::instruction::S_RRRICI or
             suffix == abi::instruction::U_RRRICI) {
    ss << to_string_s_rrici(instruction);
  } else if (suffix == abi::instruction::RIR) {
    ss << to_string_rir(instruction);
  } else if (suffix == abi::instruction::RIRC) {
    ss << to_string_rirc(instruction);
  } else if (suffix == abi::instruction::RIRCI) {
    ss << to_string_rirci(instruction);
  } else if (suffix == abi::instruction::ZIR) {
    ss << to_string_zir(instruction);
  } else if (suffix == abi::instruction::ZIRC) {
    ss << to_string_zirc(instruction);
  } else if (suffix == abi::instruction::ZIRCI) {
    ss << to_string_zirci(instruction);
  } else if (suffix == abi::instruction::S_RIRC or
             suffix == abi::instruction::U_RIRC) {
    ss << to_string_rirc(instruction);
  } else if (suffix == abi::instruction::S_RIRCI or
             suffix == abi::instruction::U_RIRCI) {
    ss << to_string_rirci(instruction);
  } else if (suffix == abi::instruction::R) {
    ss << to_string_r(instruction);
  } else if (suffix == abi::instruction::RCI) {
    ss << to_string_rci(instruction);
  } else if (suffix == abi::instruction::Z) {
    ss << to_string_z(instruction);
  } else if (suffix == abi::instruction::ZCI) {
    ss << to_string_zci(instruction);
  } else if (suffix == abi::instruction::S_R or
             suffix == abi::instruction::U_R) {
    ss << to_string_s_r(instruction);
  } else if (suffix == abi::instruction::S_RCI or
             suffix == abi::instruction::U_RCI) {
    ss << to_string_s_rci(instruction);
  } else if (suffix == abi::instruction::CI) {
    ss << to_string_ci(instruction);
  } else if (suffix == abi::instruction::I) {
    ss << to_string_i(instruction);
  } else if (suffix == abi::instruction::DDCI) {
    ss << to_string_ddci(instruction);
  } else if (suffix == abi::instruction::ERRI) {
    ss << to_string_erri(instruction);
  } else if (suffix == abi::instruction::S_ERRI or
             suffix == abi::instruction::U_ERRI) {
    ss << to_string_s_erri(instruction);
  } else if (suffix == abi::instruction::EDRI) {
    ss << to_string_edri(instruction);
  } else if (suffix == abi::instruction::ERII) {
    ss << to_string_erii(instruction);
  } else if (suffix == abi::instruction::ERIR) {
    ss << to_string_erir(instruction);
  } else if (suffix == abi::instruction::ERID) {
    ss << to_string_erid(instruction);
  } else if (suffix == abi::instruction::DMA_RRI) {
    ss << to_string_dma_rri(instruction);
  } else {
    throw std::invalid_argument("");
  }

  return ss.str();
}

std::string InstructionConverter::to_string_rici(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << instruction->imm()->value() << ", ";
  ss << ConditionConverter::to_string(instruction->condition()) << ", ";
  ss << instruction->pc()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_rri(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->rc()) << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << instruction->imm()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_rric(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->rc()) << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << instruction->imm()->value() << ", ";
  ss << ConditionConverter::to_string(instruction->condition());
  return ss.str();
}

std::string InstructionConverter::to_string_rrici(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->rc()) << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << instruction->imm()->value() << ", ";
  ss << ConditionConverter::to_string(instruction->condition()) << ", ";
  ss << instruction->pc()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_rrif(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->rc()) << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << instruction->imm()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_rrr(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->rc()) << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << RegConverter::to_string(instruction->rb());
  return ss.str();
}

std::string InstructionConverter::to_string_rrrc(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->rc()) << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << RegConverter::to_string(instruction->rb()) << ", ";
  ss << ConditionConverter::to_string(instruction->condition());
  return ss.str();
}

std::string InstructionConverter::to_string_rrrci(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->rc()) << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << RegConverter::to_string(instruction->rb()) << ", ";
  ss << ConditionConverter::to_string(instruction->condition()) << ", ";
  ss << instruction->pc()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_zri(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << instruction->imm()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_zric(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << instruction->imm()->value() << ", ";
  ss << ConditionConverter::to_string(instruction->condition());
  return ss.str();
}

std::string InstructionConverter::to_string_zrici(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << instruction->imm()->value() << ", ";
  ss << ConditionConverter::to_string(instruction->condition()) << ", ";
  ss << instruction->pc()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_zrif(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << instruction->imm()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_zrr(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << RegConverter::to_string(instruction->rb());
  return ss.str();
}

std::string InstructionConverter::to_string_zrrc(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << RegConverter::to_string(instruction->rb()) << ", ";
  ss << ConditionConverter::to_string(instruction->condition());
  return ss.str();
}

std::string InstructionConverter::to_string_zrrci(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << RegConverter::to_string(instruction->rb()) << ", ";
  ss << ConditionConverter::to_string(instruction->condition()) << ", ";
  ss << instruction->pc()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_s_rri(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->dc()) << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << instruction->imm()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_s_rric(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->dc()) << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << instruction->imm()->value() << ", ";
  ss << ConditionConverter::to_string(instruction->condition());
  return ss.str();
}

std::string InstructionConverter::to_string_s_rrici(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->dc()) << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << instruction->imm()->value() << ", ";
  ss << ConditionConverter::to_string(instruction->condition()) << ", ";
  ss << instruction->pc()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_s_rrif(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->dc()) << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << instruction->imm()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_s_rrr(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->dc()) << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << RegConverter::to_string(instruction->rb());
  return ss.str();
}

std::string InstructionConverter::to_string_s_rrrc(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->dc()) << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << RegConverter::to_string(instruction->rb()) << ", ";
  ss << ConditionConverter::to_string(instruction->condition());
  return ss.str();
}

std::string InstructionConverter::to_string_s_rrrci(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->dc()) << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << RegConverter::to_string(instruction->rb()) << ", ";
  ss << ConditionConverter::to_string(instruction->condition()) << ", ";
  ss << instruction->pc()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_rr(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->rc()) << ", ";
  ss << RegConverter::to_string(instruction->ra());
  return ss.str();
}

std::string InstructionConverter::to_string_rrc(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->rc()) << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << ConditionConverter::to_string(instruction->condition());
  return ss.str();
}

std::string InstructionConverter::to_string_rrci(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->rc()) << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << ConditionConverter::to_string(instruction->condition()) << ", ";
  ss << instruction->pc()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_zr(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->ra());
  return ss.str();
}

std::string InstructionConverter::to_string_zrc(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << ConditionConverter::to_string(instruction->condition());
  return ss.str();
}

std::string InstructionConverter::to_string_zrci(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << ConditionConverter::to_string(instruction->condition()) << ", ";
  ss << instruction->pc()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_s_rr(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->dc()) << ", ";
  ss << RegConverter::to_string(instruction->ra());
  return ss.str();
}

std::string InstructionConverter::to_string_s_rrc(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->dc()) << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << ConditionConverter::to_string(instruction->condition());
  return ss.str();
}

std::string InstructionConverter::to_string_s_rrci(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->dc()) << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << ConditionConverter::to_string(instruction->condition()) << ", ";
  ss << instruction->pc()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_drdici(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->dc()) << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << RegConverter::to_string(instruction->db()) << ", ";
  ss << instruction->imm()->value() << ", ";
  ss << ConditionConverter::to_string(instruction->condition()) << ", ";
  ss << instruction->pc()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_rrri(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->rc()) << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << RegConverter::to_string(instruction->rb()) << ", ";
  ss << instruction->imm()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_rrrici(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->rc()) << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << RegConverter::to_string(instruction->rb()) << ", ";
  ss << instruction->imm()->value() << ", ";
  ss << ConditionConverter::to_string(instruction->condition()) << ", ";
  ss << instruction->pc()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_zrri(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << RegConverter::to_string(instruction->rb()) << ", ";
  ss << instruction->imm()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_zrrici(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << RegConverter::to_string(instruction->rb()) << ", ";
  ss << instruction->imm()->value() << ", ";
  ss << ConditionConverter::to_string(instruction->condition()) << ", ";
  ss << instruction->pc()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_s_rrri(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->dc()) << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << RegConverter::to_string(instruction->rb()) << ", ";
  ss << instruction->imm()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_s_rrrici(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->dc()) << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << RegConverter::to_string(instruction->rb()) << ", ";
  ss << instruction->imm()->value() << ", ";
  ss << ConditionConverter::to_string(instruction->condition()) << ", ";
  ss << instruction->pc()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_rir(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->rc()) << ", ";
  ss << instruction->imm()->value() << ", ";
  ss << RegConverter::to_string(instruction->ra());
  return ss.str();
}

std::string InstructionConverter::to_string_rirc(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->rc()) << ", ";
  ss << instruction->imm()->value() << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << ConditionConverter::to_string(instruction->condition());
  return ss.str();
}

std::string InstructionConverter::to_string_rirci(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->rc()) << ", ";
  ss << instruction->imm()->value() << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << ConditionConverter::to_string(instruction->condition()) << ", ";
  ss << instruction->pc()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_zir(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << instruction->imm()->value() << ", ";
  ss << RegConverter::to_string(instruction->ra());
  return ss.str();
}

std::string InstructionConverter::to_string_zirc(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << instruction->imm()->value() << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << ConditionConverter::to_string(instruction->condition());
  return ss.str();
}

std::string InstructionConverter::to_string_zirci(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << instruction->imm()->value() << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << ConditionConverter::to_string(instruction->condition()) << ", ";
  ss << instruction->pc()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_s_rirc(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->dc()) << ", ";
  ss << instruction->imm()->value() << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << ConditionConverter::to_string(instruction->condition());
  return ss.str();
}

std::string InstructionConverter::to_string_s_rirci(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->dc()) << ", ";
  ss << instruction->imm()->value() << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << ConditionConverter::to_string(instruction->condition()) << ", ";
  ss << instruction->pc()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_r(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->rc());
  return ss.str();
}

std::string InstructionConverter::to_string_rci(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->rc()) << ", ";
  ss << ConditionConverter::to_string(instruction->condition()) << ", ";
  ss << instruction->pc()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_z(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  return ss.str();
}

std::string InstructionConverter::to_string_zci(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << ConditionConverter::to_string(instruction->condition()) << ", ";
  ss << instruction->pc()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_s_r(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->dc());
  return ss.str();
}

std::string InstructionConverter::to_string_s_rci(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->dc()) << ", ";
  ss << ConditionConverter::to_string(instruction->condition()) << ", ";
  ss << instruction->pc()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_ci(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << ConditionConverter::to_string(instruction->condition()) << ", ";
  ss << instruction->pc()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_i(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << instruction->imm()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_ddci(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->dc()) << ", ";
  ss << RegConverter::to_string(instruction->db()) << ", ";
  ss << ConditionConverter::to_string(instruction->condition()) << ", ";
  ss << instruction->pc()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_erri(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << EndianConverter::to_string(instruction->endian()) << ", ";
  ss << RegConverter::to_string(instruction->rc()) << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << instruction->off()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_s_erri(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << EndianConverter::to_string(instruction->endian()) << ", ";
  ss << RegConverter::to_string(instruction->dc()) << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << instruction->off()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_edri(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << EndianConverter::to_string(instruction->endian()) << ", ";
  ss << RegConverter::to_string(instruction->dc()) << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << instruction->off()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_erii(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << EndianConverter::to_string(instruction->endian()) << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << instruction->off()->value() << ", ";
  ss << instruction->imm()->value();
  return ss.str();
}

std::string InstructionConverter::to_string_erir(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << EndianConverter::to_string(instruction->endian()) << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << instruction->off()->value() << ", ";
  ss << RegConverter::to_string(instruction->rb());
  return ss.str();
}

std::string InstructionConverter::to_string_erid(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << EndianConverter::to_string(instruction->endian()) << ", ";
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << instruction->off()->value() << ", ";
  ss << RegConverter::to_string(instruction->db());
  return ss.str();
}

std::string InstructionConverter::to_string_dma_rri(
    abi::instruction::Instruction *instruction) {
  std::stringstream ss;
  ss << RegConverter::to_string(instruction->ra()) << ", ";
  ss << RegConverter::to_string(instruction->rb()) << ", ";
  ss << instruction->imm()->value();
  return ss.str();
}

}  // namespace upmem_sim::converter
