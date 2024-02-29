import pytest

from abi.binary.liveness import Liveness


@pytest.fixture
def liveness() -> Liveness:
    return Liveness()


def test_def(liveness: Liveness):
    liveness.checkout_def("foo")

    assert "foo" in liveness.defs()


def test_use(liveness: Liveness):
    liveness.checkout_def("foo")
    liveness.add_use("bar")

    assert "bar" in liveness.uses("foo")


def test_global_symbol(liveness: Liveness):
    liveness.add_global_symbol("foo")

    assert "foo" in liveness.global_symbols()


def test_local_symbol(liveness: Liveness):
    liveness.checkout_def("foo")

    assert "foo" in liveness.local_symbols()


def test_unresolved_symbol(liveness: Liveness):
    liveness.checkout_def("foo")
    liveness.add_use("bar")

    assert "bar" in liveness.unresolved_symbols()
