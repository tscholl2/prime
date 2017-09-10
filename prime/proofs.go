package prime

import (
	"math/big"
)

func NextPrimeProof(N *big.Int) *big.Int {
	p := NextPrime(N)
	for !SimpleProof(p) {
		p = NextPrime(p)
	}
	return p
}

func SimpleProof(N *big.Int) bool {
	if SmallPrimeTest(N) == IsPrime {
		return true
	}
	z := new(big.Int)
	if N.Cmp(big.NewInt(int64(primes10[len(primes10)-1]))) < 1 {
		return false
	}
	s := new(big.Int)
	s.Sqrt(N)
	d := big.NewInt(2)
	for d.Cmp(s) != 1 {
		if z.Mod(N, d); z.Cmp(zero) == 0 {
			return false
		}
		d.Add(d, one)
	}
	return true
}

func factorProof(N *big.Int) factorization {
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
		p = NextPrimeProof(big.NewInt(0).Add(p, two))
	}
	if N.Cmp(one) == 1 {
		F[N] = 1
	}
	return F
}

// Returns a,k such that k>=2 and n = a^k or nil,0 if none exist.
func isPerfectPower(n *big.Int) (a *big.Int, k int) {
	maxK := n.BitLen() + 1
	for k := 2; k < maxK; k++ {
		a = binarySearchKthRoot(n, k)
		if a != nil {
			return a, k
		}
	}
	return nil, 0
}

// Returns a such that n = a^k or nil if no such a exists.
func binarySearchKthRoot(n *big.Int, k int) (a *big.Int) {
	if k == 1 {
		return n
	}
	K := big.NewInt(int64(k))
	test := func(a *big.Int) int {
		return new(big.Int).Exp(a, K, nil).Cmp(n)
	}
	a = new(big.Int)
	low := big.NewInt(2)
	if test(low) == 0 {
		return low
	}
	high := new(big.Int).Set(n) // TODO rsh by log_2(n)/k ish
	for {
		a.Rsh(a.Add(low, high), 1)
		t := test(a)
		if t == 0 {
			return a
		}
		if t == 1 {
			high.Set(a)
		} else {
			low.Set(a)
		}
		if a.Sub(high, low).Cmp(one) == 0 {
			return nil
		}
	}
}
