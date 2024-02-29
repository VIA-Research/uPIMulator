#include "simulator/dpu/cycle_rule.h"

namespace upmem_sim::simulator::dpu {

CycleRule::CycleRule(util::ArgumentParser *argument_parser)
    : input_q_(new basic::Queue<abi::instruction::Instruction>(1)),
      wait_q_(new basic::TimerQueue<abi::instruction::Instruction>(1)),
      ready_q_(new basic::Queue<abi::instruction::Instruction>(1)),
      stat_factory_(new util::StatFactory("CycleRule")) {
  int num_tasklets =
      static_cast<int>(argument_parser->get_int_parameter("num_tasklets"));
  prev_write_gp_regs_.resize(num_tasklets);
  cur_read_gp_regs_.resize(num_tasklets);
}

CycleRule::~CycleRule() {
  delete input_q_;
  delete wait_q_;
  delete ready_q_;

  for (auto &prev_write_gp_regs : prev_write_gp_regs_) {
    for (auto &gp_reg : prev_write_gp_regs) {
      delete gp_reg;
    }
  }

  for (auto &cur_read_gp_regs : cur_read_gp_regs_) {
    for (auto &gp_reg : cur_read_gp_regs) {
      delete gp_reg;
    }
  }

  delete stat_factory_;
}

util::StatFactory *CycleRule::stat_factory() {
  auto stat_factory = new util::StatFactory("");
  stat_factory->merge(stat_factory_);
  return stat_factory;
}

void CycleRule::push(abi::instruction::Instruction *instruction) {
  assert(instruction != nullptr);
  input_q_->push(instruction);
}

void CycleRule::cycle() {
  service_input_q();
  service_ready_q();

  wait_q_->cycle();
}

void CycleRule::service_input_q() {
  if (input_q_->can_pop() and wait_q_->can_push()) {
    abi::instruction::Instruction *instruction = input_q_->pop();
    int extra_cycle = calculate_extra_cycles(instruction);

    wait_q_->push(instruction, extra_cycle);

    stat_factory_->increment("cycle_rule", extra_cycle);
    stat_factory_->increment(
        std::to_string(instruction->thread()->id()) + "_cycle_rule",
        extra_cycle);
  }
}

void CycleRule::service_ready_q() {
  if (wait_q_->can_pop() and ready_q_->can_push()) {
    abi::instruction::Instruction *instruction = wait_q_->pop();
    ready_q_->push(instruction);

    for (auto &gp_reg : prev_write_gp_regs_[instruction->thread()->id()]) {
      delete gp_reg;
    }
    prev_write_gp_regs_[instruction->thread()->id()].clear();

    for (auto &gp_reg : cur_read_gp_regs_[instruction->thread()->id()]) {
      delete gp_reg;
    }
    cur_read_gp_regs_[instruction->thread()->id()].clear();

    prev_write_gp_regs_[instruction->thread()->id()] =
        collect_write_gp_regs(instruction);
  }
}

int CycleRule::calculate_extra_cycles(
    abi::instruction::Instruction *instruction) {
  cur_read_gp_regs_[instruction->thread()->id()] =
      collect_read_gp_regs(instruction);

  auto [even_counter, odd_counter] = calculate_counters(instruction);

  return even_counter / 2 + odd_counter / 2;
}

std::tuple<int, int> CycleRule::calculate_counters(
    abi::instruction::Instruction *instruction) {
  std::set<upmem_sim::abi::reg::GPReg *> registers =
      merge(prev_write_gp_regs_[instruction->thread()->id()],
            cur_read_gp_regs_[instruction->thread()->id()]);

  int even_counter = 0;
  int odd_counter = 0;
  for (auto &register_ : registers) {
    if (register_->index() % 2 == 0) {
      even_counter += 1;
    } else {
      odd_counter += 1;
    }
  }
  return {even_counter, odd_counter};
}

std::set<abi::reg::GPReg *> CycleRule::collect_read_gp_regs(
    abi::instruction::Instruction *instruction) {
  abi::instruction::Suffix suffix = instruction->suffix();
  if (suffix == abi::instruction::RICI or suffix == abi::instruction::RRI or
      suffix == abi::instruction::RRIC or suffix == abi::instruction::RRICI or
      suffix == abi::instruction::RRIF or suffix == abi::instruction::ZRI or
      suffix == abi::instruction::ZRIC or suffix == abi::instruction::ZRICI or
      suffix == abi::instruction::ZRIF or suffix == abi::instruction::S_RRI or
      suffix == abi::instruction::U_RRI or suffix == abi::instruction::S_RRIC or
      suffix == abi::instruction::U_RRIC or
      suffix == abi::instruction::S_RRICI or
      suffix == abi::instruction::U_RRICI or
      suffix == abi::instruction::S_RRIF or
      suffix == abi::instruction::U_RRIF or suffix == abi::instruction::RR or
      suffix == abi::instruction::RRC or suffix == abi::instruction::RRCI or
      suffix == abi::instruction::ZR or suffix == abi::instruction::ZRC or
      suffix == abi::instruction::ZRCI or suffix == abi::instruction::S_RR or
      suffix == abi::instruction::U_RR or suffix == abi::instruction::S_RRC or
      suffix == abi::instruction::U_RRC or suffix == abi::instruction::S_RRCI or
      suffix == abi::instruction::U_RRCI or suffix == abi::instruction::RIR or
      suffix == abi::instruction::RIRC or suffix == abi::instruction::RIRCI or
      suffix == abi::instruction::ZIR or suffix == abi::instruction::ZIRC or
      suffix == abi::instruction::ZIRCI or suffix == abi::instruction::S_RIRC or
      suffix == abi::instruction::U_RIRC or
      suffix == abi::instruction::S_RIRCI or
      suffix == abi::instruction::U_RIRCI or suffix == abi::instruction::ERRI or
      suffix == abi::instruction::S_ERRI or
      suffix == abi::instruction::U_ERRI or suffix == abi::instruction::EDRI or
      suffix == abi::instruction::ERII) {
    if (instruction->ra()->is_gp_reg()) {
      return {new abi::reg::GPReg(instruction->ra()->gp_reg()->index())};
    } else {
      return {};
    }
  } else if (suffix == abi::instruction::RRR or
             suffix == abi::instruction::RRRC or
             suffix == abi::instruction::RRRCI or
             suffix == abi::instruction::ZRR or
             suffix == abi::instruction::ZRRC or
             suffix == abi::instruction::ZRRCI or
             suffix == abi::instruction::S_RRR or
             suffix == abi::instruction::U_RRR or
             suffix == abi::instruction::S_RRRC or
             suffix == abi::instruction::U_RRRC or
             suffix == abi::instruction::S_RRRCI or
             suffix == abi::instruction::U_RRRCI or
             suffix == abi::instruction::RRRI or
             suffix == abi::instruction::RRRICI or
             suffix == abi::instruction::ZRRI or
             suffix == abi::instruction::ZRRICI or
             suffix == abi::instruction::S_RRRI or
             suffix == abi::instruction::U_RRRI or
             suffix == abi::instruction::S_RRRICI or
             suffix == abi::instruction::U_RRRICI or
             suffix == abi::instruction::ERIR or
             suffix == abi::instruction::DMA_RRI) {
    if (instruction->ra()->is_gp_reg() and instruction->rb()->is_gp_reg()) {
      return {new abi::reg::GPReg(instruction->ra()->gp_reg()->index()),
              new abi::reg::GPReg(instruction->rb()->gp_reg()->index())};
    } else if (instruction->ra()->is_gp_reg()) {
      return {new abi::reg::GPReg(instruction->ra()->gp_reg()->index())};
    } else if (instruction->rb()->is_gp_reg()) {
      return {new abi::reg::GPReg(instruction->rb()->gp_reg()->index())};
    } else {
      return {};
    }
  } else if (suffix == abi::instruction::DRDICI or
             suffix == abi::instruction::ERID) {
    if (instruction->ra()->is_gp_reg()) {
      return {new abi::reg::GPReg(instruction->ra()->gp_reg()->index()),
              new abi::reg::GPReg(instruction->db()->even_reg()->index()),
              new abi::reg::GPReg(instruction->db()->odd_reg()->index())};
    } else {
      return {new abi::reg::GPReg(instruction->db()->even_reg()->index()),
              new abi::reg::GPReg(instruction->db()->odd_reg()->index())};
    }
  } else if (suffix == abi::instruction::R or suffix == abi::instruction::RCI or
             suffix == abi::instruction::Z or suffix == abi::instruction::ZCI or
             suffix == abi::instruction::S_R or
             suffix == abi::instruction::U_R or
             suffix == abi::instruction::S_RCI or
             suffix == abi::instruction::U_RCI or
             suffix == abi::instruction::CI or suffix == abi::instruction::I) {
    return {};
  } else if (suffix == abi::instruction::DDCI) {
    return {new abi::reg::GPReg(instruction->db()->even_reg()->index()),
            new abi::reg::GPReg(instruction->db()->odd_reg()->index())};
  } else {
    throw std::invalid_argument("");
  }
}

std::set<abi::reg::GPReg *> CycleRule::collect_write_gp_regs(
    abi::instruction::Instruction *instruction) {
  abi::instruction::Suffix suffix = instruction->suffix();
  if (suffix == abi::instruction::RICI or suffix == abi::instruction::ZRI or
      suffix == abi::instruction::ZRIC or suffix == abi::instruction::ZRICI or
      suffix == abi::instruction::ZRIF or suffix == abi::instruction::ZRR or
      suffix == abi::instruction::ZRRC or suffix == abi::instruction::ZRRCI or
      suffix == abi::instruction::ZR or suffix == abi::instruction::ZRC or
      suffix == abi::instruction::ZRCI or suffix == abi::instruction::ZRRI or
      suffix == abi::instruction::ZRRICI or suffix == abi::instruction::ZIR or
      suffix == abi::instruction::ZIRC or suffix == abi::instruction::ZIRCI or
      suffix == abi::instruction::Z or suffix == abi::instruction::ZCI or
      suffix == abi::instruction::CI or suffix == abi::instruction::I or
      suffix == abi::instruction::ERII or suffix == abi::instruction::ERIR or
      suffix == abi::instruction::ERID or suffix == abi::instruction::DMA_RRI) {
    return {};
  } else if (suffix == abi::instruction::RRI or
             suffix == abi::instruction::RRIC or
             suffix == abi::instruction::RRICI or
             suffix == abi::instruction::RRIF or
             suffix == abi::instruction::RRR or
             suffix == abi::instruction::RRRC or
             suffix == abi::instruction::RRRCI or
             suffix == abi::instruction::RR or
             suffix == abi::instruction::RRC or
             suffix == abi::instruction::RRCI or
             suffix == abi::instruction::RRRI or
             suffix == abi::instruction::RRRICI or
             suffix == abi::instruction::RIR or
             suffix == abi::instruction::RIRC or
             suffix == abi::instruction::RIRCI or
             suffix == abi::instruction::R or suffix == abi::instruction::RCI or
             suffix == abi::instruction::ERRI) {
    return {new abi::reg::GPReg(instruction->rc()->index())};
  } else if (suffix == abi::instruction::S_RRI or
             suffix == abi::instruction::U_RRI or
             suffix == abi::instruction::S_RRIC or
             suffix == abi::instruction::U_RRIC or
             suffix == abi::instruction::S_RRICI or
             suffix == abi::instruction::U_RRICI or
             suffix == abi::instruction::S_RRIF or
             suffix == abi::instruction::U_RRIF or
             suffix == abi::instruction::S_RRR or
             suffix == abi::instruction::U_RRR or
             suffix == abi::instruction::S_RRRC or
             suffix == abi::instruction::U_RRRC or
             suffix == abi::instruction::S_RRRCI or
             suffix == abi::instruction::U_RRRCI or
             suffix == abi::instruction::S_RR or
             suffix == abi::instruction::U_RR or
             suffix == abi::instruction::S_RRC or
             suffix == abi::instruction::U_RRC or
             suffix == abi::instruction::S_RRCI or
             suffix == abi::instruction::U_RRCI or
             suffix == abi::instruction::DRDICI or
             suffix == abi::instruction::S_RRRI or
             suffix == abi::instruction::U_RRRI or
             suffix == abi::instruction::S_RRRICI or
             suffix == abi::instruction::U_RRRICI or
             suffix == abi::instruction::S_RIRC or
             suffix == abi::instruction::U_RIRC or
             suffix == abi::instruction::S_RIRCI or
             suffix == abi::instruction::U_RIRCI or
             suffix == abi::instruction::S_R or
             suffix == abi::instruction::U_R or
             suffix == abi::instruction::DDCI or
             suffix == abi::instruction::S_ERRI or
             suffix == abi::instruction::U_ERRI or
             suffix == abi::instruction::EDRI) {
    return {new abi::reg::GPReg(instruction->dc()->even_reg()->index()),
            new abi::reg::GPReg(instruction->dc()->odd_reg()->index())};
  } else {
    throw std::invalid_argument("");
  }
}

std::set<abi::reg::GPReg *> CycleRule::merge(
    std::set<abi::reg::GPReg *> regs1, std::set<abi::reg::GPReg *> regs2) {
  std::set<abi::reg::GPReg *> regs;

  auto finder = [&regs](RegIndex index) {
    for (auto &reg : regs) {
      if (reg->index() == index) {
        return true;
      }
    }
    return false;
  };

  for (auto &reg : regs1) {
    if (not finder(reg->index())) {
      regs.insert(reg);
    }
  }

  for (auto &reg : regs2) {
    if (not finder(reg->index())) {
      regs.insert(reg);
    }
  }

  return std::move(regs);
}

}  // namespace upmem_sim::simulator::dpu
