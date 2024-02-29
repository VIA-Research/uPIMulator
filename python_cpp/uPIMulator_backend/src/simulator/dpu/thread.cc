#include "simulator/dpu/thread.h"

namespace upmem_sim::simulator::dpu {

Thread::~Thread() {
  assert(state_ == ZOMBIE);
  delete reg_file_;
}

}  // namespace upmem_sim::simulator::dpu