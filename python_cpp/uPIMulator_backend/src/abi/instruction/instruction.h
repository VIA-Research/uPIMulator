#ifndef UPMEM_SIM_ABI_ISA_INSTRUCTION_INSTRUCTION_H_
#define UPMEM_SIM_ABI_ISA_INSTRUCTION_INSTRUCTION_H_

#include <set>

#include "abi/instruction/op_code.h"
#include "abi/instruction/suffix.h"
#include "abi/isa/condition.h"
#include "abi/isa/endian.h"
#include "abi/reg/pair_reg.h"
#include "abi/reg/src_reg.h"
#include "abi/word/immediate.h"
#include "simulator/dpu/thread.h"

namespace upmem_sim::abi::instruction {

class Instruction {
 public:
  static std::set<OpCode> acquire_rici_op_codes() { return {ACQUIRE}; }
  static std::set<OpCode> release_rici_op_codes() { return {RELEASE}; }
  static std::set<OpCode> boot_rici_op_codes() { return {BOOT, RESUME}; }
  static std::set<OpCode> rici_op_codes() {
    std::set<OpCode> acquire_rici_op_codes =
        Instruction::acquire_rici_op_codes();
    std::set<OpCode> release_rici_op_codes =
        Instruction::release_rici_op_codes();
    std::set<OpCode> boot_rici_op_codes = Instruction::boot_rici_op_codes();

    std::set<OpCode> rici_op_codes = {};
    rici_op_codes.insert(acquire_rici_op_codes.begin(),
                         acquire_rici_op_codes.end());
    rici_op_codes.insert(release_rici_op_codes.begin(),
                         release_rici_op_codes.end());
    rici_op_codes.insert(boot_rici_op_codes.begin(), boot_rici_op_codes.end());
    return rici_op_codes;
  }

  static std::set<OpCode> add_rri_op_codes() {
    return {ADD, ADDC, AND, OR, XOR};
  }
  static std::set<OpCode> asr_rri_op_codes() {
    return {ASR, LSL, LSL1, LSL1X, LSLX, LSR, LSR1, LSR1X, LSRX, ROL, ROR};
  }
  static std::set<OpCode> call_rri_op_codes() { return {CALL}; }
  static std::set<OpCode> rri_op_codes() {
    std::set<OpCode> add_rri_op_codes = Instruction::add_rri_op_codes();
    std::set<OpCode> asr_rri_op_codes = Instruction::asr_rri_op_codes();
    std::set<OpCode> call_rri_op_codes = Instruction::call_rri_op_codes();

    std::set<OpCode> rri_op_codes = {};
    rri_op_codes.insert(add_rri_op_codes.begin(), add_rri_op_codes.end());
    rri_op_codes.insert(asr_rri_op_codes.begin(), asr_rri_op_codes.end());
    rri_op_codes.insert(call_rri_op_codes.begin(), call_rri_op_codes.end());
    return rri_op_codes;
  }

  static std::set<OpCode> add_rric_op_codes() {
    return {ADD, ADDC, AND, ANDN, NAND, NOR, NXOR, OR, ORN, XOR, HASH};
  }
  static std::set<OpCode> asr_rric_op_codes() {
    return {ASR, LSL, LSL1, LSL1X, LSLX, LSR, LSR1, LSR1X, LSRX, ROL, ROR};
  }
  static std::set<OpCode> sub_rric_op_codes() { return {SUB, SUBC}; }
  static std::set<OpCode> rric_op_codes() {
    std::set<OpCode> add_rric_op_codes = Instruction::add_rric_op_codes();
    std::set<OpCode> asr_rric_op_codes = Instruction::asr_rric_op_codes();
    std::set<OpCode> sub_rric_op_codes = Instruction::sub_rric_op_codes();

    std::set<OpCode> rric_op_codes = {};
    rric_op_codes.insert(add_rric_op_codes.begin(), add_rric_op_codes.end());
    rric_op_codes.insert(asr_rric_op_codes.begin(), asr_rric_op_codes.end());
    rric_op_codes.insert(sub_rric_op_codes.begin(), sub_rric_op_codes.end());
    return rric_op_codes;
  }

