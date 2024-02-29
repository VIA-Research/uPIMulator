#include "instruction.h"

#include <sstream>

#include "abi/cc/acquire_cc.h"
#include "abi/cc/add_nz_cc.h"
#include "abi/cc/boot_cc.h"
#include "abi/cc/count_nz_cc.h"
#include "abi/cc/div_cc.h"
#include "abi/cc/div_nz_cc.h"
#include "abi/cc/ext_sub_set_cc.h"
#include "abi/cc/false_cc.h"
#include "abi/cc/imm_shift_nz_cc.h"
#include "abi/cc/log_nz_cc.h"
#include "abi/cc/log_set_cc.h"
#include "abi/cc/mul_nz_cc.h"
#include "abi/cc/release_cc.h"
#include "abi/cc/shift_nz_cc.h"
#include "abi/cc/sub_nz_cc.h"
#include "abi/cc/sub_set_cc.h"
#include "abi/cc/true_cc.h"
#include "abi/cc/true_false_cc.h"

namespace upmem_sim::abi::instruction {

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::SrcReg *ra,
                         int64_t imm, isa::Condition condition, int64_t pc)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(ra != nullptr);

  if (suffix_ == RICI) {
    init_rici(ra, imm, condition, pc);
  } else if (suffix_ == ZRICI) {
    init_zrici(ra, imm, condition, pc);
  } else {
    throw std::invalid_argument("");
  }
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::GPReg *rc,
                         reg::SrcReg *ra, int64_t imm)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(rc != nullptr);
  assert(ra != nullptr);

  init_rri(rc, ra, imm);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::GPReg *rc,
                         reg::SrcReg *ra, int64_t imm, isa::Condition condition)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(rc != nullptr);
  assert(ra != nullptr);

  if (suffix_ == RRIC) {
    init_rric(rc, ra, imm, condition);
  } else if (suffix_ == RRIF) {
    init_rrif(rc, ra, imm, condition);
  } else {
    throw std::invalid_argument("");
  }
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::GPReg *rc,
                         reg::SrcReg *ra, int64_t imm, isa::Condition condition,
                         int64_t pc)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(rc != nullptr);
  assert(ra != nullptr);

  init_rrici(rc, ra, imm, condition, pc);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::GPReg *rc,
                         reg::SrcReg *ra, reg::SrcReg *rb)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(rc != nullptr);
  assert(ra != nullptr);
  assert(rb != nullptr);

  init_rrr(rc, ra, rb);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::GPReg *rc,
                         reg::SrcReg *ra, reg::SrcReg *rb,
                         isa::Condition condition)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(rc != nullptr);
  assert(ra != nullptr);
  assert(rb != nullptr);

  init_rrrc(rc, ra, rb, condition);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::GPReg *rc,
                         reg::SrcReg *ra, reg::SrcReg *rb,
                         isa::Condition condition, int64_t pc)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(rc != nullptr);
  assert(ra != nullptr);
  assert(rb != nullptr);

  init_rrrci(rc, ra, rb, condition, pc);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::SrcReg *ra,
                         int64_t imm)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(ra != nullptr);

  init_zri(ra, imm);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::SrcReg *ra,
                         int64_t imm, isa::Condition condition)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(ra != nullptr);

  if (suffix_ == ZRIC) {
    init_zric(ra, imm, condition);
  } else if (suffix_ == ZRIF) {
    init_zrif(ra, imm, condition);
  } else {
    throw std::invalid_argument("");
  }
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::SrcReg *ra,
                         reg::SrcReg *rb)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(ra != nullptr);
  assert(rb != nullptr);

  init_zrr(ra, rb);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::SrcReg *ra,
                         reg::SrcReg *rb, isa::Condition condition)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(ra != nullptr);
  assert(rb != nullptr);

  init_zrrc(ra, rb, condition);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::SrcReg *ra,
                         reg::SrcReg *rb, isa::Condition condition, int64_t pc)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(ra != nullptr);
  assert(rb != nullptr);

  init_zrrci(ra, rb, condition, pc);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc,
                         reg::SrcReg *ra, int64_t imm)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(dc != nullptr);
  assert(ra != nullptr);

  init_s_rri(dc, ra, imm);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc,
                         reg::SrcReg *ra, int64_t imm, isa::Condition condition)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(dc != nullptr);
  assert(ra != nullptr);

  if (suffix == S_RRIC or suffix == U_RRIC) {
    init_s_rric(dc, ra, imm, condition);
  } else if (suffix == S_RRIF or suffix == U_RRIF) {
    init_s_rrif(dc, ra, imm, condition);
  } else {
    throw std::invalid_argument("");
  }
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc,
                         reg::SrcReg *ra, int64_t imm, isa::Condition condition,
                         int64_t pc)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(dc != nullptr);
  assert(ra != nullptr);

  init_s_rrici(dc, ra, imm, condition, pc);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc,
                         reg::SrcReg *ra, reg::SrcReg *rb)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(dc != nullptr);
  assert(ra != nullptr);
  assert(rb != nullptr);

  init_s_rrr(dc, ra, rb);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc,
                         reg::SrcReg *ra, reg::SrcReg *rb,
                         isa::Condition condition)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(dc != nullptr);
  assert(ra != nullptr);
  assert(rb != nullptr);

  init_s_rrrc(dc, ra, rb, condition);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc,
                         reg::SrcReg *ra, reg::SrcReg *rb,
                         isa::Condition condition, int64_t pc)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(dc != nullptr);
  assert(ra != nullptr);
  assert(rb != nullptr);

  init_s_rrrci(dc, ra, rb, condition, pc);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::GPReg *rc,
                         reg::SrcReg *ra)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(rc != nullptr);
  assert(ra != nullptr);

  init_rr(rc, ra);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::GPReg *rc,
                         reg::SrcReg *ra, isa::Condition condition)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(rc != nullptr);
  assert(ra != nullptr);

  init_rrc(rc, ra, condition);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::GPReg *rc,
                         reg::SrcReg *ra, isa::Condition condition, int64_t pc)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(rc != nullptr);
  assert(ra != nullptr);

  init_rrci(rc, ra, condition, pc);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::SrcReg *ra)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(ra != nullptr);

  init_zr(ra);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::SrcReg *ra,
                         isa::Condition condition)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(ra != nullptr);

  init_zrc(ra, condition);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::SrcReg *ra,
                         isa::Condition condition, int64_t pc)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(ra != nullptr);

  init_zrci(ra, condition, pc);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc,
                         reg::SrcReg *ra)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(dc != nullptr);
  assert(ra != nullptr);

  init_s_rr(dc, ra);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc,
                         reg::SrcReg *ra, isa::Condition condition)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(dc != nullptr);
  assert(ra != nullptr);

  init_s_rrc(dc, ra, condition);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc,
                         reg::SrcReg *ra, isa::Condition condition, int64_t pc)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(dc != nullptr);
  assert(ra != nullptr);

  init_s_rrci(dc, ra, condition, pc);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc,
                         reg::SrcReg *ra, reg::PairReg *db, int64_t imm,
                         isa::Condition condition, int64_t pc)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(dc != nullptr);
  assert(ra != nullptr);
  assert(db != nullptr);

  init_drdici(dc, ra, db, imm, condition, pc);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::GPReg *rc,
                         reg::SrcReg *ra, reg::SrcReg *rb, int64_t imm)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(rc != nullptr);
  assert(ra != nullptr);
  assert(rb != nullptr);

  init_rrri(rc, ra, rb, imm);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::GPReg *rc,
                         reg::SrcReg *ra, reg::SrcReg *rb, int64_t imm,
                         isa::Condition condition, int64_t pc)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(rc != nullptr);
  assert(ra != nullptr);
  assert(rb != nullptr);

  init_rrrici(rc, ra, rb, imm, condition, pc);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::SrcReg *ra,
                         reg::SrcReg *rb, int64_t imm)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(ra != nullptr);
  assert(rb != nullptr);

  if (suffix == ZRRI) {
    init_zrri(ra, rb, imm);
  } else if (suffix == DMA_RRI) {
    init_dma_rri(ra, rb, imm);
  } else {
    throw std::invalid_argument("");
  }
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::SrcReg *ra,
                         reg::SrcReg *rb, int64_t imm, isa::Condition condition,
                         int64_t pc)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(ra != nullptr);
  assert(rb != nullptr);

  init_zrrici(ra, rb, imm, condition, pc);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc,
                         reg::SrcReg *ra, reg::SrcReg *rb, int64_t imm)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(dc != nullptr);
  assert(ra != nullptr);
  assert(rb != nullptr);

  init_s_rrri(dc, ra, rb, imm);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc,
                         reg::SrcReg *ra, reg::SrcReg *rb, int64_t imm,
                         isa::Condition condition, int64_t pc)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(dc != nullptr);
  assert(ra != nullptr);
  assert(rb != nullptr);

  init_s_rrrici(dc, ra, rb, imm, condition, pc);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::GPReg *rc,
                         int64_t imm, reg::SrcReg *ra)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(rc != nullptr);
  assert(ra != nullptr);

  init_rir(rc, imm, ra);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::GPReg *rc,
                         int64_t imm, reg::SrcReg *ra, isa::Condition condition)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(rc != nullptr);
  assert(ra != nullptr);

  init_rirc(rc, imm, ra, condition);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::GPReg *rc,
                         int64_t imm, reg::SrcReg *ra, isa::Condition condition,
                         int64_t pc)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(rc != nullptr);
  assert(ra != nullptr);

  init_rirci(rc, imm, ra, condition, pc);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, int64_t imm,
                         reg::SrcReg *ra)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(ra != nullptr);

  init_zir(imm, ra);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, int64_t imm,
                         reg::SrcReg *ra, isa::Condition condition)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(ra != nullptr);

  init_zirc(imm, ra, condition);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, int64_t imm,
                         reg::SrcReg *ra, isa::Condition condition, int64_t pc)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(ra != nullptr);

  init_zirci(imm, ra, condition, pc);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc,
                         int64_t imm, reg::SrcReg *ra, isa::Condition condition)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(dc != nullptr);
  assert(ra != nullptr);

  init_s_rirc(dc, imm, ra, condition);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc,
                         int64_t imm, reg::SrcReg *ra, isa::Condition condition,
                         int64_t pc)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(dc != nullptr);
  assert(ra != nullptr);

  init_s_rirci(dc, imm, ra, condition, pc);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::GPReg *rc)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(rc != nullptr);

  init_r(rc);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::GPReg *rc,
                         isa::Condition condition, int64_t pc)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(rc != nullptr);

  init_rci(rc, condition, pc);
}

