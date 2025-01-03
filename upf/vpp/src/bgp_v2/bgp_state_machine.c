#include <bgp/bgp.h>
#include <vlib/vlib.h>
#include <bgp/bgp.h>

void bgp_handle_route_update(bgp_main_t *bmp, bgp_neighbor_t *neighbor) {
    u8 *route_update_data_pointer = NULL;
    u16 route_update_data_length = 0;

    // Construct the route update data
    if (bgp_construct_route_update(bmp, neighbor, &route_update_data_pointer, &route_update_data_length) == 0) {
        // Allocate memory for the BGP message
        bgp_message_t *message = clib_mem_alloc(sizeof(bgp_message_t));
        if (!message) {
            clib_warning("Failed to allocate memory for BGP message");
            clib_mem_free(route_update_data_pointer); // Free route update data if allocation fails
            return;
        }

        // Populate the BGP message
        message->type = BGP_MSG_UPDATE;
        message->data = route_update_data_pointer;
        message->length = route_update_data_length;

        // Enqueue the message for transmission
        if (bgp_enqueue_message(neighbor, message) < 0) {
            clib_warning("Failed to enqueue route update message for neighbor %U",
                         format_ip4_address, &neighbor->neighbor_ip);
            clib_mem_free(route_update_data_pointer); // Free the message data
            clib_mem_free(message);                   // Free the message structure
        } else {
            clib_warning("Route update message enqueued for neighbor %U",
                         format_ip4_address, &neighbor->neighbor_ip);
        }
    } else {
        clib_warning("Failed to construct route update data");
    }
}


/**
 * Handles sending and receiving Keepalive messages for a BGP neighbor.
 */
void bgp_handle_keepalive(bgp_neighbor_t *neighbor) {
    bgp_main_t *bmp = &bgp_main;

    // Check if the keepalive timer has expired
    if (neighbor->keepalive_timer > 0) {
        neighbor->keepalive_timer--;
    } else {
        // Send a Keepalive message
        bgp_message_t *keepalive_message = clib_mem_alloc(sizeof(bgp_message_t));
        keepalive_message->type = BGP_MSG_KEEPALIVE;
        keepalive_message->data = NULL;  // No payload for Keepalive messages
        keepalive_message->length = sizeof(bgp_message_header_t); // Only the header

        if (bgp_enqueue_message(neighbor, keepalive_message) < 0) {
            clib_mem_free(keepalive_message); // Free memory if enqueue fails
            clib_warning("Queue is full for neighbor %U, dropping Keepalive message.",
                         format_ip4_address, &neighbor->neighbor_ip);
        } else {
            clib_warning("Sent Keepalive to neighbor %U.",
                         format_ip4_address, &neighbor->neighbor_ip);
        }

        // Reset the keepalive timer
        neighbor->keepalive_timer = bmp->keepalive_time;
    }

    // Process received Keepalive messages (if any)
    while (!queue_is_empty(&neighbor->output_queue)) {

        bgp_message_t *received_message = queue_dequeue(&neighbor->output_queue);
        if (received_message && received_message->type == BGP_MSG_KEEPALIVE) {
            clib_warning("Received Keepalive from neighbor %U.",
                         format_ip4_address, &neighbor->neighbor_ip);
            clib_mem_free(received_message); // Free the processed message
        }
    }
}


/* Handle entering a specific state */
void bgp_enter_state(bgp_main_t *bmp, bgp_neighbor_t *neighbor, bgp_state_t new_state) {
    clib_warning("Neighbor %U transitioning to state: %s",
                 format_ip4_address, &neighbor->neighbor_ip, bgp_state_to_string(new_state));

    // Underlined change: Use a while loop for iterative state transitions
    while (neighbor->state != new_state) {
        neighbor->state = new_state;

        switch (new_state) {
            case BGP_STATE_IDLE:
                // Handle transition to IDLE state
                bgp_clear_session_resources(neighbor);
                return;

            case BGP_STATE_CONNECT:
                // Handle transition to CONNECT state
                if (bgp_send_open_message(bmp, neighbor) != 0) {
                    clib_warning("Failed to send OPEN message to neighbor %U", format_ip4_address, &neighbor->neighbor_ip);
                    // Underlined change: Transition back to IDLE on failure
                    new_state = BGP_STATE_IDLE;
                } else {
                    // Underlined change: Transition to OPEN_SENT on success
                    clib_warning("BGP_STATE_OPEN_SENT sent to neighbor %U", format_ip4_address, &neighbor->neighbor_ip);
                    new_state = BGP_STATE_OPEN_SENT;
                }
                break;

            case BGP_STATE_OPEN_SENT:
                // Handle transition to OPEN_SENT state
                // Wait for OPEN message from peer
                // Underlined change: Placeholder transition to OPEN_CONFIRM for demonstration
                new_state = BGP_STATE_OPEN_CONFIRM;
                break;

            case BGP_STATE_OPEN_CONFIRM:
                // Handle transition to OPEN_CONFIRM state
                // Wait for KEEPALIVE message from peer
                // Underlined change: Placeholder transition to ESTABLISHED for demonstration
                new_state = BGP_STATE_ESTABLISHED;
                break;

            case BGP_STATE_ESTABLISHED:
                // Handle transition to ESTABLISHED state
                // Start exchanging UPDATE messages
                return;

            default:
                clib_warning("Unknown BGP state: %d", new_state);
                return;
        }
    }
}