  static std::set<OpCode> add_rrici_op_codes() { return {ADD, ADDC}; }
  static std::set<OpCode> and_rrici_op_codes() {
    return {AND, ANDN, NAND, NOR, NXOR, OR, ORN, XOR, HASH};
  }
  static std::set<OpCode> asr_rrici_op_codes() {
    return {ASR, LSL, LSL1, LSL1X, LSLX, LSR, LSR1, LSR1X, LSRX, ROL, ROR};
  }
  static std::set<OpCode> sub_rrici_op_codes() { return {SUB, SUBC}; }
  static std::set<OpCode> rrici_op_codes() {
    std::set<OpCode> add_rrici_op_codes = Instruction::add_rrici_op_codes();
    std::set<OpCode> and_rrici_op_codes = Instruction::and_rrici_op_codes();
    std::set<OpCode> asr_rrici_op_codes = Instruction::asr_rrici_op_codes();
    std::set<OpCode> sub_rrici_op_codes = Instruction::sub_rrici_op_codes();

    std::set<OpCode> rrici_op_codes = {};
    rrici_op_codes.insert(add_rrici_op_codes.begin(), add_rrici_op_codes.end());
    rrici_op_codes.insert(and_rrici_op_codes.begin(), and_rrici_op_codes.end());
    rrici_op_codes.insert(asr_rrici_op_codes.begin(), asr_rrici_op_codes.end());
    rrici_op_codes.insert(sub_rrici_op_codes.begin(), sub_rrici_op_codes.end());
    return rrici_op_codes;
  }

  static std::set<OpCode> rrif_op_codes() {
    return {ADD, ADDC, AND, ANDN, NAND, NOR, NXOR,
            OR,  ORN,  SUB, SUBC, XOR,  HASH};
  }

  static std::set<OpCode> rrr_op_codes() {
    return {ADD,       ADDC,      AND,       ANDN,      ASR,       CMPB4,
            LSL,       LSL1,      LSL1X,     LSLX,      LSR,       LSR1,
            LSR1X,     LSRX,      MUL_SH_SH, MUL_SH_SL, MUL_SH_UH, MUL_SH_UL,
            MUL_SL_SH, MUL_SL_SL, MUL_SL_UH, MUL_SL_UL, MUL_UH_UH, MUL_UH_UL,
            MUL_UL_UH, MUL_UL_UL, NAND,      NOR,       NXOR,      OR,
            ORN,       ROL,       ROR,       RSUB,      RSUBC,     SUB,
            SUBC,      XOR,       HASH,      CALL};
  }

  static std::set<OpCode> add_rrrc_op_codes() {
    return {ADD,       ADDC,      AND,       ANDN,      ASR,       CMPB4,
            LSL,       LSL1,      LSL1X,     LSLX,      LSR,       LSR1,
            LSR1X,     LSRX,      MUL_SH_SH, MUL_SH_SL, MUL_SH_UH, MUL_SH_UL,
            MUL_SL_SH, MUL_SL_SL, MUL_SL_UH, MUL_SL_UL, MUL_UH_UH, MUL_UH_UL,
            MUL_UL_UH, MUL_UL_UL, NAND,      NOR,       NXOR,      ROL,
            ROR,       OR,        ORN,       XOR,       HASH,      CALL};
  }
  static std::set<OpCode> rsub_rrrc_op_codes() { return {RSUB, RSUBC}; }
  static std::set<OpCode> sub_rrrc_op_codes() { return {SUB, SUBC}; }
  static std::set<OpCode> rrrc_op_codes() {
    std::set<OpCode> add_rrrc_op_codes = Instruction::add_rrrc_op_codes();
    std::set<OpCode> rsub_rrrc_op_codes = Instruction::rsub_rrrc_op_codes();
    std::set<OpCode> sub_rrrc_op_codes = Instruction::sub_rrrc_op_codes();

    std::set<OpCode> rrrc_op_codes = {};
    rrrc_op_codes.insert(add_rrrc_op_codes.begin(), add_rrrc_op_codes.end());
    rrrc_op_codes.insert(rsub_rrrc_op_codes.begin(), rsub_rrrc_op_codes.end());
    rrrc_op_codes.insert(sub_rrrc_op_codes.begin(), sub_rrrc_op_codes.end());
    return rrrc_op_codes;
  }

