#ifndef UPMEM_PROFILER_BASIC_REG_FILE_PARSER_H_
#define UPMEM_PROFILER_BASIC_REG_FILE_PARSER_H_

#include <fstream>
#include <map>

#include "main.h"

namespace upmem_profiler::basic {

using RegFile = std::map<RegIndex, int64_t>;

class RegFileParser {
public:
  static RegFile parse(std::ifstream &ifs);
};

} // namespace upmem_profiler::basic

#endif
