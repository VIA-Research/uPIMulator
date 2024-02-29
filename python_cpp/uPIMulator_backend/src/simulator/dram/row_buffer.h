#ifndef UPMEM_SIM_SIMULATOR_DRAM_ROW_BUFFER_H_
#define UPMEM_SIM_SIMULATOR_DRAM_ROW_BUFFER_H_

#include <map>
#include <string>

#include "simulator/basic/queue.h"
#include "simulator/basic/timer_queue.h"
#include "simulator/dram/memory_command.h"
#include "simulator/dram/mram.h"
#include "util/argument_parser.h"
#include "util/stat_factory.h"

namespace upmem_sim::simulator::dram {

class RowBuffer {
 public:
  explicit RowBuffer(util::ArgumentParser *argument_parser);
  ~RowBuffer();

  util::StatFactory *stat_factory();

  void connect_mram(MRAM *mram);

  bool empty() {
    return input_q_->empty() and ready_q_->empty() and
           activation_q_->empty() and io_q_->empty() and bus_q_->empty() and
           precharge_q_->empty();
  }

  bool can_push() { return input_q_->can_push(); }
  void push(MemoryCommand *memory_command);
  bool can_pop() { return ready_q_->can_pop(); }
  MemoryCommand *pop();

  void flush();

  void cycle();

 protected:
  void service_input_q();
  void service_activation_q();
  void service_io_q();
  void service_bus_q();
  void service_precharge_q();

  std::vector<int> read_from_mram();
  std::vector<int> read_from_row_buffer(Address address, Address size);

  void write_to_mram();
  void write_to_row_buffer(Address address, Address size,
                           std::vector<int> bytes);

  int index(Address address);

 private:
  std::map<std::string, int> timing_parameters_;
  Address wordline_size_;

  MRAM *mram_;
  abi::word::DataAddressWord *row_address_;
  std::vector<int> row_buffer_;

  basic::Queue<MemoryCommand> *input_q_;
  basic::Queue<MemoryCommand> *ready_q_;

  basic::TimerQueue<MemoryCommand> *activation_q_;
  basic::TimerQueue<MemoryCommand> *io_q_;
  basic::TimerQueue<MemoryCommand> *bus_q_;
  basic::TimerQueue<MemoryCommand> *precharge_q_;

  util::StatFactory *stat_factory_;
};

}  // namespace upmem_sim::simulator::dram

#endif
