#include "converter/flag_converter.h"

#include <stdexcept>

namespace upmem_sim::converter {

std::string FlagConverter::to_string(abi::isa::Flag flag) {
  if (flag == abi::isa::ZERO) {
    return "zero";
  } else if (flag == abi::isa::CARRY) {
    return "carry";
  } else {
    throw std::invalid_argument("");
  }
}

}  // namespace upmem_sim::converter