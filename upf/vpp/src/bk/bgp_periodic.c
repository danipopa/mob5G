#include <vlib/vlib.h>
#include <bgp/bgp.h>

static uword bgp_periodic_process(vlib_main_t *vm, vlib_node_runtime_t *rt, vlib_frame_t *f) {
    bgp_main_t *bmp = &bgp_main;
    while (1) {
        if (bmp->periodic_timer_enabled) {
            // Perform periodic tasks such as keep-alives
        }
        vlib_process_wait_for_event_or_clock(vm, 1.0);
    }
    return 0;
}

void bgp_create_periodic_process(bgp_main_t *bmp) {
    if (bmp->periodic_node_index > 0) return;
    bmp->periodic_node_index = vlib_process_create(bmp->vlib_main, "bgp-periodic-process", bgp_periodic_process, 16);
}

