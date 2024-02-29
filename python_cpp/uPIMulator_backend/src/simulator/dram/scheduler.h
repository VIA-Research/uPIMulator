#ifndef UPMEM_SIM_SIMULATOR_DRAM_SCHEDULER_H_
#define UPMEM_SIM_SIMULATOR_DRAM_SCHEDULER_H_

#include "simulator/basic/queue.h"
#include "simulator/dpu/dma_command.h"
#include "simulator/dram/memory_command.h"
#include "util/argument_parser.h"
#include "util/stat_factory.h"

namespace upmem_sim::simulator::dram {

class Scheduler {
 public:
  using MemoryReference = std::tuple<dpu::DMACommand *, Address, Address>;

  explicit Scheduler(util::ArgumentParser *argument_parser);
  ~Scheduler();

  util::StatFactory *stat_factory();

  bool empty() {
    return input_q_->empty() and ready_q_->empty() and reorder_buffer_.empty();
  }

  bool can_push() { return input_q_->can_push(); }
  void push(dpu::DMACommand *dma_command);
  bool can_pop() { return ready_q_->can_pop(); }
  MemoryCommand *pop();

  void flush();

  virtual void cycle() = 0;

 protected:
  basic::Queue<dpu::DMACommand> *input_q_;
  basic::Queue<MemoryCommand> *ready_q_;
  abi::word::DataAddressWord *row_address_;

  std::vector<MemoryReference> reorder_buffer_;

  Address wordline_size_;

  util::StatFactory *stat_factory_;
};

}  // namespace upmem_sim::simulator::dram

#endif