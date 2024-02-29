#include "simulator/sram/lock.h"

namespace upmem_sim::simulator::sram {

void Lock::acquire(ThreadID id) {
  assert(can_acquire());
  id_ = new ThreadID(id);
}

void Lock::release(ThreadID id) {
  assert(can_release(id));

  delete id_;
  id_ = nullptr;
}

}  // namespace upmem_sim::simulator::sram
