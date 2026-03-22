package utils

import (
	"math"
	"math/big"
)

// sqrtPriceX96 → price (token1 / token0)
func SqrtPriceX96ToPrice(sqrtPriceX96 *big.Int) float64 {
	f := new(big.Float).SetInt(sqrtPriceX96)

	// 除以 2^96
	denom := new(big.Float).SetFloat64(math.Pow(2, 96))
	f.Quo(f, denom)

	// 平方
	f.Mul(f, f)

	price, _ := f.Float64()
	return price
}
