/* package parse_pb

This package interprets a stream of PB bytes
*/
package parse_pb

import (
	"fmt"
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
