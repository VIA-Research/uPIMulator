#ifndef UPMEM_SIM_SIMULATOR_REG_SP_REG_H_
#define UPMEM_SIM_SIMULATOR_REG_SP_REG_H_

#include "abi/reg/sp_reg.h"
#include "abi/word/data_word.h"

namespace upmem_sim::simulator::reg {

class SPReg {
 public:
  explicit SPReg(ThreadID id);
  ~SPReg();

  int64_t read(abi::reg::SPReg sp_reg,
               abi::word::Representation representation);
  void cycle() = delete;

 private:
  abi::word::DataWord *zero_;
  abi::word::DataWord *one_;
  abi::word::DataWord *lneg_;
  abi::word::DataWord *mneg_;
  abi::word::DataWord *id_;
  abi::word::DataWord *id2_;
  abi::word::DataWord *id4_;
  abi::word::DataWord *id8_;
};

}  // namespace upmem_sim::simulator::reg

#endif
