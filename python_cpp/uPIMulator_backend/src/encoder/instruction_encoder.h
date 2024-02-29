#ifndef UPMEM_SIM_ENCODER_INSTRUCTION_ENCODER_H_
#define UPMEM_SIM_ENCODER_INSTRUCTION_ENCODER_H_

#include <cmath>

#include "abi/instruction/instruction.h"
#include "abi/instruction/op_code.h"
#include "abi/instruction/suffix.h"
#include "abi/word/instruction_word.h"
#include "encoder/byte_stream.h"
#include "util/config_loader.h"

namespace upmem_sim::encoder {

class InstructionEncoder {
 public:
  static abi::instruction::Instruction *decode(ByteStream *byte_stream);

 protected:
  static abi::instruction::OpCode decode_op_code(
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Suffix decode_suffix(
      abi::word::InstructionWord *instruction_word);
  static abi::reg::GPReg *decode_gp_reg(
      abi::word::InstructionWord *instruction_word, int begin, int end);
  static abi::reg::SrcReg *decode_src_reg(
      abi::word::InstructionWord *instruction_word, int begin, int end);
  static abi::reg::PairReg *decode_pair_reg(
      abi::word::InstructionWord *instruction_word, int begin, int end);
  static int64_t decode_imm(abi::word::InstructionWord *instruction_word,
                            int begin, int end,
                            abi::word::Representation representation);
  static int64_t decode_off(abi::word::InstructionWord *instruction_word,
                            int begin, int end,
                            abi::word::Representation representation);
  static abi::isa::Condition decode_condition(
      abi::word::InstructionWord *instruction_word, int begin, int end);
  static int64_t decode_pc(abi::word::InstructionWord *instruction_word,
                           int begin, int end);
  static abi::isa::Endian decode_endian(
      abi::word::InstructionWord *instruction_word, int begin, int end);

  static abi::instruction::Instruction *decode_rici(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_rri(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_rric(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_rrici(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_rrif(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_rrr(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_rrrc(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_rrrci(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_zri(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_zric(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_zrici(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_zrif(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_zrr(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_zrrc(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_zrrci(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_s_rri(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_s_rric(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_s_rrici(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_s_rrif(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_s_rrr(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_s_rrrc(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_s_rrrci(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_rr(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_rrc(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_rrci(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_zr(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_zrc(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_zrci(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_s_rr(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_s_rrc(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_s_rrci(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_drdici(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_rrri(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_rrrici(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_zrri(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_zrrici(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_s_rrri(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_s_rrrici(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_rir(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_rirc(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_rirci(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_zir(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_zirc(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_zirci(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_s_rirc(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_s_rirci(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_r(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_rci(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_z(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_zci(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_s_r(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_s_rci(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_ci(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_i(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_ddci(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_erri(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_s_erri(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_edri(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_erii(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_erir(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_erid(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);
  static abi::instruction::Instruction *decode_dma_rri(
      abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
      abi::word::InstructionWord *instruction_word);

  static int op_code_begin() { return 0; }
  static int op_code_end() { return op_code_begin() + op_code_width(); }

  static int suffix_begin() { return op_code_end(); }
  static int suffix_end() { return suffix_begin() + suffix_width(); }

  static int op_code_width() {
    return ceil(log2(1.0 + abi::instruction::SDMA));
  }
  static int suffix_width() {
    return ceil(log2(1.0 + abi::instruction::DMA_RRI));
  }
  static int register_width() {
    return ceil(log2(util::ConfigLoader::num_gp_registers() + abi::reg::ID8));
  }
  static int condition_width() { return ceil(log2(1.0 + abi::isa::LARGE)); }
  static int pc_width() { return util::ConfigLoader::iram_address_width(); }
  static int endian_width() { return ceil(log2(1.0 + abi::isa::BIG)); }
};

}  // namespace upmem_sim::encoder

#endif
