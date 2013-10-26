/* package parse_pb

This package interprets a stream of PB bytes
*/
package parse_pb

import (
	"fmt"
	"io"
	"log"
	"strings"
)
/* ---------------------------------------------------------------------------
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
** -------------------------------------------------------------------------*/

type PBObjectMessageList struct {
	PBList
	Sequence int
	ObjectName string
}

// PBMesageList represents a python function call
type PBMessageList struct {
	Sequence int
	ObjectNumber int
	Name string
	PositionalArgs PBList
	KeywordArgs PBList
}

func NewPBMessageList(sequence int, objectNumber int, name string, 
	positionalArgs PBList, keywordArgs PBList) PBMessageList {
	return PBMessageList{Sequence: sequence, ObjectNumber: objectNumber,
		Name: name, PositionalArgs: positionalArgs, KeywordArgs: keywordArgs}
}

func (message PBMessageList) String() string {
	var printPositionalArgs []string
	for _, x := range message.PositionalArgs.Value {
		printPositionalArgs = append(printPositionalArgs, x.String())
	}
	
	var printKeywordArgs []string
	for _, x := range message.KeywordArgs.Value {
		printKeywordArgs = append(printKeywordArgs, x.String())
	}
	
	return fmt.Sprintf("Message(%s (%s) %s %s)", message.Name, 
		message.Sequence, strings.Join(printPositionalArgs, ","), 
		strings.Join(printKeywordArgs, ","))
}

func (message PBMessageList) Marshal(writer io.Writer) error {
	marshaledLen, err := marshalBase128Int(7) 
	if err != nil {
		return err
	}
	if _, err := writer.Write(marshaledLen); err != nil {
		return err
	}

	if _, err := writer.Write([]byte{PB_LIST}); err != nil {
		return err
	}
	messageVocabItem := NewPBVocab(VocabMessage)
	if err := messageVocabItem.Marshal(writer); err != nil {
		return err
	}

	sequenceItem := NewPBInt(message.Sequence)
	if err := sequenceItem.Marshal(writer); err != nil {
		return err
	}

	objectNumberItem := NewPBInt(message.ObjectNumber)
	if err := objectNumberItem.Marshal(writer); err != nil {
		return err
	}

	nameItem := NewPBString(message.Name)
	if err := nameItem.Marshal(writer); err != nil {
		return err
	}

	x := NewPBInt(1) // I don't know what this number is
	if err := x.Marshal(writer); err != nil {
		return err
	} 

	// positional arguments
	if _, err := writer.Write([]byte{PB_LIST}); err != nil {
		return err
	}
	tupleVocabItem := NewPBVocab(VocabTuple)
	if err := tupleVocabItem.Marshal(writer); err != nil {
		return err
	} 
	for _, positionalItem := range message.PositionalArgs.Value {
		if err := positionalItem.Marshal(writer); err != nil {
			return err
		} 
	}

	// keyword arguments
	if _, err := writer.Write([]byte{PB_LIST}); err != nil {
		return err
	}
	dictVocabItem := NewPBVocab(VocabDictionary)
	if err := dictVocabItem.Marshal(writer); err != nil {
		return err
	} 
	for _, keywordItem := range message.KeywordArgs.Value {
		if err := keywordItem.Marshal(writer); err != nil {
			return err
		} 
	}
	
	return nil
}

func parseMessageList(item PBList) ParseItem {

	// we expect a sequence number after PBVocab
	// PB_LIST(PB_VOCAB(Message),PB_INT(1)...
	intResult, ok := item.Value[1].(PBInt)
	if !ok {
		return item
	}
	sequence := intResult.Value

    // if the item after the Sequence is a PBString, we assume this is a request
	// for a remote object reference
	// PB_LIST(PB_VOCAB(Message),PB_INT(1),PB_STRING("root"),...

	// otherwise, we expect an int, which is the index of the remote object 
	// reference
    // PB_LIST(PB_VOCAB(Message),PB_INT(2),PB_INT(1),PB_STRING("respond"),...

	var objectNumber int

	switch nextEntry := item.Value[2].(type) {
	case PBString:
		return PBObjectMessageList{PBList: item, Sequence: sequence,
			ObjectName: string(nextEntry.Value)}
	case PBInt:
		objectNumber = nextEntry.Value
	default:
		return item
	}

	// now we expect a string identifying the message
	stringResult, ok := item.Value[3].(PBString)
	if !ok {
		log.Printf("ERROR: expecting string for item[3] %s", 
			item.String())
		return item
	}
	name := string(stringResult.Value)

	// now we expect a 1, I don't know what it represents
	_, ok = item.Value[4].(PBInt)
	if !ok {
		log.Printf("ERROR: expecting int for item[4] %s", 
			item.String())
		return item
	}

	// now we expect a PBList of positional arguments
	positionalList, ok := item.Value[5].(PBList)
	if !ok {
		log.Printf("ERROR: expecting list for item[5] %s", 
			item.String())
		return item
	}

	// we expect a tuple
	positionalVocab, ok := positionalList.Value[0].(PBVocab)
	if !ok {
		log.Printf("ERROR: expecting vocab as positionalList[0] %s", 
			item.String())
		return item
	}
	if positionalVocab.Value != VocabTuple {
		log.Printf("ERROR: expecting VocabTuple as positionalList[0] %s", 
			item.String())
		return item
	}

	// now we expect a PBList of keyword arguments
	keywordList, ok := item.Value[6].(PBList)
	if !ok {
		log.Printf("ERROR: expecting list for item[6] %s", 
			item.String())
		return item
	}

	// we expect a dictionary
	keywordVocab, ok := keywordList.Value[0].(PBVocab)
	if !ok {
		log.Printf("ERROR: expecting vocab as keywordList[0] %s", 
			item.String())
		return item
	}
	if keywordVocab.Value != VocabDictionary {
		log.Printf("ERROR: expecting VocabDict as keywordList[0] %s", 
			item.String())
		return item
	}

	return NewPBMessageList(sequence, objectNumber, name, 
		NewPBList(positionalList.Value[1:]...), 
		NewPBList(keywordList.Value[1:]...)) 	
}
