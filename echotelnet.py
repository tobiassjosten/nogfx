#!/usr/bin/env python3

import logging
import socket
import select

HOSTNAME = 'localhost'
PORT = '4000'

MAXIMUM_QUEUED_CONNECTIONS = 5
RECEIVING_BUFFER_SIZE = 4096

logger = logging.getLogger(__name__)

def start_server(hostname, port):
    # Get all possible binding addresses for given hostname and port.
    possible_addresses = socket.getaddrinfo(
        hostname,
        port,
        family=socket.AF_UNSPEC,
        type=socket.SOCK_STREAM,
        flags=socket.AI_PASSIVE
    )
    server_socket = None
    # Look for an address that will actually bind.
    for family, socket_type, protocol, name, address in possible_addresses:
        try:
            # Create socket.
            server_socket = socket.socket(family, socket_type, protocol)
            # Make socket port reusable.
            server_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
            # Bind socket to the address.
            server_socket.bind(address)
        except OSError:
            # Try another address.
            continue
        break
    if server_socket is None:
        logger.error("No suitable address available.")
        return
    # Listen for incoming connections.
    server_socket.listen(MAXIMUM_QUEUED_CONNECTIONS)
    logger.info("Listening on %s port %d." % server_socket.getsockname()[:2])
    monitored_sockets = [server_socket]
    try:
        while True:
            # Wait for any of the monitored sockets to become readable.
            ready_to_read_sockets = select.select(
                monitored_sockets,
                tuple(),
                tuple()
            )[0]
            for ready_socket in ready_to_read_sockets:
                if ready_socket == server_socket:
                    # If server socket is readable, accept new client
                    # connection.
                    client_socket, client_address = server_socket.accept()
                    monitored_sockets.append(client_socket)
                    logger.info("New connection #%d on %s:%d." % (
                        client_socket.fileno(),
                        client_address[0],
                        client_address[1]
                    ))
                else:
                    message = ready_socket.recv(RECEIVING_BUFFER_SIZE)
                    if message:
                        # Client send correct message. Echo it.
                        logger.error(message)
                        ready_socket.sendall("you sent: ".encode('UTF-8') + message)
                    else:
                        # Client connection is lost. Handle it.
                        logger.info(
                            "Lost connection #%d." % ready_socket.fileno()
                        )
                        monitored_sockets.remove(ready_socket)
    except KeyboardInterrupt:
        pass
    logger.info("Shutdown initiated.")
    # Close client connections.
    monitored_sockets.remove(server_socket)
    for client_socket in monitored_sockets:
        logger.info("Closing connection #%d." % client_socket.fileno())
        client_socket.close()
    # Close server socket.
    logger.info("Shutting server down...")
    server_socket.close()

if __name__ == '__main__':
    # Configure logging.
    logger.setLevel(logging.INFO)
    logger.addHandler(logging.StreamHandler())
    # Start server.
    start_server(HOSTNAME, PORT)
