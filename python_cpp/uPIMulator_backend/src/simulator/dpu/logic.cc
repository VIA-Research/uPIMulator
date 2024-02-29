#include "simulator/dpu/logic.h"

#include <cmath>
#include <functional>
#include <iostream>

#include "converter/instruction_converter.h"
#include "converter/reg_file_converter.h"
#include "simulator/dpu/alu.h"

namespace upmem_sim::simulator::dpu {

Logic::~Logic() {
  delete pipeline_;
  delete cycle_rule_;
  delete wait_instruction_q_;

  delete stat_factory_;
}

util::StatFactory *Logic::stat_factory() {
  auto stat_factory = new util::StatFactory("");

  util::StatFactory *cycle_rule_stat_factory = cycle_rule_->stat_factory();
  util::StatFactory *scheduler_stat_factory =
      scheduler_->stat_factory();

  stat_factory->merge(stat_factory_);
  stat_factory->merge(cycle_rule_stat_factory);
  stat_factory->merge(scheduler_stat_factory);

  delete cycle_rule_stat_factory;
  delete scheduler_stat_factory;

  return stat_factory;
}

void Logic::connect_scheduler(RevolverScheduler *scheduler) {
  assert(scheduler != nullptr);
  assert(scheduler_ == nullptr);

  scheduler_ = scheduler;
}

void Logic::connect_atomic(sram::Atomic *atomic) {
  assert(atomic != nullptr);
  assert(atomic_ == nullptr);

  atomic_ = atomic;
}

void Logic::connect_iram(sram::IRAM *iram) {
  assert(iram != nullptr);
  assert(iram_ == nullptr);

  iram_ = iram;
}

void Logic::connect_operand_collector(OperandCollector *operand_collector) {
  assert(operand_collector != nullptr);
  assert(operand_collector_ == nullptr);

  operand_collector_ = operand_collector;
}

void Logic::connect_dma(DMA *dma) {
  assert(dma != nullptr);
  assert(dma_ == nullptr);

  dma_ = dma;
}

void Logic::cycle() {
  stat_factory_->overwrite("mram_address", -1);        
  stat_factory_->overwrite("mram_access_thread", -1);  
  stat_factory_->overwrite("mram_access_size", -1); 

  service_scheduler();

  service_pipeline();
  service_cycle_rule();
  service_logic();

  service_dma();

  pipeline_->cycle();
  cycle_rule_->cycle();

  stat_factory_->increment("logic_cycle");
}

void Logic::service_scheduler() {
  if (pipeline_->can_push() and cycle_rule_->can_push() and
      wait_instruction_q_->can_push()) {
    // if (pipeline_->can_push() and wait_instruction_q_->can_push()) {
    Thread *thread = scheduler_->schedule();
    if (thread != nullptr) {
      int chosen_thread_id = thread->id();
      abi::instruction::Instruction *instruction =
          iram_->read(thread->reg_file()->read_pc_reg());
      instruction->set_thread(thread);
      pipeline_->push(instruction);

      if (instruction->suffix() != abi::instruction::DMA_RRI) {
        if (verbose_ >= 1) {
          std::cout << "{" << dpu_id_ << "}";
          std::cout << converter::InstructionConverter::to_string(instruction) << std::endl;
        }
        
        execute_instruction(instruction);

        if (verbose_ >= 2) {
          std::cout << converter::RegFileConverter::to_string(instruction->thread()->reg_file()) << std::endl;
        }
      } else {
        scheduler_->block(thread->id());
        instruction->thread()->reg_file()->increment_pc_reg();
        wait_instruction_q_->push(instruction);
      }

      stat_factory_->increment("num_instructions");
      stat_factory_->increment(std::to_string(instruction->thread()->id()) +
                               "_num_instructions");

      for (auto &th : scheduler_->threads()) {
        if (th->state() == Thread::BLOCK) {
          th->update_thread_status("WAIT_DATA", 1);
        } else if (th->state() == Thread::RUNNABLE and
                   chosen_thread_id != th->id()) {
          th->update_thread_status("WAIT_SCHEDULE", 1);
        }
      }
    } else {
      for (auto &th : scheduler_->threads()) {
        if (th->state() == Thread::BLOCK) {
          th->update_thread_status("WAIT_DATA", 1);
        }
      }
    }

    stat_factory_->increment("active_tasklets_" +
                  std::to_string(scheduler_->get_issuable_threads()));

  } else {
    stat_factory_->increment("backpressuer");
    stat_factory_->increment("active_tasklets_" + std::to_string(0));
  }
}

void Logic::service_pipeline() {
  if (pipeline_->can_pop() and cycle_rule_->can_push()) {
    abi::instruction::Instruction *instruction = pipeline_->pop();

    if (instruction != nullptr) {
      abi::instruction::Suffix suffix = instruction->suffix();
      if (suffix == abi::instruction::ERRI or
          suffix == abi::instruction::EDRI or
          suffix == abi::instruction::ERII or
          suffix == abi::instruction::ERIR or
          suffix == abi::instruction::ERID) {  // spm access
        instruction->thread()->update_thread_status("SPM_ACCESS",
                                                    num_pipeline_stages_);
      } else if (suffix == abi::instruction::RICI) {  // sync
        instruction->thread()->update_thread_status("WAIT_SYNC",
                                                    num_pipeline_stages_);
      } else if (suffix != abi::instruction::DMA_RRI) {
        instruction->thread()->update_thread_status("ARITHMETIC",
                                                    num_pipeline_stages_);
      }

      cycle_rule_->push(instruction);
    }
  }
}

void Logic::service_cycle_rule() {
  if (cycle_rule_->can_pop()) {
    abi::instruction::Instruction *instruction = cycle_rule_->pop();

    if (instruction->suffix() != abi::instruction::DMA_RRI) {
      delete instruction;
    } else {
      if (verbose_ >= 1) {
        std::cout << "{" << dpu_id_ << "}";
        std::cout << converter::InstructionConverter::to_string(instruction) << std::endl;
      }
      execute_instruction(instruction);

      if (verbose_ >= 2) {
        std::cout << converter::RegFileConverter::to_string(instruction->thread()->reg_file()) << std::endl;
      }
    }
  }
}

void Logic::service_logic() {}

void Logic::service_dma() {
  if (wait_instruction_q_->can_pop() and dma_->can_pop()) {
    DMACommand *dma_command = dma_->pop();
    abi::instruction::Instruction *instruction = wait_instruction_q_->pop();

    assert(dma_command->instruction() == instruction);

    scheduler_->awake(instruction->thread()->id());

    delete dma_command;
    delete instruction;
  }
}

void Logic::execute_instruction(abi::instruction::Instruction *instruction) {
  abi::instruction::Suffix suffix = instruction->suffix();

  if (suffix == abi::instruction::RICI) {
    execute_rici(instruction);
  } else if (suffix == abi::instruction::RRI) {
    execute_rri(instruction);
  } else if (suffix == abi::instruction::RRIC) {
    execute_rric(instruction);
  } else if (suffix == abi::instruction::RRICI) {
    execute_rrici(instruction);
  } else if (suffix == abi::instruction::RRIF) {
    execute_rrif(instruction);
  } else if (suffix == abi::instruction::RRR) {
    execute_rrr(instruction);
  } else if (suffix == abi::instruction::RRRC) {
    execute_rrrc(instruction);
  } else if (suffix == abi::instruction::RRRCI) {
    execute_rrrci(instruction);
  } else if (suffix == abi::instruction::ZRI) {
    execute_zri(instruction);
  } else if (suffix == abi::instruction::ZRIC) {
    execute_zric(instruction);
  } else if (suffix == abi::instruction::ZRICI) {
    execute_zrici(instruction);
  } else if (suffix == abi::instruction::ZRIF) {
    execute_zrif(instruction);
  } else if (suffix == abi::instruction::ZRR) {
    execute_zrr(instruction);
  } else if (suffix == abi::instruction::ZRRC) {
    execute_zrrc(instruction);
  } else if (suffix == abi::instruction::ZRRCI) {
    execute_zrrci(instruction);
  } else if (suffix == abi::instruction::S_RRI) {
    execute_s_rri(instruction);
  } else if (suffix == abi::instruction::S_RRIC) {
    execute_s_rric(instruction);
  } else if (suffix == abi::instruction::S_RRICI) {
    execute_s_rrici(instruction);
  } else if (suffix == abi::instruction::S_RRIF) {
    execute_s_rrif(instruction);
  } else if (suffix == abi::instruction::S_RRR) {
    execute_s_rrr(instruction);
  } else if (suffix == abi::instruction::S_RRRC) {
    execute_s_rrrc(instruction);
  } else if (suffix == abi::instruction::S_RRRCI) {
    execute_s_rrrci(instruction);
  } else if (suffix == abi::instruction::U_RRI) {
    execute_u_rri(instruction);
  } else if (suffix == abi::instruction::U_RRIC) {
    execute_u_rric(instruction);
  } else if (suffix == abi::instruction::U_RRICI) {
    execute_u_rrici(instruction);
  } else if (suffix == abi::instruction::U_RRIF) {
    execute_u_rrif(instruction);
  } else if (suffix == abi::instruction::U_RRR) {
    execute_u_rrr(instruction);
  } else if (suffix == abi::instruction::U_RRRC) {
    execute_u_rrrc(instruction);
  } else if (suffix == abi::instruction::U_RRRCI) {
    execute_u_rrrci(instruction);
  } else if (suffix == abi::instruction::RR) {
    execute_rr(instruction);
  } else if (suffix == abi::instruction::RRC) {
    execute_rrc(instruction);
  } else if (suffix == abi::instruction::RRCI) {
    execute_rrci(instruction);
  } else if (suffix == abi::instruction::ZR) {
    execute_zr(instruction);
  } else if (suffix == abi::instruction::ZRC) {
    execute_zrc(instruction);
  } else if (suffix == abi::instruction::ZRCI) {
    execute_zrci(instruction);
  } else if (suffix == abi::instruction::S_RR) {
    execute_s_rr(instruction);
  } else if (suffix == abi::instruction::S_RRC) {
    execute_s_rrc(instruction);
  } else if (suffix == abi::instruction::S_RRCI) {
    execute_s_rrci(instruction);
  } else if (suffix == abi::instruction::U_RR) {
    execute_u_rr(instruction);
  } else if (suffix == abi::instruction::U_RRC) {
    execute_u_rrc(instruction);
  } else if (suffix == abi::instruction::U_RRCI) {
    execute_u_rrci(instruction);
  } else if (suffix == abi::instruction::DRDICI) {
    execute_drdici(instruction);
  } else if (suffix == abi::instruction::RRRI) {
    execute_rrri(instruction);
  } else if (suffix == abi::instruction::RRRICI) {
    execute_rrrici(instruction);
  } else if (suffix == abi::instruction::ZRRI) {
    execute_zrri(instruction);
  } else if (suffix == abi::instruction::ZRRICI) {
    execute_zrrici(instruction);
  } else if (suffix == abi::instruction::S_RRRI) {
    execute_s_rrri(instruction);
  } else if (suffix == abi::instruction::S_RRRICI) {
    execute_s_rrrici(instruction);
  } else if (suffix == abi::instruction::U_RRRI) {
    execute_u_rrri(instruction);
  } else if (suffix == abi::instruction::U_RRRICI) {
    execute_u_rrrici(instruction);
  } else if (suffix == abi::instruction::RIR) {
    execute_rir(instruction);
  } else if (suffix == abi::instruction::RIRC) {
    execute_rirc(instruction);
  } else if (suffix == abi::instruction::RIRCI) {
    execute_rirci(instruction);
  } else if (suffix == abi::instruction::ZIR) {
    execute_zir(instruction);
  } else if (suffix == abi::instruction::ZIRC) {
    execute_zirc(instruction);
  } else if (suffix == abi::instruction::ZIRCI) {
    execute_zirci(instruction);
  } else if (suffix == abi::instruction::S_RIRC) {
    execute_s_rirc(instruction);
  } else if (suffix == abi::instruction::S_RIRCI) {
    execute_s_rirci(instruction);
  } else if (suffix == abi::instruction::U_RIRC) {
    execute_u_rirc(instruction);
  } else if (suffix == abi::instruction::U_RIRCI) {
    execute_u_rirci(instruction);
  } else if (suffix == abi::instruction::R) {
    execute_r(instruction);
  } else if (suffix == abi::instruction::RCI) {
    execute_rci(instruction);
  } else if (suffix == abi::instruction::Z) {
    execute_z(instruction);
  } else if (suffix == abi::instruction::ZCI) {
    execute_zci(instruction);
  } else if (suffix == abi::instruction::S_R) {
    execute_s_r(instruction);
  } else if (suffix == abi::instruction::S_RCI) {
    execute_s_rci(instruction);
  } else if (suffix == abi::instruction::U_R) {
    execute_u_r(instruction);
  } else if (suffix == abi::instruction::U_RCI) {
    execute_u_rci(instruction);
  } else if (suffix == abi::instruction::CI) {
    execute_ci(instruction);
  } else if (suffix == abi::instruction::I) {
    execute_i(instruction);
  } else if (suffix == abi::instruction::DDCI) {
    execute_ddci(instruction);
  } else if (suffix == abi::instruction::ERRI) {
    execute_erri(instruction);
  } else if (suffix == abi::instruction::S_ERRI) {
    execute_s_erri(instruction);
  } else if (suffix == abi::instruction::U_ERRI) {
    execute_u_erri(instruction);
  } else if (suffix == abi::instruction::EDRI) {
    execute_edri(instruction);
  } else if (suffix == abi::instruction::ERII) {
    execute_erii(instruction);
  } else if (suffix == abi::instruction::ERIR) {
    execute_erir(instruction);
  } else if (suffix == abi::instruction::ERID) {
    execute_erid(instruction);
  } else if (suffix == abi::instruction::DMA_RRI) {
    execute_dma_rri(instruction);
  } else {
    throw std::invalid_argument("");
  }
}

void Logic::execute_rici(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rici_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RICI);

  abi::instruction::OpCode op_code = instruction->op_code();
  if (abi::instruction::Instruction::acquire_rici_op_codes().count(op_code)) {
    execute_acquire_rici(instruction);
  } else if (abi::instruction::Instruction::release_rici_op_codes().count(
                 op_code)) {
    execute_release_rici(instruction);
  } else if (abi::instruction::Instruction::boot_rici_op_codes().count(
                 op_code)) {
    execute_boot_rici(instruction);
  } else {
    throw std::invalid_argument("");
  }
}

void Logic::execute_acquire_rici(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::acquire_rici_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RICI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::UNSIGNED);
  int64_t imm = instruction->imm()->value();
  Address atomic_address = ALU::atomic_address_hash(ra, imm);

  bool can_acquire = atomic_->can_acquire(atomic_address);
  if (can_acquire) {
    atomic_->acquire(atomic_address, instruction->thread()->id());
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_acquire_cc(instruction, not can_acquire);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, not can_acquire, false);
}