Instruction::Instruction(OpCode op_code, Suffix suffix)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  init_z();
}

Instruction::Instruction(OpCode op_code, Suffix suffix,
                         isa::Condition condition, int64_t pc)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  if (suffix == ZCI) {
    init_zci(condition, pc);
  } else if (suffix == CI) {
    init_ci(condition, pc);
  } else {
    throw std::invalid_argument("");
  }
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(dc != nullptr);

  init_s_r(dc);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc,
                         isa::Condition condition, int64_t pc)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(dc != nullptr);

  init_s_rci(dc, condition, pc);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, int64_t imm)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  init_i(imm);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc,
                         reg::PairReg *db, isa::Condition condition, int64_t pc)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(dc != nullptr);
  assert(db != nullptr);

  init_ddci(dc, db, condition, pc);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, isa::Endian endian,
                         reg::GPReg *rc, reg::SrcReg *ra, int64_t off)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(rc != nullptr);
  assert(ra != nullptr);

  init_erri(endian, rc, ra, off);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, isa::Endian endian,
                         reg::PairReg *dc, reg::SrcReg *ra, int64_t off)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(dc != nullptr);
  assert(ra != nullptr);

  if (suffix == S_ERRI or suffix == U_ERRI) {
    init_s_erri(endian, dc, ra, off);
  } else if (suffix == EDRI) {
    init_edri(endian, dc, ra, off);
  } else {
    throw std::invalid_argument("");
  }
}

