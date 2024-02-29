#ifndef UPMEM_SIM_ABI_ISA_REG_SRC_REG_H_
#define UPMEM_SIM_ABI_ISA_REG_SRC_REG_H_

#include "gp_reg.h"
#include "sp_reg.h"

namespace upmem_sim::abi::reg {

class SrcReg {
 public:
  explicit SrcReg(GPReg *reg)
      : gp_reg_(new GPReg(reg->index())), sp_reg_(nullptr) {}
  explicit SrcReg(SPReg *reg) : gp_reg_(nullptr), sp_reg_(new SPReg(*reg)) {}
  ~SrcReg();

  bool is_gp_reg() { return gp_reg_ != nullptr; }
  bool is_sp_reg() { return sp_reg_ != nullptr; }

  GPReg *gp_reg();
  SPReg *sp_reg();

 private:
  GPReg *gp_reg_;
  SPReg *sp_reg_;
};

}  // namespace upmem_sim::abi::reg

#endif
