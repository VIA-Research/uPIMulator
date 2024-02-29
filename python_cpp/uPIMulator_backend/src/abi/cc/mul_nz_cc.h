#ifndef UPMEM_SIM_ABI_ISA_CC_MUL_NZ_CC_H_
#define UPMEM_SIM_ABI_ISA_CC_MUL_NZ_CC_H_

#include "_base_cc.h"
#include "abi/isa/condition.h"

namespace upmem_sim::abi::cc {

class MulNZCC : public _BaseCC {
 public:
  explicit MulNZCC(isa::Condition condition)
      : _BaseCC({isa::Z, isa::NZ, isa::XZ, isa::XNZ, isa::SZ, isa::SNZ,
                 isa::SPL, isa::SMI, isa::LARGE, isa::SMALL, isa::TRUE},
                condition) {}
  ~MulNZCC() = default;
};

}  // namespace upmem_sim::abi::cc

#endif
