package powers

import "math/big"

var (
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
	return r.n == nil || r.n.BitLen() == 0
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

func (r *fpn) sub(s, t *fpn) *fpn {
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
	r.n.Sub(s.n, t.n)
	r.a = m
	return r.normalize()
}

func (r *fpn) cmp(s *fpn) int {
	if r.n == nil {
		r.n = new(big.Int)
	}
	r.normalize()
	s.normalize()
	m := min(r.a, s.a)
	r.n.Lsh(r.n, uint(r.a-m))
	r.a = m
	s.n.Lsh(s.n, uint(s.a-m))
	s.a = m
	return r.n.Cmp(s.n)
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

// returns integer x such that |r - x| < 1
// also returns an approximation of the error |r - x| < b
func (r *fpn) nearestInt() (x *big.Int, b float64) {
	if r.isZero() {
		x.SetInt64(0)
		return
	}
	if r.a >= 0 {
		x.Lsh(r.n, uint(r.a))
		return
	}
	if r.n.Sign() == -1 {
		s := &fpn{new(big.Int), r.a}
		s.n.Neg(r.n)
		x, b = s.nearestInt()
		x.Neg(x)
		return
	}
	if r.n.BitLen() > -r.a {
		x.Rsh(r.n, uint(-r.a))
	} else {
		x.SetInt64(0)
	}
	// first 64 bits past decimal place
	shft := min(64, -r.a)
	decpart := new(big.Int)
	for i := -r.a - shft; i < -r.a; i++ {
		decpart.SetBit(decpart, -r.a-i, r.n.Bit(i))
	}
	// TODO finish
}
