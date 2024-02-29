#include <cuComplex.h>
#include <cufft.h>
#include <vector>
#include <stdio.h>
#include <cuda.h>
#include <thrust/functional.h>
#include <thrust/device_ptr.h>
#include <thrust/transform_scan.h>
#include <thrust/sequence.h>
#include <float.h>
#include <chrono>

using std::vector;

static const int THREADS_PER_BLOCK = 1024;

// Holds matrix profile and index values together
typedef union  {
  float floats[2];                 // floats[0] = lowest
  unsigned int ints[2];                     // ints[1] = lowIdx
  unsigned long long int ulong;    // for atomic update
} mp_entry;

struct MPIDXCombine
{
	__host__ __device__
	mp_entry operator()(double x, unsigned int idx){
		mp_entry item;
		item.floats[0] = (float) x;
		item.ints[1] = idx;
		return item;
	}
};

//Atomically updates the MP/idxs using a single 64-bit integer. We lose a small amount of precision in the output, if we do not do this we are unable
// to atomically update both the matrix profile and the indexes without using a critical section and dedicated locks.
__device__ inline unsigned long long int MPatomicMin(volatile unsigned long long int* address, double val, unsigned int idx)
{
	float fval = (float)val;
	mp_entry loc, loctest;
	loc.floats[0] = fval;
	loc.ints[1] = idx;
	loctest.ulong = *address;
	while (loctest.floats[0] > fval){
		loctest.ulong = atomicCAS((unsigned long long int*) address, loctest.ulong,  loc.ulong);
	}
	return loctest.ulong;
}

//This macro checks return value of the CUDA runtime call and exits
//the application if the call failed.
#define gpuErrchk(ans) { gpuAssert((ans), __FILE__, __LINE__); }
inline void gpuAssert(cudaError_t code, const char *file, int line, bool abort=true)
{
   if (code != cudaSuccess)
   {
      fprintf(stderr,"GPUassert: %s %s %d\n", cudaGetErrorString(code), file, line);
      if (abort) exit(code);
   }
}

//This kernel computes a sliding mean with specified window size and a corresponding prefix sum array (A)
template<class DTYPE>
__global__ void sliding_mean(DTYPE* pref_sum,  size_t window, size_t size, DTYPE* means)
{
	const DTYPE coeff = 1.0 / (DTYPE) window;
	size_t a = blockIdx.x * blockDim.x + threadIdx.x;
	size_t b = blockIdx.x * blockDim.x + threadIdx.x + window;

	if(a == 0){
		means[a] = pref_sum[window - 1] * coeff;
	}
	if(a < size - 1){
		means[a + 1] = (pref_sum[b] - pref_sum[a]) * coeff;
	}
}

//This kernel computes a sliding standard deviaiton with specified window size, the corresponding means of each element, and the prefix squared sum at each element
template<class DTYPE>
__global__ void sliding_std(DTYPE* squares, size_t window, size_t size, DTYPE* means, DTYPE* stds){
	const DTYPE coeff = 1 / (DTYPE) window;
	size_t a = blockIdx.x * blockDim.x + threadIdx.x;
	size_t b = blockIdx.x * blockDim.x + threadIdx.x + window;
	if(a == 0){
		stds[a] = sqrt((squares[window - 1] * coeff) - (means[a] * means[a]));
	}
	else if(b < size + window) {
		stds[a] = sqrt(((squares[b - 1] - squares[a - 1]) * coeff) - (means[a] * means[a]));
    }
}


template<class DTYPE>
__global__ void elementwise_multiply_inplace(const DTYPE* A, DTYPE *B, const int size)
{
    int tid = blockIdx.x * blockDim.x + threadIdx.x;
    if(tid < size) {
       B[tid] *= A[tid];
    }
}

template<>
__global__ void elementwise_multiply_inplace(const cuDoubleComplex* A, cuDoubleComplex* B, const int size)
{
    int tid = blockIdx.x * blockDim.x + threadIdx.x;
    if(tid < size) {
       B[tid] = cuCmul(A[tid], B[tid]);
    }
}


