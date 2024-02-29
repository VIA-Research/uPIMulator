#ifndef UPMEM_SIM_ABI_ISA_CC_CONST_CC_GEU_H_
#define UPMEM_SIM_ABI_ISA_CC_CONST_CC_GEU_H_

#include "_base_cc.h"
#include "abi/isa/condition.h"

namespace upmem_sim::abi::cc {

class ConstCCGEU : public _BaseCC {
 public:
  explicit ConstCCGEU(isa::Condition condition)
      : _BaseCC({isa::GEU}, condition) {}
  ~ConstCCGEU() = default;
};

}  // namespace upmem_sim::abi::cc

#endif
