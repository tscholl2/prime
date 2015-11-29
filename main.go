package main

import (
	"encoding/base64"
	"fmt"

	"github.com/alexflint/go-arg"
	"github.com/tscholl2/goprime"
)

func main() {
	args := struct {
		Bits    uint `arg:"-b,help:number of bits"`
		Format  uint `arg:"-f,help:give output in base f"`
		Newline bool `arg:"-n,help:set to remove newline"`
	}{256, 10, false}
	arg.MustParse(&args)
	if args.Bits == 0 {
		fmt.Println("error: bits must be positive integer")
		return
	}
	p := goprime.RandPrime(int(args.Bits))
	switch args.Format {
	case 10:
		fmt.Printf("%d", p)
	case 16:
		fmt.Printf("%x", p)
	case 64:
		fmt.Printf("%s", base64.StdEncoding.EncodeToString(p.Bytes()))
	default:
		fmt.Printf("error: unknown base")
	}
	if !args.Newline {
		fmt.Printf("\n")
	}
}
