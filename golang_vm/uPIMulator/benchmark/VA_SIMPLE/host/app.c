#include <dpu.h>
#include <stdlib.h>
#include <assert.h>

#define VECTOR_SIZE 1024

#define NUM_DPUS 1
#define NUM_TASKLETS 1
#define DATA_PREP_PARAMS 1024

struct vector_t {
  int size;
  int* data;
};

struct vector_t* vector_init(int size);
void vector_prep(struct vector_t* vector);
void vector_add(struct vector_t* a, struct vector_t *b, struct vector_t* c);
void vector_equal(struct vector_t* a, struct vector_t* b);

int main() {
  struct dpu_set_t dpu_set;
  struct dpu_set_t dpu;

  dpu_alloc(NUM_DPUS, nullptr, &dpu_set);

  dpu_load(dpu_set, "../device/device", nullptr);

  struct vector_t *A = vector_init(VECTOR_SIZE);
  struct vector_t *B = vector_init(VECTOR_SIZE);
  struct vector_t *C_host = vector_init(VECTOR_SIZE);
  struct vector_t *C_device = vector_init(VECTOR_SIZE);

  vector_prep(A);
  vector_prep(B);

  vector_add(A, B, C_host);

  int size_per_dpu = (VECTOR_SIZE / NUM_DPUS) * sizeof(int);
  int i;

  DPU_FOREACH(dpu_set, dpu, i) { dpu_prepare_xfer(dpu, &size_per_dpu); }
  dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, "size_per_dpu", 0, sizeof(int), DPU_XFER_DEFAULT);

  DPU_FOREACH(dpu_set, dpu, i) { dpu_prepare_xfer(dpu, &(A->data[(size_per_dpu / sizeof(int)) * i])); }
  dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, DPU_MRAM_HEAP_POINTER_NAME, 0, size_per_dpu, DPU_XFER_DEFAULT);

  DPU_FOREACH(dpu_set, dpu, i) { dpu_prepare_xfer(dpu, &(B->data[(size_per_dpu / sizeof(int)) * i])); }
  dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, DPU_MRAM_HEAP_POINTER_NAME, size_per_dpu, size_per_dpu, DPU_XFER_DEFAULT);

  dpu_launch(dpu_set, DPU_SYNCHRONOUS);

  DPU_FOREACH(dpu_set, dpu, i) { dpu_prepare_xfer(dpu, &(C_device->data[(size_per_dpu / sizeof(int)) * i])); }
  dpu_push_xfer(dpu_set, DPU_XFER_FROM_DPU, DPU_MRAM_HEAP_POINTER_NAME, 2 * size_per_dpu, size_per_dpu, DPU_XFER_DEFAULT);

  vector_equal(C_host, C_device);

  return 0;
}

struct vector_t* vector_init(int size) {
  struct vector_t* vector = malloc(sizeof(struct vector_t));

  vector->size = size;
  vector->data = malloc(VECTOR_SIZE * sizeof(int));

  return vector;
}

void vector_prep(struct vector_t* vector) {
  for (int i = 0; i < vector->size; i++) {
    vector->data[i] = i;
  }
}

void vector_add(struct vector_t* a, struct vector_t* b, struct vector_t* c) {
  assert(a->size == b->size);
  assert(a->size == c->size);

  for (int i = 0; i < c->size; i++) {
    c->data[i] = a->data[i] + b->data[i];
  }
}

void vector_equal(struct vector_t *a, struct vector_t *b) {
  assert(a->size == b->size);

  for (int i = 0; i < a->size; i++) {
    assert(a->data[i] == b->data[i]);
  }
}
