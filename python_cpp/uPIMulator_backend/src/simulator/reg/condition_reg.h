#ifndef UPMEM_SIM_SIMULATOR_REG_CONDITION_REG_H_
#define UPMEM_SIM_SIMULATOR_REG_CONDITION_REG_H_

#include <vector>

#include "abi/isa/condition.h"

namespace upmem_sim::simulator::reg {

class ConditionReg {
 public:
  explicit ConditionReg() { bits_.resize(abi::isa::LARGE + 1); }
  ~ConditionReg() = default;

  bool condition(abi::isa::Condition condition);
  void set_condition(abi::isa::Condition condition);
  void clear_condition(abi::isa::Condition condition);
  void clear_conditions();
  void cycle() = delete;

 private:
  std::vector<bool> bits_;
};

}  // namespace upmem_sim::simulator::reg

#endif
