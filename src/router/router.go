package router

import (
	"fmt"
	"hexod/src/types"
	"hexod/src/utils"
	"sync"
	"time"
)

const (
	CHANNEL_BLOCK_TIMEOUT = 300 * time.Millisecond
)

/*
每个事件一个ch，并没有涉及一个通用的event interface，
其目的就是不想断言，断言有一定的性能损耗
*/
type Router struct {
	pinMtx sync.RWMutex
	pinM   map[string]*Pin
}

func (r *Router) InputBookTicker(bookTicker *types.BookTickerEvent) {
	symbol := bookTicker.Symbol
	r.pinMtx.RLock()
	pin, ok := r.pinM[symbol]
	r.pinMtx.RUnlock()

	if !ok {
		err := fmt.Errorf("unknown %s strategy", symbol)
		log.WithField("symbol", symbol).WithError(err).Error("input book ticker error")
		return
	}

	timeoutTimer := utils.NewTimer(CHANNEL_BLOCK_TIMEOUT)
	defer utils.FreeTimer(timeoutTimer)
	select {
	case pin.bCh <- bookTicker:
	case <-timeoutTimer.C:
		log.WithField("bookTicker", *bookTicker).Error("input book ticker timeout")
	}
}

func (r *Router) InputOrderBook(orderbook *types.OrderBookEvent) {
	symbol := orderbook.Symbol
	r.pinMtx.RLock()
	pin, ok := r.pinM[symbol]
	r.pinMtx.RUnlock()
	if !ok {
		err := fmt.Errorf("unknown %s strategy", symbol)
		log.WithError(err).Error("iinput orderbook error")
		return
	}

	timeoutTimer := utils.NewTimer(CHANNEL_BLOCK_TIMEOUT)
	defer utils.FreeTimer(timeoutTimer)
	select {
	case pin.dCh <- orderbook:
	case <-timeoutTimer.C:
		log.WithField("orderBook", *orderbook).Error("input depth timeout")
	}
}

func (r *Router) Register(symbol string, pin *Pin) {
	r.pinMtx.Lock()
	defer r.pinMtx.Unlock()
	r.pinM[symbol] = pin
}

func (r *Router) UnRegister(symbol string) {
	log.WithField("symbol", symbol).Info("router unregister")
	r.pinMtx.Lock()
	defer r.pinMtx.Unlock()
	delete(r.pinM, symbol)
}
