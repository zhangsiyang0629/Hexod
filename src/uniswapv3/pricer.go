package uniswapv3

import (
	"context"
	"hexod/src/utils"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Price struct {
	Value float64
	Ts    int64
}

type Pricer struct {
	cli      *Client
	latestPx atomic.Value // float64
	ctx      context.Context
	cancel   context.CancelFunc
}

func (p *Pricer) Start() {
	go p.watchSwap()  // 主驱动
	go p.watchBlock() // fallback
	go p.pollSlot0()  // 保底
}

func (p *Pricer) watchSwap() {
	logs := make(chan types.Log)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{p.cli.addr},
	}
	sub, err := p.cli.rpc.SubscribeFilterLogs(p.ctx, query, logs)
	if err != nil {
		log.WithField("cli", p.cli.ToString()).WithError(err).Error("sub swap event error")
		return
	}

	for {
		select {
		case vLog := <-logs:
			price := p.decodeSwapPrice(vLog)
			p.setPrice(price)

		case err := <-sub.Err():
			log.WithError(err).Error("swap sub dropped")
			return
		}
	}
}

func (p *Pricer) decodeSwapPrice(vLog types.Log) float64 {
	// sqrtPriceX96 在 topics/data 中，需要 ABI decode
	data := vLog.Data

	// slot:
	// amount0
	// amount1
	// sqrtPriceX96 (第3个)
	sqrtPriceX96 := new(big.Int).SetBytes(data[64:96])
	return utils.SqrtPriceX96ToPrice(sqrtPriceX96)
}

func (p *Pricer) watchBlock() {
	headers := make(chan *types.Header)
	sub, err := p.cli.rpc.SubscribeNewHead(p.ctx, headers)
	if err != nil {
		log.WithField("cli", p.cli.ToString()).WithError(err).Error("sub block event error")
		return
	}

	for {
		select {
		case <-headers:
			price, err := p.cli.GetPrice()
			if err != nil {
				log.WithError(err).Error("watchBlock:get price error")
			} else {
				p.setPrice(price)
			}

		case err := <-sub.Err():
			log.WithError(err).Error("block sub dropped")
			return
		}
	}
}

func (p *Pricer) pollSlot0() {
	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-ticker.C:
			price, err := p.cli.GetPrice()
			if err != nil {
				log.WithError(err).Error("pollSlot0:get price error")
			} else {
				p.setPrice(price)
			}
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *Pricer) setPrice(px float64) {
	now := time.Now().UnixNano()
	old := p.latestPx.Load()
	if old != nil {
		if now <= old.(Price).Ts {
			return
		}
	}
	p.latestPx.Store(Price{px, now})
}

func (p *Pricer) GetPrice() float64 {
	v := p.latestPx.Load()
	if v == nil {
		return 0
	}
	return v.(Price).Value
}
