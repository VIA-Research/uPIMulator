#ifndef UPMEM_SIM_ABI_WORD_DATA_ADDRESS_WORD_H_
#define UPMEM_SIM_ABI_WORD_DATA_ADDRESS_WORD_H_

#include "abi/word/_base_word.h"
#include "util/config_loader.h"

namespace upmem_sim::abi::word {

class DataAddressWord : public _BaseWord {
 public:
  DataAddressWord() : _BaseWord(util::ConfigLoader::mram_address_width()) {
    assert(util::ConfigLoader::atomic_address_width() ==
               util::ConfigLoader::wram_address_width() and
           util::ConfigLoader::atomic_address_width() ==
               util::ConfigLoader::mram_address_width());
  }
  ~DataAddressWord() = default;

  Address address() { return value(UNSIGNED); }
};

}  // namespace upmem_sim::abi::word

#endif
