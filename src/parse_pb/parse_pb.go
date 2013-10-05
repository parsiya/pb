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
		PB_LIST:    "PB_LIST",
		PB_INT:     "PB_INT",
		PB_STRING:  "PB_STRING",
		PB_NEG:     "PB_NEG",
		PB_FLOAT:   "PB_FLOAT",
		PB_LONGINT: "PB_LONGINT",
		PB_LONGNEG: "PB_LONGNEG",
		PB_VOCAB:   "PB_VOCAB"}
)

type parseItem interface {
	Type() byte
	String() string
}

type Parser struct {
	reader *bufio.Reader
}

func NewParser(reader io.Reader) (*Parser, error) {
	parser := Parser{reader: bufio.NewReader(reader)}
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
		return parser.parseList(unmarshallBase128Int(intBuffer), c)
	case PB_INT:
		return parser.parseInt(unmarshallBase128Int(intBuffer), c)
	case PB_STRING:
		return parser.parseString(unmarshallBase128Int(intBuffer), c)
	case PB_NEG:
		return parser.parseNeg(-unmarshallBase128Int(intBuffer), c)
	case PB_FLOAT:
		return parser.parseFloat(c)
	case PB_VOCAB:
		return parser.parseVocab(unmarshallBase128Int(intBuffer), c)
	}
	return parser.parseUnknown(intBuffer, c)
}

func (parser *Parser) parseList(i int, _ byte) (parseItem, error) {
	size := i
	list := PBList{value: make([]parseItem, size)}
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
	return PBInt{value: i}, nil
}

func (parser *Parser) parseString(i int, _ byte) (parseItem, error) {
	data, err := parser.readAll(i)
	if err != nil {
		return nil, err
	}
	return PBString{value: data}, nil
}

func (parser *Parser) parseNeg(i int, _ byte) (parseItem, error) {
	return PBNeg{value: i}, nil
}

func (parser *Parser) parseFloat(_ byte) (parseItem, error) {
	marshalledFloat, err := parser.readAll(PackedFloatSliceSize)
	if err != nil {
		return nil, err
	}
	return PBFloat{value: unmarshallPackedFloat(marshalledFloat)}, nil
}

func (parser *Parser) parseVocab(i int, _ byte) (parseItem, error) {
	return PBVocab{value: i}, nil
}

func (parser *Parser) parseUnknown(intBuffer []byte, c byte) (
	parseItem, error) {
	return PBUnknown{intBuffer: intBuffer, c: c}, nil
}

// dumpByte dumps one byte
func dumpByte(c byte) string {
	constant, ok := pb_constants[c]
	if ok {
		return constant
	}
	return fmt.Sprintf("%d", c)
}

// readAll reads all the bytes needed to fill and array
func (parser *Parser) readAll(size int) ([]byte, error) {
	buffer := make([]byte, size)

	offset := 0
	bytesToRead := size
	for bytesToRead > 0 {
		n, err := parser.reader.Read(buffer[offset : offset+bytesToRead])
		if err != nil {
			return nil, err
		}
		offset = offset + n
		bytesToRead = bytesToRead - n
	}

	return buffer, nil
}
