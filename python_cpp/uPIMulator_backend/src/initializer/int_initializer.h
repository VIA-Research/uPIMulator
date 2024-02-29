#ifndef UPMEM_SIM_INITIALIZER_INT_INITIALIZER_H_
#define UPMEM_SIM_INITIALIZER_INT_INITIALIZER_H_

#include <cstdint>

#include "abi/word/representation.h"

namespace upmem_sim::initializer {

class IntInitializer {
 public:
  static int64_t value_by_range(int64_t min_value, int64_t max_value);
  static int64_t value_by_width(abi::word::Representation representation,
                                int width);
};

}  // namespace upmem_sim::initializer

#endif
