#include "simulator/sram/atomic.h"

namespace upmem_sim::simulator::sram {

Atomic::Atomic() {
  address_ = new abi::word::DataAddressWord();
  address_->set_value(util::ConfigLoader::atomic_offset());

  size_ = util::ConfigLoader::atomic_size();

  locks_.resize(size_);
  for (int i = 0; i < size_; i++) {
    locks_[i] = new Lock();
  }
}

Atomic::~Atomic() {
  delete address_;

  for (int i = 0; i < size_; i++) {
    delete locks_[i];
  }
}

int Atomic::index(Address address) {
  assert(address >= this->address());
  assert(address < this->address() + size_);
  assert(this->address() <= address and address < this->address() + size_);

  return static_cast<int>(address - this->address());
}

}  // namespace upmem_sim::simulator::sram
