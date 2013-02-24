// goembed generates a Go source file from an input file.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var (
	packageFlag = flag.String("package", "", "Go package name")
	varFlag     = flag.String("var", "", "Go var name")
)

func main() {
	flag.Parse()

	raw, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Reading stdin: %v", err)
	}

	fmt.Printf("package %s\n\n", *packageFlag)
	fmt.Printf("var %s = []byte{ // %d bytes\n", *varFlag, len(raw))

	const perLine = 16
	for len(raw) > 0 {
		n := perLine
		if n > len(raw) {
			n = len(raw)
		}
		for _, b := range raw[:n] {
			fmt.Printf(" 0x%02x,", b)
		}
		fmt.Print("\n")
		raw = raw[n:]
	}
	fmt.Print("}\n")
}
