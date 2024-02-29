# -*- coding: utf-8 -*-

import time

import numpy as np

# Local Imports
from cu_lib_import import binary_search as gpu_search

# Set an array label to create
arr_len = 2048576
num_querys = 16777216

# Dummy array created
arr = np.arange(0, arr_len, 1).astype("i8")

# Random search querys created
querys = np.random.randint(1, arr_len, num_querys)

# GPU search function call
t0 = time.time()
res_gpu = gpu_search(arr, len(arr), querys, len(querys))
print("Total GPU Time: %i ms" % ((time.time() - t0) * 1e003))
