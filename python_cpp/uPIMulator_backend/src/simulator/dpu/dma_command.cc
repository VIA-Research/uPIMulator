#include "simulator/dpu/dma_command.h"

#include <iostream>

#include "abi/word/data_word.h"

namespace upmem_sim::simulator::dpu {

DMACommand::DMACommand(Operation operation, Address mram_address, Address size)
    : operation_(operation),
      wram_address_(nullptr),
      mram_address_(new abi::word::DataAddressWord()),
      size_(size),
      instruction_(nullptr) {
  assert(operation == READ);
  assert(mram_address >= util::ConfigLoader::mram_offset());
  assert(mram_address + size_ <=
         util::ConfigLoader::mram_offset() + util::ConfigLoader::mram_size());
  assert(size % util::ConfigLoader::min_access_granularity() == 0);

  mram_address_->set_value(mram_address);
  bytes_.resize(size);
  acks_.resize(size);
}

DMACommand::DMACommand(Operation operation, Address wram_address,
                       Address mram_address, Address size,
                       abi::instruction::Instruction *instruction)
    : operation_(operation),
      wram_address_(new abi::word::DataAddressWord()),
      mram_address_(new abi::word::DataAddressWord()),
      size_(size),
      instruction_(instruction) {
  assert(operation == READ);
  assert(wram_address >= util::ConfigLoader::wram_offset());
  assert(wram_address + size_ <=
         util::ConfigLoader::wram_offset() + util::ConfigLoader::wram_size());
  assert(mram_address >= util::ConfigLoader::mram_offset());
  assert(mram_address + size_ <=
         util::ConfigLoader::mram_offset() + util::ConfigLoader::mram_size());
  assert(size % util::ConfigLoader::min_access_granularity() == 0);
  assert(instruction->op_code() == abi::instruction::LDMA);

  wram_address_->set_value(wram_address);
  mram_address_->set_value(mram_address);
  bytes_.resize(size);
  acks_.resize(size);
}

DMACommand::DMACommand(Operation operation, Address mram_address, Address size,
                       std::vector<int> bytes)
    : operation_(operation),
      wram_address_(nullptr),
      mram_address_(new abi::word::DataAddressWord()),
      size_(size),
      bytes_(bytes),
      instruction_(nullptr) {
  assert(operation == WRITE);
  assert(mram_address >= util::ConfigLoader::mram_offset());
  assert(mram_address + size_ <=
         util::ConfigLoader::mram_offset() + util::ConfigLoader::mram_size());
  assert(size % util::ConfigLoader::min_access_granularity() == 0);
  assert(size == bytes.size());

  mram_address_->set_value(mram_address);
  bytes_.resize(size);
  acks_.resize(size);
}

DMACommand::DMACommand(Operation operation, Address wram_address,
                       Address mram_address, Address size,
                       std::vector<int> bytes,
                       abi::instruction::Instruction *instruction)
    : operation_(operation),
      wram_address_(new abi::word::DataAddressWord()),
      mram_address_(new abi::word::DataAddressWord()),
      size_(size),
      bytes_(bytes),
      instruction_(instruction) {
  assert(operation == WRITE);
  assert(wram_address >= util::ConfigLoader::wram_offset());
  assert(wram_address + size_ <=
         util::ConfigLoader::wram_offset() + util::ConfigLoader::wram_size());
  assert(mram_address >= util::ConfigLoader::mram_offset());
  assert(mram_address + size_ <=
         util::ConfigLoader::mram_offset() + util::ConfigLoader::mram_size());
  assert(size % util::ConfigLoader::min_access_granularity() == 0);
  assert(size == bytes.size());
  assert(instruction->op_code() == abi::instruction::SDMA);

  wram_address_->set_value(wram_address);
  mram_address_->set_value(mram_address);
  bytes_.resize(size);
  acks_.resize(size);
}

DMACommand::~DMACommand() {
  assert(is_ready());

  delete wram_address_;
  delete mram_address_;
}

abi::instruction::Instruction *DMACommand::instruction() {
  assert(has_instruction());
  return instruction_;
}

std::vector<int> DMACommand::bytes() {
  if (operation_ == READ) {
    assert(is_ready());
  }

  return bytes_;
}

std::vector<int> DMACommand::bytes(Address mram_address, Address size) {
  std::vector<int> bytes;
  bytes.resize(size);
  std::copy(bytes_.begin() + index(mram_address),
            bytes_.begin() + index(mram_address) + size, bytes.begin());
  return std::move(bytes);
}

void DMACommand::set_bytes(Address mram_address, Address size,
                           std::vector<int> bytes) {
  assert(size == bytes.size());

  std::copy(bytes.begin(), bytes.end(), bytes_.begin() + index(mram_address));
}

void DMACommand::ack_bytes(Address mram_address, Address size) {
  for (int i = 0; i < size; i++) {
    assert(not acks_[index(mram_address) + i]);
    acks_[index(mram_address) + i] = true;
  }
}

int DMACommand::index(Address mram_address) {
  assert(this->mram_address() <= mram_address and
         mram_address <= this->mram_address() + size_);

  return static_cast<int>(mram_address - this->mram_address());
}

}  // namespace upmem_sim::simulator::dpu
