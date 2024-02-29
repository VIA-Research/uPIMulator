#/bin/bash

directory=$1 # data_sweep_hbm, new_data_sweep
benchmark=$2
data_size=$3

for tasklet in 1 2 4 8 16
do
    eval "mkdir -p /home/dongjae/pimulator/new_experiment_result/${directory}/instruction_mix/${benchmark}/${data_size}/"
    echo "mkdir -p /home/dongjae/pimulator/new_experiment_result/${directory}/instruction_mix/${benchmark}/${data_size}/"
    eval "nohup /home/dongjae/pimulator/tool/upmem_profiler/build/src/upmem_profiler --benchmark ${benchmark} --num_tasklets ${tasklet} --logpath /home/dongjae/pimulator/new_experiment_result/${directory}/trace/${benchmark}/${data_size}/${benchmark}.${tasklet}.trace --mode instruction_mix > /home/dongjae/pimulator/new_experiment_result/${directory}/instruction_mix/${benchmark}/${data_size}/${benchmark}.${tasklet}.txt &"
    echo "nohup /home/dongjae/pimulator/tool/upmem_profiler/build/src/upmem_profiler --benchmark ${benchmark} --num_tasklets ${tasklet} --logpath /home/dongjae/pimulator/new_experiment_result/${directory}/trace/${benchmark}/${data_size}/${benchmark}.${tasklet}.trace --mode instruction_mix > /home/dongjae/pimulator/new_experiment_result/${directory}/instruction_mix/${benchmark}/${data_size}/${benchmark}.${tasklet}.txt &"
done