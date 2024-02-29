#ifndef UPMEM_PROFILER_UTIL_CONFIG_LOADER_H_
#define UPMEM_PROFILER_UTIL_CONFIG_LOADER_H_

namespace upmem_profiler::util {

class ConfigLoader {
public:
  static int num_gp_registers() { return 24; }
};

} // namespace upmem_profiler::util

#endif
