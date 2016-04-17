package prime

import (
	"crypto/rand"
	"math/big"
	"sort"
)

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
	if !StrongLucasSelfridge(N) {
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
		n := uint16(N.Uint64())
		i := sort.Search(len(primes10), func(i int) bool {
			return primes10[i] >= n
		})
		if i >= len(primes10) || n != primes10[i] {
			return 0
		}
		return 1
	}
	// quick test for N even
	if N.Bits()[0]&1 == 0 {
		return 0
	}
	// compare several small gcds for efficency
	z := new(big.Int)
	if z.GCD(nil, nil, N, prodPrimes10A).Cmp(one) == 1 {
		return 0
	}
	if z.GCD(nil, nil, N, prodPrimes10B).Cmp(one) == 1 {
		return 0
	}
	if z.GCD(nil, nil, N, prodPrimes10C).Cmp(one) == 1 {
		return 0
	}
	if z.GCD(nil, nil, N, prodPrimes10D).Cmp(one) == 1 {
		return 0
	}
	return -1
}

// StrongMillerRabin returns true if N is a
// strong Miller-Rabin pseudoprime in base a.
// That is, it returns false if a is a witness
// for compositeness of N or N is a strong
// pseudoprime base a.
//
// Use builtin ProbablyPrime if you want to do a lot
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
	nm1 := new(big.Int).Sub(N, one)
	Ad := new(big.Int).Exp(A, d, N)
	if Ad.Cmp(one) == 0 || Ad.Cmp(nm1) == 0 {
		return true
	}
	for r := uint(1); r < s; r++ {
		Ad.Exp(Ad, two, N)
		if Ad.Cmp(nm1) == 0 {
			return true
		}
	}

	// Step 3: a is not a witness
	return false
}

// StrongLucasSelfridge returns true if N is
// a strong Lucas-Selfridge pseudoprime and
// false otherwise.
//
// For more information see
// http://www.trnicely.net/misc/bpsw.html
func StrongLucasSelfridge(N *big.Int) bool {
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

// SolovayStrassen returns true if after choosing k
// random numbers in [2,...,N] we did not find an
// "Euler liar". That is, every number a we chose
// satisfied a^((n-1)/2) = Jacobi(a/N) mod N.
// See https://en.wikipedia.org/wiki/Solovay%E2%80%93Strassen_primality_test.
// Probability it passes and is not prime is 2^(-k).
func SolovayStrassen(N *big.Int, k int) bool {
	s := SmallPrimeTest(N)
	if s == 1 {
		return true
	}
	if s == 0 {
		return false
	}
	a := new(big.Int)
	b := make([]byte, N.BitLen())
	for i := 0; i < k; i++ {
		rand.Read(b)
		a.SetBytes(b)
		if !basedSolovayStrassen(N, a) {
			return false
		}
	}
	return true
}

func basedSolovayStrassen(N, a *big.Int) (t bool) {
	// we assume N is odd
	x := JacobiSymbol(a, N)
	if x == 0 {
		return
	}
	z := new(big.Int)
	z.Exp(a, z.Rsh(z.Sub(N, one), 1), N) // this step is expensive
	if (x == 1 && z.Cmp(one) == 0) || (x == -1 && z.Sub(N, z).Cmp(one) == 0) {
		t = true
	}
	return
}