// A is input unaligned sliding dot products produced by ifft
// out is the computed vector of distances
template<class DTYPE>
__global__ void normalized_aligned_distance(const DTYPE* A, DTYPE* out, DTYPE * lastzs,
  const DTYPE * AMean, const DTYPE* ASigma,
  const unsigned int windowSize, const int exclusionZone,
  const unsigned int ProfileLength, DTYPE* profile,
  unsigned int * profile_idx, const unsigned int scratch, mp_entry *profile_entry)
  {

    int thID = blockIdx.x * blockDim.x + threadIdx.x;
    int i = 1;
    int j = thID + i;

    DTYPE lastz = lastzs[thID];

    if(j > exclusionZone)
    {
//      while(j < ProfileLength)
//      {
        lastz  = lastz + (A[j + windowSize - 1] * A[i + windowSize - 1]) - (A[j - 1] * A[i - 1]);
        DTYPE distance = max(2 * (windowSize - (lastz -  AMean[j] * AMean[i] * windowSize) / (ASigma[j] * ASigma[i])), 0.0);

        if (distance < profile_entry[j].floats[0])
        {
          MPatomicMin((unsigned long long int*)&profile_entry[j], distance, i);
        }
        if (distance < profile_entry[i].floats[0])
        {
          MPatomicMin((unsigned long long int*)&profile_entry[i], distance, j);
        }
        i++;
        j++;
  //    }
    }
  }


template<class DTYPE>
__global__ void initialize_lastzs(const DTYPE* A, DTYPE* out, DTYPE * lastzs_last,
  const DTYPE * AMean, const DTYPE* ASigma,  const unsigned int windowSize, const unsigned int exclusionZone,
  const unsigned int ProfileLength, DTYPE* profile,
  unsigned int * profile_idx)
{
    int j = blockIdx.x * blockDim.x + threadIdx.x;

    if((j > exclusionZone) && (j < ProfileLength)) {
       DTYPE lastz = 0;
       for (int index = j; index < windowSize + j; index++)
       {
         lastz += A[index] * A[index-j];
       }

       DTYPE distance = max(2 * (windowSize - (lastz -  AMean[j] * AMean[0] * windowSize) / (ASigma[j] * ASigma[0])), 0.0);
       // Update the distance profile
       out[j] = distance;
       // Update the matrix profile if needed
       if(profile[j] > distance) {
         profile[j] = distance;
         profile_idx[j] = 0;
       }
       if(j < ProfileLength) lastzs_last[j] = lastz;

    }
    else if (j < ProfileLength)
    {
      out[j] = DBL_MAX;
    }
}


template<class DTYPE>
__host__ void distance_profile(const DTYPE* A, DTYPE* QT, DTYPE * lastzs,
  DTYPE *profile, unsigned int *profile_idx, const DTYPE * AMean, const DTYPE * ASigma, const int timeSeriesLength,
  const int windowSize,const int exclusionZone, const unsigned int i, mp_entry *profile_entry)
  {
    const int ProfileLength = timeSeriesLength - windowSize + 1;

    dim3 grid(ceil(ProfileLength / (float) THREADS_PER_BLOCK), 1, 1);
    dim3 block(THREADS_PER_BLOCK, 1, 1);

    normalized_aligned_distance<DTYPE><<<grid, block>>>(A, QT, lastzs, AMean, ASigma,windowSize,
      exclusionZone, ProfileLength,profile, profile_idx, i, profile_entry);
    gpuErrchk(cudaPeekAtLastError());

  }

