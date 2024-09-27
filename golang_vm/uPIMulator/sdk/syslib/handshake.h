/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef DPUSYSCORE_HANDSHAKE_H
#define DPUSYSCORE_HANDSHAKE_H

/**
 * @file handshake.h
 * @brief Synchronization with handshakes.
 *
 * This synchronization mechanism allows to synchronize 2 tasklets. One tasklet will serve as a notifier
 * and will call handshake_notify() and the other as a customer and will call handshake_wait_for(notifier).
 *
 * @internal    If the notifier called handshake_notify() before the customer, it will stop until some tasklet
 *              calls handshake_wait_for(notifier).
 *              If a tasklet called handshake_wait_for(notifier) before the notifier, it will stop until
 *              the notifier calls handshake_notify(). If afterwards (still before the notifier calls
 *              handshake_notify()) another tasklet attempts to call handshake_wait_for(notifier) with
 *              the same tasklet in the parameter, the function will do nothing and will return the number
 *              of error and set the errno to the corresponding error number.
 */

#include <sysdef.h>

/**
 * @fn handshake_notify
 * @brief Notifies a tasklet waiting for the notifier.
 *
 * The invoking tasklet is suspended until another tasklet calls handshake_wait_for(notifier).
 * When this condition is reached, the function resumes the waiting tasklet.
 */
void
handshake_notify(void);

/**
 * @fn handshake_wait_for
 * @brief Waits for the notifier tasklet
 *
 * The invoking tasklet is suspended until the notifier tasklet (indicated in the parameter) sends a
 * notification to tell the invoking tasklet that it can go ahead.
 *
 * Beware that if the notifier tasklet and the invoking tasklet are the same, the tasklet will be suspended with no
 * easy way to wake it up. The user should check this case itself if it is something that their program allows.
 *
 * If the number of the notifier is not a defined tasklet, the function behavior is undefined. If some other tasklet has
 * already called handshake_wait_for() with the same notifier in the parameter and that the notifier has not yet called
 * handshake_notify(), the function will do nothing and simply return EALREADY.
 *
 * In both cases the errno will be set to the corresponding error number.
 *
 * @param notifier a number to wait the notification from. It must be a defined tasklet.
 * @return 0 if no error was detected, EALREADY if a corresponding error was detected.
 */
int
handshake_wait_for(sysname_t notifier);

#endif /* DPUSYSCORE_HANDSHAKE_H */
