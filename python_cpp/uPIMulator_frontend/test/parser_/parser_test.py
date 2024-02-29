import os
import subprocess
from typing import Set

import pytest

from parser_.grammar_generator import GrammarGenerator
from parser_.parser import Parser
from util.path_collector import PathCollector


@pytest.fixture
def parser_path() -> str:
    return os.path.join(PathCollector().src_path_in_local(), "parser_", "parser.py",)


@pytest.fixture
def asm_filepaths() -> Set[str]:
    filepaths: Set[str] = set()
    for root_path, _, filenames in os.walk(PathCollector.asm_path_in_local()):
        for filename in filenames:
            if filename.split(".")[-1] == "S":
                filepaths.add(os.path.join(root_path, filename))
    return filepaths


@pytest.fixture
def bin_filepaths() -> Set[str]:
    filepaths: Set[str] = set()
    for root_path, _, filenames in os.walk(PathCollector.bin_path_in_local()):
        for filename in filenames:
            if filename.split(".")[-1] == "S":
                filepaths.add(os.path.join(root_path, filename))
    return filepaths


@pytest.fixture
def num_tasklets() -> int:
    return 16


def test_parse_asm(parser_path: str, asm_filepaths: Set[str]):
    assert GrammarGenerator.clean()
    assert GrammarGenerator.generate()

    for asm_filepath in asm_filepaths:
        commands = ["python", f"{parser_path}", "--file", f"{asm_filepath}"]

        result = subprocess.run(commands, capture_output=True, text=True,)

        if result.stderr != "":
            with open(asm_filepath, encoding="ISO-8859-1") as file:
                for line in file.readlines():
                    line = Parser.preprocess(line)

                    commands = ["python", f"{parser_path}", "--line", f"{line}"]

                    result = subprocess.run(commands, capture_output=True, text=True,)

                    print(" ".join(commands))
                    assert result.stderr == ""


def test_parse_bin(parser_path: str, bin_filepaths: Set[str]):
    assert GrammarGenerator.clean()
    assert GrammarGenerator.generate()

    for bin_filepath in bin_filepaths:
        commands = ["python", f"{parser_path}", "--file", f"{bin_filepath}"]

        result = subprocess.run(commands, capture_output=True, text=True,)

        if result.stderr != "":
            with open(bin_filepath, encoding="ISO-8859-1") as file:
                for line in file.readlines():
                    line = Parser.preprocess(line)

                    commands = ["python", f"{parser_path}", "--line", f"{line}"]

                    result = subprocess.run(commands, capture_output=True, text=True,)

                    print(" ".join(commands))
                    assert result.stderr == ""
