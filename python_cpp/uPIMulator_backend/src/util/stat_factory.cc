#include "util/stat_factory.h"

#include <cassert>

namespace upmem_sim::util {

std::set<std::string> StatFactory::stats() {
  std::set<std::string> stats;
  for (auto &[stat, _] : stats_) {
    stats.insert(stat);
  }
  return std::move(stats);
}

void StatFactory::increment(std::string stat, int64_t value) {
  if (stats_.count(stat)) {
    stats_[stat] += value;
  } else {
    stats_[stat] = value;
  }
}

void StatFactory::overwrite(std::string stat, int64_t value) {
  stats_[stat] = value;
}

void StatFactory::merge(StatFactory *stat_factory) {
  for (auto &[stat, value] : stat_factory->stats_) {
    if (stats_.count(stat_factory->name() + "/" + stat)) {
      stats_[stat_factory->name() + "/" + stat] += value;
    } else {
      stats_[stat_factory->name() + "/" + stat] = value;
    }
  }
}

}  // namespace upmem_sim::util
