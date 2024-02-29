#include "simulator/dpu/pipeline.h"

#include "converter/instruction_converter.h"

namespace upmem_sim::simulator::dpu {

Pipeline::Pipeline(util::ArgumentParser *argument_parser)
    : input_q_(new basic::Queue<abi::instruction::Instruction>(1)),
      ready_q_(new basic::Queue<abi::instruction::Instruction>(1)) {
  int num_pipeline_stages = static_cast<int>(
      argument_parser->get_int_parameter("num_pipeline_stages"));
  assert(num_pipeline_stages > 1);
  wait_q_ =
      new basic::Queue<abi::instruction::Instruction>(num_pipeline_stages - 1);

  while (wait_q_->can_push()) {
    wait_q_->push(nullptr);
  }
  ready_q_->push(nullptr);
}

Pipeline::~Pipeline() {
  while (wait_q_->can_pop()) {
    if (wait_q_->front() != nullptr) {
      converter::InstructionConverter::to_string(wait_q_->front());
    }

    assert(wait_q_->pop() == nullptr);
  }

  while (ready_q_->can_pop()) {
    assert(ready_q_->pop() == nullptr);
  }

  delete input_q_;
  delete wait_q_;
  delete ready_q_;
}

void Pipeline::push(abi::instruction::Instruction *instruction) {
  assert(can_push());
  assert(instruction != nullptr);

  input_q_->push(instruction);
}

void Pipeline::cycle() {
  service_input_q();
  service_wait_q();
}

bool Pipeline::empty_wait_q() {
  std::queue<abi::instruction::Instruction *> wait_q;
  while (wait_q_->can_pop()) {
    abi::instruction::Instruction *instruction = wait_q_->pop();
    wait_q.push(instruction);
  }

  while (not wait_q.empty()) {
    abi::instruction::Instruction *instruction = wait_q.front();
    wait_q.pop();

    wait_q_->push(instruction);
    if (instruction != nullptr) {
      return false;
    }
  }
  return true;
}

void Pipeline::service_input_q() {
  if (input_q_->can_pop() and wait_q_->can_push()) {
    abi::instruction::Instruction *instruction = input_q_->pop();
    wait_q_->push(instruction);
  } else if (wait_q_->can_push()) {
    wait_q_->push(nullptr);
  }
}

void Pipeline::service_wait_q() {
  if (wait_q_->can_pop() and ready_q_->can_push()) {
    abi::instruction::Instruction *instruction = wait_q_->pop();
    ready_q_->push(instruction);
  }
}

}  // namespace upmem_sim::simulator::dpu
