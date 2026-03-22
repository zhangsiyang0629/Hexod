package router

import (
	"hexod/src/types"
)

var defaultRouter *Router

func InitRouter() {
	defaultRouter = &Router{pinM: make(map[string]*Pin, 8)}
}

func Register(symbol string, pin *Pin) {
	defaultRouter.Register(symbol, pin)
}

func UnRegister(symbol string) {
	defaultRouter.UnRegister(symbol)
}

func InputBookTicker(bookTicker *types.BookTickerEvent) {
	defaultRouter.InputBookTicker(bookTicker)
}

func InputOrderBook(orderBook *types.OrderBookEvent) {
	defaultRouter.InputOrderBook(orderBook)
}
