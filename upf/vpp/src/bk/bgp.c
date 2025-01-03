#include <vnet/vnet.h>
#include <vnet/ip/ip.h>
#include <vnet/plugin/plugin.h>
#include <bgp/bgp.h>

#include <vlibapi/api.h>
#include <vlibmemory/api.h>
#include <vpp/app/version.h>
#include <stdbool.h>

#include <bgp/bgp.api_enum.h>
#include <bgp/bgp.api_types.h>

#define REPLY_MSG_ID_BASE bmp->msg_id_base
#include <vlibapi/api_helper_macros.h>

bgp_main_t bgp_main;

int bgp_enable_disable(bgp_main_t *bmp, u32 sw_if_index, int enable_disable) {
    // Add logic for enabling/disabling BGP on interfaces.
    return 0;
}

static clib_error_t *bgp_init(vlib_main_t *vm) {
    bgp_main_t *bmp = &bgp_main;
    bmp->vlib_main = vm;
    bmp->vnet_main = vnet_get_main();
    bmp->bgp_router_id = 0;  // Initialize router ID
    bmp->bgp_as_number = 0;  // Initialize AS number
    // Initialize the routes pool
    bmp->routes = 0;
    bmp->prefix_lists = NULL; // Initialize to NULL to indicate an empty list
    return 0;
}

VLIB_INIT_FUNCTION(bgp_init);


static inline int ip4_address_cmp(const ip4_address_t *a, const ip4_address_t *b) {
    return memcmp(a, b, sizeof(ip4_address_t));
}

bgp_neighbor_t *bgp_find_neighbor(bgp_main_t *bmp, ip4_address_t neighbor_ip) {
    bgp_neighbor_t *neighbor;
    pool_foreach(neighbor, bmp->neighbors) {
        if (!ip4_address_cmp(&neighbor->neighbor_ip, &neighbor_ip)) {
            return neighbor;
        }
    }
    return NULL;
}

void bgp_set_timers(bgp_main_t *bmp, u32 hold_time, u32 keepalive_time) {
    bmp->hold_time = hold_time;
    bmp->keepalive_time = keepalive_time;

    clib_warning("BGP timers set: Hold Time = %d, Keepalive Time = %d", hold_time, keepalive_time);
    // Update timers for all active BGP sessions if applicable
    // Iterate through neighbors and apply changes.
}


void bgp_set_cluster_id(bgp_main_t *bmp, u32 cluster_id) {
    bmp->cluster_id = cluster_id;

    clib_warning("BGP cluster ID set to %d", cluster_id);
    // Apply cluster ID in route reflector behavior if applicable.
}

void bgp_set_route_reflector_client(bgp_main_t *bmp, ip4_address_t neighbor_ip) {
    bgp_neighbor_t *neighbor = bgp_find_neighbor(bmp, neighbor_ip);

    if (!neighbor) {
        clib_warning("Neighbor %U not found.", format_ip4_address, &neighbor_ip);
        return;
    }

    neighbor->is_route_reflector_client = 1;
    clib_warning("Neighbor %U marked as a route reflector client.", format_ip4_address, &neighbor_ip);
}

bgp_prefix_list_t *bgp_find_or_create_prefix_list(bgp_main_t *bmp, const char *list_name) {
    bgp_prefix_list_t *list;
    vec_foreach(list, bmp->prefix_lists) {
        if (strcmp((char *)list->name, list_name) == 0) {
            return list;
        }
    }

    list = clib_mem_alloc(sizeof(bgp_prefix_list_t));
    memset(list, 0, sizeof(*list));
    strncpy(list->name, list_name, sizeof(list->name) - 1);
    vec_add1(bmp->prefix_lists, list);
    return list;
}

void bgp_set_aggregate_address(bgp_main_t *bmp, ip4_address_t *prefix, u8 mask_length, int summary_only, int as_set) {
    bgp_aggregate_t *aggregate;

    aggregate = bgp_find_or_create_aggregate(bmp, prefix, mask_length);
    aggregate->summary_only = summary_only;
    aggregate->as_set = as_set;

    clib_warning("Aggregate address %U/%d configured (summary-only=%d, as-set=%d)",
                 format_ip4_address, prefix, mask_length, summary_only, as_set);
}

void bgp_show_routes(bgp_main_t *bmp) {
    bgp_route_t *route;
    vlib_cli_output(vm, "BGP Routes:");
    pool_foreach(route, bmp->routes) {
        vlib_cli_output(vm, "Route: %U/%d -> Next Hop: %U",
                        format_ip4_address, &route->prefix, route->mask_length,
                        format_ip4_address, &route->next_hop);
    }
}

