package prime

import "math/big"

var (
	// diffs to next relatively prime number mod 210 = 2*3*5*7
	// so gcd(210,diff210[i] + i) = 1
	diffs210 = []int{
		11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 2, 1, 4, 3, 2, 1, 2, 1, 4, 3, 2, 1, 6,
		5, 4, 3, 2, 1, 2, 1, 6, 5, 4, 3, 2, 1, 4, 3, 2, 1, 2, 1, 4, 3, 2, 1, 6, 5,
		4, 3, 2, 1, 6, 5, 4, 3, 2, 1, 2, 1, 6, 5, 4, 3, 2, 1, 4, 3, 2, 1, 2, 1, 6,
		5, 4, 3, 2, 1, 4, 3, 2, 1, 6, 5, 4, 3, 2, 1, 8, 7, 6, 5, 4, 3, 2, 1, 4, 3,
		2, 1, 2, 1, 4, 3, 2, 1, 2, 1, 4, 3, 2, 1, 8, 7, 6, 5, 4, 3, 2, 1, 6, 5, 4,
		3, 2, 1, 4, 3, 2, 1, 6, 5, 4, 3, 2, 1, 2, 1, 4, 3, 2, 1, 6, 5, 4, 3, 2, 1,
		2, 1, 6, 5, 4, 3, 2, 1, 6, 5, 4, 3, 2, 1, 4, 3, 2, 1, 2, 1, 4, 3, 2, 1, 6,
		5, 4, 3, 2, 1, 2, 1, 6, 5, 4, 3, 2, 1, 4, 3, 2, 1, 2, 1, 4, 3, 2, 1, 2, 1,
		10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 2}
	// all primes < 10 bits long
	primes10 = []uint16{
		2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61,
		67, 71, 73, 79, 83, 89, 97, 101, 103, 107, 109, 113, 127, 131, 137,
		139, 149, 151, 157, 163, 167, 173, 179, 181, 191, 193, 197, 199, 211,
		223, 227, 229, 233, 239, 241, 251, 257, 263, 269, 271, 277, 281, 283,
		293, 307, 311, 313, 317, 331, 337, 347, 349, 353, 359, 367, 373, 379,
		383, 389, 397, 401, 409, 419, 421, 431, 433, 439, 443, 449, 457, 461,
		463, 467, 479, 487, 491, 499, 503, 509, 521, 523, 541, 547, 557, 563,
		569, 571, 577, 587, 593, 599, 601, 607, 613, 617, 619, 631, 641, 643,
		647, 653, 659, 661, 673, 677, 683, 691, 701, 709, 719, 727, 733, 739,
		743, 751, 757, 761, 769, 773, 787, 797, 809, 811, 821, 823, 827, 829,
		839, 853, 857, 859, 863, 877, 881, 883, 887, 907, 911, 919, 929, 937,
		941, 947, 953, 967, 971, 977, 983, 991, 997, 1009, 1013, 1019, 1021}
	// all primes < 10 bits and their product
	prodPrimes10, _ = new(big.Int).SetString("613a0497aa700632594668d2175f6874157ab081f7d649a3e936c6608f20575cb03949974ef1fb62db814d5fdf2c0d0e2d0abb2b26e8cc08403e32336e4bf96f1ffa1b71d1f4c342dc3812e17d7035b9e93905bff2c1a6de", 16)
	// all squares mod 1155 = 3 * 5 * 7 * 11
	smallSquareMod = big.NewInt(3 * 5 * 7 * 11)
	smallSquares   = []uint16{
		0, 1, 4, 9, 15, 16, 25, 36, 49, 60, 64, 70, 81, 91, 99, 100,
		114, 121, 126, 130, 135, 141, 144, 154, 165, 169, 190, 196, 210,
		214, 225, 231, 235, 240, 246, 256, 280, 289, 291, 295, 301, 309,
		319, 324, 330, 331, 345, 361, 364, 366, 375, 379, 385, 394, 396,
		399, 400, 421, 429, 441, 445, 456, 466, 471, 484, 499, 504, 511,
		520, 526, 529, 540, 550, 555, 561, 564, 576, 595, 606, 609, 610,
		616, 625, 630, 631, 639, 660, 669, 676, 694, 709, 715, 729, 730,
		735, 751, 760, 771, 774, 781, 784, 786, 795, 814, 819, 826, 834,
		840, 841, 856, 861, 870, 889, 891, 900, 924, 925, 939, 940, 946,
		949, 960, 961, 966, 984, 991, 994, 1005, 1015, 1024, 1026, 1045,
		1050, 1054, 1059, 1065, 1071, 1089, 1101, 1114, 1120, 1131, 1134,
		1149}
)
