#/bin/bash
directory=$1 # data_sweep_hbm, new_data_sweep
ptw=$2
tlb_way=$3
tlb_set=$4
benchmark=$5 
data_size=$6

for tasklet in 1 2 4 8 16
do
    eval "mkdir -p /home/dongjae/${directory}/tlb_behavior/ptw${ptw}_tlbway${tlb_way}_tlbset${tlb_set}/${benchmark}/${data_size}/"
    echo "mkdir -p /home/dongjae/${directory}/tlb_behavior/ptw${ptw}_tlbway${tlb_way}_tlbset${tlb_set}/${benchmark}/${data_size}/"
    eval "nohup /home/dongjae/upmem_profiler/build/src/upmem_profiler --benchmark ${benchmark} --num_tasklets ${tasklet} --logpath /home/dongjae/${directory}/trace/ptw${ptw}_tlbway${tlb_way}_tlbset${tlb_set}/${benchmark}/${data_size}/${benchmark}.${tasklet}.trace --mode tlb_behavior > /home/dongjae/${directory}/tlb_behavior/ptw${ptw}_tlbway${tlb_way}_tlbset${tlb_set}/${benchmark}/${data_size}/${benchmark}.${tasklet}.csv &"
    echo "nohup /home/dongjae/upmem_profiler/build/src/upmem_profiler --benchmark ${benchmark} --num_tasklets ${tasklet} --logpath /home/dongjae/${directory}/trace/ptw${ptw}_tlbway${tlb_way}_tlbset${tlb_set}/${benchmark}/${data_size}/${benchmark}.${tasklet}.trace --mode tlb_behavior > /home/dongjae/${directory}/tlb_behavior/ptw${ptw}_tlbway${tlb_way}_tlbset${tlb_set}/${benchmark}/${data_size}/${benchmark}.${tasklet}.csv &"
done