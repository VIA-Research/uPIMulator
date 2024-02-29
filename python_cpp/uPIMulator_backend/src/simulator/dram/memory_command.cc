#include "simulator/dram/memory_command.h"

#include <cassert>

#include "abi/word/data_word.h"
#include "util/config_loader.h"

namespace upmem_sim::simulator::dram {

MemoryCommand::MemoryCommand(Operation operation, Address address)
    : operation_(operation),
      address_(new abi::word::DataAddressWord()),
      size_(0),
      dma_command_(nullptr) {
  assert(operation == ACTIVATION or operation == PRECHARGE);
  assert(address >= util::ConfigLoader::mram_offset());
  assert(address + size_ <=
         util::ConfigLoader::mram_offset() + util::ConfigLoader::mram_size());

  address_->set_value(address);
}

MemoryCommand::MemoryCommand(Operation operation, Address address, Address size,
                             dpu::DMACommand *dma_command)
    : operation_(operation),
      address_(new abi::word::DataAddressWord()),
      size_(size),
      dma_command_(dma_command) {
  assert(operation == READ);
  assert(address >= util::ConfigLoader::mram_offset());
  assert(address + size_ <=
         util::ConfigLoader::mram_offset() + util::ConfigLoader::mram_size());
  assert(dma_command->operation() == dpu::DMACommand::Operation::READ);

  address_->set_value(address);
}

MemoryCommand::MemoryCommand(Operation operation, Address address, Address size,
                             std::vector<int> bytes,
                             dpu::DMACommand *dma_command)
    : operation_(operation),
      address_(new abi::word::DataAddressWord()),
      size_(size),
      bytes_(bytes),
      dma_command_(dma_command) {
  assert(operation == WRITE);
  assert(address >= util::ConfigLoader::mram_offset());
  assert(address + size_ <=
         util::ConfigLoader::mram_offset() + util::ConfigLoader::mram_size());
  assert(dma_command->operation() == dpu::DMACommand::Operation::WRITE);

  address_->set_value(address);
}

MemoryCommand::~MemoryCommand() { delete address_; }

std::vector<int> MemoryCommand::bytes() {
  assert(operation_ == READ or operation_ == WRITE);
  return bytes_;
}

dpu::DMACommand *MemoryCommand::dma_command() {
  assert(dma_command_ != nullptr);
  return dma_command_;
}

}  // namespace upmem_sim::simulator::dram
