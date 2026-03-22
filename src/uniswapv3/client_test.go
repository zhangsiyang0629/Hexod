package uniswapv3

import "testing"

func TestGetPrice(t *testing.T) {
	rpcUrl := "https://sepolia.infura.io/v3/e2f51a4b6ed64f0c88ddf637e3c0cb2b"
	poolAddress := "0xB0223eA2626d7BEe26aB0F0Fd0aD9417012719c4"
	cli := NewClient(rpcUrl, poolAddress, PoolSlotABI)
	t.Log(cli.GetPrice())
}
