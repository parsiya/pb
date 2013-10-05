/* program test_server
   package main

   connectionHandler

   This program acts as a server for testing pb communication
*/
package main

import (
   "log"
   "net"
)

func handleConnection(connection net.Conn) {
	log.Printf("DEBUG: handleConnection %s", connection.RemoteAddr())
	connection.Close()
}
