#ifndef UPMEM_SIM_SIMULATOR_CPU_INIT_THREAD_H_
#define UPMEM_SIM_SIMULATOR_CPU_INIT_THREAD_H_

#include "simulator/cpu/thread.h"
#include "simulator/rank/rank.h"

namespace upmem_sim::simulator::cpu {

class InitThread : public Thread {
 public:
  explicit InitThread(util::ArgumentParser *argument_parser)
      : Thread(argument_parser), rank_(nullptr) {}
  ~InitThread() = default;

  void connect_rank(rank::Rank *rank);

  void init();
  void launch();

  void cycle() = delete;

 protected:
  void dma_transfer_to_atomic();
  void dma_transfer_to_iram();
  void dma_transfer_to_wram();
  void dma_transfer_to_mram();

 private:
  rank::Rank *rank_;
};

}  // namespace upmem_sim::simulator::cpu

#endif
