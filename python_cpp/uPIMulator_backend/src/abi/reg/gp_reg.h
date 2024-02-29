#ifndef UPMEM_SIM_ABI_ISA_REG_GP_REG_H_
#define UPMEM_SIM_ABI_ISA_REG_GP_REG_H_

#include <cassert>
#include <string>

#include "util/config_loader.h"

namespace upmem_sim::abi::reg {

class GPReg {
 public:
  explicit GPReg(RegIndex index) : index_(index) {
    assert(0 <= index and index < util::ConfigLoader::num_gp_registers());
  }
  ~GPReg() = default;

  RegIndex index() { return index_; }

 private:
  RegIndex index_;
};

}  // namespace upmem_sim::abi::reg

#endif
