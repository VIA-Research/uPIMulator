/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

////
//////TODO : **endptr
////
////#include <stddef.h>
////#include "ctype.h"
////#include "limits.h"
////
////static long overflow(int sign)
////{
////    if(sign)
////        return INT64_MAX;
////    else
////        return INT64_MIN;
////
////
////}
////
////long strtol(const char *nptr, char **endptr, int base)
////{
////    if(nptr == NULL)
////        return 0;
////    unsigned long result = 0;
////    int sign = 1;
////    int i = 0;
//////    while(isspace(nptr[i]))
//////        i++;
////
//////    switch(nptr[i]){
//////        case '-' :
//////            sign = 0;       //we change sign only when '-' was encountered
//////        case '+' :
//////            i++;            //we increment i in both cases
//////    }
////
//////
//////    if((base == 16 ) && (nptr[i] == '0') && (nptr[i+1] == 'x'))
//////        i+=2;
//////    else if(base == 0){
//////        if(nptr[i] == '0'){
//////            base = 8;
//////            i++;
//////            if(nptr[i+1] == 'x'){
//////                base = 16;
//////                i++;
//////            }
//////        }
//////        else
//////            base = 10;
//////    }
////
//////    switch(base) {
//////        case 10 :
////            while(((nptr[i]>>4) == 0x3) && ((nptr[i] & 0xf)<=9)){
////                if((unsigned long)result>>60){
////                    return overflow(sign);
////                }
////
////                result = (unsigned long)(result<<1) + (unsigned long)(result<<3) + (unsigned long) nptr[i] - (unsigned long)
///'0';
//////                if((unsigned long)result>>63){
//////                    return overflow(sign);
//////                }
////                i++;
////            }
//////            break;
//////        case 2 :
//////            while((nptr[i] & 0xfe) == 0x30){
//////                result = (result<<1) + (unsigned long)(nptr[i] & 0x1);
//////                i++;
//////            }
//////            break;
//////        case 4 :
//////            while((nptr[i] & 0xfc) == 0x30){
//////                result = (result<<2) + (unsigned long)(nptr[i] & 0x3);
//////                i++;
//////            }
//////            break;
//////        case 8 :
//////            while((nptr[i] & 0xf8) == 0x30){
//////                result = (result<<3) + (unsigned long)(nptr[i] & 0x7);
//////                i++;
//////            }
//////            break;
//////        case 16 :
//////            while(((nptr[i]>='0') && (nptr[i]<='9')) || ((nptr[i]>='a') && (nptr[i]<='f')) || ((nptr[i]>='A') &&
///(nptr[i]<='F')) ){
//////                unsigned long digit = ((nptr[i]>='0') && (nptr[i]<='9')) ? (unsigned long)(nptr[i] - '0') : (unsigned
/// long)((nptr[i] & 0xf) + 10);
//////                result = result<<4;
//////                result = result + digit;
//////                i++;
//////            }
//////            break;
//////        default :
//////            if(base < 2 || base > 36){ //base 0 has already been replaced by 8,10 or 16
//////                //TODO!!!
//////                break;
//////            }
//////    }
////
//////    if(sign == 0)
//////        return (long)0 - (long)result;
//////    else
////        return (long) result;
////
////}
////
//
//
//
////
//
//// Copyright (c) 2014-2019 - UPMEM
//
////
//
//#include <stddef.h>
//
//#include "limits.h"
//
//#include "ctype.h"
//
// static long overflow(int sign)
//
//{
//
//    if(sign)
//
//        return INT64_MAX;
//
//    else
//
//        return INT64_MIN;
//
//}
//
// long int strtol(const char *nptr, char **endptr, int base)
//
//{
//
//    if(nptr == NULL)
//
//        return 0;
//
//    int sign = 1;
//
//    int i=0;
//
//    long result = 0;
//
//
//
//    while((nptr[i]>='0') && (nptr[i]<='9')){
//
//        if(result>>60){
//
//                    return overflow(sign);
//
//                }
//
//        result = (result<<1) + (result<<3)+ nptr[i] - '0';
//
//        i++;
//
//    }
//
//    return result;
//
//}
