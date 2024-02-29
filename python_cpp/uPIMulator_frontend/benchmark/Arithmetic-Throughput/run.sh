#!/bin/bash

for i in ADD SUB MUL DIV
do
	for j in INT32 FLOAT UINT32 INT64 DOUBLE UINT64
	do 	
		for k in 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 
		do
                    NR_DPUS=1 NR_TASKLETS=$k BL=10 OP=$i TYPE=$j make all
                    wait
                    ./bin/host_code -w 0 -e 1 -i 1048576 > profile/${i}_${j}_tl${k}.txt
                    wait
                    make clean
                    wait
		done
	done
done
