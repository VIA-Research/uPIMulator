#ifndef UPMEM_SIM_ABI_WORD_DATA_WORD_H_
#define UPMEM_SIM_ABI_WORD_DATA_WORD_H_

#include "abi/word/_base_word.h"
#include "util/config_loader.h"

namespace upmem_sim::abi::word {

class DataWord : public _BaseWord {
 public:
  DataWord() : _BaseWord(util::ConfigLoader::mram_data_width()) {
    assert(util::ConfigLoader::atomic_data_width() ==
               util::ConfigLoader::wram_data_width() and
           util::ConfigLoader::atomic_data_width() ==
               util::ConfigLoader::mram_data_width());
  }
  ~DataWord() = default;
};

}  // namespace upmem_sim::abi::word

#endif
