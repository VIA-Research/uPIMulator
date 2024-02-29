#include "simulator/reg/sp_reg.h"

#include <stdexcept>

namespace upmem_sim::simulator::reg {

SPReg::SPReg(ThreadID id)
    : zero_(new abi::word::DataWord()),
      one_(new abi::word::DataWord()),
      lneg_(new abi::word::DataWord()),
      mneg_(new abi::word::DataWord()),
      id_(new abi::word::DataWord()),
      id2_(new abi::word::DataWord()),
      id4_(new abi::word::DataWord()),
      id8_(new abi::word::DataWord()) {
  zero_->set_value(0);
  one_->set_value(1);
  lneg_->set_value(-1);
  mneg_->set_bit(mneg_->width() - 1);
  id_->set_value(id);
  id2_->set_value(2 * id);
  id4_->set_value(4 * id);
  id8_->set_value(8 * id);
}

SPReg::~SPReg() {
  delete zero_;
  delete one_;
  delete lneg_;
  delete mneg_;
  delete id_;
  delete id2_;
  delete id4_;
  delete id8_;
}

int64_t SPReg::read(abi::reg::SPReg sp_reg,
                    abi::word::Representation representation) {
  if (sp_reg == abi::reg::ZERO) {
    return zero_->value(representation);
  } else if (sp_reg == abi::reg::ONE) {
    return one_->value(representation);
  } else if (sp_reg == abi::reg::LNEG) {
    return lneg_->value(representation);
  } else if (sp_reg == abi::reg::MNEG) {
    return mneg_->value(representation);
  } else if (sp_reg == abi::reg::ID) {
    return id_->value(representation);
  } else if (sp_reg == abi::reg::ID2) {
    return id2_->value(representation);
  } else if (sp_reg == abi::reg::ID4) {
    return id4_->value(representation);
  } else if (sp_reg == abi::reg::ID8) {
    return id8_->value(representation);
  } else {
    throw std::invalid_argument("");
  }
}

}  // namespace upmem_sim::simulator::reg
