import argparse
import os

from compiler.compiler import Compiler
from iss.system import System
from linker_.linker import Linker
from util.path_collector import PathCollector


def compile(benchmark: str, num_tasklets: int) -> None:
    Compiler.clean()
    Compiler.compile_sdk(num_tasklets)
    Compiler.compile_benchmark(benchmark, num_tasklets)


def link(benchmark: str, num_dpus: int,num_tasklets: int, data_prep_param: list) -> None:
    linker = Linker(num_tasklets)
    linker.link(os.path.join(PathCollector.asm_path_in_local(), f"{benchmark}.{num_tasklets}", "main.S"), data_prep_param, num_dpus)


def iss(benchmark: str, num_tasklets: int) -> None:
    system = System(benchmark, num_tasklets)

    system.init()
    while not system.is_finished():
        system.cycle()
    system.fini()


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--benchmark", type=str, default="SEL")
    parser.add_argument("--num_tasklets", type=int, default=1)
    parser.add_argument("--mode", type=str, default="all")
    parser.add_argument("--data_prep_param", type=str, default="2048")
    parser.add_argument("--num_dpus", type=int, default="16")
    args = parser.parse_args()

    data_prep_param = [int(elem) for elem in args.data_prep_param.split(',')]

    if args.mode == "compile":
        compile(args.benchmark, args.num_tasklets)
    elif args.mode == "link":
        link(args.benchmark, args.num_dpus, args.num_tasklets, data_prep_param)
    elif args.mode == "iss":
        iss(args.benchmark, args.num_tasklets)
    elif args.mode == "all":
        compile(args.benchmark, args.num_tasklets)
        link(args.benchmark,  args.num_dpus, args.num_tasklets, data_prep_param)
        iss(args.benchmark, args.num_tasklets)
    else:
        raise ValueError
