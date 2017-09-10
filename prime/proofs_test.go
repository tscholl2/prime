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

func Test_binarySearchKthRoot(t *testing.T) {
	tests := []struct {
		name string
		N    *big.Int
		k    int
		want *big.Int
	}{
		{"4", big.NewInt(4), 2, big.NewInt(2)},
		{"9", big.NewInt(9), 2, big.NewInt(3)},
		{"27", big.NewInt(27), 3, big.NewInt(3)},
		{"125", big.NewInt(125), 3, big.NewInt(5)},
		{"124", big.NewInt(124), 5, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := binarySearchKthRoot(tt.N, tt.k)
			if (got != nil && tt.want != nil && tt.want.Cmp(got) != 0) || (tt.want != nil && got == nil) {
				t.Errorf("newtonsMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isPerfectPower(t *testing.T) {
	tests := []struct {
		name  string
		n     *big.Int
		wantA *big.Int
		wantK int
	}{
		// TODO: Add test cases.
		{"4", big.NewInt(4), big.NewInt(2), 2},
		{"125", big.NewInt(125), big.NewInt(5), 3},
		{"30^3", big.NewInt(27000), big.NewInt(2 * 3 * 5), 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotA, gotK := isPerfectPower(tt.n)
			if gotA == nil && tt.wantA != nil {
				t.Errorf("isPerfectPower() gotA = %v, want %v", gotA, tt.wantA)
			}
			if tt.wantA != nil && (tt.wantA.Cmp(gotA) != 0 || tt.wantK != gotK) {
				t.Errorf("isPerfectPower() gotA = %v, want %v", gotA, tt.wantA)
			}
		})
	}
}
