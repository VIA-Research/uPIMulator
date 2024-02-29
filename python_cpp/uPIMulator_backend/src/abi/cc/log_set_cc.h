#ifndef UPMEM_SIM_ABI_ISA_CC_LOG_SET_CC_H_
#define UPMEM_SIM_ABI_ISA_CC_LOG_SET_CC_H_

#include "_base_cc.h"
#include "abi/isa/condition.h"

namespace upmem_sim::abi::cc {

class LogSetCC : public _BaseCC {
 public:
  explicit LogSetCC(isa::Condition condition)
      : _BaseCC({isa::Z, isa::NZ, isa::XZ, isa::XNZ}, condition) {}
  ~LogSetCC() = default;
};

}  // namespace upmem_sim::abi::cc

#endif
