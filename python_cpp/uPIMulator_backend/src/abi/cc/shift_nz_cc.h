#ifndef UPMEM_SIM_ABI_ISA_CC_SHIFT_NZ_CC_H_
#define UPMEM_SIM_ABI_ISA_CC_SHIFT_NZ_CC_H_

#include "_base_cc.h"
#include "abi/isa/condition.h"

namespace upmem_sim::abi::cc {

class ShiftNZCC : public _BaseCC {
 public:
  explicit ShiftNZCC(isa::Condition condition)
      : _BaseCC({isa::Z, isa::NZ, isa::XZ, isa::XNZ, isa::E, isa::O, isa::PL,
                 isa::MI, isa::SZ, isa::SNZ, isa::SE, isa::SO, isa::SPL,
                 isa::SMI, isa::SH32, isa::NSH32, isa::TRUE},
                condition) {}
  ~ShiftNZCC() = default;
};

}  // namespace upmem_sim::abi::cc

#endif
