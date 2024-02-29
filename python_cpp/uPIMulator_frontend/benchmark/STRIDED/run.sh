#!/bin/bash

for i in COARSECOARSE FINEFINE
do
	for j in 1 
	do 	
        for k in 1 2 4 8 16
        do 	
            for l in 1 2 4 8 16 32 64 128 256 512 1024 2048 4096
            do 	
                NR_DPUS=$j NR_TASKLETS=$k BL=10 OP=$i make all
                wait
                ./bin/host_code -w 0 -e 1 -i 2097152 -s ${l} >& profile/${i}_${j}_tl${k}_s${l}.txt
                wait
                make clean
                wait
            done
        done
	done
done
