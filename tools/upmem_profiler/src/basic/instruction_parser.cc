#include "basic/instruction_parser.h"

#include <algorithm>
#include <cmath>
#include <iostream>

#include "converter/op_code_converter.h"
#include "converter/suffix_converter.h"

namespace upmem_profiler::basic {

ThreadID InstructionParser::parse_thread_id(std::string line) {
  int open_bracket_pos = line.find("[");
  int close_bracket_pos = line.find("]");

  ThreadID thread_id = std::stoi(line.substr(open_bracket_pos + 1, close_bracket_pos - 1));
  return thread_id;
}

abi::instruction::OpCode InstructionParser::parse_op_code(std::string line) {
  int close_bracket_pos = line.find("]");

  std::vector<std::string> tokens = split_by_comma(line.substr(close_bracket_pos + 1));
  std::string op_code = tokens[0];

  return converter::OpCodeConverter::to_op_code(op_code);
}

abi::instruction::Suffix InstructionParser::parse_suffix(std::string line) {
  std::vector<std::string> tokens = split_by_comma(line);
  std::string suffix = tokens[1];

  return converter::SuffixConverter::to_suffix(suffix);
}

std::tuple<ThreadID, Address> InstructionParser::parse_call_rri_instruction(std::string line, RegFile reg_file) {
  ThreadID thread_id = parse_thread_id(line);

  std::vector<std::string> tokens = split_by_comma(line);
  std::string ra = tokens[3];
  int64_t imm = std::stoi(tokens[4]);

  int64_t ra_value = lookup_reg_file(ra, thread_id, reg_file);

  Address callee_address = ra_value + imm;

  return {thread_id, callee_address};
}

std::tuple<ThreadID, Address> InstructionParser::parse_call_rrr_instruction(std::string line, RegFile reg_file) {
  ThreadID thread_id = parse_thread_id(line);

  std::vector<std::string> tokens = split_by_comma(line);
  std::string ra = tokens[3];
  std::string rb = tokens[4];

  int64_t ra_value = lookup_reg_file(ra, thread_id, reg_file);
  int64_t rb_value = lookup_reg_file(rb, thread_id, reg_file);

  Address callee_address = ra_value + rb_value;

  return {thread_id, callee_address};
}

ThreadID InstructionParser::parse_return_instruction(std::string line) { return parse_thread_id(line); }

std::vector<std::string> InstructionParser::split_by_comma(std::string line) {
  std::vector<std::string> tokens;
  int pos;
  while ((pos = line.find(",")) != std::string::npos) {
    std::string token = line.substr(0, pos);

    if (token.substr(0, 1) == " ") {
      token = token.substr(1);
    }

    line.erase(0, pos + 1);
    tokens.push_back(token);
  }
  tokens.push_back(line);
  return std::move(tokens);
}

int64_t InstructionParser::lookup_reg_file(std::string reg, ThreadID thread_id, RegFile reg_file) {
  if (reg == "zero") {
    return 0;
  } else if (reg == "one") {
    return 1;
  } else if (reg == "lneg") {
    return -1;
  } else if (reg == "mneg") {
    return static_cast<int64_t>(pow(2, 30));
  } else if (reg == "id") {
    return thread_id;
  } else if (reg == "id2") {
    return 2 * thread_id;
  } else if (reg == "id4") {
    return 4 * thread_id;
  } else if (reg == "id8") {
    return 8 * thread_id;
  } else if (reg.substr(0, 1) == "r") {
    RegIndex reg_index = std::stoi(reg.substr(1));
    return reg_file[reg_index];
  } else {
    throw std::invalid_argument("");
  }
}

} // namespace upmem_profiler::basic
