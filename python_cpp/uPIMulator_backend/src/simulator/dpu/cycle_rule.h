#ifndef UPMEM_SIM_SIMULATOR_DPU_CYCLE_RULE_H_
#define UPMEM_SIM_SIMULATOR_DPU_CYCLE_RULE_H_

#include "abi/instruction/instruction.h"
#include "simulator/basic/queue.h"
#include "simulator/basic/timer_queue.h"
#include "util/argument_parser.h"
#include "util/stat_factory.h"

namespace upmem_sim::simulator::dpu {

class CycleRule {
 public:
  explicit CycleRule(util::ArgumentParser *argument_parser);
  ~CycleRule();

  util::StatFactory *stat_factory();

  bool empty() {
    return input_q_->empty() and wait_q_->empty() and ready_q_->empty();
  }
  bool can_push() { return input_q_->can_push(); }
  void push(abi::instruction::Instruction *instruction);
  bool can_pop() { return ready_q_->can_pop(); }
  abi::instruction::Instruction *pop() { return ready_q_->pop(); }
  void cycle();

 protected:
  void service_input_q();
  void service_ready_q();

  int calculate_extra_cycles(abi::instruction::Instruction *instruction);
  std::tuple<int, int> calculate_counters(
      abi::instruction::Instruction *instruction);

  static std::set<abi::reg::GPReg *> collect_read_gp_regs(
      abi::instruction::Instruction *instruction);
  static std::set<abi::reg::GPReg *> collect_write_gp_regs(
      abi::instruction::Instruction *instruction);
  static std::set<abi::reg::GPReg *> merge(std::set<abi::reg::GPReg *> regs1,
                                           std::set<abi::reg::GPReg *> regs2);

 private:
  basic::Queue<abi::instruction::Instruction> *input_q_;
  basic::TimerQueue<abi::instruction::Instruction> *wait_q_;
  basic::Queue<abi::instruction::Instruction> *ready_q_;

  std::vector<std::set<abi::reg::GPReg *>> prev_write_gp_regs_;
  std::vector<std::set<abi::reg::GPReg *>> cur_read_gp_regs_;

  util::StatFactory *stat_factory_;
};

}  // namespace upmem_sim::simulator::dpu

#endif
