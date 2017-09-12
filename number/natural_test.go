package number

import (
	"math"
	"math/big"
	"reflect"
	"testing"
)

func Test_number_add(t *testing.T) {
	type args struct {
		y number
		z number
	}
	tests := []struct {
		name string
		x    number
		args args
		want number
	}{
		{name: "1+1", x: nil, args: args{y: number{1}, z: number{1}}, want: number{2}},
		{name: "1+2", x: nil, args: args{y: number{1}, z: number{2}}, want: number{3}},
		{name: "1+max", x: nil, args: args{y: number{1}, z: number{0x7fffffffffffffff}}, want: number{0, 1}},
		{name: "1+maxmax", x: nil, args: args{y: number{1}, z: number{0x7fffffffffffffff, 0x7fffffffffffffff}}, want: number{0, 0, 1}},
		{name: "max+max", x: nil, args: args{y: number{0x7fffffffffffffff}, z: number{0x7fffffffffffffff}}, want: number{0x7fffffffffffffff - 1, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.x.add(tt.args.y, tt.args.z); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("number.add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Benchmark_number_add(b *testing.B) {
	var x number
	y := number{0, 0x7fffffffffffffff}
	z := number{0, 0x7fffffffffffffff - 1}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.add(y, z)
	}
}

func Benchmark_big_add(b *testing.B) {
	x := new(big.Int)
	y := new(big.Int).Lsh(big.NewInt(math.MaxInt64), 64)
	z := new(big.Int).Lsh(big.NewInt(math.MaxInt64-1), 64)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.Add(y, z)
	}
}
