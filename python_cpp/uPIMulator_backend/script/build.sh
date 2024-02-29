#!/bin/bash

mkdir -p ../build
cd ../build || exit
# shellcheck disable=SC2035
rm -rf *
cmake ..
make -j
cd - || exit