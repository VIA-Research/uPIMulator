/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef _DPUSYSCORE_ERRNO_H_
#define _DPUSYSCORE_ERRNO_H_

#include <defs.h>

/**
 * @file errno.h
 * @brief Defines the system error numbers.
 */

// errno is an array indexed on the tasklet id rather than
// a single integer.
extern int __errno[];

// Mimic errno variable as an index to __errno.
// Defined in such a way that users can't override errno.
#define errno (*(__errno + me()))

/**
 * @def E2BIG
 * @brief  Argument list too long.
 */
#define E2BIG 1
/**
 * @def EACCES
 * @brief  Permission denied.
 */
#define EACCES 2
/**
 * @def EADDRINUSE
 * @brief  Address in use.
 */
#define EADDRINUSE 3
/**
 * @def EADDRNOTAVAIL
 * @brief  Address not available.
 */
#define EADDRNOTAVAIL 4
/**
 * @def EAFNOSUPPORT
 * @brief  Address family not supported.
 */
#define EAFNOSUPPORT 5
/**
 * @def EAGAIN
 * @brief  Resource unavailable, try again.
 */
#define EAGAIN 6
/**
 * @def EALREADY
 * @brief  Connection already in progress.
 */
#define EALREADY 7
/**
 * @def EBADF
 * @brief  Bad file descriptor.
 */
#define EBADF 8
/**
 * @def EBADMSG
 * @brief  Bad message.
 */
#define EBADMSG 9
/**
 * @def EBUSY
 * @brief  Device or resource busy.
 */
#define EBUSY 10
/**
 * @def ECANCELED
 * @brief  Operation canceled.
 */
#define ECANCELED 11
/**
 * @def ECHILD
 * @brief  No child processes.
 */
#define ECHILD 12
/**
 * @def ECONNABORTED
 * @brief  Connection aborted.
 */
#define ECONNABORTED 13
/**
 * @def ECONNREFUSED
 * @brief  Connection refused.
 */
#define ECONNREFUSED 14
/**
 * @def ECONNRESET
 * @brief  Connection reset.
 */
#define ECONNRESET 15
/**
 * @def EDEADLK
 * @brief  Resource deadlock would occur.
 */
#define EDEADLK 16
/**
 * @def EDESTADDRREQ
 * @brief  Destination address required.
 */
#define EDESTADDRREQ 17
/**
 * @def EDOM
 * @brief  Mathematics argument out of domain of function.
 */
#define EDOM 18
/**
 * @def EDQUOT
 * @brief  Reserved.
 */
#define EDQUOT 19
/**
 * @def EEXIST
 * @brief  File exists.
 */
#define EEXIST 20
/**
 * @def EFAULT
 * @brief  Bad address.
 */
#define EFAULT 21
/**
 * @def EFBIG
 * @brief  File too large.
 */
#define EFBIG 22
/**
 * @def EHOSTUNREACH
 * @brief  Host is unreachable.
 */
#define EHOSTUNREACH 23
/**
 * @def EIDRM
 * @brief  Identifier removed.
 */
#define EIDRM 24
/**
 * @def EILSEQ
 * @brief  Illegal byte sequence.
 */
#define EILSEQ 25
/**
 * @def EINPROGRESS
 * @brief  Operation in progress.
 */
#define EINPROGRESS 26
/**
 * @def EINTR
 * @brief  Interrupted function.
 */
#define EINTR 27
/**
 * @def EINVAL
 * @brief  Invalid argument.
 */
#define EINVAL 28
/**
 * @def EIO
 * @brief  I/O error.
 */
#define EIO 29
/**
 * @def EISCONN
 * @brief  Socket is connected.
 */
#define EISCONN 30
/**
 * @def EISDIR
 * @brief  Is a directory.
 */
#define EISDIR 31
/**
 * @def ELOOP
 * @brief  Too many levels of symbolic links.
 */
#define ELOOP 32
/**
 * @def EMFILE
 * @brief  File descriptor value too large.
 */
#define EMFILE 33
/**
 * @def EMLINK
 * @brief  Too many links.
 */
#define EMLINK 34
/**
 * @def EMSGSIZE
 * @brief  Message too large.
 */
#define EMSGSIZE 35
/**
 * @def EMULTIHOP
 * @brief  Reserved.
 */
#define EMULTIHOP 36
/**
 * @def ENAMETOOLONG
 * @brief  Filename too long.
 */
#define ENAMETOOLONG 37
/**
 * @def ENETDOWN
 * @brief  Network is down.
 */
#define ENETDOWN 38
/**
 * @def ENETRESET
 * @brief  Connection aborted by network.
 */
#define ENETRESET 39
/**
 * @def ENETUNREACH
 * @brief  Network unreachable.
 */
