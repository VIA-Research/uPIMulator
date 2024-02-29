#ifndef UPMEM_SIM_ABI_ISA_CC_NO_CC_H_
#define UPMEM_SIM_ABI_ISA_CC_NO_CC_H_

#include "_base_cc.h"
#include "abi/isa/condition.h"

namespace upmem_sim::abi::cc {

class NoCC : public _BaseCC {
 public:
  explicit NoCC(isa::Condition condition) : _BaseCC({}, condition) {}
  ~NoCC() = default;
};

}  // namespace upmem_sim::abi::cc

#endif
