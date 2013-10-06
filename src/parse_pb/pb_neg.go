/* package parse_pb

This package interprets a stream of PB bytes
*/

package parse_pb

import (
	"fmt"
	"io"
)

type PBNeg struct {
	value int
}

func NewPBNeg(value int) PBNeg {
	return PBNeg{value: value}
}

func UnmarshallPBNeg(intBuffer []byte) (PBNeg, error) {
	return NewPBNeg(-unmarshallBase128Int(intBuffer)), nil
}

func (item PBNeg) String() string {
	return fmt.Sprintf("PB_NEG(%d)", item.value)
}

func (item PBNeg) Type() byte {
	return PB_NEG
}

func (item PBNeg) Marshall(writer io.Writer) error {
	marshalledNeg, err := marshallBase128Int(-item.value) // must be positive
	if err != nil {
		return err
	}
	if _, err := writer.Write(marshalledNeg); err != nil {
		return err
	}
	if _, err := writer.Write([]byte{PB_NEG}); err != nil {
		return err
	}
	return nil
}
