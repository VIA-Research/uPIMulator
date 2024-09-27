# ⚙️ Usage
## Currently Supported Mode
uPIMulator operates in an execution-driven simulation mode, enabling cycle-accurate performance analysis of PIM-based applications.

## Workflow
The typical usage workflow comprises two primary stages:

1. **Binary Generation:** Compile, assemble, and link your application code to generate the required binary files for simulation.
2. **Cycle-Level Simulation:** Utilize the generated binary files as input to the cycle-level simulator to obtain detailed performance metrics and insights.

We are actively working on expanding uPIMulator's capabilities and may introduce additional usage modes in future releases.

## Installation & Build
### Prerequisites
- **Python:** Python 3.10 or higher
- **C++ Compiler:** C++20 compliant compiler
- **CMake:** CMake 3.16 or higher
- **Docker:** Docker Engine
- **User Permissions:** Your Ubuntu user account must be a member of the `docker` group.
- **Operating System:** Ubuntu 18.04 or later (recommended)
- **Environment Management (Optional):** We recommend using a tool like [Anaconda](https://www.anaconda.com/) for managing your Python environment.

### Installation Steps
1. **Install Linker Dependencies:** Navigate to the `uPIMulator_frontend` directory and install the required Python packages using `pip`:

   ```bash
   cd /path/to/uPIMulator/python_cpp/uPIMulator_frontend
   pip install -r requirements.txt
   ```

2. **Build the Cycle-level Simulator:** Navigate to the `uPIMulator_backend/script` directory and execute the build script:

   ```bash
   cd /path/to/uPIMulator/python_cpp/uPIMulator_backend/script
   sh build.sh
   ```

> **Note:** Replace `/path/to/uPIMulator` with the actual path to your uPIMulator repository. 

If you encounter any issues during the installation or build process, please refer to the troubleshooting section in the documentation or open an issue on our GitHub repository.

Certainly, here's the paragraph revised for a more professional and appropriate tone suitable for a GitHub README:

## Binary Files Generation
We will use the VA (vector addition) benchmark as an example to demonstrate the binary file generation phase.
Upon successful completion of the compile/assemble/link process, you will find the generated binary files within the `uPIMulator_frontend/bin` directory.

> **Pre-Generated Binaries:** To expedite the setup process, you can utilize our pre-generated binary files available at the following [link](https://drive.google.com/file/d/1kfL-xGn1F18Ezmw81IvAhxEaLiZZOLFR/view?usp=sharing).

### Compilation

```bash
cd /path/to/uPIMulator/python_cpp/uPIMulator_frontend/src
python main.py --mode compile --num_tasklets 16 --benchmark VA --num_dpus 1
```

### Assembly and Linking

```bash
cd /path/to/uPIMulator/python_cpp/uPIMulator_frontend/src
python main.py --mode link --num_tasklets 16 --benchmark VA --data_prep_param 1024 --num_dpus 1
```

> **Troubleshooting:** In the event that the linking command encounters errors, please execute the following Docker build command and then reattempt the compilation and linking steps:

```bash
docker build -t bongjoonhyun/compiler -f /path/to/uPIMulator/python_cpp/uPIMulator_frontend/docker/compiler.dockerfile .
```

> **Supported Benchmarks and Performance Note:** uPIMulator currently supports 13 PrIM benchmarks.
> Please be aware that the initial compile/assemble/link process may take approximately 30 minutes.

Certainly, let's revise the paragraph for a more professional and appropriate tone suitable for a GitHub README:

## Cycle-level Simulation
### Executing a Simulation
Initiate a cycle-level simulation by providing the following inputs to the `uPIMulator` executable:

- Benchmark name
- Number of tasklets
- Absolute path to the directory containing the generated binary files (`bindir`)

You can further customize the simulation behavior using various command-line options. 

### Simulation Output
The simulation results will be printed directly to the standard output (`stdout`).

### Example Command

```bash
cd /path/to/uPIMulator/uPIMulator_backend/build/
./src/uPIMulator --benchmark VA --num_tasklets 16 --bindir /path/to/uPIMulator/uPIMulator_frontend/bin/1_dpus/ --logdir .
```

> **Important Note:** Ensure that you provide the absolute path to the `bindir` when executing the simulation.

Feel free to explore different benchmark configurations and utilize the command-line options to tailor the simulation to your specific requirements.

Certainly, let's revise the paragraph for a more professional and appropriate tone suitable for a GitHub README:

# 📄 Reproducing Figures from the Paper
This section provides instructions for replicating the figures presented in our paper.
We offer detailed replication manuals for figures used in **Section 4 (Demystifying UPMEM-PIM with uPIMulator)**.
Please note that Figure 8 is omitted for brevity.

## Configuration of PrIM Benchmarks
- **Single DPU Focus:** For Figures 5, 6, 7, and 9, the `num_dpus` parameter must always be set to `1`, as these experiments specifically characterize the behavior of a single DPU.
- **Data Preparation Parameter:**  When generating the binary files for the PrIM benchmarks, please configure the `data_prep_param` parameter according to the following table:

| Benchmark | `data_prep_param` (Figures 5, 6, 7, 9) | `data_prep_param` (Figure 10) |
|---|---|---|
| BS       | 32768 | 131072 |
| GEMV     | 2048  | 4096   |
| HST-L    | 131072 | 524288 |
| HST-S    | 131072 | 524288 |
| MLP      | 256   | 1024   |
| RED      | 524288 | 2097152|
| SCAN-RSS | 262144 | 1048576|
| SCAN-SSA | 262144 | 1048576|
| SEL      | 524288 | 2097152|
| TRNS     | 1024  | 128    |
| TS       | 2048  | 65536  |
| UNI      | 524288 | 2097152|
| VA       | 524288 | 2097152|

## Example Command

```bash
python main.py --mode link --num_tasklets 16 --benchmark VA --data_prep_param 524288 --num_dpus 1
``` 

Please ensure you adhere to these configurations to accurately replicate the figures presented in the paper. 

## Figure 5: PrIM Compute and Memory Utilization
<img src="../assets/uPIMulator_figure5.png" width="400"/>

This figure illustrates the compute utilization (red points) and memory read bandwidth utilization (blue points) of the PrIM benchmarks when executed with 1, 4, and 16 threads (tasklets).

### Calculation Formulas

- **Compute Utilization (IPC):**  `num_instructions` / (`logic_cycle` - `communication_cycle`)
- **Memory Read Bandwidth Utilization (GB/s):** Please refer to the provided Excel sheet for the calculation at the following [link](../assets/figure5_mem_util_calculator.xlsx)

> **Note:** The required values for these calculations can be obtained from the simulation results generated by uPIMulator. 

## Figure 6: DPU Runtime Breakdown
<img src="../assets/uPIMulator_fiture6.png" width="400"/>

This figure presents a breakdown of DPU runtime, categorizing cycles into active (black) and various idle states (red, yellow, blue).
The following formulas can be used to calculate the proportions of each category:

### Calculation Formulas

- **Issuable Ratio:**  `breakdown_run` / (`logic_cycle` - `communication_cycle`)
- **Idle (Memory) Ratio:** `breakdown_dma` / (`logic_cycle` - `communication_cycle`)
- **Idle (Revolver) Ratio:** (`breakdown_etc` - `communication_cycle`) / (`logic_cycle` - `communication_cycle`)
- **Idle (RF) Ratio:** `backpressuer` / (`logic_cycle` - `communication_cycle`)

> **Note:** The values for the variables in these formulas (e.g., `breakdown_run`, `logic_cycle`) can be extracted from the simulation results generated by uPIMulator. 

## Figure 7: Issuable Tasklets
<img src="../assets/uPIMulator_figure7.png" width="400"/>

Figure 7 visualizes the number of tasklets (threads) that are ready for execution (issuable) by the DPU scheduler at each cycle.

### Replication
To reproduce this figure, utilize the provided [Excel sheet](../assets/figure7_active_tasklet_breakdown.xlsx).
The spreadsheet includes instructions on how to populate it with the relevant simulation output data, and it will automatically generate the corresponding figure.

> **Important Configuration Note:** Please ensure that the number of threads is configured to **16 tasklets** when running the simulations for this figure.
> You can achieve this by using the following command-line argument: `--num_tasklets 16`.

## Figure 9: Instruction Mix (Single DPU)
<img src="../assets/uPIMulator_figure9.png" width="400"/>

Figure 9 provides a breakdown of the instruction mix observed during single-DPU execution.
To generate this figure, follow the steps outlined below using the `upmem_profiler` tool and the accompanying Excel sheet.

### Procedure

1. **Build the Profiler**

   ```bash
   cd /path/to/uPIMulator/tools/upmem_profiler/script
   bash build.sh
   ```

2. **Extract Instructions**
   Run the simulation with the `--verbose 1` flag to capture detailed instruction traces.

   ```bash
   cd /path/to/uPIMulator/python_cpp/uPIMulator_backend/
   ./build/src/uPIMulator --benchmark VA --num_tasklets 16 --bindir /path/to/uPIMulator/python_cpp/uPIMulator_frontend/bin/1_dpus/ --logdir . --verbose 1 > trace.txt
   ```

3. **Run the Profiler**
   Process the generated trace file using the `upmem_profiler` in `instruction_mix` mode.

   ```bash
   cd /path/to/uPIMulator/tools/upmem_profiler/
   ./build/src/upmem_profiler --logpath /path/to/uPIMulator/python_cpp/uPIMulator_backend/trace.txt --mode instruction_mix
   ```

4. **Generate the Figure**
   Utilize the profiler's output to populate the provided [Excel sheet](../assets/figure9_instruction_mix.xlsx), which will automatically generate the instruction mix figure.

> **Important Configuration Note:** Similar to Figure 7, the instruction mix analysis in Figure 9 is based on simulations with **16 tasklets**.
> Ensure that you maintain this configuration (`--num_tasklets 16`) for accurate replication. 

## Figure 10: Multi-DPU Latency Breakdown and Speedup
<img src="../assets/uPIMulator_figure10.png" width="400"/>

Figure 10 presents the latency breakdown and speedup achieved in multi-DPU scenarios.

### Configuring the Number of DPUs
You can adjust the number of DPUs by modifying the `num_dpus` parameter in both the `uPIMulator_frontend` and `uPIMulator_backend`.

### Generating the Latency Breakdown
To obtain the latency breakdown data for plotting, utilize the `upmem_reg_model` tool located in the `tools/upmem_reg_model/` directory.
This tool implements a communication model between the host and DPUs based on linear regression.

### Procedure

1. **Prepare Input Excel:**
   - We provide a sample input Excel file as a template.
   - Append a new row to this file, specifying the benchmark name, number of DPUs, and the `data_prep_param` used in your simulation.
   - Fill in the relevant time values (in milliseconds) obtained from your simulation results, such as kernel execution time.

2. **Run the Regression Model:**

   ```bash
   cd /path/to/uPIMulator/tools/upmem_reg_model/src
   python main.py --input_excel_filepath /path/to/your/input_excel_file --output_excel_filepath /path/to/your/output_excel_file
   ```

3. **Access the Output:** 
   - The linear regression results will be available in the specified output Excel file.
   - Use this data to create the latency breakdown plots as shown in Figure 10.

Please ensure that you follow these steps carefully to accurately reproduce the multi-DPU latency breakdown and speedup analysis presented in the paper.

# 🌋 Adding Custom Benchmarks
uPIMulator empowers you to go beyond the provided PrIM benchmark suite by incorporating your own custom benchmarks.
This is particularly beneficial if you have access to UPMEM-PIM hardware and want to evaluate your code's performance in a simulated environment.

## Requirements
To successfully integrate a new benchmark, ensure it adheres to the following:

1. **UPMEM-C Language:**  The benchmark must be implemented in UPMEM-C, a C-like language tailored for UPMEM-PIM programming.
Consult the [UPMEM SDK documentation](https://sdk.upmem.com/2021.4.0/) for detailed programming guidelines.

2. **File Structure and Naming:**  
   - Maintain the same file hierarchy as the PrIM benchmarks, including a `dpu` subdirectory.
   - Utilize the same `Makefile` structure used by the PrIM suite for automated compilation.

3. **Communication Variables:**
   - **Host-to-DPU:** Use the `DPU_INPUT_ARGUMENTS` variable for data transfer from the host to DPUs via WRAM.
   - **DPU-to-Host:**  Use the `DPU_RESULTS` variable for data transfer from DPUs to the host via WRAM.
   > **Important:** Deviating from these variable names will result in the linker ignoring the communication during the linking phase.

## Data Preparation
Once your benchmark meets the above criteria, it's ready for data input.
Since UPMEM PIM-enabled memory directly utilizes physical addresses, exercise caution when feeding input/output data. 

You'll need to provide a Python script to handle data preparation for your benchmark.
This script should reside in the `uPIMulator_frontend/src/assembler/data_prep` directory and be recognized by `uPIMulator_frontend/src/assembler/assembler.py`.

> **Key Considerations for Data Preparation Scripts: 
> - Data transferred from the host to DPUs using `dpu_push_xfer` must be organized within the `input_dpu_mram_heap_pointer_name` variable in your data preparation script.
> - Similarly, data transferred from DPUs to the host using `dpu_push_xfer` should be placed within the `output_dpu_mram_heap_pointer_name` variable.

## Reference Examples
We have included data preparation scripts for the 13 supported PrIM benchmarks.
These serve as excellent references for structuring your custom data preparation scripts.

By following these guidelines, you can seamlessly integrate your benchmarks into uPIMulator for comprehensive performance evaluation and analysis. 

If you have any questions or encounter any difficulties during the integration process, don't hesitate to reach out to us for support. 