// Reduction kernel, upper layer
// This reduction was adapted from the nvidia whitepaper:
// http://developer.download.nvidia.com/compute/cuda/1.1-Beta/x86_website/projects/reduction/doc/reduction.pdf
template <class DTYPE, unsigned int blockSize>
__global__ void reduce(const DTYPE *g_idata, DTYPE *g_odata, unsigned int *g_oloc,  unsigned int ProfileLength) {
	__shared__ DTYPE sdata[blockSize];
	__shared__ DTYPE sloc[blockSize];
	unsigned int tid = threadIdx.x;
	unsigned int i = blockIdx.x*(blockSize*2) + tid;
	unsigned int gridSize = blockSize*2*gridDim.x;
	DTYPE temp;
	unsigned int temploc;
	sdata[tid] = DBL_MAX;
	while (i < ProfileLength) {
		if (i + blockSize < ProfileLength)
		{
			if (g_idata[i] < g_idata[i+blockSize])
			{
				temp=g_idata[i];
				temploc=i;
			}
			else
			{
				temp=g_idata[i+blockSize];
				temploc = i+blockSize;
			}
		}
		else
		{
			temp = g_idata[i];
			temploc = i;
		}
		if (sdata[tid] > temp)
		{
			sdata[tid] = temp;
			sloc[tid] = temploc;
		}
		i += gridSize;
	}
	__syncthreads();
	if (blockSize >= 1024) {
		if (tid < 512 && sdata[tid] > sdata[tid + 512])
		{
			sdata[tid] = sdata[tid + 512];
			sloc[tid] = sloc[tid + 512];
		}
		__syncthreads();
	}
	if (blockSize >= 512 ) {
		if (tid < 256 && sdata[tid] > sdata[tid + 256])
		{
			sdata[tid] = sdata[tid + 256];
			sloc[tid] = sloc[tid + 256];
		}
		__syncthreads();
	}
	if (blockSize >= 256) {
		if (tid < 128 && sdata[tid] > sdata[tid + 128])
		{
			sdata[tid] = sdata[tid + 128];
			sloc[tid] = sloc[tid + 128];
		}
		__syncthreads();
	}
	if (blockSize >= 128) {
		if (tid < 64 && sdata[tid] > sdata[tid + 64])
		{
			sdata[tid] = sdata[tid + 64];
			sloc[tid] = sloc[tid + 64];
		}
		__syncthreads();
	}

	if (blockSize >= 64) {
		if (tid < 32 && sdata[tid] > sdata[tid + 32])
		{
			sdata[tid] = sdata[tid + 32];
			sloc[tid] = sloc[tid + 32];
		}
		__syncthreads();
	}

	if (blockSize >= 32) {
		if (tid < 16 && sdata[tid] > sdata[tid + 16])
		{
			sdata[tid] = sdata[tid + 16];
			sloc[tid] = sloc[tid + 16];
		}
		__syncthreads();
	}

	if (blockSize >= 16) {
		if (tid < 8 && sdata[tid] > sdata[tid + 8])
		{
			sdata[tid] = sdata[tid + 8];
			sloc[tid] = sloc[tid + 8];
		}
		__syncthreads();
	}

	if (blockSize >= 8) {
		if (tid < 4 && sdata[tid] > sdata[tid + 4])
		{
			sdata[tid] = sdata[tid + 4];
			sloc[tid] = sloc[tid + 4];
		}
		__syncthreads();
	}

	if (blockSize >= 4) {
		if (tid < 2 && sdata[tid] > sdata[tid + 2])
		{
			sdata[tid] = sdata[tid + 2];
			sloc[tid] = sloc[tid + 2];
		}
		__syncthreads();
	}

	if (blockSize >= 2) {
		if (tid == 0)
		{
			if (sdata[0] <= sdata[1])
			{
				g_odata[blockIdx.x] = sdata[0];
				g_oloc[blockIdx.x] = sloc[0];
			}
			else
			{
				g_odata[blockIdx.x] = sdata[1];
				g_oloc[blockIdx.x] = sloc[1];
			}
		}
	}
	else
	{
		if (tid == 0)
		{
			g_odata[blockIdx.x] = sdata[0];
			g_oloc[blockIdx.x] = sloc[0];
		}
	}
}

