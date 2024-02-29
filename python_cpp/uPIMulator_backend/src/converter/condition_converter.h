#ifndef UPMEM_SIM_CONVERTER_CONDITION_CONVERTER_H_
#define UPMEM_SIM_CONVERTER_CONDITION_CONVERTER_H_

#include <string>

#include "abi/isa/condition.h"

namespace upmem_sim::converter {

class ConditionConverter {
 public:
  static std::string to_string(abi::isa::Condition condition);
};

}  // namespace upmem_sim::converter

#endif
