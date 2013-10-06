/* packed_float_test.go
   package parse_pb

   test marshaling and unmarshaling PB_NEG
*/
package parse_pb

import (
   	"bytes"
	"testing"
)

func TestPBString(t *testing.T) {
	var testValues = []string{"", "a", "ab", "abcdefghijklmnopqrstuvwxyz"}
	for n, testValue := range testValues {
		var buffer bytes.Buffer
		pbString := NewPBString(testValue)
		if err := pbString.Marshal(&buffer); err != nil {
			t.Errorf("#%d error pbString.Marshal %s %s",
				n+1, pbString.String(), err)
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
		resultString, ok := result.(PBString)
		if !ok {
			t.Errorf("#%d unable to convert result %s", n+1, result.String)
			continue
		}
		if string(resultString.value) != testValue {
			t.Errorf("#%d expecting %d found %s",
				n+1, testValue, result.String())
		}	
	}
}
