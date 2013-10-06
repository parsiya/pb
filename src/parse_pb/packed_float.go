/* packed_float.go
   package parse_pb

   parsing and generating twisted numeric marshaling (packed float)
*/
package parse_pb

import (
	"math"
)

const (
	PackedFloatSliceSize = 8
)

// marshalPackedFloat converts 'number' into a slice of 8 bytes, which form
// the representation of the number: see IEEE 754
func marshalPackedFloat(number float64) []byte {
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

// unmarshalPackedFloat converts a slice of bytes to a float
// see see IEEE 754
func unmarshalPackedFloat(marshaledFloat []byte) float64 {
	var uintVal uint64

	for _, b := range marshaledFloat {
		uintVal = uintVal << 8
		uintVal = uintVal | uint64(b)
	}

	return math.Float64frombits(uintVal)
}
