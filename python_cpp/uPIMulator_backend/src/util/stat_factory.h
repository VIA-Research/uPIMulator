#ifndef UPMEM_SIM_UTIL_STAT_FACTORY_H_
#define UPMEM_SIM_UTIL_STAT_FACTORY_H_

#include <map>
#include <set>
#include <string>

namespace upmem_sim::util {

class StatFactory {
 public:
  explicit StatFactory(std::string name) : name_(name) {}
  ~StatFactory() = default;

  std::string name() { return name_; }

  std::set<std::string> stats();
  int64_t value(std::string stat) { return stats_[stat]; }

  void increment(std::string stat) { increment(stat, 1); }
  void increment(std::string stat, int64_t value);
  void overwrite(std::string stat, int64_t value);

  void merge(StatFactory *stat_factory);

 private:
  std::string name_;
  std::map<std::string, int64_t> stats_;
};

}  // namespace upmem_sim::util

#endif
