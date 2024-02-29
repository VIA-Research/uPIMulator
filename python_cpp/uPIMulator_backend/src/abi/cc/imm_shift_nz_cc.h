#ifndef UPMEM_SIM_ABI_ISA_CC_IMM_SHIFT_NZ_CC_H_
#define UPMEM_SIM_ABI_ISA_CC_IMM_SHIFT_NZ_CC_H_

#include "_base_cc.h"
#include "abi/isa/condition.h"

namespace upmem_sim::abi::cc {

class ImmShiftNZCC : public _BaseCC {
 public:
  explicit ImmShiftNZCC(isa::Condition condition)
      : _BaseCC({isa::Z, isa::NZ, isa::XZ, isa::XNZ, isa::E, isa::O, isa::PL,
                 isa::MI, isa::SZ, isa::SNZ, isa::SPL, isa::SMI, isa::SE,
                 isa::SO, isa::TRUE},
                condition) {}
  ~ImmShiftNZCC() = default;
};

}  // namespace upmem_sim::abi::cc

#endif
