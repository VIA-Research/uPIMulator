#ifndef UPMEM_SIM_ABI_ISA_REG_PAIR_REG_H_
#define UPMEM_SIM_ABI_ISA_REG_PAIR_REG_H_

#include <cassert>

#include "gp_reg.h"

namespace upmem_sim::abi::reg {

class PairReg {
 public:
  explicit PairReg(RegIndex index);
  ~PairReg();

  GPReg *even_reg() { return even_reg_; }
  GPReg *odd_reg() { return odd_reg_; }

 private:
  GPReg *even_reg_;
  GPReg *odd_reg_;
};

}  // namespace upmem_sim::abi::reg

#endif
