package router

import (
	"hexod/src/types"
)

type Pin struct {
	bCh chan *types.BookTickerEvent
	dCh chan *types.OrderBookEvent
}

func (p *Pin) SetBookTickerCh(bCh chan *types.BookTickerEvent) {
	p.bCh = bCh
}

func (p *Pin) SetOrderBookEventCh(dCh chan *types.OrderBookEvent) {
	p.dCh = dCh
}
