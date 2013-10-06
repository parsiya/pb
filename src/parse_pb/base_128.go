/* base_128.go
   package parse_pb

   parsing and generating twisted numeric marshaling (base 128)
*/
package parse_pb

import (
	"fmt"
)

// marshalBase128Int converts 'number' into a slice of bytes, which form
// the reresentation of the number in base 128
func marshalBase128Int(number int) ([]byte, error) {
	if number < 0 {
		return nil, fmt.Errorf("attempt to marshal negative number %s", number)
	}
	if number == 0 {
		return []byte{0}, nil
	}

	var result []byte
	for ; number > 0; number = number >> 7 {
		b := number & 0x7F
		result = append(result, byte(b))
	}
	return result, nil
}

// unmarshalBase128Int converts a slice of bytes to an int
// the bytes are assumed to form a representation of the number in base 128
func unmarshalBase128Int(marshaledInt []byte) int {
	var result int

	exponent := 1
	for _, b := range marshaledInt {
		result = result + (int(b) * exponent)
		exponent = exponent << 7
	}

	return result
}
