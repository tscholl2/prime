package float

import (
	"math"
	"math/big"
)

var (
	zero = big.NewInt(0)
	one  = big.NewInt(1)
)

// floating point numbers
// this struct represents n * 2^a
// note (n,a) is equivalent to (2n,a-1).
type fpn struct {
	n *big.Int
	a int
}

func (f *fpn) normalize() *fpn {
	s := trailingZeroBits(f.n)
	f.a = f.a + int(s)
	f.n.Rsh(f.n, s)
	return f
}

func divb(n *big.Int, a int, k *big.Int, b uint) *fpn {
	g := n.BitLen() - k.BitLen() - int(b)
	floor := new(big.Int)
	switch {
	case g > 0:
		floor.Div(n, floor.Lsh(k, uint(g)))
	case g < 0:
		floor.Div(floor.Lsh(n, uint(-g)), k)
	default:
		floor.Div(n, k)
	}
	return &fpn{n: floor, a: g + a}
}

// attempt to find a number within b bits of
// y^(-1/k)
func nrootbyk(y *big.Float, k, b int64) *big.Float {
	// TODO
	return nil
}

func algB(y *big.Int, k, b int) *big.Float {
	if b < 1 || b > int(math.Log(8*float64(k))) {
		panic("no")
	}
	// initialization stuff
	g := y.BitLen()
	a := -g / k
	//B := int64(math.Ceil(math.Log(66 * float64(2*k+1))))
	// 1. set z <- 2^a + 2^(a-1), j <- 1
	z := new(big.Int).Add(new(big.Int).Lsh(one, uint(a)), new(big.Int).Lsh(one, uint(a-1)))
	j := 1
	// 2. if j = b stop
	if j == b {
		return new(big.Float).SetInt(z)
	}
	// 3. compute r <- truncB(powB(z,k),truncB(y))
	r := 0.5
	// 4.
	if r < 993.0/1024.0 {
		z = z.Add(z, new(big.Int).Lsh(one, uint(a-j-1)))
	}

	//TODO
	return nil
}

/*

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

*/

func trailingZeroBits(x *big.Int) (i uint) {
	if x.Sign() < 0 {
		panic("unknown bits of negative")
	}
	if x.Sign() == 0 || x.Bit(0) == 1 {
		return 0
	}
	for i = 1; i < uint(x.BitLen()) && x.Bit(int(i)) != 1; i++ {
	}
	return
}
