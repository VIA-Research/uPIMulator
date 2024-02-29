#ifndef UPMEM_PROFILER_BASIC_INSTRUCTION_PASER_H_
#define UPMEM_PROFILER_BASIC_INSTRUCTION_PASER_H_

#include <map>
#include <string>
#include <tuple>
#include <vector>

#include "abi/instruction/op_code.h"
#include "abi/instruction/suffix.h"
#include "basic/reg_file_parser.h"
#include "main.h"

namespace upmem_profiler::basic {

class InstructionParser {
public:
  static ThreadID parse_thread_id(std::string line);
  static abi::instruction::OpCode parse_op_code(std::string line);
  static abi::instruction::Suffix parse_suffix(std::string line);

  static bool is_instruction(std::string line) {
    return line.find("[") != std::string::npos and line.find("]") != std::string::npos;
  }
  static bool is_call_rri_instruction(std::string line) { return line.find("call, rri") != std::string::npos; }
  static bool is_call_rrr_instruction(std::string line) { return line.find("call, rrr") != std::string::npos; }
  static bool is_return_instruction(std::string line) { return line.find("call, zri, r23, 0") != std::string::npos; }

  static std::tuple<ThreadID, Address> parse_call_rri_instruction(std::string line, RegFile reg_file);
  static std::tuple<ThreadID, Address> parse_call_rrr_instruction(std::string line, RegFile reg_file);
  static ThreadID parse_return_instruction(std::string line);

protected:
  static std::vector<std::string> split_by_comma(std::string line);
  static int64_t lookup_reg_file(std::string reg, ThreadID thread_id, RegFile reg_file);
};

} // namespace upmem_profiler::basic

#endif
