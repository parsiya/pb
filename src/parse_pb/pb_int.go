/* package parse_pb

This package interprets a stream of PB bytes
*/

package parse_pb

import (
	"fmt"
)

type PBInt struct {
	value int
}

func NewPBInt(value int) PBInt {
	return PBInt{value: value}
}

func (item PBInt) String() string {
	return fmt.Sprintf("PB_INT(%d)", item.value)
}

func (item PBInt) Type() byte {
	return PB_INT
}
