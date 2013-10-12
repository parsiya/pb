/* packed_float_test.go
   package parse_pb

   test marshaling and unmarshaling PB_VOCAB
*/
package parse_pb

import (
   	"bytes"
	"testing"
)

func TestPBVocab(t *testing.T) {
	var testValues = []int{VocabNone, VocabVersion, VocabCopy, VocabPassword, 
		VocabMessage, VocabAnswer, VocabUncache}
	for n, testValue := range testValues {
		var buffer bytes.Buffer
		pbVocab := NewPBVocab(testValue)
		if err := pbVocab.Marshal(&buffer); err != nil {
			t.Errorf("#%d error pbVocab.Marshal %s %s",
				n+1, pbVocab.String(), err)
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
		resultVocab, ok := result.(PBVocab)
		if !ok {
			t.Errorf("#%d unable to convert result %s", n+1, result.String)
			continue
		}
		if resultVocab.Value != testValue {
			t.Errorf("#%d expecting %d found %s",
				n+1, testValue, result.String())
		}	
	}
}
