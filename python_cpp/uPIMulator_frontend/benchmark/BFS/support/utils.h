
#ifndef _UTILS_H_
#define _UTILS_H_

#define PRINT_ERROR(fmt, ...)       fprintf(stderr, "\033[0;31mERROR:\033[0m   " fmt "\n", ##__VA_ARGS__)
#define PRINT_WARNING(fmt, ...)     fprintf(stderr, "\033[0;35mWARNING:\033[0m " fmt "\n", ##__VA_ARGS__)
#define PRINT_INFO(cond, fmt, ...)  if(cond) printf("\033[0;32mINFO:\033[0m    " fmt "\n", ##__VA_ARGS__);
#define PRINT(fmt, ...)             printf(fmt "\n", ##__VA_ARGS__)

#endif

