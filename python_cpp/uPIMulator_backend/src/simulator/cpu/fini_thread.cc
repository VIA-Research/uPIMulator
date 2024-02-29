#include "simulator/cpu/fini_thread.h"

namespace upmem_sim::simulator::cpu {

void FiniThread::connect_rank(rank::Rank *rank) {
  assert(rank != nullptr);
  assert(rank_ == nullptr);

  rank_ = rank;
}

void FiniThread::cycle() {
  for (auto &dpu : rank_->dpus()) {
    for (auto &thread : dpu->scheduler()->threads()) {
      if (thread->reg_file()->read_pc_reg() == sys_end_pointer() and
          thread->state() == dpu::Thread::SLEEP) {
        dpu->scheduler()->shutdown(thread->id());
      }
    }
  }
}

}  // namespace upmem_sim::simulator::cpu
