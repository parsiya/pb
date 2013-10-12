/* package parse_pb

This package interprets a stream of PB bytes
*/

package parse_pb

import (
	"fmt"
	"io"
)

type PBInt struct {
	Value int
}

func NewPBInt(Value int) PBInt {
	return PBInt{Value: Value}
}

func UnmarshalPBInt(intBuffer []byte) (PBInt, error) {
	return NewPBInt(unmarshalBase128Int(intBuffer)), nil
}

func (item PBInt) String() string {
	return fmt.Sprintf("PB_INT(%d)", item.Value)
}

func (item PBInt) Marshal(writer io.Writer) error {
	marshaledInt, err := marshalBase128Int(item.Value)
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
