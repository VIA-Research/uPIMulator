#ifndef UPMEM_SIM_CONVERTER_FLAG_CONVERTER_H_
#define UPMEM_SIM_CONVERTER_FLAG_CONVERTER_H_

#include <string>

#include "abi/isa/flag.h"

namespace upmem_sim::converter {

class FlagConverter {
 public:
  static std::string to_string(abi::isa::Flag flag);
};

}  // namespace upmem_sim::converter

#endif
