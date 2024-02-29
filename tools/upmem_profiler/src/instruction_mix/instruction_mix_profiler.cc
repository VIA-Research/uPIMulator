#include "instruction_mix/instruction_mix_profiler.h"

#include <cassert>
#include <fstream>
#include <iostream>

#include "basic/instruction_parser.h"

namespace upmem_profiler::instruciton_mix {

InstructionMixProfiler::InstructionMixProfiler(util::ArgumentParser *argument_parser) {
  std::string log_file = argument_parser->get_string_parameter("logpath");
  instructions_.resize(argument_parser->get_int_parameter("num_tasklets"));

  total_inst_cnt_ = 0;

  std::ifstream ifs(log_file);
  std::string line;
  while (std::getline(ifs, line)) {
    if (basic::InstructionParser::is_instruction(line)) {
      ThreadID thread_id = basic::InstructionParser::parse_thread_id(line);
      abi::instruction::OpCode op_code = basic::InstructionParser::parse_op_code(line);
      abi::instruction::Suffix suffix = basic::InstructionParser::parse_suffix(line);

      instructions_[thread_id].push_back({op_code, suffix});
      total_inst_cnt_++;
    }

  }

  register_mix("synchronization", abi::instruction::ACQUIRE, abi::instruction::RICI);
  register_mix("synchronization", abi::instruction::RELEASE, abi::instruction::RICI);

  register_mix("arithmetic", abi::instruction::ADD, abi::instruction::RRI);
  register_mix("arithmetic", abi::instruction::ADD, abi::instruction::RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ADD, abi::instruction::RRICI);
  register_mix("arithmetic", abi::instruction::ADD, abi::instruction::RRIF);
  register_mix("arithmetic", abi::instruction::ADD, abi::instruction::RRR);
  register_mix("arithmetic", abi::instruction::ADD, abi::instruction::RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ADD, abi::instruction::RRRCI);
  //register_mix("arithmetic", abi::instruction::ADD, abi::instruction::); // add:ssi
  //register_mix("arithmetic", abi::instruction::ADD, abi::instruction::); // add:sss
  register_mix("arithmetic", abi::instruction::ADD, abi::instruction::ZRI);
  register_mix("arithmetic", abi::instruction::ADD, abi::instruction::ZRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ADD, abi::instruction::ZRICI);
  register_mix("arithmetic", abi::instruction::ADD, abi::instruction::ZRIF);
  register_mix("arithmetic", abi::instruction::ADD, abi::instruction::ZRR);
  register_mix("arithmetic", abi::instruction::ADD, abi::instruction::ZRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ADD, abi::instruction::ZRRCI);
  register_mix("arithmetic", abi::instruction::ADD, abi::instruction::S_RRI);
  register_mix("arithmetic", abi::instruction::ADD, abi::instruction::S_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ADD, abi::instruction::S_RRICI);
  register_mix("arithmetic", abi::instruction::ADD, abi::instruction::S_RRIF);
  register_mix("arithmetic", abi::instruction::ADD, abi::instruction::S_RRR);
  register_mix("arithmetic", abi::instruction::ADD, abi::instruction::S_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ADD, abi::instruction::S_RRRCI);
  register_mix("arithmetic", abi::instruction::ADD, abi::instruction::U_RRI);
  register_mix("arithmetic", abi::instruction::ADD, abi::instruction::U_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ADD, abi::instruction::U_RRICI);
  register_mix("arithmetic", abi::instruction::ADD, abi::instruction::U_RRIF);
  register_mix("arithmetic", abi::instruction::ADD, abi::instruction::U_RRR);
  register_mix("arithmetic", abi::instruction::ADD, abi::instruction::U_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ADD, abi::instruction::U_RRRCI);

  register_mix("arithmetic", abi::instruction::ADDC, abi::instruction::RRI);
  register_mix("arithmetic", abi::instruction::ADDC, abi::instruction::RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ADDC, abi::instruction::RRICI);
  register_mix("arithmetic", abi::instruction::ADDC, abi::instruction::RRIF);
  register_mix("arithmetic", abi::instruction::ADDC, abi::instruction::RRR);
  register_mix("arithmetic", abi::instruction::ADDC, abi::instruction::RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ADDC, abi::instruction::RRRCI);
  register_mix("arithmetic", abi::instruction::ADDC, abi::instruction::ZRI);
  register_mix("arithmetic", abi::instruction::ADDC, abi::instruction::ZRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ADDC, abi::instruction::ZRICI);
  register_mix("arithmetic", abi::instruction::ADDC, abi::instruction::ZRIF);
  register_mix("arithmetic", abi::instruction::ADDC, abi::instruction::ZRR);
  register_mix("arithmetic", abi::instruction::ADDC, abi::instruction::ZRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ADDC, abi::instruction::ZRRCI);
  register_mix("arithmetic", abi::instruction::ADDC, abi::instruction::S_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ADDC, abi::instruction::S_RRICI);
  register_mix("arithmetic", abi::instruction::ADDC, abi::instruction::S_RRIF);
  register_mix("arithmetic", abi::instruction::ADDC, abi::instruction::S_RRR);
  register_mix("arithmetic", abi::instruction::ADDC, abi::instruction::S_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ADDC, abi::instruction::S_RRRCI);
  register_mix("arithmetic", abi::instruction::ADDC, abi::instruction::U_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ADDC, abi::instruction::U_RRICI);
  register_mix("arithmetic", abi::instruction::ADDC, abi::instruction::U_RRIF);
  register_mix("arithmetic", abi::instruction::ADDC, abi::instruction::U_RRR);
  register_mix("arithmetic", abi::instruction::ADDC, abi::instruction::U_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ADDC, abi::instruction::U_RRRCI);

  register_mix("arithmetic", abi::instruction::AND, abi::instruction::RRI);
  register_mix("arithmetic", abi::instruction::AND, abi::instruction::RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::AND, abi::instruction::RRICI);
  register_mix("arithmetic", abi::instruction::AND, abi::instruction::RRIF);
  register_mix("arithmetic", abi::instruction::AND, abi::instruction::RRR);
  register_mix("arithmetic", abi::instruction::AND, abi::instruction::RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::AND, abi::instruction::RRRCI);
  register_mix("arithmetic", abi::instruction::AND, abi::instruction::ZRI);
  register_mix("arithmetic", abi::instruction::AND, abi::instruction::ZRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::AND, abi::instruction::ZRICI);
  register_mix("arithmetic", abi::instruction::AND, abi::instruction::ZRIF);
  register_mix("arithmetic", abi::instruction::AND, abi::instruction::ZRR);
  register_mix("arithmetic", abi::instruction::AND, abi::instruction::ZRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::AND, abi::instruction::ZRRCI);
  //register_mix("arithmetic", abi::instruction::AND, abi::instruction::); // and.s:rki
  register_mix("arithmetic", abi::instruction::AND, abi::instruction::S_RRI);
  register_mix("arithmetic", abi::instruction::AND, abi::instruction::S_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::AND, abi::instruction::S_RRICI);
  register_mix("arithmetic", abi::instruction::AND, abi::instruction::S_RRIF);
  register_mix("arithmetic", abi::instruction::AND, abi::instruction::S_RRR);
  register_mix("arithmetic", abi::instruction::AND, abi::instruction::S_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::AND, abi::instruction::S_RRRCI);
  //register_mix("arithmetic", abi::instruction::AND, abi::instruction::); // and.u:rki
  register_mix("arithmetic", abi::instruction::AND, abi::instruction::U_RRI);
  register_mix("arithmetic", abi::instruction::AND, abi::instruction::U_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::AND, abi::instruction::U_RRICI);
  register_mix("arithmetic", abi::instruction::AND, abi::instruction::U_RRIF);
  register_mix("arithmetic", abi::instruction::AND, abi::instruction::U_RRR);
  register_mix("arithmetic", abi::instruction::AND, abi::instruction::U_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::AND, abi::instruction::U_RRRCI);

  register_mix("arithmetic", abi::instruction::ANDN, abi::instruction::RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ANDN, abi::instruction::RRICI);
  register_mix("arithmetic", abi::instruction::ANDN, abi::instruction::RRIF);
  register_mix("arithmetic", abi::instruction::ANDN, abi::instruction::RRR);
  register_mix("arithmetic", abi::instruction::ANDN, abi::instruction::RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ANDN, abi::instruction::RRRCI);
  register_mix("arithmetic", abi::instruction::ANDN, abi::instruction::ZRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ANDN, abi::instruction::ZRICI);
  register_mix("arithmetic", abi::instruction::ANDN, abi::instruction::ZRIF);
  register_mix("arithmetic", abi::instruction::ANDN, abi::instruction::ZRR);
  register_mix("arithmetic", abi::instruction::ANDN, abi::instruction::ZRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ANDN, abi::instruction::ZRRCI);
  register_mix("arithmetic", abi::instruction::ANDN, abi::instruction::S_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ANDN, abi::instruction::S_RRICI);
  register_mix("arithmetic", abi::instruction::ANDN, abi::instruction::S_RRIF);
  register_mix("arithmetic", abi::instruction::ANDN, abi::instruction::S_RRR);
  register_mix("arithmetic", abi::instruction::ANDN, abi::instruction::S_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ANDN, abi::instruction::S_RRRCI);
  register_mix("arithmetic", abi::instruction::ANDN, abi::instruction::U_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ANDN, abi::instruction::U_RRICI);
  register_mix("arithmetic", abi::instruction::ANDN, abi::instruction::U_RRIF);
  register_mix("arithmetic", abi::instruction::ANDN, abi::instruction::U_RRR);
  register_mix("arithmetic", abi::instruction::ANDN, abi::instruction::U_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ANDN, abi::instruction::U_RRRCI);

  register_mix("arithmetic", abi::instruction::ASR, abi::instruction::RRI);
  register_mix("arithmetic", abi::instruction::ASR, abi::instruction::RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ASR, abi::instruction::RRICI);
  register_mix("arithmetic", abi::instruction::ASR, abi::instruction::RRR);
  register_mix("arithmetic", abi::instruction::ASR, abi::instruction::RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ASR, abi::instruction::RRRCI);
  register_mix("arithmetic", abi::instruction::ASR, abi::instruction::ZRI);
  register_mix("arithmetic", abi::instruction::ASR, abi::instruction::ZRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ASR, abi::instruction::ZRICI);
  register_mix("arithmetic", abi::instruction::ASR, abi::instruction::ZRR);
  register_mix("arithmetic", abi::instruction::ASR, abi::instruction::ZRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ASR, abi::instruction::ZRRCI);
  register_mix("arithmetic", abi::instruction::ASR, abi::instruction::S_RRI);
  register_mix("arithmetic", abi::instruction::ASR, abi::instruction::S_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ASR, abi::instruction::S_RRICI);
  register_mix("arithmetic", abi::instruction::ASR, abi::instruction::S_RRR);
  register_mix("arithmetic", abi::instruction::ASR, abi::instruction::S_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ASR, abi::instruction::S_RRRCI);
  register_mix("arithmetic", abi::instruction::ASR, abi::instruction::U_RRI);
  register_mix("arithmetic", abi::instruction::ASR, abi::instruction::U_RRIC);
  register_mix("arithmetic", abi::instruction::ASR, abi::instruction::U_RRICI);
  register_mix("arithmetic", abi::instruction::ASR, abi::instruction::U_RRR);
  register_mix("arithmetic", abi::instruction::ASR, abi::instruction::U_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ASR, abi::instruction::U_RRRCI);

  register_mix("arithmetic", abi::instruction::CAO, abi::instruction::RR);
  register_mix("arithmetic", abi::instruction::CAO, abi::instruction::RRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::CAO, abi::instruction::RRCI);
  register_mix("arithmetic", abi::instruction::CAO, abi::instruction::ZR);
  register_mix("arithmetic", abi::instruction::CAO, abi::instruction::ZRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::CAO, abi::instruction::ZRCI);
  register_mix("arithmetic", abi::instruction::CAO, abi::instruction::S_RR);
  register_mix("arithmetic", abi::instruction::CAO, abi::instruction::S_RRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::CAO, abi::instruction::S_RRCI);
  register_mix("arithmetic", abi::instruction::CAO, abi::instruction::U_RR);
  register_mix("arithmetic", abi::instruction::CAO, abi::instruction::U_RRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::CAO, abi::instruction::U_RRCI);

  register_mix("arithmetic", abi::instruction::CLO, abi::instruction::RR);
  register_mix("arithmetic", abi::instruction::CLO, abi::instruction::RRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::CLO, abi::instruction::RRCI);
  register_mix("arithmetic", abi::instruction::CLO, abi::instruction::ZR);
  register_mix("arithmetic", abi::instruction::CLO, abi::instruction::ZRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::CLO, abi::instruction::ZRCI);
  register_mix("arithmetic", abi::instruction::CLO, abi::instruction::S_RR);
  register_mix("arithmetic", abi::instruction::CLO, abi::instruction::S_RRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::CLO, abi::instruction::S_RRCI);
  register_mix("arithmetic", abi::instruction::CLO, abi::instruction::U_RR);
  register_mix("arithmetic", abi::instruction::CLO, abi::instruction::U_RRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::CLO, abi::instruction::U_RRCI);

  register_mix("arithmetic", abi::instruction::CLS, abi::instruction::RR);
  register_mix("arithmetic", abi::instruction::CLS, abi::instruction::RRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::CLS, abi::instruction::RRCI);
  register_mix("arithmetic", abi::instruction::CLS, abi::instruction::ZR);
  register_mix("arithmetic", abi::instruction::CLS, abi::instruction::ZRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::CLS, abi::instruction::ZRCI);
  register_mix("arithmetic", abi::instruction::CLS, abi::instruction::S_RR);
  register_mix("arithmetic", abi::instruction::CLS, abi::instruction::S_RRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::CLS, abi::instruction::S_RRCI);
  register_mix("arithmetic", abi::instruction::CLS, abi::instruction::U_RR);
  register_mix("arithmetic", abi::instruction::CLS, abi::instruction::U_RRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::CLS, abi::instruction::U_RRCI);

  register_mix("arithmetic", abi::instruction::CLZ, abi::instruction::RR);
  register_mix("arithmetic", abi::instruction::CLZ, abi::instruction::RRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::CLZ, abi::instruction::RRCI);
  register_mix("arithmetic", abi::instruction::CLZ, abi::instruction::ZR);
  register_mix("arithmetic", abi::instruction::CLZ, abi::instruction::ZRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::CLZ, abi::instruction::ZRCI);
  register_mix("arithmetic", abi::instruction::CLZ, abi::instruction::S_RR);
  register_mix("arithmetic", abi::instruction::CLZ, abi::instruction::S_RRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::CLZ, abi::instruction::S_RRCI);
  register_mix("arithmetic", abi::instruction::CLZ, abi::instruction::U_RR);
  register_mix("arithmetic", abi::instruction::CLZ, abi::instruction::U_RRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::CLZ, abi::instruction::U_RRCI);

  register_mix("arithmetic", abi::instruction::CMPB4, abi::instruction::RRR);
  register_mix("arithmetic", abi::instruction::CMPB4, abi::instruction::RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::CMPB4, abi::instruction::RRRCI);
  register_mix("arithmetic", abi::instruction::CMPB4, abi::instruction::ZRR);
  register_mix("arithmetic", abi::instruction::CMPB4, abi::instruction::ZRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::CMPB4, abi::instruction::ZRRCI);
  register_mix("arithmetic", abi::instruction::CMPB4, abi::instruction::S_RRR);
  register_mix("arithmetic", abi::instruction::CMPB4, abi::instruction::S_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::CMPB4, abi::instruction::S_RRRCI);
  register_mix("arithmetic", abi::instruction::CMPB4, abi::instruction::U_RRR);
  register_mix("arithmetic", abi::instruction::CMPB4, abi::instruction::U_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::CMPB4, abi::instruction::U_RRRCI);


  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::DIV_STEP, abi::instruction::DRDICI);

