#include <iostream>

#include "instruction_mix/instruction_mix_profiler.h"
#include "util/argument_parser.h"

namespace upmem_profiler {

util::ArgumentParser *init_argument_parser() {
  auto argument_parser = new util::ArgumentParser();

  argument_parser->add_option("mode", util::ArgumentParser::STRING, "instruction_mix");

  argument_parser->add_option("labelpath", util::ArgumentParser::STRING, "/home/dongjaelee/upmem_profiler/bin/1024/VA.16/labels.bin");
  argument_parser->add_option("logpath", util::ArgumentParser::STRING, "/home/dongjae/data_sweep_hbm_mmu/trace/ptw1_tlbway16_tlbset1/VA/131072/VA.16.trace");
  argument_parser->add_option("num_tasklets", util::ArgumentParser::INT, "16");

  return argument_parser;
}

} // namespace upmem_profiler

int main(int argc, char **argv) {
  upmem_profiler::util::ArgumentParser *argument_parser = upmem_profiler::init_argument_parser();
  argument_parser->parse(argc, argv);

  std::string mode = argument_parser->get_string_parameter("mode");
  if (mode == "instruction_mix") {
    auto instruction_mix_profiler = new upmem_profiler::instruciton_mix::InstructionMixProfiler(argument_parser);
    instruction_mix_profiler->profile();
    delete instruction_mix_profiler;
  } else {
    throw std::invalid_argument("");
  }

  delete argument_parser;

  return 0;
}
