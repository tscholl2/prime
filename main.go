package main

import (
	"bytes"
	"encoding/ascii85"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tscholl2/prime/prime"
)

func main() {
	flag.CommandLine.Usage = func() {
		fmt.Println(`prime: generate a prime number and print to stdout
Example: 'prime -f 16 -b 128' prints: 83b19881300529d1fd4dac680415c60f
Example: 'prime -f 64 -b 256' prints: zUwiK96O4sy6pm3LtQM5YtRP9L4RGxsU/zCNliZZXn0=
Example: 'prime -f  0 -b 1024 > p.bytes' saves the raw bytes to the file 'p.bytes'
Options:`)
		flag.CommandLine.PrintDefaults()
	}
	var b, f int
	flag.IntVar(&b, "b", 128, "number of bits [supports: 2,...,128,...]")
	flag.IntVar(&f, "f", 10, "format of output [supports: 0,2-36,64,85]")
	flag.Parse()
	if b <= 1 {
		log.Fatalf("bits must be positive integer > 1, not %d", b)
	}
	p := prime.RandPrime(b)
	var s string
	switch {
	case f == 0:
		os.Stdout.Write(p.Bytes())
		return
	case 2 <= f && f <= 36:
		s = p.Text(f)
	case f == 64:
		s = string(base64.StdEncoding.EncodeToString(p.Bytes()))
	case f == 85:
		buf := bytes.NewBuffer(nil)
		enc := ascii85.NewEncoder(buf)
		enc.Write(p.Bytes())
		enc.Close()
		s = string(buf.Bytes())
	default:
		log.Fatal("unknown base")
	}
	fmt.Println(s)
}
