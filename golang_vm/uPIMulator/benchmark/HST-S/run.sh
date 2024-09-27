#!/bin/bash

for i in 1 
do
	for b in 64 128 256 512 1024 2048 4096
	do
    	for k in 1 2 4 8 16
	    do
            NR_DPUS=$i NR_TASKLETS=$k BL=10 make all
            wait
            ./bin/host_code -w 2 -e 5 -b ${b} -x 1 > profile/HSTS_${b}_tl${k}_dpu${i}.txt
            wait
            make clean
            wait
		done
	done
done
