import os
import subprocess
from typing import List

from util.path_collector import PathCollector


class DockerClient:
    def __init__(self):
        pass

    @staticmethod
    def repo_name() -> str:
        return "bongjoonhyun"

    @staticmethod
    def build(image: str) -> bool:
        result = subprocess.run(
            [
                "docker",
                "build",
                "-t",
                f"{DockerClient.repo_name()}/{image}",
                "-f",
                f"{DockerClient._dockerfile(image)}",
                ".",
            ],
            capture_output=True,
            text=True,
            check=True,
        )

        if result.stderr != "":
            print(result.stderr)

        return result.stderr == ""

    @staticmethod
    def build_all() -> bool:
        for dockerfile in DockerClient._dockerfiles_path():
            if not DockerClient.build(DockerClient._image(dockerfile)):
                return False
        return True

    @staticmethod
    def run(image: str, commands: List[str]) -> bool:
        DockerClient.build(image)

        result = subprocess.run(
            [
                "docker",
                "run",
                "--rm",
                "-v",
                f"{PathCollector().root_path_in_local()}:/root/{PathCollector().project_name()}",
                f"{DockerClient.repo_name()}/{image}:latest",
                "bash",
                "-c",
                " && ".join(commands),
            ],
            capture_output=True,
            text=True,
            check=True,
        )

        if result.stderr != "":
            print(result.stderr)

        return result.stderr == ""

    @staticmethod
    def _docker_path() -> str:
        return os.path.join(PathCollector().root_path_in_local(), "docker")

    @staticmethod
    def _dockerfiles_path() -> List[str]:
        return [os.path.join(DockerClient._docker_path(), f) for f in os.listdir(DockerClient._docker_path())]

    @staticmethod
    def _image(dockerfile: str) -> str:
        return os.path.basename(dockerfile).split(".")[0]

    @staticmethod
    def _dockerfile(image: str) -> str:
        return os.path.join(DockerClient._docker_path(), f"{image}.dockerfile")
