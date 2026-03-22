package types

import (
	"fmt"
	"hash/fnv"
	"hexod/src/constant"
)

type BookTickerEvent struct {
	Symbol     string
	Exchange   constant.ExchangeType
	AskPrice   float64
	AskQty     float64
	BidPrice   float64
	BidQty     float64
	ExchangeTs int64
	Ts         int64
	TraceId    string
}

func (e *BookTickerEvent) Hash() uint32 {
	h := fnv.New32a()
	s := fmt.Sprintf("%f%f%f%f%d", e.AskPrice, e.AskQty, e.BidPrice, e.BidQty, e.ExchangeTs)
	h.Write([]byte(s))
	return h.Sum32()
}

type OrderBookItem struct {
	Price       float64
	Qty         float64
	OriginQty   string
	OriginPrice string
}

type OrderBookEvent struct {
	Symbol     string
	Exchange   constant.ExchangeType
	Asks       []OrderBookItem
	Bids       []OrderBookItem
	ExchangeTs int64
	Ts         int64
	TraceId    string
}

type KlineEvent struct {
	Symbol   string
	Open     float64
	High     float64
	Low      float64
	Close    float64
	Vol      float64
	VolQuote float64
	Confirm  bool
	Ts       int64
}

type EventType int

const (
	CLOSE_PLACE_ORDER EventType = iota
	KLINE
	RESET_KLINE_FETCH
)

type GeneralEvent struct {
	Symbol    string
	Exchange  constant.ExchangeType
	EventType EventType
	EventDat  interface{}
}

type CommonEvent GeneralEvent
