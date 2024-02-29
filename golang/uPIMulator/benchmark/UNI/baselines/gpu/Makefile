all:
	/usr/local/cuda/bin/nvcc unique.cu -I/usr/local/cuda/include -lm -o unique -D COARSENING=32 -D THREADS=512 -D INT64 

clean:
	rm unique
