/* packed_float_test.go
   package parse_pb

   test marshaling and unmarshaling PB_LIST
*/
package parse_pb

import (
   	"bytes"
	"testing"
)

func TestPBList(t *testing.T) {
	pbString := NewPBString("pb")
	pbInt := NewPBInt(42)
	pbNeg := NewPBNeg(-666)
	pbFloat := NewPBFloat(3.14159)
	pbVocab := NewPBVocab(VocabMessage)
	pbList := NewPBList(pbString, pbInt, pbNeg, pbFloat, pbVocab)

	var buffer bytes.Buffer
	if err := pbList.Marshal(&buffer); err != nil {
		t.Errorf("error pbList.Marshal %s %s", pbList.String(), err)
	}
	parser, err := NewParser(&buffer)
	if err != nil {
		t.Errorf("error NewParser %s", err)
	}
	result, err := parser.Step()
	if err != nil {
		t.Errorf("error parser.Step %s", err)
	}	
	resultList, ok := result.(PBList)
	if !ok {
		t.Errorf("unable to convert result %s", result.String)
	}

	// 2013-10-06 dougfort: kind of a cheesy equality test, 
	// but it seems ok to me (and easy)
	if resultList.String() != pbList.String() {
		t.Errorf("expecting %s found %s",
			pbList.String(), resultList.String())
	}	
}
