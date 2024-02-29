#ifndef UPMEM_SIM_UTIL_ARGUMENT_PARSER_H_
#define UPMEM_SIM_UTIL_ARGUMENT_PARSER_H_

#include <map>
#include <set>
#include <string>

namespace upmem_sim::util {

class ArgumentParser {
 public:
  using Option = std::string;
  using Parameter = std::string;

  enum OptionType { STRING = 0, INT };

  explicit ArgumentParser() = default;
  ~ArgumentParser() = default;

  void add_option(Option option, OptionType option_type,
                  Parameter default_parameter);
  void parse(int argc, char **argv);

  std::set<Option> options();
  OptionType option_type(Option option) { return option_types_[option]; }

  std::string get_string_parameter(Option option);
  int64_t get_int_parameter(Option option);

 private:
  std::map<Option, OptionType> option_types_;
  std::map<Option, Parameter> default_parameters_;
  std::map<Option, Parameter> custom_parameters_;
};

}  // namespace upmem_sim::util

#endif
