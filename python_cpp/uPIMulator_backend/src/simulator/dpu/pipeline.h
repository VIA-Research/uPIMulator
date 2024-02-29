#ifndef UPMEM_SIM_SIMULATOR_DPU_PIPELINE_H_
#define UPMEM_SIM_SIMULATOR_DPU_PIPELINE_H_

#include "abi/instruction/instruction.h"
#include "simulator/basic/queue.h"
#include "util/argument_parser.h"

namespace upmem_sim::simulator::dpu {

class Pipeline {
 public:
  explicit Pipeline(util::ArgumentParser *argument_parser);
  ~Pipeline();

  bool empty() {
    return empty_input_q() and empty_wait_q() and empty_ready_q();
  }
  bool can_push() { return input_q_->can_push(); }
  void push(abi::instruction::Instruction *instruction);
  bool can_pop() { return ready_q_->can_pop(); }
  abi::instruction::Instruction *pop() { return ready_q_->pop(); }
  void cycle();

 protected:
  bool empty_input_q() { return input_q_->empty(); }
  bool empty_wait_q();
  bool empty_ready_q() {
    return ready_q_->empty() or ready_q_->front() == nullptr;
  }

  void service_input_q();
  void service_wait_q();

 private:
  basic::Queue<abi::instruction::Instruction> *input_q_;
  basic::Queue<abi::instruction::Instruction> *wait_q_;
  basic::Queue<abi::instruction::Instruction> *ready_q_;
};

}  // namespace upmem_sim::simulator::dpu

#endif
