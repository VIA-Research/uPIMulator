#ifndef UPMEM_SIM_SIMULATOR_SRAM_ATOMIC_H_
#define UPMEM_SIM_SIMULATOR_SRAM_ATOMIC_H_

#include <vector>

#include "abi/word/data_address_word.h"
#include "simulator/sram/lock.h"
#include "util/config_loader.h"

namespace upmem_sim::simulator::sram {

class Atomic {
 public:
  explicit Atomic();
  ~Atomic();

  Address address() { return address_->address(); }
  Address size() { return size_; }

  bool can_acquire(Address address) {
    return locks_[index(address)]->can_acquire();
  }
  void acquire(Address address, ThreadID id) {
    locks_[index(address)]->acquire(id);
  }
  bool can_release(Address address, ThreadID id) {
    return locks_[index(address)]->can_release(id);
  }
  void release(Address address, ThreadID id) {
    locks_[index(address)]->release(id);
  }
  void cycle() = delete;

 protected:
  int index(Address address);

 private:
  abi::word::DataAddressWord *address_;
  Address size_;

  std::vector<Lock *> locks_;
};

}  // namespace upmem_sim::simulator::sram

#endif
