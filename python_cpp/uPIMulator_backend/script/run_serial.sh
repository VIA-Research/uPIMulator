#/bin/bash

num_dpus=$1
benchmark=$2
data_size=$3

echo "${benchmark} - ${data_size}"
eval "mkdir -p /home/dongjae/pimulator/hpca_2024_log/multi_dpus_validation/${benchmark}/${data_size}/${num_dpus}_dpus/${benchmark}.16."
sim_name="/home/dongjae/pimulator/simulator/upmem_sim_multi_dpus/build/src/upmem_sim"
bin_dir="/home/dongjae/pimulator/bin_files_validation/bin_multi/${benchmark}/${data_size}/${num_dpus}_dpus/"
log_dir="/home/dongjae/pimulator/hpca_2024_log/multi_dpus_validation/${benchmark}/${data_size}/${num_dpus}_dpus/${benchmark}.16." # Note that you should add '/' at the last position
cmd="nohup ${sim_name} --benchmark ${benchmark} --num_tasklets 16 --bindir ${bin_dir} --logdir ${log_dir} --num_dpus ${num_dpus}"
echo ${cmd}
eval ${cmd}

