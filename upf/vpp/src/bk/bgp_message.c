/*
 * bgp_message.c - BGP protocol message handling
 *
 * Copyright (c) <current-year> <your-organization>
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at:
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#include <vlib/vlib.h>
#include <bgp/bgp.h>
#include <vppinfra/clib.h>

/**
 * BGP Message Types
 */
typedef enum {
    BGP_MSG_OPEN = 1,
    BGP_MSG_UPDATE = 2,
    BGP_MSG_NOTIFICATION = 3,
    BGP_MSG_KEEPALIVE = 4
} bgp_message_type_t;

/**
 * Parse BGP Header
 */
static int bgp_parse_header(const u8 *data, size_t len, bgp_message_type_t *type, u16 *length) {
    if (len < 19) {
        clib_warning("Invalid BGP header length: %zu", len);
        return -1;
    }

    // Validate marker
    for (int i = 0; i < 16; i++) {
        if (data[i] != 0xFF) {
            clib_warning("Invalid BGP marker");
            return -1;
        }
    }

    // Extract length
    *length = (data[16] << 8) | data[17];
    if (*length < 19 || *length > 4096) {
        clib_warning("Invalid BGP message length: %u", *length);
        return -1;
    }

    // Extract type
    *type = data[18];
    return 0;
}

/**
 * Handle BGP OPEN Message
 */
static void bgp_handle_open_message(bgp_main_t *bmp, const u8 *data, size_t len) {
    if (len < 29) {
        clib_warning("Invalid BGP OPEN message length: %zu", len);
        return;
    }

    u8 version = data[0];
    u16 as_number = (data[1] << 8) | data[2];
    u16 hold_time = (data[3] << 8) | data[4];
    u32 bgp_id = (data[5] << 24) | (data[6] << 16) | (data[7] << 8) | data[8];
    u8 opt_param_len = data[9];

    clib_warning("Received BGP OPEN: Version=%u, AS=%u, HoldTime=%u, BGP ID=%u, OptLen=%u",
                 version, as_number, hold_time, bgp_id, opt_param_len);

    // TODO: Process optional parameters
}

/**
 * Handle BGP KEEPALIVE Message
 */
static void bgp_handle_keepalive_message(bgp_main_t *bmp) {
    clib_warning("Received BGP KEEPALIVE");
    // TODO: Reset session hold timer
}

/**
 * Handle BGP UPDATE Message
 */
static void bgp_handle_update_message(bgp_main_t *bmp, const u8 *data, size_t len) {
    clib_warning("Received BGP UPDATE: Length=%zu", len);
    // TODO: Parse and process UPDATE message
}

/**
 * Handle BGP NOTIFICATION Message
 */
static void bgp_handle_notification_message(bgp_main_t *bmp, const u8 *data, size_t len) {
    clib_warning("Received BGP NOTIFICATION: Length=%zu", len);
    // TODO: Handle NOTIFICATION (log and possibly reset session)
}

/**
 * Process BGP Message
 */
void bgp_process_message(void *message, size_t length) {
    bgp_main_t *bmp = &bgp_main;
    bgp_message_type_t type;
    u16 msg_length;

    if (bgp_parse_header(message, length, &type, &msg_length) < 0) {
        clib_warning("Failed to parse BGP message");
        return;
    }

    const u8 *payload = (const u8 *)message + 19;
    size_t payload_len = msg_length - 19;

    switch (type) {
    case BGP_MSG_OPEN:
        bgp_handle_open_message(bmp, payload, payload_len);
        break;

    case BGP_MSG_KEEPALIVE:
        bgp_handle_keepalive_message(bmp);
        break;

    case BGP_MSG_UPDATE:
        bgp_handle_update_message(bmp, payload, payload_len);
        break;

    case BGP_MSG_NOTIFICATION:
        bgp_handle_notification_message(bmp, payload, payload_len);
        break;

    default:
        clib_warning("Unknown BGP message type: %u", type);
        break;
    }
}

/**
 * Generate BGP KEEPALIVE Message
 */
void bgp_send_keepalive_message(bgp_main_t *bmp) {
    u8 keepalive[19] = {0};

    // Set marker
    memset(keepalive, 0xFF, 16);

    // Set length
    keepalive[16] = 0;
    keepalive[17] = 19;

    // Set type
    keepalive[18] = BGP_MSG_KEEPALIVE;

    // TODO: Send message to peer
    clib_warning("Sending BGP KEEPALIVE");
}


