#include "encoder/instruction_encoder.h"

#include <stdexcept>

namespace upmem_sim::encoder {

abi::instruction::Instruction *InstructionEncoder::decode(
    ByteStream *byte_stream) {
  auto instruction_word = new abi::word::InstructionWord();
  instruction_word->from_byte_stream(byte_stream);

  abi::instruction::OpCode op_code = decode_op_code(instruction_word);
  abi::instruction::Suffix suffix = decode_suffix(instruction_word);

  abi::instruction::Instruction *instruction = nullptr;
  if (suffix == abi::instruction::RICI) {
    instruction = decode_rici(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::RRI) {
    instruction = decode_rri(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::RRIC) {
    instruction = decode_rric(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::RRICI) {
    instruction = decode_rrici(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::RRIF) {
    instruction = decode_rrif(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::RRR) {
    instruction = decode_rrr(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::RRRC) {
    instruction = decode_rrrc(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::RRRCI) {
    instruction = decode_rrrci(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::ZRI) {
    instruction = decode_zri(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::ZRIC) {
    instruction = decode_zric(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::ZRICI) {
    instruction = decode_zrici(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::ZRIF) {
    instruction = decode_zrif(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::ZRR) {
    instruction = decode_zrr(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::ZRRC) {
    instruction = decode_zrrc(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::ZRRCI) {
    instruction = decode_zrrci(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::S_RRI or
             suffix == abi::instruction::U_RRI) {
    instruction = decode_s_rri(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::S_RRIC or
             suffix == abi::instruction::U_RRIC) {
    instruction = decode_s_rric(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::S_RRICI or
             suffix == abi::instruction::U_RRICI) {
    instruction = decode_s_rrici(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::S_RRIF or
             suffix == abi::instruction::U_RRIF) {
    instruction = decode_s_rrif(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::S_RRR or
             suffix == abi::instruction::U_RRR) {
    instruction = decode_s_rrr(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::S_RRRC or
             suffix == abi::instruction::U_RRRC) {
    instruction = decode_s_rrrc(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::S_RRRCI or
             suffix == abi::instruction::U_RRRCI) {
    instruction = decode_s_rrrci(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::RR) {
    instruction = decode_rr(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::RRC) {
    instruction = decode_rrc(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::RRCI) {
    instruction = decode_rrci(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::ZR) {
    instruction = decode_zr(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::ZRC) {
    instruction = decode_zrc(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::ZRCI) {
    instruction = decode_zrci(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::S_RR or
             suffix == abi::instruction::U_RR) {
    instruction = decode_s_rr(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::S_RRC or
             suffix == abi::instruction::U_RRC) {
    instruction = decode_s_rrc(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::S_RRCI or
             suffix == abi::instruction::U_RRCI) {
    instruction = decode_s_rrci(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::DRDICI) {
    instruction = decode_drdici(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::RRRI) {
    instruction = decode_rrri(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::RRRICI) {
    instruction = decode_rrrici(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::ZRRI) {
    instruction = decode_zrri(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::ZRRICI) {
    instruction = decode_zrrici(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::S_RRRI or
             suffix == abi::instruction::U_RRRI) {
    instruction = decode_s_rrri(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::S_RRRICI or
             suffix == abi::instruction::U_RRRICI) {
    instruction = decode_s_rrrici(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::RIR) {
    instruction = decode_rir(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::RIRC) {
    instruction = decode_rirc(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::RIRCI) {
    instruction = decode_rirci(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::ZIR) {
    instruction = decode_zir(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::ZIRC) {
    instruction = decode_zirc(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::ZIRCI) {
    instruction = decode_zirci(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::S_RIRC or
             suffix == abi::instruction::U_RIRC) {
    instruction = decode_s_rirc(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::S_RIRCI or
             suffix == abi::instruction::U_RIRCI) {
    instruction = decode_s_rirci(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::R) {
    instruction = decode_r(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::RCI) {
    instruction = decode_rci(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::Z) {
    instruction = decode_z(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::ZCI) {
    instruction = decode_zci(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::S_R or
             suffix == abi::instruction::U_R) {
    instruction = decode_s_r(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::S_RCI or
             suffix == abi::instruction::U_RCI) {
    instruction = decode_s_rci(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::CI) {
    instruction = decode_ci(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::I) {
    instruction = decode_i(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::DDCI) {
    instruction = decode_ddci(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::ERRI) {
    instruction = decode_erri(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::S_ERRI or
             suffix == abi::instruction::U_ERRI) {
    instruction = decode_s_erri(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::EDRI) {
    instruction = decode_edri(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::ERII) {
    instruction = decode_erii(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::ERIR) {
    instruction = decode_erir(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::ERID) {
    instruction = decode_erid(op_code, suffix, instruction_word);
  } else if (suffix == abi::instruction::DMA_RRI) {
    instruction = decode_dma_rri(op_code, suffix, instruction_word);
  } else {
    throw std::invalid_argument("");
  }

  delete instruction_word;
  return instruction;
}

abi::instruction::OpCode InstructionEncoder::decode_op_code(
    abi::word::InstructionWord *instruction_word) {
  return static_cast<abi::instruction::OpCode>(instruction_word->bit_slice(
      abi::word::UNSIGNED, op_code_begin(), op_code_end()));
}

abi::instruction::Suffix InstructionEncoder::decode_suffix(
    abi::word::InstructionWord *instruction_word) {
  return static_cast<abi::instruction::Suffix>(instruction_word->bit_slice(
      abi::word::UNSIGNED, suffix_begin(), suffix_end()));
}

abi::reg::GPReg *InstructionEncoder::decode_gp_reg(
    abi::word::InstructionWord *instruction_word, int begin, int end) {
  auto index = static_cast<RegIndex>(
      instruction_word->bit_slice(abi::word::UNSIGNED, begin, end));
  return new abi::reg::GPReg(index);
}

abi::reg::SrcReg *InstructionEncoder::decode_src_reg(
    abi::word::InstructionWord *instruction_word, int begin, int end) {
  auto index = static_cast<RegIndex>(
      instruction_word->bit_slice(abi::word::UNSIGNED, begin, end));
  if (index < util::ConfigLoader::num_gp_registers()) {
    return new abi::reg::SrcReg(new abi::reg::GPReg(index));
  } else {
    return new abi::reg::SrcReg(
        new abi::reg::SPReg(static_cast<abi::reg::SPReg>(
            index - util::ConfigLoader::num_gp_registers())));
  }
}

abi::reg::PairReg *InstructionEncoder::decode_pair_reg(
    abi::word::InstructionWord *instruction_word, int begin, int end) {
  auto index = static_cast<RegIndex>(
      instruction_word->bit_slice(abi::word::UNSIGNED, begin, end));
  return new abi::reg::PairReg(index);
}

int64_t InstructionEncoder::decode_imm(
    abi::word::InstructionWord *instruction_word, int begin, int end,
    abi::word::Representation representation) {
  return instruction_word->bit_slice(representation, begin, end);
}

int64_t InstructionEncoder::decode_off(
    abi::word::InstructionWord *instruction_word, int begin, int end,
    abi::word::Representation representation) {
  return decode_imm(instruction_word, begin, end, representation);
}

abi::isa::Condition InstructionEncoder::decode_condition(
    abi::word::InstructionWord *instruction_word, int begin, int end) {
  return static_cast<abi::isa::Condition>(
      instruction_word->bit_slice(abi::word::UNSIGNED, begin, end));
}

int64_t InstructionEncoder::decode_pc(
    abi::word::InstructionWord *instruction_word, int begin, int end) {
  return decode_imm(instruction_word, begin, end, abi::word::UNSIGNED);
}

abi::isa::Endian InstructionEncoder::decode_endian(
    abi::word::InstructionWord *instruction_word, int begin, int end) {
  return static_cast<abi::isa::Endian>(
      instruction_word->bit_slice(abi::word::UNSIGNED, begin, end));
}

abi::instruction::Instruction *InstructionEncoder::decode_rici(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rici_op_codes().count(op_code));
  assert(suffix == abi::instruction::RICI);

  int ra_begin = suffix_end();
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int imm_begin = ra_end;
  int imm_end = imm_begin + 16;
  int64_t imm =
      decode_imm(instruction_word, imm_begin, imm_end, abi::word::SIGNED);

  int condition_begin = imm_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  int pc_begin = condition_end;
  int pc_end = pc_begin + pc_width();
  int64_t pc = decode_pc(instruction_word, pc_begin, pc_end);

  return new abi::instruction::Instruction(op_code, suffix, ra, imm, condition,
                                           pc);
}

abi::instruction::Instruction *InstructionEncoder::decode_rri(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rri_op_codes().count(op_code));
  assert(suffix == abi::instruction::RRI);

  int rc_begin = suffix_end();
  int rc_end = rc_begin + register_width();
  abi::reg::GPReg *rc = decode_gp_reg(instruction_word, rc_begin, rc_end);

  int ra_begin = rc_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int imm_begin = ra_end;
  int imm_end;
  int64_t imm;
  if (abi::instruction::Instruction::add_rri_op_codes().count(op_code)) {
    imm_end = imm_begin + 32;
    imm = decode_imm(instruction_word, imm_begin, imm_end, abi::word::UNSIGNED);
  } else if (abi::instruction::Instruction::asr_rri_op_codes().count(op_code)) {
    imm_end = imm_begin + 5;
    imm = decode_imm(instruction_word, imm_begin, imm_end, abi::word::UNSIGNED);
  } else if (abi::instruction::Instruction::call_rri_op_codes().count(
                 op_code)) {
    imm_end = imm_begin + 24;
    imm = decode_imm(instruction_word, imm_begin, imm_end, abi::word::SIGNED);
  } else {
    throw std::invalid_argument("");
  }

  return new abi::instruction::Instruction(op_code, suffix, rc, ra, imm);
}

abi::instruction::Instruction *InstructionEncoder::decode_rric(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rric_op_codes().count(op_code));
  assert(suffix == abi::instruction::RRIC);

  int rc_begin = suffix_end();
  int rc_end = rc_begin + register_width();
  abi::reg::GPReg *rc = decode_gp_reg(instruction_word, rc_begin, rc_end);

  int ra_begin = rc_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int imm_begin = ra_end;
  int imm_end;
  int64_t imm;
  if (abi::instruction::Instruction::add_rric_op_codes().count(op_code) or
      abi::instruction::Instruction::sub_rric_op_codes().count(op_code)) {
    imm_end = imm_begin + 24;
    imm = decode_imm(instruction_word, imm_begin, imm_end, abi::word::UNSIGNED);
  } else if (abi::instruction::Instruction::asr_rri_op_codes().count(op_code)) {
    imm_end = imm_begin + 5;
    imm = decode_imm(instruction_word, imm_begin, imm_end, abi::word::UNSIGNED);
  } else {
    throw std::invalid_argument("");
  }

  int condition_begin = imm_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  return new abi::instruction::Instruction(op_code, suffix, rc, ra, imm,
                                           condition);
}

abi::instruction::Instruction *InstructionEncoder::decode_rrici(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rrici_op_codes().count(op_code));
  assert(suffix == abi::instruction::RRICI);

  int rc_begin = suffix_end();
  int rc_end = rc_begin + register_width();
  abi::reg::GPReg *rc = decode_gp_reg(instruction_word, rc_begin, rc_end);

  int ra_begin = rc_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int imm_begin = ra_end;
  int imm_end;
  int64_t imm;
  if (abi::instruction::Instruction::add_rrici_op_codes().count(op_code) or
      abi::instruction::Instruction::and_rrici_op_codes().count(op_code) or
      abi::instruction::Instruction::sub_rrici_op_codes().count(op_code)) {
    imm_end = imm_begin + 8;
    imm = decode_imm(instruction_word, imm_begin, imm_end, abi::word::UNSIGNED);
  } else if (abi::instruction::Instruction::asr_rrici_op_codes().count(
                 op_code)) {
    imm_end = imm_begin + 5;
    imm = decode_imm(instruction_word, imm_begin, imm_end, abi::word::UNSIGNED);
  } else {
    throw std::invalid_argument("");
  }

  int condition_begin = imm_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  int pc_begin = condition_end;
  int pc_end = pc_begin + pc_width();
  int64_t pc = decode_pc(instruction_word, pc_begin, pc_end);

  return new abi::instruction::Instruction(op_code, suffix, rc, ra, imm,
                                           condition, pc);
}

abi::instruction::Instruction *InstructionEncoder::decode_rrif(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rrif_op_codes().count(op_code));
  assert(suffix == abi::instruction::RRIF);

  int rc_begin = suffix_end();
  int rc_end = rc_begin + register_width();
  abi::reg::GPReg *rc = decode_gp_reg(instruction_word, rc_begin, rc_end);

  int ra_begin = rc_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int imm_begin = ra_end;
  int imm_end = imm_begin + 24;
  int64_t imm =
      decode_imm(instruction_word, imm_begin, imm_end, abi::word::SIGNED);

  int condition_begin = imm_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  return new abi::instruction::Instruction(op_code, suffix, rc, ra, imm,
                                           condition);
}

abi::instruction::Instruction *InstructionEncoder::decode_rrr(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rrr_op_codes().count(op_code));
  assert(suffix == abi::instruction::RRR);

  int rc_begin = suffix_end();
  int rc_end = rc_begin + register_width();
  abi::reg::GPReg *rc = decode_gp_reg(instruction_word, rc_begin, rc_end);

  int ra_begin = rc_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int rb_begin = ra_end;
  int rb_end = rb_begin + register_width();
  abi::reg::SrcReg *rb = decode_src_reg(instruction_word, rb_begin, rb_end);

  return new abi::instruction::Instruction(op_code, suffix, rc, ra, rb);
}

abi::instruction::Instruction *InstructionEncoder::decode_rrrc(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rrrc_op_codes().count(op_code));
  assert(suffix == abi::instruction::RRRC);

  int rc_begin = suffix_end();
  int rc_end = rc_begin + register_width();
  abi::reg::GPReg *rc = decode_gp_reg(instruction_word, rc_begin, rc_end);

  int ra_begin = rc_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int rb_begin = ra_end;
  int rb_end = rb_begin + register_width();
  abi::reg::SrcReg *rb = decode_src_reg(instruction_word, rb_begin, rb_end);

  int condition_begin = rb_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  return new abi::instruction::Instruction(op_code, suffix, rc, ra, rb,
                                           condition);
}

abi::instruction::Instruction *InstructionEncoder::decode_rrrci(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rrrci_op_codes().count(op_code));
  assert(suffix == abi::instruction::RRRCI);

  int rc_begin = suffix_end();
  int rc_end = rc_begin + register_width();
  abi::reg::GPReg *rc = decode_gp_reg(instruction_word, rc_begin, rc_end);

  int ra_begin = rc_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int rb_begin = ra_end;
  int rb_end = rb_begin + register_width();
  abi::reg::SrcReg *rb = decode_src_reg(instruction_word, rb_begin, rb_end);

  int condition_begin = rb_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  int pc_begin = condition_end;
  int pc_end = pc_begin + pc_width();
  int64_t pc = decode_pc(instruction_word, pc_begin, pc_end);

  return new abi::instruction::Instruction(op_code, suffix, rc, ra, rb,
                                           condition, pc);
}

abi::instruction::Instruction *InstructionEncoder::decode_zri(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rri_op_codes().count(op_code));
  assert(suffix == abi::instruction::ZRI);

  int ra_begin = suffix_end();
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int imm_begin = ra_end;
  int imm_end;
  int64_t imm;
  if (abi::instruction::Instruction::add_rri_op_codes().count(op_code)) {
    imm_end = imm_begin + 32;
    imm = decode_imm(instruction_word, imm_begin, imm_end, abi::word::UNSIGNED);
  } else if (abi::instruction::Instruction::asr_rri_op_codes().count(op_code)) {
    imm_end = imm_begin + 5;
    imm = decode_imm(instruction_word, imm_begin, imm_end, abi::word::UNSIGNED);
  } else if (abi::instruction::Instruction::call_rri_op_codes().count(
                 op_code)) {
    imm_end = imm_begin + 28;
    imm = decode_imm(instruction_word, imm_begin, imm_end, abi::word::SIGNED);
  } else {
    throw std::invalid_argument("");
  }

  return new abi::instruction::Instruction(op_code, suffix, ra, imm);
}

abi::instruction::Instruction *InstructionEncoder::decode_zric(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rric_op_codes().count(op_code));
  assert(suffix == abi::instruction::ZRIC);

  int ra_begin = suffix_end();
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int imm_begin = ra_end;
  int imm_end;
  int64_t imm;
  if (abi::instruction::Instruction::add_rric_op_codes().count(op_code) or
      abi::instruction::Instruction::sub_rric_op_codes().count(op_code)) {
    imm_end = imm_begin + 27;
    imm = decode_imm(instruction_word, imm_begin, imm_end, abi::word::UNSIGNED);
  } else if (abi::instruction::Instruction::asr_rri_op_codes().count(op_code)) {
    imm_end = imm_begin + 5;
    imm = decode_imm(instruction_word, imm_begin, imm_end, abi::word::UNSIGNED);
  } else {
    throw std::invalid_argument("");
  }

  int condition_begin = imm_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  return new abi::instruction::Instruction(op_code, suffix, ra, imm, condition);
}

abi::instruction::Instruction *InstructionEncoder::decode_zrici(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rrici_op_codes().count(op_code));
  assert(suffix == abi::instruction::ZRICI);

  int ra_begin = suffix_end();
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int imm_begin = ra_end;
  int imm_end;
  int64_t imm;
  if (abi::instruction::Instruction::add_rrici_op_codes().count(op_code) or
      abi::instruction::Instruction::and_rrici_op_codes().count(op_code) or
      abi::instruction::Instruction::sub_rrici_op_codes().count(op_code)) {
    imm_end = imm_begin + 11;
    imm = decode_imm(instruction_word, imm_begin, imm_end, abi::word::UNSIGNED);
  } else if (abi::instruction::Instruction::asr_rrici_op_codes().count(
                 op_code)) {
    imm_end = imm_begin + 5;
    imm = decode_imm(instruction_word, imm_begin, imm_end, abi::word::UNSIGNED);
  } else {
    throw std::invalid_argument("");
  }

  int condition_begin = imm_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  int pc_begin = condition_end;
  int pc_end = pc_begin + pc_width();
  int64_t pc = decode_pc(instruction_word, pc_begin, pc_end);

  return new abi::instruction::Instruction(op_code, suffix, ra, imm, condition,
                                           pc);
}

abi::instruction::Instruction *InstructionEncoder::decode_zrif(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rrif_op_codes().count(op_code));
  assert(suffix == abi::instruction::ZRIF);

  int ra_begin = suffix_end();
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int imm_begin = ra_end;
  int imm_end = imm_begin + 24;
  int64_t imm =
      decode_imm(instruction_word, imm_begin, imm_end, abi::word::SIGNED);

  int condition_begin = imm_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  return new abi::instruction::Instruction(op_code, suffix, ra, imm, condition);
}

abi::instruction::Instruction *InstructionEncoder::decode_zrr(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rrr_op_codes().count(op_code));
  assert(suffix == abi::instruction::ZRR);

  int ra_begin = suffix_end();
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int rb_begin = ra_end;
  int rb_end = rb_begin + register_width();
  abi::reg::SrcReg *rb = decode_src_reg(instruction_word, rb_begin, rb_end);

  return new abi::instruction::Instruction(op_code, suffix, ra, rb);
}

abi::instruction::Instruction *InstructionEncoder::decode_zrrc(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rrrc_op_codes().count(op_code));
  assert(suffix == abi::instruction::ZRRC);

  int ra_begin = suffix_end();
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int rb_begin = ra_end;
  int rb_end = rb_begin + register_width();
  abi::reg::SrcReg *rb = decode_src_reg(instruction_word, rb_begin, rb_end);

  int condition_begin = rb_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  return new abi::instruction::Instruction(op_code, suffix, ra, rb, condition);
}

abi::instruction::Instruction *InstructionEncoder::decode_zrrci(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rrrci_op_codes().count(op_code));
  assert(suffix == abi::instruction::ZRRCI);

  int ra_begin = suffix_end();
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int rb_begin = ra_end;
  int rb_end = rb_begin + register_width();
  abi::reg::SrcReg *rb = decode_src_reg(instruction_word, rb_begin, rb_end);

  int condition_begin = rb_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  int pc_begin = condition_end;
  int pc_end = pc_begin + pc_width();
  int64_t pc = decode_pc(instruction_word, pc_begin, pc_end);

  return new abi::instruction::Instruction(op_code, suffix, ra, rb, condition,
                                           pc);
}

abi::instruction::Instruction *InstructionEncoder::decode_s_rri(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rri_op_codes().count(op_code));
  assert(suffix == abi::instruction::S_RRI or
         suffix == abi::instruction::U_RRI);

  int dc_begin = suffix_end();
  int dc_end = dc_begin + register_width();
  abi::reg::PairReg *dc = decode_pair_reg(instruction_word, dc_begin, dc_end);

  int ra_begin = dc_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int imm_begin = ra_end;
  int imm_end;
  int64_t imm;
  if (abi::instruction::Instruction::add_rri_op_codes().count(op_code)) {
    imm_end = imm_begin + 32;
    imm = decode_imm(instruction_word, imm_begin, imm_end, abi::word::UNSIGNED);
  } else if (abi::instruction::Instruction::asr_rri_op_codes().count(op_code)) {
    imm_end = imm_begin + 5;
    imm = decode_imm(instruction_word, imm_begin, imm_end, abi::word::UNSIGNED);
  } else if (abi::instruction::Instruction::call_rri_op_codes().count(
                 op_code)) {
    imm_end = imm_begin + 24;
    imm = decode_imm(instruction_word, imm_begin, imm_end, abi::word::SIGNED);
  } else {
    throw std::invalid_argument("");
  }

  return new abi::instruction::Instruction(op_code, suffix, dc, ra, imm);
}

abi::instruction::Instruction *InstructionEncoder::decode_s_rric(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rric_op_codes().count(op_code));
  assert(suffix == abi::instruction::S_RRIC or
         suffix == abi::instruction::U_RRIC);

  int dc_begin = suffix_end();
  int dc_end = dc_begin + register_width();
  abi::reg::PairReg *dc = decode_pair_reg(instruction_word, dc_begin, dc_end);

  int ra_begin = dc_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int imm_begin = ra_end;
  int imm_end;
  int64_t imm;
  if (abi::instruction::Instruction::add_rric_op_codes().count(op_code) or
      abi::instruction::Instruction::sub_rric_op_codes().count(op_code)) {
    imm_end = imm_begin + 24;
    imm = decode_imm(instruction_word, imm_begin, imm_end, abi::word::UNSIGNED);
  } else if (abi::instruction::Instruction::asr_rri_op_codes().count(op_code)) {
    imm_end = imm_begin + 5;
    imm = decode_imm(instruction_word, imm_begin, imm_end, abi::word::UNSIGNED);
  } else {
    throw std::invalid_argument("");
  }

  int condition_begin = imm_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  return new abi::instruction::Instruction(op_code, suffix, dc, ra, imm,
                                           condition);
}

abi::instruction::Instruction *InstructionEncoder::decode_s_rrici(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rrici_op_codes().count(op_code));
  assert(suffix == abi::instruction::S_RRICI or
         suffix == abi::instruction::U_RRICI);

  int dc_begin = suffix_end();
  int dc_end = dc_begin + register_width();
  abi::reg::PairReg *dc = decode_pair_reg(instruction_word, dc_begin, dc_end);

  int ra_begin = dc_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int imm_begin = ra_end;
  int imm_end;
  int64_t imm;
  if (abi::instruction::Instruction::add_rrici_op_codes().count(op_code) or
      abi::instruction::Instruction::and_rrici_op_codes().count(op_code) or
      abi::instruction::Instruction::sub_rrici_op_codes().count(op_code)) {
    imm_end = imm_begin + 8;
    imm = decode_imm(instruction_word, imm_begin, imm_end, abi::word::UNSIGNED);
  } else if (abi::instruction::Instruction::asr_rrici_op_codes().count(
                 op_code)) {
    imm_end = imm_begin + 5;
    imm = decode_imm(instruction_word, imm_begin, imm_end, abi::word::UNSIGNED);
  } else {
    throw std::invalid_argument("");
  }

  int condition_begin = imm_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  int pc_begin = condition_end;
  int pc_end = pc_begin + pc_width();
  int64_t pc = decode_pc(instruction_word, pc_begin, pc_end);

  return new abi::instruction::Instruction(op_code, suffix, dc, ra, imm,
                                           condition, pc);
}

abi::instruction::Instruction *InstructionEncoder::decode_s_rrif(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rrif_op_codes().count(op_code));
  assert(suffix == abi::instruction::S_RRIF or
         suffix == abi::instruction::U_RRIF);

  int dc_begin = suffix_end();
  int dc_end = dc_begin + register_width();
  abi::reg::PairReg *dc = decode_pair_reg(instruction_word, dc_begin, dc_end);

  int ra_begin = dc_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int imm_begin = ra_end;
  int imm_end = imm_begin + 24;
  int64_t imm =
      decode_imm(instruction_word, imm_begin, imm_end, abi::word::SIGNED);

  int condition_begin = imm_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  return new abi::instruction::Instruction(op_code, suffix, dc, ra, imm,
                                           condition);
}

abi::instruction::Instruction *InstructionEncoder::decode_s_rrr(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rrr_op_codes().count(op_code));
  assert(suffix == abi::instruction::S_RRR or
         suffix == abi::instruction::U_RRR);

  int dc_begin = suffix_end();
  int dc_end = dc_begin + register_width();
  abi::reg::PairReg *dc = decode_pair_reg(instruction_word, dc_begin, dc_end);

  int ra_begin = dc_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int rb_begin = ra_end;
  int rb_end = rb_begin + register_width();
  abi::reg::SrcReg *rb = decode_src_reg(instruction_word, rb_begin, rb_end);

  return new abi::instruction::Instruction(op_code, suffix, dc, ra, rb);
}

abi::instruction::Instruction *InstructionEncoder::decode_s_rrrc(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rrrc_op_codes().count(op_code));
  assert(suffix == abi::instruction::S_RRRC or
         suffix == abi::instruction::U_RRRC);

  int dc_begin = suffix_end();
  int dc_end = dc_begin + register_width();
  abi::reg::PairReg *dc = decode_pair_reg(instruction_word, dc_begin, dc_end);

  int ra_begin = dc_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int rb_begin = ra_end;
  int rb_end = rb_begin + register_width();
  abi::reg::SrcReg *rb = decode_src_reg(instruction_word, rb_begin, rb_end);

  int condition_begin = rb_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  return new abi::instruction::Instruction(op_code, suffix, dc, ra, rb,
                                           condition);
}

abi::instruction::Instruction *InstructionEncoder::decode_s_rrrci(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rrrci_op_codes().count(op_code));
  assert(suffix == abi::instruction::S_RRRCI or
         suffix == abi::instruction::U_RRRCI);

  int dc_begin = suffix_end();
  int dc_end = dc_begin + register_width();
  abi::reg::PairReg *dc = decode_pair_reg(instruction_word, dc_begin, dc_end);

  int ra_begin = dc_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int rb_begin = ra_end;
  int rb_end = rb_begin + register_width();
  abi::reg::SrcReg *rb = decode_src_reg(instruction_word, rb_begin, rb_end);

  int condition_begin = rb_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  int pc_begin = condition_end;
  int pc_end = pc_begin + pc_width();
  int64_t pc = decode_pc(instruction_word, pc_begin, pc_end);

  return new abi::instruction::Instruction(op_code, suffix, dc, ra, rb,
                                           condition, pc);
}

abi::instruction::Instruction *InstructionEncoder::decode_rr(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rr_op_codes().count(op_code));
  assert(suffix == abi::instruction::RR);

  int rc_begin = suffix_end();
  int rc_end = rc_begin + register_width();
  abi::reg::GPReg *rc = decode_gp_reg(instruction_word, rc_begin, rc_end);

  int ra_begin = rc_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  return new abi::instruction::Instruction(op_code, suffix, rc, ra);
}

abi::instruction::Instruction *InstructionEncoder::decode_rrc(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rrc_op_codes().count(op_code));
  assert(suffix == abi::instruction::RRC);

  int rc_begin = suffix_end();
  int rc_end = rc_begin + register_width();
  abi::reg::GPReg *rc = decode_gp_reg(instruction_word, rc_begin, rc_end);

  int ra_begin = rc_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int condition_begin = ra_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  return new abi::instruction::Instruction(op_code, suffix, rc, ra, condition);
}

abi::instruction::Instruction *InstructionEncoder::decode_rrci(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rrci_op_codes().count(op_code));
  assert(suffix == abi::instruction::RRCI);

  int rc_begin = suffix_end();
  int rc_end = rc_begin + register_width();
  abi::reg::GPReg *rc = decode_gp_reg(instruction_word, rc_begin, rc_end);

  int ra_begin = rc_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int condition_begin = ra_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  int pc_begin = condition_end;
  int pc_end = pc_begin + pc_width();
  int64_t pc = decode_pc(instruction_word, pc_begin, pc_end);

  return new abi::instruction::Instruction(op_code, suffix, rc, ra, condition,
                                           pc);
}

abi::instruction::Instruction *InstructionEncoder::decode_zr(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rr_op_codes().count(op_code));
  assert(suffix == abi::instruction::ZR);

  int ra_begin = suffix_end();
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  return new abi::instruction::Instruction(op_code, suffix, ra);
}

abi::instruction::Instruction *InstructionEncoder::decode_zrc(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rrc_op_codes().count(op_code));
  assert(suffix == abi::instruction::ZRC);

  int ra_begin = suffix_end();
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int condition_begin = ra_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  return new abi::instruction::Instruction(op_code, suffix, ra, condition);
}

abi::instruction::Instruction *InstructionEncoder::decode_zrci(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rrci_op_codes().count(op_code));
  assert(suffix == abi::instruction::ZRCI);

  int ra_begin = suffix_end();
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int condition_begin = ra_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  int pc_begin = condition_end;
  int pc_end = pc_begin + pc_width();
  int64_t pc = decode_pc(instruction_word, pc_begin, pc_end);

  return new abi::instruction::Instruction(op_code, suffix, ra, condition, pc);
}

abi::instruction::Instruction *InstructionEncoder::decode_s_rr(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rr_op_codes().count(op_code));
  assert(suffix == abi::instruction::S_RR or suffix == abi::instruction::U_RR);

  int dc_begin = suffix_end();
  int dc_end = dc_begin + register_width();
  abi::reg::PairReg *dc = decode_pair_reg(instruction_word, dc_begin, dc_end);

  int ra_begin = dc_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  return new abi::instruction::Instruction(op_code, suffix, dc, ra);
}

abi::instruction::Instruction *InstructionEncoder::decode_s_rrc(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rrc_op_codes().count(op_code));
  assert(suffix == abi::instruction::S_RRC or
         suffix == abi::instruction::U_RRC);

  int dc_begin = suffix_end();
  int dc_end = dc_begin + register_width();
  abi::reg::PairReg *dc = decode_pair_reg(instruction_word, dc_begin, dc_end);

  int ra_begin = dc_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int condition_begin = ra_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  return new abi::instruction::Instruction(op_code, suffix, dc, ra, condition);
}

abi::instruction::Instruction *InstructionEncoder::decode_s_rrci(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rrci_op_codes().count(op_code));
  assert(suffix == abi::instruction::S_RRCI or
         suffix == abi::instruction::U_RRCI);

  int dc_begin = suffix_end();
  int dc_end = dc_begin + register_width();
  abi::reg::PairReg *dc = decode_pair_reg(instruction_word, dc_begin, dc_end);

  int ra_begin = dc_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int condition_begin = ra_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  int pc_begin = condition_end;
  int pc_end = pc_begin + pc_width();
  int64_t pc = decode_pc(instruction_word, pc_begin, pc_end);

  return new abi::instruction::Instruction(op_code, suffix, dc, ra, condition,
                                           pc);
}

abi::instruction::Instruction *InstructionEncoder::decode_drdici(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::drdici_op_codes().count(op_code));
  assert(suffix == abi::instruction::DRDICI);

  int dc_begin = suffix_end();
  int dc_end = dc_begin + register_width();
  abi::reg::PairReg *dc = decode_pair_reg(instruction_word, dc_begin, dc_end);

  int ra_begin = dc_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int db_begin = ra_end;
  int db_end = db_begin + register_width();
  abi::reg::PairReg *db = decode_pair_reg(instruction_word, db_begin, db_end);

  int imm_begin = db_end;
  int imm_end = imm_begin + 5;
  int64_t imm =
      decode_imm(instruction_word, imm_begin, imm_end, abi::word::UNSIGNED);

  int condition_begin = imm_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  int pc_begin = condition_end;
  int pc_end = pc_begin + pc_width();
  int64_t pc = decode_pc(instruction_word, pc_begin, pc_end);

  return new abi::instruction::Instruction(op_code, suffix, dc, ra, db, imm,
                                           condition, pc);
}

abi::instruction::Instruction *InstructionEncoder::decode_rrri(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rrri_op_codes().count(op_code));
  assert(suffix == abi::instruction::RRRI);

  int rc_begin = suffix_end();
  int rc_end = rc_begin + register_width();
  abi::reg::GPReg *rc = decode_gp_reg(instruction_word, rc_begin, rc_end);

  int ra_begin = rc_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int rb_begin = ra_end;
  int rb_end = rb_begin + register_width();
  abi::reg::SrcReg *rb = decode_src_reg(instruction_word, rb_begin, rb_end);

  int imm_begin = rb_end;
  int imm_end = imm_begin + 5;
  int64_t imm =
      decode_imm(instruction_word, imm_begin, imm_end, abi::word::UNSIGNED);

  return new abi::instruction::Instruction(op_code, suffix, rc, ra, rb, imm);
}

abi::instruction::Instruction *InstructionEncoder::decode_rrrici(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rrrici_op_codes().count(op_code));
  assert(suffix == abi::instruction::RRRICI);

  int rc_begin = suffix_end();
  int rc_end = rc_begin + register_width();
  abi::reg::GPReg *rc = decode_gp_reg(instruction_word, rc_begin, rc_end);

  int ra_begin = rc_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int rb_begin = ra_end;
  int rb_end = rb_begin + register_width();
  abi::reg::SrcReg *rb = decode_src_reg(instruction_word, rb_begin, rb_end);

  int imm_begin = rb_end;
  int imm_end = imm_begin + 5;
  int64_t imm =
      decode_imm(instruction_word, imm_begin, imm_end, abi::word::UNSIGNED);

  int condition_begin = imm_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  int pc_begin = condition_end;
  int pc_end = pc_begin + pc_width();
  int64_t pc = decode_pc(instruction_word, pc_begin, pc_end);

  return new abi::instruction::Instruction(op_code, suffix, rc, ra, rb, imm,
                                           condition, pc);
}

abi::instruction::Instruction *InstructionEncoder::decode_zrri(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rrri_op_codes().count(op_code));
  assert(suffix == abi::instruction::ZRRI);

  int ra_begin = suffix_end();
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int rb_begin = ra_end;
  int rb_end = rb_begin + register_width();
  abi::reg::SrcReg *rb = decode_src_reg(instruction_word, rb_begin, rb_end);

  int imm_begin = rb_end;
  int imm_end = imm_begin + 5;
  int64_t imm =
      decode_imm(instruction_word, imm_begin, imm_end, abi::word::UNSIGNED);

  return new abi::instruction::Instruction(op_code, suffix, ra, rb, imm);
}

abi::instruction::Instruction *InstructionEncoder::decode_zrrici(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rrrici_op_codes().count(op_code));
  assert(suffix == abi::instruction::ZRRICI);

  int ra_begin = suffix_end();
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int rb_begin = ra_end;
  int rb_end = rb_begin + register_width();
  abi::reg::SrcReg *rb = decode_src_reg(instruction_word, rb_begin, rb_end);

  int imm_begin = rb_end;
  int imm_end = imm_begin + 5;
  int64_t imm =
      decode_imm(instruction_word, imm_begin, imm_end, abi::word::UNSIGNED);

  int condition_begin = imm_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  int pc_begin = condition_end;
  int pc_end = pc_begin + pc_width();
  int64_t pc = decode_pc(instruction_word, pc_begin, pc_end);

  return new abi::instruction::Instruction(op_code, suffix, ra, rb, imm,
                                           condition, pc);
}

abi::instruction::Instruction *InstructionEncoder::decode_s_rrri(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rrri_op_codes().count(op_code));
  assert(suffix == abi::instruction::S_RRRI or
         suffix == abi::instruction::U_RRRI);

  int dc_begin = suffix_end();
  int dc_end = dc_begin + register_width();
  abi::reg::PairReg *dc = decode_pair_reg(instruction_word, dc_begin, dc_end);

  int ra_begin = dc_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int rb_begin = ra_end;
  int rb_end = rb_begin + register_width();
  abi::reg::SrcReg *rb = decode_src_reg(instruction_word, rb_begin, rb_end);

  int imm_begin = rb_end;
  int imm_end = imm_begin + 5;
  int64_t imm =
      decode_imm(instruction_word, imm_begin, imm_end, abi::word::UNSIGNED);

  return new abi::instruction::Instruction(op_code, suffix, dc, ra, rb, imm);
}

abi::instruction::Instruction *InstructionEncoder::decode_s_rrrici(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rrrici_op_codes().count(op_code));
  assert(suffix == abi::instruction::S_RRICI or
         suffix == abi::instruction::U_RRRICI);

  int dc_begin = suffix_end();
  int dc_end = dc_begin + register_width();
  abi::reg::PairReg *dc = decode_pair_reg(instruction_word, dc_begin, dc_end);

  int ra_begin = dc_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int rb_begin = ra_end;
  int rb_end = rb_begin + register_width();
  abi::reg::SrcReg *rb = decode_src_reg(instruction_word, rb_begin, rb_end);

  int imm_begin = rb_end;
  int imm_end = imm_begin + 5;
  int64_t imm =
      decode_imm(instruction_word, imm_begin, imm_end, abi::word::UNSIGNED);

  int condition_begin = imm_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  int pc_begin = condition_end;
  int pc_end = pc_begin + pc_width();
  int64_t pc = decode_pc(instruction_word, pc_begin, pc_end);

  return new abi::instruction::Instruction(op_code, suffix, dc, ra, rb, imm,
                                           condition, pc);
}

abi::instruction::Instruction *InstructionEncoder::decode_rir(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rir_op_codes().count(op_code));
  assert(suffix == abi::instruction::RIR);

  int rc_begin = suffix_end();
  int rc_end = rc_begin + register_width();
  abi::reg::GPReg *rc = decode_gp_reg(instruction_word, rc_begin, rc_end);

  int imm_begin = rc_end;
  int imm_end = imm_begin + 32;
  int64_t imm =
      decode_imm(instruction_word, imm_begin, imm_end, abi::word::UNSIGNED);

  int ra_begin = imm_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  return new abi::instruction::Instruction(op_code, suffix, rc, imm, ra);
}

abi::instruction::Instruction *InstructionEncoder::decode_rirc(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rirc_op_codes().count(op_code));
  assert(suffix == abi::instruction::RIRC);

  int rc_begin = suffix_end();
  int rc_end = rc_begin + register_width();
  abi::reg::GPReg *rc = decode_gp_reg(instruction_word, rc_begin, rc_end);

  int imm_begin = rc_end;
  int imm_end = imm_begin + 24;
  int64_t imm =
      decode_imm(instruction_word, imm_begin, imm_end, abi::word::SIGNED);

  int ra_begin = imm_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int condition_begin = ra_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  return new abi::instruction::Instruction(op_code, suffix, rc, imm, ra,
                                           condition);
}

abi::instruction::Instruction *InstructionEncoder::decode_rirci(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rirci_op_codes().count(op_code));
  assert(suffix == abi::instruction::RIRCI);

  int rc_begin = suffix_end();
  int rc_end = rc_begin + register_width();
  abi::reg::GPReg *rc = decode_gp_reg(instruction_word, rc_begin, rc_end);

  int imm_begin = rc_end;
  int imm_end = imm_begin + 8;
  int64_t imm =
      decode_imm(instruction_word, imm_begin, imm_end, abi::word::SIGNED);

  int ra_begin = imm_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int condition_begin = ra_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  int pc_begin = condition_end;
  int pc_end = pc_begin + pc_width();
  int64_t pc = decode_pc(instruction_word, pc_begin, pc_end);

  return new abi::instruction::Instruction(op_code, suffix, rc, imm, ra,
                                           condition, pc);
}

abi::instruction::Instruction *InstructionEncoder::decode_zir(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rir_op_codes().count(op_code));
  assert(suffix == abi::instruction::ZIR);

  int imm_begin = suffix_end();
  int imm_end = imm_begin + 32;
  int64_t imm =
      decode_imm(instruction_word, imm_begin, imm_end, abi::word::UNSIGNED);

  int ra_begin = imm_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  return new abi::instruction::Instruction(op_code, suffix, imm, ra);
}

abi::instruction::Instruction *InstructionEncoder::decode_zirc(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rirc_op_codes().count(op_code));
  assert(suffix == abi::instruction::ZIRC);

  int imm_begin = suffix_end();
  int imm_end = imm_begin + 24;
  int64_t imm =
      decode_imm(instruction_word, imm_begin, imm_end, abi::word::SIGNED);

  int ra_begin = imm_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int condition_begin = ra_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  return new abi::instruction::Instruction(op_code, suffix, imm, ra, condition);
}

abi::instruction::Instruction *InstructionEncoder::decode_zirci(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rirci_op_codes().count(op_code));
  assert(suffix == abi::instruction::ZIRCI);

  int imm_begin = suffix_end();
  int imm_end = imm_begin + 8;
  int64_t imm =
      decode_imm(instruction_word, imm_begin, imm_end, abi::word::SIGNED);

  int ra_begin = imm_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int condition_begin = ra_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  int pc_begin = condition_end;
  int pc_end = pc_begin + pc_width();
  int64_t pc = decode_pc(instruction_word, pc_begin, pc_end);

  return new abi::instruction::Instruction(op_code, suffix, imm, ra, condition,
                                           pc);
}

abi::instruction::Instruction *InstructionEncoder::decode_s_rirc(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rirc_op_codes().count(op_code));
  assert(suffix == abi::instruction::S_RIRC or
         suffix == abi::instruction::U_RIRC);

  int dc_begin = suffix_end();
  int dc_end = dc_begin + register_width();
  abi::reg::PairReg *dc = decode_pair_reg(instruction_word, dc_begin, dc_end);

  int imm_begin = dc_end;
  int imm_end = imm_begin + 24;
  int64_t imm =
      decode_imm(instruction_word, imm_begin, imm_end, abi::word::SIGNED);

  int ra_begin = imm_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int condition_begin = ra_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  return new abi::instruction::Instruction(op_code, suffix, dc, imm, ra,
                                           condition);
}

abi::instruction::Instruction *InstructionEncoder::decode_s_rirci(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rirci_op_codes().count(op_code));
  assert(suffix == abi::instruction::S_RIRCI or
         suffix == abi::instruction::U_RIRCI);

  int dc_begin = suffix_end();
  int dc_end = dc_begin + register_width();
  abi::reg::PairReg *dc = decode_pair_reg(instruction_word, dc_begin, dc_end);

  int imm_begin = dc_end;
  int imm_end = imm_begin + 8;
  int64_t imm =
      decode_imm(instruction_word, imm_begin, imm_end, abi::word::SIGNED);

  int ra_begin = imm_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int condition_begin = ra_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  int pc_begin = condition_end;
  int pc_end = pc_begin + pc_width();
  int64_t pc = decode_pc(instruction_word, pc_begin, pc_end);

  return new abi::instruction::Instruction(op_code, suffix, dc, imm, ra,
                                           condition, pc);
}

abi::instruction::Instruction *InstructionEncoder::decode_r(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::r_op_codes().count(op_code));
  assert(suffix == abi::instruction::R);

  int rc_begin = suffix_end();
  int rc_end = rc_begin + register_width();
  abi::reg::GPReg *rc = decode_gp_reg(instruction_word, rc_begin, rc_end);

  return new abi::instruction::Instruction(op_code, suffix, rc);
}

abi::instruction::Instruction *InstructionEncoder::decode_rci(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rci_op_codes().count(op_code));
  assert(suffix == abi::instruction::RCI);

  int rc_begin = suffix_end();
  int rc_end = rc_begin + register_width();
  abi::reg::GPReg *rc = decode_gp_reg(instruction_word, rc_begin, rc_end);

  int condition_begin = rc_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  int pc_begin = condition_end;
  int pc_end = pc_begin + pc_width();
  int64_t pc = decode_pc(instruction_word, pc_begin, pc_end);

  return new abi::instruction::Instruction(op_code, suffix, rc, condition, pc);
}

abi::instruction::Instruction *InstructionEncoder::decode_z(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::r_op_codes().count(op_code) or
         op_code == abi::instruction::NOP);
  assert(suffix == abi::instruction::Z);

  return new abi::instruction::Instruction(op_code, suffix);
}

abi::instruction::Instruction *InstructionEncoder::decode_zci(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rci_op_codes().count(op_code));
  assert(suffix == abi::instruction::ZCI);

  int condition_begin = suffix_end();
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  int pc_begin = condition_end;
  int pc_end = pc_begin + pc_width();
  int64_t pc = decode_pc(instruction_word, pc_begin, pc_end);

  return new abi::instruction::Instruction(op_code, suffix, condition, pc);
}

abi::instruction::Instruction *InstructionEncoder::decode_s_r(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::r_op_codes().count(op_code));
  assert(suffix == abi::instruction::S_R or suffix == abi::instruction::U_R);

  int dc_begin = suffix_end();
  int dc_end = dc_begin + register_width();
  abi::reg::PairReg *dc = decode_pair_reg(instruction_word, dc_begin, dc_end);

  return new abi::instruction::Instruction(op_code, suffix, dc);
}

abi::instruction::Instruction *InstructionEncoder::decode_s_rci(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::rci_op_codes().count(op_code));
  assert(suffix == abi::instruction::S_RCI or
         suffix == abi::instruction::U_RCI);

  int dc_begin = suffix_end();
  int dc_end = dc_begin + register_width();
  abi::reg::PairReg *dc = decode_pair_reg(instruction_word, dc_begin, dc_end);

  int condition_begin = dc_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  int pc_begin = condition_end;
  int pc_end = pc_begin + pc_width();
  int64_t pc = decode_pc(instruction_word, pc_begin, pc_end);

  return new abi::instruction::Instruction(op_code, suffix, dc, condition, pc);
}

abi::instruction::Instruction *InstructionEncoder::decode_ci(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::ci_op_codes().count(op_code));
  assert(suffix == abi::instruction::CI);

  int condition_begin = suffix_end();
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  int pc_begin = condition_end;
  int pc_end = pc_begin + pc_width();
  int64_t pc = decode_pc(instruction_word, pc_begin, pc_end);

  return new abi::instruction::Instruction(op_code, suffix, condition, pc);
}

abi::instruction::Instruction *InstructionEncoder::decode_i(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::i_op_codes().count(op_code));
  assert(suffix == abi::instruction::I);

  int imm_begin = suffix_end();
  int imm_end = imm_begin + 24;
  int64_t imm =
      decode_imm(instruction_word, imm_begin, imm_end, abi::word::SIGNED);

  return new abi::instruction::Instruction(op_code, suffix, imm);
}

abi::instruction::Instruction *InstructionEncoder::decode_ddci(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::ddci_op_codes().count(op_code));
  assert(suffix == abi::instruction::DDCI);

  int dc_begin = suffix_end();
  int dc_end = dc_begin + register_width();
  abi::reg::PairReg *dc = decode_pair_reg(instruction_word, dc_begin, dc_end);

  int db_begin = dc_end;
  int db_end = db_begin + register_width();
  abi::reg::PairReg *db = decode_pair_reg(instruction_word, db_begin, db_end);

  int condition_begin = db_end;
  int condition_end = condition_begin + condition_width();
  abi::isa::Condition condition =
      decode_condition(instruction_word, condition_begin, condition_end);

  int pc_begin = condition_end;
  int pc_end = pc_begin + pc_width();
  int64_t pc = decode_pc(instruction_word, pc_begin, pc_end);

  return new abi::instruction::Instruction(op_code, suffix, dc, db, condition,
                                           pc);
}

abi::instruction::Instruction *InstructionEncoder::decode_erri(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::erri_op_codes().count(op_code));
  assert(suffix == abi::instruction::ERRI);

  int endian_begin = suffix_end();
  int endian_end = endian_begin + endian_width();
  abi::isa::Endian endian =
      decode_endian(instruction_word, endian_begin, endian_end);

  int rc_begin = endian_end;
  int rc_end = rc_begin + register_width();
  abi::reg::GPReg *rc = decode_gp_reg(instruction_word, rc_begin, rc_end);

  int ra_begin = rc_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int off_begin = ra_end;
  int off_end = off_begin + 24;
  int64_t off =
      decode_off(instruction_word, off_begin, off_end, abi::word::SIGNED);

  return new abi::instruction::Instruction(op_code, suffix, endian, rc, ra,
                                           off);
}

abi::instruction::Instruction *InstructionEncoder::decode_s_erri(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::erri_op_codes().count(op_code));
  assert(suffix == abi::instruction::S_ERRI or
         suffix == abi::instruction::U_ERRI);

  int endian_begin = suffix_end();
  int endian_end = endian_begin + endian_width();
  abi::isa::Endian endian =
      decode_endian(instruction_word, endian_begin, endian_end);

  int dc_begin = endian_end;
  int dc_end = dc_begin + register_width();
  abi::reg::PairReg *dc = decode_pair_reg(instruction_word, dc_begin, dc_end);

  int ra_begin = dc_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int off_begin = ra_end;
  int off_end = off_begin + 24;
  int64_t off =
      decode_off(instruction_word, off_begin, off_end, abi::word::SIGNED);

  return new abi::instruction::Instruction(op_code, suffix, endian, dc, ra,
                                           off);
}

abi::instruction::Instruction *InstructionEncoder::decode_edri(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::edri_op_codes().count(op_code));
  assert(suffix == abi::instruction::EDRI);

  int endian_begin = suffix_end();
  int endian_end = endian_begin + endian_width();
  abi::isa::Endian endian =
      decode_endian(instruction_word, endian_begin, endian_end);

  int dc_begin = endian_end;
  int dc_end = dc_begin + register_width();
  abi::reg::PairReg *dc = decode_pair_reg(instruction_word, dc_begin, dc_end);

  int ra_begin = dc_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int off_begin = ra_end;
  int off_end = off_begin + 24;
  int64_t off =
      decode_off(instruction_word, off_begin, off_end, abi::word::SIGNED);

  return new abi::instruction::Instruction(op_code, suffix, endian, dc, ra,
                                           off);
}

abi::instruction::Instruction *InstructionEncoder::decode_erii(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::erii_op_codes().count(op_code));
  assert(suffix == abi::instruction::ERII);

  int endian_begin = suffix_end();
  int endian_end = endian_begin + endian_width();
  abi::isa::Endian endian =
      decode_endian(instruction_word, endian_begin, endian_end);

  int ra_begin = endian_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int off_begin = ra_end;
  int off_end = off_begin + 24;
  int64_t off =
      decode_off(instruction_word, off_begin, off_end, abi::word::SIGNED);

  int imm_begin = off_end;
  int imm_end = imm_begin + 16;
  int64_t imm =
      decode_imm(instruction_word, imm_begin, imm_end, abi::word::SIGNED);

  return new abi::instruction::Instruction(op_code, suffix, endian, ra, off,
                                           imm);
}

abi::instruction::Instruction *InstructionEncoder::decode_erir(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::erir_op_codes().count(op_code));
  assert(suffix == abi::instruction::ERIR);

  int endian_begin = suffix_end();
  int endian_end = endian_begin + endian_width();
  abi::isa::Endian endian =
      decode_endian(instruction_word, endian_begin, endian_end);

  int ra_begin = endian_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int off_begin = ra_end;
  int off_end = off_begin + 24;
  int64_t off =
      decode_off(instruction_word, off_begin, off_end, abi::word::SIGNED);

  int rb_begin = off_end;
  int rb_end = rb_begin + register_width();
  abi::reg::SrcReg *rb = decode_src_reg(instruction_word, rb_begin, rb_end);

  return new abi::instruction::Instruction(op_code, suffix, endian, ra, off,
                                           rb);
}

abi::instruction::Instruction *InstructionEncoder::decode_erid(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::erid_op_codes().count(op_code));
  assert(suffix == abi::instruction::ERID);

