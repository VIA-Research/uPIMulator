#include "converter/reg_converter.h"

#include <stdexcept>

namespace upmem_sim::converter {

std::string RegConverter::to_string(abi::reg::SPReg sp_reg) {
  if (sp_reg == abi::reg::ZERO) {
    return "zero";
  } else if (sp_reg == abi::reg::ONE) {
    return "one";
  } else if (sp_reg == abi::reg::LNEG) {
    return "lneg";
  } else if (sp_reg == abi::reg::MNEG) {
    return "mneg";
  } else if (sp_reg == abi::reg::ID) {
    return "id";
  } else if (sp_reg == abi::reg::ID2) {
    return "id2";
  } else if (sp_reg == abi::reg::ID4) {
    return "id4";
  } else if (sp_reg == abi::reg::ID8) {
    return "id8";
  } else {
    throw std::invalid_argument("");
  }
}

std::string RegConverter::to_string(abi::reg::SrcReg *src_reg) {
  if (src_reg->is_gp_reg()) {
    return to_string(src_reg->gp_reg());
  } else if (src_reg->is_sp_reg()) {
    return to_string(*src_reg->sp_reg());
  } else {
    throw std::invalid_argument("");
  }
}

}  // namespace upmem_sim::converter
