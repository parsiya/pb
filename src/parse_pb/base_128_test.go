/* base_128_test.go
   package parse_pb

   test parsing and generating twisted numeric marshalling
*/
package parse_pb

import (
	"bytes"
	"testing"
)

type testItem struct {
	Number        int
	MarshalledInt []byte
}

var (
	testItems = []testItem{
		testItem{0, []byte{0}},
		testItem{1, []byte{1}},
		testItem{2, []byte{2}},
		testItem{3, []byte{3}},
		testItem{4, []byte{4}},
		testItem{5, []byte{5}},
		testItem{6, []byte{6}},
		testItem{7, []byte{7}},
		testItem{88, []byte{88}},
		testItem{89, []byte{89}},
		testItem{100, []byte{100}},
		testItem{128, []byte{0, 1}},
		testItem{256, []byte{0, 2}},
		testItem{1073741824, []byte{0, 0, 0, 0, 4}},
		testItem{807421171, []byte{115, 9, 1, 1, 3}},
		testItem{324850306, []byte{2, 37, 115, 26, 1}},
		testItem{229367883, []byte{75, 64, 47, 109}},
		testItem{379260387, []byte{99, 27, 108, 52, 1}},
		testItem{464119173, []byte{5, 75, 39, 93, 1}},
		testItem{558673591, []byte{55, 93, 50, 10, 2}},
		testItem{655777377, []byte{97, 60, 89, 56, 2}},
		testItem{200334256, []byte{48, 55, 67, 95}},
		testItem{981368999, []byte{39, 1, 122, 83, 3}},
		testItem{732654087, []byte{7, 84, 45, 93, 2}}}
)

func TestMarshallInt(t *testing.T) {

	for n, testItem := range testItems {
		result, err := marshallBase128Int(testItem.Number)
		if err != nil {
			t.Errorf("#%d error marshallBase128Int(%v) %s",
				n+1, testItem.Number, err)
		}
		if bytes.Compare(result, testItem.MarshalledInt) != 0 {
			t.Errorf("#%d marshallBase128Int(%v) = %v, want %v",
				n+1, testItem.Number, result, testItem.MarshalledInt)
		}
	}
}

func TestUnmarshallInt(t *testing.T) {

	for n, testItem := range testItems {
		result := unmarshallBase128Int(testItem.MarshalledInt)
		if result != testItem.Number {
			t.Errorf("#%d unmarshallBase128Int(%v) = %v, want %v",
				n+1, testItem.MarshalledInt, result, testItem.Number)
		}
	}
}
