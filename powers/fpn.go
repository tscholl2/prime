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
