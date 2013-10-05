/* program test_server
   package main

   connectionHandler

   This program acts as a server for testing pb communication
*/
package main

import (
   "parse_pb"
   "io"
   "log"
   "net"
)

var (
	pbString = parse_pb.NewPBString("pb")
	noneString = parse_pb.NewPBString("none")
	greeting = parse_pb.NewPBList(pbString, noneString)
)

func handleConnection(connection net.Conn) {
	log.Printf("INFO: handleConnection %s", connection.RemoteAddr())
	defer connection.Close()

	parser, err := parse_pb.NewParser(connection)
	if err != nil {
		log.Fatalf("CRITICAL: parse_pb.NewParser failed %s", err)
	}
	for {
		result, err := parser.Step()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("CRITICAL: error on parser.Step %s", err)
		}
		log.Printf("DEBUG: %s", result)
	}
	log.Printf("INFO: handleConnection ends %s", connection.RemoteAddr())
}
