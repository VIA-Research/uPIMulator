import os
from typing import Set

from util.docker_client import DockerClient
from util.path_collector import PathCollector


class Compiler:
    def __init__(self):
        pass

    @staticmethod
    def compile_benchmark(benchmark: str, num_tasklets: int):
        source = Compiler._source_command()
        cd = f"cd {PathCollector.benchmark_path_in_docker()}/{benchmark}"
        make_clean = "make clean"
        make = f"make NR_TASKLETS={num_tasklets}"
        mkdir = f"mkdir -p {PathCollector.asm_path_in_docker()}/{benchmark}.{num_tasklets}"
        mv = f"mv {PathCollector.benchmark_path_in_docker()}/{benchmark}/bin/dpu_code.S {PathCollector.asm_path_in_docker()}/{benchmark}.{num_tasklets}/main.S"
        commands = [source, cd, make_clean, make, mkdir, mv]

        return DockerClient.run(Compiler._docker_image(), commands)

    @staticmethod
    def compile_sdk(num_tasklets: int) -> bool:
        for filepath in Compiler._sdk_filepaths_in_docker():
            library_name, filename = filepath.split("/")[-2], filepath.split("/")[-1]

            common_flags = f"-O3 -S -w -DNR_TASKLETS=${num_tasklets}"
            include_flags = f"-I{Compiler._misc_path()} -I{Compiler._stdlib_path()} -I{Compiler._syslib_path()}"
            output_flag = f"-o {PathCollector.asm_path_in_docker()}/{library_name}/{filename[:-2]}.S"

            source = Compiler._source_command()
            mkdir = f"mkdir -p {PathCollector.asm_path_in_docker()}/{library_name}"
            dpu_upmem_dpu_rte_clang = (
                f"{Compiler._dpu_upmem_dpurte_clang()} {common_flags} {include_flags} {output_flag} {filepath}"
            )
            commands = [source, mkdir, dpu_upmem_dpu_rte_clang]

            if not DockerClient.run(Compiler._docker_image(), commands):
                return False
        return True

    @staticmethod
    def clean() -> bool:
        rm = f"rm -rf {PathCollector.asm_path_in_docker()}"
        return DockerClient.run(Compiler._docker_image(), [rm])

    @staticmethod
    def _sdk_filepaths_in_docker() -> Set[str]:
        filepaths: Set[str] = set()
        for root_path, _, filenames in os.walk(PathCollector.sdk_path_in_local()):
            for filename in filenames:
                if filename.split(".")[-1] == "c":
                    library_name = root_path.split(os.path.sep)[-1]
                    filepaths.add(f"{PathCollector.sdk_path_in_docker()}/{library_name}/{filename}")
        return filepaths

    @staticmethod
    def _docker_image() -> str:
        return "compiler"

    @staticmethod
    def _source_command() -> str:
        return f"source {PathCollector.upmem_sdk_path_in_docker()}/upmem_env.sh"

    @staticmethod
    def _upmem_include_path() -> str:
        return f"{PathCollector.upmem_sdk_path_in_docker()}/include"

    @staticmethod
    def _misc_path() -> str:
        return f"{PathCollector.sdk_path_in_docker()}/misc"

    @staticmethod
    def _stdlib_path() -> str:
        return f"{PathCollector.sdk_path_in_docker()}/stdlib"

    @staticmethod
    def _syslib_path() -> str:
        return f"{PathCollector.sdk_path_in_docker()}/syslib"

    @staticmethod
    def _dpu_upmem_dpurte_clang() -> str:
        return "dpu-upmem-dpurte-clang"
