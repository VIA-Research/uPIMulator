import os
from typing import Set

import pytest

from linker_.linker import Linker
from util.path_collector import PathCollector


@pytest.fixture
def asm_benchmark_filepaths() -> Set[str]:
    benchmarks = os.listdir(PathCollector.benchmark_path_in_local())
    filepaths: Set[str] = set()
    for root_path, _, filenames in os.walk(PathCollector.asm_path_in_local()):
        benchmark = root_path.split(os.path.sep)[-1].split(".")[0]

        if benchmark in benchmarks:
            for filename in filenames:
                if filename.split(".")[-1] == "S":
                    filepaths.add(os.path.join(root_path, filename))
    return filepaths


@pytest.fixture
def num_tasklets() -> int:
    return 1


# TODO(bongjoon.hyun@gmail.com): instantiate linker based on num_tasklets
def test_link(asm_benchmark_filepaths: Set[str], num_tasklets: int):
    linker = Linker(num_tasklets)
    for asm_benchmark_filepath in asm_benchmark_filepaths:
        linker.link(asm_benchmark_filepath)