void Logic::execute_release_rici(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::release_rici_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RICI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::UNSIGNED);
  int64_t imm = instruction->imm()->value();
  Address atomic_address = ALU::atomic_address_hash(ra, imm);

  bool can_release =
      atomic_->can_release(atomic_address, instruction->thread()->id());
  if (can_release) {
    atomic_->release(atomic_address, instruction->thread()->id());
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_acquire_cc(instruction, not can_release);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, not can_release, false);
}

void Logic::execute_boot_rici(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::boot_rici_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RICI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::UNSIGNED);
  int64_t imm = instruction->imm()->value();
  Address thread_id = ALU::atomic_address_hash(ra, imm);

  instruction->thread()->reg_file()->clear_conditions();

  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::BOOT) {
    bool can_boot = scheduler_->boot(static_cast<ThreadID>(thread_id));
    set_boot_cc(instruction, ra, not can_boot);
    set_flags(instruction, not can_boot, false);
  } else if (op_code == abi::instruction::RESUME) {
    bool can_resume = scheduler_->awake(static_cast<ThreadID>(thread_id));
    set_boot_cc(instruction, ra, not can_resume);
    set_flags(instruction, not can_resume, false);
  } else {
    throw std::invalid_argument("");
  }

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }
}

void Logic::execute_rri(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rri_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RRI);

  abi::instruction::OpCode op_code = instruction->op_code();
  if (abi::instruction::Instruction::add_rri_op_codes().count(op_code)) {
    execute_add_rri(instruction);
  } else if (abi::instruction::Instruction::asr_rri_op_codes().count(op_code)) {
    execute_asr_rri(instruction);
  } else if (abi::instruction::Instruction::call_rri_op_codes().count(
                 op_code)) {
    execute_call_rri(instruction);
  } else {
    throw std::invalid_argument("");
  }
}

void Logic::execute_add_rri(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::add_rri_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RRI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::ADD) {
    std::tie(result, carry, overflow) = ALU::add(ra, imm);
  } else if (op_code == abi::instruction::ADDC) {
    std::tie(result, carry, overflow) = ALU::addc(
        ra, imm, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else if (op_code == abi::instruction::AND) {
    result = ALU::and_(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::OR) {
    result = ALU::or_(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::XOR) {
    result = ALU::xor_(ra, imm);
    carry = false;
    overflow = false;
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), result);
  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, carry);
}

void Logic::execute_asr_rri(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::asr_rri_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RRI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::ASR) {
    result = ALU::asr(ra, imm);
  } else if (op_code == abi::instruction::LSL) {
    result = ALU::lsl(ra, imm);
  } else if (op_code == abi::instruction::LSL1) {
    result = ALU::lsl1(ra, imm);
  } else if (op_code == abi::instruction::LSL1X) {
    result = ALU::lsl1x(ra, imm);
  } else if (op_code == abi::instruction::LSLX) {
    result = ALU::lslx(ra, imm);
  } else if (op_code == abi::instruction::LSR) {
    result = ALU::lsr(ra, imm);
  } else if (op_code == abi::instruction::LSR1) {
    result = ALU::lsr1(ra, imm);
  } else if (op_code == abi::instruction::LSR1X) {
    result = ALU::lsr1x(ra, imm);
  } else if (op_code == abi::instruction::LSRX) {
    result = ALU::lsrx(ra, imm);
  } else if (op_code == abi::instruction::ROL) {
    result = ALU::rol(ra, imm);
  } else if (op_code == abi::instruction::ROR) {
    result = ALU::ror(ra, imm);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), result);
  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, false);
}

void Logic::execute_call_rri(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::call_rri_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RRI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  Address result;
  bool carry;
  bool overflow;
  if (imm == 0) {
    std::tie(result, carry, overflow) = ALU::add(ra, imm);
  } else {
    std::tie(result, carry, overflow) =
        ALU::add(ra * abi::word::InstructionWord().size(), imm);
  }

  instruction->thread()->reg_file()->clear_conditions();

  Address pc = instruction->thread()->reg_file()->read_pc_reg();
  instruction->thread()->reg_file()->write_gp_reg(
      instruction->rc(), pc + abi::word::InstructionWord().size());

  instruction->thread()->reg_file()->write_pc_reg(result);

  set_flags(instruction, result, carry);
}

void Logic::execute_rric(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rric_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RRIC);

  abi::instruction::OpCode op_code = instruction->op_code();
  if (abi::instruction::Instruction::add_rric_op_codes().count(op_code)) {
    execute_add_rric(instruction);
  } else if (abi::instruction::Instruction::asr_rric_op_codes().count(
                 op_code)) {
    execute_asr_rric(instruction);
  } else if (abi::instruction::Instruction::sub_rric_op_codes().count(
                 op_code)) {
    execute_sub_rric(instruction);
  } else {
    throw std::invalid_argument("");
  }
}

void Logic::execute_add_rric(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::add_rric_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RRIC);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::ADD) {
    std::tie(result, carry, overflow) = ALU::add(ra, imm);
  } else if (op_code == abi::instruction::ADDC) {
    std::tie(result, carry, overflow) = ALU::addc(
        ra, imm, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else if (op_code == abi::instruction::AND) {
    result = ALU::and_(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::ANDN) {
    result = ALU::andn(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::NAND) {
    result = ALU::nand(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::NOR) {
    result = ALU::nor(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::NXOR) {
    result = ALU::nxor(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::OR) {
    result = ALU::or_(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::ORN) {
    result = ALU::orn(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::XOR) {
    result = ALU::xor_(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::HASH) {
    result = ALU::hash(ra, imm);
    carry = false;
    overflow = false;
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_log_set_cc(instruction, result);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), 1);
  } else {
    instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), 0);
  }

  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, carry);
}

void Logic::execute_asr_rric(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::asr_rric_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RRIC);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::ASR) {
    result = ALU::asr(ra, imm);
  } else if (op_code == abi::instruction::LSL) {
    result = ALU::lsl(ra, imm);
  } else if (op_code == abi::instruction::LSL1) {
    result = ALU::lsl1(ra, imm);
  } else if (op_code == abi::instruction::LSL1X) {
    result = ALU::lsl1x(ra, imm);
  } else if (op_code == abi::instruction::LSLX) {
    result = ALU::lslx(ra, imm);
  } else if (op_code == abi::instruction::LSR) {
    result = ALU::lsr(ra, imm);
  } else if (op_code == abi::instruction::LSR1) {
    result = ALU::lsr1(ra, imm);
  } else if (op_code == abi::instruction::LSR1X) {
    result = ALU::lsr1x(ra, imm);
  } else if (op_code == abi::instruction::LSRX) {
    result = ALU::lsrx(ra, imm);
  } else if (op_code == abi::instruction::ROL) {
    result = ALU::rol(ra, imm);
  } else if (op_code == abi::instruction::ROR) {
    result = ALU::ror(ra, imm);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_log_set_cc(instruction, result);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), 1);
  } else {
    instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), 0);
  }

  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, false);
}

void Logic::execute_sub_rric(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::sub_rric_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RRIC);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::SUB) {
    std::tie(result, carry, overflow) = ALU::sub(ra, imm);
  } else if (op_code == abi::instruction::SUBC) {
    std::tie(result, carry, overflow) = ALU::subc(
        ra, imm, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_ext_sub_set_cc(instruction, ra, imm, result, carry, overflow);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), 1);
  } else {
    instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), 0);
  }

  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, false);
}

void Logic::execute_rrici(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rrici_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RRICI);

  abi::instruction::OpCode op_code = instruction->op_code();
  if (abi::instruction::Instruction::add_rrici_op_codes().count(op_code)) {
    execute_add_rrici(instruction);
  } else if (abi::instruction::Instruction::and_rrici_op_codes().count(
                 op_code)) {
    execute_and_rrici(instruction);
  } else if (abi::instruction::Instruction::asr_rrici_op_codes().count(
                 op_code)) {
    execute_asr_rrici(instruction);
  } else if (abi::instruction::Instruction::sub_rrici_op_codes().count(
                 op_code)) {
    execute_sub_rrici(instruction);
  } else {
    throw std::invalid_argument("");
  }
}

void Logic::execute_add_rrici(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::add_rrici_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RRICI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::ADD) {
    std::tie(result, carry, overflow) = ALU::add(ra, imm);
  } else if (op_code == abi::instruction::ADDC) {
    std::tie(result, carry, overflow) = ALU::addc(
        ra, imm, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_add_nz_cc(instruction, ra, result, carry, overflow);

  instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), result);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, carry);
}

void Logic::execute_and_rrici(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::and_rrici_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RRICI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::AND) {
    result = ALU::and_(ra, imm);
  } else if (op_code == abi::instruction::ANDN) {
    result = ALU::andn(ra, imm);
  } else if (op_code == abi::instruction::NAND) {
    result = ALU::nand(ra, imm);
  } else if (op_code == abi::instruction::NOR) {
    result = ALU::nor(ra, imm);
  } else if (op_code == abi::instruction::NXOR) {
    result = ALU::nxor(ra, imm);
  } else if (op_code == abi::instruction::OR) {
    result = ALU::or_(ra, imm);
  } else if (op_code == abi::instruction::ORN) {
    result = ALU::orn(ra, imm);
  } else if (op_code == abi::instruction::XOR) {
    result = ALU::xor_(ra, imm);
  } else if (op_code == abi::instruction::HASH) {
    result = ALU::hash(ra, imm);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_log_nz_cc(instruction, ra, result);

  instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), result);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, false);
}

void Logic::execute_asr_rrici(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::asr_rrici_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RRICI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::ASR) {
    result = ALU::asr(ra, imm);
  } else if (op_code == abi::instruction::LSL) {
    result = ALU::lsl(ra, imm);
  } else if (op_code == abi::instruction::LSL1) {
    result = ALU::lsl1(ra, imm);
  } else if (op_code == abi::instruction::LSL1X) {
    result = ALU::lsl1x(ra, imm);
  } else if (op_code == abi::instruction::LSLX) {
    result = ALU::lslx(ra, imm);
  } else if (op_code == abi::instruction::LSR) {
    result = ALU::lsr(ra, imm);
  } else if (op_code == abi::instruction::LSR1) {
    result = ALU::lsr1(ra, imm);
  } else if (op_code == abi::instruction::LSR1X) {
    result = ALU::lsr1x(ra, imm);
  } else if (op_code == abi::instruction::LSRX) {
    result = ALU::lsrx(ra, imm);
  } else if (op_code == abi::instruction::ROL) {
    result = ALU::rol(ra, imm);
  } else if (op_code == abi::instruction::ROR) {
    result = ALU::ror(ra, imm);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_imm_shift_nz_cc(instruction, ra, result);

  instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), result);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, false);
}

void Logic::execute_sub_rrici(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::sub_rrici_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RRICI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::SUB) {
    std::tie(result, carry, overflow) = ALU::sub(ra, imm);
  } else if (op_code == abi::instruction::SUBC) {
    std::tie(result, carry, overflow) = ALU::subc(
        ra, imm, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_sub_nz_cc(instruction, ra, imm, result, carry, overflow);

  instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), result);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, carry);
}

