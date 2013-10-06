/* packed_float_test.go
   package parse_pb

   test marshaling and unmarshaling PB_FLOAT
*/
package parse_pb

import (
   	"bytes"
	"testing"
)

func TestPBFloat(t *testing.T) {
	var testValues = []float64{0.0, -1.0, 1/4, -1/1000, 1/1001, -1/1024, 
		-1/1024*1024, 1/1024*1024*1024}
	for n, testValue := range testValues {
		var buffer bytes.Buffer
		pbFloat := NewPBFloat(testValue)
		if err := pbFloat.Marshal(&buffer); err != nil {
			t.Errorf("#%d error pbFloat.Marshal %s %s",
				n+1, pbFloat.String(), err)
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
		resultFloat, ok := result.(PBFloat)
		if !ok {
			t.Errorf("#%d unable to convert result %s", n+1, result.String)
			continue
		}
		if resultFloat.value != testValue {
			t.Errorf("#%d expecting %d found %s",
				n+1, testValue, result.String())
		}	
	}
}
