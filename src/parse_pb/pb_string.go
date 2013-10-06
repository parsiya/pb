/* package parse_pb

This package interprets a stream of PB bytes
*/

package parse_pb

import (
	"fmt"
	"io"
)

type PBString struct {
	value []byte
}

func NewPBString(data string) PBString {
	return PBString{value: []byte(data)}
}

func UnmarshalPBString(intBuffer []byte, parser *Parser) (PBString, error) {
	size := unmarshalBase128Int(intBuffer)
	data, err := parser.readAll(size)
	if err != nil {
		return PBString{}, err
	}
	return PBString{value: data}, nil
}

func (item PBString) Type() byte {
	return PB_STRING
}

func (item PBString) String() string {
	return fmt.Sprintf("PB_STRING(%q)", item.value)
}

func (item PBString) Marshal(writer io.Writer) error {
	marshaledLen, err := marshalBase128Int(len(item.value))
	if err != nil {
		return err
	}
	if _, err := writer.Write(marshaledLen); err != nil {
		return err
	}
	if _, err := writer.Write([]byte{PB_STRING}); err != nil {
		return err
	}
	if _, err := writer.Write(item.value); err != nil {
		return err
	}
	return nil
}
