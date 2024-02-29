#ifndef UPMEM_SIM_ABI_ISA_CC_ADD_NZ_CC_H_
#define UPMEM_SIM_ABI_ISA_CC_ADD_NZ_CC_H_

#include "_base_cc.h"
#include "abi/isa/condition.h"

namespace upmem_sim::abi::cc {

class AddNZCC : public _BaseCC {
 public:
  explicit AddNZCC(isa::Condition condition)
      : _BaseCC({isa::C,    isa::NC,   isa::Z,    isa::NZ,   isa::XZ,
                 isa::XNZ,  isa::OV,   isa::NOV,  isa::PL,   isa::MI,
                 isa::SZ,   isa::SNZ,  isa::SPL,  isa::SMI,  isa::NC5,
                 isa::NC6,  isa::NC7,  isa::NC8,  isa::NC9,  isa::NC10,
                 isa::NC11, isa::NC12, isa::NC13, isa::NC14, isa::TRUE},
                condition) {}
  ~AddNZCC() = default;
};

}  // namespace upmem_sim::abi::cc

#endif
