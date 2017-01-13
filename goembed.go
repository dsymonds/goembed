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

	// Generate []byte(<big string constant>) instead of []byte{<list of byte values>}.
	// The latter causes a memory explosion in the compiler (60 MB of input chews over 9 GB RAM).
	// Doing a string conversion avoids some of that, but incurs a slight startup cost.
	if !*gzipFlag {
		fmt.Printf(`var %s = []byte("`, *varFlag)
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
		fmt.Printf(`var %s_gzip = []byte("`, *varFlag)
		raw = gz
	}

	for _, b := range raw {
		fmt.Printf("\\x%02x", b)
	}
	fmt.Println(`")`)
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
