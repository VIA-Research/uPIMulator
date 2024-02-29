import os

from parser_.grammar_generator import GrammarGenerator
from util.path_collector import PathCollector


def test_generate():
    assert GrammarGenerator.generate()

    for filename in GrammarGenerator.generated_filenames():
        assert os.path.exists(os.path.join(PathCollector.src_path_in_local(), "parser_", "grammar", filename))
