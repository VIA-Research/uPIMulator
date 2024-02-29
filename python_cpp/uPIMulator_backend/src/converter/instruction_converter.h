#ifndef UPMEM_SIM_CONVERTER_INSTRUCTION_CONVERTER_H_
#define UPMEM_SIM_CONVERTER_INSTRUCTION_CONVERTER_H_

#include <string>

#include "abi/instruction/instruction.h"

namespace upmem_sim::converter {

class InstructionConverter {
 public:
  static std::string to_string(abi::instruction::Instruction *instruction);

 protected:
  static std::string to_string_rici(abi::instruction::Instruction *instruction);

  static std::string to_string_rri(abi::instruction::Instruction *instruction);
  static std::string to_string_rric(abi::instruction::Instruction *instruction);
  static std::string to_string_rrici(
      abi::instruction::Instruction *instruction);
  static std::string to_string_rrif(abi::instruction::Instruction *instruction);
  static std::string to_string_rrr(abi::instruction::Instruction *instruction);
  static std::string to_string_rrrc(abi::instruction::Instruction *instruction);
  static std::string to_string_rrrci(
      abi::instruction::Instruction *instruction);

  static std::string to_string_zri(abi::instruction::Instruction *instruction);
  static std::string to_string_zric(abi::instruction::Instruction *instruction);
  static std::string to_string_zrici(
      abi::instruction::Instruction *instruction);
  static std::string to_string_zrif(abi::instruction::Instruction *instruction);
  static std::string to_string_zrr(abi::instruction::Instruction *instruction);
  static std::string to_string_zrrc(abi::instruction::Instruction *instruction);
  static std::string to_string_zrrci(
      abi::instruction::Instruction *instruction);

  static std::string to_string_s_rri(
      abi::instruction::Instruction *instruction);
  static std::string to_string_s_rric(
      abi::instruction::Instruction *instruction);
  static std::string to_string_s_rrici(
      abi::instruction::Instruction *instruction);
  static std::string to_string_s_rrif(
      abi::instruction::Instruction *instruction);
  static std::string to_string_s_rrr(
      abi::instruction::Instruction *instruction);
  static std::string to_string_s_rrrc(
      abi::instruction::Instruction *instruction);
  static std::string to_string_s_rrrci(
      abi::instruction::Instruction *instruction);

  static std::string to_string_rr(abi::instruction::Instruction *instruction);
  static std::string to_string_rrc(abi::instruction::Instruction *instruction);
  static std::string to_string_rrci(abi::instruction::Instruction *instruction);

  static std::string to_string_zr(abi::instruction::Instruction *instruction);
  static std::string to_string_zrc(abi::instruction::Instruction *instruction);
  static std::string to_string_zrci(abi::instruction::Instruction *instruction);

  static std::string to_string_s_rr(abi::instruction::Instruction *instruction);
  static std::string to_string_s_rrc(
      abi::instruction::Instruction *instruction);
  static std::string to_string_s_rrci(
      abi::instruction::Instruction *instruction);

  static std::string to_string_drdici(
      abi::instruction::Instruction *instruction);

  static std::string to_string_rrri(abi::instruction::Instruction *instruction);
  static std::string to_string_rrrici(
      abi::instruction::Instruction *instruction);

  static std::string to_string_zrri(abi::instruction::Instruction *instruction);
  static std::string to_string_zrrici(
      abi::instruction::Instruction *instruction);

  static std::string to_string_s_rrri(
      abi::instruction::Instruction *instruction);
  static std::string to_string_s_rrrici(
      abi::instruction::Instruction *instruction);

  static std::string to_string_rir(abi::instruction::Instruction *instruction);
  static std::string to_string_rirc(abi::instruction::Instruction *instruction);
  static std::string to_string_rirci(
      abi::instruction::Instruction *instruction);

  static std::string to_string_zir(abi::instruction::Instruction *instruction);
  static std::string to_string_zirc(abi::instruction::Instruction *instruction);
  static std::string to_string_zirci(
      abi::instruction::Instruction *instruction);

  static std::string to_string_s_rirc(
      abi::instruction::Instruction *instruction);
  static std::string to_string_s_rirci(
      abi::instruction::Instruction *instruction);

  static std::string to_string_r(abi::instruction::Instruction *instruction);
  static std::string to_string_rci(abi::instruction::Instruction *instruction);

  static std::string to_string_z(abi::instruction::Instruction *instruction);
  static std::string to_string_zci(abi::instruction::Instruction *instruction);

  static std::string to_string_s_r(abi::instruction::Instruction *instruction);
  static std::string to_string_s_rci(
      abi::instruction::Instruction *instruction);

  static std::string to_string_ci(abi::instruction::Instruction *instruction);
  static std::string to_string_i(abi::instruction::Instruction *instruction);

  static std::string to_string_ddci(abi::instruction::Instruction *instruction);

  static std::string to_string_erri(abi::instruction::Instruction *instruction);

  static std::string to_string_s_erri(
      abi::instruction::Instruction *instruction);

  static std::string to_string_edri(abi::instruction::Instruction *instruction);

  static std::string to_string_erii(abi::instruction::Instruction *instruction);
  static std::string to_string_erir(abi::instruction::Instruction *instruction);
  static std::string to_string_erid(abi::instruction::Instruction *instruction);

  static std::string to_string_dma_rri(
      abi::instruction::Instruction *instruction);
};

}  // namespace upmem_sim::converter

#endif
