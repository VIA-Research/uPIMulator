#!/bin/bash

for i in 1 2 4 8 16 32 64 
do
	for j in 1 
	do 	
		for k in SERIAL PUSH BROADCAST
		do
			for l in 1 4 16 64 256 1024 4096 16384 65536 262144 1048576 4194304 
			do
				NR_DPUS=$i NR_TASKLETS=$j BL=10 TRANSFER=$k make all
				wait
				./bin/host_code -w 5 -e 20 -i ${l} >& profile/${i}_tl${j}_TR${k}_i${l}.txt
				wait
				make clean
				wait
			done
		done
	done
done
