package powers

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	s1 := &fpn{big.NewInt(1), 0}
	s2 := &fpn{big.NewInt(1), -1}
	s12 := new(fpn).add(s1, s2)
	r := &fpn{big.NewInt(3), -1}
	require.Equal(t, r.n, s12.n)
	require.Equal(t, r.a, s12.a)
}

func TestCmp(t *testing.T) {
	cases := []struct {
		r *fpn
		s *fpn
		c int
	}{
		{new(fpn), new(fpn), 0},
		{new(fpn), &fpn{big.NewInt(1), 0}, -1},
		{&fpn{big.NewInt(1), 0}, new(fpn), 1},
		{&fpn{big.NewInt(3), -1}, &fpn{big.NewInt(10), -2}, -1},
		{&fpn{big.NewInt(8), 0}, &fpn{big.NewInt(1), 3}, 0},
	}
	for _, c := range cases {
		require.Equal(t, c.c, c.r.cmp(c.s), fmt.Sprintln("r = ", c.r.n, ", ", c.r.a, "s = ", c.s.n, ", ", c.s.a))
	}
}

func TestLeq1(t *testing.T) {
	cases := []struct {
		r    *fpn
		leq1 bool
	}{
		{new(fpn), true},
		{&fpn{big.NewInt(1), -1}, true},
		{&fpn{big.NewInt(1), 0}, true},
		{&fpn{big.NewInt(1), 1}, false},
		{&fpn{big.NewInt(-2), -2}, true},
		{&fpn{big.NewInt(3), -2}, true},
		{&fpn{big.NewInt(5), -2}, false},
		{&fpn{big.NewInt(1023), -10}, true},
		{&fpn{big.NewInt(1023), -9}, false},
	}
	for _, c := range cases {
		require.Equal(t, c.leq1, c.r.leq1(), fmt.Sprintln("n = ", c.r.n, ", a = ", c.r.a))
	}
}

func TestRound(t *testing.T) {
	cases := []struct {
		r *fpn
		n *big.Int
	}{
		{new(fpn), big.NewInt(0)},
		{&fpn{big.NewInt(1), 0}, big.NewInt(1)},
		{&fpn{big.NewInt(3), -1}, big.NewInt(2)},
		{&fpn{big.NewInt(3), -2}, big.NewInt(1)},
		{&fpn{big.NewInt(-3), -8}, big.NewInt(0)},
		{&fpn{big.NewInt(100), -3}, big.NewInt(13)},
	}
	for _, c := range cases {
		require.Equal(t, c.r.round(), c.n, fmt.Sprintln("n = ", c.r.n, ", a = ", c.r.a))
	}
}
