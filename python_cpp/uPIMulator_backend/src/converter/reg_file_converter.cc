#include "converter/reg_file_converter.h"

#include <sstream>

#include "converter/condition_converter.h"
#include "converter/flag_converter.h"
#include "converter/reg_converter.h"

namespace upmem_sim::converter {

std::string RegFileConverter::to_string(simulator::reg::RegFile *reg_file) {
  std::stringstream ss;
  for (RegIndex index = 0; index < util::ConfigLoader::num_gp_registers();
       index++) {
    auto gp_reg = new abi::reg::GPReg(index);
    ss << RegConverter::to_string(gp_reg) << ": "
       << reg_file->read_gp_reg(gp_reg, abi::word::SIGNED) << std::endl;
    delete gp_reg;
  }

  /*
  ss << "pc: " << reg_file->read_pc_reg() << std::endl;

  for (abi::isa::Condition condition = abi::isa::TRUE; condition !=
  abi::isa::LARGE; condition = static_cast<abi::isa::Condition>(condition + 1))
  { ss << ConditionConverter::to_string(condition) << ": " <<
  reg_file->condition(condition) << std::endl;
  }

  for (abi::isa::Flag flag = abi::isa::ZERO; flag != abi::isa::CARRY; flag =
  static_cast<abi::isa::Flag>(flag + 1)) { ss << FlagConverter::to_string(flag)
  << ": " << reg_file->flag(flag) << std::endl;
  }
  */

  return ss.str();
}

}  // namespace upmem_sim::converter
