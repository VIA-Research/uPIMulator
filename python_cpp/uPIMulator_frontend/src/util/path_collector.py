import os
from pathlib import Path


class PathCollector:
    def __init__(self):
        pass

    @staticmethod
    def project_name() -> str:
        return "uPIMulator_frontend"

    @staticmethod
    def root_path_in_local() -> str:
        dirpath = Path(__file__).resolve().parent.as_posix()
        dirnames = dirpath.split(os.path.sep)
        for i in range(len(dirnames)):
            if dirnames[i] == PathCollector.project_name():
                return os.path.sep.join(dirnames[: i + 1])
        raise ValueError

    @staticmethod
    def root_path_in_docker() -> str:
        return f"/root/{PathCollector.project_name()}"

    @staticmethod
    def asm_path_in_local() -> str:
        return os.path.join(PathCollector.root_path_in_local(), "asm")

    @staticmethod
    def asm_path_in_docker() -> str:
        return f"{PathCollector.root_path_in_docker()}/asm"

    @staticmethod
    def benchmark_path_in_local() -> str:
        return os.path.join(PathCollector.root_path_in_local(), "benchmark")

    @staticmethod
    def benchmark_path_in_docker() -> str:
        return f"{PathCollector.root_path_in_docker()}/benchmark"

    @staticmethod
    def bin_path_in_local() -> str:
        return os.path.join(PathCollector.root_path_in_local(), "bin")

    @staticmethod
    def bin_path_in_docker() -> str:
        return f"{PathCollector.root_path_in_docker()}/bin"

    @staticmethod
    def trace_path_in_local() -> str:
        return os.path.join(PathCollector.root_path_in_local(), "trace")

    @staticmethod
    def trace_path_in_docker() -> str:
        return f"{PathCollector.root_path_in_docker()}/trace"

    @staticmethod
    def sdk_path_in_local() -> str:
        return os.path.join(PathCollector.root_path_in_local(), "sdk")

    @staticmethod
    def sdk_path_in_docker() -> str:
        return f"{PathCollector.root_path_in_docker()}/sdk"

    @staticmethod
    def src_path_in_local() -> str:
        return os.path.join(PathCollector.root_path_in_local(), "src")

    @staticmethod
    def src_path_in_docker() -> str:
        return f"{PathCollector.root_path_in_docker()}/src"

    @staticmethod
    def upmem_sdk_path_in_local() -> str:
        raise ValueError

    @staticmethod
    def upmem_sdk_path_in_docker() -> str:
        return "/root/upmem-2021.3.0-Linux-x86_64"
