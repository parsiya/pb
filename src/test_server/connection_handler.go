/* program test_server
   package main

   connectionHandler

   This program acts as a server for testing pb communication
*/
package main

import (
   "parse_pb"
   "test_server/output_generator"
   "test_server/client_handler"
   "fmt"
   "io"
   "log"
   "net"
   "strings"
   "strconv"
)

type State struct {
	Connection net.Conn
	Parser *parse_pb.Parser
	ClientHandler client_handler.ClientHandler
}

// state function inspired by Rob Pike's video 'Lexical Scanning in Go'
type stateFunction func(state *State) stateFunction

// handleConnection loops through state functions to handle messages
// to and from the client
func handleConnection(connection net.Conn) {
	log.Printf("INFO: handleConnection %s", connection.RemoteAddr())
	defer connection.Close()

	parser, err := parse_pb.NewParser(connection)
	if err != nil {
		log.Fatalf("CRITICAL: parse_pb.NewParser failed %s", err)
	}

	outputGenerator := output_generator.New(connection)
	defer outputGenerator.Close()

	clientHandler := client_handler.New(outputGenerator)

	state := State{Connection: connection, Parser: parser, 
		ClientHandler: clientHandler}
	defer state.ClientHandler.Close()

	// main loop of state functions
	for f := startState; f != nil; f = f(&state) {
	}

	log.Printf("INFO: handleConnection ends %s", connection.RemoteAddr())
}

// startState sends the initial greeting to the client and waits for the 
// client's protocol choice
func startState(state *State) stateFunction {

	/* -----------------------------------------------------------------------
	** PB_LIST(
	**    PB_STRING("pb"),
	**    PB_STRING("none"))
	** ---------------------------------------------------------------------*/
	pbString := parse_pb.NewPBString("pb")
	noneString := parse_pb.NewPBString("none")
	greeting := parse_pb.NewPBList(pbString, noneString)

	// send the initial greeting to the client
	if err := greeting.Marshal(state.Connection); err != nil {
		log.Fatalf("CRITICAL: greeting.Marshal %s %s", greeting.String(), err)
	}

	/* -----------------------------------------------------------------------
	** PB_LIST(
	**    PB_VOCAB(Version),
	**    PB_INT(6))
	** ---------------------------------------------------------------------*/
	version := parse_pb.NewPBList(parse_pb.NewPBVocab(parse_pb.VocabVersion),
		parse_pb.NewPBInt(6))

	// send the initial greeting to the client
	if err := version.Marshal(state.Connection); err != nil {
		log.Fatalf("CRITICAL: version.Marshal %s %s", version.String(), err)
	}

	// we expect the user to send the string 'pb'
	rawResult, err := state.Parser.Step()
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
	rawResult, err = state.Parser.Step()
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
func loginState(state *State) stateFunction {
	// now we expect the user to send the login message
	rawResult, err := state.Parser.Step()
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
		log.Printf("ERROR: (expecting root request) %s", rawRootRequest)
		return nil
	}

	userName, deviceId, err := parseRootRequest(rootRequest)
	if err != nil {
		log.Printf("ERROR: error parsing root request %s", err)
		return nil
	}
	state.ClientHandler.SetUserNameAndDeviceId(userName, deviceId)

	// now we expect the user to send a response to the challenge
	rawResponse, err := state.Parser.Step()
	if err != nil {
		log.Printf("ERROR: (expecting challenge response) %s", err)
		return nil
	}
	responseList, ok := rawResponse.(parse_pb.PBList)
	if !ok {
		log.Printf("ERROR: (expecting response list) %s", rawResponse)
	}
	responseItem := responseList.Reparse()
	response, ok := responseItem.(parse_pb.PBMessageList)
	if !ok {
		log.Printf("ERROR: (expecting challenge response) %s", responseItem)
		return nil
	}

	err = parseResponse(response)
	if err != nil {
		log.Printf("ERROR: error parsing response %s", err)
		return nil
	}

	responseAnswer := constructResponseAnswer()

	// send challenge answer to the client
	if err := responseAnswer.Marshal(state.Connection); err != nil {
		log.Fatalf("CRITICAL: answer.Marshal %s %s", 
			responseAnswer.String(), err)
	}

	return runState
}

