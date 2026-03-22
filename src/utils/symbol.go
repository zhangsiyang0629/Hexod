package utils

import (
	"fmt"
	"strings"
)

func OkInstId2Symbol(instId string) string {
	tmp := strings.Split(instId, "-")
	if len(tmp) == 2 {
		return fmt.Sprintf("%s_%s", tmp[0], tmp[1])
	} else if len(tmp) == 3 {
		return fmt.Sprintf("%s_%s_PERP", tmp[0], tmp[1])
	}
	panic("bad instId:" + instId)
}

func Symbol2OkInstId(symbol string) string {
	tmp := strings.Split(symbol, "_")
	if len(tmp) == 2 {
		return fmt.Sprintf("%s-%s", tmp[0], tmp[1])
	} else if len(tmp) == 3 {
		return fmt.Sprintf("%s-%s-SWAP", tmp[0], tmp[1])
	}
	panic("bad symbol:" + symbol)
}

func BaseQuote(symbol string) (string, string) {
	tmp := strings.Split(symbol, "_")
	return tmp[0], tmp[1]
}
