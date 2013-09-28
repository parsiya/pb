/* package parse_pb

This package interprets a stream of PB bytes
*/
package parse_pb

/*
the "Banana" protocol works by serializing lists.
look at banana.py line 150 (dataReceived).  define same constants.
*/
import (
	"bufio"
   	"fmt"
   	"io"
   	"strings"
)

const (
	PB_LIST      byte = 128
	PB_INT       byte = 129
	PB_STRING    byte = 130
	PB_NEG       byte = 131
	PB_FLOAT     byte = 132
	PB_LONGINT   byte = 133
	PB_LONGNEG   byte = 134
	PB_VOCAB     byte = 135
	HIGH_BIT_SET byte = 128
)

var (
	pb_constants = map[byte]string{
		128 : "PB_LIST",
		129 : "PB_INT",
		130 : "PB_STRING",
		131 : "PB_NEG",
		132 : "PB_FLOAT",
		133 : "PB_LONGINT",
		134 : "PB_LONGNEG",
		135 : "PB_VOCAB"}
)

type parseItem interface {
	Type() byte
	String() string
}

type PBList struct {
	value []parseItem
}

func (item PBList) Type() byte {
	return PB_LIST
}

func (item PBList) String() string {
	var printValues []string 
	for _, x := range item.value {
		printValues = append(printValues, x.String()) 
	}
	return fmt.Sprintf("PB_LIST(%s)", strings.Join(printValues, ","))
} 

type PBInt struct {
	value int
}

func (item PBInt) String() string {
	return fmt.Sprintf("PB_INT(%d)", item.value)
} 

func (item PBInt) Type() byte {
	return PB_INT
}

type PBString struct {
	value string
}

func (item PBString) String() string {
	return fmt.Sprintf("PB_STRING(%s)", item.value)
} 

func (item PBString) Type() byte {
	return PB_STRING
}

type PBNeg struct {
	value int
}

func (item PBNeg) String() string {
	return fmt.Sprintf("PB_NEG(%d)", item.value)
} 

func (item PBNeg) Type() byte {
	return PB_NEG
}

type PBFloat struct {
	value float64
}

func (item PBFloat) String() string {
	return fmt.Sprintf("PB_FLOAT(%f)", item.value)
} 

func (item PBFloat) Type() byte {
	return PB_FLOAT
}

type PBVocab struct {
	value int
}

func (item PBVocab) Type() byte {
	return PB_VOCAB
}

func (item PBVocab) String() string {
	return fmt.Sprintf("PB_VOCAB(%d)", item.value)
} 

type PBUnknown struct {
	intBuffer []byte
	c byte
}

func (item PBUnknown) String() string {
	return fmt.Sprintf("* %v %s *", item.intBuffer, dumpByte(item.c))
}

func (item PBUnknown) Type() byte {
	return 0
}

type Parser struct {
	reader *bufio.Reader
}

func NewParser(reader io.Reader) (*Parser, error) {
	parser := Parser{reader : bufio.NewReader(reader)}
	return &parser, nil
}

func (parser *Parser) Step() (parseItem, error) {
	var b byte
	var intBuffer []byte
	var c byte
	var err error
	for {
		if b, err = parser.reader.ReadByte(); err != nil {
			return nil, err
		}
		if (b & HIGH_BIT_SET) != 0 {
			c = b
			break
		}
		intBuffer = append(intBuffer, b)
	}
	switch c {
	case PB_LIST:
		return parser.parseList(parseBase128Int(intBuffer), c)
	case PB_INT:
		return parser.parseInt(parseBase128Int(intBuffer), c)
	case PB_STRING:
		return parser.parseString(parseBase128Int(intBuffer), c)
	case PB_NEG:
		return parser.parseNeg(-parseBase128Int(intBuffer), c)
	case PB_FLOAT:
		buffer := make([]byte, 8)
		if _, err := parser.reader.Read(buffer); err != nil {
			return nil, err
		}
		return parser.parseFloat(c, buffer)
	case PB_VOCAB:
		return parser.parseVocab(parseBase128Int(intBuffer), c)
	}
	return parser.parseUnknown(intBuffer, c)
}

func (parser *Parser) parseList(i int, _ byte) (parseItem, error) {
	size := i
	list := PBList{value : make([]parseItem, size)}
	for j := 0; j < size; j++ {
		value, err := parser.Step()
		if err != nil {
			return nil, err
		}
		list.value[j] = value 
	}

	return list, nil
}

func (parser *Parser) parseInt(i int, _ byte) (parseItem, error) {
	return PBInt{value : i}, nil
}

func (parser *Parser) parseString(i int, _ byte) (parseItem, error) {
	size := i
	buffer := make([]byte, size)
	n, err := parser.reader.Read(buffer)
	if err != nil {
		return nil, err
	}
	if n != size {
		return nil, 
			fmt.Errorf("parserString read %d bytes expecting %d", n, size)
	}
	return PBString{value : string(buffer)}, nil
}

func (parser *Parser) parseNeg(i int, _ byte) (parseItem, error) {
	return PBNeg{value : i}, nil
}

func (parser *Parser) parseFloat(_ byte, _ []byte) (parseItem, error) {
	return PBFloat{value : 0.9}, nil
}

func (parser *Parser) parseVocab(i int, _ byte) (parseItem, error) {
	return PBVocab{value : i}, nil
}

func (parser *Parser) parseUnknown(intBuffer []byte, c byte) (
	parseItem, error) {
	return PBUnknown{intBuffer : intBuffer, c : c}, nil
}

// dumpByte dumps one byte
func dumpByte(c byte) string {
	constant, ok := pb_constants[c]
	if ok {
		return constant
	}
	return fmt.Sprintf("%d", c)
}