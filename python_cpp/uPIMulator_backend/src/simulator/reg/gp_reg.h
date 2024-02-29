#ifndef UPMEM_SIM_SIMULATOR_REG_GP_REG_H_
#define UPMEM_SIM_SIMULATOR_REG_GP_REG_H_

#include "abi/reg/gp_reg.h"
#include "abi/word/data_word.h"

namespace upmem_sim::simulator::reg {

class GPReg {
 public:
  explicit GPReg(RegIndex index)
      : gp_reg_(new abi::reg::GPReg(index)), word_(new abi::word::DataWord()) {}
  ~GPReg();

  RegIndex index() { return gp_reg_->index(); }
  int64_t read(abi::word::Representation representation) {
    return word_->value(representation);
  }
  void write(int64_t value) { word_->set_value(value); }
  void cycle() = delete;

 private:
  abi::reg::GPReg *gp_reg_;
  abi::word::DataWord *word_;
};

}  // namespace upmem_sim::simulator::reg

#endif
