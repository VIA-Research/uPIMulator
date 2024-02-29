#ifndef UPMEM_SIM_ABI_ISA_CC_ACQUIRE_CC_H_
#define UPMEM_SIM_ABI_ISA_CC_ACQUIRE_CC_H_

#include "_base_cc.h"
#include "abi/isa/condition.h"

namespace upmem_sim::abi::cc {

class AcquireCC : public _BaseCC {
 public:
  explicit AcquireCC(isa::Condition condition)
      : _BaseCC({isa::Z, isa::NZ, isa::TRUE}, condition) {}
  ~AcquireCC() = default;
};

}  // namespace upmem_sim::abi::cc

#endif
