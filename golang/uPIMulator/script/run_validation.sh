#/bin/bash

root_dir=/home/via/dongjae/upimulator_beta/golang/uPIMulator

eval "mkdir -p ${root_dir}/bin/"

for num_tasklets in 1 2 4 8 16
do
    eval "${root_dir}/build/uPIMulator --root_dirpath ${root_dir}/ --bin_dirpath ${root_dir}/bin --benchmark BS --num_channels 1 --num_ranks_per_channel 1 --num_dpus_per_rank 1 --num_tasklets ${num_tasklets} --data_prep_params 8192"
    eval "mv ${root_dir}/bin/log.txt ${root_dir}/validation_log/BS_${num_tasklets}.txt"
    eval "rm ${root_dir}/bin/*"
    echo "BS-${num_tasklets} done \n"

    eval "${root_dir}/build/uPIMulator --root_dirpath ${root_dir}/ --bin_dirpath ${root_dir}/bin --benchmark GEMV --num_channels 1 --num_ranks_per_channel 1 --num_dpus_per_rank 1 --num_tasklets ${num_tasklets} --data_prep_params 256"
    eval "mv ${root_dir}/bin/log.txt ${root_dir}/validation_log/GEMV_${num_tasklets}.txt"
    eval "rm ${root_dir}/bin/*"
    echo "GEMV-${num_tasklets} done \n"

    eval "${root_dir}/build/uPIMulator --root_dirpath ${root_dir}/ --bin_dirpath ${root_dir}/bin --benchmark HST-L --num_channels 1 --num_ranks_per_channel 1 --num_dpus_per_rank 1 --num_tasklets ${num_tasklets} --data_prep_params 32768"
    eval "mv ${root_dir}/bin/log.txt ${root_dir}/validation_log/HST-L_${num_tasklets}.txt"
    eval "rm ${root_dir}/bin/*"
    echo "HST-L-${num_tasklets} done \n"

    eval "${root_dir}/build/uPIMulator --root_dirpath ${root_dir}/ --bin_dirpath ${root_dir}/bin --benchmark HST-S --num_channels 1 --num_ranks_per_channel 1 --num_dpus_per_rank 1 --num_tasklets ${num_tasklets} --data_prep_params 32768"
    eval "mv ${root_dir}/bin/log.txt ${root_dir}/validation_log/HST-S_${num_tasklets}.txt"
    eval "rm ${root_dir}/bin/*"
    echo "HST-S-${num_tasklets} done \n"

    eval "${root_dir}/build/uPIMulator --root_dirpath ${root_dir}/ --bin_dirpath ${root_dir}/bin --benchmark MLP --num_channels 1 --num_ranks_per_channel 1 --num_dpus_per_rank 1 --num_tasklets ${num_tasklets} --data_prep_params 64"
    eval "mv ${root_dir}/bin/log.txt ${root_dir}/validation_log/MLP_${num_tasklets}.txt"
    eval "rm ${root_dir}/bin/*"
    echo "MLP-${num_tasklets} done \n"

    eval "${root_dir}/build/uPIMulator --root_dirpath ${root_dir}/ --bin_dirpath ${root_dir}/bin --benchmark RED --num_channels 1 --num_ranks_per_channel 1 --num_dpus_per_rank 1 --num_tasklets ${num_tasklets} --data_prep_params 65536"
    eval "mv ${root_dir}/bin/log.txt ${root_dir}/validation_log/RED_${num_tasklets}.txt"
    eval "rm ${root_dir}/bin/*"
    echo "RED-${num_tasklets} done \n"

    eval "${root_dir}/build/uPIMulator --root_dirpath ${root_dir}/ --bin_dirpath ${root_dir}/bin --benchmark SCAN-RSS --num_channels 1 --num_ranks_per_channel 1 --num_dpus_per_rank 1 --num_tasklets ${num_tasklets} --data_prep_params 65536"
    eval "mv ${root_dir}/bin/log.txt ${root_dir}/validation_log/SCAN-RSS_${num_tasklets}.txt"
    eval "rm ${root_dir}/bin/*"
    echo "SCAN-RSS-${num_tasklets} done \n"

    eval "${root_dir}/build/uPIMulator --root_dirpath ${root_dir}/ --bin_dirpath ${root_dir}/bin --benchmark SCAN-SSA --num_channels 1 --num_ranks_per_channel 1 --num_dpus_per_rank 1 --num_tasklets ${num_tasklets} --data_prep_params 65536"
    eval "mv ${root_dir}/bin/log.txt ${root_dir}/validation_log/SCAN-SSA_${num_tasklets}.txt"
    eval "rm ${root_dir}/bin/*"
    echo "SCAN-SSA-${num_tasklets} done \n"

    eval "${root_dir}/build/uPIMulator --root_dirpath ${root_dir}/ --bin_dirpath ${root_dir}/bin --benchmark SEL --num_channels 1 --num_ranks_per_channel 1 --num_dpus_per_rank 1 --num_tasklets ${num_tasklets} --data_prep_params 65536"
    eval "mv ${root_dir}/bin/log.txt ${root_dir}/validation_log/SEL_${num_tasklets}.txt"
    eval "rm ${root_dir}/bin/*"
    echo "SEL-${num_tasklets} done \n"

    eval "${root_dir}/build/uPIMulator --root_dirpath ${root_dir}/ --bin_dirpath ${root_dir}/bin --benchmark TRNS --num_channels 1 --num_ranks_per_channel 1 --num_dpus_per_rank 1 --num_tasklets ${num_tasklets} --data_prep_params 128"
    eval "mv ${root_dir}/bin/log.txt ${root_dir}/validation_log/TRNS_${num_tasklets}.txt"
    eval "rm ${root_dir}/bin/*"
    echo "TRNS-${num_tasklets} done \n"

    eval "${root_dir}/build/uPIMulator --root_dirpath ${root_dir}/ --bin_dirpath ${root_dir}/bin --benchmark TS --num_channels 1 --num_ranks_per_channel 1 --num_dpus_per_rank 1 --num_tasklets ${num_tasklets} --data_prep_params 256"
    eval "mv ${root_dir}/bin/log.txt ${root_dir}/validation_log/TS_${num_tasklets}.txt"
    eval "rm ${root_dir}/bin/*"
    echo "TS-${num_tasklets} done \n"

    eval "${root_dir}/build/uPIMulator --root_dirpath ${root_dir}/ --bin_dirpath ${root_dir}/bin --benchmark UNI --num_channels 1 --num_ranks_per_channel 1 --num_dpus_per_rank 1 --num_tasklets ${num_tasklets} --data_prep_params 65536"
    eval "mv ${root_dir}/bin/log.txt ${root_dir}/validation_log/UNI_${num_tasklets}.txt"
    eval "rm ${root_dir}/bin/*"
    echo "UNI-${num_tasklets} done \n"

    eval "${root_dir}/build/uPIMulator --root_dirpath ${root_dir}/ --bin_dirpath ${root_dir}/bin --benchmark VA --num_channels 1 --num_ranks_per_channel 1 --num_dpus_per_rank 1 --num_tasklets ${num_tasklets} --data_prep_params 65536"
    eval "mv ${root_dir}/bin/log.txt ${root_dir}/validation_log/VA_${num_tasklets}.txt"
    eval "rm ${root_dir}/bin/*"
    echo "VA-${num_tasklets} done \n"


done