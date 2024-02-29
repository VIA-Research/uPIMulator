from typing import List

from iss.dpu.dispatcher import Dispatcher
from iss.dpu.dma import DMA
from iss.dpu.logic import Logic
from iss.dpu.scheduler import Scheduler
from iss.dpu.thread import Thread
from iss.dram.mram import MRAM
from iss.sram.atomic import Atomic
from iss.sram.iram import IRAM
from iss.sram.wram import WRAM
from util.config_loader import ConfigLoader


class DPU:
    def __init__(self, num_threads: int):
        assert 0 < num_threads <= ConfigLoader.max_num_tasklets()

        self._threads: List[Thread] = [Thread(i) for i in range(num_threads)]
        self._scheduler: Scheduler = Scheduler(self._threads)

        self._atomic: Atomic = Atomic()
        self._iram: IRAM = IRAM()
        self._wram: WRAM = WRAM()
        self._mram: MRAM = MRAM()
        self._dma: DMA = DMA(self._atomic, self._iram, self._wram, self._mram)
        self._dispatcher: Dispatcher = Dispatcher(self._wram)

        self._logic = Logic(
            self._scheduler, self._atomic, self._iram, self._wram, self._mram, self._dma, self._dispatcher,
        )

    def threads(self) -> List[Thread]:
        return self._threads

    def atomic(self) -> Atomic:
        return self._atomic

    def iram(self) -> IRAM:
        return self._iram

    def wram(self) -> WRAM:
        return self._wram

    def mram(self) -> MRAM:
        return self._mram

    def dma(self) -> DMA:
        return self._dma

    def is_zombie(self) -> bool:
        for thread in self._threads:
            if thread.state() != Thread.State.ZOMBIE:
                return False
        return True

    def boot(self) -> None:
        self._scheduler.boot(0)

    def cycle(self) -> None:
        self._logic.cycle()
        self._scheduler.cycle()
