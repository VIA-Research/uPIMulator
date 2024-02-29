import os

import pytest

from compiler.compiler import Compiler
from util.path_collector import PathCollector


# TODO(bongjoon.hyun@gmail.com): compile for num_tasklets from 1 to 24
@pytest.fixture
def num_tasklets() -> int:
    return 1


def test_compile_benchmarks(num_tasklets: int):
    for benchmark in os.listdir(PathCollector.benchmark_path_in_local()):
        assert Compiler.compile_benchmark(benchmark, num_tasklets)


def test_compile_sdk(num_tasklets: int):
    assert Compiler.compile_sdk(num_tasklets)
