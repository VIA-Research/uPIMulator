#include "basic/stats_parser.h"

#include <fstream>
#include <vector>

namespace upmem_profiler::basic {

Stats StatsParser::parse(std::ifstream &ifs) {
  Stats stats;

  while (true) {
    std::string line;
    std::getline(ifs, line);

    if (line == "") {
      break;
    }

    std::vector<std::string> tokens = split_by_colon(line);

    if (tokens.size() == 2) {
      std::string stat = tokens[0];
      int64_t value = std::stoi(tokens[1]);

      stats[stat] = value;
    }
  }

  return std::move(stats);
}

std::vector<std::string> StatsParser::split_by_colon(std::string line) {
  std::vector<std::string> tokens;
  int pos;
  while ((pos = line.find(":")) != std::string::npos) {
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

} // namespace upmem_profiler::basic
