#include <bgp/bgp.h>
#include <vnet/vnet.h>
// #include <bgp/bgp_messages.h>

// Compare two IPv4 addresses
int ip4_address_cmp(const ip4_address_t *a, const ip4_address_t *b) {
    return memcmp(a, b, sizeof(ip4_address_t));
}

// Print a warning and exit on memory allocation failure
void *safe_mem_alloc(size_t size) {
    void *ptr = clib_mem_alloc(size);
    if (!ptr) {
        clib_warning("Memory allocation failed for size: %zu", size);
        abort();
    }
    return ptr;
}

int unformat_fib_prefix(unformat_input_t *input, fib_prefix_t *prefix) {
    ip4_address_t ip4;
    ip6_address_t ip6;
    u8 mask_length;

    if (unformat(input, "%U/%d", unformat_ip4_address, &ip4, &mask_length)) {
        prefix->fp_proto = FIB_PROTOCOL_IP4;
        prefix->fp_len = mask_length;
        clib_memcpy(&prefix->fp_addr.ip4, &ip4, sizeof(ip4));
        return 1;
    } else if (unformat(input, "%U/%d", unformat_ip6_address, &ip6, &mask_length)) {
        prefix->fp_proto = FIB_PROTOCOL_IP6;
        prefix->fp_len = mask_length;
        clib_memcpy(&prefix->fp_addr.ip6, &ip6, sizeof(ip6));
        return 1;
    }

    return 0;  // Parsing failed
}

const char *bgp_state_to_string(bgp_state_t state) {
    switch (state) {
        case BGP_STATE_IDLE:
            return "Idle";
        case BGP_STATE_CONNECT:
            return "Connect";
        case BGP_STATE_ACTIVE:
            return "Active";
        case BGP_STATE_OPEN_SENT:
            return "OpenSent";
        case BGP_STATE_OPEN_CONFIRM:
            return "OpenConfirm";
        case BGP_STATE_ESTABLISHED:
            return "Established";
        default:
            return "Unknown";
    }
}

int bgp_enqueue_message(bgp_neighbor_t *neighbor, bgp_message_t *message) {
    if (queue_is_full(&neighbor->output_queue)) {
        clib_warning("Output queue is full for neighbor %U", format_ip4_address, &neighbor->neighbor_ip);
        return -1; // Queue is full
    }
    queue_enqueue(&neighbor->output_queue, message);
    return 0; // Success
}

bgp_message_t *bgp_dequeue_message(bgp_neighbor_t *neighbor) {
    if (queue_is_empty(&neighbor->output_queue)) {
        clib_warning("Queue is empty for neighbor %U.", format_ip4_address, &neighbor->neighbor_ip);
        return NULL;
    }
    return queue_dequeue(&neighbor->output_queue);

}


void queue_init(custom_queue_t *queue, int capacity) {
    queue->buffer = clib_mem_alloc(capacity * sizeof(bgp_message_t *));
    queue->capacity = capacity;
    queue->head = 0;
    queue->tail = 0;
    queue->count = 0;
}


int queue_enqueue(custom_queue_t *queue, bgp_message_t *message) {
    if (queue_is_full(queue)) {
        return -1; // Queue is full
    }
    queue->buffer[queue->tail] = message;
    queue->tail = (queue->tail + 1) % queue->capacity;
    queue->count++;
    return 0; // Success
}

bgp_message_t *queue_dequeue(custom_queue_t *queue) {
    if (queue->count == 0) {
        return NULL; // Queue is empty
    }
    bgp_message_t *message = queue->buffer[queue->head];
    queue->head = (queue->head + 1) % queue->capacity;
    queue->count--;
    return message;
}

void queue_free(custom_queue_t *queue) {
    clib_mem_free(queue->buffer);
    queue->buffer = NULL;
    queue->capacity = 0;
    queue->count = 0;
    queue->head = 0;
    queue->tail = 0;
}

bool queue_is_empty(custom_queue_t *queue) {
    return queue->count == 0;
}

bool queue_is_full(custom_queue_t *queue) {
    return queue->count == queue->capacity;
}

/* Create a BGP OPEN message */
void *bgp_create_open_message(size_t *out_length, u32 local_as, u32 bgp_identifier) {
    bgp_open_message_t *msg = clib_mem_alloc(sizeof(bgp_open_message_t));
    memset(msg, 0, sizeof(*msg));

    // Set the marker to all ones
    memset(msg->marker, 0xFF, sizeof(msg->marker));

    // Set the length of the message
    msg->length = clib_host_to_net_u16(sizeof(bgp_open_message_t));

    // Set the message type to OPEN
    msg->type = BGP_MSG_OPEN;

    // Set the BGP version
    msg->version = 4;  // BGP-4

    // Set the local AS number
    msg->my_as = clib_host_to_net_u16(local_as);

    // Set the hold time (default to 180 seconds)
    msg->hold_time = clib_host_to_net_u16(180);

    // Set the BGP identifier
    msg->bgp_identifier = clib_host_to_net_u32(bgp_identifier);

    // Set the optional parameters length to 0 (no optional parameters)
    msg->opt_param_len = 0;

    // Set the output length
    *out_length = sizeof(bgp_open_message_t);

    return msg;
}


/**
 * Constructs a BGP route update message.
 */
int bgp_construct_route_update(bgp_main_t *bmp, bgp_neighbor_t *neighbor, u8 **data, u16 *length) {
    // Example implementation: This should encode withdrawn routes, path attributes, and NLRI.
    // For now, we'll mock a simple route update message.
    
    const char *mock_data = "Mock Route Update Data"; // Replace with actual data encoding
    size_t mock_length = strlen(mock_data) + 1;

    *data = clib_mem_alloc(mock_length);
    if (!*data) {
        clib_warning("Failed to allocate memory for route update data");
        return -1;
    }

    memcpy(*data, mock_data, mock_length);
    *length = (u16)mock_length;

    return 0; // Success
}




// Request a full RIB update for the neighbor
void bgp_request_full_update(bgp_neighbor_t *neighbor, bool rib_in) {
    clib_warning("Requested full RIB %s update for neighbor %U",
                 rib_in ? "IN" : "OUT",
                 format_ip4_address, &neighbor->neighbor_ip);

    // TODO: Add logic to trigger a full update
    // For RIB IN: Update inbound routes and attributes from the neighbor
    // For RIB OUT: Recompute advertised routes
}

// Recompute the outbound RIB for the neighbor
void bgp_recompute_rib_out(bgp_neighbor_t *neighbor) {
    clib_warning("Recomputing RIB OUT for neighbor %U",
                 format_ip4_address, &neighbor->neighbor_ip);

    // TODO: Add logic to recompute RIB OUT
    // This might involve updating route advertisements to the neighbor
}