#data_size=(32768 2048 128 128 256 524288 65536 65536 524288 1024 2048 524288 524288) # 100ms
#data_size=(8192 512 32 32 64 131072 16384 16384 131072 256 512 131072 131072)
# data_size=(1024)
# idx=0

# #for benchmark in "BS" "GEMV" "HST-L" "HST-S" "MLP" "RED" "SCAN-RSS" "SCAN-SSA" "SEL" "TRNS" "TS" "UNI" "VA"
# for benchmark in "TRNS"
# do
#     for tasklet in 16
#     do
#         eval "mkdir -p /home/dongjae/new_data_sweep/active_tasklet/${benchmark}/${data_size[${idx}]}/"
#         eval "nohup /home/dongjae/upmem_profiler/build/src/upmem_profiler --benchmark ${benchmark} --num_tasklets ${tasklet} --labelpath /home/dongjae/data_sweep/bin/${benchmark}/${data_size[${idx}]}/${benchmark}.${tasklet}/labels.bin --logpath /home/dongjae/new_data_sweep/trace/${benchmark}/${data_size[${idx}]}/${benchmark}.${tasklet}.trace --mode active_tasklet > /home/dongjae/new_data_sweep/active_tasklet/${benchmark}/${data_size[${idx}]}/${benchmark}.${tasklet}.txt &"
#     done
#     let "idx=idx+1"
# done


case=$1 # data_sweep_hbm, new_data_sweep
benchmark=$2 
data_size=$3

for tasklet in 1 2 4 8 16
do
    echo "mkdir -p /home/dongjae/${case}/active_tasklet/${benchmark}/${data_size}/"
    echo "nohup /home/dongjae/upmem_profiler/build/src/upmem_profiler --benchmark ${benchmark} --num_tasklets ${tasklet} --labelpath /home/dongjae/data_sweep/bin/${benchmark}/${data_size}/${benchmark}.${tasklet}/labels.bin --logpath /home/dongjae/${case}/trace/${benchmark}/${data_size}/${benchmark}.${tasklet}.trace --mode active_tasklet > /home/dongjae/${case}/active_tasklet/${benchmark}/${data_size}/${benchmark}.${tasklet}.txt &"
    eval "mkdir -p /home/dongjae/${case}/active_tasklet/${benchmark}/${data_size}/"
    eval "nohup /home/dongjae/upmem_profiler/build/src/upmem_profiler --benchmark ${benchmark} --num_tasklets ${tasklet} --labelpath /home/dongjae/data_sweep/bin/${benchmark}/${data_size}/${benchmark}.${tasklet}/labels.bin --logpath /home/dongjae/${case}/trace/${benchmark}/${data_size}/${benchmark}.${tasklet}.trace --mode active_tasklet > /home/dongjae/${case}/active_tasklet/${benchmark}/${data_size}/${benchmark}.${tasklet}.txt &"
done

echo "All profilers are running the server!"