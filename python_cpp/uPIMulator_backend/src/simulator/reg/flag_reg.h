#ifndef UPMEM_SIM_SIMULATOR_REG_FLAG_REG_H_
#define UPMEM_SIM_SIMULATOR_REG_FLAG_REG_H_

#include <vector>

#include "abi/isa/flag.h"

namespace upmem_sim::simulator::reg {

class FlagReg {
 public:
  explicit FlagReg() { bits_.resize(abi::isa::CARRY + 1); }
  ~FlagReg() = default;

  bool flag(abi::isa::Flag flag) { return bits_[flag]; }
  void set_flag(abi::isa::Flag flag) { bits_[flag] = true; }
  void clear_flag(abi::isa::Flag flag) { bits_[flag] = false; }
  void cycle() = delete;

 private:
  std::vector<bool> bits_;
};

}  // namespace upmem_sim::simulator::reg

#endif
