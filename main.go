package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"

	"github.com/tscholl2/goprime"
)

func main() {
	flag.CommandLine.Usage = func() {
		fmt.Println(`generate a prime and print to stdout
Example: 'prime -b 128 -f 16' ---> f72f3e72feb9cc3e030a6a5c8873dd49
Example: 'prime -b 256 -f 64' ---> /ICHR9xHQ843HisAED723qiesV1LtzxHEzSFPOPL3bc=
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
		log.Fatalf("error: bits must be positive integer > 1, not %d\n", args.Bits)
	}
	p := goprime.RandPrime(int(args.Bits))
	switch args.Format {
	case 2:
		s := ""
		for p.Sign() != 0 {
			s = fmt.Sprintf("%d", p.Bits()[0]&1) + s
			p.Rsh(p, 1)
		}
		fmt.Println(s)
	case 10:
		fmt.Printf("%d\n", p)
	case 16:
		fmt.Printf("%x\n", p)
	case 64:
		fmt.Printf("%s\n", base64.StdEncoding.EncodeToString(p.Bytes()))
	default:
		fmt.Printf("error: unknown base")
	}
}