Instruction::Instruction(OpCode op_code, Suffix suffix, isa::Endian endian,
                         reg::SrcReg *ra, int64_t off, int64_t imm)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(ra != nullptr);

  init_erii(endian, ra, off, imm);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, isa::Endian endian,
                         reg::SrcReg *ra, int64_t off, reg::SrcReg *rb)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(ra != nullptr);
  assert(rb != nullptr);

  init_erir(endian, ra, off, rb);
}

Instruction::Instruction(OpCode op_code, Suffix suffix, isa::Endian endian,
                         reg::SrcReg *ra, int64_t off, reg::PairReg *db)
    : op_code_(op_code),
      suffix_(suffix),
      rc_(nullptr),
      ra_(nullptr),
      rb_(nullptr),
      dc_(nullptr),
      db_(nullptr),
      condition_(nullptr),
      imm_(nullptr),
      off_(nullptr),
      pc_(nullptr),
      endian_(nullptr),
      thread_(nullptr) {
  assert(ra != nullptr);
  assert(db != nullptr);

  init_erid(endian, ra, off, db);
}

Instruction::~Instruction() {
  delete rc_;
  delete dc_;
  delete db_;
  delete condition_;
  delete imm_;
  delete off_;
  delete pc_;
  delete endian_;
}

reg::GPReg *Instruction::rc() {
  assert(rc_ != nullptr);
  return rc_;
}

reg::SrcReg *Instruction::ra() {
  assert(ra_ != nullptr);
  return ra_;
}

reg::SrcReg *Instruction::rb() {
  assert(rb_ != nullptr);
  return rb_;
}

reg::PairReg *Instruction::dc() {
  assert(dc_ != nullptr);
  return dc_;
}

reg::PairReg *Instruction::db() {
  assert(db_ != nullptr);
  return db_;
}

isa::Condition Instruction::condition() {
  assert(condition_ != nullptr);
  return *condition_;
}

abi::word::Immediate *Instruction::imm() {
  assert(imm_ != nullptr);
  return imm_;
}

abi::word::Immediate *Instruction::off() {
  assert(off_ != nullptr);
  return off_;
}

abi::word::Immediate *Instruction::pc() {
  assert(pc_ != nullptr);
  return pc_;
}

isa::Endian Instruction::endian() {
  assert(endian_ != nullptr);
  return *endian_;
}

simulator::dpu::Thread *Instruction::thread() {
  assert(thread_ != nullptr);
  return thread_;
}

void Instruction::set_thread(simulator::dpu::Thread *thread) {
  assert(thread != nullptr);
  assert(thread_ == nullptr);

  thread_ = thread;
}

void Instruction::init_rici(reg::SrcReg *ra, int64_t imm,
                            isa::Condition condition, int64_t pc) {
  assert(Instruction::rici_op_codes().count(op_code_));
  assert(suffix_ == RICI);

  ra_ = ra;
  imm_ = new word::Immediate(word::SIGNED, 16, imm);

  if (Instruction::acquire_rici_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::AcquireCC(condition).condition());
  } else if (Instruction::release_rici_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::ReleaseCC(condition).condition());
  } else if (Instruction::boot_rici_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::BootCC(condition).condition());
  } else {
    throw std::invalid_argument("");
  }

  pc_ = new word::Immediate(word::UNSIGNED,
                            util::ConfigLoader::iram_address_width(), pc);
}

void Instruction::init_rri(reg::GPReg *rc, reg::SrcReg *ra, int64_t imm) {
  assert(Instruction::rri_op_codes().count(op_code_));
  assert(suffix_ == RRI);

  rc_ = rc;
  ra_ = ra;

  if (Instruction::add_rri_op_codes().count(op_code_)) {
    imm_ = new word::Immediate(word::UNSIGNED, 32, imm);
  } else if (Instruction::asr_rri_op_codes().count(op_code_)) {
    imm_ = new word::Immediate(word::UNSIGNED, 5, imm);
  } else if (Instruction::call_rri_op_codes().count(op_code_)) {
    imm_ = new word::Immediate(word::SIGNED, 24, imm);
  } else {
    throw std::invalid_argument("");
  }
}

void Instruction::init_rric(reg::GPReg *rc, reg::SrcReg *ra, int64_t imm,
                            isa::Condition condition) {
  assert(Instruction::rric_op_codes().count(op_code_));
  assert(suffix_ == RRIC);

  rc_ = rc;
  ra_ = ra;

  if (Instruction::add_rric_op_codes().count(op_code_) or
      Instruction::sub_rric_op_codes().count(op_code_)) {
    imm_ = new word::Immediate(word::SIGNED, 24, imm);
  } else if (Instruction::asr_rric_op_codes().count(op_code_)) {
    imm_ = new word::Immediate(word::UNSIGNED, 5, imm);
  } else {
    throw std::invalid_argument("");
  }

  if (Instruction::add_rric_op_codes().count(op_code_) or
      Instruction::asr_rric_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::LogSetCC(condition).condition());
  } else if (Instruction::sub_rric_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::ExtSubSetCC(condition).condition());
  } else {
    throw std::invalid_argument("");
  }
}

void Instruction::init_rrici(reg::GPReg *rc, reg::SrcReg *ra, int64_t imm,
                             isa::Condition condition, int64_t pc) {
  assert(Instruction::rrici_op_codes().count(op_code_));
  assert(suffix_ == RRICI);

  rc_ = rc;
  ra_ = ra;

  if (Instruction::add_rrici_op_codes().count(op_code_) or
      Instruction::and_rrici_op_codes().count(op_code_) or
      Instruction::sub_rrici_op_codes().count(op_code_)) {
    imm_ = new word::Immediate(word::SIGNED, 8, imm);
  } else if (Instruction::asr_rrici_op_codes().count(op_code_)) {
    imm_ = new word::Immediate(word::UNSIGNED, 5, imm);
  } else {
    throw std::invalid_argument("");
  }

  if (Instruction::add_rrici_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::AddNZCC(condition).condition());
  } else if (Instruction::and_rrici_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::LogNZCC(condition).condition());
  } else if (Instruction::asr_rrici_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::ImmShiftNZCC(condition).condition());
  } else if (Instruction::sub_rrici_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::SubNZCC(condition).condition());
  } else {
    throw std::invalid_argument("");
  }

  pc_ = new word::Immediate(word::UNSIGNED,
                            util::ConfigLoader::iram_address_width(), pc);
}