  register_mix("arithmetic", abi::instruction::EXTSB, abi::instruction::RR);
  register_mix("arithmetic", abi::instruction::EXTSB, abi::instruction::RRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::EXTSB, abi::instruction::RRCI);
  register_mix("arithmetic", abi::instruction::EXTSB, abi::instruction::ZR);
  register_mix("arithmetic", abi::instruction::EXTSB, abi::instruction::ZRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::EXTSB, abi::instruction::ZRCI);
  register_mix("arithmetic", abi::instruction::EXTSB, abi::instruction::S_RR);
  register_mix("arithmetic", abi::instruction::EXTSB, abi::instruction::S_RRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::EXTSB, abi::instruction::S_RRCI);

  register_mix("arithmetic", abi::instruction::EXTSH, abi::instruction::RR);
  register_mix("arithmetic", abi::instruction::EXTSH, abi::instruction::RRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::EXTSH, abi::instruction::RRCI);
  register_mix("arithmetic", abi::instruction::EXTSH, abi::instruction::ZR);
  register_mix("arithmetic", abi::instruction::EXTSH, abi::instruction::ZRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::EXTSH, abi::instruction::ZRCI);
  register_mix("arithmetic", abi::instruction::EXTSH, abi::instruction::S_RR);
  register_mix("arithmetic", abi::instruction::EXTSH, abi::instruction::S_RRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::EXTSH, abi::instruction::S_RRCI);

  register_mix("arithmetic", abi::instruction::EXTUB, abi::instruction::RR);
  register_mix("arithmetic", abi::instruction::EXTUB, abi::instruction::RRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::EXTUB, abi::instruction::RRCI);
  register_mix("arithmetic", abi::instruction::EXTUB, abi::instruction::ZR);
  register_mix("arithmetic", abi::instruction::EXTUB, abi::instruction::ZRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::EXTUB, abi::instruction::ZRCI);
  register_mix("arithmetic", abi::instruction::EXTUB, abi::instruction::U_RR);
  register_mix("arithmetic", abi::instruction::EXTUB, abi::instruction::U_RRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::EXTUB, abi::instruction::U_RRCI);

  register_mix("arithmetic", abi::instruction::EXTUH, abi::instruction::RR);
  register_mix("arithmetic", abi::instruction::EXTUH, abi::instruction::RRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::EXTUH, abi::instruction::RRCI);
  register_mix("arithmetic", abi::instruction::EXTUH, abi::instruction::ZR);
  register_mix("arithmetic", abi::instruction::EXTUH, abi::instruction::ZRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::EXTUH, abi::instruction::ZRCI);
  register_mix("arithmetic", abi::instruction::EXTUH, abi::instruction::U_RR);
  register_mix("arithmetic", abi::instruction::EXTUH, abi::instruction::U_RRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::EXTUH, abi::instruction::U_RRCI);

