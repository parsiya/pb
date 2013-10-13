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

// state function inspired by Rob Pike's video 'Lexical Scanning in Go'
type stateFunction func(parser * parse_pb.Parser, 
	connection net.Conn) stateFunction

// handleConnection loops through state functions to handle messages
// to and from the client
func handleConnection(connection net.Conn) {
	log.Printf("INFO: handleConnection %s", connection.RemoteAddr())
	defer connection.Close()

	parser, err := parse_pb.NewParser(connection)
	if err != nil {
		log.Fatalf("CRITICAL: parse_pb.NewParser failed %s", err)
	}

	for state := startState; state != nil; state = state(parser, connection) {
	}

	log.Printf("INFO: handleConnection ends %s", connection.RemoteAddr())
}

// startState sends the initial greeting to the client and waits for the 
// client's protocol choice
func startState(parser * parse_pb.Parser, connection net.Conn) stateFunction {
	pbString := parse_pb.NewPBString("pb")
	noneString := parse_pb.NewPBString("none")
	greeting := parse_pb.NewPBList(pbString, noneString)

	// send the initial greeting to the client
	if err := greeting.Marshal(connection); err != nil {
		log.Fatalf("CRITICAL: greeting.Marshal %s %s", greeting.String(), err)
	}

	// we expect the user to send the string 'pb'
	rawResult, err := parser.Step()
	if err != nil {
		log.Printf("ERROR: (expecting 'pb') %s", err)
		return nil
	}
	clientProtocolResult, ok := rawResult.(parse_pb.PBString)
	if !ok {
		log.Printf("ERROR: (expecting 'pb') %s", clientProtocolResult)
		return nil
	} 
	clientProtocol := string(clientProtocolResult.Value)
	log.Printf("DEBUG: client protocol %s", clientProtocol)

	// now we expect the user to send the version
	rawResult, err = parser.Step()
	if err != nil {
		log.Printf("ERROR: (expecting version) %s", err)
		return nil
	}
	versionResult, ok := rawResult.(parse_pb.PBList)
	if !ok {
		log.Printf("ERROR: (expecting version) %s", versionResult)
		return nil
	}
	rawVersionItem := versionResult.Reparse()
	versionItem, ok := rawVersionItem.(parse_pb.PBVersionList)
	if !ok {
		log.Printf("ERROR: (expecting version) %s", rawVersionItem)
		return nil
	}
	log.Printf("DEBUG: client version %d", versionItem.Version)

	return loginState
}

// loginState handles message traffic for authenticating the client
func loginState(parser * parse_pb.Parser, connection net.Conn) stateFunction {
	// now we expect the user to send the login message
	rawResult, err := parser.Step()
	if err != nil {
		log.Printf("ERROR: (expecting login) %s", err)
		return nil
	}
	result, ok := rawResult.(parse_pb.PBList)
	if !ok {
		log.Printf("ERROR: (expecting login) %s", result)
	}

	// we expect a request for the root object
	rawRootRequest := result.Reparse()
	rootRequest, ok := rawRootRequest.(parse_pb.PBObjectMessageList)
	if !ok {
		log.Printf("ERROR: (expecting version) %s", rawRootRequest)
		return nil
	}

	/*------------------------------------------------------------------------
	** we're expecting this:
	** PB_LIST(
	** 0 PB_VOCAB(Message),
	** 1 PB_INT(1),
	** 2 PB_STRING("root"),
	** 3 PB_VOCAB(Login),
	** 4 PB_INT(1),
	** 5 PB_LIST(
	**   0 PB_VOCAB(Tuple),
	**   1 PB_LIST(
	**     0 PB_STRING("unicode"),
	**     1 PB_STRING("FunctionalTestUser1@1"))),
	** 6 PB_LIST(
	**   0 PB_VOCAB(Dictionary)))
	** TODO: parse the whole thing
	** XXX: need to handle the non unicode case
	**-----------------------------------------------------------------------*/
	var userString parse_pb.PBString

	internalList, ok := rootRequest.PBList.Value[5].(parse_pb.PBList)
	if !ok {
		log.Printf("ERROR: (internal list) %s", rootRequest)
		return nil
	} 
	// if we have a list, it's twisted unicode, if a string its the userid
	switch internalItem := internalList.Value[1].(type) {
	case parse_pb.PBList:
		userString, ok = internalItem.Value[1].(parse_pb.PBString)
		if !ok {
			log.Printf("ERROR: (expecting userString) %s", internalItem)
			return nil
		}
	case parse_pb.PBString:
		userString = internalItem
	default:
		log.Printf("ERROR: unexpected type in root request %s %s", 
			internalList, rootRequest)
		return nil
	}

	log.Printf("DEBUG: userString = %s", string(userString.Value))
	log.Printf("DEBUG: %s", result)
	return nil
}

// runState handles message traffic for a fully connected client
func runState(parser * parse_pb.Parser, connection net.Conn) stateFunction {
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
	return nil
}

