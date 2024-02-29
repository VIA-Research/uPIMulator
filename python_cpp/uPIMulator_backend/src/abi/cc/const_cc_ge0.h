#ifndef UPMEM_SIM_ABI_ISA_CC_CONST_CC_GE0_H_
#define UPMEM_SIM_ABI_ISA_CC_CONST_CC_GE0_H_

#include "_base_cc.h"
#include "abi/isa/condition.h"

namespace upmem_sim::abi::cc {

class ConstCCGE0 : public _BaseCC {
 public:
  explicit ConstCCGE0(isa::Condition condition)
      : _BaseCC({isa::PL}, condition) {}
  ~ConstCCGE0() = default;
};

}  // namespace upmem_sim::abi::cc

#endif
