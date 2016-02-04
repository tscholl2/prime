/*
Package prime is another prime library for various prime number
related functions. Mostly written just to learn Go and
computational number theory in another language.
*/
package prime

import (
	"crypto/rand"
	"math"
	"math/big"
	"sort"
)

var (
	one = big.NewInt(1)
	two = big.NewInt(2)
)

func randBig(bits int) *big.Int {
	if bits <= 0 {
		return new(big.Int)
	}
	bytes := make([]byte, bits/8)
	rand.Read(bytes)
	return new(big.Int).SetBit(new(big.Int).SetBytes(bytes), bits-1, 1)
}

// RandPrime returns a random prime
// of a given bit size. For small bits
// it just gives something close.
func RandPrime(bits int) (p *big.Int) {
	if bits <= 10 {
		start := sort.Search(len(primes10), func(i int) bool {
			return big.NewInt(int64(primes10[i])).BitLen() >= bits
		})
		slice := primes10[start:]
		end := sort.Search(len(slice), func(i int) bool {
			return big.NewInt(int64(slice[i])).BitLen() > bits
		})
		set := slice[:end]
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(set))))
		return big.NewInt(int64(set[int(n.Int64())]))
	}
	for {
		N := randBig(bits)
		p := NextPrime(N)
		if p.BitLen() == bits || bits < 5 {
			return p
		}
	}
}

// NextPrime returns a number, p >= N with
// high probability that p is the next prime
// occuring after N.
func NextPrime(N *big.Int) (p *big.Int) {
	if N.Sign() <= 0 {
		return big.NewInt(2)
	}
	if N.BitLen() <= 10 {
		n := uint16(N.Int64())
		i := sort.Search(len(primes10), func(i int) bool { return primes10[i] >= n })
		if i < len(primes10) {
			return big.NewInt(int64(primes10[i]))
		}
	}
	m := len(diffs210)
	i := int(new(big.Int).Mod(N, big.NewInt(int64(m))).Int64())
	p = new(big.Int).Set(N)
	for {
		if BPSW(p) {
			return
		}
		p.Add(p, big.NewInt(int64(diffs210[i])))
		i = (i + diffs210[i]) % m
	}
}

// BPSW runs the Baillie-PSW primality test on N.
// It returns true if N is probably a prime, and false
// otherwise.
//
// For more see http://www.trnicely.net/misc/bpsw.html
func BPSW(N *big.Int) bool {
	//Step 0: parse input
	if N.Sign() <= 0 {
		panic("BPSW is for positive integers only")
	}

	// Step 1: check  all small primes
	// returns 1 if prime, 0 if composite, -1 else
	switch SmallPrimeTest(N) {
	case 1:
		return true
	case 0:
		return false
	}

	// Step 2: Miller-Rabin test
	// returns false if composite
	if !StrongMillerRabin(N, 2) {
		return false
	}

	// Step 3: Lucas-Selfridge test
	// returns false if composite
	if !StrongLucasSelfridgeTest(N) {
		return false
	}

	// Step 4: If didn't fail other tests
	// return true, i.e. this passed
	return true
}

// SmallPrimeTest returns
// 1 if N is prime
// 0 if N is composite
// -1 if undetermined
func SmallPrimeTest(N *big.Int) int {
	if N.Sign() <= 0 {
		panic("SmallPrimeTest for positive integers only")
	}
	if N.BitLen() <= 10 {
		n := uint16(N.Int64())
		i := sort.Search(len(primes10), func(i int) bool { return primes10[i] >= n })
		if i >= len(primes10) || n != primes10[i] {
			return 0
		}
		return 1
	}
	if new(big.Int).GCD(nil, nil, N, prodPrimes10).Cmp(one) == 1 {
		return 0
	}
	return -1
}

// StrongMillerRabin returns true if N is a
// strong Miller-Rabin psuedoprime in base a.
// That is, it returns false if a is a witness
// for compositeness of N or N is a strong
// pseudoprime base a.
//
// ProbablyPrime if you want to do a lot
// of random tests, this is for one specific
// base value.
func StrongMillerRabin(N *big.Int, a int64) bool {
	// Step 0: parse input
	if N.Sign() < 0 || N.Bit(0) == 0 || a < 2 {
		panic("MR is for positive odd integers with a >= 2")
	}
	A := big.NewInt(a)
	if new(big.Int).GCD(nil, nil, N, A).Cmp(one) != 0 {
		return false
	}

	// Step 1: find d,s, so that n - 1 = d*2^s
	// with d odd
	d := new(big.Int).Sub(N, one)
	s := trailingZeroBits(d)
	d.Rsh(d, s)

	// Step 2: compute powers a^d
	// and then a^(d*2^r) for 0<r<s
	var nm1, Ad big.Int
	Ad.Exp(A, d, N)
	nm1.Sub(N, one)
	if Ad.Cmp(one) == 0 || Ad.Cmp(&nm1) == 0 {
		return true
	}
	for r := uint(1); r < s; r++ {
		Ad.Exp(&Ad, two, N)
		if Ad.Cmp(&nm1) == 0 {
			return true
		}
	}

	// Step 3: a is not a witness
	return false
}