  int endian_begin = suffix_end();
  int endian_end = endian_begin + endian_width();
  abi::isa::Endian endian =
      decode_endian(instruction_word, endian_begin, endian_end);

  int ra_begin = endian_end;
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int off_begin = ra_end;
  int off_end = off_begin + 24;
  int64_t off =
      decode_off(instruction_word, off_begin, off_end, abi::word::SIGNED);

  int db_begin = off_end;
  int db_end = db_begin + register_width();
  abi::reg::PairReg *db = decode_pair_reg(instruction_word, db_begin, db_end);

  return new abi::instruction::Instruction(op_code, suffix, endian, ra, off,
                                           db);
}

abi::instruction::Instruction *InstructionEncoder::decode_dma_rri(
    abi::instruction::OpCode op_code, abi::instruction::Suffix suffix,
    abi::word::InstructionWord *instruction_word) {
  assert(abi::instruction::Instruction::dma_rri_op_codes().count(op_code));
  assert(suffix == abi::instruction::DMA_RRI);

  int ra_begin = suffix_end();
  int ra_end = ra_begin + register_width();
  abi::reg::SrcReg *ra = decode_src_reg(instruction_word, ra_begin, ra_end);

  int rb_begin = ra_end;
  int rb_end = rb_begin + register_width();
  abi::reg::SrcReg *rb = decode_src_reg(instruction_word, rb_begin, rb_end);

  int imm_begin = rb_end;
  int imm_end = imm_begin + 8;
  int64_t imm =
      decode_imm(instruction_word, imm_begin, imm_end, abi::word::UNSIGNED);

  return new abi::instruction::Instruction(op_code, suffix, ra, rb, imm);
}

}  // namespace upmem_sim::encoder