void bgp_show_neighbor(bgp_main_t *bmp, ip4_address_t neighbor_ip) {
    bgp_neighbor_t *neighbor = bgp_find_neighbor(bmp, neighbor_ip);

    if (!neighbor) {
        clib_warning("Neighbor %U not found.", format_ip4_address, &neighbor_ip);
        return;
    }

    clib_warning("Neighbor %U: State=%s, AS=%d, Hold Time=%d, Keepalive Time=%d",
                 format_ip4_address, &neighbor_ip,
                 bgp_state_to_string(neighbor->state),
                 neighbor->remote_as,
                 neighbor->hold_time, neighbor->keepalive_time);
}

void bgp_show_statistics(bgp_main_t *bmp) {
    clib_warning("BGP Statistics:");
    clib_warning("  Messages Sent: Open=%d, Update=%d, Keepalive=%d, Notification=%d",
                 bmp->stats.msg_sent_open,
                 bmp->stats.msg_sent_update,
                 bmp->stats.msg_sent_keepalive,
                 bmp->stats.msg_sent_notification);

    clib_warning("  Messages Received: Open=%d, Update=%d, Keepalive=%d, Notification=%d",
                 bmp->stats.msg_received_open,
                 bmp->stats.msg_received_update,
                 bmp->stats.msg_received_keepalive,
                 bmp->stats.msg_received_notification);
}

void bgp_show_rib_in(bgp_main_t *bmp, ip4_address_t neighbor_ip) {
    bgp_rib_entry_t *entry;

    clib_warning("BGP RIB In for Neighbor %U:", format_ip4_address, &neighbor_ip);
    pool_foreach(entry, bmp->rib_in) {
        if (!ip4_address_cmp(&entry->neighbor_ip, &neighbor_ip)) {
            clib_warning("  Route: %U/%d -> Next Hop: %U",
                         format_ip4_address, &entry->prefix, entry->mask_length,
                         format_ip4_address, &entry->next_hop);
        }
    }
}


void bgp_show_rib_out(bgp_main_t *bmp, ip4_address_t neighbor_ip) {
    bgp_rib_entry_t *entry;

    clib_warning("BGP RIB Out for Neighbor %U:", format_ip4_address, &neighbor_ip);
    pool_foreach(entry, bmp->rib_out) {
        if (!ip4_address_cmp(&entry->neighbor_ip, &neighbor_ip)) {
            clib_warning("  Route: %U/%d -> Next Hop: %U",
                         format_ip4_address, &entry->prefix, entry->mask_length,
                         format_ip4_address, &entry->next_hop);
        }
    }
}


bgp_neighbor_t *bgp_find_neighbor(bgp_main_t *bmp, ip4_address_t neighbor_ip) {
    bgp_neighbor_t *neighbor;

    pool_foreach(neighbor, bmp->neighbors) {
        if (!ip4_address_cmp(&neighbor->address, &neighbor_ip)) {
            return neighbor;
        }
    }

    return NULL;
}


bgp_aggregate_t *bgp_find_or_create_aggregate(bgp_main_t *bmp, ip4_address_t prefix, u8 prefix_length) {
    bgp_aggregate_t *aggregate;
    pool_foreach(aggregate, bmp->aggregates) {
        if (ip4_address_cmp(&aggregate->prefix, &prefix) == 0 &&
            aggregate->prefix_length == prefix_length) {
            return aggregate;
        }
    }
    pool_get(bmp->aggregates, aggregate);
    aggregate->prefix = prefix;
    aggregate->prefix_length = prefix_length;
    aggregate->summary_only = 0;
    aggregate->as_set = 0;
    return aggregate;
}



void bgp_reset_neighbor(bgp_main_t *bmp, ip4_address_t neighbor_ip) {
    bgp_neighbor_t *neighbor = bgp_find_neighbor(bmp, neighbor_ip);

    if (!neighbor) {
        clib_warning("Neighbor %U not found.", format_ip4_address, &neighbor_ip);
        return;
    }

    clib_warning("Resetting BGP neighbor %U...", format_ip4_address, &neighbor_ip);

    // Transition neighbor to Idle state
    neighbor->state = BGP_STATE_IDLE;

    // Clear any queued messages for this neighbor
    vec_free(neighbor->output_queue);

    // Clear RIB-in and RIB-out for the neighbor
    bgp_clear_rib_in_for_neighbor(bmp, neighbor_ip);
    bgp_clear_rib_out_for_neighbor(bmp, neighbor_ip);

    // Restart the session by transitioning to the Connect state
    neighbor->state = BGP_STATE_CONNECT;

    // Send an event to the periodic process to attempt re-connection
    vlib_process_signal_event(bmp->vlib_main, bmp->periodic_node_index, BGP_EVENT1, (uword) neighbor);

    clib_warning("BGP neighbor %U reset complete.", format_ip4_address, &neighbor_ip);
}

