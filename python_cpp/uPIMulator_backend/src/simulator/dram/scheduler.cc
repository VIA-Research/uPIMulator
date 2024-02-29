#include "simulator/dram/scheduler.h"

namespace upmem_sim::simulator::dram {

Scheduler::Scheduler(util::ArgumentParser *argument_parser)
    : input_q_(new basic::Queue<dpu::DMACommand>(-1)),
      ready_q_(new basic::Queue<MemoryCommand>(3)),
      row_address_(nullptr),
      stat_factory_(new util::StatFactory("Scheduler")) {
  wordline_size_ = argument_parser->get_int_parameter("wordline_size");

  assert(wordline_size_ > 0);
  assert(wordline_size_ % util::ConfigLoader::min_access_granularity() == 0);
}

Scheduler::~Scheduler() {
  delete input_q_;
  delete ready_q_;
  delete row_address_;

  delete stat_factory_;
}

util::StatFactory *Scheduler::stat_factory() {
  auto stat_factory = new util::StatFactory("");

  stat_factory->merge(stat_factory_);

  return stat_factory;
}

void Scheduler::push(dpu::DMACommand *dma_command) {
  assert(dma_command != nullptr);
  input_q_->push(dma_command);
}

MemoryCommand *Scheduler::pop() {
  assert(can_pop());
  return ready_q_->pop();
}

void Scheduler::flush() {
  delete row_address_;
  row_address_ = nullptr;
}

}  // namespace upmem_sim::simulator::dram