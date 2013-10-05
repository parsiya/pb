/* packed_float_test.go
   package parse_pb

   test parsing and generating twisted numeric marshalling for floating point
*/
package parse_pb

import (
	"bytes"
	"testing"
)

type testFloatItem struct {
	Number          float64
	MarshalledFloat []byte
}

var (
	testFloatItems = []testFloatItem{
		testFloatItem{-1000000.000001, []byte{193, 46, 132, 128, 0, 0, 33, 142}},
		testFloatItem{0, []byte{0, 0, 0, 0, 0, 0, 0, 0}},
		testFloatItem{1, []byte{63, 240, 0, 0, 0, 0, 0, 0}},
		testFloatItem{0.0, []byte{0, 0, 0, 0, 0, 0, 0, 0}},
		testFloatItem{1.0, []byte{63, 240, 0, 0, 0, 0, 0, 0}},
		testFloatItem{2.0, []byte{64, 0, 0, 0, 0, 0, 0, 0}},
		testFloatItem{0.48854207016041895, []byte{63, 223, 68, 69, 245, 131, 200, 230}},
		testFloatItem{0.3713929303120891, []byte{63, 215, 196, 230, 218, 105, 252, 100}},
		testFloatItem{0.5145094827339313, []byte{63, 224, 118, 220, 151, 58, 95, 31}},
		testFloatItem{0.6077756973010146, []byte{63, 227, 114, 230, 4, 230, 196, 57}},
		testFloatItem{0.4479096449865517, []byte{63, 220, 170, 141, 55, 49, 238, 246}},
		testFloatItem{0.35020940986921945, []byte{63, 214, 105, 212, 186, 136, 241, 200}},
		testFloatItem{0.6530723096639013, []byte{63, 228, 229, 247, 230, 125, 191, 158}},
		testFloatItem{0.3938815320394976, []byte{63, 217, 53, 90, 226, 166, 233, 198}},
		testFloatItem{0.7054660635505499, []byte{63, 230, 147, 45, 144, 236, 102, 38}},
		testFloatItem{0.0915535444595782, []byte{63, 183, 112, 13, 151, 73, 105, 200}}}
)

func TestMarshallPackedFloat(t *testing.T) {

	for n, testFloatItem := range testFloatItems {
		result := marshallPackedFloat(testFloatItem.Number)
		if bytes.Compare(result, testFloatItem.MarshalledFloat) != 0 {
			t.Errorf("#%d marshallPackedFloat(%v) = %v, want %v",
				n+1, testFloatItem.Number, result, testFloatItem.MarshalledFloat)
		}
	}
}

func TestUnmarshallPackedFloat(t *testing.T) {

	for n, testFloatItem := range testFloatItems {
		result := unmarshallPackedFloat(testFloatItem.MarshalledFloat)
		if result != testFloatItem.Number {
			t.Errorf("#%d unmarshallPackedFloat(%v) = %v, want %v",
				n+1, testFloatItem.MarshalledFloat, result, testFloatItem.Number)
		}
	}
}
