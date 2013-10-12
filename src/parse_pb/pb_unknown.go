/* package parse_pb

This package interprets a stream of PB bytes
*/

package parse_pb

import (
	"fmt"
	"io"
)

type PBUnknown struct {
	intBuffer []byte
	c         byte
}

func NewPBUnknown(c byte, intBuffer []byte) PBUnknown {
	return PBUnknown{intBuffer: intBuffer, c: c}
}

func UnmarshalUnknown(intBuffer []byte, c byte) (PBUnknown, error) {
	return NewPBUnknown(c, intBuffer), nil
}

func (item PBUnknown) String() string {
	return fmt.Sprintf("* %v %s *", item.intBuffer, dumpByte(item.c))
}

func (item PBUnknown) Marshal(writer io.Writer) error {
	return fmt.Errorf("attempt to marshal unknown item %s", item.String())
}
