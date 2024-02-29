#ifndef UPMEM_SIM_SIMULATOR_REG_PC_REG_H_
#define UPMEM_SIM_SIMULATOR_REG_PC_REG_H_

#include "abi/word/instruction_address_word.h"
#include "abi/word/instruction_word.h"

namespace upmem_sim::simulator::reg {

class PCReg {
 public:
  explicit PCReg() : word_(new abi::word::InstructionAddressWord()) {}
  ~PCReg() { delete word_; }

  int64_t read() { return word_->address(); }
  void write(int64_t value) { word_->set_value(value); }
  void increment() { write(read() + abi::word::InstructionWord().size()); }
  void cycle() = delete;

 private:
  abi::word::InstructionAddressWord *word_;
};

}  // namespace upmem_sim::simulator::reg

#endif
