#!/bin/bash

# WRAM
for i in streaming random
do
	for j in 1 
	do 	
        for k in 1 2 4 8 16  
        do 	
            NR_DPUS=$j NR_TASKLETS=$k BL=10 MEM=WRAM OP=$i make all
            wait
            ./bin/host_code -w 0 -e 1 -i 2097152 >& profile/${i}_${j}_tl${k}_s1_WRAM.txt
            wait
            make clean
            wait
        done
	done
done

for i in strided 
do
	for j in 1 
	do 	
        for k in 1 2 4 8 16  
        do 	
            for l in 1 2 4 8 16 32 64
            do 	
                NR_DPUS=$j NR_TASKLETS=$k BL=10 MEM=WRAM OP=$i make all
                wait
                ./bin/host_code -w 0 -e 1 -i 2097152 -s ${l} >& profile/${i}_${j}_tl${k}_s${l}_WRAM.txt
                wait
                make clean
                wait
	        done
        done
	done
done
