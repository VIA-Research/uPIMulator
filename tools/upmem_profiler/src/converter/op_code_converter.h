#ifndef UPMEM_PROFILER_CONVERTER_OP_CODE_CONVERTER_H_
#define UPMEM_PROFILER_CONVERTER_OP_CODE_CONVERTER_H_

#include <string>

#include "abi/instruction/op_code.h"

namespace upmem_profiler::converter {

class OpCodeConverter {
public:
  static abi::instruction::OpCode to_op_code(std::string op_code);
};

} // namespace upmem_profiler::converter

#endif
