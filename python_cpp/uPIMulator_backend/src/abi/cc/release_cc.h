#ifndef UPMEM_SIM_ABI_ISA_CC_RELEASE_CC_H_
#define UPMEM_SIM_ABI_ISA_CC_RELEASE_CC_H_

#include "_base_cc.h"
#include "abi/isa/condition.h"

namespace upmem_sim::abi::cc {

class ReleaseCC : public _BaseCC {
 public:
  explicit ReleaseCC(isa::Condition condition)
      : _BaseCC({isa::NZ}, condition) {}
  ~ReleaseCC() = default;
};

}  // namespace upmem_sim::abi::cc

#endif
