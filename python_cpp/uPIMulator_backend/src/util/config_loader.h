#ifndef UPMEM_SIM_UTIL_CONFIG_LOADER_H_
#define UPMEM_SIM_UTIL_CONFIG_LOADER_H_

#include "main.h"

namespace upmem_sim::util {

class ConfigLoader {
 public:
  static int atomic_address_width() { return 32; }
  static int atomic_data_width() { return 32; }
  static Address atomic_offset() { return 0; }
  static Address atomic_size() { return 256; }

  static int iram_address_width() { return 32; }
  static int iram_data_width() { return 96; }
  static Address iram_offset() { return 384 * 1024; }
  static Address iram_size() { return 48 * 1024; }

  static int wram_address_width() { return 32; }
  static int wram_data_width() { return 32; }
  static Address wram_offset() { return 512; }
  static Address wram_size() { return 128 * 1024; }

  static Address stack_size() { return 2 * 1024; }
  static Address heap_size() { return 4 * 1024; }

  static int mram_address_width() { return 32; }
  static int mram_data_width() { return 32; }
  static Address mram_offset() { return 512 * 1024; }
  static Address mram_size() { return 64 * 1024 * 1024; }

  static int num_gp_registers() { return 24; }
  static int max_num_tasklets() { return 24; }
  static int min_access_granularity() { return 8; }
};

}  // namespace upmem_sim::util

#endif
