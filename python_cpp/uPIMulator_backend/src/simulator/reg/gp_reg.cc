#include "simulator/reg/gp_reg.h"

namespace upmem_sim::simulator::reg {

GPReg::~GPReg() {
  delete gp_reg_;
  delete word_;
}

}  // namespace upmem_sim::simulator::reg
