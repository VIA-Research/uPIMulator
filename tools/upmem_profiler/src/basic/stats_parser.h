#ifndef UPMEM_PROFILER_BASIC_STATS_PARSER_H_
#define UPMEM_PROFILER_BASIC_STATS_PARSER_H_

#include <map>
#include <set>
#include <string>
#include <vector>

namespace upmem_profiler::basic {

using Stats = std::map<std::string, int64_t>;

class StatsParser {
public:
  static Stats parse(std::ifstream &ifs);

protected:
  static std::vector<std::string> split_by_colon(std::string line);

private:
  std::set<std::string> stats_;
};

} // namespace upmem_profiler::basic

#endif
