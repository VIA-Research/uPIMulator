#include "simulator/reg/reg_file.h"

#include <stdexcept>
#include <tuple>

namespace upmem_sim::simulator::reg {

RegFile::RegFile(ThreadID id)
    : sp_reg_(new SPReg(id)),
      pc_reg_(new PCReg()),
      condition_reg_(new ConditionReg()),
      flag_reg_(new FlagReg()),
      exception_reg_(new ExceptionReg()) {
  for (RegIndex index = 0; index < util::ConfigLoader::num_gp_registers();
       index++) {
    gp_regs_.push_back(new GPReg(index));
  }
}

RegFile::~RegFile() {
  for (RegIndex index = 0; index < util::ConfigLoader::num_gp_registers();
       index++) {
    delete gp_regs_[index];
  }
  delete sp_reg_;
  delete pc_reg_;
  delete condition_reg_;
  delete flag_reg_;
  delete exception_reg_;
}

std::tuple<int64_t, int64_t> RegFile::read_pair_reg(
    abi::reg::PairReg *pair_reg, abi::word::Representation representation) {
  int64_t even = read_gp_reg(pair_reg->even_reg(), representation);
  int64_t odd = read_gp_reg(pair_reg->odd_reg(), abi::word::UNSIGNED);
  return {even, odd};
}

int64_t RegFile::read_src_reg(abi::reg::SrcReg *src_reg,
                              abi::word::Representation representation) {
  if (src_reg->is_gp_reg()) {
    return read_gp_reg(src_reg->gp_reg(), representation);
  } else {
    return read_sp_reg(*src_reg->sp_reg(), representation);
  }
}

void RegFile::write_pair_reg(abi::reg::PairReg *pair_reg, int64_t even,
                             int64_t odd) {
  write_gp_reg(pair_reg->even_reg(), even);
  write_gp_reg(pair_reg->odd_reg(), odd);
}

}  // namespace upmem_sim::simulator::reg
