#ifndef UPMEM_SIM_SIMULATOR_CPU_FINI_THREAD_H_
#define UPMEM_SIM_SIMULATOR_CPU_FINI_THREAD_H_

#include "simulator/cpu/thread.h"
#include "simulator/rank/rank.h"

namespace upmem_sim::simulator::cpu {

class FiniThread : public Thread {
 public:
  explicit FiniThread(util::ArgumentParser *argument_parser)
      : Thread(argument_parser), rank_(nullptr) {}
  ~FiniThread() = default;

  void connect_rank(rank::Rank *rank);

  void fini() = delete;

  void cycle();

 private:
  rank::Rank *rank_;
};

}  // namespace upmem_sim::simulator::cpu

#endif
