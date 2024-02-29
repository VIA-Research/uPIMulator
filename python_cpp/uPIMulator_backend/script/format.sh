#!/bin/bash

cd ../lib
find . -regex '.*\.\(cpp\|hpp\|cc\|h\)' -exec clang-format -style=Google -i {} \;
cd -

cd ../src
find . -regex '.*\.\(cpp\|hpp\|cc\|h\)' -exec clang-format -style=Google -i {} \;
cd -
