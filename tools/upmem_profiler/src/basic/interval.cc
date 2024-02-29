#include "basic/interval.h"

namespace upmem_profiler::basic {

Stats Interval::stats() {
  Stats stats;
  for (auto &[end_stat, end_value] : end_) {
    if (begin_.count(end_stat)) {
      stats[end_stat] = end_value - begin_[end_stat];
    } else {
      stats[end_stat] = end_value;
    }
  }
  return std::move(stats);
}

} // namespace upmem_profiler::basic
