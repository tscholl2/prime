package number

import (
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

func Test_number_setHex(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		z    number
		args args
		want number
	}{
		{"1", nil, args{"01"}, number{1}},
		{"0", nil, args{"00"}, number{}},
		{"7fffffffffffffff", nil, args{"7fffffffffffffff"}, number{0x7fffffffffffffff}},
		{"ffffffffffffffff", nil, args{"ffffffffffffffff"}, number{0x7fffffffffffffff, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.z.setHex(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("number.setHex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_number_lsh(t *testing.T) {
	type args struct {
		x number
		s uint
	}
	tests := []struct {
		name string
		z    number
		args args
		want number
	}{
		{"1<<1", nil, args{number{1}, 1}, number{1 << 1}},
		{"1<<63", nil, args{number{1}, 63}, number{0, 1}},
		{"1<<64", nil, args{number{1}, 64}, number{0, 1 << 1}},
		{"2<<64", nil, args{number{2}, 64}, number{0, 2 << 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.z.lsh(tt.args.x, tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("number.lsh() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Benchmark_number_add(b *testing.B) {
	x := number{}.setHex("ffffffffffffffafffffffffffffffffffffffffffffffffffffffffffffffffffafffffffffffffffffffffffffffffffffffffffffffffffffffafffffffffffffffffffffffffffffffffffffffffffffffffffafffffffffffffffffffffffffffffffffffff")
	y := number{}.setHex("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffef")
	z := number{}.add(x, y)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		z.add(x, y)
	}
}

func Benchmark_big_add(b *testing.B) {
	x, ok := new(big.Int).SetString("ffffffffffffffafffffffffffffffffffffffffffffffffffffffffffffffffffafffffffffffffffffffffffffffffffffffffffffffffffffffafffffffffffffffffffffffffffffffffffffffffffffffffffafffffffffffffffffffffffffffffffffffff", 16)
	if !ok {
		panic("uh oh")
	}
	y, ok := new(big.Int).SetString("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffef", 16)
	if !ok {
		panic("uh oh")
	}
	z := new(big.Int).Add(x, y)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		z.Add(x, y)
	}
}
