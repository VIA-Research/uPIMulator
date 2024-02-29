#ifndef UPMEM_SIM_ABI_ISA_CC_DIV_NZ_CC_H_
#define UPMEM_SIM_ABI_ISA_CC_DIV_NZ_CC_H_

#include "_base_cc.h"
#include "abi/isa/condition.h"

namespace upmem_sim::abi::cc {

class DivNZCC : public _BaseCC {
 public:
  explicit DivNZCC(isa::Condition condition)
      : _BaseCC({isa::SZ, isa::SNZ, isa::SPL, isa::SMI, isa::TRUE}, condition) {
  }
  ~DivNZCC() = default;
};

}  // namespace upmem_sim::abi::cc

#endif