//reduction kernel, lower layer
template <class DTYPE, unsigned int blockSize>
__global__ void reducelast(DTYPE *g_idata, unsigned int *g_iloc,
  unsigned int start_loc, DTYPE* profilei, unsigned int* profileidxi, unsigned int n) {

	__shared__ DTYPE sdata[blockSize];
	__shared__ DTYPE sloc[blockSize];
	unsigned int tid = threadIdx.x;
	unsigned int i = blockIdx.x*(blockSize*2) + tid;
	unsigned int gridSize = blockSize*2*gridDim.x;
	DTYPE temp;
	unsigned int temploc;
	sdata[tid] = DBL_MAX;
	DTYPE minval;
	unsigned int minloc;
	while (i < n) {
		if (i + blockSize <n)
		{
			if (g_idata[i] < g_idata[i+blockSize])
			{
				temp=g_idata[i];
				temploc=g_iloc[i];
			}
			else
			{
				temp=g_idata[i+blockSize];
				temploc = g_iloc[i+blockSize];
			}
		}
		else
		{
			temp = g_idata[i];
			temploc = g_iloc[i];
		}
		if (sdata[tid] > temp)
		{
			sdata[tid] = temp;
			sloc[tid] = temploc;
		}
		i += gridSize;
	}
	__syncthreads();
	if (blockSize >= 1024) {
		if (tid < 512 && sdata[tid] > sdata[tid + 512])
		{
			sdata[tid] = sdata[tid + 512];
			sloc[tid] = sloc[tid + 512];
		}
		__syncthreads();
	}
	if (blockSize >= 512 ) {
		if (tid < 256 && sdata[tid] > sdata[tid + 256])
		{
			sdata[tid] = sdata[tid + 256];
			sloc[tid] = sloc[tid + 256];
		}
		__syncthreads();
	}
	if (blockSize >= 256) {
		if (tid < 128 && sdata[tid] > sdata[tid + 128])
		{
			sdata[tid] = sdata[tid + 128];
			sloc[tid] = sloc[tid + 128];
		}
		__syncthreads();
	}
	if (blockSize >= 128) {
		if (tid < 64 && sdata[tid] > sdata[tid + 64])
		{
			sdata[tid] = sdata[tid + 64];
			sloc[tid] = sloc[tid + 64];
		}
		__syncthreads();
	}

	if (blockSize >= 64) {
		if (tid < 32 && sdata[tid] > sdata[tid + 32])
		{
			sdata[tid] = sdata[tid + 32];
			sloc[tid] = sloc[tid + 32];
		}
		__syncthreads();
	}

	if (blockSize >= 32) {
		if (tid < 16 && sdata[tid] > sdata[tid + 16])
		{
			sdata[tid] = sdata[tid + 16];
			sloc[tid] = sloc[tid + 16];
		}
		__syncthreads();
	}

	if (blockSize >= 16) {
		if (tid < 8 && sdata[tid] > sdata[tid + 8])
		{
			sdata[tid] = sdata[tid + 8];
			sloc[tid] = sloc[tid + 8];
		}
		__syncthreads();
	}

	if (blockSize >= 8) {
		if (tid < 4 && sdata[tid] > sdata[tid + 4])
		{
			sdata[tid] = sdata[tid + 4];
			sloc[tid] = sloc[tid + 4];
		}
		__syncthreads();
	}

	if (blockSize >= 4) {
		if (tid < 2 && sdata[tid] > sdata[tid + 2])
		{
			sdata[tid] = sdata[tid + 2];
			sloc[tid] = sloc[tid + 2];
		}
		__syncthreads();
	}

	if (blockSize >= 2) {
		if (tid == 0)
		{
			if (sdata[0] <= sdata[1])
			{
				minval = sdata[0];
				minloc = sloc[0];
			}
			else
			{
				minval = sdata[1];
				minloc = sloc[1];
			}
		}
	}
	else
	{
		if (tid == 0)
		{
			minval = sdata[0];
			minloc = sloc[0];
		}
	}

	if (tid==0)
	{
		if (minval<(*profilei))
		{
			(*profilei)=minval;
			(*profileidxi)=minloc+start_loc;
		}
	}

}

