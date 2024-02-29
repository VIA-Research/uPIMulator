#ifndef UPMEM_SIM_CONVERTER_REG_CONVERTER_H_
#define UPMEM_SIM_CONVERTER_REG_CONVERTER_H_

#include <string>

#include "abi/reg/gp_reg.h"
#include "abi/reg/pair_reg.h"
#include "abi/reg/sp_reg.h"
#include "abi/reg/src_reg.h"

namespace upmem_sim::converter {

class RegConverter {
 public:
  static std::string to_string(abi::reg::GPReg *gp_reg) {
    return "r" + std::to_string(gp_reg->index());
  }
  static std::string to_string(abi::reg::PairReg *pair_reg) {
    return "d" + std::to_string(pair_reg->even_reg()->index());
  }
  static std::string to_string(abi::reg::SPReg sp_reg);
  static std::string to_string(abi::reg::SrcReg *src_reg);
};

}  // namespace upmem_sim::converter

#endif