  register_mix("arithmetic", abi::instruction::LSL, abi::instruction::RRI);
  register_mix("arithmetic", abi::instruction::LSL, abi::instruction::RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSL, abi::instruction::RRICI);
  register_mix("arithmetic", abi::instruction::LSL, abi::instruction::RRR);
  register_mix("arithmetic", abi::instruction::LSL, abi::instruction::RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSL, abi::instruction::RRRCI);
  register_mix("arithmetic", abi::instruction::LSL, abi::instruction::ZRI);
  register_mix("arithmetic", abi::instruction::LSL, abi::instruction::ZRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSL, abi::instruction::ZRICI);
  register_mix("arithmetic", abi::instruction::LSL, abi::instruction::ZRR);
  register_mix("arithmetic", abi::instruction::LSL, abi::instruction::ZRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSL, abi::instruction::ZRRCI);
  register_mix("arithmetic", abi::instruction::LSL, abi::instruction::S_RRI);
  register_mix("arithmetic", abi::instruction::LSL, abi::instruction::S_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSL, abi::instruction::S_RRICI);
  register_mix("arithmetic", abi::instruction::LSL, abi::instruction::S_RRR);
  register_mix("arithmetic", abi::instruction::LSL, abi::instruction::S_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSL, abi::instruction::S_RRRCI);
  register_mix("arithmetic", abi::instruction::LSL, abi::instruction::U_RRI);
  register_mix("arithmetic", abi::instruction::LSL, abi::instruction::U_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSL, abi::instruction::U_RRICI);
  register_mix("arithmetic", abi::instruction::LSL, abi::instruction::U_RRR);
  register_mix("arithmetic", abi::instruction::LSL, abi::instruction::U_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSL, abi::instruction::U_RRRCI);

