package powers

import (
	"math"
	"math/big"
)

var (
	one = big.NewInt(1)
)

// return 1 +- 2^b
func onePlusBit(b int, sign int) *fpn {
	// TODO
	panic("unimplimeneted")
}

// returns ceil(log_2(k))
func logCeil(k uint) int {
	return int(math.Ceil((math.Log2(float64(k)))))
}

// return number of zero bits on the right of |x|
func trailingZeroBits(x *big.Int) (i uint) {
	if x.Sign() < 0 {
		// TODO optimize?
		return trailingZeroBits(new(big.Int).Neg(x))
	}
	if x.Sign() == 0 || x.Bit(0) == 1 {
		return 0
	}
	for i = 1; i < uint(x.BitLen()) && x.Bit(int(i)) != 1; i++ {
	}
	return
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
