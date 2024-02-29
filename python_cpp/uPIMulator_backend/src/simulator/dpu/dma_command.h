#ifndef UPMEM_SIM_SIMULATOR_DPU_DMA_COMMAND_H_
#define UPMEM_SIM_SIMULATOR_DPU_DMA_COMMAND_H_

#include <algorithm>

#include "abi/instruction/instruction.h"
#include "abi/word/data_address_word.h"
#include "abi/word/data_word.h"
#include "main.h"

namespace upmem_sim::simulator::dpu {

class DMACommand {
 public:
  enum Operation { READ = 0, WRITE };

  explicit DMACommand(Operation operation, Address mram_address, Address size);
  explicit DMACommand(Operation operation, Address wram_address,
                      Address mram_address, Address size,
                      abi::instruction::Instruction *instruction);
  explicit DMACommand(Operation operation, Address mram_address, Address size,
                      std::vector<int> bytes);
  explicit DMACommand(Operation operation, Address wram_address,
                      Address mram_address, Address size,
                      std::vector<int> bytes,
                      abi::instruction::Instruction *instruction);
  ~DMACommand();

  Operation operation() { return operation_; }
  Address wram_address() { return wram_address_->address(); }
  Address mram_address() { return mram_address_->address(); }
  Address size() { return size_; }
  bool has_instruction() { return instruction_ != nullptr; }
  abi::instruction::Instruction *instruction();

  std::vector<int> bytes();
  std::vector<int> bytes(Address mram_address, Address size);

  void set_bytes(Address mram_address, Address size, std::vector<int> bytes);
  void ack_bytes(Address mram_address, Address size);

  bool is_ready() {
    return std::all_of(acks_.begin(), acks_.end(),
                       [](bool ack) { return ack; });
  }

 protected:
  int index(Address mram_address);

 private:
  Operation operation_;
  abi::word::DataAddressWord *wram_address_;
  abi::word::DataAddressWord *mram_address_;
  Address size_;
  abi::instruction::Instruction *instruction_;
  std::vector<int> bytes_;
  std::vector<bool> acks_;
};

}  // namespace upmem_sim::simulator::dpu

#endif
