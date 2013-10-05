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

func (item PBUnknown) String() string {
	return fmt.Sprintf("* %v %s *", item.intBuffer, dumpByte(item.c))
}

func (item PBUnknown) Type() byte {
	return 0
}
