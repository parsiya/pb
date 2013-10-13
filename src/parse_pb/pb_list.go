/* package parse_pb

This package interprets a stream of PB bytes
*/
package parse_pb

import (
	"fmt"
	"io"
	"strings"
)

type PBList struct {
	Value []parseItem
}

type PBVersionList struct {
	PBList
	Version int
}

type PBObjectMessageList struct {
	PBList
	Sequence int
	ObjectName string
}

type PBMessageList struct {
	PBList
	Sequence int
	ObjectNumber int
	MessageName string
}

type PBAnswerList struct {
	PBList
	Sequence int
}

func NewPBList(items ...parseItem) PBList {
	return PBList{Value: items}
}

func UnmarshalPBList(intBuffer []byte, parser *Parser) (PBList, error) {
	size := unmarshalBase128Int(intBuffer)
	list := PBList{Value: make([]parseItem, size)}
	for j := 0; j < size; j++ {
		Value, err := parser.Step()
		if err != nil {
			return PBList{}, err
		}
		list.Value[j] = Value
	}

	return list, nil	
}

func (item PBList) String() string {
	var printValues []string
	for _, x := range item.Value {
		printValues = append(printValues, x.String())
	}
	return fmt.Sprintf("PB_LIST(%s)", strings.Join(printValues, ","))
}

func (item PBList) Marshal(writer io.Writer) error {
	marshaledLen, err := marshalBase128Int(len(item.Value)) 
	if err != nil {
		return err
	}
	if _, err := writer.Write(marshaledLen); err != nil {
		return err
	}
	if _, err := writer.Write([]byte{PB_LIST}); err != nil {
		return err
	}
	for _, subItem := range item.Value {
		if err := subItem.Marshal(writer); err != nil {
			return err
		} 
	}
	return nil
}

// Reparse attempts to match the internal structure of a PBList
// with a more fully defined object 
func (item PBList) Reparse() parseItem {

	// if you send us an empty list, you get back an empty list
	if len(item.Value) == 0 {
		return item
	}

	// We need an initial PBVocab to understand the internals
	vocabResult, ok := item.Value[0].(PBVocab)
	if !ok {
		return item
	}

	switch vocabResult.Value {
	case VocabVersion:
		return parseVersionList(item)
	case VocabMessage:
		return parseMessageList(item)
	case VocabAnswer:
		return parseAnswerList(item)
	}
	return item
}

func parseVersionList(item PBList) parseItem {

	// PB_LIST(PB_VOCAB(Version),PB_INT(6))

	intResult, ok := item.Value[1].(PBInt)
	if !ok {
		return item
	}

	return PBVersionList{PBList: item, Version: intResult.Value}
}

func parseMessageList(item PBList) parseItem {

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
	stringResult, ok := item.Value[3].(PBInt)
	if !ok {
		return item
	}
	messageName := string(stringResult.Value)

	return PBMessageList{PBList: item, Sequence: sequence, 
		ObjectNumber: objectNumber, MessageName: messageName}

}

func parseAnswerList(item PBList) parseItem {
	intResult, ok := item.Value[1].(PBInt)
	if !ok {
		return item
	}

	return PBVersionList{PBList: item, Version: intResult.Value}
}
