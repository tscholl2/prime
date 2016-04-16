package prime

import (
	"math/big"
	"math"
)

// attempt to find a number within b bits of
// y^(-1/k)
func nrootbyk(y *big.Float, k,b int64) *big.Float {
    // TODO
}

// IsKthPower determins if n = a^k for some integer a.
// It is written for positive integers >= 2.
// Returns true if there is such an a.
func IsKthPower(n *big.Int, k int64) bool {
    if n.BitLen() == 1 || k < 2 {
        return false
    }
    z := new(big.Int)
    // approximate 1/n
    y := new(big.Float).Quo(new(big.Float).SetInt64(1),new(big.Float).SetInt(n))
    // set f = floor(log(2n))
    f := int64(math.Log(float64(n.BitLen()+1))/math.Log(2))
    // set b = 3 + ceil(f/k)
    b := 3 + int64(math.Ceil(float64(f)/float64(k)))
    // compute r = nroot_b(y,k)
    
    // find x with |r - x| <= 5/8
    
    // compute sign of n - x^k
    
    // return n == x^k
    
}