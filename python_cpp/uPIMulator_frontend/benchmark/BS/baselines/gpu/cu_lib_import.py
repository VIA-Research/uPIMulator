# -*- coding: utf-8 -*-

__all__ = [
    "binary_search",
]


import os.path as path
import platform
from ctypes import *

from numpy.ctypeslib import load_library, ndpointer

## Load the DLL
if platform.system() == "Linux":
    cuda_lib = load_library("cu_binary_search.so", path.dirname(path.realpath(__file__)))
elif platform.system() == "Windows":
    cuda_lib = load_library("cu_binary_search.dll", path.dirname(path.realpath(__file__)))


## Define argtypes for all functions to import
argtype_defs = {
    "binary_search": [ndpointer("i8"), c_int, ndpointer("i8"), c_int],
}


## Import functions from DLL
for func, argtypes in argtype_defs.items():
    locals().update({func: cuda_lib[func]})
    locals()[func].argtypes = argtypes
