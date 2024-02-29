from typing import List

from util.docker_client import DockerClient
from util.path_collector import PathCollector


class GrammarGenerator:
    def __init__(self):
        pass

    @staticmethod
    def generate() -> bool:
        cd = f"cd {GrammarGenerator._grammar_path_in_docker()}"
        antlr = f"{GrammarGenerator._antlr()} -Dlanguage=Python3 {GrammarGenerator.grammar()}.g4"

        return DockerClient.run(GrammarGenerator._docker_image(), [cd, antlr])

    @staticmethod
    def clean() -> bool:
        commands = [f"cd {GrammarGenerator._grammar_path_in_docker()}"]
        for filename in GrammarGenerator.generated_filenames():
            commands.append(f"rm -f {filename}")
        return DockerClient.run("compiler", commands)

    @staticmethod
    def grammar() -> str:
        return "assembly"

    @staticmethod
    def generated_filenames() -> List[str]:
        return [
            f"{GrammarGenerator.grammar()}.interp",
            f"{GrammarGenerator.grammar()}.tokens",
            f"{GrammarGenerator.grammar()}Lexer.interp",
            f"{GrammarGenerator.grammar()}Lexer.py",
            f"{GrammarGenerator.grammar()}Lexer.tokens",
            f"{GrammarGenerator.grammar()}Listener.py",
            f"{GrammarGenerator.grammar()}Parser.py",
        ]

    @staticmethod
    def _grammar_path_in_docker() -> str:
        return f"{PathCollector.src_path_in_docker()}/parser_/grammar"

    @staticmethod
    def _docker_image() -> str:
        return "parser"

    @staticmethod
    def _class_path() -> str:
        return "/root/antlr-4.9.2-complete.jar:$CLASSPATH"

    @staticmethod
    def _antlr() -> str:
        return f'java -Xmx500M -cp "{GrammarGenerator._class_path()}" org.antlr.v4.Tool'

    @staticmethod
    def _grun() -> str:
        return f'java -Xmx500M -cp "{GrammarGenerator._class_path()}" org.antlr.v4.gui.TestRig'
