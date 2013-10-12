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

func (item PBList) Type() byte {
	return PB_LIST
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
