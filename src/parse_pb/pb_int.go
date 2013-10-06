/* package parse_pb

This package interprets a stream of PB bytes
*/

package parse_pb

import (
	"fmt"
	"io"
)

type PBInt struct {
	value int
}

func NewPBInt(value int) PBInt {
	return PBInt{value: value}
}

func (item PBInt) Type() byte {
	return PB_INT
}

func (item PBInt) String() string {
	return fmt.Sprintf("PB_INT(%d)", item.value)
}

func (item PBInt) Marshall(writer io.Writer) error {
	marshalledInt, err := marshallBase128Int(item.value) 
	if err != nil {
		return err
	}
	if _, err := writer.Write(marshalledInt); err != nil {
		return err
	}
	if _, err := writer.Write([]byte{PB_INT}); err != nil {
		return err
	}
	return nil
}