  static std::set<OpCode> add_rrrci_op_codes() { return {ADD, ADDC}; }
  static std::set<OpCode> and_rrrci_op_codes() {
    return {AND, ANDN, NAND, NOR, NXOR, OR, ORN, XOR, HASH};
  }
  static std::set<OpCode> asr_rrrci_op_codes() {
    return {ASR, CMPB4, LSL,   LSL1, LSL1X, LSLX,
            LSR, LSR1,  LSR1X, LSRX, ROL,   ROR};
  }
  static std::set<OpCode> mul_rrrci_op_codes() {
    return {MUL_SH_SH, MUL_SH_SL, MUL_SH_UH, MUL_SH_UL, MUL_SL_SH, MUL_SL_SL,
            MUL_SL_UH, MUL_SL_UL, MUL_UH_UH, MUL_UH_UL, MUL_UL_UH, MUL_UL_UL};
  }
  static std::set<OpCode> rsub_rrrci_op_codes() {
    return {RSUB, RSUBC, SUB, SUBC};
  }
  static std::set<OpCode> rrrci_op_codes() {
    std::set<OpCode> add_rrrci_op_codes = Instruction::add_rrrci_op_codes();
    std::set<OpCode> and_rrrci_op_codes = Instruction::and_rrrci_op_codes();
    std::set<OpCode> asr_rrrci_op_codes = Instruction::asr_rrici_op_codes();
    std::set<OpCode> mul_rrrci_op_codes = Instruction::mul_rrrci_op_codes();
    std::set<OpCode> rsub_rrrci_op_codes = Instruction::rsub_rrrci_op_codes();

    std::set<OpCode> rrrci_op_codes = {};
    rrrci_op_codes.insert(add_rrrci_op_codes.begin(), add_rrrci_op_codes.end());
    rrrci_op_codes.insert(and_rrrci_op_codes.begin(), and_rrrci_op_codes.end());
    rrrci_op_codes.insert(asr_rrrci_op_codes.begin(), asr_rrrci_op_codes.end());
    rrrci_op_codes.insert(mul_rrrci_op_codes.begin(), mul_rrrci_op_codes.end());
    rrrci_op_codes.insert(rsub_rrrci_op_codes.begin(),
                          rsub_rrrci_op_codes.end());
    return rrrci_op_codes;
  }

  static std::set<OpCode> rr_op_codes() {
    return {CAO, CLO, CLS, CLZ, EXTSB, EXTSH, EXTUB, EXTUH, SATS, TIME_CFG};
  }

  static std::set<OpCode> rrc_op_codes() {
    return {CAO, CLO, CLS, CLZ, EXTSB, EXTSH, EXTUB, EXTUH, SATS};
  }

  static std::set<OpCode> cao_rrci_op_codes() { return {CAO, CLO, CLS, CLZ}; }
  static std::set<OpCode> extsb_rrci_op_codes() {
    return {EXTSB, EXTSH, EXTUB, EXTUH, SATS};
  }
  static std::set<OpCode> time_cfg_rrci_op_codes() { return {TIME_CFG}; }
  static std::set<OpCode> rrci_op_codes() {
    std::set<OpCode> cao_rrci_op_codes = Instruction::cao_rrci_op_codes();
    std::set<OpCode> extsb_rrci_op_codes = Instruction::extsb_rrci_op_codes();
    std::set<OpCode> time_cfg_rrci_op_codes =
        Instruction::time_cfg_rrci_op_codes();

    std::set<OpCode> rrci_op_codes = {};
    rrci_op_codes.insert(cao_rrci_op_codes.begin(), cao_rrci_op_codes.end());
    rrci_op_codes.insert(extsb_rrci_op_codes.begin(),
                         extsb_rrci_op_codes.end());
    rrci_op_codes.insert(time_cfg_rrci_op_codes.begin(),
                         time_cfg_rrci_op_codes.end());
    return rrci_op_codes;
  }

  static std::set<OpCode> div_step_drdici_op_codes() { return {DIV_STEP}; }
  static std::set<OpCode> mul_step_drdici_op_codes() { return {MUL_STEP}; }
  static std::set<OpCode> drdici_op_codes() {
    std::set<OpCode> div_step_drdici_op_codes =
        Instruction::div_step_drdici_op_codes();
    std::set<OpCode> mul_step_drdici_op_codes =
        Instruction::mul_step_drdici_op_codes();

    std::set<OpCode> drdici_op_codes = {};
    drdici_op_codes.insert(div_step_drdici_op_codes.begin(),
                           div_step_drdici_op_codes.end());
    drdici_op_codes.insert(mul_step_drdici_op_codes.begin(),
                           mul_step_drdici_op_codes.end());
    return drdici_op_codes;
  }