void bgp_hard_reset_neighbor(bgp_main_t *bmp, ip4_address_t neighbor_ip) {
    bgp_neighbor_t *neighbor = bgp_find_neighbor(bmp, neighbor_ip);

    if (!neighbor) {
        clib_warning("Neighbor %U not found.", format_ip4_address, &neighbor_ip);
        return;
    }

    clib_warning("Performing hard reset for BGP neighbor %U...", format_ip4_address, &neighbor_ip);

    // Transition neighbor to Idle state
    neighbor->state = BGP_STATE_IDLE;

    // Clear all RIB-in and RIB-out entries
    bgp_clear_rib_in_for_neighbor(bmp, neighbor_ip);
    bgp_clear_rib_out_for_neighbor(bmp, neighbor_ip);

    // Clear any queued messages for this neighbor
    vec_free(neighbor->output_queue);

    // Clear session-specific resources (e.g., timers, state machines)
    bgp_clear_session_resources(neighbor);

    // Remove and reinitialize the neighbor structure
    bgp_remove_neighbor(bmp, neighbor_ip);
    bgp_add_neighbor(bmp, neighbor_ip);

    // Restart the session by transitioning to the Connect state
    neighbor->state = BGP_STATE_CONNECT;

    // Notify periodic process to handle re-connection
    vlib_process_signal_event(bmp->vlib_main, bmp->periodic_node_index, BGP_EVENT1, (uword) neighbor);

    clib_warning("Hard reset for BGP neighbor %U completed.", format_ip4_address, &neighbor_ip);
}


void bgp_clear_session_resources(bgp_neighbor_t *neighbor) {
    clib_warning("Clearing session resources for neighbor %U...", format_ip4_address, &neighbor->ip);

    // Free dynamically allocated resources
    vec_free(neighbor->input_queue);
    vec_free(neighbor->output_queue);

    // Reset timers
    neighbor->hold_timer = 0;
    neighbor->keepalive_timer = 0;
}


void bgp_remove_neighbor(bgp_main_t *bmp, ip4_address_t neighbor_ip) {
    bgp_neighbor_t *neighbor;

    pool_foreach(neighbor, bmp->neighbors) {
        if (!ip4_address_cmp(&neighbor->ip, &neighbor_ip)) {
            pool_put(bmp->neighbors, neighbor);
            clib_warning("Removed neighbor %U.", format_ip4_address, &neighbor_ip);
            return;
        }
    }

    clib_warning("Neighbor %U not found during removal.", format_ip4_address, &neighbor_ip);
}

void bgp_add_neighbor(bgp_main_t *bmp, ip4_address_t neighbor_ip) {

    clib_spinlock_lock(&bmp->lock);
    bgp_neighbor_t *neighbor;
    pool_get_zero(bmp->neighbors, neighbor);
    neighbor->neighbor_ip = neighbor_ip;
    neighbor->state = BGP_STATE_IDLE;
    neighbor->remote_as = remote_as;
    neighbor->hold_timer = 0;
    neighbor->keepalive_timer = 0;
    clib_warning("Added neighbor %U.", format_ip4_address, &neighbor_ip);
    clib_spinlock_unlock(&bmp->lock);
}


void bgp_soft_reset_neighbor(bgp_main_t *bmp, ip4_address_t neighbor_ip, bool inbound) {
    bgp_neighbor_t *neighbor = bgp_find_neighbor(bmp, neighbor_ip);

    if (!neighbor) {
        clib_warning("Neighbor %U not found.", format_ip4_address, &neighbor_ip);
        return;
    }

    clib_warning("Performing soft reset for BGP neighbor %U (%s)...",
                 format_ip4_address, &neighbor_ip, inbound ? "inbound" : "outbound");

    if (inbound) {
        // Reapply inbound route policies and refresh RIB-in
        bgp_clear_rib_in_for_neighbor(bmp, neighbor_ip);
        bgp_request_full_update(neighbor, /*rib_in=*/true);
    } else {
        // Reapply outbound route policies and refresh RIB-out
        bgp_clear_rib_out_for_neighbor(bmp, neighbor_ip);
        bgp_recompute_rib_out(neighbor);
    }

    clib_warning("Soft reset for BGP neighbor %U (%s) completed.",
                 format_ip4_address, &neighbor_ip, inbound ? "inbound" : "outbound");
}

void bgp_request_full_update(bgp_neighbor_t *neighbor, bool rib_in) {
    if (rib_in) {
        clib_warning("Requesting full BGP UPDATE from neighbor %U", format_ip4_address, &neighbor->ip);

        // Send a BGP ROUTE REFRESH message to the neighbor
        bgp_send_route_refresh(neighbor);

    } else {
        clib_warning("Outbound RIB refresh for neighbor %U", format_ip4_address, &neighbor->ip);

        // Trigger recomputation of outbound routes for this neighbor
        bgp_recompute_rib_out(neighbor);
    }
}

