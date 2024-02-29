# ⚙️ Usage

uPIMulator currently supports one usage mode, which is an execution-driven mode.
The flow of using the uPIMulator framework is one-folded: to compile, assemble, and link to generate binary files, and feed the binary files to the cycle-level simulator to get simulation results.

# 🏃 Getting Started

## Install & Build

### Prerequisites:
- uPIMulator requires Go compiler and SDK whose version is higher than 1.21.5.
You can download and install Go at the [official website](https://go.dev/doc/install)
- uPIMulator requires Docker.
- Your Ubuntu account must be added to the Docker group.
- uPIMulator has been tested on Ubuntu 18.04 using an Intel CPU.

### Install:
```
cd /path/to/uPIMulator/golang/uPIMulator/script
python build.py
```

## Binary Files Generation & Cycle-level Simulation

We will use the VA (vector addition) benchmark as an example for the binary file generation phase.
For the first simulation, it would take roughly 30 minutes.
One can run a simulation by giving the benchmark name, the number of tasklets, and the path to a directory where you want to put binary files, log files, etc.
One can tweak simulation parameters through command line options.
The simulation results are printed out to `stdout`.

- **Note that you need to create an empty `bin` directory before the simulation.**
- **Note that you need to specifiy absolute paths to `root_dirpath` and `bin_dirpath`.**

```
cd /path/to/uPIMulator/golang/uPIMulator
rm -rf bin
mkdir bin
./build/uPIMulator --root_dirpath /path/to/uPIMulator/golang/uPIMulator --bin_dirpath /path/to/uPIMulator/golang/uPIMulator/bin --benchmark VA --num_channels 1 --num_ranks_per_channel 1 --num_dpus_per_rank 1 --num_tasklets 16 --data_prep_params 1024
```

We provide pre-generated binary files at the [following link](https://drive.google.com/file/d/1gLt0XgFbkb4EAY-8Usw_GP1Ef6NTK08J/view).

# 📄 Reproducing Figures from the Paper
To replicate the figures illustrated in our paper, please follow the instructions below.
Note that we only provides instructions for Figure 5 and Figure 6 for brevity.

## Configuration of PrIM Benchmarks
- Parameter `num_channels`, `num_ranks_per_channel`, and `num_dpus_per_rank` must always be configured to 1, since Figure 5 and Figure 6 experiment only characterizes the behavior of 1 DPU.
- When generating the binary file of the PrIM benchmark, please configure the `data_prep_param` parameter as shown in the following table. 
(e.g., `./uPIMulator --root_dirpath /path/to/uPIMulator/ --bin_dirpath /path/to/uPIMulator/bin --benchmark VA --num_channels 1 --num_ranks_per_channel 1 --num_dpus_per_rank 1 --num_tasklets 16 --data_prep_params 524288`)

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

- **Compute utilization (i.e., IPC)** = (value of `num_instructions`) / (value of `logic_cycle`) 
- **Memory read bandwidth utilization (GB/s)** = (We provide an excel sheet to calculate the memory read bandwidth utilization in the following [link](../assets/figure5_mem_util_calculator.xlsx))

Note that the value of `XYZ` can be obtained by the simulation results.

## Figure 6
<img src="../assets/uPIMulator_fiture6.png" width="400"/>

Figure 6 illustrates the breakdown of DPU’s runtime into active (black) and idle (red, yellow, blue) cycles.
To plot the breakdown of DPU's runtime, please use the formula as follows:

- **Issuable ratio** = (value of `breakdown_run`) / (value of `logic_cycle`)
- **Idle (Memory) ratio** = (value of `breakdown_dma`) / (value of `logic_cycle`)
- **Idle (Revolver) ratio** = (value of `breakdown_etc`) / (value of `logic_cycle`)
- **Idle (RF) ratio** = (value of `backpressure`) / (value of `logic_cycle`)

Note that the value of `XYZ` can be obtained by the simulation results.

## Figure 7
<img src="../assets/uPIMulator_figure7.png" width="400"/>

Figure 7 shows the number of issuable tasklets (i.e., threads) by the DPU scheduler. To replicate the above figure, please use the following [excel sheet](../assets/figure7_active_tasklet_breakdown.xlsx) provided in our repository. By filling out the excel spreadsheet based on the simulation output (the spreadsheet explains what kind of simulation output needs to be filled out), it automatically generates the figure.

**Note that the number of threads configured in this experiment is ONLY 16 tasklets.**
(e.g., `--num_tasklets 16`)

## Figure 9
<img src="../assets/uPIMulator_figure9.png" width="400"/>

Figure 9 illustrates the instruction mix when executing with a single DPU. To plot the given figure, please both use [upmem_profiler](../tools/upmem_profiler/) and [excel sheet](../assets/figure9_instruction_mix.xlsx) and follow the procedures as explained below.

### Build Profiler

```
cd /path/to/uPIMulator/tools/upmem_profiler/script
bash build.sh
```

### Extract Instructions
When extracting instructions, you should add additional parameter, `--verbose 1`. 
```
cd /path/to/uPIMulator/golang/uPIMulator/
./build/uPIMulator --root_dirpath /path/to/uPIMulator/golang/uPIMulator --bin_dirpath /path/to/uPIMulator/golang/uPIMulator/bin --benchmark VA --num_channels 1 --num_ranks_per_channel 1 --num_dpus_per_rank 1 --num_tasklets 16 --data_prep_params 1024 > trace.txt
```

### Run Profiler
```
cd /path/to/uPIMulator/tools/upmem_profiler/
./build/src/upmem_profiler --logpath /path/to/uPIMulator/golang/uPIMulator/upmem_sim_multi_dpus/trace.txt --mode instruction_mix

```

Based on the given output from the profiler, please fill out the [excel sheet](../assets/figure9_instruction_mix.xlsx) to generate the figure.

**Same as figure 7, be cautious that figure 9 is characterized based on 16 tasklets**

## Figure 10
<img src="../assets/uPIMulator_figure10.png" width="400"/>

Figure 10 illustrates multi-DPU's latency breakdown and speedup.
You can set the number of DPUs by changing these parameters: `num_channels`, `num_ranks_per_channel`, `num_dpus_per_rank` of `uPIMulator`.
To plot the breakdown of DPU's runtime, please use [upmem_reg_model](../tools/upmem_reg_model/) under the `tools` directory.
This is our communication model between the host and DPUs that is based on a linear regression method.
You need to prepare an input Excel and feed the input Excel file to `tools/upmem_reg_model/src/main.py`.
We provide a sample input and output Excel file as an example.
You can append a row with the benchmark name, number of DPUs, data prep param that you use for your simulation in the sample input Excel file.
In addition, you need to fill in time in **ms**, for example, kernel execution time. You can convert 
Then, run `tools/upmem_reg_model/src/main.py`, as shown in commands below.

```
cd /path/to/uPIMulator/tools/upmem_reg_model/src
python main.py --input_excel_filepath /path/to/input_excel_file --output_excel_filepath /path/to/output_excel_file
```

You will be able to find an output Excel file showing the result of linear regression results.

# 🌋 Benchmark Addition

A new benchmark can be added under the `benchmark` directory.
The benchmark should be written in UPMEM-C language, which is a modified C language, just like CUDA for GPGPU.
Please refer to UPMEM SDK (https://sdk.upmem.com/2021.4.0/) for further instructions on the UPMEM programming model.
On top of the UPMEM programming model, there are some restrictions imposed by the uPIMulator linker:

1. The benchmark should include the same file hierarchy, that is it should include the `dpu` subdirectory and the same `CMakeLists.txt` used by the PrIM benchmark suite.
This is because the uPIMulator linker automatically detects benchmarks and compiles them using the `CMakeLists.txt`. 

Since UPMEM PIM-enabled memory uses physical address as-is and the uPIMulator framework does not support concurrent execution of host and PIM-enabled memory, one should "carefully" feed input/output data using physical address. One should provide input/output feed Go source code (aka data prep) to the uPIMulator framework for it to generate final binary files.
The data prep script should be placed under the `uPImulator/src/assembler` and be recognized by `uPIMulator/src/assembler/assembler.go`.
We already provide 13 PrIM benchmarks ready to be run on the uPIMulator framework.
These are a couple of information you need to be aware of when writing data prep script:

- Data sent from the host to DPUs through `dpu_push_xfer` must be stacked in the `input_dpu_mram_heap_pointer_name` of the data prep script.
- Likewise, data sent from DPUs to the host through `dpu_push_xfer` must be stacked in the `output_dpu_mram_heap_pointer_name` of the data prep.

We have given example data prep scripts for 13 benchmarks listed below.
One can follow the same style when implementing data prep scripts to add a new benchmark.
