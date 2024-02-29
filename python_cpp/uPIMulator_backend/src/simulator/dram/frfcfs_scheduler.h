#ifndef UPMEM_SIM_SIMULATOR_DRAM_FRFCFS_SCHEDULER_H_
#define UPMEM_SIM_SIMULATOR_DRAM_FRFCFS_SCHEDULER_H_

#include "simulator/dram/scheduler.h"
#include "util/argument_parser.h"

namespace upmem_sim::simulator::dram {

class FRFCFSScheduler : public Scheduler {
 public:
  explicit FRFCFSScheduler(util::ArgumentParser *argument_parser)
      : Scheduler(argument_parser) {}
  ~FRFCFSScheduler() { assert(reorder_buffer_.empty()); }

  void cycle() final;

 protected:
  void service_input_q();
  void service_dma_command(dpu::DMACommand *dma_command);
  void service_output_q();

  bool service_fr();
  bool service_fcfs();
};

}  // namespace upmem_sim::simulator::dram

#endif
