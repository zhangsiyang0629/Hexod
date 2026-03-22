// uniswapv4/abi.go
package uniswapv4

const PoolManagerQuoteABI = `[
        {
            "inputs": [
                {
                    "components": [
                        {
                            "components": [
                                {"internalType":"address","name":"token0","type":"address"},
                                {"internalType":"address","name":"token1","type":"address"},
                                {"internalType":"uint24","name":"fee","type":"uint24"},
                                {"internalType":"int24","name":"tickSpacing","type":"int24"},
								{"internalType":"address","name":"hooks","type":"address"}
                            ],
                            "internalType":"struct PoolKey",
                            "name":"poolKey",
                            "type":"tuple"
                        },
                        {"internalType":"bool","name":"zeroForOne","type":"bool"},
                        {"internalType":"uint128","name":"exactAmount","type":"uint128"},
                        {"internalType":"bytes","name":"hookData","type":"bytes"}
                    ],
                    "name":"params",
                    "type":"tuple"
                }
            ],
            "name":"quoteExactInputSingle",
            "outputs":[{"internalType":"uint256","name":"amountOut","type":"uint256"}],
            "stateMutability":"view",
            "type":"function"
        }
    ]`

const SwapABI = `[
  {
    "name": "swap",
    "type": "function",
    "stateMutability": "nonpayable",
    "inputs": [
      {
        "name": "key",
        "type": "tuple",
        "components": [
          	{"internalType":"address","name":"token0","type":"address"},
        	{"internalType":"address","name":"token1","type":"address"},
            {"internalType":"uint24","name":"fee","type":"uint24"},
            {"internalType":"int24","name":"tickSpacing","type":"int24"},
			{"internalType":"address","name":"hooks","type":"address"}
        ]
      },
      {
        "name": "params",
        "type": "tuple",
        "components": [
			{"internalType":"int256","name":"amountSpecified","type":"int256"},
			{"internalType":"int24","name":"tickSpacing","type":"int24"},
			{"internalType":"bool","name":"zeroForOne","type":"bool"},
			{"internalType":"uint160","name":"sqrtPriceLimitX96","type":"uint160"},
			{"internalType":"int24","name":"lpFeeOverride","type":"int24"}
        ]
      },
      {
        "name": "hookData",
        "type": "bytes"
      }
    ],
    "outputs": [
      {
        "name": "swapDelta",
        "type": "tuple",
        "components": [
          {
            "name": "amount0",
            "type": "int256"
          },
          {
            "name": "amount1",
            "type": "int256"
          }
        ]
      }
    ]
  },
  {
    "name": "approve",
    "type": "function",
    "inputs": [
      { "name": "spender", "type": "address" },
      { "name": "amount", "type": "uint256" }
    ],
    "outputs": [
      { "name": "", "type": "bool" }
    ]
  }
]`

const SwapCheckABI = `[
  {
    "constant": true,
    "inputs": [
      {
        "name": "account",
        "type": "address"
      }
    ],
    "name": "balanceOf",
    "outputs": [
      {
        "name": "",
        "type": "uint256"
      }
    ],
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [
      {
        "name": "owner",
        "type": "address"
      },
      {
        "name": "spender",
        "type": "address"
      }
    ],
    "name": "allowance",
    "outputs": [
      {
        "name": "",
        "type": "uint256"
      }
    ],
    "type": "function"
  }
]`
