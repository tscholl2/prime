package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tscholl2/goprime"
)

func main() {
	flag.CommandLine.Usage = func() {
		fmt.Println(`prime: generate a prime number and print to stdout
Example: 'prime -b 128 -f 16' prints: f72f3e72feb9cc3e030a6a5c8873dd49
Example: 'prime -b 256 -f 64' prints: /ICHR9xHQ843HisAED723qiesV1LtzxHEzSFPOPL3bc=
Example: 'prime -b 1024 -f 0 > p.bytes' saves the raw bytes to the file 'p.bytes'
Options:`)
		flag.CommandLine.PrintDefaults()
	}
	args := struct {
		Bits   uint
		Format uint
	}{}
	flag.UintVar(&args.Bits, "b", 256, "number of bits [supports: 2,...,255,...]")
	flag.UintVar(&args.Format, "f", 10, "format of output [supports: 2,10,16,64]")
	flag.Parse()
	if args.Bits <= 1 {
		log.Fatalf("error: bits must be positive integer > 1, not %d", args.Bits)
	}
	p := goprime.RandPrime(int(args.Bits))
	var s string
	switch args.Format {
	case 0:
		os.Stdout.Write(p.Bytes())
		return
	case 2:
		for p.Sign() != 0 {
			s = fmt.Sprintf("%d", p.Bits()[0]&1) + s
			p.Rsh(p, 1)
		}
	case 10:
		s = fmt.Sprintf("%d", p)
	case 16:
		s = fmt.Sprintf("%x", p)
	case 64:
		s = fmt.Sprintf("%s", base64.StdEncoding.EncodeToString(p.Bytes()))
	default:
		s = "error: unknown base"
	}
	fmt.Println(s)
}
