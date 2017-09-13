package number

import (
	"encoding/hex"
	"fmt"
)

type number []uint64

func (x *number) setHex(s string) *number {
	arr, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	var n number
	var c uint64
	for i := 0; i < len(arr); i++ {
		var xi uint64
		for j := 0; j < 8 && i < len(arr); j++ {
			xi = xi << 8
			xi = xi | uint64(arr[i])
			i++
		}
		c = (xi & 0x8000000000000000) >> 63
		xi = xi & 0x7fffffffffffffff
		n = append(n, xi)
	}
	*x = n
	return x
}

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
		fmt.Printf("top = %b", top)
		fmt.Printf("bot = %b", bot)
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

func (x number) add(y, z number) number {
	yLength := len(y)
	zLength := len(z)
	if yLength < zLength {
		return x.add(z, y)
	}
	// yLength >= zLength
	var c uint64
	x = make([]uint64, yLength+1) // TODO: use x's slice rather than new one
	i := 0
	for ; i < zLength; i++ {
		x[i] = y[i] + z[i] + c
		c = (x[i] & 0x8000000000000000) >> 63 // top bit
		x[i] = x[i] & 0x7fffffffffffffff      // clear top bit
	}
	for ; i < yLength; i++ {
		x[i] = y[i] + c
		c = (x[i] & 0x8000000000000000) >> 63 // top bit
		x[i] = x[i] & 0x7fffffffffffffff      // clear top bit
	}
	x[i] = c
	if c == 0 {
		x = x[:i]
	}
	return x
}
