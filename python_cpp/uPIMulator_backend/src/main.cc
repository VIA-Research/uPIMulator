#include <iostream>

#include "simulator/system.h"
#include "util/argument_parser.h"

namespace upmem_sim {

util::ArgumentParser *init_argument_parser() {
  auto argument_parser = new util::ArgumentParser();

  // NOTE(dongjae.lee@kaist.ac.kr): Explanation of verbose level
  // level 0: Only prints simulation output
  // level 1: level 0 + prints UPMEM instruction executed per each logic cycle
  // level 2: level + prints UPMEM register file values per each logic cycle
  argument_parser->add_option("verbose", util::ArgumentParser::INT, "0");

  argument_parser->add_option("benchmark", util::ArgumentParser::STRING, "TRNS");
  argument_parser->add_option("num_dpus", util::ArgumentParser::INT, "1");
  argument_parser->add_option("num_tasklets", util::ArgumentParser::INT, "16");

  argument_parser->add_option("bindir", util::ArgumentParser::STRING,
                              "/home/via/uPIMulator_frontend/bin");
  argument_parser->add_option("logdir", util::ArgumentParser::STRING,
                              "/home/via/uPIMulator_backend/log");

  argument_parser->add_option("logic_frequency", util::ArgumentParser::INT,
                              "350");
  argument_parser->add_option("memory_frequency", util::ArgumentParser::INT,
                              "2400");  // based on DDR4-2400

  argument_parser->add_option("num_pipeline_stages", util::ArgumentParser::INT,
                              "14");
  argument_parser->add_option("num_revolver_scheduling_cycles",
                              util::ArgumentParser::INT, "11");

  argument_parser->add_option("wordline_size", util::ArgumentParser::INT,
                              "1024");

  argument_parser->add_option("t_rcd", util::ArgumentParser::INT,
                              "32");  // based on DDR4-2400
  argument_parser->add_option("t_ras", util::ArgumentParser::INT,
                              "78");  // based on DDR4-2400
  argument_parser->add_option("t_rp", util::ArgumentParser::INT,
                              "32");  // based on DDR4-2400
  argument_parser->add_option("t_cl", util::ArgumentParser::INT,
                              "32");  // based on DDR4-2400
  argument_parser->add_option("t_bl", util::ArgumentParser::INT,
                              "8");  // based on DDR4-2400

  argument_parser->add_option("memory_scheduling_policy",
                              util::ArgumentParser::STRING, "frfcfs");

  argument_parser->add_option("rank_read_bandwidth", util::ArgumentParser::INT,
                              "1"); //1
  argument_parser->add_option("rank_write_bandwidth", util::ArgumentParser::INT,
                              "3"); //3

  return argument_parser;
}

}  // namespace upmem_sim

int main(int argc, char **argv) {
  upmem_sim::util::ArgumentParser *argument_parser =
      upmem_sim::init_argument_parser();
  argument_parser->parse(argc, argv);

  auto system = new upmem_sim::simulator::System(argument_parser);
  system->init();
  while (not system->is_finished()) {
    system->cycle();
  }
  system->fini();

  for (auto &option : argument_parser->options()) {
    if (argument_parser->option_type(option) ==
        upmem_sim::util::ArgumentParser::INT) {
      std::cout << option << ": " << argument_parser->get_int_parameter(option)
                << std::endl;
    } else if (argument_parser->option_type(option) ==
               upmem_sim::util::ArgumentParser::STRING) {
        std::cout << option << ": "
                << argument_parser->get_string_parameter(option) << std::endl;
    } else {
      throw std::invalid_argument("");
    }
  }

  upmem_sim::util::StatFactory *system_stat_factory = system->stat_factory();
  for (auto &stat : system_stat_factory->stats()) {
      std::cout << stat << ": " << system_stat_factory->value(stat) << std::endl;
  }
  delete system_stat_factory;

  delete argument_parser;
  delete system;

  return 0;
}
