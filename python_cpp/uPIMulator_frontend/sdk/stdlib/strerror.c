/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <stddef.h>
#include <defs.h>
#include "errno.h"

// http://www.delorie.com/gnu/docs/glibc/libc_17.html

// static char* strerror_errors[] = { "Success",
//"E2BIG", "EACCES", "EADDRINUSE", "EADDRNOTAVAIL", "EAFNOSUPPORT", "EAGAIN", "EALREADY", "EBADF",
//"EBADMSG", "EBUSY", "ECANCELED", "ECHILD", "ECONNABORTED", "ECONNREFUSED", "ECONNRESET", "EDEADLK",
//"EDESTADDRREQ", "EDOM", "EDQUOT", "EEXIST", "EFAULT", "EFBIG", "EHOSTUNREACH", "EIDRM",
//"EILSEQ", "EINPROGRESS", "EINTR", "EINVAL", "EIO", "EISCONN", "EISDIR", "ELOOP",
//"EMFILE", "EMLINK", "EMSGSIZE", "EMULTIHOP", "ENAMETOOLONG", "ENETDOWN", "ENETRESET", "ENETUNREACH",
//"ENFILE", "ENOBUFS", "ENODATA", "ENODEV", "ENOENT", "ENOEXEC", "ENOLCK", "ENOLINK",
//"ENOMEM", "ENOMSG", "ENOPROTOOPT", "ENOSPC", "ENOSR", "ENOSTR", "ENOSYS", "ENOTCONN",
//"ENOTDIR", "ENOTEMPTY", "ENOTRECOVERABLE", "ENOTSOCK", "ENOTSUP", "ENOTTY", "ENXIO", "EOPNOTSUPP",
//"EOVERFLOW", "EOWNERDEAD", "EPERM", "EPIPE", "EPROTO", "EPROTONOSUPPORT", "EPROTOTYPE", "ERANGE",
//"EROFS", "ESPIPE", "ESRCH", "ESTALE", "ETIME", "ETIMEDOUT", "ETXTBSY", "EWOULDBLOCK",
//"EXDEV"
//};

// 81 errors in total including 2 not supported ones.

const static char *strerror_errors_complete[] = { "Success",
    "Argument list too long",
    "Permission denied",
    "Address in use",
    "Address not available",
    "Address family not supported",
    "Resource unavailable, try again",
    "Connection already in progress",
    "Bad file descriptor",

    "Bad message",
    "Device or resource busy",
    "Operation canceled",
    "No child processes",
    "Connection aborted",
    "Connection refused",
    "Connection reset",
    "Resource deadlock would occur",

    "Destination address required",
    "Mathematics argument out of domain of function",
    "Disk quota exceeded",
    "File exists",
    "Bad address",
    "File too big",
    "Host unreachable",
    "Identifier removed",

    "Illegal byte sequence",
    "Operation in progress",
    "Interrupted function",
    "Invalid argument",
    "I/O error",
    "Socket is connected",
    "File is a directory",
    "Too many levels of symbolic links",

    "File descriptor value too large",
    "Too many links",
    "Message too large",
    "EMULTIHOP",
    "Filename too long",
    "Network is down",
    "Connection aborted by network",
    "Network unreachable",

    "Too many files open in system",
    "No buffer space available",
    "No message is available on the STREAM head read queue",
    "No such device",
    "No such file or directory",
    "Executable file format error",
    "No locks available",
    "ENOLINK",

    "Not enough space",
    "No message of the desired type",
    "Protocol unavailable",
    "No space left on device",
    "No STREAM resources",
    "Not a STREAM",
    "Function not supported",
    "The socket is not connected",

    "Not a directory or a symbolic link to a directory",
    "Directory not empty",
    "State not recoverable",
    "Not a socket",
    "Not supported",
    "Inappropriate I/O control operation",
    "No such device or address",
    "Operation not supported on socket",

    "Value too large to be stored in data type",
    "Previous owner died",
    "Operation not permitted",
    "Broken pipe",
    "Protocol error",
    "Protocol not supported",
    "Protocol wrong type for socket",
    "Result too large",

    "Read-only file system",
    "Invalid seek",
    "No such process",
    "ESTALE",
    "Stream ioctl() timeout",
    "Connection timed out",
    "Text file busy",
    "Operation would block",

    "Cross-device link",
    "Unknown error" };

char *
strerror(int errnum)
{
    unsigned int length = sizeof(strerror_errors_complete) / sizeof(strerror_errors_complete[0]) - 1; //-1 for "Unknown error"
    if (((unsigned int)errnum) >= length) {
        errno = EINVAL;
        return (char *)strerror_errors_complete[length];
    }
    return (char *)strerror_errors_complete[errnum];
}
