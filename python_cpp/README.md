# ⚙️ Usage

uPIMulator currently supports one usage mode, which is an execution-driven mode.
The flow of using the uPIMulator framework is two-folded: 1) to compile, assemble, and link to generate binary files, and 2) feed the binary files to the cycle-level simulator to get simulation results.

# 🏃 Getting Started

## Install & Build

### Prerequisites:
- uPIMulator requires Python 3.10 interpreter and C++20 compiler.
- uPIMulator requires CMake version higher than 3.16.
- uPIMulator requires Docker.
- Your Ubuntu account must be added to the Docker group.
- uPIMulator has been tested on Ubuntu 18.04 or higher version.
- We recommend that you use [Anaconda](https://www.anaconda.com/) or a similar tool.

### Install the Linker:

```
cd /path/to/uPIMulator/python_cpp/uPIMulator_frontend
pip install -r requirements.txt
```

### Build the Cycle-level Simulator:

```
cd /path/to/uPIMulator/python_cpp/uPIMulator_backend/script
sh build.sh
```

## Binary Files Generation

We will use the VA (vector addition) benchmark as an example for the binary file generation phase.
After the compile/assemble/link phase, you will be able to find binary files under the `uPIMulator_frontend/bin` directory.

We provide pre-generated binary files at the [following link](https://drive.google.com/file/d/1a_AilWSHiSF9WLYu--GFBSyEpguFLjSe/view).

### Compile

```
cd /path/to/uPIMulator/python_cpp/uPIMulator_frontend/src
python main.py --mode compile --num_tasklets 16 --benchmark VA --num_dpus 1
```

### Assemble and Link

```
cd /path/to/uPIMulator/python_cpp/uPIMulator_frontend/src
python main.py --mode link --num_tasklets 16 --benchmark VA --data_prep_param 1024 --num_dpus 1
```

If the above command fails, type command below and repeat the commands above.
```
docker build -t bongjoonhyun/compiler -f /path/to/uPIMulator/python_cpp/uPIMulator_frontend/docker/compiler.dockerfile .
```

Currently, we support 13 PrIM benchmarks.
For the first compile/assemble/link, it will take roughly 30 minutes.

## Cycle-level Simulation

Run a simulation by giving the benchmark name, the number of tasklets, and the path to the generated binary files as inputs to the cycle-level simulator.
Tweak simulation parameters through command line options.
The simulation results are printed out to `stdout`.

- **Note that you need to specify an absolute path to `bindir`.**

```
cd /path/to/uPIMulator/uPIMulator_backend/build/
./src/uPIMulator --benchmark VA --num_tasklets 16 --bindir /path/to/uPIMulator/uPIMulator_frontend/bin/1_dpus/ --logdir .
```

# 📄 Reproducing Figures from the Paper
To replicate the figures illustrated in our paper, please follow the instructions below.
We provide replication manuals for Figures used in *Section 4. (Demystifying UPMEM-PIM with uPIMulator)*. Note that figure 8 is omitted for brevity.

## Configuration of PrIM Benchmarks
- Parameter `num_dpus` must always be configured to 1 for Figure 5, Figure 6, Figure 7, Figure 9 expriments, since they characterize the behavior of a single DPU.
- When generating the binary file of the PrIM benchmark, please configure the `data_prep_param` parameter as shown in the following table. 
(e.g., `python main.py --mode link --num_tasklets 16 --benchmark VA --data_prep_param 524288 --num_dpus 1`)

| Benchmark|data_prep_param (fig 5, 6, 7, 9) |data_prep_param (fig 10)|
|----------|---------------------------------|------------------------|
| BS       |           32768                 |       131072           |
| GEMV     |           2048                  |       4096             |
| HST-L    |           131072                |       524288           |
| HST-S    |           131072                |       524288           |
| MLP      |           256                   |       1024             |
| RED      |           524288                |       2097152          |
| SCAN-RSS |           262144                |       1048576          |
| SCAN-SSA |           262144                |       1048576          |
| SEL      |           524288                |       2097152          |
| TRNS     |           1024                  |       128              |
| TS       |           2048                  |       65536            |
| UNI      |           524288                |       2097152          |
| VA       |           524288                |       2097152          |

## Figure 5
<img src="../assets/uPIMulator_figure5.png" width="400"/>

Figure 5 depicts PrIM’s compute utilization (red points) and memory read bandwidth utilization (blue points) when executing with 1/4/16 threads (i.e., tasklets).
The calculation formula of compute utilization and memory read bandwidth utilization is as follows:

- **Compute utilization (i.e., IPC)** = (value of `num_instructions`) / (value of `logic_cycle` - value of `communication_cycle`) 
- **Memory read bandwidth utilization (GB/s)** = (We provide an excel sheet to calculate the memory read bandwidth utilization in the following [link](../assets/figure5_mem_util_calculator.xlsx))


Note that the value of `XYZ` can be obtained by the simulation results.

## Figure 6
<img src="../assets/uPIMulator_fiture6.png" width="400"/>

Figure 6 illustrates the breakdown of DPU’s runtime into active (black) and idle (red, yellow, blue) cycles.
To plot the breakdown of DPU's runtime, please use the following formulas:

- **Issuable ratio** = (value of `breakdown_run`) / (value of `logic_cycle` - value of `communication_cycle`)
- **Idle (Memory) ratio** = (value of `breakdown_dma`) / (value of `logic_cycle` - value of `communication_cycle`)
- **Idle (Revolver) ratio** = (value of `breakdown_etc` - value of `communication_cycle`) / (value of `logic_cycle` - value of `communication_cycle`)
- **Idle (RF) ratio** = (value of `backpressuer`) / (value of `logic_cycle` - value of `communication_cycle`)

Note that the value of `XYZ` can be obtained by the simulation results.

## Figure 7
<img src="../assets/uPIMulator_figure7.png" width="400"/>

Figure 7 shows the number of issuable tasklets (i.e., threads) by the DPU scheduler. To replicate the above figure, please use the following [excel sheet](../assets/figure7_active_tasklet_breakdown.xlsx) provided in our repository. By filling out the Excel spreadsheet based on the simulation output (the spreadsheet explains what kind of simulation output needs to be filled out), the spreadsheet automatically generates the figure.

**Note that the number of threads configured in this experiment is ONLY 16 tasklets.**
(e.g., `--num_tasklets 16`)

## Figure 9
<img src="../assets/uPIMulator_figure9.png" width="400"/>

Figure 9 illustrates the instruction mix when executing with a single DPU. To plot the given figure, please use both [upmem_profiler](../tools/upmem_profiler/) and [Excel sheet](../assets/figure9_instruction_mix.xlsx) and follow the procedures as explained below.

### Build Profiler

```
cd /path/to/uPIMulator/tools/upmem_profiler/script
bash build.sh
```

### Extract Instructions
When extracting instructions, you should add one additional parameter, `--verbose 1`. 
```
cd /path/to/uPIMulator/python_cpp/uPIMulator_backend/
./build/src/uPIMulator --benchmark VA --num_tasklets 16 --bindir /path/to/uPIMulator/python_cpp/uPIMulator_frontend/bin/1_dpus/ --logdir . --verbose 1 > trace.txt
```

### Run Profiler
```
cd /path/to/uPIMulator/tools/upmem_profiler/
./build/src/upmem_profiler --logpath /path/to/uPIMulator/python_cpp/uPIMulator_backend/trace.txt --mode instruction_mix

```

Based on the given output from the profiler, please fill out the [Excel spreadsheet](../assets/figure9_instruction_mix.xlsx) to generate the figure.

**Same as figure 7, be cautious that figure 9 is characterized based on 16 tasklets**

## Figure 10
<img src="../assets/uPIMulator_figure10.png" width="400"/>

Figure 10 illustrates multi-DPU's latency breakdown and speedup.
You can set the number of fDPUs by changing this parameter: `num_dpus` of `uPIMulator_frontend` and `uPIMulator_backend`.
To plot the breakdown of DPU's runtime, please use [upmem_reg_model](../tools/upmem_reg_model/) under the `tools` directory.
This is our communication model between the host and DPUs that is based on a linear regression method.
You need to prepare an input Excel and feed the input Excel file to `tools/upmem_reg_model/src/main.py`.
We provide a sample input and output Excel file as an example.
You can append a row with the benchmark name, number of DPUs, data prep param that you use for your simulation in the sample input Excel file.
In addition, you need to fill in time in **ms**, for example, kernel execution time.
Then, run `tools/upmem_reg_model/src/main.py`, as shown in commands below.

```
cd /path/to/uPIMulator/tools/upmem_reg_model/src
python main.py --input_excel_filepath /path/to/input_excel_file --output_excel_filepath /path/to/output_excel_file
```

You will be able to find an output Excel file showing the result of linear regression results.

# 🌋 Advanced: Benchmark Addition

A new benchmark can be added under the `uPIMulator_frontend/benchmark` directory.
The benchmark should be written in UPMEM-C language, which is a modified C language, just like CUDA for GPGPU.
Please refer to UPMEM SDK (https://sdk.upmem.com/2021.4.0/) for further instructions on the UPMEM programming model.
On top of the UPMEM programming model, there are some restrictions imposed by the uPIMulator linker:

1. The benchmark should include the same file hierarchy, that is it should include the `dpu` subdirectory and the same `Makefile` used by the PrIM benchmark suite.
This is because the uPIMulator linker automatically detects benchmarks and compiles them using the `Makefile`. 
2. The benchmark can only have `DPU_INPUT_ARGUMENTS` as means to host→DPU communication through WRAM.
Likewise, the benchmark can only have `DPU_RESULTS` as means to DPU→host communication through WRAM.
If one lets the variables have different names, they will be ignored in the linking phase.

If your benchmark complies with the UPMEM programming model and the above restrictions, it is ready to be fed with input data.
Since UPMEM PIM-enabled memory uses physical addresses as-is, so one should "carefully" feed input/output data using physical addresses.
One should provide a Python script that feed input and output data (aka data prep) to the uPIMulator framework for it to generate final binary files.
The data prep script should be placed under the `uPIMulator_frontend/src/assembler/data_prep` and be recognized by `uPIMulator_frontend/src/assembler/assembler.py`.
We already provide 13 PrIM benchmarks ready to be run on the uPIMulator framework.
These are a couple of information you need to be aware of when writing data prep script:

- Data sent from the host to DPUs through `dpu_push_xfer` must be stacked in the `input_dpu_mram_heap_pointer_name` of the data prep script.
- Likewise, data sent from DPUs to the host through `dpu_push_xfer` must be stacked in the `output_dpu_mram_heap_pointer_name` of the data prep script.

We have given example data prep scripts for 13 benchmarks listed below.
One can follow the same style when implementing data prep scripts to add a new benchmark.