// runState handles message traffic for a fully connected client
func runState(state *State) stateFunction {
	for {
		rawResult, err := state.Parser.Step()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("CRITICAL: error on parser.Step %s", err)
		}

		// at this stage we only understand lists
		resultList, ok := rawResult.(parse_pb.PBList)
		if !ok {
			log.Printf("ERROR: unknown type from client %s", 
				resultList.String())
			continue
		}

		result := resultList.Reparse()

		switch result := result.(type) {
		case parse_pb.PBMessageList:
			if err := handleIncomingMessage(state, result); err != nil {
				log.Printf("ERROR: error in incoming message %s %s", 
					result.String(), err)
			}
			continue
		case parse_pb.PBAnswerList:
			if err := handleIncomingAnswer(state, result); err != nil {
				log.Printf("ERROR: error in incoming answer %s %s", 
					result.String(), err)
			}
			continue
		}
		log.Printf("ERROR: unhandled input from client %s", result.String())
	}
	return nil
}

func parseRootRequest(rootRequest parse_pb.PBObjectMessageList) (
	string, int, error) {
	var authName parse_pb.PBString

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

	internalList, ok := rootRequest.PBList.Value[5].(parse_pb.PBList)
	if !ok {
		return "", 0, fmt.Errorf("can't cast internal list %s", rootRequest)
	}

	// if we have a list, it's twisted unicode, 
	// if a string its the userid
	switch internalItem := internalList.Value[1].(type) {
	case parse_pb.PBList:
		authName, ok = internalItem.Value[1].(parse_pb.PBString)
		if !ok {
			return "", 0, fmt.Errorf("can't cast expecting authName %s", 
				internalItem)
		}
	case parse_pb.PBString:
		authName = internalItem
	default:
		return "", 0, fmt.Errorf("unexpected type in root request %s %s", 
			internalList, rootRequest)
	}

	// we expect a string of the form <user-name>@<device-id>
	splitName := strings.Split(string(authName.Value), "@")
	if len(splitName) != 2 {
		return "", 0, fmt.Errorf("Unparseable authName '%s'", authName)
	}

	deviceId, err := strconv.Atoi(splitName[1])
	if err != nil {
		return "", 0, fmt.Errorf("Unparseable device-id '%s' %s", 
			authName, err)
	}

	return splitName[0], deviceId, nil
}

func parseResponse(_ parse_pb.PBMessageList) error {
	/* -----------------------------------------------------------------------
	** PB_LIST(
	**     PB_VOCAB(Message),
	**     PB_INT(2),
	**     PB_INT(1),
	**     PB_STRING("respond"),
	**     PB_INT(1),
	**     PB_LIST(
	**         PB_VOCAB(Tuple),
	**         PB_STRING("\xfd\x83NP\xb7\xa4\x02?\xafu\xb6\xe3\f<\xf1j"),
	**         PB_LIST(
	**             PB_VOCAB(Remote),
	**             PB_INT(1))),
	**     PB_LIST(
	**         PB_VOCAB(Dictionary)))
	** ---------------------------------------------------------------------*/

	return nil
}

func constructResponseAnswer() parse_pb.PBList {
	/* -----------------------------------------------------------------------
	** PB_LIST(
	**     PB_VOCAB(Answer),
	**     PB_INT(2),
	**     PB_LIST(
	**         PB_VOCAB(Remote),
	**         PB_INT(2)))
	** ---------------------------------------------------------------------*/
	remote := parse_pb.NewPBList(parse_pb.NewPBVocab(parse_pb.VocabRemote),
		parse_pb.NewPBInt(2))
	answer := parse_pb.NewPBList(parse_pb.NewPBVocab(parse_pb.VocabAnswer),
		parse_pb.NewPBInt(2), remote)
	return answer
}

func handleIncomingMessage(state *State, message parse_pb.PBMessageList) error {
	switch message.Name {
	case "ping":
		return sendPingReply(message.Sequence, state.Connection)
	}
	return fmt.Errorf("Unknown message from client %s", message)
}

func handleIncomingAnswer(state *State, message parse_pb.PBAnswerList) error {
	return nil
}

func sendPingReply(sequence int, writer io.Writer) error {
	/* -----------------------------------------------------------------------
	** PB_LIST(
	**     PB_VOCAB(Answer),
	**     PB_INT(3),
	**     PB_LIST(
	**         PB_STRING("boolean"),
	**         PB_STRING("true")))
	** ---------------------------------------------------------------------*/
	answer := parse_pb.NewPBList(parse_pb.NewPBVocab(parse_pb.VocabAnswer),
		parse_pb.NewPBInt(sequence),parse_pb.NewPBList(
			parse_pb.NewPBString("boolean"), parse_pb.NewPBString("true")))

	log.Printf("DEBUG: sendPingReply %s", answer)
	return answer.Marshal(writer)
}
