package number

type number []uint64

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
