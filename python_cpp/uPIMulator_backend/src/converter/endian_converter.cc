#include "converter/endian_converter.h"

#include <stdexcept>

namespace upmem_sim::converter {

std::string EndianConverter::to_string(abi::isa::Endian endian) {
  if (endian == abi::isa::LITTLE) {
    return "!little";
  } else if (endian == abi::isa::BIG) {
    return "!big";
  } else {
    throw std::invalid_argument("");
  }
}

}  // namespace upmem_sim::converter