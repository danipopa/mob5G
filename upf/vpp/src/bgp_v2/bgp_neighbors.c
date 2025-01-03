#include <vnet/vnet.h>
#include <vnet/ip/ip.h>
#include <vnet/plugin/plugin.h>
#include <vlibmemory/api.h>
#include <bgp/bgp.h>
#include <vpp/app/version.h>

#define REPLY_MSG_ID_BASE bmp->msg_id_base
#include <vlibapi/api_helper_macros.h>


/**
 * Initialize a BGP neighbor and its resources.
 */

/* Enable or disable BGP on a specific interface */
int bgp_enable_disable(bgp_main_t *bmp, u32 sw_if_index, int enable_disable) {
    clib_warning("%sabling BGP on interface index %u",
                 enable_disable ? "En" : "Dis", sw_if_index);

    // Ensure the interface index is valid
    vnet_sw_interface_t *sw_interface = vnet_get_sw_interface(bmp->vnet_main, sw_if_index);
    if (!sw_interface) {
        clib_warning("Invalid interface index %u", sw_if_index);
        return -1;
    }

    // Iterate over all neighbors and enable/disable BGP sessions
    bgp_neighbor_t *neighbor;
    pool_foreach(neighbor, bmp->neighbors) {
        if (enable_disable) {
            // Enable BGP session
            bgp_start_session(bmp, neighbor);
        } else {
            // Disable BGP session
            bgp_stop_session(bmp, neighbor->neighbor_ip);
        }
    }

    return 0;
}

/* Initialize a BGP neighbor */
void bgp_neighbor_init(bgp_main_t *bmp, bgp_neighbor_t *neighbor, ip4_address_t neighbor_ip, u32 remote_as, u32 local_as, u32 bgp_identifier) {
    neighbor->neighbor_ip = neighbor_ip;
    neighbor->remote_as = remote_as;
    neighbor->local_as = local_as;
    neighbor->bgp_identifier = bgp_identifier;
    neighbor->socket = bgp_socket_init(&neighbor_ip);  // Initialize the socket
    clib_spinlock_init(&neighbor->lock);

    if (!neighbor->socket) {
        clib_warning("Failed to initialize socket for neighbor %U", format_ip4_address, &neighbor_ip);
    } else {
        clib_warning("Initialized socket %d for neighbor %U", neighbor->socket->socket_fd, format_ip4_address, &neighbor_ip);
        if (bgp_socket_connect(neighbor->socket) != 0) {
            clib_warning("Failed to connect socket %d for neighbor %U", neighbor->socket->socket_fd, format_ip4_address, &neighbor_ip);
        } else {
            clib_warning("Connected socket %d for neighbor %U", neighbor->socket->socket_fd, format_ip4_address, &neighbor_ip);
        }
    }

    // Transition to Idle state initially
    bgp_enter_state(bmp, neighbor, BGP_STATE_IDLE);
}

/* Start a BGP session */
void bgp_start_session(bgp_main_t *bmp, bgp_neighbor_t *neighbor) {
    // Transition to Connect state
    bgp_enter_state(bmp, neighbor, BGP_STATE_CONNECT);
}



/* Stop a BGP session with a neighbor */
void bgp_stop_session(bgp_main_t *bmp, ip4_address_t neighbor_ip) {
    bgp_neighbor_t *neighbor = bgp_find_neighbor(bmp, neighbor_ip);
    if (!neighbor) {
        clib_warning("Neighbor %U not found, cannot stop session.", format_ip4_address, &neighbor_ip);
        return;
    }

    // Lock the neighbor structure to ensure thread safety
    clib_spinlock_lock(&neighbor->lock);

    // Transition to Idle state
    bgp_enter_state(bmp, neighbor, BGP_STATE_IDLE);

    // Clear session-specific resources
    bgp_clear_session_resources(neighbor);

    // Unlock the neighbor structure
    clib_spinlock_unlock(&neighbor->lock);

    clib_warning("Stopped BGP session with neighbor: %U.", format_ip4_address, &neighbor_ip);
}

/* Add a new BGP neighbor */
void bgp_add_neighbor(bgp_main_t *bmp, ip4_address_t neighbor_ip, u32 remote_as) {
    clib_spinlock_lock(&bmp->lock);

    bgp_neighbor_t *neighbor;
    pool_get_zero(bmp->neighbors, neighbor);

    // Initialize the neighbor
    bgp_neighbor_init(bmp, neighbor, neighbor_ip, remote_as, bmp->bgp_as_number, bmp->bgp_router_id);

    if (!neighbor->socket) {
        clib_warning("Failed to initialize socket for neighbor %U", format_ip4_address, &neighbor_ip);
        pool_put(bmp->neighbors, neighbor);  // Remove neighbor from pool if socket initialization fails
        clib_spinlock_unlock(&bmp->lock);
        return;
    }

    clib_warning("Added BGP neighbor: %U (AS %u)", format_ip4_address, &neighbor_ip, remote_as);

    // Start the BGP session
    bgp_start_session(bmp, neighbor);

    clib_spinlock_unlock(&bmp->lock);
}

