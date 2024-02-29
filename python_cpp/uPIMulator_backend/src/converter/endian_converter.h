#ifndef UPMEM_SIM_CONVERTER_ENDIAN_CONVERTER_H_
#define UPMEM_SIM_CONVERTER_ENDIAN_CONVERTER_H_

#include <string>

#include "abi/isa/endian.h"

namespace upmem_sim::converter {

class EndianConverter {
 public:
  static std::string to_string(abi::isa::Endian endian);
};

}  // namespace upmem_sim::converter

#endif
