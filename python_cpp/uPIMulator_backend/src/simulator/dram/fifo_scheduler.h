#ifndef UPMEM_SIM_SIMULATOR_DRAM_FIFO_SCHEDULER_H_
#define UPMEM_SIM_SIMULATOR_DRAM_FIFO_SCHEDULER_H_

#include "simulator/dram/scheduler.h"
#include "util/argument_parser.h"

namespace upmem_sim::simulator::dram {

class FIFOScheduler : public Scheduler {
 public:
  explicit FIFOScheduler(util::ArgumentParser *argument_parser)
      : Scheduler(argument_parser) {}
  ~FIFOScheduler() { assert(reorder_buffer_.empty()); }

  void cycle() final;

 protected:
  void service_input_q();
  void service_dma_command(dpu::DMACommand *dma_command);
  void service_output_q();

  bool service_fcfs();
};

}  // namespace upmem_sim::simulator::dram

#endif
