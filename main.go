package main

import (
	"encoding/base64"
	"fmt"

	"github.com/alexflint/go-arg"
	"github.com/tscholl2/goprime"
)

func main() {
	args := struct {
		Bits uint `arg:"-b,help:number of bits"`
		Base uint `arg:"-f,help:format output base f"`
	}{256, 10}
	arg.MustParse(&args)
	p := goprime.RandPrime(int(args.Bits))
	switch args.Base {
	case 10:
		fmt.Printf("%d", p)
	case 16:
		fmt.Printf("%x", p)
	case 64:
		fmt.Printf("%s", base64.StdEncoding.EncodeToString(p.Bytes()))
	default:
		fmt.Printf("unknown base")
	}
}
