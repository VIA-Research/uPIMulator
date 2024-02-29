#include "util/argument_parser.h"

#include <cassert>

namespace upmem_sim::util {

void ArgumentParser::add_option(Option option, OptionType option_type,
                                Parameter default_parameter) {
  assert(not option_types_.count(option));
  assert(not default_parameters_.count(option));
  assert(not custom_parameters_.count(option));

  option_types_[option] = option_type;
  default_parameters_[option] = default_parameter;
}

void ArgumentParser::parse(int argc, char **argv) {
  assert(argc % 2 == 1);

  for (int i = 1; i < argc; i += 2) {
    for (auto &[option, _] : option_types_) {
      std::string argv_option = std::string(argv[i]);

      assert(argv_option.substr(0, 2) == "--");

      if (option == argv_option.substr(2)) {
        assert(not custom_parameters_.count(option));

        std::string argv_parameter = std::string(argv[i + 1]);
        custom_parameters_[option] = argv_parameter;
      }
    }
  }
}

std::set<ArgumentParser::Option> ArgumentParser::options() {
  std::set<Option> options;
  for (auto &[option, _] : option_types_) {
    options.insert(option);
  }
  return std::move(options);
}

std::string ArgumentParser::get_string_parameter(Option option) {
  assert(option_types_[option] == STRING);

  if (custom_parameters_.count(option)) {
    return custom_parameters_[option];
  } else {
    return default_parameters_[option];
  }
}

int64_t ArgumentParser::get_int_parameter(Option option) {
  assert(option_types_[option] == INT);

  if (custom_parameters_.count(option)) {
    return std::stoi(custom_parameters_[option]);
  } else {
    return std::stoi(default_parameters_[option]);
  }
}

}  // namespace upmem_sim::util
