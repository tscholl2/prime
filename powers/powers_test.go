package powers

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	onerat  = big.NewRat(1, 1)
	tworat  = big.NewRat(2, 1)
	halfrat = big.NewRat(1, 2)
)

func toRat(r *fpn) *big.Rat {
	f := new(big.Rat).SetInt(r.n)
	for i := 0; i < r.a; i++ {
		f.Mul(f, tworat)
	}
	for i := r.a; i < 0; i++ {
		f.Mul(f, halfrat)
	}
	return f
}

func powrat(r *big.Rat, k int) *big.Rat {
	if k > 1000 || k < -1000 {
		panic("too much")
	}
	s := new(big.Rat).Set(onerat)
	for i := 0; i < k; i++ {
		s.Mul(s, r)
	}
	for i := k; i < 0; i++ {
		s.Mul(s, new(big.Rat).Inv(r))
	}
	return s
}

func pr(r *big.Rat) string {
	f, _ := r.Float64()
	return fmt.Sprintf("%0.10f", f)
}

func TestDivb(t *testing.T) {
	cases := []struct {
		r *fpn
		k uint
		b uint
	}{
		{&fpn{big.NewInt(1024), 0}, 1, 10},
		{&fpn{big.NewInt(1024), 0}, 2, 10},
		{&fpn{big.NewInt(3255), 0}, 325, 18},
		{&fpn{big.NewInt(4327), -3}, 4, 100},
		{&fpn{big.NewInt(1111), 0}, 100, 17},
	}
	for _, c := range cases {
		s := toRat(divb(c.r, c.k, c.b))
		k := big.NewRat(int64(c.k), 1)
		sk := new(big.Rat).Mul(s, k)
		sk12b := new(big.Rat).Mul(
			new(big.Rat).Mul(s, k),
			new(big.Rat).Add(onerat, powrat(tworat, int(c.k))),
		)
		r := toRat(c.r)
		require.True(t, sk.Cmp(r) <= 0, fmt.Sprintln("sk = ", sk, ", r = ", c.r))
		require.True(t, r.Cmp(sk12b) == -1, fmt.Sprintln("r = ", c.r, ", sk12b = ", sk12b))
	}
}

func TestTruncb(t *testing.T) {
	cases := []struct {
		r *fpn
		b uint
	}{
		{&fpn{big.NewInt(1024), 0}, 10},
		{&fpn{big.NewInt(1024), 0}, 3},
		{&fpn{big.NewInt(3255), 0}, 18},
		{&fpn{big.NewInt(4327), -3}, 100},
		{&fpn{big.NewInt(1111), 0}, 17},
	}
	for _, c := range cases {
		expected := divb(c.r, 1, c.b)
		output := truncb(c.r, c.b)
		require.Equal(t, expected, output, fmt.Sprintf("r = %s, b = %d", c.r, c.b))
	}
}

func TestPowb(t *testing.T) {
	cases := []struct {
		r *fpn
		k uint
		b uint
	}{
		{&fpn{big.NewInt(1024), 0}, 1, 10},
		{&fpn{big.NewInt(1024), 0}, 2, 10},
		{&fpn{big.NewInt(3255), 0}, 30, 18},
		{&fpn{big.NewInt(4327), -3}, 4, 100},
		{&fpn{big.NewInt(1111), 0}, 5, 17},
	}
	for _, c := range cases {
		// returns s such that s <= r^k < s(1 + 2^(1 - b))^(2k - 1)
		r := toRat(c.r)
		s := toRat(powb(c.r, c.k, c.b))
		rk := powrat(r, int(c.k))
		s12bk := new(big.Rat).Mul(
			s,
			powrat(new(big.Rat).Add(onerat, powrat(tworat, int(c.k))), int(2*c.k-1)),
		)
		require.True(t, s.Cmp(rk) <= 0, fmt.Sprintf("s = %s, rk = %s", s, rk))
		require.True(t, rk.Cmp(s12bk) == -1, fmt.Sprintf("rk = %s, s12bk = %s", rk, s12bk))
	}
}

func TestAlgB(t *testing.T) {
	cases := []struct {
		r *fpn
		k uint
		b uint
	}{
		{&fpn{big.NewInt(1000), -1}, 60, 9},
		{&fpn{big.NewInt(43628), 8}, 45, 9},
		{&fpn{big.NewInt(119), -3}, 62, 8},
	}
	for _, c := range cases {
		yinv := new(big.Rat).Inv(toRat(c.r))
		sk := powrat(toRat(algB(c.r, c.k, c.b)), int(c.k))
		a1mbk := powrat(big.NewRat(1<<c.b-1, 1<<c.b), int(c.k))
		a1pbk := powrat(big.NewRat(1<<c.b+1, 1<<c.b), int(c.k))
		require.True(t, new(big.Rat).Mul(sk, a1mbk).Cmp(yinv) == -1)
		require.True(t, yinv.Cmp(new(big.Rat).Mul(sk, a1pbk)) == -1)
	}
}

func TestAlgN(t *testing.T) {
	cases := []struct {
		r *fpn
		k uint
		b uint
	}{
		{&fpn{big.NewInt(1001), -1}, 60, 11},
	}
	for _, c := range cases {
		yinv := new(big.Rat).Inv(toRat(c.r))
		sk := powrat(toRat(algN(c.r, c.k, c.b)), int(c.k))
		a1mbk := powrat(big.NewRat(1<<c.b-1, 1<<c.b), int(c.k))
		a1pbk := powrat(big.NewRat(1<<c.b+1, 1<<c.b), int(c.k))
		require.True(t, new(big.Rat).Mul(sk, a1mbk).Cmp(yinv) == -1)
		require.True(t, yinv.Cmp(new(big.Rat).Mul(sk, a1pbk)) == -1)
	}
}

func TestAlgC(t *testing.T) {
	cases := []struct {
		n    *big.Int
		x    *fpn
		k    uint
		sign int
	}{
		{big.NewInt(81), &fpn{big.NewInt(4), 0}, 3, 1},
	}
	for _, c := range cases {
		s := algC(c.n, c.x, c.k)
		require.Equal(t, c.sign, s, fmt.Sprintf("n = %s, x = %s, k = %d", c.n, c.x, c.k))
	}
}

func TestAlgX(t *testing.T) {
	cases := []struct {
		n *big.Int
		x *big.Int
		p uint
	}{
		{big.NewInt(81), big.NewInt(9), 2},
	}
	for _, c := range cases {
		x, p := algX(c.n)
		require.Equal(t, c.x, x, fmt.Sprintf("n = %s, x = %s, p = %d", c.n, c.x, c.p))
		require.Equal(t, c.p, p, fmt.Sprintf("n = %s, x = %s, p = %d", c.n, c.x, c.p))
	}
}