  static std::set<OpCode> rrri_op_codes() {
    return {LSL_ADD, LSL_SUB, LSR_ADD, ROL_ADD};
  }
  static std::set<OpCode> rrrici_op_codes() {
    return {LSL_ADD, LSL_SUB, LSR_ADD, ROL_ADD};
  }

  static std::set<OpCode> rir_op_codes() { return {SUB, SUBC}; }
  static std::set<OpCode> rirc_op_codes() { return {SUB, SUBC}; }
  static std::set<OpCode> rirci_op_codes() { return {SUB, SUBC}; }

  static std::set<OpCode> r_op_codes() { return {TIME}; }
  static std::set<OpCode> rci_op_codes() { return {TIME}; }

  static std::set<OpCode> ci_op_codes() { return {STOP}; }
  static std::set<OpCode> i_op_codes() { return {FAULT}; }

  static std::set<OpCode> movd_ddci_op_codes() { return {MOVD}; }
  static std::set<OpCode> swapd_ddci_op_codes() { return {SWAPD}; }
  static std::set<OpCode> ddci_op_codes() {
    std::set<OpCode> movd_ddci_op_codes = Instruction::movd_ddci_op_codes();
    std::set<OpCode> swapd_ddci_op_codes = Instruction::swapd_ddci_op_codes();

    std::set<OpCode> ddci_op_codes = {};
    ddci_op_codes.insert(movd_ddci_op_codes.begin(), movd_ddci_op_codes.end());
    ddci_op_codes.insert(swapd_ddci_op_codes.begin(),
                         swapd_ddci_op_codes.end());
    return ddci_op_codes;
  }

  static std::set<OpCode> erri_op_codes() { return {LBS, LBU, LHS, LHU, LW}; }
  static std::set<OpCode> edri_op_codes() { return {LD}; }

  static std::set<OpCode> erii_op_codes() {
    return {SB, SB_ID, SD, SD_ID, SH, SH_ID, SW, SW_ID, SD, SD_ID};
  }
  static std::set<OpCode> erir_op_codes() { return {SB, SH, SW}; }
  static std::set<OpCode> erid_op_codes() { return {SD}; }

  static std::set<OpCode> ldma_dma_rri_op_codes() { return {LDMA}; }
  static std::set<OpCode> ldmai_dma_rri_op_codes() { return {LDMAI}; }
  static std::set<OpCode> sdma_dma_rri_op_codes() { return {SDMA}; }
  static std::set<OpCode> dma_rri_op_codes() {
    std::set<OpCode> ldma_dma_rri_op_codes =
        Instruction::ldma_dma_rri_op_codes();
    std::set<OpCode> ldmai_dma_rri_op_codes =
        Instruction::ldmai_dma_rri_op_codes();
    std::set<OpCode> sdma_dma_rri_op_codes =
        Instruction::sdma_dma_rri_op_codes();

    std::set<OpCode> dma_rri_op_codes = {};
    dma_rri_op_codes.insert(ldma_dma_rri_op_codes.begin(),
                            ldma_dma_rri_op_codes.end());
    dma_rri_op_codes.insert(ldmai_dma_rri_op_codes.begin(),
                            ldmai_dma_rri_op_codes.end());
    dma_rri_op_codes.insert(sdma_dma_rri_op_codes.begin(),
                            sdma_dma_rri_op_codes.end());
    return dma_rri_op_codes;
  }

