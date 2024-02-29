#ifndef UPMEM_SIM_ABI_ISA_CC_LOG_NZ_CC_H_
#define UPMEM_SIM_ABI_ISA_CC_LOG_NZ_CC_H_

#include "_base_cc.h"
#include "abi/isa/condition.h"

namespace upmem_sim::abi::cc {

class LogNZCC : public _BaseCC {
 public:
  explicit LogNZCC(isa::Condition condition)
      : _BaseCC({isa::Z, isa::NZ, isa::XZ, isa::XNZ, isa::PL, isa::MI, isa::SZ,
                 isa::SNZ, isa::SPL, isa::SMI, isa::TRUE},
                condition) {}
  ~LogNZCC() = default;
};

}  // namespace upmem_sim::abi::cc

#endif
