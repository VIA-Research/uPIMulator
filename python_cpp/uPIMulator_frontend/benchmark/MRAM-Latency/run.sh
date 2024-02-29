#!/bin/bash

for i in 3 4 5 6 7 8 9 10 11
do
	for j in 1 
	do
		NR_DPUS=1 NR_TASKLETS=$j BL=$i OP=READ make all
		wait
		./bin/host_code -w 0 -e 1 -i 2097152 > profile/read_tl${j}_bl${i}.txt
		wait
		make clean 
		wait

		NR_DPUS=1 NR_TASKLETS=$j BL=$i OP=WRITE make all
		wait
		./bin/host_code -w 0 -e 1 -i 2097152 > profile/write_tl${j}_bl${i}.txt
		wait
		make clean 
		wait
        done
done
