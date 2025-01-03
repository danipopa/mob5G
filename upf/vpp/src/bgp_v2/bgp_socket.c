#include <vnet/vnet.h>
#include <vnet/ip/ip.h>
#include <vppinfra/socket.h>
#include <fcntl.h>
#include <bgp/bgp.h>
#include <arpa/inet.h>  // Include for inet_ntoa

/* Initialize the BGP socket */
bgp_socket_t *bgp_socket_init(ip4_address_t *peer_ip) {
    bgp_socket_t *sock = clib_mem_alloc(sizeof(bgp_socket_t));
    if (!sock) {
        clib_warning("Failed to allocate memory for BGP socket");
        return NULL;
    }
    memset(sock, 0, sizeof(*sock));

    sock->socket_fd = socket(AF_INET, SOCK_STREAM, 0);
    if (sock->socket_fd < 0) {
        clib_warning("Failed to create socket: %s", strerror(errno));
        clib_mem_free(sock);
        return NULL;
    }

    sock->peer_addr.sin_family = AF_INET;
    sock->peer_addr.sin_port = htons(BGP_PORT);
    sock->peer_addr.sin_addr.s_addr = peer_ip->as_u32;

    // Log the attributes of the sock structure
    clib_warning("Initialized BGP socket: socket_fd=%d, peer_addr=%s:%d",
                 sock->socket_fd,
                 inet_ntoa(*(struct in_addr *)&sock->peer_addr.sin_addr),
                 ntohs(sock->peer_addr.sin_port));

    return sock;
}

/* Establish a connection to the BGP peer */
int bgp_socket_connect(bgp_socket_t *sock) {
    int flags = fcntl(sock->socket_fd, F_GETFL, 0);
    if (flags < 0) {
        clib_warning("Failed to get socket flags: %s", strerror(errno));
        return -1;
    }

    if (fcntl(sock->socket_fd, F_SETFL, flags | O_NONBLOCK) < 0) {
        clib_warning("Failed to set socket to non-blocking mode: %s", strerror(errno));
        return -1;
    }

    if (connect(sock->socket_fd, (struct sockaddr *)&sock->peer_addr, sizeof(sock->peer_addr)) < 0) {
        if (errno != EINPROGRESS) {
            clib_warning("Failed to connect to BGP peer: %s", strerror(errno));
            return -1;
        }
    }

    clib_warning("Connecting to BGP peer: socket_fd=%d, peer_addr=%s:%d",
                sock->socket_fd,
                inet_ntoa(*(struct in_addr *)&sock->peer_addr.sin_addr),
                ntohs(sock->peer_addr.sin_port));

    
                
    return 0;
}

/* Send a BGP message */
int bgp_socket_send(int socket_fd, void *message, size_t length) {
    fd_set write_fds;
    struct timeval timeout;
    int result;

    FD_ZERO(&write_fds);
    FD_SET(socket_fd, &write_fds);

    timeout.tv_sec = 5;  // 5 seconds timeout
    timeout.tv_usec = 0;

    result = select(socket_fd + 1, NULL, &write_fds, NULL, &timeout);
    if (result <= 0) {
        clib_warning("Socket not ready for writing: %s", strerror(errno));
        return -1;
    }

    ssize_t sent = send(socket_fd, message, length, 0);
    if (sent < 0) {
        clib_warning("Failed to send message: %s", strerror(errno));
        return -1;
    }

    clib_warning("******Sending a BGP message on socket: sent(%d)", sent);
    return 0;
}

/* Close the BGP socket */
void bgp_socket_close(int socket_fd) {
    if (close(socket_fd) < 0) {
        clib_warning("Failed to close socket: %s", strerror(errno));
    }
}