void Instruction::init_rrif(reg::GPReg *rc, reg::SrcReg *ra, int64_t imm,
                            isa::Condition condition) {
  assert(Instruction::rrif_op_codes().count(op_code_));
  assert(suffix_ == RRIF);

  rc_ = rc;
  ra_ = ra;
  imm_ = new word::Immediate(word::SIGNED, 24, imm);
  condition_ = new isa::Condition(cc::FalseCC(condition).condition());
}

void Instruction::init_rrr(reg::GPReg *rc, reg::SrcReg *ra, reg::SrcReg *rb) {
  assert(Instruction::rrr_op_codes().count(op_code_));
  assert(suffix_ == RRR);

  rc_ = rc;
  ra_ = ra;
  rb_ = rb;
}

void Instruction::init_rrrc(reg::GPReg *rc, reg::SrcReg *ra, reg::SrcReg *rb,
                            isa::Condition condition) {
  assert(Instruction::rrrc_op_codes().count(op_code_));
  assert(suffix_ == RRRC);

  rc_ = rc;
  ra_ = ra;
  rb_ = rb;

  if (Instruction::add_rrrc_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::LogSetCC(condition).condition());
  } else if (Instruction::rsub_rrrc_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::SubSetCC(condition).condition());
  } else if (Instruction::sub_rrrc_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::ExtSubSetCC(condition).condition());
  } else {
    throw std::invalid_argument("");
  }
}

void Instruction::init_rrrci(reg::GPReg *rc, reg::SrcReg *ra, reg::SrcReg *rb,
                             isa::Condition condition, int64_t pc) {
  assert(Instruction::rrrci_op_codes().count(op_code_));
  assert(suffix_ == RRRCI);

  rc_ = rc;
  ra_ = ra;
  rb_ = rb;

  if (Instruction::add_rrrci_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::AddNZCC(condition).condition());
  } else if (Instruction::and_rrrci_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::LogNZCC(condition).condition());
  } else if (Instruction::asr_rrrci_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::ShiftNZCC(condition).condition());
  } else if (Instruction::mul_rrrci_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::MulNZCC(condition).condition());
  } else if (Instruction::rsub_rrrci_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::SubNZCC(condition).condition());
  } else {
    throw std::invalid_argument("");
  }

  pc_ = new word::Immediate(word::UNSIGNED,
                            util::ConfigLoader::iram_address_width(), pc);
}

void Instruction::init_zri(reg::SrcReg *ra, int64_t imm) {
  assert(Instruction::rri_op_codes().count(op_code_));
  assert(suffix_ == ZRI);

  ra_ = ra;

  if (Instruction::add_rri_op_codes().count(op_code_)) {
    imm_ = new word::Immediate(word::UNSIGNED, 32, imm);
  } else if (Instruction::asr_rri_op_codes().count(op_code_)) {
    imm_ = new word::Immediate(word::UNSIGNED, 5, imm);
  } else if (Instruction::call_rri_op_codes().count(op_code_)) {
    imm_ = new word::Immediate(word::SIGNED, 28, imm);
  } else {
    throw std::invalid_argument("");
  }
}

void Instruction::init_zric(reg::SrcReg *ra, int64_t imm,
                            isa::Condition condition) {
  assert(Instruction::rric_op_codes().count(op_code_));
  assert(suffix_ == ZRIC);

  ra_ = ra;

  if (Instruction::add_rric_op_codes().count(op_code_) or
      Instruction::sub_rric_op_codes().count(op_code_)) {
    imm_ = new word::Immediate(word::SIGNED, 27, imm);
  } else if (Instruction::asr_rric_op_codes().count(op_code_)) {
    imm_ = new word::Immediate(word::UNSIGNED, 5, imm);
  } else {
    throw std::invalid_argument("");
  }

  if (Instruction::add_rric_op_codes().count(op_code_) or
      Instruction::asr_rric_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::LogSetCC(condition).condition());
  } else if (Instruction::sub_rric_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::ExtSubSetCC(condition).condition());
  } else {
    throw std::invalid_argument("");
  }
}

void Instruction::init_zrici(reg::SrcReg *ra, int64_t imm,
                             isa::Condition condition, int64_t pc) {
  assert(Instruction::rrrci_op_codes().count(op_code_));
  assert(suffix_ == ZRICI);

  ra_ = ra;

  if (Instruction::add_rrici_op_codes().count(op_code_) or
      Instruction::and_rrici_op_codes().count(op_code_) or
      Instruction::sub_rrici_op_codes().count(op_code_)) {
    imm_ = new word::Immediate(word::SIGNED, 11, imm);
  } else if (Instruction::asr_rrici_op_codes().count(op_code_)) {
    imm_ = new word::Immediate(word::UNSIGNED, 5, imm);
  } else {
    throw std::invalid_argument("");
  }

  if (Instruction::add_rrici_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::AddNZCC(condition).condition());
  } else if (Instruction::and_rrici_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::LogNZCC(condition).condition());
  } else if (Instruction::asr_rrici_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::ImmShiftNZCC(condition).condition());
  } else if (Instruction::sub_rrici_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::SubNZCC(condition).condition());
  } else {
    throw std::invalid_argument("");
  }

  pc_ = new word::Immediate(word::UNSIGNED,
                            util::ConfigLoader::iram_address_width(), pc);
}

void Instruction::init_zrif(reg::SrcReg *ra, int64_t imm,
                            isa::Condition condition) {
  assert(Instruction::rrif_op_codes().count(op_code_));
  assert(suffix_ == ZRIF);

  ra_ = ra;
  imm_ = new word::Immediate(word::SIGNED, 27, imm);
  condition_ = new isa::Condition(cc::FalseCC(condition).condition());
}

void Instruction::init_zrr(reg::SrcReg *ra, reg::SrcReg *rb) {
  assert(Instruction::rrr_op_codes().count(op_code_));
  assert(suffix_ == ZRR);

  ra_ = ra;
  rb_ = rb;
}

