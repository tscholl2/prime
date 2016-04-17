package prime

import (
	"math"
	"math/big"
)

// attempt to find a number within b bits of
// y^(-1/k)
func nrootbyk(y *big.Float, k, b int64) *big.Float {
	// TODO
	return nil
}

func algB(y *big.Float, k, b int64) *big.Float {
	if b < 1 || b > math.Log(8*float64(k)) {
		panic("no")
	}
	// initialization stuff
	i, _ := y.Int(nil)
	g := i.BitLen()
	a := -int64(g) / k
	B := int64(math.Ceil(math.Log(66 * float64(2*k+1))))
	// 1. set z <- 2^a + 2^(a-1), j <- 1
	z := new(big.Int).Add(new(big.Int).Lsh(one, a), new(big.Int).Lsh(one, a-1))
	j := 1
}

// IsKthPower determins if n = a^k for some integer a.
// It is written for positive integers >= 2.
// Returns true if there is such an a.
func IsKthPower(n *big.Int, k int64) bool {
	if n.BitLen() == 1 || k < 2 {
		return false
	}
	z := new(big.Int)
	// approximate 1/n
	y := new(big.Float).Quo(new(big.Float).SetInt64(1), new(big.Float).SetInt(n))
	// set f = floor(log(2n))
	f := int64(math.Log(float64(n.BitLen()+1)) / math.Log(2))
	// set b = 3 + ceil(f/k)
	b := 3 + int64(math.Ceil(float64(f)/float64(k)))
	// compute r = nroot_b(y,k)

	// find x with |r - x| <= 5/8

	// compute sign of n - x^k

	// return n == x^k
	return false
}
