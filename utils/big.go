package utils

import "math/big"

func BigRatFromString(value string) (*big.Rat, bool) {
	return new(big.Rat).SetString(value)
}

func BigRatFromFloat(a float64) *big.Rat {
	return new(big.Rat).SetFloat64(a)
}

func BigRatFromInt(a int64) *big.Rat {
	return new(big.Rat).SetInt64(a)
}

func BigRatDiv(a, b *big.Rat) *big.Rat {
	ta := new(big.Rat).Set(a)
	tb := new(big.Rat).Set(b)
	return ta.Quo(ta, tb)
}

func BigRatMul(a, b *big.Rat) *big.Rat {
	ta := new(big.Rat).Set(a)
	tb := new(big.Rat).Set(b)
	return ta.Mul(ta, tb)
}

func BigRatAdd(a, b *big.Rat) *big.Rat {
	ta := new(big.Rat).Set(a)
	tb := new(big.Rat).Set(b)
	return ta.Add(ta, tb)
}

func BigRatSub(a, b *big.Rat) *big.Rat {
	ta := new(big.Rat).Set(a)
	tb := new(big.Rat).Set(b)
	return ta.Sub(ta, tb)
}

func BigFloatFromString(a string) (*big.Float, bool) {
	return new(big.Float).SetString(a)
}
func BigFloatFromFloat(a float64) *big.Float {
	return new(big.Float).SetFloat64(a)
}

func BigFloatMul(a, b *big.Float) *big.Float {
	ta := new(big.Float).Set(a)
	tb := new(big.Float).Set(b)
	return ta.Mul(ta, tb)
}

func BigFloatAdd(a, b *big.Float) *big.Float {
	ta := new(big.Float).Set(a)
	tb := new(big.Float).Set(b)
	return ta.Add(ta, tb)
}

func BigFloatSub(a, b *big.Float) *big.Float {
	ta := new(big.Float).Set(a)
	tb := new(big.Float).Set(b)
	return ta.Sub(ta, tb)
}



func BigIntFromString(a string) (*big.Int, bool) {
	return new(big.Int).SetString(a,10)
}


func BigIntFromInt(a int64) *big.Int {
	return big.NewInt(a)
}

func BigIntMul(a, b *big.Int) *big.Int {
	ta := new(big.Int).Set(a)
	tb := new(big.Int).Set(b)
	return ta.Mul(ta, tb)
}

func BigIntDiv(a, b *big.Int) *big.Int {
	ta := new(big.Int).Set(a)
	tb := new(big.Int).Set(b)
	return ta.Div(ta, tb)
}

func BigIntAdd(a, b *big.Int) *big.Int {
	ta := new(big.Int).Set(a)
	tb := new(big.Int).Set(b)
	return ta.Add(ta, tb)
}

func BigIntSub(a, b *big.Int) *big.Int {
	ta := new(big.Int).Set(a)
	tb := new(big.Int).Set(b)
	return ta.Sub(ta, tb)
}