void Instruction::init_zrrc(reg::SrcReg *ra, reg::SrcReg *rb,
                            isa::Condition condition) {
  assert(Instruction::rrrc_op_codes().count(op_code_));
  assert(suffix_ == ZRRC);

  ra_ = ra;
  rb_ = rb;

  if (Instruction::add_rrrc_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::LogSetCC(condition).condition());
  } else if (Instruction::rsub_rrrc_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::SubSetCC(condition).condition());
  } else if (Instruction::sub_rrrc_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::ExtSubSetCC(condition).condition());
  } else {
    throw std::invalid_argument("");
  }
}

void Instruction::init_zrrci(reg::SrcReg *ra, reg::SrcReg *rb,
                             isa::Condition condition, int64_t pc) {
  assert(Instruction::rrrci_op_codes().count(op_code_));
  assert(suffix_ == ZRRCI);

  ra_ = ra;
  rb_ = rb;

  if (Instruction::add_rrrci_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::AddNZCC(condition).condition());
  } else if (Instruction::and_rrrci_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::LogNZCC(condition).condition());
  } else if (Instruction::asr_rrrci_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::ShiftNZCC(condition).condition());
  } else if (Instruction::mul_rrrci_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::MulNZCC(condition).condition());
  } else if (Instruction::rsub_rrrci_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::SubNZCC(condition).condition());
  } else {
    throw std::invalid_argument("");
  }

  pc_ = new word::Immediate(word::UNSIGNED,
                            util::ConfigLoader::iram_address_width(), pc);
}

void Instruction::init_s_rri(reg::PairReg *dc, reg::SrcReg *ra, int64_t imm) {
  assert(Instruction::rri_op_codes().count(op_code_));
  assert(suffix_ == S_RRI or suffix_ == U_RRI);

  dc_ = dc;
  ra_ = ra;

  if (Instruction::add_rri_op_codes().count(op_code_)) {
    imm_ = new word::Immediate(word::UNSIGNED, 32, imm);
  } else if (Instruction::asr_rri_op_codes().count(op_code_)) {
    imm_ = new word::Immediate(word::UNSIGNED, 5, imm);
  } else if (Instruction::call_rri_op_codes().count(op_code_)) {
    imm_ = new word::Immediate(word::SIGNED, 24, imm);
  } else {
    throw std::invalid_argument("");
  }
}

void Instruction::init_s_rric(reg::PairReg *dc, reg::SrcReg *ra, int64_t imm,
                              isa::Condition condition) {
  assert(Instruction::rric_op_codes().count(op_code_));
  assert(suffix_ == S_RRIC or suffix_ == U_RRIC);

  dc_ = dc;
  ra_ = ra;

  if (Instruction::add_rric_op_codes().count(op_code_) or
      Instruction::sub_rric_op_codes().count(op_code_)) {
    imm_ = new word::Immediate(word::SIGNED, 24, imm);
  } else if (Instruction::asr_rric_op_codes().count(op_code_)) {
    imm_ = new word::Immediate(word::UNSIGNED, 5, imm);
  } else {
    throw std::invalid_argument("");
  }

  if (Instruction::add_rric_op_codes().count(op_code_) or
      Instruction::asr_rric_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::LogSetCC(condition).condition());
  } else if (Instruction::sub_rric_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::ExtSubSetCC(condition).condition());
  } else {
    throw std::invalid_argument("");
  }
}

void Instruction::init_s_rrici(reg::PairReg *dc, reg::SrcReg *ra, int64_t imm,
                               isa::Condition condition, int64_t pc) {
  assert(Instruction::rrici_op_codes().count(op_code_));
  assert(suffix_ == S_RRICI or suffix_ == U_RRICI);

  dc_ = dc;
  ra_ = ra;

  if (Instruction::add_rrici_op_codes().count(op_code_) or
      Instruction::and_rrici_op_codes().count(op_code_) or
      Instruction::sub_rrici_op_codes().count(op_code_)) {
    imm_ = new word::Immediate(word::SIGNED, 8, imm);
  } else if (Instruction::asr_rrici_op_codes().count(op_code_)) {
    imm_ = new word::Immediate(word::UNSIGNED, 5, imm);
  } else {
    throw std::invalid_argument("");
  }

  if (Instruction::add_rrici_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::AddNZCC(condition).condition());
  } else if (Instruction::and_rrici_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::LogNZCC(condition).condition());
  } else if (Instruction::asr_rrici_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::ImmShiftNZCC(condition).condition());
  } else if (Instruction::sub_rrici_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::SubNZCC(condition).condition());
  } else {
    throw std::invalid_argument("");
  }

  pc_ = new word::Immediate(word::UNSIGNED,
                            util::ConfigLoader::iram_address_width(), pc);
}

void Instruction::init_s_rrif(reg::PairReg *dc, reg::SrcReg *ra, int64_t imm,
                              isa::Condition condition) {
  assert(Instruction::rrif_op_codes().count(op_code_));
  assert(suffix_ == S_RRIF or suffix_ == U_RRIF);

  dc_ = dc;
  ra_ = ra;
  imm_ = new word::Immediate(word::SIGNED, 24, imm);
  condition_ = new isa::Condition(cc::FalseCC(condition).condition());
}

void Instruction::init_s_rrr(reg::PairReg *dc, reg::SrcReg *ra,
                             reg::SrcReg *rb) {
  assert(Instruction::rrr_op_codes().count(op_code_));
  assert(suffix_ == S_RRR or suffix_ == U_RRR);

  dc_ = dc;
  ra_ = ra;
  rb_ = rb;
}

void Instruction::init_s_rrrc(reg::PairReg *dc, reg::SrcReg *ra,
                              reg::SrcReg *rb, isa::Condition condition) {
  assert(Instruction::rrrc_op_codes().count(op_code_));
  assert(suffix_ == S_RRRC or suffix_ == U_RRRC);

  dc_ = dc;
  ra_ = ra;
  rb_ = rb;

  if (Instruction::add_rrrc_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::LogSetCC(condition).condition());
  } else if (Instruction::rsub_rrrc_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::SubSetCC(condition).condition());
  } else if (Instruction::sub_rrrc_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::ExtSubSetCC(condition).condition());
  } else {
    throw std::invalid_argument("");
  }
}

