#ifndef BINARY_SEARCH_H
#define BINARY_SEARCH_H

#ifdef _WIN32
  #include <windows.h>
  #define DLL_EXPORT __declspec(dllexport)
#else
  #define DLL_EXPORT
#endif


extern "C" {

    int DLL_EXPORT binary_search(const long int *arr, const long int len, const long int *querys, const long int num_querys);

}

#endif /* BINARY_SEARCH_H */