template<class DTYPE>
void reducemain(DTYPE* vd, unsigned int start_loc, unsigned int max_block_num, unsigned int max_thread_num, unsigned int n, DTYPE* profile, unsigned int* profileidx, unsigned int i, DTYPE* reduced_result, unsigned int* reduced_loc)
{

	if (n==0) //if this happens, there's an error
		return;
	if (max_thread_num>1024)
		max_thread_num=1024;

	unsigned int * middle_loc_pointer=reduced_loc;


	unsigned int num_threads=max_thread_num;

	unsigned int num_blocks=n/(num_threads*2);
	if (n%(num_threads*2)!=0)
		num_blocks++;
	if (num_blocks>=max_block_num)
		num_blocks=max_block_num;
	DTYPE *middle_pointer = NULL;
	unsigned int curn;
	if (num_blocks>1) //upperlevel reduction
	{
		middle_pointer=reduced_result;
		curn=num_blocks;
		switch (num_threads)
		{
			case 1024:
				reduce<DTYPE, 1024><<<num_blocks,1024>>>(vd + start_loc,reduced_result,reduced_loc,n); break;
			case 512:
				reduce<DTYPE, 512><<<num_blocks,512>>>(vd + start_loc,reduced_result,reduced_loc,n); break;
			case 256:
				reduce<DTYPE, 256><<<num_blocks,256>>>(vd+start_loc,reduced_result,reduced_loc,n); break;
			case 128:
				reduce<DTYPE, 128><<<num_blocks,128>>>(vd+start_loc,reduced_result,reduced_loc,n); break;
			case 64:
				reduce<DTYPE, 64><<<num_blocks,64>>>(vd+start_loc,reduced_result,reduced_loc,n); break;
			case 32:
				reduce<DTYPE, 32><<<num_blocks,32>>>(vd+start_loc,reduced_result,reduced_loc,n); break;
			case 16:
				reduce<DTYPE, 16><<<num_blocks,16>>>(vd+start_loc,reduced_result,reduced_loc,n); break;
			case 8:
				reduce<DTYPE, 8><<<num_blocks,8>>>(vd+start_loc,reduced_result,reduced_loc,n); break;
			case 4:
				reduce<DTYPE, 4><<<num_blocks,4>>>(vd+start_loc,reduced_result,reduced_loc,n); break;
			case 2:
				reduce<DTYPE, 2><<<num_blocks,2>>>(vd+start_loc,reduced_result,reduced_loc,n); break;
			default:
				break;
		}
	        gpuErrchk( cudaPeekAtLastError() );
	}
	else
	{
		middle_pointer=vd+start_loc;
		curn=n;
        auto ptr = thrust::device_pointer_cast(reduced_loc);
		thrust::sequence(ptr,ptr+curn);
	}


	num_threads=floor(pow(2,ceil(log(curn)/log(2))-1));
	if (num_threads>max_thread_num)
		num_threads=max_thread_num;
	switch (num_threads)
	{
		case 1024:
			reducelast<DTYPE,1024><<<1,1024>>>(middle_pointer, middle_loc_pointer, start_loc, profile+i, profileidx+i, curn);
	        	gpuErrchk( cudaPeekAtLastError() );
			break;
		case 512:
			reducelast<DTYPE,512><<<1,512>>>(middle_pointer,middle_loc_pointer, start_loc, profile+i, profileidx+i,  curn);
	        	gpuErrchk( cudaPeekAtLastError() );
			break;
		case 256:
			reducelast<DTYPE,256><<<1,256>>>(middle_pointer,middle_loc_pointer, start_loc, profile+i, profileidx+i, curn);
	        	gpuErrchk( cudaPeekAtLastError() );
			break;
		case 128:
			reducelast<DTYPE,128><<<1,128>>>(middle_pointer,middle_loc_pointer, start_loc,  profile+i, profileidx+i,curn);
	        	gpuErrchk( cudaPeekAtLastError() );
			break;
		case 64:
			reducelast<DTYPE,64><<<1,64>>>(middle_pointer,middle_loc_pointer, start_loc,  profile+i, profileidx+i,curn);
	        	gpuErrchk( cudaPeekAtLastError() );
			break;
		case 32:
			reducelast<DTYPE,32><<<1,32>>>(middle_pointer,middle_loc_pointer, start_loc, profile+i, profileidx+i, curn);
	        	gpuErrchk( cudaPeekAtLastError() );
			break;
		case 16:
			reducelast<DTYPE,16><<<1,16>>>(middle_pointer,middle_loc_pointer, start_loc, profile+i, profileidx+i,curn);
	        	gpuErrchk( cudaPeekAtLastError() );
			break;
		case 8:
			reducelast<DTYPE,8><<<1,8>>>(middle_pointer,middle_loc_pointer, start_loc, profile+i, profileidx+i,curn);
	        	gpuErrchk( cudaPeekAtLastError() );
			break;
		case 4:
			reducelast<DTYPE,4><<<1,4>>>(middle_pointer,middle_loc_pointer, start_loc,  profile+i, profileidx+i,curn);
	        	gpuErrchk( cudaPeekAtLastError() );
			break;
		case 2:
			reducelast<DTYPE,2><<<1,2>>>(middle_pointer,middle_loc_pointer, start_loc, profile+i, profileidx+i, curn);
	        	gpuErrchk( cudaPeekAtLastError() );
			break;
		case 1:
			reducelast<DTYPE,1><<<1,1>>>(middle_pointer,middle_loc_pointer, start_loc, profile+i, profileidx+i, curn);
	        	gpuErrchk( cudaPeekAtLastError() );
			break;
		case 0:
			reducelast<DTYPE,1><<<1,1>>>(middle_pointer,middle_loc_pointer, start_loc, profile+i, profileidx+i, curn);
	        	gpuErrchk( cudaPeekAtLastError() );
			break;
		default:
			break;
	}
}

template<class DTYPE>
struct square_op : public thrust::unary_function<DTYPE,DTYPE>
{
  __host__ __device__
  DTYPE operator()(DTYPE x) const
  {
    return x * x;
  }
};