void Instruction::init_s_rrrci(reg::PairReg *dc, reg::SrcReg *ra,
                               reg::SrcReg *rb, isa::Condition condition,
                               int64_t pc) {
  assert(Instruction::rrrci_op_codes().count(op_code_));
  assert(suffix_ == S_RRRCI or suffix_ == U_RRRCI);

  dc_ = dc;
  ra_ = ra;
  rb_ = rb;

  if (Instruction::add_rrrci_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::AddNZCC(condition).condition());
  } else if (Instruction::and_rrrci_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::LogNZCC(condition).condition());
  } else if (Instruction::asr_rrrci_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::ShiftNZCC(condition).condition());
  } else if (Instruction::mul_rrrci_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::MulNZCC(condition).condition());
  } else if (Instruction::rsub_rrrci_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::SubNZCC(condition).condition());
  } else {
    throw std::invalid_argument("");
  }

  pc_ = new word::Immediate(word::UNSIGNED,
                            util::ConfigLoader::iram_address_width(), pc);
}

void Instruction::init_rr(reg::GPReg *rc, reg::SrcReg *ra) {
  assert(Instruction::rr_op_codes().count(op_code_));
  assert(suffix_ == RR);

  rc_ = rc;
  ra_ = ra;
}

void Instruction::init_rrc(reg::GPReg *rc, reg::SrcReg *ra,
                           isa::Condition condition) {
  assert(Instruction::rrc_op_codes().count(op_code_));
  assert(suffix_ == RRC);

  rc_ = rc;
  ra_ = ra;
  condition_ = new isa::Condition(cc::LogSetCC(condition).condition());
}

void Instruction::init_rrci(reg::GPReg *rc, reg::SrcReg *ra,
                            isa::Condition condition, int64_t pc) {
  assert(Instruction::rrci_op_codes().count(op_code_));
  assert(suffix_ == RRCI);

  rc_ = rc;
  ra_ = ra;

  if (Instruction::cao_rrci_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::CountNZCC(condition).condition());
  } else if (Instruction::extsb_rrci_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::LogNZCC(condition).condition());
  } else if (Instruction::time_cfg_rrci_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::TrueCC(condition).condition());
  } else {
    throw std::invalid_argument("");
  }

  pc_ = new word::Immediate(word::UNSIGNED,
                            util::ConfigLoader::iram_address_width(), pc);
}

void Instruction::init_zr(reg::SrcReg *ra) {
  assert(Instruction::rr_op_codes().count(op_code_));
  assert(suffix_ == ZR);

  ra_ = ra;
}

void Instruction::init_zrc(reg::SrcReg *ra, isa::Condition condition) {
  assert(Instruction::rrc_op_codes().count(op_code_));
  assert(suffix_ == ZRC);

  ra_ = ra;
  condition_ = new isa::Condition(cc::LogSetCC(condition).condition());
}

void Instruction::init_zrci(reg::SrcReg *ra, isa::Condition condition,
                            int64_t pc) {
  assert(Instruction::rrci_op_codes().count(op_code_));
  assert(suffix_ == ZRCI);

  ra_ = ra;

  if (Instruction::cao_rrci_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::CountNZCC(condition).condition());
  } else if (Instruction::extsb_rrci_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::LogNZCC(condition).condition());
  } else if (Instruction::time_cfg_rrci_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::TrueCC(condition).condition());
  } else {
    throw std::invalid_argument("");
  }

  pc_ = new word::Immediate(word::UNSIGNED,
                            util::ConfigLoader::iram_address_width(), pc);
}

void Instruction::init_s_rr(reg::PairReg *dc, reg::SrcReg *ra) {
  assert(Instruction::rr_op_codes().count(op_code_));
  assert(suffix_ == S_RR or suffix_ == U_RR);

  dc_ = dc;
  ra_ = ra;
}

void Instruction::init_s_rrc(reg::PairReg *dc, reg::SrcReg *ra,
                             isa::Condition condition) {
  assert(Instruction::rrc_op_codes().count(op_code_));
  assert(suffix_ == S_RRC or suffix_ == U_RRC);

  dc_ = dc;
  ra_ = ra;
  condition_ = new isa::Condition(cc::LogSetCC(condition).condition());
}

void Instruction::init_s_rrci(reg::PairReg *dc, reg::SrcReg *ra,
                              isa::Condition condition, int64_t pc) {
  assert(Instruction::rrci_op_codes().count(op_code_));
  assert(suffix_ == S_RRCI or suffix_ == U_RRCI);

  dc_ = dc;
  ra_ = ra;

  if (Instruction::cao_rrci_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::CountNZCC(condition).condition());
  } else if (Instruction::extsb_rrci_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::LogNZCC(condition).condition());
  } else if (Instruction::time_cfg_rrci_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::TrueCC(condition).condition());
  } else {
    throw std::invalid_argument("");
  }

  pc_ = new word::Immediate(word::UNSIGNED,
                            util::ConfigLoader::iram_address_width(), pc);
}

void Instruction::init_drdici(reg::PairReg *dc, reg::SrcReg *ra,
                              reg::PairReg *db, int64_t imm,
                              isa::Condition condition, int64_t pc) {
  assert(Instruction::drdici_op_codes().count(op_code_));
  assert(suffix_ == DRDICI);

  dc_ = dc;
  ra_ = ra;
  db_ = db;
  imm_ = new word::Immediate(word::UNSIGNED, 5, imm);

  if (Instruction::div_step_drdici_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::DivCC(condition).condition());
  } else if (Instruction::mul_step_drdici_op_codes().count(op_code_)) {
    condition_ = new isa::Condition(cc::BootCC(condition).condition());
  } else {
    throw std::invalid_argument("");
  }

  pc_ = new word::Immediate(word::UNSIGNED,
                            util::ConfigLoader::iram_address_width(), pc);
}

void Instruction::init_rrri(reg::GPReg *rc, reg::SrcReg *ra, reg::SrcReg *rb,
                            int64_t imm) {
  assert(Instruction::rrri_op_codes().count(op_code_));
  assert(suffix_ == RRRI);

  rc_ = rc;
  ra_ = ra;
  rb_ = rb;
  imm_ = new word::Immediate(word::UNSIGNED, 5, imm);
}

