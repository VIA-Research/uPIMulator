#ifndef UPMEM_SIM_ABI_ISA_CC_CONST_CC_ZERO_H_
#define UPMEM_SIM_ABI_ISA_CC_CONST_CC_ZERO_H_

#include "_base_cc.h"
#include "abi/isa/condition.h"

namespace upmem_sim::abi::cc {

class ConstCCZero : public _BaseCC {
 public:
  explicit ConstCCZero(isa::Condition condition)
      : _BaseCC({isa::Z}, condition) {}
  ~ConstCCZero() = default;
};

}  // namespace upmem_sim::abi::cc

#endif
