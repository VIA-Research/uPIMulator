#!/bin/bash

for i in 1 
do
	for j in 1 2 4 8 12 16 
	do 	
            NR_DPUS=$i NR_TASKLETS=$j BL=10 make all
            wait
            ./bin/host_code -w 0 -e 1 -i 2097152 >& profile/gups_${i}_tl${j}.txt
            wait
            make clean
            wait
	done
done