void bgp_send_route_refresh(bgp_neighbor_t *neighbor) {
    // Construct BGP Route Refresh message
    u8 route_refresh[23] = {0};

    // Set marker
    memset(route_refresh, 0xFF, 16);

    // Set length
    route_refresh[16] = 0;
    route_refresh[17] = 23;

    // Set type (Route Refresh)
    route_refresh[18] = 5;

    // Set AFI (IPv4) and SAFI (Unicast)
    route_refresh[19] = 0;  // AFI High Byte
    route_refresh[20] = 1;  // AFI Low Byte
    route_refresh[21] = 1;  // SAFI

    // Send the message
    bgp_send_message(neighbor, route_refresh, sizeof(route_refresh));

    clib_warning("Sent BGP Route Refresh to neighbor %U", format_ip4_address, &neighbor->ip);
}

void bgp_update_prefix_list(bgp_main_t *bmp, const char *list_name, ip4_address_t *prefix, u8 mask_length, bool permit) {
    bgp_prefix_list_t *list = bgp_find_or_create_prefix_list(bmp, list_name);

    if (!list) {
        clib_warning("Failed to find or create prefix list: %s", list_name);
        return;
    }

    bgp_prefix_t *entry = clib_mem_alloc(sizeof(bgp_prefix_t));
    entry->prefix = *prefix;
    entry->mask_length = mask_length;
    entry->permit = permit;

    vec_add1(list->entries, entry);

    clib_warning("Prefix %U/%u %s added to list %s.",
                 format_ip4_address, prefix, mask_length,
                 permit ? "permit" : "deny", list_name);
}


bgp_prefix_list_t *bgp_find_or_create_prefix_list(bgp_main_t *bmp, const char *list_name) {
    bgp_prefix_list_t *list;

    // Check if the list exists
    vec_foreach(list, bmp->prefix_lists) {
        if (strcmp((char *)list->name, list_name) == 0) {
            return list;
        }
    }

    // Create a new list if it doesn't exist
    list = clib_mem_alloc(sizeof(bgp_prefix_list_t));
    memset(list, 0, sizeof(*list));
    strncpy((char *)list->name, list_name, sizeof(list->name) - 1);
    vec_add1(bmp->prefix_lists, list);

    return list;
}

void bgp_free_prefix_lists(bgp_main_t *bmp) {
    bgp_prefix_list_t *list;
    vec_foreach(list, bmp->prefix_lists) {
        vec_free(list->entries); // Free entries in each list
        clib_mem_free(list);     // Free the list itself
    }
    vec_free(bmp->prefix_lists); // Free the array of prefix lists
}

void bgp_add_aggregate(bgp_main_t *bmp, ip4_address_t prefix, u8 prefix_length, u8 summary_only, u8 as_set) {
    bgp_aggregate_t *aggregate;
    pool_get(bmp->aggregates, aggregate);
    aggregate->prefix = prefix;
    aggregate->prefix_length = prefix_length;
    aggregate->summary_only = summary_only;
    aggregate->as_set = as_set;
    clib_warning("Added aggregate: %U/%d (summary-only: %d, as-set: %d)",
                 format_ip4_address, &prefix, prefix_length, summary_only, as_set);
}

void bgp_remove_aggregate(bgp_main_t *bmp, ip4_address_t prefix, u8 prefix_length) {
    bgp_aggregate_t *aggregate;
    pool_foreach(aggregate, bmp->aggregates, {
        if (ip4_address_cmp(&aggregate->prefix, &prefix) == 0 && aggregate->prefix_length == prefix_length) {
            pool_put(bmp->aggregates, aggregate);
            clib_warning("Removed aggregate: %U/%d", format_ip4_address, &prefix, prefix_length);
            return;
        }
    });
    clib_warning("Aggregate not found: %U/%d", format_ip4_address, &prefix, prefix_length);
}

bgp_aggregate_t *bgp_find_or_create_aggregate(bgp_main_t *bmp, ip4_address_t prefix, u8 prefix_length) {
    bgp_aggregate_t *aggregate;

    // Search for an existing aggregate
    // pool_foreach(aggregate, bmp->aggregates, {
    //     if (ip4_address_cmp(&aggregate->prefix, &prefix) == 0 &&
    //         aggregate->prefix_length == prefix_length) {
    //         clib_warning("Found existing aggregate: %U/%d",
    //                      format_ip4_address, &prefix, prefix_length);
    //         return aggregate;
    //     }
    // });

    // // Create a new aggregate if not found
    // pool_get(bmp->aggregates, aggregate);
    aggregate->prefix = prefix;
    aggregate->prefix_length = prefix_length;
    aggregate->summary_only = 0; // Default value
    aggregate->as_set = 0;       // Default value
    clib_warning("Created new aggregate: %U/%d",
                 format_ip4_address, &prefix, prefix_length);
    return aggregate;
}


