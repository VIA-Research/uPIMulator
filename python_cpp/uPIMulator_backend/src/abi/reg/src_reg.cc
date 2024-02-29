#include "src_reg.h"

#include <cassert>

namespace upmem_sim::abi::reg {

SrcReg::~SrcReg() {
  delete gp_reg_;
  delete sp_reg_;
}

GPReg *SrcReg::gp_reg() {
  assert(is_gp_reg());
  return gp_reg_;
}

SPReg *SrcReg::sp_reg() {
  assert(is_sp_reg());
  return sp_reg_;
}

}  // namespace upmem_sim::abi::reg
