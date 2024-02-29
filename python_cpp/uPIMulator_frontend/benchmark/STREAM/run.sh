#!/bin/bash

# MRAM
for i in copy copyw add scale triad
do
	for j in 1 
	do 	
        for k in 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 
        do 	
            NR_DPUS=$j NR_TASKLETS=$k BL=10 MEM=MRAM OP=$i make all
            wait
            ./bin/host_code -w 0 -e 1 -i 2097152 >& profile/${i}_${j}_tl${k}_MRAM.txt
            wait
            make clean
            wait
        done
	done
done

# WRAM
for i in copyw add scale triad
do
	for j in 1 
	do 	
        for k in 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 
        do 	
            NR_DPUS=$j NR_TASKLETS=$k BL=10 MEM=WRAM OP=$i make all
            wait
            ./bin/host_code -w 0 -e 1 -i 2097152 >& profile/${i}_${j}_tl${k}_WRAM.txt
            wait
            make clean
            wait
        done
	done
done
