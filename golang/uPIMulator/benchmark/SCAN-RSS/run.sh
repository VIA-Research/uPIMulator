#!/bin/bash

for i in  2048 4096 8192 16384 65536 262144 1048576 3932160
do
		    NR_DPUS=1 NR_TASKLETS=16 BL=10 VERSION=SINGLE make all
			wait
            ./bin/host_code -w 10 -e 100 -i ${i} > profile/out${i}_tl16_bl10_dpu11
			wait
			make clean
			wait
done
