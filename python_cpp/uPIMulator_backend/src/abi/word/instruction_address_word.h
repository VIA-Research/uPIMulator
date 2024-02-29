#ifndef UPMEM_SIM_ABI_WORD_INSTRUCTION_ADDRESS_WORD_H_
#define UPMEM_SIM_ABI_WORD_INSTRUCTION_ADDRESS_WORD_H_

#include "abi/word/_base_word.h"
#include "util/config_loader.h"

namespace upmem_sim::abi::word {

class InstructionAddressWord : public _BaseWord {
 public:
  InstructionAddressWord()
      : _BaseWord(util::ConfigLoader::iram_address_width()) {}
  ~InstructionAddressWord() = default;

  Address address() { return value(UNSIGNED); }
};

}  // namespace upmem_sim::abi::word

#endif
