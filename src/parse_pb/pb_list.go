/* package parse_pb

This package interprets a stream of PB bytes
*/
package parse_pb

import (
	"fmt"
	"io"
	"strings"
)

type PBList struct {
	value []parseItem
}

func NewPBList(items ...parseItem) PBList {
	return PBList{value: items}
}

func (item PBList) Type() byte {
	return PB_LIST
}

func (item PBList) String() string {
	var printValues []string
	for _, x := range item.value {
		printValues = append(printValues, x.String())
	}
	return fmt.Sprintf("PB_LIST(%s)", strings.Join(printValues, ","))
}

func (item PBList) Marshall(writer io.Writer) error {
	marshalledLen, err := marshallBase128Int(len(item.value)) 
	if err != nil {
		return err
	}
	if _, err := writer.Write(marshalledLen); err != nil {
		return err
	}
	if _, err := writer.Write([]byte{PB_LIST}); err != nil {
		return err
	}
	for i := len(item.value)-1; i >= 0; i = i - 1 {
		if err := item.value[i].Marshall(writer); err != nil {
			return err
		} 
	}
	return nil
}
