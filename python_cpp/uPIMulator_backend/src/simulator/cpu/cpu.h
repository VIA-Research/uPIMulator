#ifndef UPMEM_SIM_SIMULATOR_CPU_CPU_H_
#define UPMEM_SIM_SIMULATOR_CPU_CPU_H_

#include "simulator/cpu/fini_thread.h"
#include "simulator/cpu/init_thread.h"
#include "simulator/cpu/sched_thread.h"
#include "simulator/rank/rank.h"

namespace upmem_sim::simulator::cpu {

class CPU {
 public:
  explicit CPU(util::ArgumentParser *argument_parser)
      : init_thread_(new InitThread(argument_parser)),
        sched_thread_(new SchedThread(argument_parser)),
        fini_thread_(new FiniThread(argument_parser)) {}
  ~CPU();

  void connect_rank(rank::Rank *rank);

  int num_executions() { return fini_thread_->num_executions(); }

  void init() { init_thread_->init(); }
  void launch() { init_thread_->launch(); }
  void sched(int execution) { sched_thread_->sched(execution); }
  void check(int execution) { sched_thread_->check(execution); }
  void fini() {}
  void cycle() { fini_thread_->cycle(); }

 private:
  InitThread *init_thread_;
  SchedThread *sched_thread_;
  FiniThread *fini_thread_;
};

}  // namespace upmem_sim::simulator::cpu

#endif
