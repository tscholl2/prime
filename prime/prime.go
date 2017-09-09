/*
Package prime is another prime library for various prime number
related functions. Mostly written just to learn Go and
computational number theory in another language.
*/
package prime

import (
	"crypto/rand"
	"math/big"
	"sort"
)

var (
	zero = big.NewInt(0)
	one  = big.NewInt(1)
	two  = big.NewInt(2)
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
// occurring after N.
func NextPrime(N *big.Int) (p *big.Int) {
	if N.Sign() <= 0 {
		return big.NewInt(2)
	}
	if N.BitLen() <= 10 {
		n := uint16(N.Int64())
		i := sort.Search(len(primes10), func(i int) bool {
			return primes10[i] >= n
		})
		if i < len(primes10) {
			return big.NewInt(int64(primes10[i]))
		}
	}
	m := len(diffs210)
	i := int(new(big.Int).Mod(N, big.NewInt(int64(m))).Int64())
	p = new(big.Int).Set(N)
	for {
		if BPSW(p) != IsComposite {
			return
		}
		p.Add(p, big.NewInt(int64(diffs210[i])))
		i = (i + diffs210[i]) % m
	}
}
