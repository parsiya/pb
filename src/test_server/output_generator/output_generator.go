/*
package output_generator generates messages in the pb protocol from 
generic messages
*/
package output_generator

import (
   "parse_pb"
	"io"
	"log"
)

const (
	outgoingChanCapacity  = 100
)

type OutputGenerator interface {
	Close()
	IssueLoginChallenge(challenge string)
	AcceptChallengeResponse()
}

type outputGenerator struct {
	outgoing chan<- interface{}
}

type issueLoginChallenge struct {
	challenge string
}

type acceptChallengeResponse struct {}

// create a new entity supporting the OutputGenerator interface
func New(writer io.Writer) OutputGenerator {
	outgoingChan := make(chan interface{}, outgoingChanCapacity)
	go run(outgoingChan, writer)
	return &outputGenerator{outgoing: outgoingChan}
}

func run(outgoing <-chan interface{}, writer io.Writer) {
	for item := range outgoing {
		switch item := item.(type) {
		case issueLoginChallenge:
			log.Printf("DEBUG: issueLoginChallenge")
			loginChallenge := constructLoginChallenge(item.challenge)
			
			if err := loginChallenge.Marshal(writer); err != nil {
				log.Fatalf("CRITICAL: loginChallenge.Marshal %s %s", 
					loginChallenge.String(), err)
				}
		case acceptChallengeResponse:
			log.Printf("DEBUG: acceptChallengeResponse")
			acceptance := constructChallengeAcceptance()
			
			if err := acceptance.Marshal(writer); err != nil {
				log.Fatalf("CRITICAL: acceptance.Marshal %s %s", 
					acceptance.String(), err)
				}

		}
	}
}

func (generator *outputGenerator) Close() {
	close(generator.outgoing)
}

func (generator *outputGenerator) IssueLoginChallenge(challenge string) {
	generator.outgoing <- issueLoginChallenge{challenge}
}

func (generator *outputGenerator) AcceptChallengeResponse() {
	generator.outgoing <- acceptChallengeResponse{}
}

func constructLoginChallenge(challengeText string) parse_pb.PBList {
	/* ----------------------------------------------------------------------- 
	** PB_LIST(
	**     PB_VOCAB(Answer),
	**     PB_INT(1),
	**     PB_LIST(
	**         PB_VOCAB(Tuple),
	**         PB_STRING("N\x86\r\xaa\r\xf3\x99Q\xe1*\xfc\x06\x1d\xf3\xf8N"),
	**         PB_LIST(
	**             PB_VOCAB(Remote),
	**             PB_INT(1))))
	**---------------------------------------------------------------------*/
	remote := parse_pb.NewPBList(parse_pb.NewPBVocab(parse_pb.VocabRemote),
		parse_pb.NewPBInt(1))

	var challenge = []parse_pb.ParseItem{
		parse_pb.NewPBVocab(parse_pb.VocabTuple),
		parse_pb.NewPBString(challengeText),
		remote}

	var dummyAnswer = []parse_pb.ParseItem{
		parse_pb.NewPBVocab(parse_pb.VocabAnswer),
		parse_pb.NewPBInt(1),
		parse_pb.NewPBList(challenge...)}

	// TODO: construct a real answer
	return parse_pb.NewPBList(dummyAnswer...)
}

func constructChallengeAcceptance() parse_pb.PBList {
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

