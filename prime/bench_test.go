package prime

import (
	cryptoRand "crypto/rand"
	"math/big"
	"math/rand"
	"testing"
)

// Benchmark constants
var benchmarkPrime, _ = new(big.Int).SetString("3ba9a88eb20cfdfe4a380607f5025cdcd0f0bbb73b6f8d45bb0d7bdcd7d485b513d4f8c3d0d572f47ea6f32b4d19978c1a578f919c126e997548b8d0acc64284287a3a321e292e1be9614bf21254011a25df84b77b7411d41e65fd50298fc4660651580b5bd3f38377e2a6260021694cb4096873762f45ba41562ed1cddaca67", 16)

// utilities

func BenchmarkJacobiSymbol(b *testing.B) {
	benchmarkOdd, _ := new(big.Int).SetString("3ba9a88eb20cfdfe4a380607f5025cdcd0f0bbb73b6f8d45bb0d7bdcd7d485b513d4f8c3d0d572f47ea6f32b4d19978c1a578f919c126e997548b8d0acc64284287a3a321e292e1be9614bf21254011a25df84b77b7411d41e65fd50298fc4660651580b5bd3f38377e2a6260021694cb4096873762f45ba41562ed1cddaca69", 16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		JacobiSymbol(benchmarkPrime, benchmarkOdd)
	}
}

func BenchmarkJacobiSymbolStandard(b *testing.B) {
	benchmarkOdd, _ := new(big.Int).SetString("3ba9a88eb20cfdfe4a380607f5025cdcd0f0bbb73b6f8d45bb0d7bdcd7d485b513d4f8c3d0d572f47ea6f32b4d19978c1a578f919c126e997548b8d0acc64284287a3a321e292e1be9614bf21254011a25df84b77b7411d41e65fd50298fc4660651580b5bd3f38377e2a6260021694cb4096873762f45ba41562ed1cddaca69", 16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		big.Jacobi(benchmarkPrime, benchmarkOdd)
	}
}

func BenchmarkIsSquareTrue(b *testing.B) {
	largetrue, _ := new(big.Int).SetString("3b17f061370666c4f11db552e1dc533fbf30531421a6292207fd136a94f9f011e672a24f0ef1422210ab44f96e43599d6576030ded2b0f9c79fc8b8efd8558f09c168e35895707d7749fb92e18d9f0653efdc05daeee522204766c6aea0f2dbc5793beabbd629e69b38f5c0c56a37fd4ceb27d667ab9d1b098dae5beec2d3bfa96be55a3b9262d5662429ba76fb4f359d5674c0d861c81", 16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsSquare(largetrue)
	}
}

func BenchmarkIsSquareFalse(b *testing.B) {
	largefalse, _ := new(big.Int).SetString("3b17f061370666c4f11db552e1dc533fbf30531421a6292207fd136a94f9f011e672a24f0ef1422210ab44f96e43599d6576030ded2b0f9c79fc8b8efd8558f09c168e35895707d7749fb92e18d9f0653efdc05daeee522204766c6aea0f2dbc5793beabbd629e69b38f5c0c56a37fd4ceb27d667ab9d1b098dae5beec2d3bfa96be55a3b9262d5662429ba76fb4f359d5674c0d861d81", 16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsSquare(largefalse)
	}
}

func BenchmarkIsSquareRandomNumber(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsSquare(randBig(1024))
	}
}

func BenchmarkTrailingZeroBits(b *testing.B) {
	for i := 0; i < b.N; i++ {
		trailingZeroBits(new(big.Int).Lsh(randBig(1024), uint(rand.Intn(100))))
	}
}

// utility primality tests

func BenchmarkSmallPrimeTest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SmallPrimeTest(benchmarkPrime)
	}
}

func BenchmarkSmallPrimeTestRandom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SmallPrimeTest(randBig(1024))
	}
}

func BenchmarkStrongMillerRabin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StrongMillerRabin(benchmarkPrime, 2)
	}
}

func BenchmarkStrongMillerRabinRandom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x := randBig(1024)
		x.SetBit(x, 0, 1)
		StrongMillerRabin(x, 2)
	}
}

func BenchmarkStrongLucasSelfridgeTest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StrongLucasSelfridgeTest(benchmarkPrime)
	}
}

// primality wrappers

func BenchmarkNextPrime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NextPrime(randBig(1024))
	}
}

func BenchmarkRandPrime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandPrime(1024)
	}
}

func BenchmarkCryptoRandPrime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cryptoRand.Prime(cryptoRand.Reader, 1024)
	}
}

// benchmark primality tests

func BenchmarkBPSW(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BPSW(benchmarkPrime)
	}
}

func BenchmarkProbablyPrime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// 20 is used in testing for primality in
		// crypto/rand.Prime function:
		// https://golang.org/src/crypto/rand/util.go?#L99
		benchmarkPrime.ProbablyPrime(20)
	}
}

// random primality tests

func BenchmarkBPSWRand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BPSW(randBig(1024))
	}
}

func BenchmarkProbablyPrimeRandom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randBig(1024).ProbablyPrime(20)
	}
}
