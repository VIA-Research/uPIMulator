#ifndef UPMEM_SIM_SIMULATOR_REG_REG_FILE_H_
#define UPMEM_SIM_SIMULATOR_REG_REG_FILE_H_

#include "abi/reg/pair_reg.h"
#include "abi/reg/src_reg.h"
#include "simulator/reg/condition_reg.h"
#include "simulator/reg/exception_reg.h"
#include "simulator/reg/flag_reg.h"
#include "simulator/reg/gp_reg.h"
#include "simulator/reg/pc_reg.h"
#include "simulator/reg/sp_reg.h"

namespace upmem_sim::simulator::reg {

class RegFile {
 public:
  explicit RegFile(ThreadID id);
  ~RegFile();

  int64_t read_gp_reg(abi::reg::GPReg *gp_reg,
                      abi::word::Representation representation) {
    return gp_regs_[gp_reg->index()]->read(representation);
  }
  int64_t read_sp_reg(abi::reg::SPReg sp_reg,
                      abi::word::Representation representation) {
    return sp_reg_->read(sp_reg, representation);
  }
  std::tuple<int64_t, int64_t> read_pair_reg(
      abi::reg::PairReg *pair_reg, abi::word::Representation representation);
  int64_t read_src_reg(abi::reg::SrcReg *src_reg,
                       abi::word::Representation representation);

  int64_t read_pc_reg() { return pc_reg_->read(); }
  bool condition(abi::isa::Condition condition) {
    return condition_reg_->condition(condition);
  }
  bool flag(abi::isa::Flag flag) { return flag_reg_->flag(flag); }
  bool exception(abi::isa::Exception exception) {
    return exception_reg_->exception(exception);
  }

  void write_gp_reg(abi::reg::GPReg *gp_reg, int64_t value) {
    gp_regs_[gp_reg->index()]->write(value);
  }
  void write_pair_reg(abi::reg::PairReg *pair_reg, int64_t even, int64_t odd);
  void write_pc_reg(int64_t value) { pc_reg_->write(value); }
  void increment_pc_reg() { pc_reg_->increment(); }

  void set_condition(abi::isa::Condition condition) {
    condition_reg_->set_condition(condition);
  }
  void clear_condition(abi::isa::Condition condition) {
    condition_reg_->clear_condition(condition);
  }
  void clear_conditions() { condition_reg_->clear_conditions(); }

  void set_flag(abi::isa::Flag flag) { flag_reg_->set_flag(flag); }
  void clear_flag(abi::isa::Flag flag) { flag_reg_->clear_flag(flag); }

  void set_exception(abi::isa::Exception exception) {
    exception_reg_->set_exception(exception);
  }
  void clear_exception(abi::isa::Exception exception) {
    exception_reg_->clear_exception(exception);
  }

  void cycle() = delete;

 private:
  std::vector<GPReg *> gp_regs_;
  SPReg *sp_reg_;
  PCReg *pc_reg_;
  ConditionReg *condition_reg_;
  FlagReg *flag_reg_;
  ExceptionReg *exception_reg_;
};

}  // namespace upmem_sim::simulator::reg

#endif
