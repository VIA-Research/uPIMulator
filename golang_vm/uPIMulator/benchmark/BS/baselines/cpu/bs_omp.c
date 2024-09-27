
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <string.h>
#include <unistd.h>
#include <getopt.h>
#include <assert.h>
#include <time.h>
#include <stdint.h>
#include "timer.h"

#define DTYPE uint64_t
/*
* @brief creates a "test file" by filling a bufferwith values
*/
void create_test_file(DTYPE * input, uint64_t  nr_elements, DTYPE * querys, uint64_t n_querys) {

  uint64_t max = UINT64_MAX;
  uint64_t min = 0;

  srand(time(NULL));

  input[0] = 1;
  for (uint64_t i = 1; i < nr_elements; i++) {
        input[i] = input[i - 1] + (rand() % 10) + 1;
  }

  for(uint64_t i = 0; i < n_querys; i++)
  {
	querys[i] = input[rand() % (nr_elements - 2)];
  }
}

/**
* @brief compute output in the host
*/
uint64_t binarySearch(DTYPE * input, uint64_t input_size, DTYPE* querys, unsigned n_querys)
{

	uint64_t found = -1;
	uint64_t q, r, l, m;
	
       #pragma omp parallel for private(q,r,l,m)
     	for(q = 0; q < n_querys; q++)
      	{
		l = 0;
		r = input_size;
		while (l <= r) 
		{
	    		m = l + (r - l) / 2;

	    		// Check if x is present at mid
	     		if (input[m] == querys[q])
			{	
		    		found += m;
				break;
			}
	    		// If x greater, ignore left half
	    		if (input[m] < querys[q])
			    	l = m + 1;

	    		// If x is smaller, ignore right half
			else
		    		r = m - 1;
		
		}
       	}

      	return found;
}

  /**
  * @brief Main of the Host Application.
  */
  int main(int argc, char **argv) {

    Timer timer;
    uint64_t input_size = atol(argv[1]);
    uint64_t n_querys = atol(argv[2]);

    printf("Vector size: %lu, num searches: %lu\n", input_size, n_querys);
	
    DTYPE * input = malloc((input_size) * sizeof(DTYPE));
    DTYPE * querys = malloc((n_querys) * sizeof(DTYPE));

    DTYPE result_host = -1;

    // Create an input file with arbitrary data.
    create_test_file(input, input_size, querys, n_querys);
	
    start(&timer, 0, 0);
    result_host = binarySearch(input, input_size - 1, querys, n_querys);   
    stop(&timer, 0);


    int status = (result_host);
    if (status) {
        printf("[OK] Execution time: ");
	print(&timer, 0, 1);
	printf("ms.\n");
    } else {
        printf("[ERROR]\n");
    }
    free(input);


    return status ? 0 : 1;
}

