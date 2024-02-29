#include "simulator/cpu/cpu.h"

namespace upmem_sim::simulator::cpu {

CPU::~CPU() {
  delete init_thread_;
  delete sched_thread_;
  delete fini_thread_;
}

void CPU::connect_rank(rank::Rank* rank) {
  init_thread_->connect_rank(rank);
  sched_thread_->connect_rank(rank);
  fini_thread_->connect_rank(rank);
}

}  // namespace upmem_sim::simulator::cpu