void Logic::execute_rrif(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rrif_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RRIF);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::ADD) {
    std::tie(result, carry, overflow) = ALU::add(ra, imm);
  } else if (op_code == abi::instruction::ADDC) {
    std::tie(result, carry, overflow) = ALU::addc(
        ra, imm, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else if (op_code == abi::instruction::AND) {
    result = ALU::and_(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::ANDN) {
    result = ALU::andn(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::NAND) {
    result = ALU::nand(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::NOR) {
    result = ALU::nor(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::NXOR) {
    result = ALU::nxor(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::OR) {
    result = ALU::or_(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::ORN) {
    result = ALU::orn(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::SUB) {
    std::tie(result, carry, overflow) = ALU::sub(ra, imm);
  } else if (op_code == abi::instruction::SUBC) {
    std::tie(result, carry, overflow) = ALU::subc(
        ra, imm, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else if (op_code == abi::instruction::XOR) {
    result = ALU::xor_(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::HASH) {
    result = ALU::hash(ra, imm);
    carry = false;
    overflow = false;
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), result);
  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, carry);
}

void Logic::execute_rrr(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rrr_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RRR);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t rb = instruction->thread()->reg_file()->read_src_reg(
      instruction->rb(), abi::word::SIGNED);

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::ADD) {
    std::tie(result, carry, overflow) = ALU::add(ra, rb);
  } else if (op_code == abi::instruction::ADDC) {
    std::tie(result, carry, overflow) = ALU::addc(
        ra, rb, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else if (op_code == abi::instruction::AND) {
    result = ALU::and_(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::ANDN) {
    result = ALU::andn(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::ASR) {
    result = ALU::asr(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::CMPB4) {
    result = ALU::cmpb4(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::LSL) {
    result = ALU::lsl(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::LSL1) {
    result = ALU::lsl1(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::LSL1X) {
    result = ALU::lsl1x(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::LSLX) {
    result = ALU::lslx(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::LSR) {
    result = ALU::lsr(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::LSR1) {
    result = ALU::lsr1(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::LSR1X) {
    result = ALU::lsr1x(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::LSRX) {
    result = ALU::lsrx(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::ROL) {
    result = ALU::rol(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::ROR) {
    result = ALU::ror(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_SH_SH) {
    result = ALU::mul_sh_sh(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_SH_SL) {
    result = ALU::mul_sh_sl(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_SH_UH) {
    result = ALU::mul_sh_uh(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_SH_UL) {
    result = ALU::mul_sh_ul(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_SL_SH) {
    result = ALU::mul_sl_sh(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_SL_SL) {
    result = ALU::mul_sl_sl(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_SL_UH) {
    result = ALU::mul_sl_uh(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_SL_UL) {
    result = ALU::mul_sl_ul(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_UH_UH) {
    result = ALU::mul_uh_uh(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_UH_UL) {
    result = ALU::mul_uh_ul(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_UL_UH) {
    result = ALU::mul_ul_uh(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_UL_UL) {
    result = ALU::mul_ul_ul(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::NAND) {
    result = ALU::nand(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::NOR) {
    result = ALU::nor(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::NXOR) {
    result = ALU::nxor(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::OR) {
    result = ALU::or_(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::ORN) {
    result = ALU::orn(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::RSUB) {
    std::tie(result, carry, overflow) = ALU::sub(rb, ra);
  } else if (op_code == abi::instruction::RSUBC) {
    std::tie(result, carry, overflow) = ALU::subc(
        rb, ra, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else if (op_code == abi::instruction::SUB) {
    std::tie(result, carry, overflow) = ALU::sub(ra, rb);
  } else if (op_code == abi::instruction::SUBC) {
    std::tie(result, carry, overflow) = ALU::subc(
        ra, rb, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else if (op_code == abi::instruction::XOR) {
    result = ALU::xor_(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::HASH) {
    result = ALU::hash(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::CALL) {
    std::tie(result, carry, overflow) = ALU::add(ra, rb);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();

  if (op_code == abi::instruction::CALL) {
    Address pc = instruction->thread()->reg_file()->read_pc_reg();
    instruction->thread()->reg_file()->write_gp_reg(
        instruction->rc(), pc + abi::word::InstructionWord().size());
    instruction->thread()->reg_file()->write_pc_reg(result);
  } else {
    instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), result);
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, carry);
}

void Logic::execute_rrrc(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rrrc_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RRRC);

  abi::instruction::OpCode op_code = instruction->op_code();
  if (abi::instruction::Instruction::add_rrrc_op_codes().count(op_code)) {
    execute_add_rrrc(instruction);
  } else if (abi::instruction::Instruction::rsub_rrrc_op_codes().count(
                 op_code)) {
    execute_rsub_rrrc(instruction);
  } else if (abi::instruction::Instruction::sub_rrrc_op_codes().count(
                 op_code)) {
    execute_sub_rrrc(instruction);
  } else {
    throw std::invalid_argument("");
  }
}

void Logic::execute_add_rrrc(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::add_rrrc_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RRRC);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t rb = instruction->thread()->reg_file()->read_src_reg(
      instruction->rb(), abi::word::SIGNED);

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::ADD) {
    std::tie(result, carry, overflow) = ALU::add(ra, rb);
  } else if (op_code == abi::instruction::ADDC) {
    std::tie(result, carry, overflow) = ALU::addc(
        ra, rb, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else if (op_code == abi::instruction::AND) {
    result = ALU::and_(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::ANDN) {
    result = ALU::andn(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::ASR) {
    result = ALU::asr(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::CMPB4) {
    result = ALU::cmpb4(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::LSL) {
    result = ALU::lsl(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::LSL1) {
    result = ALU::lsl1(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::LSL1X) {
    result = ALU::lsl1x(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::LSLX) {
    result = ALU::lslx(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::LSR) {
    result = ALU::lsr(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::LSR1) {
    result = ALU::lsr1(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::LSR1X) {
    result = ALU::lsr1x(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::LSRX) {
    result = ALU::lsrx(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::ROL) {
    result = ALU::rol(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::ROR) {
    result = ALU::ror(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_SH_SH) {
    result = ALU::mul_sh_sh(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_SH_SL) {
    result = ALU::mul_sh_sl(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_SH_UH) {
    result = ALU::mul_sh_uh(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_SH_UL) {
    result = ALU::mul_sh_ul(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_SL_SH) {
    result = ALU::mul_sl_sh(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_SL_SL) {
    result = ALU::mul_sl_sl(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_SL_UH) {
    result = ALU::mul_sl_uh(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_SL_UL) {
    result = ALU::mul_sl_ul(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_UH_UH) {
    result = ALU::mul_uh_uh(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_UH_UL) {
    result = ALU::mul_uh_ul(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_UL_UH) {
    result = ALU::mul_ul_uh(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_UL_UL) {
    result = ALU::mul_ul_ul(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::NAND) {
    result = ALU::nand(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::NOR) {
    result = ALU::nor(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::NXOR) {
    result = ALU::nxor(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::OR) {
    result = ALU::or_(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::ORN) {
    result = ALU::orn(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::XOR) {
    result = ALU::xor_(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::HASH) {
    result = ALU::hash(ra, rb);
    carry = false;
    overflow = false;
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_log_set_cc(instruction, result);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), 1);
  } else {
    instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), 0);
  }

  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, carry);
}

void Logic::execute_rsub_rrrc(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rsub_rrrc_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RRRC);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t rb = instruction->thread()->reg_file()->read_src_reg(
      instruction->rb(), abi::word::SIGNED);

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::RSUB) {
    std::tie(result, carry, overflow) = ALU::sub(rb, ra);
  } else if (op_code == abi::instruction::RSUBC) {
    std::tie(result, carry, overflow) = ALU::subc(
        rb, ra, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_sub_set_cc(instruction, ra, rb, result);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), 1);
  } else {
    instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), 0);
  }

  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, carry);
}

void Logic::execute_sub_rrrc(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::sub_rrrc_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RRRC);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t rb = instruction->thread()->reg_file()->read_src_reg(
      instruction->rb(), abi::word::SIGNED);

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::SUB) {
    std::tie(result, carry, overflow) = ALU::sub(ra, rb);
  } else if (op_code == abi::instruction::SUBC) {
    std::tie(result, carry, overflow) = ALU::subc(
        ra, rb, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_ext_sub_set_cc(instruction, ra, rb, result, carry, overflow);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), 1);
  } else {
    instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), 0);
  }

  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, carry);
}

void Logic::execute_rrrci(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rrrci_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RRRCI);

  abi::instruction::OpCode op_code = instruction->op_code();
  if (abi::instruction::Instruction::add_rrrci_op_codes().count(op_code)) {
    execute_add_rrrci(instruction);
  } else if (abi::instruction::Instruction::and_rrrci_op_codes().count(
                 op_code)) {
    execute_and_rrrci(instruction);
  } else if (abi::instruction::Instruction::asr_rrrci_op_codes().count(
                 op_code)) {
    execute_asr_rrrci(instruction);
  } else if (abi::instruction::Instruction::mul_rrrci_op_codes().count(
                 op_code)) {
    execute_mul_rrrci(instruction);
  } else if (abi::instruction::Instruction::rsub_rrrci_op_codes().count(
                 op_code)) {
    execute_rsub_rrrci(instruction);
  } else {
    throw std::invalid_argument("");
  }
}

void Logic::execute_add_rrrci(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::add_rrrci_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RRRCI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t rb = instruction->thread()->reg_file()->read_src_reg(
      instruction->rb(), abi::word::SIGNED);

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::ADD) {
    std::tie(result, carry, overflow) = ALU::add(ra, rb);
  } else if (op_code == abi::instruction::ADDC) {
    std::tie(result, carry, overflow) = ALU::addc(
        ra, rb, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_add_nz_cc(instruction, ra, result, carry, overflow);

  instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), result);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, carry);
}

void Logic::execute_and_rrrci(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::and_rrrci_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RRRCI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t rb = instruction->thread()->reg_file()->read_src_reg(
      instruction->rb(), abi::word::SIGNED);

  int64_t result;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::AND) {
    result = ALU::and_(ra, rb);
  } else if (op_code == abi::instruction::ANDN) {
    result = ALU::andn(ra, rb);
  } else if (op_code == abi::instruction::NAND) {
    result = ALU::nand(ra, rb);
  } else if (op_code == abi::instruction::NOR) {
    result = ALU::nor(ra, rb);
  } else if (op_code == abi::instruction::NXOR) {
    result = ALU::nxor(ra, rb);
  } else if (op_code == abi::instruction::OR) {
    result = ALU::or_(ra, rb);
  } else if (op_code == abi::instruction::ORN) {
    result = ALU::orn(ra, rb);
  } else if (op_code == abi::instruction::XOR) {
    result = ALU::xor_(ra, rb);
  } else if (op_code == abi::instruction::HASH) {
    result = ALU::hash(ra, rb);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_log_nz_cc(instruction, ra, result);

  instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), result);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, false);
}

void Logic::execute_asr_rrrci(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::asr_rrrci_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RRRCI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t rb = instruction->thread()->reg_file()->read_src_reg(
      instruction->rb(), abi::word::SIGNED);

  int64_t result;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::ASR) {
    result = ALU::asr(ra, rb);
  } else if (op_code == abi::instruction::CMPB4) {
    result = ALU::cmpb4(ra, rb);
  } else if (op_code == abi::instruction::LSL) {
    result = ALU::lsl(ra, rb);
  } else if (op_code == abi::instruction::LSL1) {
    result = ALU::lsl1(ra, rb);
  } else if (op_code == abi::instruction::LSL1X) {
    result = ALU::lsl1x(ra, rb);
  } else if (op_code == abi::instruction::LSLX) {
    result = ALU::lslx(ra, rb);
  } else if (op_code == abi::instruction::LSR) {
    result = ALU::lsr(ra, rb);
  } else if (op_code == abi::instruction::LSR1) {
    result = ALU::lsr1(ra, rb);
  } else if (op_code == abi::instruction::LSR1X) {
    result = ALU::lsr1x(ra, rb);
  } else if (op_code == abi::instruction::LSRX) {
    result = ALU::lsrx(ra, rb);
  } else if (op_code == abi::instruction::ROL) {
    result = ALU::rol(ra, rb);
  } else if (op_code == abi::instruction::ROR) {
    result = ALU::ror(ra, rb);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_log_nz_cc(instruction, ra, result);

  instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), result);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, false);
}

void Logic::execute_mul_rrrci(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::mul_rrrci_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RRRCI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t rb = instruction->thread()->reg_file()->read_src_reg(
      instruction->rb(), abi::word::SIGNED);

  int64_t result;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::MUL_SH_SH) {
    result = ALU::mul_sh_sh(ra, rb);
  } else if (op_code == abi::instruction::MUL_SH_SL) {
    result = ALU::mul_sh_sl(ra, rb);
  } else if (op_code == abi::instruction::MUL_SH_UH) {
    result = ALU::mul_sh_uh(ra, rb);
  } else if (op_code == abi::instruction::MUL_SH_UL) {
    result = ALU::mul_sh_ul(ra, rb);
  } else if (op_code == abi::instruction::MUL_SL_SH) {
    result = ALU::mul_sl_sh(ra, rb);
  } else if (op_code == abi::instruction::MUL_SL_SL) {
    result = ALU::mul_sl_sl(ra, rb);
  } else if (op_code == abi::instruction::MUL_SL_UH) {
    result = ALU::mul_sl_uh(ra, rb);
  } else if (op_code == abi::instruction::MUL_SL_UL) {
    result = ALU::mul_sl_ul(ra, rb);
  } else if (op_code == abi::instruction::MUL_UH_UH) {
    result = ALU::mul_uh_uh(ra, rb);
  } else if (op_code == abi::instruction::MUL_UH_UL) {
    result = ALU::mul_uh_ul(ra, rb);
  } else if (op_code == abi::instruction::MUL_UL_UH) {
    result = ALU::mul_ul_uh(ra, rb);
  } else if (op_code == abi::instruction::MUL_UL_UL) {
    result = ALU::mul_ul_ul(ra, rb);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_mul_nz_cc(instruction, ra, result);

  instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), result);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, false);
}

void Logic::execute_rsub_rrrci(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rsub_rrrci_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RRRCI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t rb = instruction->thread()->reg_file()->read_src_reg(
      instruction->rb(), abi::word::SIGNED);

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::RSUB) {
    std::tie(result, carry, overflow) = ALU::sub(rb, ra);
  } else if (op_code == abi::instruction::RSUBC) {
    std::tie(result, carry, overflow) = ALU::subc(
        rb, ra, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else if (op_code == abi::instruction::SUB) {
    std::tie(result, carry, overflow) = ALU::sub(ra, rb);
  } else if (op_code == abi::instruction::SUBC) {
    std::tie(result, carry, overflow) = ALU::subc(
        ra, rb, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_sub_nz_cc(instruction, ra, rb, result, carry, overflow);

  instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), result);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, false);
}

void Logic::execute_zri(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rri_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZRI);

  abi::instruction::OpCode op_code = instruction->op_code();
  if (abi::instruction::Instruction::add_rri_op_codes().count(op_code)) {
    execute_add_zri(instruction);
  } else if (abi::instruction::Instruction::asr_rri_op_codes().count(op_code)) {
    execute_asr_zri(instruction);
  } else if (abi::instruction::Instruction::call_rri_op_codes().count(
                 op_code)) {
    execute_call_zri(instruction);
  } else {
    throw std::invalid_argument("");
  }
}

void Logic::execute_add_zri(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::add_rri_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZRI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::ADD) {
    std::tie(result, carry, overflow) = ALU::add(ra, imm);
  } else if (op_code == abi::instruction::ADDC) {
    std::tie(result, carry, overflow) = ALU::addc(
        ra, imm, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else if (op_code == abi::instruction::AND) {
    result = ALU::and_(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::OR) {
    result = ALU::or_(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::XOR) {
    result = ALU::xor_(ra, imm);
    carry = false;
    overflow = false;
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, carry);
}

void Logic::execute_asr_zri(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::asr_rri_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZRI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::ASR) {
    result = ALU::asr(ra, imm);
  } else if (op_code == abi::instruction::LSL) {
    result = ALU::lsl(ra, imm);
  } else if (op_code == abi::instruction::LSL1) {
    result = ALU::lsl1(ra, imm);
  } else if (op_code == abi::instruction::LSL1X) {
    result = ALU::lsl1x(ra, imm);
  } else if (op_code == abi::instruction::LSLX) {
    result = ALU::lslx(ra, imm);
  } else if (op_code == abi::instruction::LSR) {
    result = ALU::lsr(ra, imm);
  } else if (op_code == abi::instruction::LSR1) {
    result = ALU::lsr1(ra, imm);
  } else if (op_code == abi::instruction::LSR1X) {
    result = ALU::lsr1x(ra, imm);
  } else if (op_code == abi::instruction::LSRX) {
    result = ALU::lsrx(ra, imm);
  } else if (op_code == abi::instruction::ROL) {
    result = ALU::rol(ra, imm);
  } else if (op_code == abi::instruction::ROR) {
    result = ALU::ror(ra, imm);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, false);
}

void Logic::execute_call_zri(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::call_rri_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZRI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  Address result;
  bool carry;
  bool overflow;
  if (imm == 0) {
    std::tie(result, carry, overflow) = ALU::add(ra, imm);
  } else {
    std::tie(result, carry, overflow) =
        ALU::add(ra * abi::word::InstructionWord().size(), imm);
  }

  instruction->thread()->reg_file()->clear_conditions();
  Address pc = instruction->thread()->reg_file()->read_pc_reg();
  instruction->thread()->reg_file()->write_pc_reg(result);

  set_flags(instruction, result, carry);
}

void Logic::execute_zric(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rric_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZRIC);

  abi::instruction::OpCode op_code = instruction->op_code();
  if (abi::instruction::Instruction::add_rric_op_codes().count(op_code)) {
    execute_add_zric(instruction);
  } else if (abi::instruction::Instruction::asr_rric_op_codes().count(
                 op_code)) {
    execute_asr_zric(instruction);
  } else if (abi::instruction::Instruction::sub_rric_op_codes().count(
                 op_code)) {
    execute_sub_zric(instruction);
  } else {
    throw std::invalid_argument("");
  }
}

void Logic::execute_add_zric(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::add_rric_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZRIC);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::ADD) {
    std::tie(result, carry, overflow) = ALU::add(ra, imm);
  } else if (op_code == abi::instruction::ADDC) {
    std::tie(result, carry, overflow) = ALU::addc(
        ra, imm, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else if (op_code == abi::instruction::AND) {
    result = ALU::and_(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::ANDN) {
    result = ALU::andn(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::NAND) {
    result = ALU::nand(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::NOR) {
    result = ALU::nor(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::NXOR) {
    result = ALU::nxor(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::OR) {
    result = ALU::or_(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::ORN) {
    result = ALU::orn(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::XOR) {
    result = ALU::xor_(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::HASH) {
    result = ALU::hash(ra, imm);
    carry = false;
    overflow = false;
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_log_set_cc(instruction, result);

  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, carry);
}

void Logic::execute_asr_zric(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::asr_rric_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZRIC);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::ASR) {
    result = ALU::asr(ra, imm);
  } else if (op_code == abi::instruction::LSL) {
    result = ALU::lsl(ra, imm);
  } else if (op_code == abi::instruction::LSL1) {
    result = ALU::lsl1(ra, imm);
  } else if (op_code == abi::instruction::LSL1X) {
    result = ALU::lsl1x(ra, imm);
  } else if (op_code == abi::instruction::LSLX) {
    result = ALU::lslx(ra, imm);
  } else if (op_code == abi::instruction::LSR) {
    result = ALU::lsr(ra, imm);
  } else if (op_code == abi::instruction::LSR1) {
    result = ALU::lsr1(ra, imm);
  } else if (op_code == abi::instruction::LSR1X) {
    result = ALU::lsr1x(ra, imm);
  } else if (op_code == abi::instruction::LSRX) {
    result = ALU::lsrx(ra, imm);
  } else if (op_code == abi::instruction::ROL) {
    result = ALU::rol(ra, imm);
  } else if (op_code == abi::instruction::ROR) {
    result = ALU::ror(ra, imm);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_log_set_cc(instruction, result);

  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, false);
}

void Logic::execute_sub_zric(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::sub_rric_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZRIC);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::SUB) {
    std::tie(result, carry, overflow) = ALU::sub(ra, imm);
  } else if (op_code == abi::instruction::SUBC) {
    std::tie(result, carry, overflow) = ALU::subc(
        ra, imm, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_ext_sub_set_cc(instruction, ra, imm, result, carry, overflow);

  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, false);
}

void Logic::execute_zrici(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rrici_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZRICI);

  abi::instruction::OpCode op_code = instruction->op_code();
  if (abi::instruction::Instruction::add_rrici_op_codes().count(op_code)) {
    execute_add_zrici(instruction);
  } else if (abi::instruction::Instruction::and_rrici_op_codes().count(
                 op_code)) {
    execute_and_zrici(instruction);
  } else if (abi::instruction::Instruction::asr_rrici_op_codes().count(
                 op_code)) {
    execute_asr_zrici(instruction);
  } else if (abi::instruction::Instruction::sub_rrici_op_codes().count(
                 op_code)) {
    execute_sub_zrici(instruction);
  } else {
    throw std::invalid_argument("");
  }
}

void Logic::execute_add_zrici(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::add_rrici_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZRICI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::ADD) {
    std::tie(result, carry, overflow) = ALU::add(ra, imm);
  } else if (op_code == abi::instruction::ADDC) {
    std::tie(result, carry, overflow) = ALU::addc(
        ra, imm, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_add_nz_cc(instruction, ra, result, carry, overflow);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, carry);
}

void Logic::execute_and_zrici(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::and_rrici_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZRICI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::AND) {
    result = ALU::and_(ra, imm);
  } else if (op_code == abi::instruction::ANDN) {
    result = ALU::andn(ra, imm);
  } else if (op_code == abi::instruction::NAND) {
    result = ALU::nand(ra, imm);
  } else if (op_code == abi::instruction::NOR) {
    result = ALU::nor(ra, imm);
  } else if (op_code == abi::instruction::NXOR) {
    result = ALU::nxor(ra, imm);
  } else if (op_code == abi::instruction::OR) {
    result = ALU::or_(ra, imm);
  } else if (op_code == abi::instruction::ORN) {
    result = ALU::orn(ra, imm);
  } else if (op_code == abi::instruction::XOR) {
    result = ALU::xor_(ra, imm);
  } else if (op_code == abi::instruction::HASH) {
    result = ALU::hash(ra, imm);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_log_nz_cc(instruction, ra, result);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, false);
}

void Logic::execute_asr_zrici(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::asr_rrici_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZRICI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::ASR) {
    result = ALU::asr(ra, imm);
  } else if (op_code == abi::instruction::LSL) {
    result = ALU::lsl(ra, imm);
  } else if (op_code == abi::instruction::LSL1) {
    result = ALU::lsl1(ra, imm);
  } else if (op_code == abi::instruction::LSL1X) {
    result = ALU::lsl1x(ra, imm);
  } else if (op_code == abi::instruction::LSLX) {
    result = ALU::lslx(ra, imm);
  } else if (op_code == abi::instruction::LSR) {
    result = ALU::lsr(ra, imm);
  } else if (op_code == abi::instruction::LSR1) {
    result = ALU::lsr1(ra, imm);
  } else if (op_code == abi::instruction::LSR1X) {
    result = ALU::lsr1x(ra, imm);
  } else if (op_code == abi::instruction::LSRX) {
    result = ALU::lsrx(ra, imm);
  } else if (op_code == abi::instruction::ROL) {
    result = ALU::rol(ra, imm);
  } else if (op_code == abi::instruction::ROR) {
    result = ALU::ror(ra, imm);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_imm_shift_nz_cc(instruction, ra, result);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, false);
}

void Logic::execute_sub_zrici(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::sub_rrici_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZRICI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::SUB) {
    std::tie(result, carry, overflow) = ALU::sub(ra, imm);
  } else if (op_code == abi::instruction::SUBC) {
    std::tie(result, carry, overflow) = ALU::subc(
        ra, imm, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_sub_nz_cc(instruction, ra, imm, result, carry, overflow);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, carry);
}

void Logic::execute_zrif(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rrif_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZRIF);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::ADD) {
    std::tie(result, carry, overflow) = ALU::add(ra, imm);
  } else if (op_code == abi::instruction::ADDC) {
    std::tie(result, carry, overflow) = ALU::addc(
        ra, imm, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else if (op_code == abi::instruction::AND) {
    result = ALU::and_(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::ANDN) {
    result = ALU::andn(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::NAND) {
    result = ALU::nand(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::NOR) {
    result = ALU::nor(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::NXOR) {
    result = ALU::nxor(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::OR) {
    result = ALU::or_(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::ORN) {
    result = ALU::orn(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::SUB) {
    std::tie(result, carry, overflow) = ALU::sub(ra, imm);
  } else if (op_code == abi::instruction::SUBC) {
    std::tie(result, carry, overflow) = ALU::subc(
        ra, imm, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else if (op_code == abi::instruction::XOR) {
    result = ALU::xor_(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::HASH) {
    result = ALU::hash(ra, imm);
    carry = false;
    overflow = false;
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, carry);
}

void Logic::execute_zrr(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rrr_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZRR);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t rb = instruction->thread()->reg_file()->read_src_reg(
      instruction->rb(), abi::word::SIGNED);

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::ADD) {
    std::tie(result, carry, overflow) = ALU::add(ra, rb);
  } else if (op_code == abi::instruction::ADDC) {
    std::tie(result, carry, overflow) = ALU::addc(
        ra, rb, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else if (op_code == abi::instruction::AND) {
    result = ALU::and_(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::ANDN) {
    result = ALU::andn(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::ASR) {
    result = ALU::asr(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::CMPB4) {
    result = ALU::cmpb4(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::LSL) {
    result = ALU::lsl(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::LSL1) {
    result = ALU::lsl1(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::LSL1X) {
    result = ALU::lsl1x(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::LSLX) {
    result = ALU::lslx(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::LSR) {
    result = ALU::lsr(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::LSR1) {
    result = ALU::lsr1(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::LSR1X) {
    result = ALU::lsr1x(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::LSRX) {
    result = ALU::lsrx(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::ROL) {
    result = ALU::rol(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::ROR) {
    result = ALU::ror(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_SH_SH) {
    result = ALU::mul_sh_sh(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_SH_SL) {
    result = ALU::mul_sh_sl(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_SH_UH) {
    result = ALU::mul_sh_uh(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_SH_UL) {
    result = ALU::mul_sh_ul(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_SL_SH) {
    result = ALU::mul_sl_sh(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_SL_SL) {
    result = ALU::mul_sl_sl(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_SL_UH) {
    result = ALU::mul_sl_uh(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_SL_UL) {
    result = ALU::mul_sl_ul(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_UH_UH) {
    result = ALU::mul_uh_uh(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_UH_UL) {
    result = ALU::mul_uh_ul(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_UL_UH) {
    result = ALU::mul_ul_uh(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_UL_UL) {
    result = ALU::mul_ul_ul(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::NAND) {
    result = ALU::nand(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::NOR) {
    result = ALU::nor(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::NXOR) {
    result = ALU::nxor(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::OR) {
    result = ALU::or_(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::ORN) {
    result = ALU::orn(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::RSUB) {
    std::tie(result, carry, overflow) = ALU::sub(rb, ra);
  } else if (op_code == abi::instruction::RSUBC) {
    std::tie(result, carry, overflow) = ALU::subc(
        rb, ra, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else if (op_code == abi::instruction::SUB) {
    std::tie(result, carry, overflow) = ALU::sub(ra, rb);
  } else if (op_code == abi::instruction::SUBC) {
    std::tie(result, carry, overflow) = ALU::subc(
        ra, rb, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else if (op_code == abi::instruction::XOR) {
    result = ALU::xor_(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::HASH) {
    result = ALU::hash(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::CALL) {
    std::tie(result, carry, overflow) = ALU::add(ra, rb);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();

  if (op_code == abi::instruction::CALL) {
    instruction->thread()->reg_file()->write_pc_reg(result);
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, carry);
}

void Logic::execute_zrrc(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rrrc_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZRRC);

  abi::instruction::OpCode op_code = instruction->op_code();
  if (abi::instruction::Instruction::add_rrrc_op_codes().count(op_code)) {
    execute_add_zrrc(instruction);
  } else if (abi::instruction::Instruction::rsub_rrrc_op_codes().count(
                 op_code)) {
    execute_rsub_zrrc(instruction);
  } else if (abi::instruction::Instruction::sub_rrrc_op_codes().count(
                 op_code)) {
    execute_sub_zrrc(instruction);
  } else {
    throw std::invalid_argument("");
  }
}

void Logic::execute_add_zrrc(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::add_rrrc_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZRRC);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t rb = instruction->thread()->reg_file()->read_src_reg(
      instruction->rb(), abi::word::SIGNED);

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::ADD) {
    std::tie(result, carry, overflow) = ALU::add(ra, rb);
  } else if (op_code == abi::instruction::ADDC) {
    std::tie(result, carry, overflow) = ALU::addc(
        ra, rb, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else if (op_code == abi::instruction::AND) {
    result = ALU::and_(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::ANDN) {
    result = ALU::andn(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::ASR) {
    result = ALU::asr(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::CMPB4) {
    result = ALU::cmpb4(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::LSL) {
    result = ALU::lsl(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::LSL1) {
    result = ALU::lsl1(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::LSL1X) {
    result = ALU::lsl1x(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::LSLX) {
    result = ALU::lslx(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::LSR) {
    result = ALU::lsr(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::LSR1) {
    result = ALU::lsr1(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::LSR1X) {
    result = ALU::lsr1x(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::LSRX) {
    result = ALU::lsrx(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::ROL) {
    result = ALU::rol(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::ROR) {
    result = ALU::ror(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_SH_SH) {
    result = ALU::mul_sh_sh(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_SH_SL) {
    result = ALU::mul_sh_sl(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_SH_UH) {
    result = ALU::mul_sh_uh(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_SH_UL) {
    result = ALU::mul_sh_ul(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_SL_SH) {
    result = ALU::mul_sl_sh(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_SL_SL) {
    result = ALU::mul_sl_sl(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_SL_UH) {
    result = ALU::mul_sl_uh(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_SL_UL) {
    result = ALU::mul_sl_ul(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_UH_UH) {
    result = ALU::mul_uh_uh(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_UH_UL) {
    result = ALU::mul_uh_ul(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_UL_UH) {
    result = ALU::mul_ul_uh(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::MUL_UL_UL) {
    result = ALU::mul_ul_ul(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::NAND) {
    result = ALU::nand(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::NOR) {
    result = ALU::nor(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::NXOR) {
    result = ALU::nxor(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::OR) {
    result = ALU::or_(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::ORN) {
    result = ALU::orn(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::XOR) {
    result = ALU::xor_(ra, rb);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::HASH) {
    result = ALU::hash(ra, rb);
    carry = false;
    overflow = false;
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_log_set_cc(instruction, result);

  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, carry);
}

void Logic::execute_rsub_zrrc(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rsub_rrrc_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZRRC);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t rb = instruction->thread()->reg_file()->read_src_reg(
      instruction->rb(), abi::word::SIGNED);

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::RSUB) {
    std::tie(result, carry, overflow) = ALU::sub(rb, ra);
  } else if (op_code == abi::instruction::RSUBC) {
    std::tie(result, carry, overflow) = ALU::subc(
        rb, ra, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_sub_set_cc(instruction, ra, rb, result);

  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, carry);
}

void Logic::execute_sub_zrrc(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::sub_rrrc_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZRRC);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t rb = instruction->thread()->reg_file()->read_src_reg(
      instruction->rb(), abi::word::SIGNED);

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::SUB) {
    std::tie(result, carry, overflow) = ALU::sub(ra, rb);
  } else if (op_code == abi::instruction::SUBC) {
    std::tie(result, carry, overflow) = ALU::subc(
        ra, rb, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_ext_sub_set_cc(instruction, ra, rb, result, carry, overflow);

  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, carry);
}

void Logic::execute_zrrci(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rrrci_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZRRCI);

  abi::instruction::OpCode op_code = instruction->op_code();
  if (abi::instruction::Instruction::add_rrrci_op_codes().count(op_code)) {
    execute_add_zrrci(instruction);
  } else if (abi::instruction::Instruction::and_rrrci_op_codes().count(
                 op_code)) {
    execute_and_zrrci(instruction);
  } else if (abi::instruction::Instruction::asr_rrrci_op_codes().count(
                 op_code)) {
    execute_asr_zrrci(instruction);
  } else if (abi::instruction::Instruction::mul_rrrci_op_codes().count(
                 op_code)) {
    execute_mul_zrrci(instruction);
  } else if (abi::instruction::Instruction::rsub_rrrci_op_codes().count(
                 op_code)) {
    execute_rsub_zrrci(instruction);
  } else {
    throw std::invalid_argument("");
  }
}

void Logic::execute_add_zrrci(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::add_rrrci_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZRRCI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t rb = instruction->thread()->reg_file()->read_src_reg(
      instruction->rb(), abi::word::SIGNED);

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::ADD) {
    std::tie(result, carry, overflow) = ALU::add(ra, rb);
  } else if (op_code == abi::instruction::ADDC) {
    std::tie(result, carry, overflow) = ALU::addc(
        ra, rb, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_add_nz_cc(instruction, ra, result, carry, overflow);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, carry);
}

void Logic::execute_and_zrrci(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::and_rrrci_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZRRCI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t rb = instruction->thread()->reg_file()->read_src_reg(
      instruction->rb(), abi::word::SIGNED);

  int64_t result;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::AND) {
    result = ALU::and_(ra, rb);
  } else if (op_code == abi::instruction::ANDN) {
    result = ALU::andn(ra, rb);
  } else if (op_code == abi::instruction::NAND) {
    result = ALU::nand(ra, rb);
  } else if (op_code == abi::instruction::NOR) {
    result = ALU::nor(ra, rb);
  } else if (op_code == abi::instruction::NXOR) {
    result = ALU::nxor(ra, rb);
  } else if (op_code == abi::instruction::OR) {
    result = ALU::or_(ra, rb);
  } else if (op_code == abi::instruction::ORN) {
    result = ALU::orn(ra, rb);
  } else if (op_code == abi::instruction::XOR) {
    result = ALU::xor_(ra, rb);
  } else if (op_code == abi::instruction::HASH) {
    result = ALU::hash(ra, rb);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_log_nz_cc(instruction, ra, result);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, false);
}

void Logic::execute_asr_zrrci(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::asr_rrrci_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZRRCI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t rb = instruction->thread()->reg_file()->read_src_reg(
      instruction->rb(), abi::word::SIGNED);

  int64_t result;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::ASR) {
    result = ALU::asr(ra, rb);
  } else if (op_code == abi::instruction::CMPB4) {
    result = ALU::cmpb4(ra, rb);
  } else if (op_code == abi::instruction::LSL) {
    result = ALU::lsl(ra, rb);
  } else if (op_code == abi::instruction::LSL1) {
    result = ALU::lsl1(ra, rb);
  } else if (op_code == abi::instruction::LSL1X) {
    result = ALU::lsl1x(ra, rb);
  } else if (op_code == abi::instruction::LSLX) {
    result = ALU::lslx(ra, rb);
  } else if (op_code == abi::instruction::LSR) {
    result = ALU::lsr(ra, rb);
  } else if (op_code == abi::instruction::LSR1) {
    result = ALU::lsr1(ra, rb);
  } else if (op_code == abi::instruction::LSR1X) {
    result = ALU::lsr1x(ra, rb);
  } else if (op_code == abi::instruction::LSRX) {
    result = ALU::lsrx(ra, rb);
  } else if (op_code == abi::instruction::ROL) {
    result = ALU::rol(ra, rb);
  } else if (op_code == abi::instruction::ROR) {
    result = ALU::ror(ra, rb);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_log_nz_cc(instruction, ra, result);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, false);
}

void Logic::execute_mul_zrrci(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::mul_rrrci_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZRRCI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t rb = instruction->thread()->reg_file()->read_src_reg(
      instruction->rb(), abi::word::SIGNED);

  int64_t result;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::MUL_SH_SH) {
    result = ALU::mul_sh_sh(ra, rb);
  } else if (op_code == abi::instruction::MUL_SH_SL) {
    result = ALU::mul_sh_sl(ra, rb);
  } else if (op_code == abi::instruction::MUL_SH_UH) {
    result = ALU::mul_sh_uh(ra, rb);
  } else if (op_code == abi::instruction::MUL_SH_UL) {
    result = ALU::mul_sh_ul(ra, rb);
  } else if (op_code == abi::instruction::MUL_SL_SH) {
    result = ALU::mul_sl_sh(ra, rb);
  } else if (op_code == abi::instruction::MUL_SL_SL) {
    result = ALU::mul_sl_sl(ra, rb);
  } else if (op_code == abi::instruction::MUL_SL_UH) {
    result = ALU::mul_sl_uh(ra, rb);
  } else if (op_code == abi::instruction::MUL_SL_UL) {
    result = ALU::mul_sl_ul(ra, rb);
  } else if (op_code == abi::instruction::MUL_UH_UH) {
    result = ALU::mul_uh_uh(ra, rb);
  } else if (op_code == abi::instruction::MUL_UH_UL) {
    result = ALU::mul_uh_ul(ra, rb);
  } else if (op_code == abi::instruction::MUL_UL_UH) {
    result = ALU::mul_ul_uh(ra, rb);
  } else if (op_code == abi::instruction::MUL_UL_UL) {
    result = ALU::mul_ul_ul(ra, rb);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_mul_nz_cc(instruction, ra, result);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, false);
}

void Logic::execute_rsub_zrrci(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rsub_rrrci_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZRRCI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t rb = instruction->thread()->reg_file()->read_src_reg(
      instruction->rb(), abi::word::SIGNED);

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::RSUB) {
    std::tie(result, carry, overflow) = ALU::sub(rb, ra);
  } else if (op_code == abi::instruction::RSUBC) {
    std::tie(result, carry, overflow) = ALU::subc(
        rb, ra, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else if (op_code == abi::instruction::SUB) {
    std::tie(result, carry, overflow) = ALU::sub(ra, rb);
  } else if (op_code == abi::instruction::SUBC) {
    std::tie(result, carry, overflow) = ALU::subc(
        ra, rb, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_sub_nz_cc(instruction, ra, rb, result, carry, overflow);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, false);
}

void Logic::execute_s_rri(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rri_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::S_RRI);

  abi::instruction::OpCode op_code = instruction->op_code();
  if (abi::instruction::Instruction::add_rri_op_codes().count(op_code)) {
    execute_add_s_rri(instruction);
  } else if (abi::instruction::Instruction::asr_rri_op_codes().count(op_code)) {
    execute_asr_s_rri(instruction);
  } else {
    throw std::invalid_argument("");
  }
}

void Logic::execute_add_s_rri(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::add_rri_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::S_RRI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::ADD) {
    std::tie(result, carry, overflow) = ALU::add(ra, imm);
  } else if (op_code == abi::instruction::ADDC) {
    std::tie(result, carry, overflow) = ALU::addc(
        ra, imm, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else if (op_code == abi::instruction::AND) {
    result = ALU::and_(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::OR) {
    result = ALU::or_(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::XOR) {
    result = ALU::xor_(ra, imm);
    carry = false;
    overflow = false;
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();

  auto [even, odd] = ALU::signed_extension(result);
  instruction->thread()->reg_file()->write_pair_reg(instruction->dc(), even,
                                                    odd);

  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, carry);
}

void Logic::execute_asr_s_rri(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::asr_rri_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::S_RRI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::ASR) {
    result = ALU::asr(ra, imm);
  } else if (op_code == abi::instruction::LSL) {
    result = ALU::lsl(ra, imm);
  } else if (op_code == abi::instruction::LSL1) {
    result = ALU::lsl1(ra, imm);
  } else if (op_code == abi::instruction::LSL1X) {
    result = ALU::lsl1x(ra, imm);
  } else if (op_code == abi::instruction::LSLX) {
    result = ALU::lslx(ra, imm);
  } else if (op_code == abi::instruction::LSR) {
    result = ALU::lsr(ra, imm);
  } else if (op_code == abi::instruction::LSR1) {
    result = ALU::lsr1(ra, imm);
  } else if (op_code == abi::instruction::LSR1X) {
    result = ALU::lsr1x(ra, imm);
  } else if (op_code == abi::instruction::LSRX) {
    result = ALU::lsrx(ra, imm);
  } else if (op_code == abi::instruction::ROL) {
    result = ALU::rol(ra, imm);
  } else if (op_code == abi::instruction::ROR) {
    result = ALU::ror(ra, imm);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();

  auto [even, odd] = ALU::signed_extension(result);
  instruction->thread()->reg_file()->write_pair_reg(instruction->dc(), even,
                                                    odd);

  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, false);
}

void Logic::execute_s_rric(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_s_rrici(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rrici_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::S_RRICI);

  abi::instruction::OpCode op_code = instruction->op_code();
  if (abi::instruction::Instruction::add_rrici_op_codes().count(op_code)) {
    execute_add_s_rrici(instruction);
  } else if (abi::instruction::Instruction::and_rrici_op_codes().count(
                 op_code)) {
    execute_and_s_rrici(instruction);
  } else if (abi::instruction::Instruction::asr_rrici_op_codes().count(
                 op_code)) {
    execute_asr_s_rrici(instruction);
  } else if (abi::instruction::Instruction::sub_rrici_op_codes().count(
                 op_code)) {
    execute_sub_s_rrici(instruction);
  } else {
    throw std::invalid_argument("");
  }
}

void Logic::execute_add_s_rrici(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::add_rrici_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::S_RRICI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::ADD) {
    std::tie(result, carry, overflow) = ALU::add(ra, imm);
  } else if (op_code == abi::instruction::ADDC) {
    std::tie(result, carry, overflow) = ALU::addc(
        ra, imm, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_add_nz_cc(instruction, ra, result, carry, overflow);

  auto [even, odd] = ALU::signed_extension(result);
  instruction->thread()->reg_file()->write_pair_reg(instruction->dc(), even,
                                                    odd);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, carry);
}

void Logic::execute_and_s_rrici(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::and_rrici_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::S_RRICI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::AND) {
    result = ALU::and_(ra, imm);
  } else if (op_code == abi::instruction::ANDN) {
    result = ALU::andn(ra, imm);
  } else if (op_code == abi::instruction::NAND) {
    result = ALU::nand(ra, imm);
  } else if (op_code == abi::instruction::NOR) {
    result = ALU::nor(ra, imm);
  } else if (op_code == abi::instruction::NXOR) {
    result = ALU::nxor(ra, imm);
  } else if (op_code == abi::instruction::OR) {
    result = ALU::or_(ra, imm);
  } else if (op_code == abi::instruction::ORN) {
    result = ALU::orn(ra, imm);
  } else if (op_code == abi::instruction::XOR) {
    result = ALU::xor_(ra, imm);
  } else if (op_code == abi::instruction::HASH) {
    result = ALU::hash(ra, imm);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_log_nz_cc(instruction, ra, result);

  auto [even, odd] = ALU::signed_extension(result);
  instruction->thread()->reg_file()->write_pair_reg(instruction->dc(), even,
                                                    odd);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, false);
}

void Logic::execute_asr_s_rrici(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::asr_rrici_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::S_RRICI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::ASR) {
    result = ALU::asr(ra, imm);
  } else if (op_code == abi::instruction::LSL) {
    result = ALU::lsl(ra, imm);
  } else if (op_code == abi::instruction::LSL1) {
    result = ALU::lsl1(ra, imm);
  } else if (op_code == abi::instruction::LSL1X) {
    result = ALU::lsl1x(ra, imm);
  } else if (op_code == abi::instruction::LSLX) {
    result = ALU::lslx(ra, imm);
  } else if (op_code == abi::instruction::LSR) {
    result = ALU::lsr(ra, imm);
  } else if (op_code == abi::instruction::LSR1) {
    result = ALU::lsr1(ra, imm);
  } else if (op_code == abi::instruction::LSR1X) {
    result = ALU::lsr1x(ra, imm);
  } else if (op_code == abi::instruction::LSRX) {
    result = ALU::lsrx(ra, imm);
  } else if (op_code == abi::instruction::ROL) {
    result = ALU::rol(ra, imm);
  } else if (op_code == abi::instruction::ROR) {
    result = ALU::ror(ra, imm);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_imm_shift_nz_cc(instruction, ra, result);

  auto [even, odd] = ALU::signed_extension(result);
  instruction->thread()->reg_file()->write_pair_reg(instruction->dc(), even,
                                                    odd);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, false);
}

void Logic::execute_sub_s_rrici(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::sub_rrici_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::S_RRICI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::SUB) {
    std::tie(result, carry, overflow) = ALU::sub(ra, imm);
  } else if (op_code == abi::instruction::SUBC) {
    std::tie(result, carry, overflow) = ALU::subc(
        ra, imm, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_sub_nz_cc(instruction, ra, imm, result, carry, overflow);

  auto [even, odd] = ALU::signed_extension(result);
  instruction->thread()->reg_file()->write_pair_reg(instruction->dc(), even,
                                                    odd);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, carry);
}

void Logic::execute_s_rrif(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rrif_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::S_RRIF);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::ADD) {
    std::tie(result, carry, overflow) = ALU::add(ra, imm);
  } else if (op_code == abi::instruction::ADDC) {
    std::tie(result, carry, overflow) = ALU::addc(
        ra, imm, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else if (op_code == abi::instruction::AND) {
    result = ALU::and_(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::ANDN) {
    result = ALU::andn(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::NAND) {
    result = ALU::nand(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::NOR) {
    result = ALU::nor(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::NXOR) {
    result = ALU::nxor(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::OR) {
    result = ALU::or_(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::ORN) {
    result = ALU::orn(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::SUB) {
    std::tie(result, carry, overflow) = ALU::sub(ra, imm);
  } else if (op_code == abi::instruction::SUBC) {
    std::tie(result, carry, overflow) = ALU::subc(
        ra, imm, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else if (op_code == abi::instruction::XOR) {
    result = ALU::xor_(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::HASH) {
    result = ALU::hash(ra, imm);
    carry = false;
    overflow = false;
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();

  auto [even, odd] = ALU::signed_extension(result);
  instruction->thread()->reg_file()->write_pair_reg(instruction->dc(), even,
                                                    odd);

  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, carry);
}

void Logic::execute_s_rrr(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_s_rrrc(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_s_rrrci(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_u_rri(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rri_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::U_RRI);

  abi::instruction::OpCode op_code = instruction->op_code();
  if (abi::instruction::Instruction::add_rri_op_codes().count(op_code)) {
    execute_add_u_rri(instruction);
  } else if (abi::instruction::Instruction::asr_rri_op_codes().count(op_code)) {
    execute_asr_u_rri(instruction);
  } else {
    throw std::invalid_argument("");
  }
}

void Logic::execute_add_u_rri(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::add_rri_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::U_RRI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::ADD) {
    std::tie(result, carry, overflow) = ALU::add(ra, imm);
  } else if (op_code == abi::instruction::ADDC) {
    std::tie(result, carry, overflow) = ALU::addc(
        ra, imm, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else if (op_code == abi::instruction::AND) {
    result = ALU::and_(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::OR) {
    result = ALU::or_(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::XOR) {
    result = ALU::xor_(ra, imm);
    carry = false;
    overflow = false;
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();

  auto [even, odd] = ALU::unsigned_extension(result);
  instruction->thread()->reg_file()->write_pair_reg(instruction->dc(), even,
                                                    odd);

  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, carry);
}

void Logic::execute_asr_u_rri(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::asr_rri_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::U_RRI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::ASR) {
    result = ALU::asr(ra, imm);
  } else if (op_code == abi::instruction::LSL) {
    result = ALU::lsl(ra, imm);
  } else if (op_code == abi::instruction::LSL1) {
    result = ALU::lsl1(ra, imm);
  } else if (op_code == abi::instruction::LSL1X) {
    result = ALU::lsl1x(ra, imm);
  } else if (op_code == abi::instruction::LSLX) {
    result = ALU::lslx(ra, imm);
  } else if (op_code == abi::instruction::LSR) {
    result = ALU::lsr(ra, imm);
  } else if (op_code == abi::instruction::LSR1) {
    result = ALU::lsr1(ra, imm);
  } else if (op_code == abi::instruction::LSR1X) {
    result = ALU::lsr1x(ra, imm);
  } else if (op_code == abi::instruction::LSRX) {
    result = ALU::lsrx(ra, imm);
  } else if (op_code == abi::instruction::ROL) {
    result = ALU::rol(ra, imm);
  } else if (op_code == abi::instruction::ROR) {
    result = ALU::ror(ra, imm);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();

  auto [even, odd] = ALU::unsigned_extension(result);
  instruction->thread()->reg_file()->write_pair_reg(instruction->dc(), even,
                                                    odd);

  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, false);
}

void Logic::execute_u_rric(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_u_rrici(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rrici_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::U_RRICI);

  abi::instruction::OpCode op_code = instruction->op_code();
  if (abi::instruction::Instruction::add_rrici_op_codes().count(op_code)) {
    execute_add_u_rrici(instruction);
  } else if (abi::instruction::Instruction::and_rrici_op_codes().count(
                 op_code)) {
    execute_and_u_rrici(instruction);
  } else if (abi::instruction::Instruction::asr_rrici_op_codes().count(
                 op_code)) {
    execute_asr_u_rrici(instruction);
  } else if (abi::instruction::Instruction::sub_rrici_op_codes().count(
                 op_code)) {
    execute_sub_u_rrici(instruction);
  } else {
    throw std::invalid_argument("");
  }
}

void Logic::execute_add_u_rrici(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::add_rrici_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::U_RRICI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::ADD) {
    std::tie(result, carry, overflow) = ALU::add(ra, imm);
  } else if (op_code == abi::instruction::ADDC) {
    std::tie(result, carry, overflow) = ALU::addc(
        ra, imm, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_add_nz_cc(instruction, ra, result, carry, overflow);

  auto [even, odd] = ALU::unsigned_extension(result);
  instruction->thread()->reg_file()->write_pair_reg(instruction->dc(), even,
                                                    odd);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, carry);
}

void Logic::execute_and_u_rrici(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::and_rrici_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::U_RRICI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::AND) {
    result = ALU::and_(ra, imm);
  } else if (op_code == abi::instruction::ANDN) {
    result = ALU::andn(ra, imm);
  } else if (op_code == abi::instruction::NAND) {
    result = ALU::nand(ra, imm);
  } else if (op_code == abi::instruction::NOR) {
    result = ALU::nor(ra, imm);
  } else if (op_code == abi::instruction::NXOR) {
    result = ALU::nxor(ra, imm);
  } else if (op_code == abi::instruction::OR) {
    result = ALU::or_(ra, imm);
  } else if (op_code == abi::instruction::ORN) {
    result = ALU::orn(ra, imm);
  } else if (op_code == abi::instruction::XOR) {
    result = ALU::xor_(ra, imm);
  } else if (op_code == abi::instruction::HASH) {
    result = ALU::hash(ra, imm);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_log_nz_cc(instruction, ra, result);

  auto [even, odd] = ALU::unsigned_extension(result);
  instruction->thread()->reg_file()->write_pair_reg(instruction->dc(), even,
                                                    odd);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, false);
}

void Logic::execute_asr_u_rrici(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::asr_rrici_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::U_RRICI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::ASR) {
    result = ALU::asr(ra, imm);
  } else if (op_code == abi::instruction::LSL) {
    result = ALU::lsl(ra, imm);
  } else if (op_code == abi::instruction::LSL1) {
    result = ALU::lsl1(ra, imm);
  } else if (op_code == abi::instruction::LSL1X) {
    result = ALU::lsl1x(ra, imm);
  } else if (op_code == abi::instruction::LSLX) {
    result = ALU::lslx(ra, imm);
  } else if (op_code == abi::instruction::LSR) {
    result = ALU::lsr(ra, imm);
  } else if (op_code == abi::instruction::LSR1) {
    result = ALU::lsr1(ra, imm);
  } else if (op_code == abi::instruction::LSR1X) {
    result = ALU::lsr1x(ra, imm);
  } else if (op_code == abi::instruction::LSRX) {
    result = ALU::lsrx(ra, imm);
  } else if (op_code == abi::instruction::ROL) {
    result = ALU::rol(ra, imm);
  } else if (op_code == abi::instruction::ROR) {
    result = ALU::ror(ra, imm);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_imm_shift_nz_cc(instruction, ra, result);

  auto [even, odd] = ALU::unsigned_extension(result);
  instruction->thread()->reg_file()->write_pair_reg(instruction->dc(), even,
                                                    odd);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, false);
}

void Logic::execute_sub_u_rrici(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::sub_rrici_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::U_RRICI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::SUB) {
    std::tie(result, carry, overflow) = ALU::sub(ra, imm);
  } else if (op_code == abi::instruction::SUBC) {
    std::tie(result, carry, overflow) = ALU::subc(
        ra, imm, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_sub_nz_cc(instruction, ra, imm, result, carry, overflow);

  auto [even, odd] = ALU::unsigned_extension(result);
  instruction->thread()->reg_file()->write_pair_reg(instruction->dc(), even,
                                                    odd);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, carry);
}

void Logic::execute_u_rrif(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rrif_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::U_RRIF);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::ADD) {
    std::tie(result, carry, overflow) = ALU::add(ra, imm);
  } else if (op_code == abi::instruction::ADDC) {
    std::tie(result, carry, overflow) = ALU::addc(
        ra, imm, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else if (op_code == abi::instruction::AND) {
    result = ALU::and_(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::ANDN) {
    result = ALU::andn(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::NAND) {
    result = ALU::nand(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::NOR) {
    result = ALU::nor(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::NXOR) {
    result = ALU::nxor(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::OR) {
    result = ALU::or_(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::ORN) {
    result = ALU::orn(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::SUB) {
    std::tie(result, carry, overflow) = ALU::sub(ra, imm);
  } else if (op_code == abi::instruction::SUBC) {
    std::tie(result, carry, overflow) = ALU::subc(
        ra, imm, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else if (op_code == abi::instruction::XOR) {
    result = ALU::xor_(ra, imm);
    carry = false;
    overflow = false;
  } else if (op_code == abi::instruction::HASH) {
    result = ALU::hash(ra, imm);
    carry = false;
    overflow = false;
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();

  auto [even, odd] = ALU::unsigned_extension(result);
  instruction->thread()->reg_file()->write_pair_reg(instruction->dc(), even,
                                                    odd);

  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, carry);
}

void Logic::execute_u_rrr(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_u_rrrc(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_u_rrrci(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_rr(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rr_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RR);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);

  int64_t result;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::CAO) {
    result = ALU::cao(ra);
  } else if (op_code == abi::instruction::CLO) {
    result = ALU::clo(ra);
  } else if (op_code == abi::instruction::CLS) {
    result = ALU::cls(ra);
  } else if (op_code == abi::instruction::CLZ) {
    result = ALU::clz(ra);
  } else if (op_code == abi::instruction::EXTSB) {
    result = ALU::extsb(ra);
  } else if (op_code == abi::instruction::EXTSH) {
    result = ALU::extsh(ra);
  } else if (op_code == abi::instruction::EXTUB) {
    result = ALU::extub(ra);
  } else if (op_code == abi::instruction::EXTUH) {
    result = ALU::extuh(ra);
  } else if (op_code == abi::instruction::SATS) {
    result = ALU::sats(ra);
  } else if (op_code == abi::instruction::TIME_CFG) {
    throw std::bad_function_call();
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), result);
  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, false);
}

void Logic::execute_rrc(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rrc_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RRC);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);

  int64_t result;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::CAO) {
    result = ALU::cao(ra);
  } else if (op_code == abi::instruction::CLO) {
    result = ALU::clo(ra);
  } else if (op_code == abi::instruction::CLS) {
    result = ALU::cls(ra);
  } else if (op_code == abi::instruction::CLZ) {
    result = ALU::clz(ra);
  } else if (op_code == abi::instruction::EXTSB) {
    result = ALU::extsb(ra);
  } else if (op_code == abi::instruction::EXTSH) {
    result = ALU::extsh(ra);
  } else if (op_code == abi::instruction::EXTUB) {
    result = ALU::extub(ra);
  } else if (op_code == abi::instruction::EXTUH) {
    result = ALU::extuh(ra);
  } else if (op_code == abi::instruction::SATS) {
    result = ALU::sats(ra);
  } else if (op_code == abi::instruction::TIME_CFG) {
    throw std::bad_function_call();
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_log_set_cc(instruction, result);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), 1);
  } else {
    instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), 0);
  }

  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, false);
}

void Logic::execute_rrci(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rrci_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RRCI);

  abi::instruction::OpCode op_code = instruction->op_code();
  if (abi::instruction::Instruction::cao_rrci_op_codes().count(op_code)) {
    execute_cao_rrci(instruction);
  } else if (abi::instruction::Instruction::extsb_rrci_op_codes().count(
                 op_code)) {
    execute_extsb_rrci(instruction);
  } else if (abi::instruction::Instruction::time_cfg_rrci_op_codes().count(
                 op_code)) {
    execute_time_cfg_rrci(instruction);
  } else {
    throw std::invalid_argument("");
  }
}

void Logic::execute_cao_rrci(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::cao_rrci_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RRCI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);

  int64_t result;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::CAO) {
    result = ALU::cao(ra);
  } else if (op_code == abi::instruction::CLO) {
    result = ALU::clo(ra);
  } else if (op_code == abi::instruction::CLS) {
    result = ALU::cls(ra);
  } else if (op_code == abi::instruction::CLZ) {
    result = ALU::clz(ra);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_count_nz_cc(instruction, ra, result);

  instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), result);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, false);
}

void Logic::execute_extsb_rrci(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::extsb_rrci_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RRCI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);

  int64_t result;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::EXTSB) {
    result = ALU::extsb(ra);
  } else if (op_code == abi::instruction::EXTSH) {
    result = ALU::extsh(ra);
  } else if (op_code == abi::instruction::EXTUB) {
    result = ALU::extub(ra);
  } else if (op_code == abi::instruction::EXTUH) {
    result = ALU::extuh(ra);
  } else if (op_code == abi::instruction::SATS) {
    result = ALU::sats(ra);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_log_nz_cc(instruction, ra, result);

  instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), result);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, false);
}

void Logic::execute_time_cfg_rrci(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_zr(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rr_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZR);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);

  int64_t result;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::CAO) {
    result = ALU::cao(ra);
  } else if (op_code == abi::instruction::CLO) {
    result = ALU::clo(ra);
  } else if (op_code == abi::instruction::CLS) {
    result = ALU::cls(ra);
  } else if (op_code == abi::instruction::CLZ) {
    result = ALU::clz(ra);
  } else if (op_code == abi::instruction::EXTSB) {
    result = ALU::extsb(ra);
  } else if (op_code == abi::instruction::EXTSH) {
    result = ALU::extsh(ra);
  } else if (op_code == abi::instruction::EXTUB) {
    result = ALU::extub(ra);
  } else if (op_code == abi::instruction::EXTUH) {
    result = ALU::extuh(ra);
  } else if (op_code == abi::instruction::SATS) {
    result = ALU::sats(ra);
  } else if (op_code == abi::instruction::TIME_CFG) {
    throw std::bad_function_call();
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, false);
}

void Logic::execute_zrc(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rrc_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZRC);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);

  int64_t result;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::CAO) {
    result = ALU::cao(ra);
  } else if (op_code == abi::instruction::CLO) {
    result = ALU::clo(ra);
  } else if (op_code == abi::instruction::CLS) {
    result = ALU::cls(ra);
  } else if (op_code == abi::instruction::CLZ) {
    result = ALU::clz(ra);
  } else if (op_code == abi::instruction::EXTSB) {
    result = ALU::extsb(ra);
  } else if (op_code == abi::instruction::EXTSH) {
    result = ALU::extsh(ra);
  } else if (op_code == abi::instruction::EXTUB) {
    result = ALU::extub(ra);
  } else if (op_code == abi::instruction::EXTUH) {
    result = ALU::extuh(ra);
  } else if (op_code == abi::instruction::SATS) {
    result = ALU::sats(ra);
  } else if (op_code == abi::instruction::TIME_CFG) {
    throw std::bad_function_call();
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_log_set_cc(instruction, result);

  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, false);
}

void Logic::execute_zrci(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rrci_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZRCI);

  abi::instruction::OpCode op_code = instruction->op_code();
  if (abi::instruction::Instruction::cao_rrci_op_codes().count(op_code)) {
    execute_cao_zrci(instruction);
  } else if (abi::instruction::Instruction::extsb_rrci_op_codes().count(
                 op_code)) {
    execute_extsb_zrci(instruction);
  } else if (abi::instruction::Instruction::time_cfg_rrci_op_codes().count(
                 op_code)) {
    execute_time_cfg_zrci(instruction);
  } else {
    throw std::invalid_argument("");
  }
}

void Logic::execute_cao_zrci(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::cao_rrci_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZRCI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);

  int64_t result;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::CAO) {
    result = ALU::cao(ra);
  } else if (op_code == abi::instruction::CLO) {
    result = ALU::clo(ra);
  } else if (op_code == abi::instruction::CLS) {
    result = ALU::cls(ra);
  } else if (op_code == abi::instruction::CLZ) {
    result = ALU::clz(ra);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_count_nz_cc(instruction, ra, result);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, false);
}

void Logic::execute_extsb_zrci(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::extsb_rrci_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZRCI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);

  int64_t result;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::EXTSB) {
    result = ALU::extsb(ra);
  } else if (op_code == abi::instruction::EXTSH) {
    result = ALU::extsh(ra);
  } else if (op_code == abi::instruction::EXTUB) {
    result = ALU::extub(ra);
  } else if (op_code == abi::instruction::EXTUH) {
    result = ALU::extuh(ra);
  } else if (op_code == abi::instruction::SATS) {
    result = ALU::sats(ra);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_log_nz_cc(instruction, ra, result);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, false);
}

void Logic::execute_time_cfg_zrci(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_s_rr(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_s_rrc(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_s_rrci(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_u_rr(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_u_rrc(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_u_rrci(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_drdici(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::drdici_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::DRDICI);

  abi::instruction::OpCode op_code = instruction->op_code();
  if (abi::instruction::Instruction::div_step_drdici_op_codes().count(
          op_code)) {
    execute_div_step_drdici(instruction);
  } else if (abi::instruction::Instruction::mul_step_drdici_op_codes().count(
                 op_code)) {
    execute_mul_step_drdici(instruction);
  } else {
    throw std::invalid_argument("");
  }
}

void Logic::execute_div_step_drdici(
    abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::div_step_drdici_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::DRDICI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t dbe = instruction->thread()->reg_file()->read_gp_reg(
      instruction->db()->even_reg(), abi::word::SIGNED);
  int64_t dbo = instruction->thread()->reg_file()->read_gp_reg(
      instruction->db()->odd_reg(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  auto dbo_data_word = new abi::word::DataWord();
  dbo_data_word->set_value(dbo);

  auto ra_shift_data_word = new abi::word::DataWord();
  ra_shift_data_word->set_value(ALU::lsl(ra, imm));

  auto [result, carry, overflow] = ALU::sub(dbo, ALU::lsl(ra, imm));

  int64_t dce;
  int64_t dco;
  if (dbo_data_word->value(abi::word::UNSIGNED) >=
      ra_shift_data_word->value(abi::word::UNSIGNED)) {
    dce = ALU::lsl1(dbe, 1);
    dco = result;
  } else {
    dce = ALU::lsl(dbe, 1);
    dco = instruction->thread()->reg_file()->read_gp_reg(
        instruction->dc()->odd_reg(), abi::word::SIGNED);
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_div_cc(instruction, ra);

  instruction->thread()->reg_file()->write_gp_reg(instruction->dc()->even_reg(),
                                                  dce);
  instruction->thread()->reg_file()->write_gp_reg(instruction->dc()->odd_reg(),
                                                  dco);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, false);

  delete dbo_data_word;
  delete ra_shift_data_word;
}

void Logic::execute_mul_step_drdici(
    abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::mul_step_drdici_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::DRDICI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t dbe = instruction->thread()->reg_file()->read_gp_reg(
      instruction->db()->even_reg(), abi::word::SIGNED);
  int64_t dbo = instruction->thread()->reg_file()->read_gp_reg(
      instruction->db()->odd_reg(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result1 = ALU::lsr(dbe, 1);
  auto [result2, carry, overflow] = ALU::sub(ALU::and_(dbe, 1), 1);

  int64_t dco;
  if (result2 == 0) {
    std::tie(dco, carry, overflow) = ALU::add(dbo, ALU::lsl(ra, imm));
  } else {
    dco = instruction->thread()->reg_file()->read_gp_reg(
        instruction->dc()->odd_reg(), abi::word::SIGNED);
  }
  int64_t dce = ALU::lsr(dbe, 1);

  instruction->thread()->reg_file()->clear_conditions();
  set_boot_cc(instruction, ra, result1);

  instruction->thread()->reg_file()->write_gp_reg(instruction->dc()->even_reg(),
                                                  dce);
  instruction->thread()->reg_file()->write_gp_reg(instruction->dc()->odd_reg(),
                                                  dco);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result1, false);
}

void Logic::execute_rrri(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rrri_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RRRI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t rb = instruction->thread()->reg_file()->read_src_reg(
      instruction->rb(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::LSL_ADD) {
    std::tie(result, carry, overflow) = ALU::lsl_add(ra, rb, imm);
  } else if (op_code == abi::instruction::LSL_SUB) {
    std::tie(result, carry, overflow) = ALU::lsl_sub(ra, rb, imm);
  } else if (op_code == abi::instruction::LSR_ADD) {
    std::tie(result, carry, overflow) = ALU::lsr_add(ra, rb, imm);
  } else if (op_code == abi::instruction::ROL_ADD) {
    std::tie(result, carry, overflow) = ALU::rol_add(ra, rb, imm);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), result);
  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, carry);
}

void Logic::execute_rrrici(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rrrici_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RRRICI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t rb = instruction->thread()->reg_file()->read_src_reg(
      instruction->rb(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::LSL_ADD) {
    std::tie(result, carry, overflow) = ALU::lsl_add(ra, rb, imm);
  } else if (op_code == abi::instruction::LSL_SUB) {
    std::tie(result, carry, overflow) = ALU::lsl_sub(ra, rb, imm);
  } else if (op_code == abi::instruction::LSR_ADD) {
    std::tie(result, carry, overflow) = ALU::lsr_add(ra, rb, imm);
  } else if (op_code == abi::instruction::ROL_ADD) {
    std::tie(result, carry, overflow) = ALU::rol_add(ra, rb, imm);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_div_nz_cc(instruction, ra);

  instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), result);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, carry);
}

void Logic::execute_zrri(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rrri_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZRRI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t rb = instruction->thread()->reg_file()->read_src_reg(
      instruction->rb(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::LSL_ADD) {
    std::tie(result, carry, overflow) = ALU::lsl_add(ra, rb, imm);
  } else if (op_code == abi::instruction::LSL_SUB) {
    std::tie(result, carry, overflow) = ALU::lsl_sub(ra, rb, imm);
  } else if (op_code == abi::instruction::LSR_ADD) {
    std::tie(result, carry, overflow) = ALU::lsr_add(ra, rb, imm);
  } else if (op_code == abi::instruction::ROL_ADD) {
    std::tie(result, carry, overflow) = ALU::rol_add(ra, rb, imm);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, carry);
}

void Logic::execute_zrrici(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rrrici_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZRRICI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t rb = instruction->thread()->reg_file()->read_src_reg(
      instruction->rb(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::LSL_ADD) {
    std::tie(result, carry, overflow) = ALU::lsl_add(ra, rb, imm);
  } else if (op_code == abi::instruction::LSL_SUB) {
    std::tie(result, carry, overflow) = ALU::lsl_sub(ra, rb, imm);
  } else if (op_code == abi::instruction::LSR_ADD) {
    std::tie(result, carry, overflow) = ALU::lsr_add(ra, rb, imm);
  } else if (op_code == abi::instruction::ROL_ADD) {
    std::tie(result, carry, overflow) = ALU::rol_add(ra, rb, imm);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_div_nz_cc(instruction, ra);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, carry);
}

void Logic::execute_s_rrri(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_s_rrrici(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_u_rrri(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_u_rrrici(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_rir(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rir_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RIR);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::SUB) {
    std::tie(result, carry, overflow) = ALU::sub(imm, ra);
  } else if (op_code == abi::instruction::SUBC) {
    std::tie(result, carry, overflow) = ALU::subc(
        imm, ra, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), result);
  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, carry);
}

void Logic::execute_rirc(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rirc_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RIRC);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::SUB) {
    std::tie(result, carry, overflow) = ALU::sub(imm, ra);
  } else if (op_code == abi::instruction::SUBC) {
    std::tie(result, carry, overflow) = ALU::subc(
        imm, ra, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_sub_set_cc(instruction, ra, imm, result);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), 1);
  } else {
    instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), 0);
  }

  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, carry);
}

void Logic::execute_rirci(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rirci_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::RIRCI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::SUB) {
    std::tie(result, carry, overflow) = ALU::sub(imm, ra);
  } else if (op_code == abi::instruction::SUBC) {
    std::tie(result, carry, overflow) = ALU::subc(
        imm, ra, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_sub_nz_cc(instruction, ra, imm, result, carry, overflow);

  instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), result);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, carry);
}

void Logic::execute_zir(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rir_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZIR);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::SUB) {
    std::tie(result, carry, overflow) = ALU::sub(imm, ra);
  } else if (op_code == abi::instruction::SUBC) {
    std::tie(result, carry, overflow) = ALU::subc(
        imm, ra, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, carry);
}

void Logic::execute_zirc(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rirc_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZIRC);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::SUB) {
    std::tie(result, carry, overflow) = ALU::sub(imm, ra);
  } else if (op_code == abi::instruction::SUBC) {
    std::tie(result, carry, overflow) = ALU::subc(
        imm, ra, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_sub_set_cc(instruction, ra, imm, result);

  instruction->thread()->reg_file()->increment_pc_reg();

  set_flags(instruction, result, carry);
}

void Logic::execute_zirci(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::rirci_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ZIRCI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  int64_t result;
  bool carry;
  bool overflow;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::SUB) {
    std::tie(result, carry, overflow) = ALU::sub(imm, ra);
  } else if (op_code == abi::instruction::SUBC) {
    std::tie(result, carry, overflow) = ALU::subc(
        imm, ra, instruction->thread()->reg_file()->flag(abi::isa::CARRY));
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  set_sub_nz_cc(instruction, ra, imm, result, carry, overflow);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }

  set_flags(instruction, result, carry);
}

void Logic::execute_s_rirc(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_s_rirci(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_u_rirc(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_u_rirci(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_r(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_rci(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_z(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::r_op_codes().count(
             instruction->op_code()) or
         instruction->op_code() == abi::instruction::NOP);
  assert(instruction->suffix() == abi::instruction::Z);

  instruction->thread()->reg_file()->increment_pc_reg();
}

void Logic::execute_zci(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_s_r(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_s_rci(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_u_r(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_u_rci(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_ci(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::ci_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::CI);

  instruction->thread()->reg_file()->clear_conditions();
  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
    scheduler_->sleep(instruction->thread()->id());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
    scheduler_->sleep(instruction->thread()->id());
  }
}

void Logic::execute_i(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_ddci(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::ddci_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::DDCI);

  abi::instruction::OpCode op_code = instruction->op_code();
  if (abi::instruction::Instruction::movd_ddci_op_codes().count(op_code)) {
    execute_movd_ddci(instruction);
  } else if (abi::instruction::Instruction::swapd_ddci_op_codes().count(
                 op_code)) {
    execute_swapd_ddci(instruction);
  } else {
    throw std::invalid_argument("");
  }
}

void Logic::execute_movd_ddci(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::movd_ddci_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::DDCI);

  int64_t dbe = instruction->thread()->reg_file()->read_gp_reg(
      instruction->db()->even_reg(), abi::word::SIGNED);
  int64_t dbo = instruction->thread()->reg_file()->read_gp_reg(
      instruction->db()->odd_reg(), abi::word::SIGNED);

  instruction->thread()->reg_file()->clear_conditions();
  instruction->thread()->reg_file()->write_gp_reg(instruction->dc()->even_reg(),
                                                  dbe);
  instruction->thread()->reg_file()->write_gp_reg(instruction->dc()->odd_reg(),
                                                  dbo);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }
}

void Logic::execute_swapd_ddci(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::swapd_ddci_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::DDCI);

  int64_t dbe = instruction->thread()->reg_file()->read_gp_reg(
      instruction->db()->even_reg(), abi::word::SIGNED);
  int64_t dbo = instruction->thread()->reg_file()->read_gp_reg(
      instruction->db()->odd_reg(), abi::word::SIGNED);

  instruction->thread()->reg_file()->clear_conditions();
  instruction->thread()->reg_file()->write_gp_reg(instruction->dc()->even_reg(),
                                                  dbo);
  instruction->thread()->reg_file()->write_gp_reg(instruction->dc()->odd_reg(),
                                                  dbe);

  if (instruction->thread()->reg_file()->condition(instruction->condition())) {
    instruction->thread()->reg_file()->write_pc_reg(instruction->pc()->value());
  } else {
    instruction->thread()->reg_file()->increment_pc_reg();
  }
}

void Logic::execute_erri(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::erri_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ERRI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t off = instruction->off()->value();

  auto [address, carry, overflow] = ALU::add(ra, off);

  int64_t result;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::LBS) {
    result = operand_collector_->lbs(address);
  } else if (op_code == abi::instruction::LBU) {
    result = operand_collector_->lbu(address);
  } else if (op_code == abi::instruction::LHS) {
    result = operand_collector_->lhs(address);
  } else if (op_code == abi::instruction::LHU) {
    result = operand_collector_->lhu(address);
  } else if (op_code == abi::instruction::LW) {
    result = operand_collector_->lw(address);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  instruction->thread()->reg_file()->write_gp_reg(instruction->rc(), result);
  instruction->thread()->reg_file()->increment_pc_reg();
}

void Logic::execute_s_erri(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_u_erri(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_edri(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::edri_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::EDRI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t off = instruction->off()->value();

  auto [address, carry, overflow] = ALU::add(ra, off);

  int64_t even;
  int64_t odd;
  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::LD) {
    std::tie(even, odd) = operand_collector_->ld(address);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  instruction->thread()->reg_file()->write_pair_reg(instruction->dc(), even,
                                                    odd);
  instruction->thread()->reg_file()->increment_pc_reg();
}

void Logic::execute_erii(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::erii_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ERII);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t off = instruction->off()->value();
  int64_t imm = instruction->imm()->value();

  auto [address, carry, overflow] = ALU::add(ra, off);

  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::SB) {
    operand_collector_->sb(address, imm);
  } else if (op_code == abi::instruction::SB_ID) {
    operand_collector_->sb(address, ALU::or_(instruction->thread()->id(), imm));
  } else if (op_code == abi::instruction::SH) {
    operand_collector_->sh(address, imm);
  } else if (op_code == abi::instruction::SH_ID) {
    operand_collector_->sh(address, ALU::or_(instruction->thread()->id(), imm));
  } else if (op_code == abi::instruction::SW) {
    operand_collector_->sw(address, imm);
  } else if (op_code == abi::instruction::SW_ID) {
    operand_collector_->sw(address, ALU::or_(instruction->thread()->id(), imm));
  } else if (op_code == abi::instruction::SD) {
    auto [even, odd] = ALU::unsigned_extension(imm);
    operand_collector_->sd(address, even, odd);
  } else if (op_code == abi::instruction::SD_ID) {
    auto [even, odd] =
        ALU::unsigned_extension(ALU::or_(instruction->thread()->id(), imm));
    operand_collector_->sd(address, even, odd);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  instruction->thread()->reg_file()->increment_pc_reg();
}

void Logic::execute_erir(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::erir_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ERIR);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t off = instruction->off()->value();
  int64_t rb = instruction->thread()->reg_file()->read_src_reg(
      instruction->rb(), abi::word::SIGNED);
  auto rb_data_word = new abi::word::DataWord();
  rb_data_word->set_value(rb);

  auto [address, carry, overflow] = ALU::add(ra, off);

  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::SB) {
    operand_collector_->sb(address,
                           rb_data_word->bit_slice(abi::word::UNSIGNED, 0, 8));
  } else if (op_code == abi::instruction::SH) {
    operand_collector_->sh(address,
                           rb_data_word->bit_slice(abi::word::UNSIGNED, 0, 16));
  } else if (op_code == abi::instruction::SW) {
    operand_collector_->sw(address,
                           rb_data_word->bit_slice(abi::word::UNSIGNED, 0, 32));
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  instruction->thread()->reg_file()->increment_pc_reg();

  delete rb_data_word;
}

void Logic::execute_erid(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::erid_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::ERID);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t off = instruction->off()->value();
  auto [even, odd] = instruction->thread()->reg_file()->read_pair_reg(
      instruction->db(), abi::word::SIGNED);

  auto [address, carry, overflow] = ALU::add(ra, off);

  abi::instruction::OpCode op_code = instruction->op_code();
  if (op_code == abi::instruction::SD) {
    operand_collector_->sd(address, even, odd);
  } else {
    throw std::invalid_argument("");
  }

  instruction->thread()->reg_file()->clear_conditions();
  instruction->thread()->reg_file()->increment_pc_reg();
}

void Logic::execute_dma_rri(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::dma_rri_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::DMA_RRI);

  abi::instruction::OpCode op_code = instruction->op_code();
  if (abi::instruction::Instruction::ldma_dma_rri_op_codes().count(op_code)) {
    execute_ldma(instruction);
  } else if (abi::instruction::Instruction::ldmai_dma_rri_op_codes().count(
                 op_code)) {
    execute_ldmai(instruction);
  } else if (abi::instruction::Instruction::sdma_dma_rri_op_codes().count(
                 op_code)) {
    execute_sdma(instruction);
  } else {
    throw std::invalid_argument("");
  }
}

void Logic::execute_ldma(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::ldma_dma_rri_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::DMA_RRI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t rb = instruction->thread()->reg_file()->read_src_reg(
      instruction->rb(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  Address wram_end_address =
      util::ConfigLoader::wram_offset() + util::ConfigLoader::wram_size();
  auto wram_end_address_width =
      static_cast<Address>(floor(log2(wram_end_address)) + 1);
  auto wram_mask = static_cast<Address>(pow(2, wram_end_address_width) - 1);
  Address wram_address = ALU::and_(ra, wram_mask);

  Address mram_end_address =
      util::ConfigLoader::mram_offset() + util::ConfigLoader::mram_size();
  auto mram_end_address_width =
      static_cast<Address>(floor(log2(mram_end_address)) + 1);
  auto mram_mask = static_cast<Address>(pow(2, mram_end_address_width) - 1);
  Address mram_address = ALU::and_(rb, mram_mask);

  Address min_access_granularity = util::ConfigLoader::min_access_granularity();
  Address size = (1 + ALU::and_(imm + ALU::and_(ALU::lsr(ra, 24), 255), 255)) *
                 min_access_granularity;

  dma_->transfer_from_mram_to_wram(wram_address, mram_address, size,
                                   instruction);

  stat_factory_->overwrite("mram_address", mram_address); 
  stat_factory_->overwrite("mram_access_thread",
                           instruction->thread()->id());
  stat_factory_->overwrite("mram_access_size", size); 

  instruction->thread()->reg_file()->clear_conditions();
}

void Logic::execute_ldmai(abi::instruction::Instruction *instruction) {
  throw std::bad_function_call();
}

void Logic::execute_sdma(abi::instruction::Instruction *instruction) {
  assert(abi::instruction::Instruction::sdma_dma_rri_op_codes().count(
      instruction->op_code()));
  assert(instruction->suffix() == abi::instruction::DMA_RRI);

  int64_t ra = instruction->thread()->reg_file()->read_src_reg(
      instruction->ra(), abi::word::SIGNED);
  int64_t rb = instruction->thread()->reg_file()->read_src_reg(
      instruction->rb(), abi::word::SIGNED);
  int64_t imm = instruction->imm()->value();

  Address wram_end_address =
      util::ConfigLoader::wram_offset() + util::ConfigLoader::wram_size();
  auto wram_end_address_width =
      static_cast<Address>(floor(log2(wram_end_address)) + 1);
  auto wram_mask = static_cast<Address>(pow(2, wram_end_address_width) - 1);
  Address wram_address = ALU::and_(ra, wram_mask);

  Address mram_end_address =
      util::ConfigLoader::mram_offset() + util::ConfigLoader::mram_size();
  auto mram_end_address_width =
      static_cast<Address>(floor(log2(mram_end_address)) + 1);
  auto mram_mask = static_cast<Address>(pow(2, mram_end_address_width) - 1);
  Address mram_address = ALU::and_(rb, mram_mask);

  Address min_access_granularity = util::ConfigLoader::min_access_granularity();
  Address size = (1 + ALU::and_(imm + ALU::and_(ALU::lsr(ra, 24), 255), 255)) *
                 min_access_granularity;

  dma_->transfer_from_wram_to_mram(wram_address, mram_address, size,
                                   instruction);

  stat_factory_->overwrite("mram_address", mram_address);
  stat_factory_->overwrite("mram_access_thread",
                           instruction->thread()->id());
  stat_factory_->overwrite("mram_access_size", size); 

  instruction->thread()->reg_file()->clear_conditions();
}

void Logic::set_acquire_cc(abi::instruction::Instruction *instruction,
                           int64_t result) {
  if (result == 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::Z);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::NZ);
  }
}

void Logic::set_add_nz_cc(abi::instruction::Instruction *instruction,
                          int64_t operand1, int64_t result, bool carry,
                          bool overflow) {
  if (result == 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::Z);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::NZ);
  }

  if (carry) {
    instruction->thread()->reg_file()->set_condition(abi::isa::C);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::NC);
  }

  if (result == 0 and instruction->thread()->reg_file()->flag(abi::isa::ZERO)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::XZ);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::XNZ);
  }

  if (overflow) {
    instruction->thread()->reg_file()->set_condition(abi::isa::OV);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::NOV);
  }

  if (result >= 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::PL);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::MI);
  }

  if (operand1 == 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::SZ);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::SNZ);
  }

  if (operand1 >= 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::SPL);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::SMI);
  }

  auto result_data_word = new abi::word::DataWord();
  result_data_word->set_value(result);

  if (result_data_word->bit(6)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::NC5);
  }
  if (result_data_word->bit(7)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::NC6);
  }
  if (result_data_word->bit(8)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::NC7);
  }
  if (result_data_word->bit(9)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::NC8);
  }
  if (result_data_word->bit(10)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::NC9);
  }
  if (result_data_word->bit(11)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::NC10);
  }
  if (result_data_word->bit(12)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::NC11);
  }
  if (result_data_word->bit(13)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::NC12);
  }
  if (result_data_word->bit(14)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::NC13);
  }
  if (result_data_word->bit(15)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::NC14);
  }

  delete result_data_word;
}

void Logic::set_boot_cc(abi::instruction::Instruction *instruction,
                        int64_t operand1, int64_t result) {
  if (result == 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::Z);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::NZ);
  }

  if (result == 0 and instruction->thread()->reg_file()->flag(abi::isa::ZERO)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::XZ);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::XNZ);
  }

  if (operand1 == 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::SZ);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::SNZ);
  }

  if (operand1 >= 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::SPL);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::SMI);
  }
}

void Logic::set_count_nz_cc(abi::instruction::Instruction *instruction,
                            int64_t operand1, int64_t result) {
  if (result == 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::Z);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::NZ);
  }

  if (result == 0 and instruction->thread()->reg_file()->flag(abi::isa::ZERO)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::XZ);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::XNZ);
  }

  if (operand1 == 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::SZ);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::SNZ);
  }

  if (operand1 >= 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::SPL);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::SMI);
  }

  if (result == abi::word::DataWord().width()) {
    instruction->thread()->reg_file()->set_condition(abi::isa::MAX);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::NMAX);
  }
}

void Logic::set_div_cc(abi::instruction::Instruction *instruction,
                       int64_t operand1) {
  if (operand1 == 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::SZ);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::SNZ);
  }

  if (operand1 >= 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::SPL);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::SMI);
  }
}

void Logic::set_div_nz_cc(abi::instruction::Instruction *instruction,
                          int64_t operand1) {
  if (operand1 == 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::SZ);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::SNZ);
  }

  if (operand1 >= 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::SPL);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::SMI);
  }
}

void Logic::set_ext_sub_set_cc(abi::instruction::Instruction *instruction,
                               int64_t operand1, int64_t operand2,
                               int64_t result, bool carry, bool overflow) {
  if (result == 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::Z);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::NZ);
  }

  if (carry) {
    instruction->thread()->reg_file()->set_condition(abi::isa::C);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::NC);
  }

  if (result == 0 and instruction->thread()->reg_file()->flag(abi::isa::ZERO)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::XZ);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::XNZ);
  }

  if (overflow) {
    instruction->thread()->reg_file()->set_condition(abi::isa::OV);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::NOV);
  }

  if (result >= 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::PL);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::MI);
  }

  if (operand1 == operand2) {
    instruction->thread()->reg_file()->set_condition(abi::isa::EQ);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::NEQ);
  }

  if (operand1 == 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::SZ);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::SNZ);
  }

  if (operand1 >= 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::SPL);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::SMI);
  }

  auto data_word1 = new abi::word::DataWord();
  data_word1->set_value(operand1);

  auto data_word2 = new abi::word::DataWord();
  data_word2->set_value(operand2);

  if (data_word1->value(abi::word::UNSIGNED) <
      data_word2->value(abi::word::UNSIGNED)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::LTU);
  }

  if (data_word1->value(abi::word::UNSIGNED) <=
      data_word2->value(abi::word::UNSIGNED)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::LEU);
  }

  if (data_word1->value(abi::word::UNSIGNED) >
      data_word2->value(abi::word::UNSIGNED)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::GTU);
  }

  if (data_word1->value(abi::word::UNSIGNED) >=
      data_word2->value(abi::word::UNSIGNED)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::GEU);
  }

  if (data_word1->value(abi::word::SIGNED) <
      data_word2->value(abi::word::SIGNED)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::LTS);
  }

  if (data_word1->value(abi::word::SIGNED) <=
      data_word2->value(abi::word::SIGNED)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::LES);
  }

  if (data_word1->value(abi::word::SIGNED) >
      data_word2->value(abi::word::SIGNED)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::GTS);
  }

  if (data_word1->value(abi::word::SIGNED) >=
      data_word2->value(abi::word::SIGNED)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::GES);
  }

  if (carry or instruction->thread()->reg_file()->flag(abi::isa::ZERO)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::XLEU);
  }

  if (carry and not instruction->thread()->reg_file()->flag(abi::isa::ZERO)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::XGTU);
  }

  if (instruction->thread()->reg_file()->flag(abi::isa::ZERO) and
      (result < 0 or overflow)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::XLES);
  }

  if (not instruction->thread()->reg_file()->flag(abi::isa::ZERO) and
      (result >= 0 or overflow)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::XGTS);
  }

  delete data_word1;
  delete data_word2;
}

void Logic::set_imm_shift_nz_cc(abi::instruction::Instruction *instruction,
                                int64_t operand1, int64_t result) {
  if (result == 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::Z);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::NZ);
  }

  if (result == 0 and instruction->thread()->reg_file()->flag(abi::isa::ZERO)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::XZ);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::XNZ);
  }

  if (result % 2 == 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::E);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::O);
  }

  if (result >= 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::PL);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::MI);
  }

  if (operand1 == 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::SZ);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::SNZ);
  }

  if (operand1 % 2 == 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::SE);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::SO);
  }

  if (operand1 >= 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::SPL);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::SMI);
  }
}

void Logic::set_log_nz_cc(abi::instruction::Instruction *instruction,
                          int64_t operand1, int64_t result) {
  if (result == 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::Z);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::NZ);
  }

  if (result == 0 and instruction->thread()->reg_file()->flag(abi::isa::ZERO)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::XZ);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::XNZ);
  }

  if (result >= 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::PL);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::MI);
  }

  if (operand1 == 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::SZ);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::SNZ);
  }

  if (operand1 >= 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::SPL);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::SMI);
  }
}

void Logic::set_log_set_cc(abi::instruction::Instruction *instruction,
                           int64_t result) {
  if (result == 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::Z);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::NZ);
  }

  if (result == 0 and instruction->thread()->reg_file()->flag(abi::isa::ZERO)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::XZ);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::XNZ);
  }
}

void Logic::set_mul_nz_cc(abi::instruction::Instruction *instruction,
                          int64_t operand1, int64_t result) {
  if (result == 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::Z);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::NZ);
  }

  if (result == 0 and instruction->thread()->reg_file()->flag(abi::isa::ZERO)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::XZ);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::XNZ);
  }

  if (operand1 == 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::SZ);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::SNZ);
  }

  if (operand1 >= 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::SPL);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::SMI);
  }

  if (result < 256) {
    instruction->thread()->reg_file()->set_condition(abi::isa::SMALL);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::LARGE);
  }
}

void Logic::set_sub_nz_cc(abi::instruction::Instruction *instruction,
                          int64_t operand1, int64_t operand2, int64_t result,
                          bool carry, bool overflow) {
  if (result == 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::Z);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::NZ);
  }

  if (carry) {
    instruction->thread()->reg_file()->set_condition(abi::isa::C);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::NC);
  }

  if (result == 0 and instruction->thread()->reg_file()->flag(abi::isa::ZERO)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::XZ);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::XNZ);
  }

  if (overflow) {
    instruction->thread()->reg_file()->set_condition(abi::isa::OV);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::NOV);
  }

  if (result >= 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::PL);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::MI);
  }

  if (operand1 == operand2) {
    instruction->thread()->reg_file()->set_condition(abi::isa::EQ);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::NEQ);
  }

  if (operand1 >= 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::SPL);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::SMI);
  }

  auto data_word1 = new abi::word::DataWord();
  data_word1->set_value(operand1);

  auto data_word2 = new abi::word::DataWord();
  data_word2->set_value(operand2);

  if (data_word1->value(abi::word::UNSIGNED) <
      data_word2->value(abi::word::UNSIGNED)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::LTU);
  }

  if (data_word1->value(abi::word::UNSIGNED) <=
      data_word2->value(abi::word::UNSIGNED)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::LEU);
  }

  if (data_word1->value(abi::word::UNSIGNED) >
      data_word2->value(abi::word::UNSIGNED)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::GTU);
  }

  if (data_word1->value(abi::word::UNSIGNED) >=
      data_word2->value(abi::word::UNSIGNED)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::GEU);
  }

  if (data_word1->value(abi::word::SIGNED) <
      data_word2->value(abi::word::SIGNED)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::LTS);
  }

  if (data_word1->value(abi::word::SIGNED) <=
      data_word2->value(abi::word::SIGNED)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::LES);
  }

  if (data_word1->value(abi::word::SIGNED) >
      data_word2->value(abi::word::SIGNED)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::GTS);
  }

  if (data_word1->value(abi::word::SIGNED) >=
      data_word2->value(abi::word::SIGNED)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::GES);
  }

  if (carry or instruction->thread()->reg_file()->flag(abi::isa::ZERO)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::XLEU);
  }

  if (carry and not instruction->thread()->reg_file()->flag(abi::isa::ZERO)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::XGTU);
  }

  if (instruction->thread()->reg_file()->flag(abi::isa::ZERO) and
      (result < 0 or overflow)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::XLES);
  }

  if (not instruction->thread()->reg_file()->flag(abi::isa::ZERO) and
      (result >= 0 or overflow)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::XGTS);
  }

  delete data_word1;
  delete data_word2;
}

void Logic::set_sub_set_cc(abi::instruction::Instruction *instruction,
                           int64_t operand1, int64_t operand2, int64_t result) {
  if (result == 0) {
    instruction->thread()->reg_file()->set_condition(abi::isa::Z);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::NZ);
  }

  if (result == 0 and instruction->thread()->reg_file()->flag(abi::isa::ZERO)) {
    instruction->thread()->reg_file()->set_condition(abi::isa::XZ);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::XNZ);
  }

  if (operand1 == operand2) {
    instruction->thread()->reg_file()->set_condition(abi::isa::EQ);
  } else {
    instruction->thread()->reg_file()->set_condition(abi::isa::NEQ);
  }
}

void Logic::set_flags(abi::instruction::Instruction *instruction,
                      int64_t result, bool carry) {
  if (result == 0) {
    instruction->thread()->reg_file()->set_flag(abi::isa::ZERO);
  } else {
    instruction->thread()->reg_file()->clear_flag(abi::isa::ZERO);
  }

  if (carry) {
    instruction->thread()->reg_file()->set_flag(abi::isa::CARRY);
  } else {
    instruction->thread()->reg_file()->clear_flag(abi::isa::CARRY);
  }
}

}  // namespace upmem_sim::simulator::dpu
