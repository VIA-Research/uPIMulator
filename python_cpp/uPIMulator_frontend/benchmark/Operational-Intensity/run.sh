#!/bin/bash

for i in ADD SUB MUL DIV
do
	for j in CHAR SHORT INT32 FLOAT INT64 DOUBLE
	do 	
		for k in 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16
		do
			for l in 0.001953125 0.00390625 0.0078125 0.015625 0.03125 0.0625 0.125 0.25 0.5 1 2 4 8 16 32 64 128 256 512 
			do
				NR_DPUS=1 NR_TASKLETS=$k BL=10 OP=$i TYPE=$j make all
				wait
				./bin/host_code -w 0 -e 1 -i 1048576 -p ${l} >& profile/${i}_${j}_tl${k}_p${l}.txt
				wait
				make clean
				wait
			done
		done
	done
done
