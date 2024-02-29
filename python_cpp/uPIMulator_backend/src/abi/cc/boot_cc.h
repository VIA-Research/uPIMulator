#ifndef UPMEM_SIM_ABI_ISA_CC_BOOT_CC_H_
#define UPMEM_SIM_ABI_ISA_CC_BOOT_CC_H_

#include "_base_cc.h"
#include "abi/isa/condition.h"

namespace upmem_sim::abi::cc {

class BootCC : public _BaseCC {
 public:
  explicit BootCC(isa::Condition condition)
      : _BaseCC({isa::Z, isa::NZ, isa::XZ, isa::XNZ, isa::SZ, isa::SNZ,
                 isa::SPL, isa::SMI, isa::TRUE, isa::FALSE},
                condition) {}
  ~BootCC() = default;
};

}  // namespace upmem_sim::abi::cc

#endif