void Instruction::init_rrrici(reg::GPReg *rc, reg::SrcReg *ra, reg::SrcReg *rb,
                              int64_t imm, isa::Condition condition,
                              int64_t pc) {
  assert(Instruction::rrrici_op_codes().count(op_code_));
  assert(suffix_ == RRRICI);

  rc_ = rc;
  ra_ = ra;
  rb_ = rb;
  imm_ = new word::Immediate(word::UNSIGNED, 5, imm);
  condition_ = new isa::Condition(cc::DivNZCC(condition).condition());
  pc_ = new word::Immediate(word::UNSIGNED,
                            util::ConfigLoader::iram_address_width(), pc);
}

void Instruction::init_zrri(reg::SrcReg *ra, reg::SrcReg *rb, int64_t imm) {
  assert(Instruction::rrri_op_codes().count(op_code_));
  assert(suffix_ == ZRRI);

  ra_ = ra;
  rb_ = rb;
  imm_ = new word::Immediate(word::UNSIGNED, 5, imm);
}

void Instruction::init_zrrici(reg::SrcReg *ra, reg::SrcReg *rb, int64_t imm,
                              isa::Condition condition, int64_t pc) {
  assert(Instruction::rrrici_op_codes().count(op_code_));
  assert(suffix_ == ZRRICI);

  ra_ = ra;
  rb_ = rb;
  imm_ = new word::Immediate(word::UNSIGNED, 5, imm);
  condition_ = new isa::Condition(cc::DivNZCC(condition).condition());
  pc_ = new word::Immediate(word::UNSIGNED,
                            util::ConfigLoader::iram_address_width(), pc);
}

void Instruction::init_s_rrri(reg::PairReg *dc, reg::SrcReg *ra,
                              reg::SrcReg *rb, int64_t imm) {
  assert(Instruction::rrri_op_codes().count(op_code_));
  assert(suffix_ == S_RRRI or suffix_ == U_RRRI);

  dc_ = dc;
  ra_ = ra;
  rb_ = rb;
  imm_ = new word::Immediate(word::UNSIGNED, 5, imm);
}

void Instruction::init_s_rrrici(reg::PairReg *dc, reg::SrcReg *ra,
                                reg::SrcReg *rb, int64_t imm,
                                isa::Condition condition, int64_t pc) {
  assert(Instruction::rrrici_op_codes().count(op_code_));
  assert(suffix_ == S_RRRICI or suffix_ == U_RRRICI);

  dc_ = dc;
  ra_ = ra;
  rb_ = rb;
  imm_ = new word::Immediate(word::UNSIGNED, 5, imm);
  condition_ = new isa::Condition(cc::DivNZCC(condition).condition());
  pc_ = new word::Immediate(word::UNSIGNED,
                            util::ConfigLoader::iram_address_width(), pc);
}

void Instruction::init_rir(reg::GPReg *rc, int64_t imm, reg::SrcReg *ra) {
  assert(Instruction::rir_op_codes().count(op_code_));
  assert(suffix_ == RIR);

  rc_ = rc;
  imm_ = new word::Immediate(word::UNSIGNED, 32, imm);
  ra_ = ra;
}

void Instruction::init_rirc(reg::GPReg *rc, int64_t imm, reg::SrcReg *ra,
                            isa::Condition condition) {
  assert(Instruction::rirc_op_codes().count(op_code_));
  assert(suffix_ == RIRC);

  rc_ = rc;
  imm_ = new word::Immediate(word::SIGNED, 24, imm);
  ra_ = ra;
  condition_ = new isa::Condition(cc::SubSetCC(condition).condition());
}

void Instruction::init_rirci(reg::GPReg *rc, int64_t imm, reg::SrcReg *ra,
                             isa::Condition condition, int64_t pc) {
  assert(Instruction::rirci_op_codes().count(op_code_));
  assert(suffix_ == RIRCI);

  rc_ = rc;
  imm_ = new word::Immediate(word::SIGNED, 8, imm);
  ra_ = ra;
  condition_ = new isa::Condition(cc::SubNZCC(condition).condition());
  pc_ = new word::Immediate(word::UNSIGNED,
                            util::ConfigLoader::iram_address_width(), pc);
}

void Instruction::init_zir(int64_t imm, reg::SrcReg *ra) {
  assert(Instruction::rir_op_codes().count(op_code_));
  assert(suffix_ == ZIR);

  imm_ = new word::Immediate(word::UNSIGNED, 32, imm);
  ra_ = ra;
}

void Instruction::init_zirc(int64_t imm, reg::SrcReg *ra,
                            isa::Condition condition) {
  assert(Instruction::rirc_op_codes().count(op_code_));
  assert(suffix_ == ZIRC);

  imm_ = new word::Immediate(word::SIGNED, 27, imm);
  ra_ = ra;
  condition_ = new isa::Condition(cc::SubSetCC(condition).condition());
}

void Instruction::init_zirci(int64_t imm, reg::SrcReg *ra,
                             isa::Condition condition, int64_t pc) {
  assert(Instruction::rirci_op_codes().count(op_code_));
  assert(suffix_ == ZIRCI);

  imm_ = new word::Immediate(word::SIGNED, 11, imm);
  ra_ = ra;
  condition_ = new isa::Condition(cc::SubNZCC(condition).condition());
  pc_ = new word::Immediate(word::UNSIGNED,
                            util::ConfigLoader::iram_address_width(), pc);
}

void Instruction::init_s_rirc(reg::PairReg *dc, int64_t imm, reg::SrcReg *ra,
                              isa::Condition condition) {
  assert(Instruction::rirc_op_codes().count(op_code_));
  assert(suffix_ == S_RIRC or suffix_ == U_RIRC);

  dc_ = dc;
  imm_ = new word::Immediate(word::SIGNED, 24, imm);
  ra_ = ra;
  condition_ = new isa::Condition(cc::SubSetCC(condition).condition());
}

