#ifndef UPMEM_SIM_SIMULATOR_SRAM_LOCK_H_
#define UPMEM_SIM_SIMULATOR_SRAM_LOCK_H_

#include <cassert>

#include "main.h"

namespace upmem_sim::simulator::sram {

class Lock {
 public:
  explicit Lock() : id_(nullptr) {}
  ~Lock() { assert(id_ == nullptr); }

  bool can_acquire() { return id_ == nullptr; }
  void acquire(ThreadID id);
  bool can_release(ThreadID id) { return id_ == nullptr or *id_ == id; }
  void release(ThreadID id);
  void cycle() = delete;

 private:
  ThreadID *id_;
};

}  // namespace upmem_sim::simulator::sram

#endif
