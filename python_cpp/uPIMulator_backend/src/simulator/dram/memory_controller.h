#ifndef UPMEM_SIM_SIMULATOR_DRAM_MEMORY_CONTROLLER_H_
#define UPMEM_SIM_SIMULATOR_DRAM_MEMORY_CONTROLLER_H_

#include "simulator/dram/mram.h"
#include "simulator/dram/row_buffer.h"
#include "simulator/dram/scheduler.h"

namespace upmem_sim::simulator::dram {

class MemoryController {
 public:
  explicit MemoryController(util::ArgumentParser *argument_parser);
  ~MemoryController();

  util::StatFactory *stat_factory();

  void connect_mram(MRAM *mram);

  bool empty() {
    return input_q_->empty() and wait_q_->empty() and
           memory_command_q_->empty() and ready_q_->empty() and
           scheduler_->empty() and row_buffer_->empty();
  }

  bool can_push() { return input_q_->can_push(); }
  void push(dpu::DMACommand *dma_command);
  bool can_pop() { return ready_q_->can_pop(); }
  dpu::DMACommand *pop();
  dpu::DMACommand *front();

  std::vector<int> read(Address address, Address size);

  void write(Address address, Address size, std::vector<int> bytes);
  void write(Address address, Address size, encoder::ByteStream *byte_stream);

  void flush();

  void cycle();

 protected:
  void service_input_q();
  void service_scheduler();
  void service_memory_command_q();
  void service_row_buffer();
  void service_wait_q();

 private:
  Address wordline_size_;

  Scheduler *scheduler_;
  RowBuffer *row_buffer_;
  MRAM *mram_;

  basic::Queue<dpu::DMACommand> *input_q_;
  basic::Queue<dpu::DMACommand> *wait_q_;
  basic::Queue<MemoryCommand> *memory_command_q_;
  basic::Queue<dpu::DMACommand> *ready_q_;

  util::StatFactory *stat_factory_;
};

}  // namespace upmem_sim::simulator::dram

#endif
