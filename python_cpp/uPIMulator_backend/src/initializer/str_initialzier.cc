#include <cassert>
#include <cmath>
#include <random>

#include "initializer/str_initializer.h"

namespace upmem_sim::initializer {

std::string StrInitializer::identifier(int width) {
  assert(width > 0);

  std::string characters =
      "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789._";

  std::random_device random_device;
  std::mt19937 generator(random_device());
  std::uniform_int_distribution<int64_t> distribution;

  std::string identifier = "";
  for (int i = 0; i < width; i++) {
    int index = static_cast<int>(distribution(generator) % characters.length());
    identifier += characters.substr(index, 1);
  }

  return identifier;
}

}  // namespace upmem_sim::initializer
