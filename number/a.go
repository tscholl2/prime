package main

import (
	"fmt"
	"math"
)

type number []uint64

// left shift by i <= 63
func (z *number) lsh63(x *number, s uint8) *number {
	var bot, top uint64
	mask := (((uint64(1) << s) - 1) << (63 - s))
	if z == nil || cap(*z) < len(*x)+1 {
		z = new(number)
		*z = make(number, len(*x)+1)
	}
	*z = (*z)[:len(*x)+1]
	for i, xi := range *x {
		top = xi & mask // top s bits
		fmt.Printf("top = %b\n", top)
		fmt.Printf("bot = %b\n", bot)
		(*z)[i] = ((xi << s) | bot) & 0x7fffffffffffffff
		bot = top >> (63 - s) // bottom s bits
	}
	if bot != 0 {
		*z = append(*z, bot)
	}
	if (*z)[len(*z)-1] == 0 {
		*z = (*z)[:len(*z)-1]
	}
	return z
}

func main() {
	x := new(number)
	y := &number{0, math.MaxUint64 >> 2, math.MaxUint64 >> 1, 0, 1}
	z := x.lsh63(y, 3)
	for _, a := range *z {
		fmt.Printf("%b\n", a)
	}
}
