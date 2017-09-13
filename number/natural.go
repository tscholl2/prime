package number

import (
	"encoding/hex"
	"math"
)

type number []uint64

func (z number) setHex(s string) number {
	arr, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	if len(arr) == 0 {
		return z[:0]
	}
	m := int(math.Ceil(8 * float64(len(arr)) / 63))
	if cap(z) < m {
		z = make(number, m)
	}
	z = z[:m]
	for _, b := range arr {
		z = z.lsh(z, 8)
		z[0] = z[0] | uint64(b)
	}
	return z.norm()
}

func (z number) lsh(x number, s uint) number {
	t := uint8(s % 63)
	blocks := (int(s) - int(t)) / 63
	w := z.lsh63(x, t)
	if cap(z) < len(w)+blocks {
		z = make(number, len(w)+blocks)
	}
	z = z[:len(w)+blocks]
	for i := 0; i < len(w); i++ {
		z[i+blocks] = w[i]
	}
	return z
}

// left shift by i <= 63
func (z number) lsh63(x number, s uint8) number {
	var bot, top uint64
	mask := (((uint64(1) << s) - 1) << (63 - s))
	if z == nil || cap(z) < len(x)+1 {
		z = make(number, len(x)+1)
	}
	z = z[:len(x)+1]
	for i, xi := range x {
		top = xi & mask // top s bits
		z[i] = ((xi << s) | bot) & 0x7fffffffffffffff
		bot = top >> (63 - s) // bottom s bits
	}
	if bot != 0 {
		z = append(z, bot)
	}
	if z[len(z)-1] == 0 {
		z = z[:len(z)-1]
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
	if cap(x) < yLength+1 {
		x = make(number, yLength+1) // TODO: use x's slice rather than new one
	}
	x = x[:yLength+1]
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

func (z number) norm() number {
	if z == nil {
		return number{}
	}
	for i := len(z) - 1; i >= 0 && z[i] == 0; i-- {
		z = z[:i]
	}
	return z
}
