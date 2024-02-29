#ifndef UPMEM_SIM_ABI_ISA_CC_SUB_SET_CC_H_
#define UPMEM_SIM_ABI_ISA_CC_SUB_SET_CC_H_

#include "_base_cc.h"
#include "abi/isa/condition.h"

namespace upmem_sim::abi::cc {

class SubSetCC : public _BaseCC {
 public:
  explicit SubSetCC(isa::Condition condition)
      : _BaseCC({isa::Z, isa::NZ, isa::XZ, isa::XNZ, isa::EQ, isa::NEQ},
                condition) {}
  ~SubSetCC() = default;
};

}  // namespace upmem_sim::abi::cc

#endif