  explicit Instruction(OpCode op_code, Suffix suffix, reg::SrcReg *ra,
                       int64_t imm, isa::Condition condition, int64_t pc);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::GPReg *rc,
                       reg::SrcReg *ra, int64_t imm);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::GPReg *rc,
                       reg::SrcReg *ra, int64_t imm, isa::Condition condition);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::GPReg *rc,
                       reg::SrcReg *ra, int64_t imm, isa::Condition condition,
                       int64_t pc);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::GPReg *rc,
                       reg::SrcReg *ra, reg::SrcReg *rb);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::GPReg *rc,
                       reg::SrcReg *ra, reg::SrcReg *rb,
                       isa::Condition condition);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::GPReg *rc,
                       reg::SrcReg *ra, reg::SrcReg *rb,
                       isa::Condition condition, int64_t pc);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::SrcReg *ra,
                       int64_t imm);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::SrcReg *ra,
                       int64_t imm, isa::Condition condition);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::SrcReg *ra,
                       reg::SrcReg *rb);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::SrcReg *ra,
                       reg::SrcReg *rb, isa::Condition condition);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::SrcReg *ra,
                       reg::SrcReg *rb, isa::Condition condition, int64_t pc);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc,
                       reg::SrcReg *ra, int64_t imm);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc,
                       reg::SrcReg *ra, int64_t imm, isa::Condition condition);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc,
                       reg::SrcReg *ra, int64_t imm, isa::Condition condition,
                       int64_t pc);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc,
                       reg::SrcReg *ra, reg::SrcReg *rb);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc,
                       reg::SrcReg *ra, reg::SrcReg *rb,
                       isa::Condition condition);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc,
                       reg::SrcReg *ra, reg::SrcReg *rb,
                       isa::Condition condition, int64_t pc);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::GPReg *rc,
                       reg::SrcReg *ra);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::GPReg *rc,
                       reg::SrcReg *ra, isa::Condition condition);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::GPReg *rc,
                       reg::SrcReg *ra, isa::Condition condition, int64_t pc);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::SrcReg *ra);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::SrcReg *ra,
                       isa::Condition condition);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::SrcReg *ra,
                       isa::Condition condition, int64_t pc);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc,
                       reg::SrcReg *ra);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc,
                       reg::SrcReg *ra, isa::Condition condition);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc,
                       reg::SrcReg *ra, isa::Condition condition, int64_t pc);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc,
                       reg::SrcReg *ra, reg::PairReg *db, int64_t imm,
                       isa::Condition condition, int64_t pc);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::GPReg *rc,
                       reg::SrcReg *ra, reg::SrcReg *rb, int64_t imm);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::GPReg *rc,
                       reg::SrcReg *ra, reg::SrcReg *rb, int64_t imm,
                       isa::Condition condition, int64_t pc);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::SrcReg *ra,
                       reg::SrcReg *rb, int64_t imm);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::SrcReg *ra,
                       reg::SrcReg *rb, int64_t imm, isa::Condition condition,
                       int64_t pc);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc,
                       reg::SrcReg *ra, reg::SrcReg *rb, int64_t imm);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc,
                       reg::SrcReg *ra, reg::SrcReg *rb, int64_t imm,
                       isa::Condition condition, int64_t pc);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::GPReg *rc,
                       int64_t imm, reg::SrcReg *ra);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::GPReg *rc,
                       int64_t imm, reg::SrcReg *ra, isa::Condition condition);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::GPReg *rc,
                       int64_t imm, reg::SrcReg *ra, isa::Condition condition,
                       int64_t pc);
  explicit Instruction(OpCode op_code, Suffix suffix, int64_t imm,
                       reg::SrcReg *ra);
  explicit Instruction(OpCode op_code, Suffix suffix, int64_t imm,
                       reg::SrcReg *ra, isa::Condition condition);
  explicit Instruction(OpCode op_code, Suffix suffix, int64_t imm,
                       reg::SrcReg *ra, isa::Condition condition, int64_t pc);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc,
                       int64_t imm, reg::SrcReg *ra, isa::Condition condition);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc,
                       int64_t imm, reg::SrcReg *ra, isa::Condition condition,
                       int64_t pc);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::GPReg *rc);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::GPReg *rc,
                       isa::Condition condition, int64_t pc);
  explicit Instruction(OpCode op_code, Suffix suffix);
  explicit Instruction(OpCode op_code, Suffix suffix, isa::Condition condition,
                       int64_t pc);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc,
                       isa::Condition condition, int64_t pc);
  explicit Instruction(OpCode op_code, Suffix suffix, int64_t imm);
  explicit Instruction(OpCode op_code, Suffix suffix, reg::PairReg *dc,
                       reg::PairReg *db, isa::Condition condition, int64_t pc);
  explicit Instruction(OpCode op_code, Suffix suffix, isa::Endian endian,
                       reg::GPReg *rc, reg::SrcReg *ra, int64_t off);
  explicit Instruction(OpCode op_code, Suffix suffix, isa::Endian endian,
                       reg::PairReg *dc, reg::SrcReg *ra, int64_t off);
  explicit Instruction(OpCode op_code, Suffix suffix, isa::Endian endian,
                       reg::SrcReg *ra, int64_t off, int64_t imm);
  explicit Instruction(OpCode op_code, Suffix suffix, isa::Endian endian,
                       reg::SrcReg *ra, int64_t off, reg::SrcReg *rb);
  explicit Instruction(OpCode op_code, Suffix suffix, isa::Endian endian,
                       reg::SrcReg *ra, int64_t off, reg::PairReg *db);

  ~Instruction();

  OpCode op_code() { return op_code_; }
  Suffix suffix() { return suffix_; }

  reg::GPReg *rc();
  reg::SrcReg *ra();
  reg::SrcReg *rb();

  reg::PairReg *dc();
  reg::PairReg *db();

  isa::Condition condition();

  abi::word::Immediate *imm();
  abi::word::Immediate *off();
  abi::word::Immediate *pc();
  isa::Endian endian();

  simulator::dpu::Thread *thread();
  void set_thread(simulator::dpu::Thread *thread);

 protected:
  void init_rici(reg::SrcReg *ra, int64_t imm, isa::Condition condition,
                 int64_t pc);
  void init_rri(reg::GPReg *rc, reg::SrcReg *ra, int64_t imm);
  void init_rric(reg::GPReg *rc, reg::SrcReg *ra, int64_t imm,
                 isa::Condition condition);
  void init_rrici(reg::GPReg *rc, reg::SrcReg *ra, int64_t imm,
                  isa::Condition condition, int64_t pc);
  void init_rrif(reg::GPReg *rc, reg::SrcReg *ra, int64_t imm,
                 isa::Condition condition);
  void init_rrr(reg::GPReg *rc, reg::SrcReg *ra, reg::SrcReg *rb);
  void init_rrrc(reg::GPReg *rc, reg::SrcReg *ra, reg::SrcReg *rb,
                 isa::Condition condition);
  void init_rrrci(reg::GPReg *rc, reg::SrcReg *ra, reg::SrcReg *rb,
                  isa::Condition condition, int64_t pc);
  void init_zri(reg::SrcReg *ra, int64_t imm);
  void init_zric(reg::SrcReg *ra, int64_t imm, isa::Condition condition);
  void init_zrici(reg::SrcReg *ra, int64_t imm, isa::Condition condition,
                  int64_t pc);
  void init_zrif(reg::SrcReg *ra, int64_t imm, isa::Condition condition);
  void init_zrr(reg::SrcReg *ra, reg::SrcReg *rb);
  void init_zrrc(reg::SrcReg *ra, reg::SrcReg *rb, isa::Condition condition);
  void init_zrrci(reg::SrcReg *ra, reg::SrcReg *rb, isa::Condition condition,
                  int64_t pc);
  void init_s_rri(reg::PairReg *dc, reg::SrcReg *ra, int64_t imm);
  void init_s_rric(reg::PairReg *dc, reg::SrcReg *ra, int64_t imm,
                   isa::Condition condition);
  void init_s_rrici(reg::PairReg *dc, reg::SrcReg *ra, int64_t imm,
                    isa::Condition condition, int64_t pc);
  void init_s_rrif(reg::PairReg *dc, reg::SrcReg *ra, int64_t imm,
                   isa::Condition condition);
  void init_s_rrr(reg::PairReg *dc, reg::SrcReg *ra, reg::SrcReg *rb);
  void init_s_rrrc(reg::PairReg *dc, reg::SrcReg *ra, reg::SrcReg *rb,
                   isa::Condition condition);
  void init_s_rrrci(reg::PairReg *dc, reg::SrcReg *ra, reg::SrcReg *rb,
                    isa::Condition condition, int64_t pc);
  void init_rr(reg::GPReg *rc, reg::SrcReg *ra);
  void init_rrc(reg::GPReg *rc, reg::SrcReg *ra, isa::Condition condition);
  void init_rrci(reg::GPReg *rc, reg::SrcReg *ra, isa::Condition condition,
                 int64_t pc);
  void init_zr(reg::SrcReg *ra);
  void init_zrc(reg::SrcReg *ra, isa::Condition condition);
  void init_zrci(reg::SrcReg *ra, isa::Condition condition, int64_t pc);
  void init_s_rr(reg::PairReg *dc, reg::SrcReg *ra);
  void init_s_rrc(reg::PairReg *dc, reg::SrcReg *ra, isa::Condition condition);
  void init_s_rrci(reg::PairReg *dc, reg::SrcReg *ra, isa::Condition condition,
                   int64_t pc);
  void init_drdici(reg::PairReg *dc, reg::SrcReg *ra, reg::PairReg *db,
                   int64_t imm, isa::Condition condition, int64_t pc);
  void init_rrri(reg::GPReg *rc, reg::SrcReg *ra, reg::SrcReg *rb, int64_t imm);
  void init_rrrici(reg::GPReg *rc, reg::SrcReg *ra, reg::SrcReg *rb,
                   int64_t imm, isa::Condition condition, int64_t pc);
  void init_zrri(reg::SrcReg *ra, reg::SrcReg *rb, int64_t imm);
  void init_zrrici(reg::SrcReg *ra, reg::SrcReg *rb, int64_t imm,
                   isa::Condition condition, int64_t pc);
  void init_s_rrri(reg::PairReg *dc, reg::SrcReg *ra, reg::SrcReg *rb,
                   int64_t imm);
  void init_s_rrrici(reg::PairReg *dc, reg::SrcReg *ra, reg::SrcReg *rb,
                     int64_t imm, isa::Condition condition, int64_t pc);
  void init_rir(reg::GPReg *rc, int64_t imm, reg::SrcReg *ra);
  void init_rirc(reg::GPReg *rc, int64_t imm, reg::SrcReg *ra,
                 isa::Condition condition);
  void init_rirci(reg::GPReg *rc, int64_t imm, reg::SrcReg *ra,
                  isa::Condition condition, int64_t pc);
  void init_zir(int64_t imm, reg::SrcReg *ra);
  void init_zirc(int64_t imm, reg::SrcReg *ra, isa::Condition condition);
  void init_zirci(int64_t imm, reg::SrcReg *ra, isa::Condition condition,
                  int64_t pc);
  void init_s_rirc(reg::PairReg *dc, int64_t imm, reg::SrcReg *ra,
                   isa::Condition condition);
  void init_s_rirci(reg::PairReg *dc, int64_t imm, reg::SrcReg *ra,
                    isa::Condition condition, int64_t pc);
  void init_r(reg::GPReg *rc);
  void init_rci(reg::GPReg *rc, isa::Condition condition, int64_t pc);
  void init_z();
  void init_zci(isa::Condition condition, int64_t pc);
  void init_s_r(reg::PairReg *dc);
  void init_s_rci(reg::PairReg *dc, isa::Condition condition, int64_t pc);
  void init_ci(isa::Condition condition, int64_t pc);
  void init_i(int64_t imm);
  void init_ddci(reg::PairReg *dc, reg::PairReg *db, isa::Condition condition,
                 int64_t pc);
  void init_erri(isa::Endian endian, reg::GPReg *rc, reg::SrcReg *ra,
                 int64_t off);
  void init_s_erri(isa::Endian endian, reg::PairReg *dc, reg::SrcReg *ra,
                   int64_t off);
  void init_edri(isa::Endian endian, reg::PairReg *dc, reg::SrcReg *ra,
                 int64_t off);
  void init_erii(isa::Endian endian, reg::SrcReg *ra, int64_t off, int64_t imm);
  void init_erir(isa::Endian endian, reg::SrcReg *ra, int64_t off,
                 reg::SrcReg *rb);
  void init_erid(isa::Endian endian, reg::SrcReg *ra, int64_t off,
                 reg::PairReg *db);
  void init_dma_rri(reg::SrcReg *ra, reg::SrcReg *rb, int64_t imm);

 private:
  OpCode op_code_;
  Suffix suffix_;

  reg::GPReg *rc_;
  reg::SrcReg *ra_;
  reg::SrcReg *rb_;

  reg::PairReg *dc_;
  reg::PairReg *db_;

  isa::Condition *condition_;

  abi::word::Immediate *imm_;
  abi::word::Immediate *off_;
  abi::word::Immediate *pc_;

  isa::Endian *endian_;

  simulator::dpu::Thread *thread_;
};

}  // namespace upmem_sim::abi::instruction

#endif
