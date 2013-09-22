/* program dump_pb
   package main

   This program reads raw pb bytes from stdin and writes a formatted dump 
   to stdout
*/
package main

import (
   "bufio"
   "fmt"
   "io"
   "log"
   "os"
)

/*
the "Banana" protocol works by serializing lists.
look at banana.py line 150 (dataReceived).  define same constants.
*/
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

// main entry point for ergo_prpoxy
func main() {
	log.Printf("INFO: program starts")
	reader := bufio.NewReader(os.Stdin)
	var c byte
	var err error
	for {
		if c, err = reader.ReadByte(); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("CRITICAL: error on ReadByte %s", err)
		}
		constant, ok := pb_constants[c]
		if ok {
			fmt.Printf("%s\n", constant)
		} else {
			fmt.Printf("%d\n", c)
		}
	}
	log.Printf("INFO: program ends")
}
