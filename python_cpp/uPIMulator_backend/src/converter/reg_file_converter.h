#ifndef UPMEM_SIM_CONVERTER_REG_FILE_CONVERTER_H_
#define UPMEM_SIM_CONVERTER_REG_FILE_CONVERTER_H_

#include <string>

#include "simulator/reg/reg_file.h"

namespace upmem_sim::converter {

class RegFileConverter {
 public:
  static std::string to_string(simulator::reg::RegFile *reg_file);
};

}  // namespace upmem_sim::converter

#endif