/*              My BGP Commands            */

/* BGP Neighbor Clear */
static clib_error_t *
bgp_neighbor_clear_command_fn(vlib_main_t *vm, unformat_input_t *input, vlib_cli_command_t *cmd) {
    ip4_address_t neighbor_ip;
    if (!unformat(input, "%U", unformat_ip4_address, &neighbor_ip)) {
        return clib_error_return(0, "Please specify a valid neighbor IP address.");
    }

    // Logic to reset the neighbor session
    bgp_reset_neighbor(&bgp_main, neighbor_ip);

    clib_warning("BGP session with neighbor %U cleared.", format_ip4_address, &neighbor_ip);
    return 0;
}

VLIB_CLI_COMMAND(bgp_neighbor_clear_command, static) = {
    .path = "bgp neighbor clear",
    .short_help = "bgp neighbor clear <ip-address>",
    .function = bgp_neighbor_clear_command_fn,
};

/* Set BGP Router ID */
static clib_error_t *
bgp_set_router_id_command_fn(vlib_main_t *vm, unformat_input_t *input, vlib_cli_command_t *cmd)
{
    bgp_main_t *bmp = &bgp_main;
    ip4_address_t router_id;

    if (!unformat(input, "%U", unformat_ip4_address, &router_id))
        return clib_error_return(0, "Invalid IPv4 address format");

    bmp->bgp_router_id = router_id.as_u32;
    clib_warning("BGP Router ID set to %U", format_ip4_address, &router_id);
    return 0;
}

VLIB_CLI_COMMAND(bgp_set_router_id_command, static) = {
    .path = "set bgp router-id",
    .short_help = "set bgp router-id <IPv4>",
    .function = bgp_set_router_id_command_fn,
};


/*  Set BGP AS Number */
static clib_error_t *
bgp_set_as_command_fn(vlib_main_t *vm, unformat_input_t *input, vlib_cli_command_t *cmd)
{
    bgp_main_t *bmp = &bgp_main;
    u32 as_number;

    if (!unformat(input, "%d", &as_number))
        return clib_error_return(0, "Invalid AS number");

    bmp->bgp_as_number = as_number;
    clib_warning("BGP AS number set to %u", as_number);
    return 0;
}

VLIB_CLI_COMMAND(bgp_set_as_command, static) = {
    .path = "set bgp as",
    .short_help = "set bgp as <AS-number>",
    .function = bgp_set_as_command_fn,
};


/* Add BGP Neighbor */
static clib_error_t *
bgp_add_neighbor_command_fn(vlib_main_t *vm, unformat_input_t *input, vlib_cli_command_t *cmd)
{
    //bgp_main_t *bmp = &bgp_main;
    ip4_address_t neighbor_ip;
    u32 remote_as;

    if (!unformat(input, "%U remote-as %d", unformat_ip4_address, &neighbor_ip, &remote_as))
        return clib_error_return(0, "Usage: set bgp neighbor <IPv4> remote-as <AS-number>");

    // TODO: Add logic to store neighbor configuration
    clib_warning("Added BGP neighbor %U with remote AS %u", format_ip4_address, &neighbor_ip, remote_as);
    return 0;
}

VLIB_CLI_COMMAND(bgp_add_neighbor_command, static) = {
    .path = "set bgp neighbor",
    .short_help = "set bgp neighbor <IPv4> remote-as <AS-number>",
    .function = bgp_add_neighbor_command_fn,
};


/* Enable BGP on Interface */
static clib_error_t *
bgp_enable_interface_command_fn(vlib_main_t *vm, unformat_input_t *input, vlib_cli_command_t *cmd)
{
    bgp_main_t *bmp = &bgp_main;
    u32 sw_if_index = ~0;

    if (!unformat(input, "%U", unformat_vnet_sw_interface, bmp->vnet_main, &sw_if_index))
        return clib_error_return(0, "Invalid interface name");

    // TODO: Add logic to enable BGP on the interface
    clib_warning("Enabled BGP on interface index %u", sw_if_index);
    return 0;
}

VLIB_CLI_COMMAND(bgp_enable_interface_command, static) = {
    .path = "set bgp enable",
    .short_help = "set bgp enable <interface-name>",
    .function = bgp_enable_interface_command_fn,
};