template<class DTYPE>
void compute_statistics(const DTYPE *T, DTYPE *means, DTYPE *stds, DTYPE *scratch, size_t n, size_t m)
{
    square_op<DTYPE> sqr;
    dim3 grid(ceil(n / (double) THREADS_PER_BLOCK), 1,1);
    dim3 block(THREADS_PER_BLOCK, 1, 1);
    thrust::device_ptr<const DTYPE> dev_ptr_T = thrust::device_pointer_cast(T);
    thrust::device_ptr<DTYPE> dev_ptr_scratch = thrust::device_pointer_cast(scratch);

	thrust::inclusive_scan(dev_ptr_T, dev_ptr_T + n + m - 1, dev_ptr_scratch, thrust::plus<DTYPE>());
    sliding_mean<DTYPE><<<grid, block>>>(scratch, m, n, means);
	thrust::transform_inclusive_scan(dev_ptr_T, dev_ptr_T + n + m - 1, dev_ptr_scratch, sqr,thrust::plus<DTYPE>());
    sliding_std<DTYPE><<<grid,block>>>(scratch, m, n, means, stds);
}



template<class DTYPE>
void STREAMP(DTYPE* T, const int timeSeriesLength, const int windowSize, DTYPE* profile, unsigned int* profile_idxs,  mp_entry *profile_with_idx)
{

  int exclusionZone = windowSize / 4;
  size_t ProfileLength = timeSeriesLength - windowSize + 1;
  DTYPE * AMean, * ASigma, *QT, *lastzs, *reduced_result;

  dim3 block(THREADS_PER_BLOCK,1,1);
  dim3 grid(ceil(ProfileLength / (float) THREADS_PER_BLOCK), 1, 1);

  unsigned int *reduced_loc;

  //clock_t start, now;
  const unsigned int max_block_num=2048;
  const unsigned int max_thread_num=1024;
  unsigned int middle_loc_size=max_block_num>max_thread_num?max_block_num:max_thread_num;
 // printf("size = %d, window = %d, exclusion = %d\n", ProfileLength, windowSize, exclusionZone);

  //start = clock();

  cudaMalloc(&QT, ProfileLength * sizeof(DTYPE));
  cudaMalloc(&AMean, ProfileLength * sizeof(DTYPE));
  cudaMalloc(&ASigma, ProfileLength * sizeof(DTYPE));
  cudaMalloc(&lastzs, ProfileLength * sizeof(DTYPE));

  cudaMalloc(&reduced_result, max_block_num * sizeof(DTYPE));
  cudaMalloc(&reduced_loc, middle_loc_size * sizeof(unsigned int));

  //now = clock();
  //printf("Allocate memory took %lf sec\n", (now - start) / (double) CLOCKS_PER_SEC);

  // Precompute statistics
  //start = clock();

  //Use QT vector as scratch space as we don't need it yet
  compute_statistics(T, AMean, ASigma, QT, ProfileLength, windowSize);
  //now = clock();
 // printf("Precompute statistics took %lf sec\n", (now - start) / (double) CLOCKS_PER_SEC);

  // Initialize profile and lastzs_last
 // start = clock();
   auto begin = std::chrono::high_resolution_clock::now();

  initialize_lastzs<DTYPE><<<grid, block>>>(T, QT, lastzs, AMean, ASigma,  windowSize, exclusionZone,
    ProfileLength, profile, profile_idxs);

  reducemain(QT, 0, 2048, 1024, ProfileLength, profile, profile_idxs, 0, reduced_result, reduced_loc);

  MPIDXCombine combiner;
  auto ptr_prof = thrust::device_pointer_cast(profile);
  auto ptr_idx = thrust::device_pointer_cast(profile_idxs);
  auto ptr_comb = thrust::device_pointer_cast(profile_with_idx);
  thrust::transform(ptr_prof, ptr_prof + ProfileLength, ptr_idx, ptr_comb, combiner);

  cudaDeviceSynchronize();

  // compute the distance profile
  distance_profile<DTYPE>(T, QT, lastzs, profile, profile_idxs, AMean, ASigma, timeSeriesLength,
    windowSize, exclusionZone, 1, profile_with_idx);

  cudaDeviceSynchronize();
  //now = clock();
  auto end = std::chrono::high_resolution_clock::now(); 
  std::cout << "STREAMP time: "<< (float) std::chrono::duration_cast<std::chrono::microseconds>(end-begin).count() / 1000 << " ms." << std::endl;

  cudaFree(QT);
  cudaFree(AMean);
  cudaFree(ASigma);
  cudaFree(lastzs);
}