  register_mix("arithmetic", abi::instruction::LSL_ADD, abi::instruction::RRRI);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSL_ADD, abi::instruction::RRRICI);
  register_mix("arithmetic", abi::instruction::LSL_ADD, abi::instruction::ZRRI);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSL_ADD, abi::instruction::ZRRICI);
  register_mix("arithmetic", abi::instruction::LSL_ADD, abi::instruction::S_RRRI);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSL_ADD, abi::instruction::S_RRRICI);
  register_mix("arithmetic", abi::instruction::LSL_ADD, abi::instruction::U_RRRI);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSL_ADD, abi::instruction::U_RRRICI);

  register_mix("arithmetic", abi::instruction::LSL_SUB, abi::instruction::RRRI);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSL_SUB, abi::instruction::RRRICI);
  register_mix("arithmetic", abi::instruction::LSL_SUB, abi::instruction::ZRRI);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSL_SUB, abi::instruction::ZRRICI);
  register_mix("arithmetic", abi::instruction::LSL_SUB, abi::instruction::S_RRRI);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSL_SUB, abi::instruction::S_RRRICI);
  register_mix("arithmetic", abi::instruction::LSL_SUB, abi::instruction::U_RRRI);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSL_SUB, abi::instruction::U_RRRICI);

  register_mix("arithmetic", abi::instruction::LSL1, abi::instruction::RRI);
  register_mix("arithmetic", abi::instruction::LSL1, abi::instruction::RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSL1, abi::instruction::RRICI);
  register_mix("arithmetic", abi::instruction::LSL1, abi::instruction::RRR);
  register_mix("arithmetic", abi::instruction::LSL1, abi::instruction::RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSL1, abi::instruction::RRRCI);
  register_mix("arithmetic", abi::instruction::LSL1, abi::instruction::ZRI);
  register_mix("arithmetic", abi::instruction::LSL1, abi::instruction::ZRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSL1, abi::instruction::ZRICI);
  register_mix("arithmetic", abi::instruction::LSL1, abi::instruction::ZRR);
  register_mix("arithmetic", abi::instruction::LSL1, abi::instruction::ZRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSL1, abi::instruction::ZRRCI);
  register_mix("arithmetic", abi::instruction::LSL1, abi::instruction::S_RRI);
  register_mix("arithmetic", abi::instruction::LSL1, abi::instruction::S_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSL1, abi::instruction::S_RRICI);
  register_mix("arithmetic", abi::instruction::LSL1, abi::instruction::S_RRR);
  register_mix("arithmetic", abi::instruction::LSL1, abi::instruction::S_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSL1, abi::instruction::S_RRRCI);
  register_mix("arithmetic", abi::instruction::LSL1, abi::instruction::U_RRI);
  register_mix("arithmetic", abi::instruction::LSL1, abi::instruction::U_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSL1, abi::instruction::U_RRICI);
  register_mix("arithmetic", abi::instruction::LSL1, abi::instruction::U_RRR);
  register_mix("arithmetic", abi::instruction::LSL1, abi::instruction::U_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSL1, abi::instruction::U_RRRCI);

  register_mix("arithmetic", abi::instruction::LSL1X, abi::instruction::RRI);
  register_mix("arithmetic", abi::instruction::LSL1X, abi::instruction::RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSL1X, abi::instruction::RRICI);
  register_mix("arithmetic", abi::instruction::LSL1X, abi::instruction::RRR);
  register_mix("arithmetic", abi::instruction::LSL1X, abi::instruction::RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSL1X, abi::instruction::RRRCI);
  register_mix("arithmetic", abi::instruction::LSL1X, abi::instruction::ZRI);
  register_mix("arithmetic", abi::instruction::LSL1X, abi::instruction::ZRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSL1X, abi::instruction::ZRICI);
  register_mix("arithmetic", abi::instruction::LSL1X, abi::instruction::ZRR);
  register_mix("arithmetic", abi::instruction::LSL1X, abi::instruction::ZRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSL1X, abi::instruction::ZRRCI);
  register_mix("arithmetic", abi::instruction::LSL1X, abi::instruction::S_RRI);
  register_mix("arithmetic", abi::instruction::LSL1X, abi::instruction::S_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSL1X, abi::instruction::S_RRICI);
  register_mix("arithmetic", abi::instruction::LSL1X, abi::instruction::S_RRR);
  register_mix("arithmetic", abi::instruction::LSL1X, abi::instruction::S_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSL1X, abi::instruction::S_RRRCI);
  register_mix("arithmetic", abi::instruction::LSL1X, abi::instruction::U_RRI);
  register_mix("arithmetic", abi::instruction::LSL1X, abi::instruction::U_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSL1X, abi::instruction::U_RRICI);
  register_mix("arithmetic", abi::instruction::LSL1X, abi::instruction::U_RRR);
  register_mix("arithmetic", abi::instruction::LSL1X, abi::instruction::U_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSL1X, abi::instruction::U_RRRCI);

  register_mix("arithmetic", abi::instruction::LSLX, abi::instruction::RRI);
  register_mix("arithmetic", abi::instruction::LSLX, abi::instruction::RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSLX, abi::instruction::RRICI);
  register_mix("arithmetic", abi::instruction::LSLX, abi::instruction::RRR);
  register_mix("arithmetic", abi::instruction::LSLX, abi::instruction::RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSLX, abi::instruction::RRRCI);
  register_mix("arithmetic", abi::instruction::LSLX, abi::instruction::ZRI);
  register_mix("arithmetic", abi::instruction::LSLX, abi::instruction::ZRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSLX, abi::instruction::ZRICI);
  register_mix("arithmetic", abi::instruction::LSLX, abi::instruction::ZRR);
  register_mix("arithmetic", abi::instruction::LSLX, abi::instruction::ZRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSLX, abi::instruction::ZRRCI);
  register_mix("arithmetic", abi::instruction::LSLX, abi::instruction::S_RRI);
  register_mix("arithmetic", abi::instruction::LSLX, abi::instruction::S_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSLX, abi::instruction::S_RRICI);
  register_mix("arithmetic", abi::instruction::LSLX, abi::instruction::S_RRR);
  register_mix("arithmetic", abi::instruction::LSLX, abi::instruction::S_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSLX, abi::instruction::S_RRRCI);
  register_mix("arithmetic", abi::instruction::LSLX, abi::instruction::U_RRI);
  register_mix("arithmetic", abi::instruction::LSLX, abi::instruction::U_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSLX, abi::instruction::U_RRICI);
  register_mix("arithmetic", abi::instruction::LSLX, abi::instruction::U_RRR);
  register_mix("arithmetic", abi::instruction::LSLX, abi::instruction::U_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSLX, abi::instruction::U_RRRCI);

  register_mix("arithmetic", abi::instruction::LSR, abi::instruction::RRI);
  register_mix("arithmetic", abi::instruction::LSR, abi::instruction::RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSR, abi::instruction::RRICI);
  register_mix("arithmetic", abi::instruction::LSR, abi::instruction::RRR);
  register_mix("arithmetic", abi::instruction::LSR, abi::instruction::RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSR, abi::instruction::RRRCI);
  register_mix("arithmetic", abi::instruction::LSR, abi::instruction::ZRI);
  register_mix("arithmetic", abi::instruction::LSR, abi::instruction::ZRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSR, abi::instruction::ZRICI);
  register_mix("arithmetic", abi::instruction::LSR, abi::instruction::ZRR);
  register_mix("arithmetic", abi::instruction::LSR, abi::instruction::ZRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSR, abi::instruction::ZRRCI);
  register_mix("arithmetic", abi::instruction::LSR, abi::instruction::S_RRI);
  register_mix("arithmetic", abi::instruction::LSR, abi::instruction::S_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSR, abi::instruction::S_RRICI);
  register_mix("arithmetic", abi::instruction::LSR, abi::instruction::S_RRR);
  register_mix("arithmetic", abi::instruction::LSR, abi::instruction::S_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSR, abi::instruction::S_RRRCI);
  register_mix("arithmetic", abi::instruction::LSR, abi::instruction::U_RRI);
  register_mix("arithmetic", abi::instruction::LSR, abi::instruction::U_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSR, abi::instruction::U_RRICI);
  register_mix("arithmetic", abi::instruction::LSR, abi::instruction::U_RRR);
  register_mix("arithmetic", abi::instruction::LSR, abi::instruction::U_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSR, abi::instruction::U_RRRCI);

  register_mix("arithmetic", abi::instruction::LSR_ADD, abi::instruction::RRRI);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSR_ADD, abi::instruction::RRRICI);
  register_mix("arithmetic", abi::instruction::LSR_ADD, abi::instruction::ZRRI);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSR_ADD, abi::instruction::ZRRICI);
  register_mix("arithmetic", abi::instruction::LSR_ADD, abi::instruction::S_RRRI);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSR_ADD, abi::instruction::S_RRRICI);
  register_mix("arithmetic", abi::instruction::LSR_ADD, abi::instruction::U_RRRI);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSR_ADD, abi::instruction::U_RRRICI);

  register_mix("arithmetic", abi::instruction::LSR1, abi::instruction::RRI);
  register_mix("arithmetic", abi::instruction::LSR1, abi::instruction::RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSR1, abi::instruction::RRICI);
  register_mix("arithmetic", abi::instruction::LSR1, abi::instruction::RRR);
  register_mix("arithmetic", abi::instruction::LSR1, abi::instruction::RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSR1, abi::instruction::RRRCI);
  register_mix("arithmetic", abi::instruction::LSR1, abi::instruction::ZRI);
  register_mix("arithmetic", abi::instruction::LSR1, abi::instruction::ZRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSR1, abi::instruction::ZRICI);
  register_mix("arithmetic", abi::instruction::LSR1, abi::instruction::ZRR);
  register_mix("arithmetic", abi::instruction::LSR1, abi::instruction::ZRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSR1, abi::instruction::ZRRCI);
  register_mix("arithmetic", abi::instruction::LSR1, abi::instruction::S_RRI);
  register_mix("arithmetic", abi::instruction::LSR1, abi::instruction::S_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSR1, abi::instruction::S_RRICI);
  register_mix("arithmetic", abi::instruction::LSR1, abi::instruction::S_RRR);
  register_mix("arithmetic", abi::instruction::LSR1, abi::instruction::S_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSR1, abi::instruction::S_RRRCI);
  register_mix("arithmetic", abi::instruction::LSR1, abi::instruction::U_RRI);
  register_mix("arithmetic", abi::instruction::LSR1, abi::instruction::U_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSR1, abi::instruction::U_RRICI);
  register_mix("arithmetic", abi::instruction::LSR1, abi::instruction::U_RRR);
  register_mix("arithmetic", abi::instruction::LSR1, abi::instruction::U_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSR1, abi::instruction::U_RRRCI);

  register_mix("arithmetic", abi::instruction::LSR1X, abi::instruction::RRI);
  register_mix("arithmetic", abi::instruction::LSR1X, abi::instruction::RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSR1X, abi::instruction::RRICI);
  register_mix("arithmetic", abi::instruction::LSR1X, abi::instruction::RRR);
  register_mix("arithmetic", abi::instruction::LSR1X, abi::instruction::RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSR1X, abi::instruction::RRRCI);
  register_mix("arithmetic", abi::instruction::LSR1X, abi::instruction::ZRI);
  register_mix("arithmetic", abi::instruction::LSR1X, abi::instruction::ZRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSR1X, abi::instruction::ZRICI);
  register_mix("arithmetic", abi::instruction::LSR1X, abi::instruction::ZRR);
  register_mix("arithmetic", abi::instruction::LSR1X, abi::instruction::ZRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSR1X, abi::instruction::ZRRCI);
  register_mix("arithmetic", abi::instruction::LSR1X, abi::instruction::S_RRI);
  register_mix("arithmetic", abi::instruction::LSR1X, abi::instruction::S_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSR1X, abi::instruction::S_RRICI);
  register_mix("arithmetic", abi::instruction::LSR1X, abi::instruction::S_RRR);
  register_mix("arithmetic", abi::instruction::LSR1X, abi::instruction::S_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSR1X, abi::instruction::S_RRRCI);
  register_mix("arithmetic", abi::instruction::LSR1X, abi::instruction::U_RRI);
  register_mix("arithmetic", abi::instruction::LSR1X, abi::instruction::U_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSR1X, abi::instruction::U_RRICI);
  register_mix("arithmetic", abi::instruction::LSR1X, abi::instruction::U_RRR);
  register_mix("arithmetic", abi::instruction::LSR1X, abi::instruction::U_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSR1X, abi::instruction::U_RRRCI);

  register_mix("arithmetic", abi::instruction::LSRX, abi::instruction::RRI);
  register_mix("arithmetic", abi::instruction::LSRX, abi::instruction::RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSRX, abi::instruction::RRICI);
  register_mix("arithmetic", abi::instruction::LSRX, abi::instruction::RRR);
  register_mix("arithmetic", abi::instruction::LSRX, abi::instruction::RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSRX, abi::instruction::RRRCI);
  register_mix("arithmetic", abi::instruction::LSRX, abi::instruction::ZRI);
  register_mix("arithmetic", abi::instruction::LSRX, abi::instruction::ZRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSRX, abi::instruction::ZRICI);
  register_mix("arithmetic", abi::instruction::LSRX, abi::instruction::ZRR);
  register_mix("arithmetic", abi::instruction::LSRX, abi::instruction::ZRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSRX, abi::instruction::ZRRCI);
  register_mix("arithmetic", abi::instruction::LSRX, abi::instruction::S_RRI);
  register_mix("arithmetic", abi::instruction::LSRX, abi::instruction::S_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSRX, abi::instruction::S_RRICI);
  register_mix("arithmetic", abi::instruction::LSRX, abi::instruction::S_RRR);
  register_mix("arithmetic", abi::instruction::LSRX, abi::instruction::S_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSRX, abi::instruction::S_RRRCI);
  register_mix("arithmetic", abi::instruction::LSRX, abi::instruction::U_RRI);
  register_mix("arithmetic", abi::instruction::LSRX, abi::instruction::U_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSRX, abi::instruction::U_RRICI);
  register_mix("arithmetic", abi::instruction::LSRX, abi::instruction::U_RRR);
  register_mix("arithmetic", abi::instruction::LSRX, abi::instruction::U_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::LSRX, abi::instruction::U_RRRCI);

  register_mix("heavy_arithmetic", abi::instruction::MUL_SH_SH, abi::instruction::RRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SH_SH, abi::instruction::RRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_SH_SH, abi::instruction::RRRCI);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SH_SH, abi::instruction::ZRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SH_SH, abi::instruction::ZRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_SH_SH, abi::instruction::ZRRCI);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SH_SH, abi::instruction::S_RRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SH_SH, abi::instruction::S_RRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_SH_SH, abi::instruction::S_RRRCI);

  register_mix("heavy_arithmetic", abi::instruction::MUL_SH_SL, abi::instruction::RRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SH_SL, abi::instruction::RRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_SH_SL, abi::instruction::RRRCI);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SH_SL, abi::instruction::ZRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SH_SL, abi::instruction::ZRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_SH_SL, abi::instruction::ZRRCI);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SH_SL, abi::instruction::S_RRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SH_SL, abi::instruction::S_RRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_SH_SL, abi::instruction::S_RRRCI);

  register_mix("heavy_arithmetic", abi::instruction::MUL_SH_UH, abi::instruction::RRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SH_UH, abi::instruction::RRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_SH_UH, abi::instruction::RRRCI);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SH_UH, abi::instruction::ZRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SH_UH, abi::instruction::ZRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_SH_UH, abi::instruction::ZRRCI);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SH_UH, abi::instruction::S_RRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SH_UH, abi::instruction::S_RRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_SH_UH, abi::instruction::S_RRRCI);

  register_mix("heavy_arithmetic", abi::instruction::MUL_SH_UL, abi::instruction::RRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SH_UL, abi::instruction::RRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_SH_UL, abi::instruction::RRRCI);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SH_UL, abi::instruction::ZRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SH_UL, abi::instruction::ZRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_SH_UL, abi::instruction::ZRRCI);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SH_UL, abi::instruction::S_RRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SH_UL, abi::instruction::S_RRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_SH_UL, abi::instruction::S_RRRCI);

  register_mix("heavy_arithmetic", abi::instruction::MUL_SL_SH, abi::instruction::RRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SL_SH, abi::instruction::RRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_SL_SH, abi::instruction::RRRCI);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SL_SH, abi::instruction::ZRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SL_SH, abi::instruction::ZRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_SL_SH, abi::instruction::ZRRCI);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SL_SH, abi::instruction::S_RRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SL_SH, abi::instruction::S_RRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_SL_SH, abi::instruction::S_RRRCI);

  register_mix("heavy_arithmetic", abi::instruction::MUL_SL_SL, abi::instruction::RRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SL_SL, abi::instruction::RRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_SL_SL, abi::instruction::RRRCI);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SL_SL, abi::instruction::ZRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SL_SL, abi::instruction::ZRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_SL_SL, abi::instruction::ZRRCI);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SL_SL, abi::instruction::S_RRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SL_SL, abi::instruction::S_RRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_SL_SL, abi::instruction::S_RRRCI);

  register_mix("heavy_arithmetic", abi::instruction::MUL_SL_UH, abi::instruction::RRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SL_UH, abi::instruction::RRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_SL_UH, abi::instruction::RRRCI);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SL_UH, abi::instruction::ZRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SL_UH, abi::instruction::ZRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_SL_UH, abi::instruction::ZRRCI);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SL_UH, abi::instruction::S_RRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SL_UH, abi::instruction::S_RRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_SL_UH, abi::instruction::S_RRRCI);

  register_mix("heavy_arithmetic", abi::instruction::MUL_SL_UL, abi::instruction::RRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SL_UL, abi::instruction::RRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_SL_UL, abi::instruction::RRRCI);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SL_UL, abi::instruction::ZRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SL_UL, abi::instruction::ZRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_SL_UL, abi::instruction::ZRRCI);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SL_UL, abi::instruction::S_RRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_SL_UL, abi::instruction::S_RRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_SL_UL, abi::instruction::S_RRRCI);

  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_STEP, abi::instruction::DRDICI);

  register_mix("heavy_arithmetic", abi::instruction::MUL_UH_UH, abi::instruction::RRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_UH_UH, abi::instruction::RRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_UH_UH, abi::instruction::RRRCI);
  register_mix("heavy_arithmetic", abi::instruction::MUL_UH_UH, abi::instruction::ZRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_UH_UH, abi::instruction::ZRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_UH_UH, abi::instruction::ZRRCI);
  register_mix("heavy_arithmetic", abi::instruction::MUL_UH_UH, abi::instruction::U_RRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_UH_UH, abi::instruction::U_RRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_UH_UH, abi::instruction::U_RRRCI);

  register_mix("heavy_arithmetic", abi::instruction::MUL_UH_UL, abi::instruction::RRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_UH_UL, abi::instruction::RRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_UH_UL, abi::instruction::RRRCI);
  register_mix("heavy_arithmetic", abi::instruction::MUL_UH_UL, abi::instruction::ZRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_UH_UL, abi::instruction::ZRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_UH_UL, abi::instruction::ZRRCI);
  register_mix("heavy_arithmetic", abi::instruction::MUL_UH_UL, abi::instruction::U_RRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_UH_UL, abi::instruction::U_RRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_UH_UL, abi::instruction::U_RRRCI);

  register_mix("heavy_arithmetic", abi::instruction::MUL_UL_UH, abi::instruction::RRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_UL_UH, abi::instruction::RRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_UL_UH, abi::instruction::RRRCI);
  register_mix("heavy_arithmetic", abi::instruction::MUL_UL_UH, abi::instruction::ZRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_UL_UH, abi::instruction::ZRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_UL_UH, abi::instruction::ZRRCI);
  register_mix("heavy_arithmetic", abi::instruction::MUL_UL_UH, abi::instruction::U_RRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_UL_UH, abi::instruction::U_RRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_UL_UH, abi::instruction::U_RRRCI);

  register_mix("heavy_arithmetic", abi::instruction::MUL_UL_UL, abi::instruction::RRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_UL_UL, abi::instruction::RRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_UL_UL, abi::instruction::RRRCI);
  register_mix("heavy_arithmetic", abi::instruction::MUL_UL_UL, abi::instruction::ZRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_UL_UL, abi::instruction::ZRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_UL_UL, abi::instruction::ZRRCI);
  register_mix("heavy_arithmetic", abi::instruction::MUL_UL_UL, abi::instruction::U_RRR);
  register_mix("heavy_arithmetic", abi::instruction::MUL_UL_UL, abi::instruction::U_RRRC);
  register_mix("heavy_arithmetic_and_cond_branch", abi::instruction::MUL_UL_UL, abi::instruction::U_RRRCI);

  register_mix("arithmetic", abi::instruction::NAND, abi::instruction::RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::NAND, abi::instruction::RRICI);
  register_mix("arithmetic", abi::instruction::NAND, abi::instruction::RRIF);
  register_mix("arithmetic", abi::instruction::NAND, abi::instruction::RRR);
  register_mix("arithmetic", abi::instruction::NAND, abi::instruction::RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::NAND, abi::instruction::RRRCI);
  register_mix("arithmetic", abi::instruction::NAND, abi::instruction::ZRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::NAND, abi::instruction::ZRICI);
  register_mix("arithmetic", abi::instruction::NAND, abi::instruction::ZRIF);
  register_mix("arithmetic", abi::instruction::NAND, abi::instruction::ZRR);
  register_mix("arithmetic", abi::instruction::NAND, abi::instruction::ZRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::NAND, abi::instruction::ZRRCI);
  register_mix("arithmetic", abi::instruction::NAND, abi::instruction::S_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::NAND, abi::instruction::S_RRICI);
  register_mix("arithmetic", abi::instruction::NAND, abi::instruction::S_RRIF);
  register_mix("arithmetic", abi::instruction::NAND, abi::instruction::S_RRR);
  register_mix("arithmetic", abi::instruction::NAND, abi::instruction::S_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::NAND, abi::instruction::S_RRRCI);
  register_mix("arithmetic", abi::instruction::NAND, abi::instruction::U_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::NAND, abi::instruction::U_RRICI);
  register_mix("arithmetic", abi::instruction::NAND, abi::instruction::U_RRIF);
  register_mix("arithmetic", abi::instruction::NAND, abi::instruction::U_RRR);
  register_mix("arithmetic", abi::instruction::NAND, abi::instruction::U_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::NAND, abi::instruction::U_RRRCI);

  register_mix("arithmetic", abi::instruction::NOR, abi::instruction::RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::NOR, abi::instruction::RRICI);
  register_mix("arithmetic", abi::instruction::NOR, abi::instruction::RRIF);
  register_mix("arithmetic", abi::instruction::NOR, abi::instruction::RRR);
  register_mix("arithmetic", abi::instruction::NOR, abi::instruction::RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::NOR, abi::instruction::RRRCI);
  register_mix("arithmetic", abi::instruction::NOR, abi::instruction::ZRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::NOR, abi::instruction::ZRICI);
  register_mix("arithmetic", abi::instruction::NOR, abi::instruction::ZRIF);
  register_mix("arithmetic", abi::instruction::NOR, abi::instruction::ZRR);
  register_mix("arithmetic", abi::instruction::NOR, abi::instruction::ZRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::NOR, abi::instruction::ZRRCI);
  register_mix("arithmetic", abi::instruction::NOR, abi::instruction::S_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::NOR, abi::instruction::S_RRICI);
  register_mix("arithmetic", abi::instruction::NOR, abi::instruction::S_RRIF);
  register_mix("arithmetic", abi::instruction::NOR, abi::instruction::S_RRR);
  register_mix("arithmetic", abi::instruction::NOR, abi::instruction::S_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::NOR, abi::instruction::S_RRRCI);
  register_mix("arithmetic", abi::instruction::NOR, abi::instruction::U_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::NOR, abi::instruction::U_RRICI);
  register_mix("arithmetic", abi::instruction::NOR, abi::instruction::U_RRIF);
  register_mix("arithmetic", abi::instruction::NOR, abi::instruction::U_RRR);
  register_mix("arithmetic", abi::instruction::NOR, abi::instruction::U_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::NOR, abi::instruction::U_RRRCI);

  register_mix("arithmetic", abi::instruction::NXOR, abi::instruction::RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::NXOR, abi::instruction::RRICI);
  register_mix("arithmetic", abi::instruction::NXOR, abi::instruction::RRIF);
  register_mix("arithmetic", abi::instruction::NXOR, abi::instruction::RRR);
  register_mix("arithmetic", abi::instruction::NXOR, abi::instruction::RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::NXOR, abi::instruction::RRRCI);
  register_mix("arithmetic", abi::instruction::NXOR, abi::instruction::ZRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::NXOR, abi::instruction::ZRICI);
  register_mix("arithmetic", abi::instruction::NXOR, abi::instruction::ZRIF);
  register_mix("arithmetic", abi::instruction::NXOR, abi::instruction::ZRR);
  register_mix("arithmetic", abi::instruction::NXOR, abi::instruction::ZRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::NXOR, abi::instruction::ZRRCI);
  register_mix("arithmetic", abi::instruction::NXOR, abi::instruction::S_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::NXOR, abi::instruction::S_RRICI);
  register_mix("arithmetic", abi::instruction::NXOR, abi::instruction::S_RRIF);
  register_mix("arithmetic", abi::instruction::NXOR, abi::instruction::S_RRR);
  register_mix("arithmetic", abi::instruction::NXOR, abi::instruction::S_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::NXOR, abi::instruction::S_RRRCI);
  register_mix("arithmetic", abi::instruction::NXOR, abi::instruction::U_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::NXOR, abi::instruction::U_RRICI);
  register_mix("arithmetic", abi::instruction::NXOR, abi::instruction::U_RRIF);
  register_mix("arithmetic", abi::instruction::NXOR, abi::instruction::U_RRR);
  register_mix("arithmetic", abi::instruction::NXOR, abi::instruction::U_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::NXOR, abi::instruction::U_RRRCI);

  register_mix("arithmetic", abi::instruction::OR, abi::instruction::RRI);
  register_mix("arithmetic", abi::instruction::OR, abi::instruction::RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::OR, abi::instruction::RRICI);
  register_mix("arithmetic", abi::instruction::OR, abi::instruction::RRIF);
  register_mix("arithmetic", abi::instruction::OR, abi::instruction::RRR);
  register_mix("arithmetic", abi::instruction::OR, abi::instruction::RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::OR, abi::instruction::RRRCI);
  register_mix("arithmetic", abi::instruction::OR, abi::instruction::ZRI);
  register_mix("arithmetic", abi::instruction::OR, abi::instruction::ZRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::OR, abi::instruction::ZRICI);
  register_mix("arithmetic", abi::instruction::OR, abi::instruction::ZRIF);
  register_mix("arithmetic", abi::instruction::OR, abi::instruction::ZRR);
  register_mix("arithmetic", abi::instruction::OR, abi::instruction::ZRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::OR, abi::instruction::ZRRCI);
  register_mix("arithmetic", abi::instruction::OR, abi::instruction::S_RRI);
  register_mix("arithmetic", abi::instruction::OR, abi::instruction::S_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::OR, abi::instruction::S_RRICI);
  register_mix("arithmetic", abi::instruction::OR, abi::instruction::S_RRIF);
  register_mix("arithmetic", abi::instruction::OR, abi::instruction::S_RRR);
  register_mix("arithmetic", abi::instruction::OR, abi::instruction::S_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::OR, abi::instruction::S_RRRCI);
  register_mix("arithmetic", abi::instruction::OR, abi::instruction::U_RRI);
  register_mix("arithmetic", abi::instruction::OR, abi::instruction::U_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::OR, abi::instruction::U_RRICI);
  register_mix("arithmetic", abi::instruction::OR, abi::instruction::U_RRIF);
  register_mix("arithmetic", abi::instruction::OR, abi::instruction::U_RRR);
  register_mix("arithmetic", abi::instruction::OR, abi::instruction::U_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::OR, abi::instruction::U_RRRCI);

  register_mix("arithmetic", abi::instruction::ORN, abi::instruction::RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ORN, abi::instruction::RRICI);
  register_mix("arithmetic", abi::instruction::ORN, abi::instruction::RRIF);
  register_mix("arithmetic", abi::instruction::ORN, abi::instruction::RRR);
  register_mix("arithmetic", abi::instruction::ORN, abi::instruction::RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ORN, abi::instruction::RRRCI);
  register_mix("arithmetic", abi::instruction::ORN, abi::instruction::ZRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ORN, abi::instruction::ZRICI);
  register_mix("arithmetic", abi::instruction::ORN, abi::instruction::ZRIF);
  register_mix("arithmetic", abi::instruction::ORN, abi::instruction::ZRR);
  register_mix("arithmetic", abi::instruction::ORN, abi::instruction::ZRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ORN, abi::instruction::ZRRCI);
  register_mix("arithmetic", abi::instruction::ORN, abi::instruction::S_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ORN, abi::instruction::S_RRICI);
  register_mix("arithmetic", abi::instruction::ORN, abi::instruction::S_RRIF);
  register_mix("arithmetic", abi::instruction::ORN, abi::instruction::S_RRR);
  register_mix("arithmetic", abi::instruction::ORN, abi::instruction::S_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ORN, abi::instruction::S_RRRCI);
  register_mix("arithmetic", abi::instruction::ORN, abi::instruction::U_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ORN, abi::instruction::U_RRICI);
  register_mix("arithmetic", abi::instruction::ORN, abi::instruction::U_RRIF);
  register_mix("arithmetic", abi::instruction::ORN, abi::instruction::U_RRR);
  register_mix("arithmetic", abi::instruction::ORN, abi::instruction::U_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ORN, abi::instruction::U_RRRCI);

  register_mix("arithmetic", abi::instruction::ROL, abi::instruction::RRI);
  register_mix("arithmetic", abi::instruction::ROL, abi::instruction::RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ROL, abi::instruction::RRICI);
  register_mix("arithmetic", abi::instruction::ROL, abi::instruction::RRR);
  register_mix("arithmetic", abi::instruction::ROL, abi::instruction::RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ROL, abi::instruction::RRRCI);
  register_mix("arithmetic", abi::instruction::ROL, abi::instruction::ZRI);
  register_mix("arithmetic", abi::instruction::ROL, abi::instruction::ZRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ROL, abi::instruction::ZRICI);
  register_mix("arithmetic", abi::instruction::ROL, abi::instruction::ZRR);
  register_mix("arithmetic", abi::instruction::ROL, abi::instruction::ZRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ROL, abi::instruction::ZRRCI);
  register_mix("arithmetic", abi::instruction::ROL, abi::instruction::S_RRI);
  register_mix("arithmetic", abi::instruction::ROL, abi::instruction::S_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ROL, abi::instruction::S_RRICI);
  register_mix("arithmetic", abi::instruction::ROL, abi::instruction::S_RRR);
  register_mix("arithmetic", abi::instruction::ROL, abi::instruction::S_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ROL, abi::instruction::S_RRRCI);
  register_mix("arithmetic", abi::instruction::ROL, abi::instruction::U_RRI);
  register_mix("arithmetic", abi::instruction::ROL, abi::instruction::U_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ROL, abi::instruction::U_RRICI);
  register_mix("arithmetic", abi::instruction::ROL, abi::instruction::U_RRR);
  register_mix("arithmetic", abi::instruction::ROL, abi::instruction::U_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ROL, abi::instruction::U_RRRCI);

  register_mix("arithmetic", abi::instruction::ROL_ADD, abi::instruction::RRRI);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ROL_ADD, abi::instruction::RRRICI);
  register_mix("arithmetic", abi::instruction::ROL_ADD, abi::instruction::ZRRI);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ROL_ADD, abi::instruction::ZRRICI);
  register_mix("arithmetic", abi::instruction::ROL_ADD, abi::instruction::S_RRRI);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ROL_ADD, abi::instruction::S_RRRICI);
  register_mix("arithmetic", abi::instruction::ROL_ADD, abi::instruction::U_RRRI);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ROL_ADD, abi::instruction::U_RRRICI);

  register_mix("arithmetic", abi::instruction::ROR, abi::instruction::RRI);
  register_mix("arithmetic", abi::instruction::ROR, abi::instruction::RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ROR, abi::instruction::RRICI);
  register_mix("arithmetic", abi::instruction::ROR, abi::instruction::RRR);
  register_mix("arithmetic", abi::instruction::ROR, abi::instruction::RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ROR, abi::instruction::RRRCI);
  register_mix("arithmetic", abi::instruction::ROR, abi::instruction::ZRI);
  register_mix("arithmetic", abi::instruction::ROR, abi::instruction::ZRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ROR, abi::instruction::ZRICI);
  register_mix("arithmetic", abi::instruction::ROR, abi::instruction::ZRR);
  register_mix("arithmetic", abi::instruction::ROR, abi::instruction::ZRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ROR, abi::instruction::ZRRCI);
  register_mix("arithmetic", abi::instruction::ROR, abi::instruction::S_RRI);
  register_mix("arithmetic", abi::instruction::ROR, abi::instruction::S_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ROR, abi::instruction::S_RRICI);
  register_mix("arithmetic", abi::instruction::ROR, abi::instruction::S_RRR);
  register_mix("arithmetic", abi::instruction::ROR, abi::instruction::S_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ROR, abi::instruction::S_RRRCI);
  register_mix("arithmetic", abi::instruction::ROR, abi::instruction::U_RRI);
  register_mix("arithmetic", abi::instruction::ROR, abi::instruction::U_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ROR, abi::instruction::U_RRICI);
  register_mix("arithmetic", abi::instruction::ROR, abi::instruction::U_RRR);
  register_mix("arithmetic", abi::instruction::ROR, abi::instruction::U_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::ROR, abi::instruction::U_RRRCI);

  register_mix("arithmetic", abi::instruction::RSUB, abi::instruction::RRR);
  register_mix("arithmetic", abi::instruction::RSUB, abi::instruction::RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::RSUB, abi::instruction::RRRCI);
  register_mix("arithmetic", abi::instruction::RSUB, abi::instruction::ZRR);
  register_mix("arithmetic", abi::instruction::RSUB, abi::instruction::ZRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::RSUB, abi::instruction::ZRRCI);
  register_mix("arithmetic", abi::instruction::RSUB, abi::instruction::S_RRR);
  register_mix("arithmetic", abi::instruction::RSUB, abi::instruction::S_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::RSUB, abi::instruction::S_RRRCI);
  register_mix("arithmetic", abi::instruction::RSUB, abi::instruction::U_RRR);
  register_mix("arithmetic", abi::instruction::RSUB, abi::instruction::U_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::RSUB, abi::instruction::U_RRRCI);

  register_mix("arithmetic", abi::instruction::RSUBC, abi::instruction::RRR);
  register_mix("arithmetic", abi::instruction::RSUBC, abi::instruction::RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::RSUBC, abi::instruction::RRRCI);
  register_mix("arithmetic", abi::instruction::RSUBC, abi::instruction::ZRR);
  register_mix("arithmetic", abi::instruction::RSUBC, abi::instruction::ZRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::RSUBC, abi::instruction::ZRRCI);
  register_mix("arithmetic", abi::instruction::RSUBC, abi::instruction::S_RRR);
  register_mix("arithmetic", abi::instruction::RSUBC, abi::instruction::S_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::RSUBC, abi::instruction::S_RRRCI);
  register_mix("arithmetic", abi::instruction::RSUBC, abi::instruction::U_RRR);
  register_mix("arithmetic", abi::instruction::RSUBC, abi::instruction::U_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::RSUBC, abi::instruction::U_RRRCI);

  register_mix("arithmetic", abi::instruction::SUB, abi::instruction::RIR);
  register_mix("arithmetic", abi::instruction::SUB, abi::instruction::RIRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::SUB, abi::instruction::RIRCI);
  //register_mix("arithmetic", abi::instruction::SUB, abi::instruction::); //sub:rirf
  register_mix("arithmetic", abi::instruction::SUB, abi::instruction::RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::SUB, abi::instruction::RRICI);
  register_mix("arithmetic", abi::instruction::SUB, abi::instruction::RRIF);
  register_mix("arithmetic", abi::instruction::SUB, abi::instruction::RRR);
  register_mix("arithmetic", abi::instruction::SUB, abi::instruction::RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::SUB, abi::instruction::RRRCI);
  //register_mix("arithmetic", abi::instruction::SUB, abi::instruction::); // sub:ssi
  //register_mix("arithmetic", abi::instruction::SUB, abi::instruction::); // sub:sss
  register_mix("arithmetic", abi::instruction::SUB, abi::instruction::ZIR);
  register_mix("arithmetic", abi::instruction::SUB, abi::instruction::ZIRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::SUB, abi::instruction::ZIRCI);
  //register_mix("arithmetic", abi::instruction::SUB, abi::instruction::); // sub:zirf
  register_mix("arithmetic", abi::instruction::SUB, abi::instruction::ZRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::SUB, abi::instruction::ZRICI);
  register_mix("arithmetic", abi::instruction::SUB, abi::instruction::ZRIF);
  register_mix("arithmetic", abi::instruction::SUB, abi::instruction::ZRR);
  register_mix("arithmetic", abi::instruction::SUB, abi::instruction::ZRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::SUB, abi::instruction::ZRRCI);
  register_mix("arithmetic", abi::instruction::SUB, abi::instruction::S_RIRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::SUB, abi::instruction::S_RIRCI);
  //register_mix("arithmetic", abi::instruction::SUB, abi::instruction::); // sub.s:rirf
  register_mix("arithmetic", abi::instruction::SUB, abi::instruction::S_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::SUB, abi::instruction::S_RRICI);
  register_mix("arithmetic", abi::instruction::SUB, abi::instruction::S_RRIF);
  register_mix("arithmetic", abi::instruction::SUB, abi::instruction::S_RRR);
  register_mix("arithmetic", abi::instruction::SUB, abi::instruction::S_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::SUB, abi::instruction::S_RRRCI);
  register_mix("arithmetic", abi::instruction::SUB, abi::instruction::U_RIRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::SUB, abi::instruction::U_RIRCI);
  //register_mix("arithmetic", abi::instruction::SUB, abi::instruction::U_RIRF); // sub.u:rirf
  register_mix("arithmetic", abi::instruction::SUB, abi::instruction::U_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::SUB, abi::instruction::U_RRICI);
  register_mix("arithmetic", abi::instruction::SUB, abi::instruction::U_RRIF);
  register_mix("arithmetic", abi::instruction::SUB, abi::instruction::U_RRR);
  register_mix("arithmetic", abi::instruction::SUB, abi::instruction::U_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::SUB, abi::instruction::U_RRRCI);

  register_mix("arithmetic", abi::instruction::SUBC, abi::instruction::RIR);
  register_mix("arithmetic", abi::instruction::SUBC, abi::instruction::RIRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::SUBC, abi::instruction::RIRCI);
  //register_mix("arithmetic", abi::instruction::SUBC, abi::instruction::); //sub:rirf
  register_mix("arithmetic", abi::instruction::SUBC, abi::instruction::RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::SUBC, abi::instruction::RRICI);
  register_mix("arithmetic", abi::instruction::SUBC, abi::instruction::RRIF);
  register_mix("arithmetic", abi::instruction::SUBC, abi::instruction::RRR);
  register_mix("arithmetic", abi::instruction::SUBC, abi::instruction::RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::SUBC, abi::instruction::RRRCI);
  register_mix("arithmetic", abi::instruction::SUBC, abi::instruction::ZIR);
  register_mix("arithmetic", abi::instruction::SUBC, abi::instruction::ZIRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::SUBC, abi::instruction::ZIRCI);
  //register_mix("arithmetic", abi::instruction::SUBC, abi::instruction::); // sub:zirf
  register_mix("arithmetic", abi::instruction::SUBC, abi::instruction::ZRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::SUBC, abi::instruction::ZRICI);
  register_mix("arithmetic", abi::instruction::SUBC, abi::instruction::ZRIF);
  register_mix("arithmetic", abi::instruction::SUBC, abi::instruction::ZRR);
  register_mix("arithmetic", abi::instruction::SUBC, abi::instruction::ZRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::SUBC, abi::instruction::ZRRCI);
  register_mix("arithmetic", abi::instruction::SUBC, abi::instruction::S_RIRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::SUBC, abi::instruction::S_RIRCI);
  //register_mix("arithmetic", abi::instruction::SUBC, abi::instruction::); // sub.s:rirf
  register_mix("arithmetic", abi::instruction::SUBC, abi::instruction::S_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::SUBC, abi::instruction::S_RRICI);
  register_mix("arithmetic", abi::instruction::SUBC, abi::instruction::S_RRIF);
  register_mix("arithmetic", abi::instruction::SUBC, abi::instruction::S_RRR);
  register_mix("arithmetic", abi::instruction::SUBC, abi::instruction::S_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::SUBC, abi::instruction::S_RRRCI);
  register_mix("arithmetic", abi::instruction::SUBC, abi::instruction::U_RIRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::SUBC, abi::instruction::U_RIRCI);
  //register_mix("arithmetic", abi::instruction::SUBC, abi::instruction::U_RIRF); // sub.u:rirf
  register_mix("arithmetic", abi::instruction::SUBC, abi::instruction::U_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::SUBC, abi::instruction::U_RRICI);
  register_mix("arithmetic", abi::instruction::SUBC, abi::instruction::U_RRIF);
  register_mix("arithmetic", abi::instruction::SUBC, abi::instruction::U_RRR);
  register_mix("arithmetic", abi::instruction::SUBC, abi::instruction::U_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::SUBC, abi::instruction::U_RRRCI);

  register_mix("arithmetic", abi::instruction::XOR, abi::instruction::RRI);
  register_mix("arithmetic", abi::instruction::XOR, abi::instruction::RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::XOR, abi::instruction::RRICI);
  register_mix("arithmetic", abi::instruction::XOR, abi::instruction::RRIF);
  register_mix("arithmetic", abi::instruction::XOR, abi::instruction::RRR);
  register_mix("arithmetic", abi::instruction::XOR, abi::instruction::RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::XOR, abi::instruction::RRRCI);
  register_mix("arithmetic", abi::instruction::XOR, abi::instruction::ZRI);
  register_mix("arithmetic", abi::instruction::XOR, abi::instruction::ZRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::XOR, abi::instruction::ZRICI);
  register_mix("arithmetic", abi::instruction::XOR, abi::instruction::ZRIF);
  register_mix("arithmetic", abi::instruction::XOR, abi::instruction::ZRR);
  register_mix("arithmetic", abi::instruction::XOR, abi::instruction::ZRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::XOR, abi::instruction::ZRRCI);
  register_mix("arithmetic", abi::instruction::XOR, abi::instruction::S_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::XOR, abi::instruction::S_RRICI);
  register_mix("arithmetic", abi::instruction::XOR, abi::instruction::S_RRIF);
  register_mix("arithmetic", abi::instruction::XOR, abi::instruction::S_RRR);
  register_mix("arithmetic", abi::instruction::XOR, abi::instruction::S_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::XOR, abi::instruction::S_RRRCI);
  register_mix("arithmetic", abi::instruction::XOR, abi::instruction::U_RRIC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::XOR, abi::instruction::U_RRICI);
  register_mix("arithmetic", abi::instruction::XOR, abi::instruction::U_RRIF);
  register_mix("arithmetic", abi::instruction::XOR, abi::instruction::U_RRR);
  register_mix("arithmetic", abi::instruction::XOR, abi::instruction::U_RRRC);
  register_mix("arithmetic_and_cond_branch", abi::instruction::XOR, abi::instruction::U_RRRCI);

  register_mix("system", abi::instruction::BOOT, abi::instruction::RICI);
  register_mix("system", abi::instruction::RESUME, abi::instruction::RICI);
  register_mix("system", abi::instruction::STOP, abi::instruction::CI);
  register_mix("system", abi::instruction::FAULT, abi::instruction::I);

  register_mix("system", abi::instruction::TIME, abi::instruction::R);
  register_mix("system_and_cond_branch", abi::instruction::TIME, abi::instruction::RCI);
  register_mix("system", abi::instruction::TIME, abi::instruction::Z);
  register_mix("system_and_cond_branch", abi::instruction::TIME, abi::instruction::ZCI);
  register_mix("system", abi::instruction::TIME, abi::instruction::S_R);
  register_mix("system_and_cond_branch", abi::instruction::TIME, abi::instruction::S_RCI);
  register_mix("system", abi::instruction::TIME, abi::instruction::U_R);
  register_mix("system_and_cond_branch", abi::instruction::TIME, abi::instruction::U_RCI);

  register_mix("system", abi::instruction::TIME_CFG, abi::instruction::RR);
  register_mix("system_and_cond_branch", abi::instruction::TIME_CFG, abi::instruction::RRCI);
  register_mix("system", abi::instruction::TIME_CFG, abi::instruction::ZR);
  register_mix("system_and_cond_branch", abi::instruction::TIME_CFG, abi::instruction::ZRCI);
  register_mix("system", abi::instruction::TIME_CFG, abi::instruction::S_RR);
  register_mix("system_and_cond_branch", abi::instruction::TIME_CFG, abi::instruction::S_RRCI);
  register_mix("system", abi::instruction::TIME_CFG, abi::instruction::U_RR);
  register_mix("system_and_cond_branch", abi::instruction::TIME_CFG, abi::instruction::U_RRCI);

  //register_mix("nop", abi::instruction::NOP, "");

  register_mix("call", abi::instruction::CALL, abi::instruction::RRI);
  register_mix("call", abi::instruction::CALL, abi::instruction::RRR);
  register_mix("call", abi::instruction::CALL, abi::instruction::ZRI);
  register_mix("call", abi::instruction::CALL, abi::instruction::ZRR);

  register_mix("sats", abi::instruction::SATS, abi::instruction::RR);
  register_mix("sats", abi::instruction::SATS, abi::instruction::RRC);
  register_mix("sats", abi::instruction::SATS, abi::instruction::RRCI);
  register_mix("sats", abi::instruction::SATS, abi::instruction::ZR);
  register_mix("sats", abi::instruction::SATS, abi::instruction::ZRC);
  register_mix("sats", abi::instruction::SATS, abi::instruction::ZRCI);
  register_mix("sats", abi::instruction::SATS, abi::instruction::S_RR);
  register_mix("sats", abi::instruction::SATS, abi::instruction::S_RRC);
  register_mix("sats", abi::instruction::SATS, abi::instruction::S_RRCI);
  register_mix("sats", abi::instruction::SATS, abi::instruction::U_RR);
  register_mix("sats", abi::instruction::SATS, abi::instruction::U_RRC);
  register_mix("sats", abi::instruction::SATS, abi::instruction::U_RRCI);

  register_mix("hash", abi::instruction::HASH, abi::instruction::RRIC);
  register_mix("hash_and_cond_branch", abi::instruction::HASH, abi::instruction::RRICI);
  register_mix("hash", abi::instruction::HASH, abi::instruction::RRIF);
  register_mix("hash", abi::instruction::HASH, abi::instruction::RRR);
  register_mix("hash", abi::instruction::HASH, abi::instruction::RRRC);
  register_mix("hash_and_cond_branch", abi::instruction::HASH, abi::instruction::RRRCI);
  register_mix("hash", abi::instruction::HASH, abi::instruction::ZRIC);
  register_mix("hash_and_cond_branch", abi::instruction::HASH, abi::instruction::ZRICI);
  register_mix("hash", abi::instruction::HASH, abi::instruction::ZRIF);
  register_mix("hash", abi::instruction::HASH, abi::instruction::ZRR);
  register_mix("hash", abi::instruction::HASH, abi::instruction::ZRRC);
  register_mix("hash_and_cond_branch", abi::instruction::HASH, abi::instruction::ZRRCI);
  register_mix("hash", abi::instruction::HASH, abi::instruction::S_RRIC);
  register_mix("hash_and_cond_branch", abi::instruction::HASH, abi::instruction::S_RRICI);
  register_mix("hash", abi::instruction::HASH, abi::instruction::S_RRIF);
  register_mix("hash", abi::instruction::HASH, abi::instruction::S_RRR);
  register_mix("hash", abi::instruction::HASH, abi::instruction::S_RRRC);
  register_mix("hash_and_cond_branch", abi::instruction::HASH, abi::instruction::S_RRRCI);
  register_mix("hash", abi::instruction::HASH, abi::instruction::U_RRIC);
  register_mix("hash_and_cond_branch", abi::instruction::HASH, abi::instruction::U_RRICI);
  register_mix("hash", abi::instruction::HASH, abi::instruction::U_RRIF);
  register_mix("hash", abi::instruction::HASH, abi::instruction::U_RRR);
  register_mix("hash", abi::instruction::HASH, abi::instruction::U_RRRC);
  register_mix("hash_and_cond_branch", abi::instruction::HASH, abi::instruction::U_RRRCI);
  
  //register_mix("reg_move_and_cond_branch", abi::instruction::MOVD, abi::instruction::RRCI);
  register_mix("reg_move_and_cond_branch", abi::instruction::MOVD, abi::instruction::DDCI);
  register_mix("reg_move_and_cond_branch", abi::instruction::SWAPD, abi::instruction::RRCI);

  register_mix("scratchpad_access", abi::instruction::LBS, abi::instruction::ERRI);
  register_mix("scratchpad_access", abi::instruction::LBS, abi::instruction::EDRI);// lbs:ersi
  register_mix("scratchpad_access", abi::instruction::LBS, abi::instruction::S_ERRI);

  register_mix("scratchpad_access", abi::instruction::LBU, abi::instruction::ERRI);
  register_mix("scratchpad_access", abi::instruction::LBU, abi::instruction::EDRI); // lbu:ersi
  register_mix("scratchpad_access", abi::instruction::LBU, abi::instruction::U_ERRI);

  register_mix("scratchpad_access", abi::instruction::LD, abi::instruction::ERRI);
  register_mix("scratchpad_access", abi::instruction::LD, abi::instruction::EDRI); // ld:ersi

  register_mix("scratchpad_access", abi::instruction::LHS, abi::instruction::ERRI);
  register_mix("scratchpad_access", abi::instruction::LHS, abi::instruction::EDRI); // lhs:ersi
  register_mix("scratchpad_access", abi::instruction::LHS, abi::instruction::S_ERRI);

  register_mix("scratchpad_access", abi::instruction::LHU, abi::instruction::ERRI);
  register_mix("scratchpad_access", abi::instruction::LHU, abi::instruction::EDRI); // lhu:ersi
  register_mix("scratchpad_access", abi::instruction::LHU, abi::instruction::U_ERRI);

  register_mix("scratchpad_access", abi::instruction::LW, abi::instruction::ERRI);
  register_mix("scratchpad_access", abi::instruction::LW, abi::instruction::EDRI); // lw:ersi
  register_mix("scratchpad_access", abi::instruction::LW, abi::instruction::S_ERRI);
  register_mix("scratchpad_access", abi::instruction::LW, abi::instruction::U_ERRI);

  register_mix("scratchpad_access", abi::instruction::SB, abi::instruction::ERII);
  register_mix("scratchpad_access", abi::instruction::SB, abi::instruction::ERIR);
  //register_mix("scratchpad_access", abi::instruction::SB, abi::instruction::ESII); // sb:esii
  register_mix("scratchpad_access", abi::instruction::SB, abi::instruction::ERID); // sb:esir

  register_mix("scratchpad_access", abi::instruction::SB_ID, abi::instruction::ERII);

  register_mix("scratchpad_access", abi::instruction::SD, abi::instruction::ERII);
  register_mix("scratchpad_access", abi::instruction::SD, abi::instruction::ERIR);
  //register_mix("scratchpad_access", abi::instruction::SD, abi::instruction::); // sd:esii
  register_mix("scratchpad_access", abi::instruction::SD, abi::instruction::ERID); // sd:esir

  register_mix("scratchpad_access", abi::instruction::SD_ID, abi::instruction::ERII);

  register_mix("scratchpad_access", abi::instruction::SH, abi::instruction::ERII);
  register_mix("scratchpad_access", abi::instruction::SH, abi::instruction::ERIR);
  //register_mix("scratchpad_access", abi::instruction::SH, abi::instruction::); // sh:esii
  register_mix("scratchpad_access", abi::instruction::SH, abi::instruction::ERID); // sh:esir

  register_mix("scratchpad_access", abi::instruction::SH_ID, abi::instruction::ERII);

  register_mix("scratchpad_access", abi::instruction::SW, abi::instruction::ERII);
  register_mix("scratchpad_access", abi::instruction::SW, abi::instruction::ERIR);
  //register_mix("scratchpad_access", abi::instruction::SW, abi::instruction::); // sw:esii
  register_mix("scratchpad_access", abi::instruction::SW, abi::instruction::ERID); // sw:esir

  register_mix("scratchpad_access", abi::instruction::SW_ID, abi::instruction::ERII);

  register_mix("mainmemory_access", abi::instruction::LDMA, abi::instruction::DMA_RRI);
  register_mix("mainmemory_access", abi::instruction::LDMAI, abi::instruction::DMA_RRI);
  register_mix("mainmemory_access", abi::instruction::SDMA, abi::instruction::DMA_RRI);

}