/* Send BGP OPEN message */
int bgp_send_open_message(bgp_main_t *bmp, bgp_neighbor_t *neighbor) {
    if (!neighbor->socket) {
        clib_warning("Neighbor %U has no valid socket", format_ip4_address, &neighbor->neighbor_ip);
        return -1;
    }

    size_t msg_length;
    void *msg = bgp_create_open_message(&msg_length, neighbor->local_as, neighbor->bgp_identifier);

    if (!msg) {
        clib_warning("Failed to create OPEN message for neighbor %U", format_ip4_address, &neighbor->neighbor_ip);
        return -1;
    }

    // Send the message over the TCP connection using bgp_socket_send
    if (bgp_socket_send(neighbor->socket->socket_fd, msg, msg_length) != 0) {
        clib_warning("Failed to send OPEN message to neighbor %U", format_ip4_address, &neighbor->neighbor_ip);
        clib_mem_free(msg);
        return -1;
    }

    clib_mem_free(msg);
    return 0;
}

/* Handle exiting a specific state */
void bgp_exit_state(bgp_main_t *bmp, bgp_neighbor_t *neighbor, bgp_state_t old_state) {
    clib_warning("Neighbor %U exiting state: %d",
                 format_ip4_address, &neighbor->neighbor_ip, old_state);

    switch (old_state) {
        case BGP_STATE_IDLE:
            // No cleanup needed
            break;

        case BGP_STATE_ESTABLISHED:
            // Clear any routing-related state
            bgp_stop_route_exchange(bmp, neighbor);
            break;

        default:
            // Handle other states if necessary
            break;
    }
}

/* Transition neighbor to a new state */
void bgp_transition_state(bgp_main_t *bmp, bgp_neighbor_t *neighbor, bgp_state_t new_state) {
    if (neighbor->state != new_state) {
        bgp_exit_state(bmp, neighbor, neighbor->state);
        bgp_enter_state(bmp, neighbor, new_state);
    }
}

/* Periodic processing of state transitions */
void bgp_process_state(bgp_main_t *bmp, bgp_neighbor_t *neighbor) {
    switch (neighbor->state) {
        case BGP_STATE_IDLE:
            bgp_transition_state(bmp, neighbor, BGP_STATE_CONNECT);
            break;

        case BGP_STATE_CONNECT:
            if (bgp_tcp_is_connected(neighbor)) {
                bgp_transition_state(bmp, neighbor, BGP_STATE_OPEN_SENT);
            } else {
                bgp_transition_state(bmp, neighbor, BGP_STATE_IDLE);
            }
            break;

        case BGP_STATE_OPEN_SENT:
            if (bgp_received_open(neighbor)) {
                bgp_transition_state(bmp, neighbor, BGP_STATE_OPEN_CONFIRM);
            } else if (bgp_received_notification(neighbor)) {
                bgp_transition_state(bmp, neighbor, BGP_STATE_IDLE);
            }
            break;

        case BGP_STATE_OPEN_CONFIRM:
            if (bgp_received_keepalive(neighbor)) {
                bgp_transition_state(bmp, neighbor, BGP_STATE_ESTABLISHED);
            } else if (bgp_received_notification(neighbor)) {
                bgp_transition_state(bmp, neighbor, BGP_STATE_IDLE);
            }
            break;

        case BGP_STATE_ESTABLISHED:
            // Handle keepalive and route updates
            bgp_handle_keepalive(neighbor);
            bgp_handle_route_update(bmp, neighbor);

            // Process queued messages for the neighbor
            while (!queue_is_empty(&neighbor->output_queue)) {
                bgp_message_t *message = queue_dequeue(&neighbor->output_queue);
                if (message) {
                    // Send the message to the neighbor
                    if (bgp_socket_send(neighbor->socket->socket_fd , message->data, message->length) < 0) {
                        clib_warning("Failed to send message to neighbor %U.",
                                    format_ip4_address, &neighbor->neighbor_ip);
                    }

                    // Free the message after sending
                    clib_mem_free(message);
                }
            }

            break;

        default:
            clib_warning("Neighbor %U in unknown state: %d",
                         format_ip4_address, &neighbor->neighbor_ip, neighbor->state);
            break;
    }
}

/* Timer events for managing neighbor state */
void bgp_handle_timers(bgp_main_t *bmp, bgp_neighbor_t *neighbor) {
    // Handle hold timer expiration
    if (neighbor->hold_timer > 0 && --neighbor->hold_timer == 0) {
        clib_warning("Hold timer expired for neighbor %U. Resetting session.",
                     format_ip4_address, &neighbor->neighbor_ip);
        bgp_transition_state(bmp, neighbor, BGP_STATE_IDLE);
    }

    // Handle keepalive timer expiration
    if (neighbor->keepalive_timer > 0 && --neighbor->keepalive_timer == 0) {
        clib_warning("Keepalive timer expired for neighbor %U. Sending keepalive.",
                     format_ip4_address, &neighbor->neighbor_ip);
        bgp_send_keepalive_message(neighbor);
        neighbor->keepalive_timer = bmp->keepalive_time;
    }
}

