all:
	/usr/local/cuda/bin/nvcc select.cu -I/usr/local/cuda/include -lm -o select -D COARSENING=32 -D THREADS=512 -D INT64 

clean:
	rm select