// StrongLucasSelfridgeTest returns true if N is
// a strong Lucas-Selfridge pseudoprime and
// false otherwise.
//
// For more information see
// http://www.trnicely.net/misc/bpsw.html
func StrongLucasSelfridgeTest(N *big.Int) bool {
	// Step 0: parse input
	if N.Sign() < 0 || N.Bit(0) == 0 {
		panic("LS is for positive odd integers only")
	}

	// Step 1: check if N is a perfect square
	if IsSquare(N) {
		return false
	}

	// Step 2: find the first element D in the
	// sequence {5, -7, 9, -11, 13, ...} such that
	// Jacobi(D,N) = -1 (Selfridge's algorithm)
	D := big.NewInt(5)
	for JacobiSymbol(D, N) != -1 {
		if D.Sign() < 0 {
			D.Sub(D, two)
		} else {
			D.Add(D, two)
		}
		D.Neg(D)
	}
	P := big.NewInt(1) // Selfridge's choice, also set on wiki package
	// http://en.wikipedia.org/wiki/Lucas_pseudoprime#Implementing_a_Lucas_probable_prime_test
	Q := new(big.Int).Sub(one, D)
	Q.Rsh(Q, 2) // divide by 4
	Q.Mod(Q, N)
	if new(big.Int).GCD(nil, nil, N, Q).Cmp(one) != 0 {
		// sanity check
		return false
	}

	// Step 3: Find d so N+1 = 2^s*d with d odd
	d := new(big.Int).Add(N, one)
	s := trailingZeroBits(d)
	d.Rsh(d, s)

	// Step 4: Calculate the U's and V's
	// return true if we have any of the equalities (mod N)
	// U_d=0, V_d=0, V_2d=0, V_4d=0, V_8d=0,...,V_{2^(s-1)d}
	divideBy2ModN := func(x *big.Int) *big.Int {
		if x.Bit(0) != 0 {
			x.Add(x, N)
		}
		return x.Rsh(x, 1)
	}
	var tmp, PxUk, DxUk, PxVk big.Int
	Uk := big.NewInt(0)         // U_0 = 0
	Vk := new(big.Int).Set(two) // V_0 = 2
	Qk := new(big.Int).Set(one) // Q^0 = 1
	// follow repeated squaring algorithm
	for i := d.BitLen() - 1; i > -1; i-- {
		// double everything
		Uk.Mul(Uk, Vk)
		Uk.Mod(Uk, N) // now U_{2k}
		Vk.Mul(Vk, Vk)
		Vk.Sub(Vk, tmp.Lsh(Qk, 1))
		Vk.Mod(Vk, N) // now V_{2k}
		Qk.Mul(Qk, Qk)
		Qk.Mod(Qk, N) // now Q^{2k}
		if d.Bit(i) == 1 {
			// if bit is set then increment by 1
			Qk.Mul(Qk, Q)
			Qk.Mod(Qk, N) // now Q^{2k+1}
			PxUk.Mul(P, Uk)
			PxUk.Mod(&PxUk, N)
			DxUk.Mul(D, Uk)
			DxUk.Mod(&DxUk, N)
			PxVk.Mul(P, Vk)
			PxVk.Mod(&PxVk, N)
			Uk.Mod(divideBy2ModN(tmp.Add(&PxUk, Vk)), N)    // now U_{2k+1}
			Vk.Mod(divideBy2ModN(tmp.Add(&DxUk, &PxVk)), N) // now V_{2k+1}
		}
	}
	// U_k, V_k, Q^k are now all with k=d
	if Uk.Sign() == 0 {
		// if U_d = 0
		return true
	}
	// Now we look at powers V_{{2^r}d} for r = 0..s-1
	var r uint
	for r = 0; r < s; r++ {
		if Vk.Sign() == 0 {
			// if V_{2^rd} = 0
			return true
		}
		Vk.Mul(Vk, Vk)
		Vk.Sub(Vk, tmp.Lsh(Qk, 1))
		Vk.Mod(Vk, N) // V_{2^{r+1}d}
		Qk.Mul(Qk, Qk)
		Qk.Mod(Qk, N) // Q_{2^(r+1)d}
	}

	// Step 5: return false because it didn't pass the test
	return false
}

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
	if x.Sign() == 0 {
		return 0
	}
	if x.Bit(0) == 1 {
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
	d := N.BitLen()
	if d < 64 {
		n := N.Int64()
		a := int64(math.Sqrt(float64(n)))
		if a*a == n {
			return true
		}
		return false
	}

	// Step 0.5: Check N has an even number of factors of 2
	s := trailingZeroBits(N)
	if s%2 != 0 {
		return false
	}

	// Step 1: Check if it is a square mod something small
	var n big.Int
	n.Set(N)
	n.Mod(n.Rsh(&n, s), smallSquareMod)
	if n.Sign() > 0 {
		m := uint16(n.Bits()[0]) // ok because smallSquareMod < 2^16
		i := sort.Search(len(smallSquares), func(i int) bool { return smallSquares[i] >= m })
		if i < 0 || i >= len(smallSquares) || m != smallSquares[i] {
			return false
		}
	}

	// Step 2: make a random guess for sqrt
	// with the right order of magnitude
	bytes := make([]byte, d/16)
	rand.Read(bytes)
	var x, y, delta big.Int
	x.SetBytes(bytes)
	y.Set(&x)

	// Step 3: run newtons method until it
	// stabilized (same value or one off), see wiki article
	// http://en.wikipedia.org/wiki/Integer_square_root
	// if it doesn't converge it should alternate between +-1
	// so return false in that case
	// convergence is fast, should take log(number of digits)
	// with some coefficient... 5 seems like it works
	for i := 0; ; i++ {
		// Set y = [(x + [N/x])/2]
		y.Rsh(y.Add(&y, x.Div(N, &x)), 1) // note: at this point y = x
		if i > int(math.Log(float64(d)))*5 {
			delta.Sub(&x, &y)
			if len(delta.Bits()) == 0 || delta.Bits()[0] == 1 {
				// if |x - y| <= 1
				return delta.Mul(&x, &x).Cmp(N) == 0
			}
		}
		x.Set(&y)
	}
}
