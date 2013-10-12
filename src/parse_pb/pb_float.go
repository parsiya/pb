/* package parse_pb

This package interprets a stream of PB bytes
*/

package parse_pb

import (
	"fmt"
	"io"
)

type PBFloat struct {
	Value float64
}

func NewPBFloat(Value float64) PBFloat {
	return PBFloat{Value: Value}
}

func UnmarshalPBFloat(parser *Parser) (PBFloat, error) {
	marshaledFloat, err := parser.readAll(PackedFloatSliceSize)
	if err != nil {
		return PBFloat{}, err
	}
	return PBFloat{Value: unmarshalPackedFloat(marshaledFloat)}, nil
}

func (item PBFloat) Type() byte {
	return PB_FLOAT
}

func (item PBFloat) String() string {
	return fmt.Sprintf("PB_FLOAT(%f)", item.Value)
}

func (item PBFloat) Marshal(writer io.Writer) error {
	if _, err := writer.Write([]byte{PB_FLOAT}); err != nil {
		return err
	}
	if _, err := writer.Write(marshalPackedFloat(item.Value)); err != nil {
		return err
	}
	return nil
}
