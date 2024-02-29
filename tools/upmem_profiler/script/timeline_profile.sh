data_size=(32768 2048 128 128 256 524288 65536 65536 524288 1024 2048 524288 524288)
idx=0

for benchmark in "BS" "GEMV" "HST-L" "HST-S" "MLP" "RED" "SCAN-RSS" "SCAN-SSA" "SEL" "TRNS" "TS" "UNI" "VA"
do
    for tasklet in 1 2 4 8 16
    do
        eval "mkdir -p /home/dongjae/data_sweep/timeline/${benchmark}/${data_size[${idx}]}/"
        eval "nohup /home/dongjae/upmem_profiler/build/src/upmem_profiler --benchmark ${benchmark} --num_tasklets ${tasklet} --labelpath /home/dongjae/data_sweep/bin/${benchmark}/${data_size[${idx}]}/${benchmark}.${tasklet}/labels.bin --logpath /home/dongjae/data_sweep/trace/${benchmark}/${data_size[${idx}]}/${benchmark}.${tasklet}.trace --mode timeline > /home/dongjae/data_sweep/timeline/${benchmark}/${data_size[${idx}]}/${benchmark}.${tasklet}.txt &"
    done
    let "idx=idx+1"
done


echo "All profilers are running the server!"