/* package parse_pb

This package interprets a stream of PB bytes
*/

package parse_pb

import (
	"fmt"
)

type PBString struct {
	value []byte
}

func (item PBString) String() string {
	return fmt.Sprintf("PB_STRING(%q)", item.value)
}

func (item PBString) Type() byte {
	return PB_STRING
}