/* Advertise Network */
static clib_error_t *
bgp_advertise_network_command_fn(vlib_main_t *vm, unformat_input_t *input, vlib_cli_command_t *cmd)
{
    //bgp_main_t *bmp = &bgp_main;
    fib_prefix_t prefix;

    //if (!unformat(input, "%U", unformat_fib_prefix, &prefix))
    //    return clib_error_return(0, "Invalid network prefix format");

    // TODO: Add logic to advertise the network
    clib_warning("Advertised BGP network %U", format_fib_prefix, &prefix);
    return 0;
}

VLIB_CLI_COMMAND(bgp_advertise_network_command, static) = {
    .path = "set bgp advertise network",
    .short_help = "set bgp advertise network <prefix>",
    .function = bgp_advertise_network_command_fn,
};

/* Show BGP Config */
static clib_error_t *
bgp_show_config_command_fn(vlib_main_t *vm, unformat_input_t *input, vlib_cli_command_t *cmd)
{
    bgp_main_t *bmp = &bgp_main;

    vlib_cli_output(vm, "BGP Configuration:");
    vlib_cli_output(vm, "  Router ID: %U", format_ip4_address, &bmp->bgp_router_id);
    vlib_cli_output(vm, "  AS Number: %u", bmp->bgp_as_number);

    // TODO: Display configured neighbors and advertised networks
    // For now, placeholder messages:
    vlib_cli_output(vm, "  Neighbors: <To be implemented>");
    vlib_cli_output(vm, "  Advertised Networks: <To be implemented>");

    return 0;
}

VLIB_CLI_COMMAND(bgp_show_config_command, static) = {
    .path = "show bgp config",
    .short_help = "show bgp config",
    .function = bgp_show_config_command_fn,
};


/* Show BGP Summary */
static clib_error_t *
bgp_show_summary_command_fn(vlib_main_t *vm, unformat_input_t *input, vlib_cli_command_t *cmd)
{
    bgp_main_t *bmp = &bgp_main;

    vlib_cli_output(vm, "BGP Summary:");
    vlib_cli_output(vm, "  Router ID: %U", format_ip4_address, &bmp->bgp_router_id);
    vlib_cli_output(vm, "  AS Number: %u", bmp->bgp_as_number);

    // TODO: Display neighbor status and session details
    // Placeholder messages:
    vlib_cli_output(vm, "  Neighbors:");
    vlib_cli_output(vm, "    Neighbor: <IP Address>, Remote AS: <AS Number>, State: <State>");

    return 0;
}

VLIB_CLI_COMMAND(bgp_show_summary_command, static) = {
    .path = "show bgp summary",
    .short_help = "show bgp summary",
    .function = bgp_show_summary_command_fn,
};

/*BGP Neighbor Reset */
static clib_error_t *
bgp_neighbor_reset_command_fn(vlib_main_t *vm, unformat_input_t *input, vlib_cli_command_t *cmd) {
    ip4_address_t neighbor_ip;
    if (!unformat(input, "%U", unformat_ip4_address, &neighbor_ip)) {
        return clib_error_return(0, "Please specify a valid neighbor IP address.");
    }

    // Logic to perform a hard reset
    bgp_hard_reset_neighbor(&bgp_main, neighbor_ip);

    clib_warning("BGP session with neighbor %U reset.", format_ip4_address, &neighbor_ip);
    return 0;
}

VLIB_CLI_COMMAND(bgp_neighbor_reset_command, static) = {
    .path = "bgp neighbor reset",
    .short_help = "bgp neighbor reset <ip-address>",
    .function = bgp_neighbor_reset_command_fn,
};



/*BGP Neighbor Soft Reset*/
static clib_error_t *
bgp_neighbor_soft_reset_command_fn(vlib_main_t *vm, unformat_input_t *input, vlib_cli_command_t *cmd) {
    ip4_address_t neighbor_ip;
    if (!unformat(input, "%U", unformat_ip4_address, &neighbor_ip)) {
        return clib_error_return(0, "Please specify a valid neighbor IP address.");
    }

    // Logic to perform a soft reset
    bgp_soft_reset_neighbor(&bgp_main, neighbor_ip);

    clib_warning("BGP session with neighbor %U soft-reset.", format_ip4_address, &neighbor_ip);
    return 0;
}

VLIB_CLI_COMMAND(bgp_neighbor_soft_reset_command, static) = {
    .path = "bgp neighbor soft-reset",
    .short_help = "bgp neighbor soft-reset <ip-address>",
    .function = bgp_neighbor_soft_reset_command_fn,
};

