#include <stdlib.h>
#include <stdio.h>
#include "../../support/timer.h"
#include "gemv_utils.h"

int main(int argc, char *argv[])
{
  const size_t rows = 20480;
  const size_t cols = 8192;

  double **A, *b, *x;

  b = (double*) malloc(sizeof(double)*rows);
  x = (double*) malloc(sizeof(double)*cols);

  allocate_dense(rows, cols, &A);

  make_hilbert_mat(rows,cols, &A);

#pragma omp parallel
    {
#pragma omp for
    for (size_t i = 0; i < cols; i++) {
      x[i] = (double) i+1 ;
    }

#pragma omp for
    for (size_t i = 0; i < rows; i++) {
      b[i] = (double) 0.0;
    }
    }

  Timer timer;
  start(&timer, 0, 0);


   gemv(A, x, rows, cols, &b);
   
   stop(&timer, 0);


    printf("Kernel ");
    print(&timer, 0, 1);
    printf("\n");

#if 0
  print_vec(x, rows);
  print_mat(A, rows, cols);
  print_vec(b, rows);
#endif

  printf("sum(x) = %f, sum(Ax) = %f\n", sum_vec(x,cols), sum_vec(b,rows));
  return 0;
}

void gemv(double** A, double* x, size_t rows, size_t cols, double** b) {
#pragma omp parallel for
  for (size_t i = 0; i < rows; i ++ )
  for (size_t j = 0; j < cols; j ++ ) {
    (*b)[i] = (*b)[i] + A[i][j]*x[j];
  }
}

void make_hilbert_mat(size_t rows, size_t cols, double*** A) {
#pragma omp parallel for
  for (size_t i = 0; i < rows; i++) {
    for (size_t j = 0; j < cols; j++) {
      (*A)[i][j] = 1.0/( (double) i + (double) j + 1.0);
    }
  }
}

double sum_vec(double* vec, size_t rows) {
  double sum = 0.0;
#pragma omp parallel for reduction(+:sum)
  for (int i = 0; i < rows; i++) sum = sum + vec[i];
  return sum;
}
