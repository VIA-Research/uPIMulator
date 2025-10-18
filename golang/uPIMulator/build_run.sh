#!/bin/bash
rm -rf bin
mkdir bin
pushd script || exit 1
python3 build.py || exit 1
popd || exit 1
bash run.sh
