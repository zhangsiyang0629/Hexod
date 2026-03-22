package constant

import "fmt"

const (
	Exchange_PionexSpot   = "pionexSpot"
	Exchange_OkSpot       = "okSpot"
	Exchange_OkFuture     = "okFuture"
	Exchange_BiSpot       = "biSpot"
	Exchange_GateSpot     = "gateSpot"
	Exchange_BitgetSpot   = "bitgetSpot"
	Exchange_BitgetFuture = "bitgetFuture"
	Exchange_BitmartSpot  = "bitmartSpot"
)

type ExchangeType int

func (e ExchangeType) Name() string {
	switch e {
	case PionexSpot:
		return Exchange_PionexSpot
	case OkSpot:
		return Exchange_OkSpot
	case OkFuture:
		return Exchange_OkFuture
	case BiSpot:
		return Exchange_BiSpot
	case GateSpot:
		return Exchange_GateSpot
	case BitgetSpot:
		return Exchange_BitgetSpot
	case BitgetFuture:
		return Exchange_BitgetFuture
	case BitmartSpot:
		return Exchange_BitmartSpot
	}
	return "unknown"
}

const (
	PionexSpot ExchangeType = iota
	OkSpot
	OkFuture
	GateSpot
	BiSpot
	BitgetSpot
	BitgetFuture
	BitmartSpot
)

func MustConverToExchangeType(name string) ExchangeType {
	switch name {
	case Exchange_PionexSpot:
		return PionexSpot
	case Exchange_OkSpot:
		return OkSpot
	case Exchange_OkFuture:
		return OkFuture
	case Exchange_BiSpot:
		return BiSpot
	case Exchange_GateSpot:
		return GateSpot
	case Exchange_BitgetSpot:
		return BitgetSpot
	case Exchange_BitgetFuture:
		return BitgetFuture
	case Exchange_BitmartSpot:
		return BitmartSpot
	}
	err := fmt.Errorf("unknonw exchange name:%s", name)
	panic(err)
}
