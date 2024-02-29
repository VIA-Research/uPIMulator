from iss.cpu.fini_thread import FiniThread
from iss.cpu.init_thread import InitThread
from iss.cpu.sched_thread import SchedThread
from iss.dpu.dpu import DPU


class CPU:
    def __init__(self, benchmark: str, num_tasklets: int, dpu: DPU):
        self._init_thread: InitThread = InitThread(benchmark, num_tasklets, dpu)
        self._sched_thread: SchedThread = SchedThread(benchmark, num_tasklets, dpu)
        self._fini_thread: FiniThread = FiniThread(benchmark, num_tasklets, dpu)

    def num_executions(self) -> int:
        return self._sched_thread.num_executions()
        
    def init(self) -> None:
        self._init_thread.init()

    def launch(self) -> None:
        self._init_thread.launch()

    def sched(self, execution: int) -> None:
        self._sched_thread.sched(execution)

    def check(self, execution: int) -> None:
        self._sched_thread.check(execution)

    def fini(self) -> None:
        self._fini_thread.fini()

    def cycle(self) -> None:
        self._init_thread.cycle()
        self._sched_thread.cycle()
        self._fini_thread.cycle()
