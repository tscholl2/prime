package prime

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextPrimeProof(t *testing.T) {
	cases := []struct {
		in, want *big.Int
	}{
		{big.NewInt(0), big.NewInt(2)},
		{big.NewInt(4), big.NewInt(5)},
		{big.NewInt(17), big.NewInt(17)},
		{big.NewInt(170), big.NewInt(173)},
		{big.NewInt(1700), big.NewInt(1709)},
		{big.NewInt(17000), big.NewInt(17011)},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, NextPrimeProof(c.in), fmt.Sprintf("in=%d", c.in))
	}
}

func TestSimpleProof(t *testing.T) {
	tests := []struct {
		name string
		N    *big.Int
		want bool
	}{
		{name: "2", N: big.NewInt(2), want: true},
		{name: "4", N: big.NewInt(4), want: false},
		{name: "1001", N: big.NewInt(1001), want: false},
		{name: "1021", N: big.NewInt(1021), want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SimpleProof(tt.N); got != tt.want {
				t.Errorf("SimpleProof() = %v, want %v", got, tt.want)
			}
		})
	}
}
