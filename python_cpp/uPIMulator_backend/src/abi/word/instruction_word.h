#ifndef UPMEM_SIM_ABI_WORD_INSTRUCTION_WORD_H_
#define UPMEM_SIM_ABI_WORD_INSTRUCTION_WORD_H_

#include "abi/word/_base_word.h"
#include "util/config_loader.h"

namespace upmem_sim::abi::word {

class InstructionWord : public _BaseWord {
 public:
  InstructionWord() : _BaseWord(util::ConfigLoader::iram_data_width()) {}
  ~InstructionWord() = default;
};

}  // namespace upmem_sim::abi::word

#endif
