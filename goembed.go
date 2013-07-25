// goembed generates a Go source file from an input file.
package main

import (
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

var (
	packageFlag = flag.String("package", "", "Go package name")
	varFlag     = flag.String("var", "", "Go var name")
	gzipFlag    = flag.Bool("gzip", false, "Whether to gzip contents")
)

func main() {
	flag.Parse()

	raw, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Reading stdin: %v", err)
	}

	fmt.Printf("package %s\n\n", *packageFlag)

	if !*gzipFlag {
		fmt.Printf("var %s = []byte{ // %d bytes\n", *varFlag, len(raw))
	} else {
		var buf bytes.Buffer
		gzw, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
		if _, err := gzw.Write(raw); err != nil {
			log.Fatal(err)
		}
		if err := gzw.Close(); err != nil {
			log.Fatal(err)
		}
		gz := buf.Bytes()

		if err := gzipPrologue.Execute(os.Stdout, *varFlag); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("var %s []byte // set in init\n\n", *varFlag)
		fmt.Printf("var %s_gzip = []byte{ // %d compressed bytes (%d uncompressed bytes)\n", *varFlag, len(gz), len(raw))
		raw = gz
	}

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

var gzipPrologue = template.Must(template.New("").Parse(`
import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func init() {
	r, err := gzip.NewReader(bytes.NewReader({{.}}_gzip))
	if err != nil {
		panic(err)
	}
	defer r.Close()
	{{.}}, err = ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
}
`))
