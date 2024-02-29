#include "converter/suffix_converter.h"

#include <stdexcept>

namespace upmem_profiler::converter {

abi::instruction::Suffix SuffixConverter::to_suffix(std::string suffix) {
  if (suffix == "rici") {
    return abi::instruction::RICI;
  } else if (suffix == "rri") {
    return abi::instruction::RRI;
  } else if (suffix == "rric") {
    return abi::instruction::RRIC;
  } else if (suffix == "rrici") {
    return abi::instruction::RRICI;
  } else if (suffix == "rrif") {
    return abi::instruction::RRIF;
  } else if (suffix == "rrr") {
    return abi::instruction::RRR;
  } else if (suffix == "rrrc") {
    return abi::instruction::RRRC;
  } else if (suffix == "rrrci") {
    return abi::instruction::RRRCI;
  } else if (suffix == "zri") {
    return abi::instruction::ZRI;
  } else if (suffix == "zric") {
    return abi::instruction::ZRIC;
  } else if (suffix == "zrici") {
    return abi::instruction::ZRICI;
  } else if (suffix == "zrif") {
    return abi::instruction::ZRIF;
  } else if (suffix == "zrr") {
    return abi::instruction::ZRR;
  } else if (suffix == "zrrc") {
    return abi::instruction::ZRRC;
  } else if (suffix == "zrrci") {
    return abi::instruction::ZRRCI;
  } else if (suffix == "s_rri") {
    return abi::instruction::S_RRI;
  } else if (suffix == "s_rric") {
    return abi::instruction::S_RRIC;
  } else if (suffix == "s_rrici") {
    return abi::instruction::S_RRICI;
  } else if (suffix == "s_rrif") {
    return abi::instruction::S_RRIF;
  } else if (suffix == "s_rrr") {
    return abi::instruction::S_RRR;
  } else if (suffix == "s_rrrc") {
    return abi::instruction::S_RRRC;
  } else if (suffix == "s_rrrci") {
    return abi::instruction::S_RRRCI;
  } else if (suffix == "u_rri") {
    return abi::instruction::U_RRI;
  } else if (suffix == "u_rric") {
    return abi::instruction::U_RRIC;
  } else if (suffix == "u_rrici") {
    return abi::instruction::U_RRICI;
  } else if (suffix == "u_rrif") {
    return abi::instruction::U_RRIF;
  } else if (suffix == "u_rrr") {
    return abi::instruction::U_RRR;
  } else if (suffix == "u_rrrc") {
    return abi::instruction::U_RRRC;
  } else if (suffix == "u_rrrci") {
    return abi::instruction::U_RRRCI;
  } else if (suffix == "rr") {
    return abi::instruction::RR;
  } else if (suffix == "rrc") {
    return abi::instruction::RRC;
  } else if (suffix == "rrci") {
    return abi::instruction::RRCI;
  } else if (suffix == "zr") {
    return abi::instruction::ZR;
  } else if (suffix == "zrc") {
    return abi::instruction::ZRC;
  } else if (suffix == "zrci") {
    return abi::instruction::ZRCI;
  } else if (suffix == "s_rr") {
    return abi::instruction::S_RR;
  } else if (suffix == "s_rrc") {
    return abi::instruction::S_RRC;
  } else if (suffix == "s_rrci") {
    return abi::instruction::S_RRCI;
  } else if (suffix == "u_rr") {
    return abi::instruction::U_RR;
  } else if (suffix == "u_rrc") {
    return abi::instruction::U_RRC;
  } else if (suffix == "u_rrci") {
    return abi::instruction::U_RRCI;
  } else if (suffix == "drdici") {
    return abi::instruction::DRDICI;
  } else if (suffix == "rrri") {
    return abi::instruction::RRRI;
  } else if (suffix == "rrrici") {
    return abi::instruction::RRRICI;
  } else if (suffix == "zrri") {
    return abi::instruction::ZRRI;
  } else if (suffix == "zrrici") {
    return abi::instruction::ZRRICI;
  } else if (suffix == "s_rrri") {
    return abi::instruction::S_RRRI;
  } else if (suffix == "s_rrrici") {
    return abi::instruction::S_RRRICI;
  } else if (suffix == "u_rrri") {
    return abi::instruction::U_RRRI;
  } else if (suffix == "u_rrrici") {
    return abi::instruction::U_RRRICI;
  } else if (suffix == "rir") {
    return abi::instruction::RIR;
  } else if (suffix == "rirc") {
    return abi::instruction::RIRC;
  } else if (suffix == "rirci") {
    return abi::instruction::RIRCI;
  } else if (suffix == "zir") {
    return abi::instruction::ZIR;
  } else if (suffix == "zirc") {
    return abi::instruction::ZIRC;
  } else if (suffix == "zirci") {
    return abi::instruction::ZIRCI;
  } else if (suffix == "s_rirc") {
    return abi::instruction::S_RIRC;
  } else if (suffix == "s_rirci") {
    return abi::instruction::S_RIRCI;
  } else if (suffix == "u_rirc") {
    return abi::instruction::U_RIRC;
  } else if (suffix == "u_rirci") {
    return abi::instruction::U_RIRCI;
  } else if (suffix == "r") {
    return abi::instruction::R;
  } else if (suffix == "rci") {
    return abi::instruction::RCI;
  } else if (suffix == "z") {
    return abi::instruction::Z;
  } else if (suffix == "zci") {
    return abi::instruction::ZCI;
  } else if (suffix == "s_r") {
    return abi::instruction::S_R;
  } else if (suffix == "s_rci") {
    return abi::instruction::S_RCI;
  } else if (suffix == "u_r") {
    return abi::instruction::U_R;
  } else if (suffix == "u_rci") {
    return abi::instruction::U_RCI;
  } else if (suffix == "ci") {
    return abi::instruction::CI;
  } else if (suffix == "i") {
    return abi::instruction::I;
  } else if (suffix == "ddci") {
    return abi::instruction::DDCI;
  } else if (suffix == "erri") {
    return abi::instruction::ERRI;
  } else if (suffix == "s_erri") {
    return abi::instruction::S_ERRI;
  } else if (suffix == "u_erri") {
    return abi::instruction::U_ERRI;
  } else if (suffix == "edri") {
    return abi::instruction::EDRI;
  } else if (suffix == "erii") {
    return abi::instruction::ERII;
  } else if (suffix == "erir") {
    return abi::instruction::ERIR;
  } else if (suffix == "erid") {
    return abi::instruction::ERID;
  } else if (suffix == "dma_rri") {
    return abi::instruction::DMA_RRI;
  } else {
    throw std::invalid_argument("");
  }
}

} // namespace upmem_profiler::converter
