#ifndef UPMEM_SIM_SIMULATOR_DPU_REVOLVER_SCHEDULER_H_
#define UPMEM_SIM_SIMULATOR_DPU_REVOLVER_SCHEDULER_H_

#include <vector>

#include "simulator/basic/queue.h"
#include "simulator/dpu/thread.h"
#include "util/argument_parser.h"
#include "util/stat_factory.h" 

namespace upmem_sim::simulator::dpu {

class RevolverScheduler {
 public:
  explicit RevolverScheduler(util::ArgumentParser *argument_parser,
                             std::vector<Thread *> threads);
  ~RevolverScheduler();

  util::StatFactory *stat_factory();
  std::vector<Thread *> threads() { return threads_; }

  Thread *schedule();

  bool boot(ThreadID id);
  bool sleep(ThreadID id);
  bool block(ThreadID id);
  bool awake(ThreadID id);
  bool shutdown(ThreadID id);

  void cycle();

  int get_issuable_threads() { return issuable_threads_; };

 private:
  int num_revolver_scheduling_cycles_;
  int issuable_threads_;

  std::vector<Thread *> threads_;
  basic::Queue<Thread> *thread_q_;

  util::StatFactory *stat_factory_;
};

}  // namespace upmem_sim::simulator::dpu

#endif
