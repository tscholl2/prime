package prime

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tests

func TestRandBig(t *testing.T) {
	for i := 0; i < 1000; i++ {
		x := randBig(i)
		require.Equal(t, x.BitLen(), i)
		require.True(t, x.Sign() >= 0)
	}
}

func TestRandPrime(t *testing.T) {
	for i := 2; i < 100; i++ {
		p := RandPrime(i)
		require.Equal(t, i, p.BitLen())
	}
}

func TestTrailingZeroBits(t *testing.T) {
	cases := []struct {
		in   *big.Int
		want uint
	}{
		{big.NewInt(0), 0},
		{big.NewInt(1), 0},
		{big.NewInt(2), 1},
		{big.NewInt(3), 0},
		{big.NewInt(4), 2},
		{big.NewInt(6), 1},
		{big.NewInt(8), 3},
		{big.NewInt(15), 0},
		{big.NewInt(16), 4},
		{big.NewInt(32), 5},
		{big.NewInt(3571), 0},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, trailingZeroBits(c.in), fmt.Sprintf("in=%d", c.in))
	}
}

func TestIsSquare(t *testing.T) {
	n1true, _ := new(big.Int).SetString("240e16068a04dea390a1f96b3f05a1", 16)
	n1false, _ := new(big.Int).SetString("240e16068a04dea390a1f96b3f05a2", 16)
	n2true, _ := new(big.Int).SetString("fa8bf08953f8b2c1f941de3fd45b952967a055ff7826e4a436b660db443b024eaeed6fdf0640", 16)
	n2false, _ := new(big.Int).SetString("fa8bf08953f8b2c1f941de3fd45b952967a055ff7826e4a436b660db443b024eaeed6fdf0641", 16)
	n3false, _ := new(big.Int).SetString("1e04ded686bffea61355f4c9c76f1e66fba27b9fa8b00f3c5884d3eff369677ad5817d783aa58db408de1310e55cd5e72a8176340", 16)
	n3true, _ := new(big.Int).SetString("1e04ded686bffea61355f4c9c76f1e66fba27b9fa8b00f3c5884d3eff369677ad5817d783aa58db408de1310e55cd5e72a8176341", 16)
	n3false2, _ := new(big.Int).SetString("1e04ded686bffea61355f4c9c76f1e66fba27b9fa8b00f3c5884d3eff369677ad5817d783aa58db408de1310e55cd5e72a8176342", 16)
	n4true, _ := new(big.Int).SetString("7afee5555433fa458dc6e8e62f1cc4533b3488893e4067830385d9b27fbf724f0ca5e4e94a1c46afb09138c1965d8aa8938bebd89ae3b4f13aecd85839f3b5db1c7b9692bc0ef2595cf8640", 16)
	n4false, _ := new(big.Int).SetString("7afee5555433fa458dc6e8e62f1cc4533b3488893e4067830385d9b27fbf724f0ca5e4e94a1c46afb09138c1965d8aa8938bebd89ae3b4f13aecd85839f3b5db1c7b9692bc0ef2595cf8641", 16)
	largetrue, _ := new(big.Int).SetString("3b17f061370666c4f11db552e1dc533fbf30531421a6292207fd136a94f9f011e672a24f0ef1422210ab44f96e43599d6576030ded2b0f9c79fc8b8efd8558f09c168e35895707d7749fb92e18d9f0653efdc05daeee522204766c6aea0f2dbc5793beabbd629e69b38f5c0c56a37fd4ceb27d667ab9d1b098dae5beec2d3bfa96be55a3b9262d5662429ba76fb4f359d5674c0d861c81", 16)
	largefalse, _ := new(big.Int).SetString("3b17f061370666c4f11db552e1dc533fbf30531421a6292207fd136a94f9f011e672a24f0ef1422210ab44f96e43599d6576030ded2b0f9c79fc8b8efd8558f09c168e35895707d7749fb92e18d9f0653efdc05daeee522204766c6aea0f2dbc5793beabbd629e69b38f5c0c56a37fd4ceb27d667ab9d1b098dae5beec2d3bfa96be55a3b9262d5662429ba76fb4f359d5674c0d861d81", 16)
	for i := 0; i < 100; i++ {
		// randomness in the sqrt function
		// needs a lot of testing to find
		// edge cases more easily
		cases := []struct {
			in   *big.Int
			want bool
		}{
			{big.NewInt(-1436278), false},
			{big.NewInt(0), true},
			{big.NewInt(1), true},
			{big.NewInt(15), false},
			{big.NewInt(16), true},
			{big.NewInt(3571), false},
			{big.NewInt(13627856 * 13627856), true},
			{big.NewInt(13627856), false},
			{n1true, true},
			{n2true, true},
			{n3true, true},
			{n4true, true},
			{n1false, false},
			{n2false, false},
			{n3false, false},
			{n3false2, false},
			{n4false, false},
			{largetrue, true},
			{largefalse, false},
		}
		for _, c := range cases {
			assert.Equal(t, c.want, IsSquare(c.in), fmt.Sprintf("in=%d", c.in))
		}
	}
	// random tests
	for i := 0; i < 100; i++ {
		x := randBig(1000)
		require.True(t, IsSquare(new(big.Int).Mul(x, x)))
		require.False(t, IsSquare(new(big.Int).Mul(x, new(big.Int).Add(x, one))))
	}
}

