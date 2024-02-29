#/bin/bash
directory=$1 # data_sweep_hbm, new_data_sweep
case=$2
benchmark=$3 
data_size=$4

for tasklet in 1 2 4 8 16
do
    eval "mkdir -p /home/dongjae/${directory}/mram_access_pattern/${benchmark}/${data_size}/"
    echo "mkdir -p /home/dongjae/${directory}/mram_access_pattern/${benchmark}/${data_size}/"
    eval "nohup /home/dongjae/upmem_profiler/build/src/upmem_profiler --benchmark ${benchmark} --num_tasklets ${tasklet} --logpath /home/dongjae/${directory}/trace/${benchmark}/${data_size}/${benchmark}.${tasklet}.trace --mode mram_access_pattern --case_study ${case} > /home/dongjae/${directory}/mram_access_pattern/${benchmark}/${data_size}/${benchmark}.${tasklet}.csv &"
    echo "nohup /home/dongjae/upmem_profiler/build/src/upmem_profiler --benchmark ${benchmark} --num_tasklets ${tasklet} --logpath /home/dongjae/${directory}/trace/${benchmark}/${data_size}/${benchmark}.${tasklet}.trace --mode mram_access_pattern --case_study ${case} > /home/dongjae/${directory}/mram_access_pattern/${benchmark}/${data_size}/${benchmark}.${tasklet}.csv &"
done