/* Remove a BGP neighbor */
void bgp_remove_neighbor(bgp_main_t *bmp, ip4_address_t neighbor_ip) {
    bgp_neighbor_t *neighbor = bgp_find_neighbor(bmp, neighbor_ip);
    if (!neighbor) {
        clib_warning("Neighbor %U not found, cannot remove.", format_ip4_address, &neighbor_ip);
        return;
    }

    // Lock the neighbor structure to ensure thread safety
    clib_spinlock_lock(&neighbor->lock);

    // Transition to Idle state
    bgp_enter_state(bmp, neighbor, BGP_STATE_IDLE);

    // Clear session-specific resources
    bgp_clear_session_resources(neighbor);

    // Remove the neighbor from the pool
    pool_put(bmp->neighbors, neighbor);

    // Unlock the neighbor structure
    clib_spinlock_unlock(&neighbor->lock);

    clib_warning("Removed BGP neighbor: %U.", format_ip4_address, &neighbor_ip);
}

// /* Find a BGP neighbor by its IP */
// bgp_neighbor_t *bgp_find_neighbor(bgp_main_t *bmp, ip4_address_t neighbor_ip) {
//     bgp_neighbor_t *neighbor;
//     pool_foreach(neighbor, bmp->neighbors) {
//         if (ip4_address_cmp(&neighbor->neighbor_ip, &neighbor_ip) == 0) {
//             return neighbor;
//         }
//     };
//     return NULL;
// }

// /* Reset a BGP neighbor session */
// void bgp_reset_neighbor(bgp_main_t *bmp, ip4_address_t neighbor_ip) {
//     bgp_neighbor_t *neighbor = bgp_find_neighbor(bmp, neighbor_ip);

//     if (!neighbor) {
//         clib_warning("Neighbor %U not found, cannot reset session.", format_ip4_address, &neighbor_ip);
//         return;
//     }

//     clib_warning("Resetting BGP neighbor: %U...", format_ip4_address, &neighbor_ip);

//     // Transition to Idle state
//     neighbor->state = BGP_STATE_IDLE;

//     // Clear queued messages for this neighbor
//     vec_free(neighbor->output_queue);

//     // Clear RIB-in and RIB-out entries for this neighbor
//     bgp_clear_rib_in_for_neighbor(bmp, neighbor_ip);
//     bgp_clear_rib_out_for_neighbor(bmp, neighbor_ip);

//     // Restart the session by transitioning to Connect state
//     neighbor->state = BGP_STATE_CONNECT;

//     // Notify the periodic process to handle re-connection
//     vlib_process_signal_event(bmp->vlib_main, bmp->periodic_node_index, BGP_EVENT1, (uword)neighbor);

//     clib_warning("BGP neighbor %U reset complete.", format_ip4_address, &neighbor_ip);
// }


// void bgp_hard_reset_neighbor(bgp_main_t *bmp, ip4_address_t neighbor_ip) {
//     bgp_neighbor_t *neighbor = bgp_find_neighbor(bmp, neighbor_ip);

//     if (!neighbor) {
//         clib_warning("Neighbor %U not found.", format_ip4_address, &neighbor_ip);
//         return;
//     }

//     clib_warning("Performing hard reset for BGP neighbor %U...", format_ip4_address, &neighbor_ip);

//     // Transition neighbor to Idle state
//     neighbor->state = BGP_STATE_IDLE;

//     // Clear RIB-in and RIB-out for the neighbor
//     bgp_clear_rib_in_for_neighbor(bmp, neighbor_ip);
//     bgp_clear_rib_out_for_neighbor(bmp, neighbor_ip);

//     // Clear session-specific resources
//     bgp_clear_session_resources(neighbor);

//     // Remove the neighbor
//     bgp_remove_neighbor(bmp, neighbor);

//     // Reinitialize the neighbor
//     u32 remote_as = neighbor->remote_as; // Save the remote AS number
//     bgp_add_neighbor(bmp, neighbor_ip, remote_as);

//     // Transition to the Connect state to restart the session
//     neighbor = bgp_find_neighbor(bmp, neighbor_ip);
//     if (neighbor) {
//         neighbor->state = BGP_STATE_CONNECT;

//         // Notify periodic process to handle reconnection
//         vlib_process_signal_event(bmp->vlib_main, bmp->periodic_node_index, BGP_EVENT1, (uword) neighbor);
//         clib_warning("Hard reset for BGP neighbor %U completed.", format_ip4_address, &neighbor_ip);
//     } else {
//         clib_warning("Failed to reinitialize neighbor %U after hard reset.", format_ip4_address, &neighbor_ip);
//     }
// }

// void bgp_soft_reset_neighbor(bgp_main_t *bmp, ip4_address_t neighbor_ip, bool inbound) {
//     bgp_neighbor_t *neighbor = bgp_find_neighbor(bmp, neighbor_ip);

//     if (!neighbor) {
//         clib_warning("Neighbor %U not found.", format_ip4_address, &neighbor_ip);
//         return;
//     }

//     clib_warning("Performing soft reset for BGP neighbor %U (%s)...",
//                  format_ip4_address, &neighbor_ip, inbound ? "inbound" : "outbound");

