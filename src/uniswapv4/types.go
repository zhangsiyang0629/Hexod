// uniswapv4/types.go
package uniswapv4

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type PoolKey struct {
	Token0      common.Address
	Token1      common.Address
	Fee         *big.Int
	TickSpacing *big.Int
	Hooks       common.Address
}

type QuoteExactSingleParams struct {
	PoolKey     PoolKey
	ZeroForOne  bool
	ExactAmount *big.Int // uint128 → *big.Int
	HookData    []byte
}

type SwapParams struct {
	AmountSpecified   *big.Int // uint128 → *big.Int
	TickSpacing       *big.Int
	ZeroForOne        bool
	SqrtPriceLimitX96 *big.Int
	LpFeeOverride     *big.Int
}

type BalanceDelta struct {
	Amount0 *big.Int
	Amount1 *big.Int
}
