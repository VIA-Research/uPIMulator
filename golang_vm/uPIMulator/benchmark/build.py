import os
import shutil
import subprocess
import argparse


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--num_dpus", type=int, default=1)
    parser.add_argument("--num_tasklets", type=int, default=1)
    args = parser.parse_args()

    benchmark_dir_path = os.path.dirname(__file__)

    build_dir_path = os.path.join(benchmark_dir_path, "build")

    if os.path.exists(build_dir_path):
        shutil.rmtree(build_dir_path)
    os.makedirs(build_dir_path)

    subprocess.run(
        [
            "cmake",
            "-D",
            f"NR_DPUS={args.num_dpus}",
            "-D"
            f"NR_TASKLETS={args.num_tasklets}",
            "-S",
            benchmark_dir_path,
            "-B",
            build_dir_path,
            "-G",
            "Ninja",
        ]
    )
    subprocess.run(["ninja", "-C", build_dir_path])
