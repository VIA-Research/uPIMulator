#include "simulator/sram/wram.h"

namespace upmem_sim::simulator::sram {

WRAM::WRAM() {
  address_ = new abi::word::DataAddressWord();
  address_->set_value(util::ConfigLoader::wram_offset());

  size_ = util::ConfigLoader::wram_size();

  assert(address() % abi::word::DataWord().size() == 0);
  assert(size_ % abi::word::DataWord().size() == 0);

  cells_.resize(num_data_words());
  for (int i = 0; i < num_data_words(); i++) {
    cells_[i] = new abi::word::DataWord();
  }
}

WRAM::~WRAM() {
  delete address_;

  for (int i = 0; i < num_data_words(); i++) {
    delete cells_[i];
  }
}

int64_t WRAM::read(Address address) {
  return cells_[index(address)]->value(abi::word::UNSIGNED);
}

void WRAM::write(Address address, int64_t value) {
  cells_[index(address)]->set_value(value);
}

void WRAM::write(Address address, encoder::ByteStream *byte_stream) {
  cells_[index(address)]->from_byte_stream(byte_stream);
}

int WRAM::index(Address address) {
  assert(address >= this->address());
  assert(address + abi::word::DataWord().size() <= this->address() + size_);
  assert((address - this->address()) % abi::word::DataWord().size() == 0);

  return static_cast<int>((address - this->address()) /
                          abi::word::DataWord().size());
}

}  // namespace upmem_sim::simulator::sram