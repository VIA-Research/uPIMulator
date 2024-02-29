#ifndef UPMEM_SIM_ABI_ISA_CC_COUNT_NZ_CC_H_
#define UPMEM_SIM_ABI_ISA_CC_COUNT_NZ_CC_H_

#include "_base_cc.h"
#include "abi/isa/condition.h"

namespace upmem_sim::abi::cc {

class CountNZCC : public _BaseCC {
 public:
  explicit CountNZCC(isa::Condition condition)
      : _BaseCC({isa::Z, isa::NZ, isa::XZ, isa::XNZ, isa::SZ, isa::SNZ,
                 isa::SPL, isa::SMI, isa::MAX, isa::NMAX, isa::TRUE},
                condition) {}
  ~CountNZCC() = default;
};

}  // namespace upmem_sim::abi::cc

#endif