#define ENETUNREACH 40
/**
 * @def ENFILE
 * @brief  Too many files open in system.
 */
#define ENFILE 41
/**
 * @def ENOBUFS
 * @brief  No buffer space available.
 */
#define ENOBUFS 42
/**
 * @def ENODATA
 * @brief  No message is available on the STREAM head read queue.
 */
#define ENODATA 43
/**
 * @def ENODEV
 * @brief  No such device.
 */
#define ENODEV 44
/**
 * @def ENOENT
 * @brief  No such file or directory.
 */
#define ENOENT 45
/**
 * @def ENOEXEC
 * @brief  Executable file format error.
 */
#define ENOEXEC 46
/**
 * @def ENOLCK
 * @brief  No locks available.
 */
#define ENOLCK 47
/**
 * @def ENOLINK
 * @brief  Reserved.
 */
#define ENOLINK 48
/**
 * @def ENOMEM
 * @brief  Not enough space.
 */
#define ENOMEM 49
/**
 * @def ENOMSG
 * @brief  No message of the desired type.
 */
#define ENOMSG 50
/**
 * @def ENOPROTOOPT
 * @brief  Protocol not available.
 */
#define ENOPROTOOPT 51
/**
 * @def ENOSPC
 * @brief  No space left on device.
 */
#define ENOSPC 52
/**
 * @def ENOSR
 * @brief  No STREAM resources.
 */
#define ENOSR 53
/**
 * @def ENOSTR
 * @brief  Not a STREAM.
 */
#define ENOSTR 54
/**
 * @def ENOSYS
 * @brief  Function not supported.
 */
#define ENOSYS 55
/**
 * @def ENOTCONN
 * @brief  The socket is not connected.
 */
#define ENOTCONN 56
/**
 * @def ENOTDIR
 * @brief  Not a directory or a symbolic link to a directory.
 */
#define ENOTDIR 57
/**
 * @def ENOTEMPTY
 * @brief  Directory not empty.
 */
#define ENOTEMPTY 58
/**
 * @def ENOTRECOVERABLE
 * @brief  State not recoverable.
 */
#define ENOTRECOVERABLE 59
/**
 * @def ENOTSOCK
 * @brief  Not a socket.
 */
#define ENOTSOCK 60
/**
 * @def ENOTSUP
 * @brief  Not supported.
 */
#define ENOTSUP 61
/**
 * @def ENOTTY
 * @brief  Inappropriate I/O control operation.
 */
#define ENOTTY 62
/**
 * @def ENXIO
 * @brief  No such device or address.
 */
#define ENXIO 63
/**
 * @def EOPNOTSUPP
 * @brief  Operation not supported on socket.
 */
#define EOPNOTSUPP ENOTSUP
/**
 * @def EOVERFLOW
 * @brief  Value too large to be stored in data type.
 */
#define EOVERFLOW 65
/**
 * @def EOWNERDEAD
 * @brief  Previous owner died.
 */
#define EOWNERDEAD 66
/**
 * @def EPERM
 * @brief  Operation not permitted.
 */
#define EPERM 67
/**
 * @def EPIPE
 * @brief  Broken pipe.
 */
#define EPIPE 68
/**
 * @def EPROTO
 * @brief  Protocol error.
 */
#define EPROTO 69
/**
 * @def EPROTONOSUPPORT
 * @brief  Protocol not supported.
 */
#define EPROTONOSUPPORT 70
/**
 * @def EPROTOTYPE
 * @brief  Protocol wrong type for socket.
 */
#define EPROTOTYPE 71
/**
 * @def ERANGE
 * @brief  Result too large.
 */
#define ERANGE 72
/**
 * @def EROFS
 * @brief  Read-only file system.
 */
#define EROFS 73
/**
 * @def ESPIPE
 * @brief  Invalid seek.
 */
#define ESPIPE 74
/**
 * @def ESRCH
 * @brief  No such process.
 */
#define ESRCH 75
/**
 * @def ESTALE
 * @brief  Reserved.
 */
#define ESTALE 76
/**
 * @def ETIME
 * @brief  Stream ioctl() timeout.
 */
#define ETIME 77
/**
 * @def ETIMEDOUT
 * @brief  Connection timed out.
 */
#define ETIMEDOUT 78
/**
 * @def ETXTBSY
 * @brief  Text file busy.
 */
#define ETXTBSY 79
/**
 * @def EWOULDBLOCK
 * @brief  Operation would block.
 */
#define EWOULDBLOCK ENOTSUP
/**
 * @def EXDEV
 * @brief  Cross-device link.
 */
#define EXDEV 81

#endif /* _DPUSYSCORE_ERRNO_H_ */
