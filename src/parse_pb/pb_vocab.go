/* package parse_pb

This package interprets a stream of PB bytes
*/

package parse_pb

import (
	"fmt"
)

type PBVocab struct {
	value int
}

func (item PBVocab) Type() byte {
	return PB_VOCAB
}

func (item PBVocab) String() string {
	name, ok := pb_vocabulary[item.value]
	if !ok {
		return fmt.Sprintf("PB_VOCAB(%d)", item.value)
	}
	return fmt.Sprintf("PB_VOCAB(%s)", name)
}