//Reads input time series from file
template<class DTYPE>
void readFile(const char* filename, vector<DTYPE>& v, const char *format_str)
{
	FILE* f = fopen( filename, "r");
	if(f == NULL){
		printf("Unable to open %s for reading, please make sure it exists\n", filename);
		exit(0);
	}
	DTYPE num;
	while(!feof(f)){
			fscanf(f, format_str, &num);
			v.push_back(num);
    }
	v.pop_back();
	fclose(f);
}

int main(int argc, char **argv)
{
  if (argc != 4) {
    printf("Usage: <subseq length> <input file> <output file>\n");
    exit(0);
  }

  int nDevices;
  double *T, *profile;
  unsigned int *idxs;
  mp_entry *profile_with_idx;
  int windowSize = atoi(argv[1]);
  char *filename = argv[2];
  //clock_t start, now;
  vector<double> T_host;


  cudaGetDeviceCount(&nDevices);
  vector<cudaDeviceProp> device_info(nDevices);

 /* printf("Number of CUDA devices: %d\n",nDevices);

  for (int i = 0; i < nDevices; ++i) {
    cudaGetDeviceProperties(&device_info.at(i), i);
    printf("Device Number: %d\n", i);
    printf("  Device name: %s\n", device_info.at(i).name);
    printf("  Memory Clock Rate (KHz): %d\n",
    device_info.at(i).memoryClockRate);
    printf("  Memory Bus Width (bits): %d\n",
    device_info.at(i).memoryBusWidth);
    printf("  Peak Memory Bandwidth (GB/s): %f\n\n",
      2.0*device_info.at(i).memoryClockRate*(device_info.at(i).memoryBusWidth/8)/1.0e6);
  }*/

 // std::cout << "Enter the device number to use: " << '\n';
  //std::cin >> selectedDevice;

  //cudaSetDevice(selectedDevice);
  cudaSetDevice(0);
  cudaFree(0);

  //start = clock();
  readFile<double>(filename, T_host, "%lf");
  //now = clock();

 // printf("Time taken to read date from file: %lf seconds\n", (now - start) / (double) CLOCKS_PER_SEC);

  vector<double> profile_host(T_host.size() - windowSize + 1, DBL_MAX);
  vector<unsigned int> index_host(profile_host.size(), 0);
  vector<mp_entry> profile_with_idx_h(profile_host.size());

  //start = clock();
  cudaMalloc(&T, T_host.size() * sizeof(double));
  cudaMemcpy(T, T_host.data(), T_host.size() * sizeof(double), cudaMemcpyHostToDevice);
  cudaMalloc(&profile, profile_host.size() * sizeof(double));
  cudaMemcpy(profile, profile_host.data(), profile_host.size() * sizeof(double), cudaMemcpyHostToDevice);
  cudaMalloc(&idxs, index_host.size() * sizeof(unsigned int));
  cudaMalloc(&profile_with_idx, profile_host.size() * sizeof(mp_entry));
  //now = clock();

 // printf("Time taken to allocate T and profile and transfer to device: %lf seconds\n", (now - start) / (double) CLOCKS_PER_SEC);

  // Do SCRIMP
  STREAMP<double>(T, T_host.size(), windowSize, profile, idxs, profile_with_idx);

  //start = clock();
  cudaMemcpy(&profile_with_idx_h[0], profile_with_idx, profile_host.size() * sizeof(mp_entry), cudaMemcpyDeviceToHost);
  //now = clock();

  //printf("Time taken to copy result to host: %lf seconds\n", (now - start) / (double) CLOCKS_PER_SEC);

  //printf("writing result to files\n");
  FILE* f1 = fopen( argv[3], "w");
  for(int i = 0; i < profile_host.size(); ++i){
    fprintf(f1, "%.10f %u\n", sqrt(profile_with_idx_h[i].floats[0]) , profile_with_idx_h[i].ints[1]);
  }

  fclose(f1);
  gpuErrchk(cudaDeviceSynchronize());
  gpuErrchk(cudaDeviceReset());

  cudaFree(T);
  cudaFree(profile);
  cudaFree(profile_with_idx);
}
