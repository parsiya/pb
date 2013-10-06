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

func UnmarshalPBInt(intBuffer []byte) (PBInt, error) {
	return NewPBInt(unmarshalBase128Int(intBuffer)), nil
}

func (item PBInt) Type() byte {
	return PB_INT
}

func (item PBInt) String() string {
	return fmt.Sprintf("PB_INT(%d)", item.value)
}

func (item PBInt) Marshal(writer io.Writer) error {
	marshaledInt, err := marshalBase128Int(item.value)
	if err != nil {
		return err
	}
	if _, err := writer.Write(marshaledInt); err != nil {
		return err
	}
	if _, err := writer.Write([]byte{PB_INT}); err != nil {
		return err
	}
	return nil
}
