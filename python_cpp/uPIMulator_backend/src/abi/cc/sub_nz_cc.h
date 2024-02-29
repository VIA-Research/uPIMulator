#ifndef UPMEM_SIM_ABI_ISA_CC_SUB_NZ_CC_H_
#define UPMEM_SIM_ABI_ISA_CC_SUB_NZ_CC_H_

#include "_base_cc.h"
#include "abi/isa/condition.h"

namespace upmem_sim::abi::cc {

class SubNZCC : public _BaseCC {
 public:
  explicit SubNZCC(isa::Condition condition)
      : _BaseCC({isa::C,    isa::NC,   isa::Z,   isa::NZ,  isa::XZ,   isa::XNZ,
                 isa::OV,   isa::NOV,  isa::MI,  isa::PL,  isa::EQ,   isa::NEQ,
                 isa::SPL,  isa::SMI,  isa::GES, isa::GEU, isa::GTS,  isa::GTU,
                 isa::LES,  isa::LEU,  isa::LTS, isa::LTU, isa::XGTS, isa::XGTU,
                 isa::XLES, isa::XLEU, isa::TRUE},
                condition) {}
  ~SubNZCC() = default;
};

}  // namespace upmem_sim::abi::cc

#endif
