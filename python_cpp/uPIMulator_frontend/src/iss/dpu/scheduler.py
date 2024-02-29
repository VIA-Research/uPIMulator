from queue import Queue
from typing import List, Optional

from iss.dpu.thread import Thread


class Scheduler:
    def __init__(self, threads: List[Thread]):
        self._threads = threads
        self._queue: Queue[Thread] = self._init_queue()

    def schedule(self) -> Optional[Thread]:
        for _ in range(len(self._threads)):
            thread = self._queue.get()
            self._queue.put(thread)

            if thread.state() == Thread.State.RUNNABLE:
                return thread
        return None

    def boot(self, id_: int) -> bool:
        if id_ < len(self._threads):
            assert self._threads[id_].id_() == id_

            if self._threads[id_].state() == Thread.State.EMBRYO:
                self._threads[id_].set_thread_state(Thread.State.RUNNABLE)
                return True
            elif self._threads[id_].state() == Thread.State.ZOMBIE:
                self._threads[id_].set_thread_state(Thread.State.RUNNABLE)
                return True
            else:
                raise ValueError
        else:
            return True

    def awake(self, id_: int) -> bool:
        assert self._threads[id_].id_() == id_

        if self._threads[id_].state() == Thread.State.RUNNABLE:
            return True
        elif self._threads[id_].state() == Thread.State.SLEEP:
            self._threads[id_].set_thread_state(Thread.State.RUNNABLE)
            return True
        else:
            raise ValueError

    def cycle(self) -> None:
        assert self._queue.full()

    def _init_queue(self) -> Queue[Thread]:
        queue: Queue[Thread] = Queue(len(self._threads))
        for thread in self._threads:
            queue.put(thread)

        return queue
