#ifndef UPMEM_SIM_SIMULATOR_DRAM_MEMORY_COMMAND_H_
#define UPMEM_SIM_SIMULATOR_DRAM_MEMORY_COMMAND_H_

#include "abi/word/data_address_word.h"
#include "abi/word/data_word.h"
#include "simulator/dpu/dma_command.h"

namespace upmem_sim::simulator::dram {

class MemoryCommand {
 public:
  enum Operation {
    ACTIVATION = 0,
    READ,
    WRITE,
    PRECHARGE,
  };

  explicit MemoryCommand(Operation operation, Address address);
  explicit MemoryCommand(Operation operation, Address address, Address size,
                         dpu::DMACommand *dma_command);
  explicit MemoryCommand(Operation operation, Address address, Address size,
                         std::vector<int> bytes, dpu::DMACommand *dma_command);
  ~MemoryCommand();

  Operation operation() { return operation_; }
  Address address() { return address_->address(); }
  Address size() { return size_; }
  std::vector<int> bytes();
  void set_bytes(std::vector<int> bytes) { bytes_ = bytes; }
  dpu::DMACommand *dma_command();

 private:
  Operation operation_;
  abi::word::DataAddressWord *address_;
  Address size_;
  std::vector<int> bytes_;
  dpu::DMACommand *dma_command_;
};

}  // namespace upmem_sim::simulator::dram

#endif
