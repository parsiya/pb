/* package parse_pb

This package interprets a stream of PB bytes
*/

package parse_pb

import (
	"fmt"
	"io"
)

type PBVocab struct {
	value int
}

func NewPBVocab(vocab int) PBVocab {
	return PBVocab{value: vocab}
}

func UnmarshalPBVocab(intBuffer []byte) (PBVocab, error) {
	return NewPBVocab(unmarshalBase128Int(intBuffer)), nil
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

func (item PBVocab) Marshal(writer io.Writer) error {
	marshaledVocab, err := marshalBase128Int(item.value)
	if err != nil {
		return err
	}
	if _, err := writer.Write(marshaledVocab); err != nil {
		return err
	}
	if _, err := writer.Write([]byte{PB_VOCAB}); err != nil {
		return err
	}
	return nil
}
