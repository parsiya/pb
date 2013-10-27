/*
package client_handler implements code to manage one client, independant of
server protocol.
*/
package client_handler

import (
	"log"
)

const (
	incomingChanCapacity = 100
	outgoingChanCapacity  = 100
)

// ClientHandler is the interface to the client handler
type ClientHandler interface {
	SetUserNameAndDeviceId(userName string, deviceId int)
	Close()
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

// create a new entity supporting the ClientHandler interface
func New(outgoingChan chan<- interface{}) ClientHandler {
	incomingChan := make(chan interface{}, incomingChanCapacity)
	go run(incomingChan)
	return &clientHandler{incoming: incomingChan}
}

func run(incoming <-chan interface{}, outgoing chan<- interface{}) {
	var userName string
	var deviceId int 
			
	for item := range incoming {
		switch item := item.(type) {
		case setUserNameAndDeviceIdRequest:
			userName = item.userName
			deviceId = item.deviceId
			log.Printf("DEBUG: userName = %s, deviceId = %d", userName, 
				deviceId)
		}
	}
}

func (handler *clientHandler) SetUserNameAndDeviceId(userName string, 
	deviceId int) {
	handler.incoming <- setUserNameAndDeviceIdRequest{userName, deviceId}
}

func (handler *clientHandler) Close() {
	close(handler.incoming)
}

