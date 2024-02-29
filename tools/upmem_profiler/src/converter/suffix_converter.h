#ifndef UPMEM_PROFILER_CONVERTER_SUFFIX_CONVERTER_H_
#define UPMEM_PROFILER_CONVERTER_SUFFIX_CONVERTER_H_

#include <string>

#include "abi/instruction/suffix.h"

namespace upmem_profiler::converter {

class SuffixConverter {
public:
  static abi::instruction::Suffix to_suffix(std::string suffix);
};

} // namespace upmem_profiler::converter

#endif
