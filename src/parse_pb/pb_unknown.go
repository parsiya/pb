/* package parse_pb

This package interprets a stream of PB bytes
*/

package parse_pb

import (
	"fmt"
)

type PBUnknown struct {
	intBuffer []byte
	c         byte
}

func NewPBUnknown(c byte, intBuffer []byte) PBUnknown {
	return PBUnknown{intBuffer: intBuffer, c: c}
}

func (item PBUnknown) String() string {
	return fmt.Sprintf("* %v %s *", item.intBuffer, dumpByte(item.c))
}

func (item PBUnknown) Type() byte {
	return 0
}
