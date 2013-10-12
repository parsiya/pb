/* package parse_pb

This package interprets a stream of PB bytes
*/

package parse_pb

import (
	"fmt"
	"io"
)

type PBNeg struct {
	Value int
}

func NewPBNeg(Value int) PBNeg {
	return PBNeg{Value: Value}
}

func UnmarshalPBNeg(intBuffer []byte) (PBNeg, error) {
	return NewPBNeg(-unmarshalBase128Int(intBuffer)), nil
}

func (item PBNeg) String() string {
	return fmt.Sprintf("PB_NEG(%d)", item.Value)
}

func (item PBNeg) Marshal(writer io.Writer) error {
	marshaledNeg, err := marshalBase128Int(-item.Value) // must be positive
	if err != nil {
		return err
	}
	if _, err := writer.Write(marshaledNeg); err != nil {
		return err
	}
	if _, err := writer.Write([]byte{PB_NEG}); err != nil {
		return err
	}
	return nil
}
