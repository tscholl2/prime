package powers

import (
	"fmt"
	"math"
	"math/big"
)

var (
	zero           = big.NewInt(0)
	one            = big.NewInt(1)
	negone         = big.NewInt(-1)
	neg993over1024 = &fpn{big.NewInt(-993), -10}
)

// floating point numbers
// this struct represents n * 2^a
// note (n,a) is equivalent to (2n,a-1).
type fpn struct {
	n *big.Int
	a int
}

func (r *fpn) normalize() *fpn {
	if r == nil || r.n == nil || r.n.Bit(0) == 1 {
		return r
	}
	i := trailingZeroBits(r.n)
	r.a = r.a + int(i)
	r.n.Rsh(r.n, i)
	return r
}

func (r *fpn) isZero() bool {
	return r == nil || r.n == nil || r.n.BitLen() == 0
}

func (r *fpn) mul(s, t *fpn) *fpn {
	if r.n == nil {
		r.n = new(big.Int)
	}
	r.n.Mul(s.n, t.n)
	r.a = s.a + t.a
	return r
}

func (r *fpn) add(s, t *fpn) *fpn {
	if r.n == nil {
		r.n = new(big.Int)
	}
	s.normalize()
	t.normalize()
	m := min(s.a, t.a)
	s.n.Lsh(s.n, uint(s.a-m))
	s.a = m
	t.n.Lsh(t.n, uint(t.a-m))
	t.a = m
	r.n.Add(s.n, t.n)
	r.a = m
	return r.normalize()
}

// returns true if r <= 1
func (r *fpn) leq1() bool {
	if r == nil || r.n == nil || r.n.Sign() <= 0 {
		return true
	}
	r.normalize()
	if r.a > 0 {
		return false
	}
	if r.a == 0 {
		return r.n.BitLen() <= 1
	}
	return r.n.BitLen()-1 < -r.a
}

// returns true if r <= 993/1024
func (r *fpn) leq993over1024() bool {
	if r == nil || r.n == nil || r.n.Sign() <= 0 {
		return true
	}
	if !r.leq1() {
		return false
	}
	return new(fpn).add(r, neg993over1024).n.Sign() <= 0
}

// returns ceil(log_2(k))
func logCeil(k uint) int {
	return int(math.Ceil((math.Log2(float64(k)))))
}

// returns s such that s <= r/k < s(1 + 2^(1 - b)) where r = n*2^a
func divb(n *big.Int, a int, k uint, b uint) *fpn {
	if k <= 0 {
		panic("no")
	}
	g := n.BitLen() - logCeil(k) - int(b)
	floor := new(big.Int)
	bk := big.NewInt(int64(k))
	switch {
	case g > 0:
		floor.Div(n, floor.Lsh(bk, uint(g)))
	case g < 0:
		floor.Div(floor.Lsh(n, uint(-g)), bk)
	default:
		floor.Div(n, bk)
	}
	return &fpn{n: floor, a: g + a}
}

// returns s such that s <= r < s(1 + 2^(1 - b)) where r = n*2^a
func truncb(n *big.Int, a int, b uint) *fpn {
	g := n.BitLen() - int(b)
	n2 := new(big.Int)
	switch {
	case g > 0:
		n2.Rsh(n, uint(g))
	case g < 0:
		n2.Lsh(n, uint(-g))
	}
	if n.BitLen() < g {
		return &fpn{new(big.Int), 0}
	}
	return &fpn{n2, a + g}
}

// returns s such that s <= r^k < s(1 + 2^(1 - b))^(2k - 1)
func powb(r *fpn, k uint, b uint) *fpn {
	if k == 0 {
		panic("no")
	}
	switch {
	case k == 1:
		return truncb(r.n, r.a, b)
	case k%2 == 0:
		s := powb(r, k>>1, b)
		if s == nil || s.n == nil {
			fmt.Println(r.n, r.a, k, b)
		}
		s.mul(s, s)
		return truncb(s.n, s.a, b)
	}
	s1 := powb(r, k-1, b)
	s2 := truncb(r.n, r.a, b)
	s3 := new(fpn).mul(s1, s2)
	return truncb(s3.n, s3.a, b)
}

func algB(y *fpn, k uint, b uint) *fpn {
	if b == 0 || int(b) > int(math.Log(8*float64(k))) || y.isZero() {
		panic("no")
	}
	// initialization stuff
	g := y.n.BitLen() + y.a                        // 2^(g-1) <= y < 2^g, works if y > 0
	a := int(math.Floor(float64(-g) / float64(k))) // floor(-g/k)
	B := uint(logCeil(66 * (2*k + 1)))             // ceil(log(66(2k+1)))
	// 1. set z = 2^a + 2^(a-1) and j = 1
	z := &fpn{big.NewInt(3), a - 1}
	for j := uint(1); j < b; j++ {
		// 2. if j = b stop
		if j == b {
			return z
		}
		// 3. compute r <- truncB(powB(z,k),truncB(y))
		pz := new(fpn).mul(powb(z, k, B), truncb(y.n, y.a, B))
		r := truncb(pz.n, pz.a, B)
		// 4. if r <= 993/1024 then set z = z + 2^{a - j - 1}
		if r.leq993over1024() {
			z.add(z, &fpn{one, a - int(j) - 1})
		}
		// 5. if r > 1 then set z = z - 2^{a - j - 1}
		if !r.leq1() {
			z.add(z, &fpn{negone, a - int(j) - 1})
		}
		// 6. set j = j + 1 and go to step 2
		// done in the for loop
	}
	panic("hmmmm")
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
