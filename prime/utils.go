package prime

import (
	"math"
	"math/big"
)

// JacobiSymbol returns the jacobi symbol ( N / D ) of
// N (numerator) over D (denominator).
// See http://en.wikipedia.org/wiki/Jacobi_symbol
func JacobiSymbol(N *big.Int, D *big.Int) int {
	//Step 0: parse input / easy cases
	if D.Sign() <= 0 || D.Bit(0) == 0 {
		// we will assume D is positive
		// wolfram is ok with negative denominator
		// im not sure what is standard though
		panic("JacobiSymbol defined for positive odd denominator only")
	}
	var n, d, tmp big.Int
	n.Set(N)
	d.Set(D)
	j := 1
	for {
		// Step 1: Reduce the numerator mod the denominator
		n.Mod(&n, &d)
		if n.Sign() == 0 {
			// if n,d not relatively prime
			return 0
		}
		if len(n.Bits()) >= len(d.Bits())-1 {
			// n > d/2 so swap n with d-n
			// and multiply j by JacobiSymbol(-1 / d)
			n.Sub(&d, &n)
			if d.Bits()[0]&3 == 3 {
				// if d = 3 mod 4
				j = -1 * j
			}
		}

		// Step 2: extract factors of 2
		s := trailingZeroBits(&n)
		n.Rsh(&n, s)
		if s&1 == 1 {
			switch d.Bits()[0] & 7 {
			case 3, 5: // d = 3,5 mod 8
				j = -1 * j
			}
		}

		// Step 3: check numerator
		if len(n.Bits()) == 1 && n.Bits()[0] == 1 {
			// if n = 1 were done
			return j
		}

		// Step 4: flip and go back to step 1
		if n.Bits()[0]&3 != 1 { // n = 3 mod 4
			if d.Bits()[0]&3 != 1 { // d = 3 mod 4
				j = -1 * j
			}
		}
		tmp.Set(&n)
		n.Set(&d)
		d.Set(&tmp)
	}
}

// counts the number of zeros at the end of the
// binary expansion. So 2=10 ---> 1, 4=100 ---> 2
// 3=111 ---> 0, see test for more examples
// also 0 ---> 0 and 1 ---> 0
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

// IsSquare returns true if N = m^2
// for some positive integer m.
// It uses newtons method and other checks.
func IsSquare(N *big.Int) bool {
	// Step -1: check inputs
	if N.Sign() <= 0 {
		// 0 is a square
		if N.Sign() == 0 {
			return true
		}
		// negative numbers are not
		return false
	}

	// Step 0: Easy case
	if N.BitLen() < 62 { // need padding, 63 is too close to limit
		n := N.Int64()
		a := int64(math.Sqrt(float64(n)))
		if a*a == n {
			return true
		}
		return false
	}

	// Step 1.1: check if it is a square mod small power of 2
	if _, ok := squaresMod128[uint8(N.Uint64())]; !ok {
		return false
	}

	// Setp 1.2: check if it is a square mod a small number
	_z := uint16(new(big.Int).Mod(N, smallSquareMod).Uint64())
	if _, ok := smallSquares[_z]; !ok {
		return false
	}

	// Step 2: run newtons method, see
	// Cohen's book computational alg. number theory
	// Ch. 1, algorithm 1.7.1
	z := new(big.Int)
	x := new(big.Int).Lsh(one, uint(N.BitLen()+2)>>1)
	y := new(big.Int)
	for {
		// Set y = [(x + [N/x])/2]
		y := y.Rsh(z.Add(x, z.Div(N, x)), 1)
		// if y < x, set x to y
		// else return x
		if y.Cmp(x) == -1 {
			x.Set(y)
		} else {
			return z.Mul(x, x).Cmp(N) == 0
		}
	}
}

type factorization = map[*big.Int]uint64

func factor(N *big.Int) factorization {
	F := make(factorization)
	// Find power of 2 dividing F
	var e int64
	for ; N.Bit(int(e)) == 0; e++ {
		count, _ := F[two]
		F[two] = count + 1
	}
	N.Rsh(N, uint(e))
	// Find upper limit
	s := new(big.Int)
	s.Sqrt(N)
	p := big.NewInt(3)
	// while p < s and N > 1
	q := new(big.Int)
	r := new(big.Int)
	for p.Cmp(s) != 1 && N.Cmp(one) == 1 {
		for e = 1; ; e++ {
			q.QuoRem(N, p, r)
			if r.BitLen() == 0 {
				N.Set(q)
				s.Sqrt(N)
				F[p] = uint64(e)
			} else {
				break
			}
		}
		p = NextPrime(big.NewInt(0).Add(p, two))
	}
	if N.Cmp(one) == 1 {
		F[N] = 1
	}
	return F
}

func lcm(f1, f2 factorization) (f factorization) {
	a1 := make(factorization)
	a2 := make(factorization)
	for p, e := range f1 {
		a1[p] = e
	}
	for p, e := range f2 {
		a2[p] = e
	}
	for p1, e1 := range f1 {
		for p2, e2 := range f2 {
			if p1.Cmp(p2) == 1 {
				if e1 > e2 {
					f[p1] = e1
				} else {
					f[p1] = e2
				}
				delete(a1, p1)
				delete(a2, p2)
			}
		}
		if _, ok := a1[p1]; ok {
			f[p1] = e1
		}
	}
	for p2, e2 := range a2 {
		f[p2] = e2
	}
	return
}
