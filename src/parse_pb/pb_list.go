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
	Value []ParseItem
}

type PBVersionList struct {
	PBList
	Version int
}

type PBAnswerList struct {
	PBList
	Sequence int
}

func NewPBList(items ...ParseItem) PBList {
	return PBList{Value: items}
}

func UnmarshalPBList(intBuffer []byte, parser *Parser) (PBList, error) {
	size := unmarshalBase128Int(intBuffer)
	list := PBList{Value: make([]ParseItem, size)}
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
func (item PBList) Reparse() ParseItem {

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

func parseVersionList(item PBList) ParseItem {

	// PB_LIST(PB_VOCAB(Version),PB_INT(6))

	intResult, ok := item.Value[1].(PBInt)
	if !ok {
		return item
	}

	return PBVersionList{PBList: item, Version: intResult.Value}
}

func parseAnswerList(item PBList) ParseItem {
	intResult, ok := item.Value[1].(PBInt)
	if !ok {
		return item
	}

	return PBVersionList{PBList: item, Version: intResult.Value}
}
