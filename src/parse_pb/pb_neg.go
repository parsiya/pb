/* package parse_pb

This package interprets a stream of PB bytes
*/

package parse_pb

import (
	"fmt"
)

type PBNeg struct {
	value int
}

func NewPBNeg(value int) PBNeg {
	return PBNeg{value: value}
}

func (item PBNeg) String() string {
	return fmt.Sprintf("PB_NEG(%d)", item.value)
}

func (item PBNeg) Type() byte {
	return PB_NEG
}
