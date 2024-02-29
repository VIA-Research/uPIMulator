#include "basic/reg_file_parser.h"

#include <iostream>

#include "util/config_loader.h"

namespace upmem_profiler::basic {

std::map<RegIndex, int64_t> RegFileParser::parse(std::ifstream &ifs) {
  std::map<RegIndex, int64_t> reg_file;
  for (RegIndex index = 0; index < util::ConfigLoader::num_gp_registers(); index++) {
    std::string reg;
    int64_t value;

    ifs >> reg >> value;

    RegIndex reg_index = std::stoi(reg.substr(1, reg.length() - 1));

    reg_file[reg_index] = value;
  }
  return std::move(reg_file);
}

} // namespace upmem_profiler::basic
