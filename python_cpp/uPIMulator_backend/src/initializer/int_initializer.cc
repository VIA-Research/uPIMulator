#include "initializer/int_initializer.h"

#include <cassert>
#include <cmath>
#include <random>

namespace upmem_sim::initializer {

int64_t IntInitializer::value_by_range(int64_t min_value, int64_t max_value) {
  assert(min_value < max_value);

  std::random_device random_device;
  std::mt19937 generator(random_device());
  std::uniform_int_distribution<int64_t> distribution;

  int64_t range = max_value - min_value;
  return distribution(generator) % range + min_value;
}

int64_t IntInitializer::value_by_width(abi::word::Representation representation,
                                       int width) {
  if (representation == abi::word::UNSIGNED) {
    return value_by_range(0, static_cast<int64_t>(pow(2, width)));
  } else {
    return value_by_range(-static_cast<int64_t>(pow(2, width - 1)),
                          static_cast<int64_t>(pow(2, width - 1)));
  }
}

}  // namespace upmem_sim::initializer
