/* package parse_pb

This package interprets a stream of PB bytes
*/

package parse_pb

import (
	"fmt"
)

type PBFloat struct {
	value float64
}

func (item PBFloat) String() string {
	return fmt.Sprintf("PB_FLOAT(%f)", item.value)
}

func (item PBFloat) Type() byte {
	return PB_FLOAT
}
