package powers

import (
	"math"
	"math/big"
)

// returns s such that s <= r/k < s(1 + 2^(1 - b)) where r = n*2^a
func divb(r *fpn, k uint, b uint) *fpn {
	if k <= 0 {
		panic("no")
	}
	g := r.n.BitLen() - logCeil(k) - int(b)
	floor := new(big.Int)
	bk := big.NewInt(int64(k))
	switch {
	case g > 0:
		floor.Div(r.n, floor.Lsh(bk, uint(g)))
	case g < 0:
		floor.Div(floor.Lsh(r.n, uint(-g)), bk)
	default:
		floor.Div(r.n, bk)
	}
	return &fpn{floor, g + r.a}
}

// returns s such that s <= r < s(1 + 2^(1 - b)) where r = n*2^a
func truncb(r *fpn, b uint) *fpn {
	g := r.n.BitLen() - int(b)
	n2 := new(big.Int)
	switch {
	case g > 0:
		n2.Rsh(r.n, uint(g))
	case g < 0:
		n2.Lsh(r.n, uint(-g))
	}
	if r.n.BitLen() < g {
		return &fpn{new(big.Int), 0}
	}
	return &fpn{n2, r.a + g}
}

// alg P
// given positive fpn r and positive ints k,b return powb(r,k)
func powb(r *fpn, k, b uint) *fpn {
	if r.isZero() {
		return new(fpn)
	}
	if k == 0 || b == 0 {
		panic("no")
	}
	// 1. if k = 1, return truncb(r)
	if k == 1 {
		return truncb(r, b)
	}
	// 2. if k is even, compute powb(r,k/2) by algP
	// and return truncb(powb(r,k/2)^2)
	if k%2 == 0 {
		s := powb(r, k/2, b)
		return truncb(s.mul(s, s), b)
	}
	// 3. compute powb(r,k-1) by algP and return truncb(powb(r,k-1) * truncb(r))
	s := powb(r, k-1, b)
	return truncb(s.mul(s, truncb(r, b)), b)
}

func algB(y *fpn, k, b uint) *fpn {
	if y.isZero() {
		return new(fpn)
	}
	if b < 1 || b > 3+uint(logCeil(k)) {
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
		// done out of loop
		// 3. compute r <- truncB(powB(z,k),truncB(y))
		pz := new(fpn).mul(powb(z, k, B), truncb(y, B))
		r := truncb(pz, B)
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
	return z
}

func algN(y *fpn, k, b uint) *fpn {
	if y.isZero() {
		return new(fpn)
	}
	if b < 4+uint(logCeil(k)) {
		panic("no")
	}
	// initialization stuff
	log2k := uint(logCeil(k)) + 1
	bb := log2k + uint(math.Ceil(float64(int(b)-int(log2k))/2)) // b' = ceil(lg(k)) + ceil((b - ceil(lg(k)))/2)
	B := 2*bb + 4 - log2k                                       // B = 2b' + 4 - ceil(lg(k))
	// 1. z = nrootb'(y,k) by B if b' <= lg(k) + 3, else algN
	var z *fpn
	if bb <= log2k+3 {
		z = algB(y, k, bb)
	} else {
		z = algN(y, k, bb)
	}
	// 2. set r2 = mul(truncB(z),k+1)
	var r2, r3, r4 fpn
	r2.mul(truncb(z, B), &fpn{big.NewInt(int64(k + 1)), 0})
	// 3. set r3 = truncB(powB(z,k+1) * truncB(y))
	r3 = *truncb(r3.mul(powb(z, k+1, B), truncb(y, B)), B)
	// 4. set r4 = divB(r2 - r3,k)
	r4 = *divb(r4.sub(&r2, &r3), k, B)
	return &r4
}

// given positive ints n,x,k return sign of n - x^k
func algC(n *big.Int, x *fpn, k uint) int {
	if n.Sign() <= 0 || x.isZero() || k < 1 {
		panic("no")
	}
	nf := &fpn{n, 0}
	// initialization
	f := n.BitLen() - 1 // f = floor(lg(2n))
	// 1. b = 1
	for b := 1; b < f; b = min(2*b, f) {
		// 2. r = pow_{b + ceil(lg(8k))}(x,k)
		r := powb(x, k, 3+uint(b+logCeil(k)))
		// 3. if n < r, return -1 and stop
		if nf.cmp(r) == -1 {
			return -1
		}
		// 4. if r(1 + 2^(-b)) <= n return 1
		if r.mul(r, new(fpn).add(&fpn{one, 0}, &fpn{one, -b})).cmp(nf) <= 0 {
			return 1
		}
		// 5. if b >= f return 0
		// 6. b = min(2b,f), goto step 2
	}
	return 0
}