void Instruction::init_s_rirci(reg::PairReg *dc, int64_t imm, reg::SrcReg *ra,
                               isa::Condition condition, int64_t pc) {
  assert(Instruction::rirci_op_codes().count(op_code_));
  assert(suffix_ == S_RIRCI or suffix_ == U_RIRCI);

  dc_ = dc;
  imm_ = new word::Immediate(word::SIGNED, 8, imm);
  ra_ = ra;
  condition_ = new isa::Condition(cc::SubNZCC(condition).condition());
  pc_ = new word::Immediate(word::UNSIGNED,
                            util::ConfigLoader::iram_address_width(), pc);
}

void Instruction::init_r(reg::GPReg *rc) {
  assert(Instruction::r_op_codes().count(op_code_));
  assert(suffix_ == R);

  rc_ = rc;
}

void Instruction::init_rci(reg::GPReg *rc, isa::Condition condition,
                           int64_t pc) {
  assert(Instruction::rci_op_codes().count(op_code_));
  assert(suffix_ == RCI);

  rc_ = rc;
  condition_ = new isa::Condition(cc::TrueCC(condition).condition());
  pc_ = new word::Immediate(word::UNSIGNED,
                            util::ConfigLoader::iram_address_width(), pc);
}

void Instruction::init_z() {
  assert(Instruction::r_op_codes().count(op_code_) or op_code_ == NOP);
  assert(suffix_ == Z);
}

void Instruction::init_zci(isa::Condition condition, int64_t pc) {
  assert(Instruction::rci_op_codes().count(op_code_));
  assert(suffix_ == ZCI);

  condition_ = new isa::Condition(cc::TrueCC(condition).condition());
  pc_ = new word::Immediate(word::UNSIGNED,
                            util::ConfigLoader::iram_address_width(), pc);
}

void Instruction::init_s_r(reg::PairReg *dc) {
  assert(Instruction::r_op_codes().count(op_code_));
  assert(suffix_ == S_R or suffix_ == U_R);

  dc_ = dc;
}

void Instruction::init_s_rci(reg::PairReg *dc, isa::Condition condition,
                             int64_t pc) {
  assert(Instruction::rci_op_codes().count(op_code_));
  assert(suffix_ == S_RCI or suffix_ == U_RCI);

  dc_ = dc;
  condition_ = new isa::Condition(cc::TrueCC(condition).condition());
  pc_ = new word::Immediate(word::UNSIGNED,
                            util::ConfigLoader::iram_address_width(), pc);
}

void Instruction::init_ci(isa::Condition condition, int64_t pc) {
  assert(Instruction::ci_op_codes().count(op_code_));
  assert(suffix_ == CI);

  condition_ = new isa::Condition(cc::BootCC(condition).condition());
  pc_ = new word::Immediate(word::UNSIGNED,
                            util::ConfigLoader::iram_address_width(), pc);
}

void Instruction::init_i(int64_t imm) {
  assert(Instruction::i_op_codes().count(op_code_));
  assert(suffix_ == I);

  imm_ = new word::Immediate(word::SIGNED, 24, imm);
}

void Instruction::init_ddci(reg::PairReg *dc, reg::PairReg *db,
                            isa::Condition condition, int64_t pc) {
  assert(Instruction::ddci_op_codes().count(op_code_));
  assert(suffix_ == DDCI);

  dc_ = dc;
  db_ = db;
  condition_ = new isa::Condition(cc::TrueFalseCC(condition).condition());
  pc_ = new word::Immediate(word::UNSIGNED,
                            util::ConfigLoader::iram_address_width(), pc);
}

void Instruction::init_erri(isa::Endian endian, reg::GPReg *rc, reg::SrcReg *ra,
                            int64_t off) {
  assert(Instruction::erri_op_codes().count(op_code_));
  assert(suffix_ == ERRI);

  endian_ = new isa::Endian(endian);
  rc_ = rc;
  ra_ = ra;
  off_ = new word::Immediate(word::SIGNED, 24, off);
}

void Instruction::init_s_erri(isa::Endian endian, reg::PairReg *dc,
                              reg::SrcReg *ra, int64_t off) {
  assert(Instruction::erri_op_codes().count(op_code_));
  assert(suffix_ == S_ERRI or suffix_ == U_ERRI);

  endian_ = new isa::Endian(endian);
  dc_ = dc;
  ra_ = ra;
  off_ = new word::Immediate(word::SIGNED, 24, off);
}

void Instruction::init_edri(isa::Endian endian, reg::PairReg *dc,
                            reg::SrcReg *ra, int64_t off) {
  assert(Instruction::edri_op_codes().count(op_code_));
  assert(suffix_ == EDRI);

  endian_ = new isa::Endian(endian);
  dc_ = dc;
  ra_ = ra;
  off_ = new word::Immediate(word::SIGNED, 24, off);
}

void Instruction::init_erii(isa::Endian endian, reg::SrcReg *ra, int64_t off,
                            int64_t imm) {
  assert(Instruction::erii_op_codes().count(op_code_));
  assert(suffix_ == ERII);

  endian_ = new isa::Endian(endian);
  ra_ = ra;
  off_ = new word::Immediate(word::SIGNED, 24, off);
  imm_ = new word::Immediate(word::SIGNED, 16, imm);
}

void Instruction::init_erir(isa::Endian endian, reg::SrcReg *ra, int64_t off,
                            reg::SrcReg *rb) {
  assert(Instruction::erir_op_codes().count(op_code_));
  assert(suffix_ == ERIR);

  endian_ = new isa::Endian(endian);
  ra_ = ra;
  off_ = new word::Immediate(word::SIGNED, 24, off);
  rb_ = rb;
}

void Instruction::init_erid(isa::Endian endian, reg::SrcReg *ra, int64_t off,
                            reg::PairReg *db) {
  assert(Instruction::erid_op_codes().count(op_code_));
  assert(suffix_ == ERID);

  endian_ = new isa::Endian(endian);
  ra_ = ra;
  off_ = new word::Immediate(word::SIGNED, 24, off);
  db_ = db;
}

void Instruction::init_dma_rri(reg::SrcReg *ra, reg::SrcReg *rb, int64_t imm) {
  assert(Instruction::dma_rri_op_codes().count(op_code_));
  assert(suffix_ == DMA_RRI);

  ra_ = ra;
  rb_ = rb;
  imm_ = new word::Immediate(word::UNSIGNED, 8, imm);
}

}  // namespace upmem_sim::abi::instruction
