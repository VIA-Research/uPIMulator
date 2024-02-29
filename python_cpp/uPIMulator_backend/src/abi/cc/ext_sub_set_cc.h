#ifndef UPMEM_SIM_ABI_ISA_CC_EXT_SUB_SET_CC_H_
#define UPMEM_SIM_ABI_ISA_CC_EXT_SUB_SET_CC_H_

#include "_base_cc.h"
#include "abi/isa/condition.h"

namespace upmem_sim::abi::cc {

class ExtSubSetCC : public _BaseCC {
 public:
  explicit ExtSubSetCC(isa::Condition condition)
      : _BaseCC({isa::C,    isa::NC,   isa::Z,    isa::NZ,   isa::XZ,  isa::XNZ,
                 isa::OV,   isa::NOV,  isa::EQ,   isa::NEQ,  isa::PL,  isa::MI,
                 isa::SZ,   isa::SNZ,  isa::SPL,  isa::SMI,  isa::GES, isa::GEU,
                 isa::GTS,  isa::GTU,  isa::LES,  isa::LEU,  isa::LTS, isa::LTU,
                 isa::XGTS, isa::XGTU, isa::XLES, isa::XLEU, isa::TRUE},
                condition) {}
  ~ExtSubSetCC() = default;
};

}  // namespace upmem_sim::abi::cc

#endif