void InstructionMixProfiler::register_mix(std::string mix, abi::instruction::OpCode op_code,
                                          abi::instruction::Suffix suffix) {
  assert(mix != "etc");

  mixes_[mix].insert({op_code, suffix});
}

void InstructionMixProfiler::profile() {
  std::vector<std::map<std::string, int64_t>> instruction_mixes;
  //auto dbg = instructions_;
  instruction_mixes.resize(instructions_.size());
  for (ThreadID thread_id = 0; thread_id < instructions_.size(); thread_id++) {
    uint64_t non_etc_count = 0;
    for (auto &[op_code, suffix] : instructions_[thread_id]) {
      bool dbg_flag = false;
      for (auto &[mix, mix_specs] : mixes_) {
        for (auto &mix_spec : mix_specs) {
          auto [mix_op_code, mix_suffix] = mix_spec;
          if (op_code == mix_op_code and suffix == mix_suffix) {
            instruction_mixes[thread_id][mix]++;
            non_etc_count++;
            dbg_flag = true;
          }
        }
        //instruction_mixes[thread_id]["etc"] = instructions_[thread_id].size() - non_etc_count;
      }
      if (dbg_flag == false) {
        std::cout << op_code << " " << suffix << std::endl;
      }
    }
    instruction_mixes[thread_id]["etc"] = instructions_[thread_id].size() - non_etc_count;
  }

  // for (ThreadID thread_id = 0; thread_id < instruction_mixes.size(); thread_id++) {
  //   std::cout << "ThreadID: " << thread_id << std::endl;

  //   for (auto &[mix, count] : instruction_mixes[thread_id]) {
  //     std::cout << thread_id  << "_" <<mix << ": " << count << std::endl;
  //   }

  //   std::cout << std::endl;
  // }

  std::cout << "INSTRUCTION_MIX: " << std::endl;

  for (auto &type_ : inst_type_) {
    uint64_t cnt_inst = 0;
    for (ThreadID thread_id = 0; thread_id < instruction_mixes.size(); thread_id++) { 
      cnt_inst += instruction_mixes[thread_id][type_];
    }
    std::cout << type_ << "," << ((double)cnt_inst / (double)total_inst_cnt_) << std::endl;
  }
}

} // namespace upmem_profiler::instruciton_mix
