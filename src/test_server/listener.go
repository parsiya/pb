/* program test_server
   package main

   listener

   This program acts as a server for testing pb communication
*/
package main

import (
   "log"
   "net"
)

// listener listens for connections on the specified address
func listen(address string) {	
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("CRITICAL: listener.Listen error %s", err)
	}
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatalf("CRITICAL: listener.Accept error %s", err)
		}
		go handleConnection(connection)
	}
}
