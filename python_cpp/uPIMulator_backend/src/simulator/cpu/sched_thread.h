#ifndef UPMEM_SIM_SIMULATOR_CPU_SCHED_THREAD_H_
#define UPMEM_SIM_SIMULATOR_CPU_SCHED_THREAD_H_

#include "simulator/cpu/thread.h"
#include "simulator/rank/rank.h"

namespace upmem_sim::simulator::cpu {

class SchedThread : public Thread {
 public:
  explicit SchedThread(util::ArgumentParser *argument_parser)
      : Thread(argument_parser), rank_(nullptr) {}
  ~SchedThread() = default;

  void connect_rank(rank::Rank *rank);

  void sched(int execution);
  void check(int execution);

  void cycle() = delete;

 protected:
  void dma_transfer_input_dpu_mram_heap_pointer_name(int execution);
  void dma_transfer_dpu_input_arguments(int execution);

  void dma_transfer_output_dpu_mram_heap_pointer_name(int execution);
  void dma_transfer_dpu_results(int execution);

 private:
  rank::Rank *rank_;
};

}  // namespace upmem_sim::simulator::cpu

#endif
