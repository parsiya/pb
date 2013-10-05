/* program dump_pb
   package main

   This program reads raw pb bytes from stdin and writes a formatted dump 
   to stdout
*/
package main

import (
   "parse_pb"
   "io"
   "log"
   "os"
)


// main entry point for ergo_prpoxy
func main() {
	log.Printf("INFO: program starts")
	parser, err := parse_pb.NewParser(os.Stdin)
	if err != nil {
		log.Fatalf("CRITICAL: parse_pb.NewParser failed %s", err)
	}
	for {
		result, err := parser.Step()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("CRITICAL: error on parser.Step %s", err)
		}
		log.Printf("DEBUG: %s", result)
	}
	log.Printf("INFO: program ends")
}
