package uniswapv4

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

func SortTokens(a, b common.Address) (common.Address, common.Address, bool) {
	if strings.ToLower(a.Hex()) < strings.ToLower(b.Hex()) {
		return a, b, true
	}
	return b, a, false
}

func FloatToWei(amount float64) *big.Int {
	wei := new(big.Float).Mul(big.NewFloat(amount), big.NewFloat(1e18))
	weiInt := new(big.Int)
	wei.Int(weiInt) // 转换成整数
	return weiInt
}
