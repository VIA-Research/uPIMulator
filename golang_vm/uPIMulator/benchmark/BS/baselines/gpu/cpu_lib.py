# -*- coding: utf-8 -*-


def binary_search(arr, search):

    L = 0
    R = len(arr)

    while L <= R:

        if L > R:
            return -1  # Error code 1

        m = (L + R) / 2
        if arr[m] < search:
            L = m + 1
        elif arr[m] > search:
            R = m - 1
        else:
            return m

    return -2  # Error code 2
