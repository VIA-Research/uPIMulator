#ifndef UPMEM_SIM_ABI_ISA_CC__BASE_CC_H_
#define UPMEM_SIM_ABI_ISA_CC__BASE_CC_H_

#include <set>

#include "abi/isa/condition.h"

namespace upmem_sim::abi::cc {

class _BaseCC {
 public:
  explicit _BaseCC(std::set<isa::Condition> conditions,
                   isa::Condition condition);
  ~_BaseCC() = default;

  isa::Condition condition() { return condition_; }

 private:
  isa::Condition condition_;
};

}  // namespace upmem_sim::abi::cc

#endif