/*Set Route Filter */
static clib_error_t *
bgp_set_route_filter_command_fn(vlib_main_t *vm, unformat_input_t *input, vlib_cli_command_t *cmd) {
    ip4_address_t neighbor_ip;
    char *filter_name;

    if (!unformat(input, "%U %s", unformat_ip4_address, &neighbor_ip, &filter_name)) {
        return clib_error_return(0, "Usage: set bgp neighbor <ip-address> route-filter <filter-name>");
    }

    // Apply the route filter
    bgp_set_route_filter(&bgp_main, neighbor_ip, filter_name);

    clib_warning("Applied route filter %s to neighbor %U.", filter_name, format_ip4_address, &neighbor_ip);
    return 0;
}

VLIB_CLI_COMMAND(bgp_set_route_filter_command, static) = {
    .path = "set bgp neighbor route-filter",
    .short_help = "set bgp neighbor <ip-address> route-filter <filter-name>",
    .function = bgp_set_route_filter_command_fn,
};

/* Prefix List Management */
static clib_error_t *
bgp_set_prefix_list_command_fn(vlib_main_t *vm, unformat_input_t *input, vlib_cli_command_t *cmd) {
    char *list_name;
    ip4_address_t prefix;
    u8 mask_length;
    u8 permit = 1;

    if (!unformat(input, "%s permit %U/%d", &list_name, unformat_ip4_address, &prefix, &mask_length)) {
        if (!unformat(input, "%s deny %U/%d", &list_name, unformat_ip4_address, &prefix, &mask_length)) {
            return clib_error_return(0, "Usage: set bgp prefix-list <name> [permit|deny] <prefix>/<mask-length>");
        }
        permit = 0;
    }

    bgp_update_prefix_list(&bgp_main, list_name, &prefix, mask_length, permit);
    clib_warning("Updated prefix list %s: %s %U/%d", list_name, permit ? "permit" : "deny", format_ip4_address, &prefix, mask_length);
    return 0;
}

VLIB_CLI_COMMAND(bgp_set_prefix_list_command, static) = {
    .path = "set bgp prefix-list",
    .short_help = "set bgp prefix-list <name> [permit|deny] <prefix>/<mask-length>",
    .function = bgp_set_prefix_list_command_fn,
};


/* Set Timers */
static clib_error_t *
bgp_set_timers_command_fn(vlib_main_t *vm, unformat_input_t *input, vlib_cli_command_t *cmd) {
    u16 hold_time, keepalive_time;

    if (!unformat(input, "%d %d", &hold_time, &keepalive_time)) {
        return clib_error_return(0, "Usage: set bgp timers <hold-time> <keepalive-time>");
    }

    bgp_set_timers(&bgp_main, hold_time, keepalive_time);
    clib_warning("Set BGP timers: Hold=%d, Keepalive=%d", hold_time, keepalive_time);
    return 0;
}

VLIB_CLI_COMMAND(bgp_set_timers_command, static) = {
    .path = "set bgp timers",
    .short_help = "set bgp timers <hold-time> <keepalive-time>",
    .function = bgp_set_timers_command_fn,
};


/* Set BGP Cluster ID */
static clib_error_t *
bgp_set_cluster_id_command_fn(vlib_main_t *vm, unformat_input_t *input, vlib_cli_command_t *cmd) {
    u32 cluster_id;

    if (!unformat(input, "%d", &cluster_id)) {
        return clib_error_return(0, "Usage: set bgp cluster-id <id>");
    }

    bgp_set_cluster_id(&bgp_main, cluster_id);
    clib_warning("Set BGP cluster ID to %d", cluster_id);
    return 0;
}

VLIB_CLI_COMMAND(bgp_set_cluster_id_command, static) = {
    .path = "set bgp cluster-id",
    .short_help = "set bgp cluster-id <id>",
    .function = bgp_set_cluster_id_command_fn,
};


/* Route Reflector Client */
static clib_error_t *
bgp_set_route_reflector_client_command_fn(vlib_main_t *vm, unformat_input_t *input, vlib_cli_command_t *cmd) {
    ip4_address_t neighbor_ip;

    if (!unformat(input, "%U", unformat_ip4_address, &neighbor_ip)) {
        return clib_error_return(0, "Usage: set bgp neighbor <ip-address> route-reflector-client");
    }

    bgp_set_route_reflector_client(&bgp_main, neighbor_ip);
    clib_warning("Set neighbor %U as route reflector client.", format_ip4_address, &neighbor_ip);
    return 0;
}

VLIB_CLI_COMMAND(bgp_set_route_reflector_client_command, static) = {
    .path = "set bgp neighbor route-reflector-client",
    .short_help = "set bgp neighbor <ip-address> route-reflector-client",
    .function = bgp_set_route_reflector_client_command_fn,
};

