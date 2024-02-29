#include "simulator/reg/condition_reg.h"

#include <cassert>

namespace upmem_sim::simulator::reg {

bool ConditionReg::condition(abi::isa::Condition condition) {
  if (condition == abi::isa::TRUE) {
    return true;
  } else if (condition == abi::isa::FALSE) {
    return false;
  } else {
    return bits_[condition];
  }
}

void ConditionReg::set_condition(abi::isa::Condition condition) {
  assert(condition != abi::isa::TRUE and condition != abi::isa::FALSE);
  bits_[condition] = true;
}

void ConditionReg::clear_condition(abi::isa::Condition condition) {
  assert(condition != abi::isa::TRUE and condition != abi::isa::FALSE);
  bits_[condition] = false;
}

void ConditionReg::clear_conditions() {
  for (abi::isa::Condition condition = abi::isa::TRUE;
       condition <= abi::isa::LARGE;
       condition = static_cast<abi::isa::Condition>(condition + 1)) {
    if (condition == abi::isa::TRUE or condition == abi::isa::FALSE) {
      continue;
    } else {
      clear_condition(condition);
    }
  }
}

}  // namespace upmem_sim::simulator::reg
