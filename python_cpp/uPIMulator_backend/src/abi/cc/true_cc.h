#ifndef UPMEM_SIM_ABI_ISA_CC_TRUE_CC_H_
#define UPMEM_SIM_ABI_ISA_CC_TRUE_CC_H_

#include "_base_cc.h"
#include "abi/isa/condition.h"

namespace upmem_sim::abi::cc {

class TrueCC : public _BaseCC {
 public:
  explicit TrueCC(isa::Condition condition) : _BaseCC({isa::TRUE}, condition) {}
  ~TrueCC() = default;
};

}  // namespace upmem_sim::abi::cc

#endif
