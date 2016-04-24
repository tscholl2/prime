package float

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDivB(t *testing.T) {
	cases := []struct {
		r     *fpn
		k     *big.Int
		b     uint
		wantN *big.Int
		wantA int
	}{
		{&fpn{big.NewInt(1024), 0}, big.NewInt(1), 0, one, 10},
	}
	for _, c := range cases {
		r := divb(c.r.n, c.r.a, c.k, c.b).normalize()
		assert.Equal(t, r.n, c.wantN)
		assert.Equal(t, r.a, c.wantA)
	}
}
