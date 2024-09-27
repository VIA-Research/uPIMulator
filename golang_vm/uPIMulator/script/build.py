import os
import shutil
import subprocess


if __name__ == "__main__":
    script_dirpath = os.path.dirname(__file__)

    build_dirpath = os.path.join(script_dirpath, "..", "build")
    src_dirpath = os.path.join(script_dirpath, "..", "src")

    if os.path.exists(build_dirpath):
        shutil.rmtree(build_dirpath)
    os.makedirs(build_dirpath)

    binary_filepath = os.path.join(build_dirpath, "uPIMulator")

    subprocess.run(["go", "build", "-C", src_dirpath, "-o", binary_filepath])
