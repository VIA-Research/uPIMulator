#include "pair_reg.h"

#include <cassert>

namespace upmem_sim::abi::reg {

PairReg::PairReg(RegIndex index) {
  assert(index % 2 == 0);

  even_reg_ = new GPReg(index);
  odd_reg_ = new GPReg(index + 1);
}

PairReg::~PairReg() {
  delete even_reg_;
  delete odd_reg_;
}

}  // namespace upmem_sim::abi::reg
