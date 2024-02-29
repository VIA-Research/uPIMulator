import pytest

from initializer.int_initializer import IntInitializer
from iss.sram.atomic import Atomic
from util.config_loader import ConfigLoader


@pytest.fixture
def atomic() -> Atomic:
    return Atomic()


def test_atomic(atomic: Atomic):
    for _ in range(100):
        address = IntInitializer.value_by_range(
            ConfigLoader.atomic_offset(), ConfigLoader.atomic_offset() + ConfigLoader.atomic_size() - 1,
        )

        id_ = IntInitializer.value_by_range(0, ConfigLoader.max_num_tasklets() - 1)

        assert atomic.can_acquire(address)
        atomic.acquire(address, id_)
        assert not atomic.can_acquire(address)

        assert atomic.can_release(address, id_)
        atomic.release(address, id_)
        assert atomic.can_acquire(address)