/* Aggregate Address  */
static clib_error_t *
bgp_set_aggregate_address_command_fn(vlib_main_t *vm, unformat_input_t *input, vlib_cli_command_t *cmd) {
    ip4_address_t prefix;
    u8 mask_length;
    int summary_only = 0, as_set = 0;

    if (!unformat(input, "%U/%d", unformat_ip4_address, &prefix, &mask_length)) {
        return clib_error_return(0, "Usage: set bgp aggregate-address <prefix>/<mask-length> [summary-only|as-set]");
    }

    while (unformat_check_input(input) != UNFORMAT_END_OF_INPUT) {
        if (unformat(input, "summary-only")) {
            summary_only = 1;
        } else if (unformat(input, "as-set")) {
            as_set = 1;
        } else {
            break;
        }
    }

    bgp_set_aggregate_address(&bgp_main, &prefix, mask_length, summary_only, as_set);
    clib_warning("Configured aggregate address %U/%d (summary-only=%d, as-set=%d)",
                 format_ip4_address, &prefix, mask_length, summary_only, as_set);
    return 0;
}

VLIB_CLI_COMMAND(bgp_set_aggregate_address_command, static) = {
    .path = "set bgp aggregate-address",
    .short_help = "set bgp aggregate-address <prefix>/<mask-length> [summary-only|as-set]",
    .function = bgp_set_aggregate_address_command_fn,
};


/* Show BGP Routes */
static clib_error_t *
bgp_show_routes_command_fn(vlib_main_t *vm, unformat_input_t *input, vlib_cli_command_t *cmd) {
    bgp_main_t *bmp = &bgp_main;

    // Iterate through the route database and print routes
    clib_warning("BGP routes:");
    bgp_show_routes(&bgp_main);

    return 0;
}

VLIB_CLI_COMMAND(bgp_show_routes_command, static) = {
    .path = "show bgp routes",
    .short_help = "show bgp routes",
    .function = bgp_show_routes_command_fn,
};

/* Show BGP Neighbor */
static clib_error_t *
bgp_show_neighbor_command_fn(vlib_main_t *vm, unformat_input_t *input, vlib_cli_command_t *cmd) {
    ip4_address_t neighbor_ip;

    if (!unformat(input, "%U", unformat_ip4_address, &neighbor_ip)) {
        return clib_error_return(0, "Usage: show bgp neighbor <ip-address>");
    }

    bgp_show_neighbor(&bgp_main, neighbor_ip);
    return 0;
}

VLIB_CLI_COMMAND(bgp_show_neighbor_command, static) = {
    .path = "show bgp neighbor",
    .short_help = "show bgp neighbor <ip-address>",
    .function = bgp_show_neighbor_command_fn,
};


/*  Show BGP Statistics */
static clib_error_t *
bgp_show_statistics_command_fn(vlib_main_t *vm, unformat_input_t *input, vlib_cli_command_t *cmd) {
    bgp_show_statistics(&bgp_main);
    return 0;
}

VLIB_CLI_COMMAND(bgp_show_statistics_command, static) = {
    .path = "show bgp statistics",
    .short_help = "show bgp statistics",
    .function = bgp_show_statistics_command_fn,
};

/* Show RIB In */
static clib_error_t *
bgp_show_rib_in_command_fn(vlib_main_t *vm, unformat_input_t *input, vlib_cli_command_t *cmd) {
    ip4_address_t neighbor_ip;

    if (!unformat(input, "%U", unformat_ip4_address, &neighbor_ip)) {
        return clib_error_return(0, "Usage: show bgp rib-in neighbor <ip-address>");
    }

    bgp_show_rib_in(&bgp_main, neighbor_ip);
    return 0;
}

VLIB_CLI_COMMAND(bgp_show_rib_in_command, static) = {
    .path = "show bgp rib-in neighbor",
    .short_help = "show bgp rib-in neighbor <ip-address>",
    .function = bgp_show_rib_in_command_fn,
};

/* Show RIB Out */
static clib_error_t *
bgp_show_rib_out_command_fn(vlib_main_t *vm, unformat_input_t *input, vlib_cli_command_t *cmd) {
    ip4_address_t neighbor_ip;

    if (!unformat(input, "%U", unformat_ip4_address, &neighbor_ip)) {
        return clib_error_return(0, "Usage: show bgp rib-out neighbor <ip-address>");
    }

    bgp_show_rib_out(&bgp_main, neighbor_ip);
    return 0;
}

VLIB_CLI_COMMAND(bgp_show_rib_out_command, static) = {
    .path = "show bgp rib-out neighbor",
    .short_help = "show bgp rib-out neighbor <ip-address>",
    .function = bgp_show_rib_out_command_fn,
};


/* Plugin registration */
VLIB_PLUGIN_REGISTER() = {
    .version = VPP_BUILD_VER,
    .description = "BGP Plugin",
};