//     if (inbound) {
//         // Reapply inbound route policies and refresh RIB-in
//         bgp_clear_rib_in_for_neighbor(bmp, neighbor_ip);
//         bgp_request_full_update(neighbor, /*rib_in=*/true);
//     } else {
//         // Reapply outbound route policies and refresh RIB-out
//         bgp_clear_rib_out_for_neighbor(bmp, neighbor_ip);
//         bgp_recompute_rib_out(neighbor);
//     }

//     clib_warning("Soft reset for BGP neighbor %U (%s) completed.",
//                  format_ip4_address, &neighbor_ip, inbound ? "inbound" : "outbound");
// }

void bgp_clear_session_resources(bgp_neighbor_t *neighbor) {
    if (neighbor->socket) {
        bgp_socket_close(neighbor->socket->socket_fd); // Close the neighbor's socket
        neighbor->socket = NULL;
    }
    queue_free(&neighbor->output_queue); // Free the neighbor's message queue
    clib_warning("Cleared session resources for neighbor %U", format_ip4_address, &neighbor->neighbor_ip);
    
    // Reset timers
    neighbor->hold_timer = 0;
    neighbor->keepalive_timer = 0;

    // Additional resource cleanup as needed
}

// void bgp_clear_rib_in_for_neighbor(bgp_main_t *bmp, ip4_address_t neighbor_ip) {
//     // Placeholder for clearing inbound RIB
//     clib_warning("Cleared inbound RIB for neighbor %U", format_ip4_address, &neighbor_ip);
// }

// void bgp_clear_rib_out_for_neighbor(bgp_main_t *bmp, ip4_address_t neighbor_ip) {
//     // Placeholder for clearing outbound RIB
//     clib_warning("Cleared outbound RIB for neighbor %U", format_ip4_address, &neighbor_ip);
// }

// /* Send BGP OPEN message */
// int bgp_send_open_message(bgp_neighbor_t *neighbor) {
//     size_t msg_length;
//     void *msg = bgp_create_open_message(&msg_length, neighbor->local_as, neighbor->bgp_identifier);

//     if (!msg) {
//         clib_warning("Failed to create OPEN message for neighbor %U", format_ip4_address, &neighbor->neighbor_ip);
//         return -1;
//     }

//     // Send the message over the TCP connection
//     ssize_t sent = send(neighbor->socket, msg, msg_length, 0);
//     if (sent != msg_length) {
//         clib_warning("Failed to send OPEN message to neighbor %U", format_ip4_address, &neighbor->neighbor_ip);
//         clib_mem_free(msg);
//         bgp_enter_state(neighbor, BGP_STATE_IDLE);  // Transition back to IDLE on failure
//         return -1;
//     }

//     clib_mem_free(msg);
//     bgp_enter_state(neighbor, BGP_STATE_OPEN_SENT);  // Transition to OPEN_SENT on success
//     return 0;
// }

// /* Stop route exchange with a BGP neighbor */
// void bgp_stop_route_exchange(bgp_neighbor_t *neighbor) {
//     clib_warning("Stopping route exchange with neighbor %U (AS %u)",
//                  format_ip4_address, &neighbor->neighbor_ip, neighbor->remote_as);

//     // Lock the neighbor structure to ensure thread safety
//     clib_spinlock_lock(&neighbor->lock);

//     // Transition to Idle state
//     bgp_enter_state(neighbor, BGP_STATE_IDLE);

//     // Clear session-specific resources
//     bgp_clear_session_resources(neighbor);

//     // Unlock the neighbor structure
//     clib_spinlock_unlock(&neighbor->lock);

//     clib_warning("Stopped route exchange with neighbor %U (AS %u)",
//                  format_ip4_address, &neighbor->neighbor_ip, neighbor->remote_as);
// }


// bool bgp_tcp_is_connected(bgp_neighbor_t *neighbor) {
//     clib_warning("Checking TCP connection for neighbor %U", format_ip4_address, &neighbor->neighbor_ip);
//     return true; // Replace with actual connection check logic
// }

// bool bgp_received_open(bgp_neighbor_t *neighbor) {
//     clib_warning("Checking if OPEN message received from neighbor %U", format_ip4_address, &neighbor->neighbor_ip);
//     return false; // Replace with actual logic
// }

// bool bgp_received_notification(bgp_neighbor_t *neighbor) {
//     clib_warning("Checking if NOTIFICATION message received from neighbor %U", format_ip4_address, &neighbor->neighbor_ip);
//     return false; // Replace with actual logic
// }

// bool bgp_received_keepalive(bgp_neighbor_t *neighbor) {
//     clib_warning("Checking if KEEPALIVE message received from neighbor %U", format_ip4_address, &neighbor->neighbor_ip);
//     return false; // Replace with actual logic
// }

// void bgp_send_keepalive_message(bgp_neighbor_t *neighbor) {
//     clib_warning("Sending KEEPALIVE message to neighbor %U", format_ip4_address, &neighbor->neighbor_ip);
//     // TODO: Add logic to construct and send a KEEPALIVE message
// }