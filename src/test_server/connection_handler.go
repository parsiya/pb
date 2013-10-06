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

	// send the initial greeting to the client
	if err := greeting.Marshal(connection); err != nil {
		log.Fatalf("CRITICAL: greeting.Marshal %s %s", greeting.String(), err)
	}

	for {
		result, err := parser.Step()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("CRITICAL: error on parser.Step %s", err)
		}
		switch result := result.(type) {
		case parse_pb.PBString:
			log.Printf("DEBUG: received PBString %s", result.String())
			continue
		case parse_pb.PBList:
			log.Printf("DEBUG: received PBList %s", result.String())
			continue
		}
		log.Fatalf("CRITICAL: unexpected input from client %s", result.String())
	}
	log.Printf("INFO: handleConnection ends %s", connection.RemoteAddr())
}
