/*
package client_handler implements code to manage one client, independant of
server protocol.
*/
package client_handler

import (
   	"test_server/output_generator"
	"log"
)

const (
	incomingChanCapacity = 100
)

// ClientHandler is the interface to the client handler
type ClientHandler interface {
	Close()
	SetUserNameAndDeviceId(userName string, deviceId int)
	ReportChallengeResponse(responseValue []byte)
}

type clientHandler struct {
	incoming chan<- interface{}
	userName string
	device int
}

type setUserNameAndDeviceIdRequest struct {
	userName string
	deviceId int
}

type reportChallengeResponse struct {
	responseValue []byte
}

// create a new entity supporting the ClientHandler interface
func New(outputGenerator output_generator.OutputGenerator) ClientHandler {
	incomingChan := make(chan interface{}, incomingChanCapacity)
	go run(incomingChan, outputGenerator)
	return &clientHandler{incoming: incomingChan}
}

func run(incoming <-chan interface{}, 
	outputGenerator output_generator.OutputGenerator) {
	var userName string
	var deviceId int 
			
	for item := range incoming {
		switch item := item.(type) {
		case setUserNameAndDeviceIdRequest:
			userName = item.userName
			deviceId = item.deviceId
			log.Printf("DEBUG: userName = %s, deviceId = %d", userName, 
				deviceId)
			outputGenerator.IssueLoginChallenge(
				"N\x86\r\xaa\r\xf3\x99Q\xe1*\xfc\x06\x1d\xf3\xf8N")
		case reportChallengeResponse:
			log.Printf("DEBUG: challenge response = %q", item.responseValue)
			outputGenerator.AcceptChallengeResponse()
		}
	}
}

func (handler *clientHandler) Close() {
	close(handler.incoming)
}

func (handler *clientHandler) SetUserNameAndDeviceId(userName string, 
	deviceId int) {
	handler.incoming <- setUserNameAndDeviceIdRequest{userName, deviceId}
}

func (handler *clientHandler) ReportChallengeResponse(responseValue []byte) {
	handler.incoming <- reportChallengeResponse{responseValue}
}
