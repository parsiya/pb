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

func (item PBString) Type() byte {
	return PB_STRING
}

func (item PBString) String() string {
	return fmt.Sprintf("PB_STRING(%q)", item.value)
}

func (item PBString) Marshall(writer io.Writer) error {
	marshalledLen, err := marshallBase128Int(len(item.value)) 
	if err != nil {
		return err
	}
	if _, err := writer.Write(marshalledLen); err != nil {
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
