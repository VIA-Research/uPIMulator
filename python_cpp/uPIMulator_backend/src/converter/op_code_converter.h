#ifndef UPMEM_SIM_CONVERTER_OP_CODE_CONVERTER_H_
#define UPMEM_SIM_CONVERTER_OP_CODE_CONVERTER_H_

#include <string>

#include "abi/instruction/op_code.h"

namespace upmem_sim::converter {

class OpCodeConverter {
 public:
  static std::string to_string(abi::instruction::OpCode op_code);
};

}  // namespace upmem_sim::converter

#endif
