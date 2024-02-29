#ifndef UPMEM_SIM_SIMULATOR_SYSTEM_H_
#define UPMEM_SIM_SIMULATOR_SYSTEM_H_

#include "simulator/cpu/cpu.h"
#include "simulator/dpu/dpu.h"
#include "simulator/rank/rank.h"

namespace upmem_sim::simulator {

class System {
 public:
  explicit System(util::ArgumentParser *argument_parser);
  ~System();

  util::StatFactory *stat_factory();

  bool is_finished() { return execuion_ == cpu_->num_executions(); }

  void init();
  void fini() { cpu_->fini(); }
  void cycle();

 protected:
  bool is_zombie() { return rank_->is_zombie(); }

 private:
  cpu::CPU *cpu_;
  rank::Rank *rank_;

  int execuion_;

  util::StatFactory *stat_factory_;

  // Tae
  std::string benchmark;
};

}  // namespace upmem_sim::simulator

#endif
