/* packed_float.go
   package parse_pb

   parsing and generating twisted numeric marshalling (packed float)
*/
package parse_pb

import (
	"math"
)

const (
	PackedFloatSliceSize = 8
)

// marshallPackedFloat converts 'number' into a slice of 8 bytes, which form
// the representation of the number: see IEEE 754
func marshallPackedFloat(number float64) []byte {
	uintVal := math.Float64bits(number)
	if uintVal == 0 {
		return []byte{0, 0, 0, 0, 0, 0, 0, 0}
	}
	result := make([]byte, PackedFloatSliceSize)
	index := PackedFloatSliceSize - 1
	for ; uintVal > 0; uintVal = uintVal >> 8 {
		b := uintVal & 0xFF
		result[index] = byte(b)
		index = index - 1
	}
	return result
}

// unmarshallPackedFloat converts a slice of bytes to a float
// see see IEEE 754
func unmarshallPackedFloat(marshalledFloat []byte) float64 {
	var uintVal uint64

	for _, b := range marshalledFloat {
		uintVal = uintVal << 8
		uintVal = uintVal | uint64(b)
	}

	return math.Float64frombits(uintVal)
}
