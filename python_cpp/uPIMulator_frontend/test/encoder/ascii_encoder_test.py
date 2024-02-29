from encoder.ascii_encoder import AsciiEncoder
from initializer.int_initializer import IntInitializer
from initializer.str_initializer import StrInitializer


def test_identifier():
    for _ in range(100):
        width = IntInitializer.value_by_range(1, 64)
        identifier = StrInitializer.identifier(width)

        bytes_ = [AsciiEncoder.encode(character) for character in identifier]

        assert identifier == "".join([AsciiEncoder.decode(ascii_code) for ascii_code in bytes_])


def test_hello_world():
    hello_world = "Hello, World!"

    bytes_ = [AsciiEncoder.encode(character) for character in hello_world]

    assert hello_world == "".join([AsciiEncoder.decode(ascii_code) for ascii_code in bytes_])