func TestStrongLucasSelfridgeTest(t *testing.T) {
	n, _ := new(big.Int).SetString("319889369713946602502766595032347", 10)
	//http://www.sciencedirect.com/science/article/pii/S0747717185710425
	cases := []struct {
		in   *big.Int
		want bool
	}{
		{big.NewInt(3 * 5 * 11 * 13 * 17), false}, // smooth number
		{big.NewInt(3), true},                     // some small primes
		{big.NewInt(5), true},
		{big.NewInt(11), true},
		{big.NewInt(797), true},
		{big.NewInt(3571 * 3571), false}, // perfect square
		{big.NewInt(3571), true},         // large prime
		{big.NewInt(5459), true},         // NOT prime! a strong Lucas psuedoprime
		{n, true},                        //also a strong lsps!, BPSW says composite though
		{big.NewInt(364387 * 362751), false},
		{big.NewInt(364387 * 362753), false},
		{big.NewInt(364387 * 362755), false},
		{big.NewInt(364387 * 362757), false},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, StrongLucasSelfridgeTest(c.in), fmt.Sprintf("in=%d", c.in))
	}
}

func TestMillerRabin(t *testing.T) {
	cases := []struct {
		inN  *big.Int
		inA  int64
		want bool
	}{
		{big.NewInt(221), 174, true},
		{big.NewInt(221), 137, false},
		{big.NewInt(7), 2, true},
		{big.NewInt(11), 2, true},
		{big.NewInt(13), 2, true},
		{big.NewInt(1709), 2, true},
		{big.NewInt(2005), 2, false},
		{big.NewInt(2047), 2, true}, // NOT prime!
		{big.NewInt(173), 6, true},
		{big.NewInt(175), 5, false}, // not relatively prime
		{big.NewInt(217), 6, true},  // NOT prime!
	}
	for _, c := range cases {
		assert.Equal(t, c.want, StrongMillerRabin(c.inN, c.inA), fmt.Sprintf("N=%d, A=%d", c.inN, c.inA))
	}
}

func TestNextPrime(t *testing.T) {
	cases := []struct {
		in, want *big.Int
	}{
		{big.NewInt(0), big.NewInt(2)},
		{big.NewInt(4), big.NewInt(5)},
		{big.NewInt(17), big.NewInt(17)},
		{big.NewInt(170), big.NewInt(173)},
		{big.NewInt(1700), big.NewInt(1709)},
		{big.NewInt(17000), big.NewInt(17011)},
		{big.NewInt(170000), big.NewInt(170003)},
		{big.NewInt(1700000), big.NewInt(1700021)},
		{new(big.Int).SetBytes([]byte{0x93, 0x5a, 0x53, 0xf3, 0x89}),
			new(big.Int).SetBytes([]byte{0x93, 0x5a, 0x53, 0xf3, 0x8d})},
		{new(big.Int).SetBytes([]byte{0x1, 0xd2, 0x19, 0x3a, 0x34, 0x58, 0xd0, 0x22, 0x96, 0x33, 0x9c, 0xbb}),
			new(big.Int).SetBytes([]byte{0x1, 0xd2, 0x19, 0x3a, 0x34, 0x58, 0xd0, 0x22, 0x96, 0x33, 0x9c, 0xc1})},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, NextPrime(c.in), fmt.Sprintf("in=%d", c.in))
	}
}

func TestJacobiSymbol(t *testing.T) {
	cases := []struct {
		N, D *big.Int
		want int
	}{
		{big.NewInt(15), big.NewInt(45), 0},
		{big.NewInt(19), big.NewInt(45), 1},
		{big.NewInt(8), big.NewInt(21), -1},
		{big.NewInt(5), big.NewInt(21), 1},
		{big.NewInt(1001), big.NewInt(9907), -1},
		{big.NewInt(-7), big.NewInt(5459), -1},
		{big.NewInt(7), big.NewInt(5459), 1},
		{big.NewInt(21), big.NewInt(3333), 0},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, JacobiSymbol(c.N, c.D), fmt.Sprintf("N=%d, D=%d", c.N, c.D))
	}
}
