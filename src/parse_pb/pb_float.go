/* package parse_pb

This package interprets a stream of PB bytes
*/

package parse_pb

import (
	"fmt"
	"io"
)

type PBFloat struct {
	value float64
}

func NewPBFloat(value float64) PBFloat {
	return PBFloat{value: value}
}

func UnmarshalPBFloat(parser *Parser) (PBFloat, error) {
	marshaledFloat, err := parser.readAll(PackedFloatSliceSize)
	if err != nil {
		return PBFloat{}, err
	}
	return PBFloat{value: unmarshalPackedFloat(marshaledFloat)}, nil
}

func (item PBFloat) Type() byte {
	return PB_FLOAT
}

func (item PBFloat) String() string {
	return fmt.Sprintf("PB_FLOAT(%f)", item.value)
}

func (item PBFloat) Marshal(writer io.Writer) error {
	if _, err := writer.Write([]byte{PB_FLOAT}); err != nil {
		return err
	}
	if _, err := writer.Write(marshalPackedFloat(item.value)); err != nil {
		return err
	}
	return nil
}
