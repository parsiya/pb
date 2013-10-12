/* packed_float_test.go
   package parse_pb

   test marshaling and unmarshaling PB_INT
*/
package parse_pb

import (
   	"bytes"
	"testing"
)

func TestPBInt(t *testing.T) {
	var testValues = []int{0, 1, 4, 1000, 1001, 1024, 1024*1024, 1024*1024*1024}
	for n, testValue := range testValues {
		var buffer bytes.Buffer
		pbInt := NewPBInt(testValue)
		if err := pbInt.Marshal(&buffer); err != nil {
			t.Errorf("#%d error pbInt.Marshal %s %s",
				n+1, pbInt.String(), err)
			continue
		}
		parser, err := NewParser(&buffer)
		if err != nil {
			t.Errorf("#%d error NewParser %s", n+1, err)
			continue
		}
		result, err := parser.Step()
		if err != nil {
			t.Errorf("#%d error parser.Step %s", n+1, err)
			continue
		}	
		resultInt, ok := result.(PBInt)
		if !ok {
			t.Errorf("#%d unable to convert result %s", n+1, result.String)
			continue
		}
		if resultInt.Value != testValue {
			t.Errorf("#%d expecting %d found %s",
				n+1, testValue, result.String())
		}	
	}
}
