#ifndef UPMEM_SIM_SIMULATOR_REG_EXCEPTION_REG_H_
#define UPMEM_SIM_SIMULATOR_REG_EXCEPTION_REG_H_

#include <vector>

#include "abi/isa/exception.h"

namespace upmem_sim::simulator::reg {

class ExceptionReg {
 public:
  explicit ExceptionReg() { bits_.resize(abi::isa::NOT_PROFILING + 1); }
  ~ExceptionReg() = default;

  bool exception(abi::isa::Exception exception) { return bits_[exception]; }
  void set_exception(abi::isa::Exception exception) { bits_[exception] = true; }
  void clear_exception(abi::isa::Exception exception) {
    bits_[exception] = false;
  }
  void cycle() = delete;

 private:
  std::vector<bool> bits_;
};

}  // namespace upmem_sim::simulator::reg

#endif
