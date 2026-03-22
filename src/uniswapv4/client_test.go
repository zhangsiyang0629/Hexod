package uniswapv4

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestGetPrice(t *testing.T) {
	rpcUrl := "wss://sepolia.infura.io/ws/v3/e2f51a4b6ed64f0c88ddf637e3c0cb2b"
	quoter := "0x61b3f2011a92d183c7dbadbda940a7555ccf9227"
	cli := NewClient(rpcUrl, V4CliConfig{QuotaAddr: quoter})

	tokenA := common.HexToAddress("0x7cc9c9f858ae836361a408b66d1198652fcb07b4") // ETH
	tokenB := common.HexToAddress("0x02b3d7aad00aefdb8f3fe1d5e3fdb7ed78cc41a3") // USDT
	token0, token1, normal := SortTokens(tokenA, tokenB)
	t.Log(normal)
	t.Log(token0)
	t.Log(token1)

	key := PoolKey{
		Token0:      token0,
		Token1:      token1,
		Fee:         big.NewInt(3000), // 0.3%
		TickSpacing: big.NewInt(30),
		Hooks:       common.Address{},
	}
	price := cli.GetPrice(key)
	if normal {
		price = 1 / price
	}
	t.Log(price)
}

func TestSwapCheck(t *testing.T) {
	rpcUrl := "wss://sepolia.infura.io/ws/v3/e2f51a4b6ed64f0c88ddf637e3c0cb2b"
	addr := "0xE03A1074c86CFeDd5C142C4F04F1a1536e203543"
	tokenAddr := "0x02b3d7aad00aefdb8f3fe1d5e3fdb7ed78cc41a3"
	senderAddr := "0x599Cbb31cD468Df8022502e69710b01cCf53EfA4"
	cli := NewClient(rpcUrl, V4CliConfig{PollAddr: addr, ChainID: 11155111})
	err := cli.SwapCheck("24554b3e74e003207cf94ffeb1542a9484b96b8833e6a30e128ef38c95514311", tokenAddr, senderAddr)
	if err != nil {
		t.Error(err)
	}
}

func TestSwap(t *testing.T) {
	rpcUrl := "wss://sepolia.infura.io/ws/v3/e2f51a4b6ed64f0c88ddf637e3c0cb2b"
	addr := "0xE03A1074c86CFeDd5C142C4F04F1a1536e203543"
	cli := NewClient(rpcUrl, V4CliConfig{PollAddr: addr, ChainID: 11155111})
	tokenAddr0 := "0x02b3d7aad00aefdb8f3fe1d5e3fdb7ed78cc41a3"
	tokenAddr1 := "0x7cc9c9f858ae836361a408b66d1198652fcb07b4"
	key := PoolKey{
		Token0:      common.HexToAddress(tokenAddr0),
		Token1:      common.HexToAddress(tokenAddr1),
		Fee:         big.NewInt(3000), // 0.3%
		TickSpacing: big.NewInt(30),
		Hooks:       common.Address{},
	}
	err := cli.Swap("24554b3e74e003207cf94ffeb1542a9484b96b8833e6a30e128ef38c95514311", key, true,
		tokenAddr0, tokenAddr1, 1)
	if err != nil {
		t.Error(err)
	}
}
