/* packed_float_test.go
   package parse_pb

   test marshaling and unmarshaling PB_NEG
*/
package parse_pb

import (
   	"bytes"
	"testing"
)

func TestPBNeg(t *testing.T) {
	var testValues = []int{-1, -4, -1000, -1001, -1024, -1024*1024, 
		-1024*1024*1024}
	for n, testValue := range testValues {
		var buffer bytes.Buffer
		pbNeg := NewPBNeg(testValue)
		if err := pbNeg.Marshal(&buffer); err != nil {
			t.Errorf("#%d error pbNeg.Marshal %s %s",
				n+1, pbNeg.String(), err)
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
		resultNeg, ok := result.(PBNeg)
		if !ok {
			t.Errorf("#%d unable to convert result %s", n+1, result.String)
			continue
		}
		if resultNeg.value != testValue {
			t.Errorf("#%d expecting %d found %s",
				n+1, testValue, result.String())
		}	
	}
}
