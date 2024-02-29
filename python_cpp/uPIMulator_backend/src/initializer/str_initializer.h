#ifndef UPMEM_SIM_INITIALIZER_STR_INITIALIZER_H_
#define UPMEM_SIM_INITIALIZER_STR_INITIALIZER_H_

#include <string>

namespace upmem_sim::initializer {

class StrInitializer {
 public:
  static std::string identifier(int width);
};

}  // namespace upmem_sim::initializer

#endif
