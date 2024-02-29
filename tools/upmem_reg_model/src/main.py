import argparse

from io_.excel_reader import ExcelReader
from io_.excel_writer import ExcelWriter
from regression.model import Model

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--input_excel_filepath",
        type=str,
        default="/Users/bongjoon/upimulator_beta/tools/upmem_reg_model/data/input.xlsx",
    )
    parser.add_argument(
        "--output_excel_filepath",
        type=str,
        default="/Users/bongjoon/upimulator_beta/tools/upmem_reg_model/data/output.xlsx",
    )

    args = parser.parse_args()

    excel_reader = ExcelReader(args.input_excel_filepath)

    min_cpu_to_dpu_bw = 0.0
    min_cpu_to_dpu_mae = 1e10
    for i in range(1, 10000):
        cpu_to_dpu_bw = i * (1.0 / 10000)

        model = Model(excel_reader.benchmarks(), cpu_to_dpu_bw, 0.035)

        cpu_to_dpu_r2, cpu_to_dpu_mae = model.regress_cpu_to_dpu()
        if cpu_to_dpu_mae < min_cpu_to_dpu_mae:
            min_cpu_to_dpu_bw = cpu_to_dpu_bw
            min_cpu_to_dpu_mae = cpu_to_dpu_mae

    min_dpu_to_cpu_bw = 0.0
    min_dpu_to_cpu_mae = 1e10
    for i in range(1, 10000):
        dpu_to_cpu_bw = i * (1.0 / 10000)

        model = Model(excel_reader.benchmarks(), 0.105, dpu_to_cpu_bw)

        dpu_to_cpu_r2, dpu_to_cpu_mae = model.regress_dpu_to_cpu()
        if dpu_to_cpu_mae < min_dpu_to_cpu_mae:
            min_dpu_to_cpu_bw = dpu_to_cpu_bw
            min_dpu_to_cpu_mae = dpu_to_cpu_mae

    min_cpu_to_dpu_bw = 0.2957
    min_dpu_to_cpu_bw = 0.0627
    model = Model(excel_reader.benchmarks(), min_cpu_to_dpu_bw, min_dpu_to_cpu_bw)

    print(f"CPU-DPU BW: {min_cpu_to_dpu_bw}")
    print(f"DPU-CPU BW: {min_dpu_to_cpu_bw}")

    print(f"Kernel: {model.regress_kernel()}")
    print(f"CPU-DPU: {model.regress_cpu_to_dpu()}")
    print(f"DPU-CPU: {model.regress_dpu_to_cpu()}")
    print(f"DPU-DPU: {model.regress_dpu_to_dpu()}")
    print(f"Communication: {model.regress_communication()}")
    print(f"Total: {model.regress_total()}")

    excel_writer = ExcelWriter(args.output_excel_filepath, model.data())
