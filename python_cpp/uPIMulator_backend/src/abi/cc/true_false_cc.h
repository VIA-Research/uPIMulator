#ifndef UPMEM_SIM_ABI_ISA_CC_TRUE_FALSE_CC_H_
#define UPMEM_SIM_ABI_ISA_CC_TRUE_FALSE_CC_H_

#include "_base_cc.h"
#include "abi/isa/condition.h"

namespace upmem_sim::abi::cc {

class TrueFalseCC : public _BaseCC {
 public:
  explicit TrueFalseCC(isa::Condition condition)
      : _BaseCC({isa::TRUE, isa::FALSE}, condition) {}
  ~TrueFalseCC() = default;
};

}  // namespace upmem_sim::abi::cc

#endif
