#ifndef UPMEM_SIM_ABI_ISA_CC_FALSE_CC_H_
#define UPMEM_SIM_ABI_ISA_CC_FALSE_CC_H_

#include "_base_cc.h"
#include "abi/isa/condition.h"

namespace upmem_sim::abi::cc {

class FalseCC : public _BaseCC {
 public:
  explicit FalseCC(isa::Condition condition)
      : _BaseCC({isa::FALSE}, condition) {}
  ~FalseCC() = default;
};

}  // namespace upmem_sim::abi::cc

#endif
