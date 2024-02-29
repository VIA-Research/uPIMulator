#/bin/bash

num_dpus=$1
benchmark=$2
data_size=$3

eval "mkdir -p /home/dongjae/pimulator/new_experiment_result/multi_dpus/trace/${benchmark}/${data_size}/${num_dpus}_dpus"

echo "${benchmark} - ${data_size}"
eval "mkdir -p /home/dongjae/pimulator/new_experiment_result/multi_dpus/log/${benchmark}/${data_size}/${num_dpus}_dpus/${benchmark}.16."
sim_name="/home/dongjae/pimulator/simulator/upmem_sim_multi_dpus/build/src/upmem_sim"
bin_dir="/home/dongjae/pimulator/bin_files/bin_multi/${benchmark}/${data_size}/${num_dpus}_dpus/"
log_dir="/home/dongjae/pimulator/new_experiment_result/multi_dpus/log/${benchmark}/${data_size}/${num_dpus}_dpus/${benchmark}.16." # Note that you should add '/' at the last position
log_path="/home/dongjae/pimulator/new_experiment_result/multi_dpus/trace/${benchmark}/${data_size}/${num_dpus}_dpus/"
cmd="nohup ${sim_name} --benchmark ${benchmark} --num_tasklets 16 --bindir ${bin_dir} --logdir ${log_dir} --num_dpus ${num_dpus} 1> ${log_path}/${benchmark}.16.trace &"
echo ${cmd}
eval ${cmd}

