#!/bin/bash

for i in 1 
do
	for j in BARRIER HAND
	do 	
		for k in 1 2 4 8 16
		do
		    PERF=1 NR_DPUS=$i NR_TASKLETS=$k BL=10 VERSION=TREE SYNC=$j make all
			wait
            ./bin/host_code -w 2 -e 10 -i 2097152 > profile/TREE_${j}_tl${k}_dpu${i}.txt
            #./bin/host_code -w 2 -e 10 -i 2048 > profile/TREE_${j}_tl${k}_dpu${i}.txt
			wait
			make clean
			wait
		done
	done
done

for i in 1 
do
    for k in 1 2 4 8 16
	do
	    PERF=1 NR_DPUS=$i NR_TASKLETS=$k BL=10 VERSION=SINGLE make all
		wait
        ./bin/host_code -w 2 -e 10 -i 2097152 > profile/SINGLE_SINGLE_tl${k}_dpu${i}.txt
        #./bin/host_code -w 2 -e 10 -i 2048 > profile/SINGLE_SINGLE_tl${k}_dpu${i}.txt
		wait
		make clean
		wait
	done
done
