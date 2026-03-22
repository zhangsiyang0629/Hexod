package uniswapv3

import (
	"context"
	"testing"
	"time"
)

func TestPricer(t *testing.T) {
	rpcUrl := "wss://sepolia.infura.io/ws/v3/e2f51a4b6ed64f0c88ddf637e3c0cb2b"
	poolAddress := "0xB0223eA2626d7BEe26aB0F0Fd0aD9417012719c4"
	cli := NewClient(rpcUrl, poolAddress, PoolSlotABI)
	px := &Pricer{
		cli: cli,
		ctx: context.Background(),
	}
	px.Start()
	for range 10 {
		t.Log(px.GetPrice())
		time.Sleep(500 * time.Millisecond)
	}
}
