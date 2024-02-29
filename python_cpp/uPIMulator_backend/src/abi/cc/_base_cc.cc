#include "_base_cc.h"

#include <cassert>

namespace upmem_sim::abi::cc {

_BaseCC::_BaseCC(std::set<isa::Condition> conditions, isa::Condition condition)
    : condition_(condition) {
  assert(conditions.count(condition) != 0);
}

}  // namespace upmem_sim::abi::cc
