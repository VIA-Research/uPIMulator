#include "simulator/dpu/operand_collector.h"

#include <tuple>

namespace upmem_sim::simulator::dpu {

void OperandCollector::connect_wram(sram::WRAM *wram) {
  assert(wram != nullptr);
  assert(wram_ == nullptr);

  wram_ = wram;
}

int64_t OperandCollector::lbs(Address address) {
  Address data_word_size = abi::word::DataWord().size();
  Address base_address = (address / data_word_size) * data_word_size;
  Address offset = address % data_word_size;

  auto data_word = new abi::word::DataWord();
  data_word->set_value(wram_->read(base_address));

  int64_t result =
      data_word->bit_slice(abi::word::SIGNED, static_cast<int>(8 * offset),
                           static_cast<int>(8 * (offset + 1)));

  delete data_word;

  return result;
}

int64_t OperandCollector::lbu(Address address) {
  Address data_word_size = abi::word::DataWord().size();
  Address base_address = (address / data_word_size) * data_word_size;
  Address offset = address % data_word_size;

  auto data_word = new abi::word::DataWord();
  data_word->set_value(wram_->read(base_address));

  int64_t result =
      data_word->bit_slice(abi::word::UNSIGNED, static_cast<int>(8 * offset),
                           static_cast<int>(8 * (offset + 1)));

  delete data_word;

  return result;
}

int64_t OperandCollector::lhs(Address address) {
  auto data_word = new abi::word::DataWord();

  data_word->set_bit_slice(0, 8, lbs(address));
  data_word->set_bit_slice(8, 16, lbs(address + 1));

  int64_t result = data_word->bit_slice(abi::word::SIGNED, 0, 16);

  delete data_word;

  return result;
}

int64_t OperandCollector::lhu(Address address) {
  auto data_word = new abi::word::DataWord();

  data_word->set_bit_slice(0, 8, lbu(address));
  data_word->set_bit_slice(8, 16, lbu(address + 1));

  int64_t result = data_word->bit_slice(abi::word::UNSIGNED, 0, 16);

  delete data_word;

  return result;
}

int64_t OperandCollector::lw(Address address) {
  auto data_word = new abi::word::DataWord();

  data_word->set_bit_slice(0, 8, lbu(address));
  data_word->set_bit_slice(8, 16, lbu(address + 1));
  data_word->set_bit_slice(16, 24, lbu(address + 2));
  data_word->set_bit_slice(24, 32, lbu(address + 3));

  int64_t result = data_word->value(abi::word::UNSIGNED);

  delete data_word;

  return result;
}

std::tuple<int64_t, int64_t> OperandCollector::ld(Address address) {
  return {lw(address + abi::word::DataWord().size()), lw(address)};
}

void OperandCollector::sb(Address address, int64_t value) {
  Address data_word_size = abi::word::DataWord().size();
  Address base_address = (address / data_word_size) * data_word_size;
  Address offset = address % data_word_size;

  auto data_word = new abi::word::DataWord();
  data_word->set_value(wram_->read(base_address));
  data_word->set_bit_slice(static_cast<int>(8 * offset),
                           static_cast<int>(8 * (offset + 1)), value);

  wram_->write(base_address, data_word->value(abi::word::UNSIGNED));

  delete data_word;
}

void OperandCollector::sh(Address address, int64_t value) {
  auto data_word = new abi::word::DataWord();
  data_word->set_value(value);

  sb(address, data_word->bit_slice(abi::word::UNSIGNED, 0, 8));
  sb(address + 1, data_word->bit_slice(abi::word::UNSIGNED, 8, 16));

  delete data_word;
}

void OperandCollector::sw(Address address, int64_t value) {
  auto data_word = new abi::word::DataWord();
  data_word->set_value(value);

  sb(address, data_word->bit_slice(abi::word::UNSIGNED, 0, 8));
  sb(address + 1, data_word->bit_slice(abi::word::UNSIGNED, 8, 16));
  sb(address + 2, data_word->bit_slice(abi::word::UNSIGNED, 16, 24));
  sb(address + 3, data_word->bit_slice(abi::word::UNSIGNED, 24, 32));

  delete data_word;
}

void OperandCollector::sd(Address address, int64_t even, int64_t odd) {
  sw(address + abi::word::DataWord().size(), even);
  sw(address, odd);
}

}  // namespace upmem_sim::simulator::dpu
