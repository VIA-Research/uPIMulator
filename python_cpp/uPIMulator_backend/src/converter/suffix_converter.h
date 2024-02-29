#ifndef UPMEM_SIM_CONVERTER_SUFFIX_CONVERTER_H_
#define UPMEM_SIM_CONVERTER_SUFFIX_CONVERTER_H_

#include <string>

#include "abi/instruction/suffix.h"

namespace upmem_sim::converter {

class SuffixConverter {
 public:
  static std::string to_string(abi::instruction::Suffix suffix);
};

}  // namespace upmem_sim::converter

#endif
