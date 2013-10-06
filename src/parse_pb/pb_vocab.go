/* package parse_pb

This package interprets a stream of PB bytes
*/

package parse_pb

import (
	"fmt"
	"io"
)

const (
	// Jelly Data Types
	VocabNone          = 1
	VocabClass         = 2
	VocabDereference   = 3
	VocabReference     = 4
	VocabDictionary    = 5
	VocabFunction      = 6
	VocabInstance      = 7
	VocabList          = 8
	VocabModule        = 9
	VocabPersistent    = 10
	VocabTuple         = 11
	VocabUnpersistable = 12

	// PB Data Types
	VocabCopy   = 13
	VocabCache  = 14
	VocabCached = 15
	VocabRemote = 16
	VocabLocal  = 17
	VocabLcache = 18

	// PB Protocol Messages
	VocabVersion      = 19
	VocabLogin        = 20
	VocabPassword     = 21
	VocabChallenge    = 22
	VocabLogged_in    = 23
	VocabNotLoggedIn  = 24
	VocabCachemessage = 25
	VocabMessage      = 26
	VocabAnswer       = 27
	VocabError        = 28
	VocabDecref       = 29
	VocabDecache      = 30
	VocabUncache      = 31
)

var (
	pb_vocabulary = map[int]string{
		VocabNone:          "None",
		VocabClass:         "Class",
		VocabDereference:   "Dereference",
		VocabReference:     "Reference",
		VocabDictionary:    "Dictionary",
		VocabFunction:      "Function",
		VocabInstance:      "Instance",
		VocabList:          "List",
		VocabModule:        "Module",
		VocabPersistent:    "Persistent",
		VocabTuple:         "Tuple",
		VocabUnpersistable: "Unpersistable",

		// PB Data Types
		VocabCopy:   "Copy",
		VocabCache:  "Cache",
		VocabCached: "Cached",
		VocabRemote: "Remote",
		VocabLocal:  "Local",
		VocabLcache: "Lcache",

		// PB Protocol Messages
		VocabVersion:      "Version",
		VocabLogin:        "Login",
		VocabPassword:     "Password",
		VocabChallenge:    "Challenge",
		VocabLogged_in:    "Logged_in",
		VocabNotLoggedIn:  "NotLoggedIn",
		VocabCachemessage: "Cachemessage",
		VocabMessage:      "Message",
		VocabAnswer:       "Answer",
		VocabError:        "Error",
		VocabDecref:       "Decref",
		VocabDecache:      "Decache",
		VocabUncache:      "Uncache"}
)

type PBVocab struct {
	value int
}

func NewPBVocab(vocab int) PBVocab {
	return PBVocab{value: vocab}
}

func UnmarshalPBVocab(intBuffer []byte) (PBVocab, error) {
	return NewPBVocab(unmarshalBase128Int(intBuffer)), nil
}

func (item PBVocab) Type() byte {
	return PB_VOCAB
}

func (item PBVocab) String() string {
	name, ok := pb_vocabulary[item.value]
	if !ok {
		return fmt.Sprintf("PB_VOCAB(%d)", item.value)
	}
	return fmt.Sprintf("PB_VOCAB(%s)", name)
}

func (item PBVocab) Marshal(writer io.Writer) error {
	marshaledVocab, err := marshalBase128Int(item.value) 
	if err != nil {
		return err
	}
	if _, err := writer.Write(marshaledVocab); err != nil {
		return err
	}
	if _, err := writer.Write([]byte{PB_VOCAB}); err != nil {
		return err
	}
	return nil
}
