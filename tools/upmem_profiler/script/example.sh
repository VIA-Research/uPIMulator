benchmark=$1
for size in 128
do
    eval "nohup /home/dongjae/upmem_profiler/build/src/upmem_profiler --num_tasklets 16 --labelpath /home/dongjae/data_sweep/bin/${benchmark}/${size}/${benchmark}.16/labels.bin --logpath /home/dongjae/data_sweep/trace/${benchmark}/${size}/${benchmark}.16.trace --mode instruction_mix > ${benchmark}_${size}.txt &"
done
