/* program test_server
   package main

   This program acts as a server for testing pb communication
*/
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	signalChannelCapacity = 1
	listenAddress = ":6666"
)

// main entry point for ergo_prpoxy
func main() {
	log.Printf("INFO: program starts")

	// set up a signal handling channel
	signalChannel := make(chan os.Signal, signalChannelCapacity)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	go listen(listenAddress)

	// Block until we ge a signal from the os
	signal := <-signalChannel
	log.Printf("INFO: terminated by signal: %v", signal)
}
