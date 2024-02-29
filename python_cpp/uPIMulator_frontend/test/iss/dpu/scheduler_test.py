from typing import List

import pytest

from iss.dpu.scheduler import Scheduler
from iss.dpu.thread import Thread
from util.config_loader import ConfigLoader


@pytest.fixture
def threads() -> List[Thread]:
    return [Thread(i) for i in range(ConfigLoader.max_num_tasklets())]


def test_single_thread():
    thread = Thread(0)
    thread.set_thread_state(Thread.State.RUNNABLE)
    scheduler = Scheduler([thread])

    for i in range(100):
        assert scheduler.schedule() == thread
        scheduler.cycle()


def test_max_threads(threads: List[Thread]):
    scheduler = Scheduler(threads)
    for thread in threads:
        thread.set_thread_state(Thread.State.RUNNABLE)

    for i in range(100):
        assert scheduler.schedule() == threads[i % len(threads)]
        scheduler.cycle()


def test_awake():
    thread = Thread(0)
    thread.set_thread_state(Thread.State.RUNNABLE)
    thread.set_thread_state(Thread.State.SLEEP)
    scheduler = Scheduler([thread])

    for _ in range(100):
        assert scheduler.schedule() is None
        scheduler.cycle()

    scheduler.awake(0)
    assert thread.state() == Thread.State.RUNNABLE
