from iss.cpu.init_thread import InitThread
from iss.dpu.dpu import DPU
from iss.dpu.thread import Thread


class FiniThread:
    def __init__(self, benchmark: str, num_tasklets: int, dpu: DPU):
        self._benchmark: str = benchmark
        self._num_tasklets: int = num_tasklets
        self._dpu: DPU = dpu

    def fini(self) -> None:
        pass

    def cycle(self) -> None:
        sys_end_pointer = InitThread.load_sys_end_pointer(self._benchmark, self._num_tasklets)
        for thread in self._dpu.threads():
            if thread.register_file().read_pc() == sys_end_pointer and thread.state() == Thread.State.SLEEP:
                thread.set_thread_state(Thread.State.ZOMBIE)
