
#ifndef _TIMER_H_
#define _TIMER_H_

#include <stdio.h>
#include <sys/time.h>

typedef struct Timer {
    struct timeval startTime;
    struct timeval endTime;
} Timer;

static void startTimer(Timer* timer) {
    gettimeofday(&(timer->startTime), NULL);
}

static void stopTimer(Timer* timer) {
    gettimeofday(&(timer->endTime), NULL);
}

static float getElapsedTime(Timer timer) {
    return ((float) ((timer.endTime.tv_sec - timer.startTime.tv_sec)
                   + (timer.endTime.tv_usec - timer.startTime.tv_usec)/1.0e6));
}

#endif

