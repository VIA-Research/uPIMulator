#ifndef UPMEM_PROFILER_BASIC_INTERVAL_H_
#define UPMEM_PROFILER_BASIC_INTERVAL_H_

#include <cstdint>

#include "basic/stats_parser.h"

namespace upmem_profiler::basic {

class Interval {
public:
  explicit Interval() = default;
  ~Interval() = default;

  void set_begin(Stats begin) { begin_ = begin; }
  void set_end(Stats end) { end_ = end; }

  int64_t begin_value(std::string stat) { return begin_[stat]; }
  int64_t end_value(std::string stat) { return end_[stat]; }

  Stats stats();

private:
  Stats begin_;
  Stats end_;
};

} // namespace upmem_profiler::basic

#endif
