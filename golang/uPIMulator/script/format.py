import os
import subprocess


if __name__ == "__main__":
    script_dirpath = os.path.dirname(__file__)

    src_dirpath = os.path.join(script_dirpath, "..", "src")

    subprocess.run(["gofmt", "-l", src_dirpath])
    subprocess.run(["golines", "-w", src_dirpath])
