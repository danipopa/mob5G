#ifndef __included_bgp_h__
#define __included_bgp_h__

#include <vnet/vnet.h>
#include <vnet/ip/ip.h>
#include <vnet/ethernet/ethernet.h>
#include <vppinfra/hash.h>
#include <vppinfra/error.h>


typedef struct {
    ip4_address_t prefix;   // The aggregated prefix
    u8 prefix_length;       // The prefix length
    u8 summary_only;        // Flag to indicate summary-only behavior
    u8 as_set;              // Flag to indicate AS_SET attribute
} bgp_aggregate_t;

typedef struct {
    ip4_address_t neighbor_ip;  // Neighbor IP address
    u32 remote_as;              // Remote AS number
    u8 is_route_reflector_client; // Flag for route reflector client
    char route_filter_name[64]; // Associated route filter

} bgp_neighbor_t;

typedef struct {
    ip4_address_t prefix;   // The IP prefix of the route
    u8 mask_length;         // The subnet mask length
    ip4_address_t next_hop; // The next-hop IP address for the route
    u32 local_pref;         // Local preference (optional, for BGP decision process)
    u32 as_path_length;     // AS path length (optional)
    char as_path[256];      // AS path string (optional, for debugging)
    u8 origin;              // Origin attribute (e.g., IGP, EGP, incomplete)
    u32 med;               // Multi-Exit Discriminator (optional)
} bgp_route_t;


typedef struct {
    ip4_address_t prefix;
    u8 mask_length;
    bool permit;
} bgp_prefix_t;

typedef struct {
    char name[64];
    bgp_prefix_t **entries; // Array of prefixes
} bgp_prefix_list_t;


typedef struct {
    u16 msg_id_base;
    u8 periodic_timer_enabled;
    u32 periodic_node_index;

    vlib_main_t *vlib_main;
    vnet_main_t *vnet_main;
    ethernet_main_t *ethernet_main;

    // BGP-specific state
    u32 bgp_router_id;
    u32 bgp_as_number;
    u32 hold_time; 
    u32 keepalive_time;
    u32 cluster_id;

    clib_spinlock_t lock;         // Spinlock for thread safety
    bgp_neighbor_t *neighbors;    // Pool of BGP neighbors

    // Add other structures for peer/session management
    bgp_prefix_list_t **prefix_lists; // Array of prefix lists
    // Pool of BGP routes
    bgp_route_t *routes; 
} bgp_main_t;


extern bgp_main_t bgp_main;

extern vlib_node_registration_t bgp_node;
extern vlib_node_registration_t bgp_periodic_node;

// Periodic function events
#define BGP_EVENT1 1
#define BGP_EVENT2 2
#define BGP_EVENT_PERIODIC_ENABLE_DISABLE 3

void bgp_create_periodic_process(bgp_main_t *);
void bgp_process_message(void *message, size_t length);

// Function declarations
bgp_neighbor_t *bgp_find_neighbor(bgp_main_t *bmp, ip4_address_t neighbor_ip);
bgp_neighbor_t *bgp_add_neighbor(bgp_main_t *bmp, ip4_address_t neighbor_ip, u32 remote_as);
void bgp_remove_neighbor(bgp_main_t *bmp, bgp_neighbor_t *neighbor);

void bgp_add_aggregate(bgp_main_t *bmp, ip4_address_t prefix, u8 prefix_length, u8 summary_only, u8 as_set);
void bgp_remove_aggregate(bgp_main_t *bmp, ip4_address_t prefix, u8 prefix_length);


#endif /* __included_bgp_h__ */

