package ws

import (
	"fmt"
	"hexod/src/constant"
	"hexod/src/utils"
	"time"
)

type okSubscribeReq struct {
	OP   string           `json:"op"`
	Args []okSubscribeArg `json:"args"`
}

type okSubscribeArg struct {
	Channel  string `json:"channel"`
	Symbol   string `json:"instId"`
	InstType string `json:"instType,omitempty"`
}

func SubSubscribeOkGeneral(cli *WsClient, symbols []string, topics []string, instType ...string) error {
	for _, s := range symbols {
		symbol := s
		for _, t := range topics {
			topic := t
			if _, err := doSubSubscribeOkGeneral(cli, symbol, topic, instType...); err != nil {
				return err
			}
			time.Sleep(100 * time.Millisecond)
		}
		time.Sleep(800 * time.Millisecond)
	}
	return nil
}

func doSubSubscribeOkGeneral(cli *WsClient, symbol string,
	topic string, instType ...string) (*okSubscribeReq, error) {
	sub := &okSubscribeReq{}
	sub.OP = "subscribe"
	arg := okSubscribeArg{
		Channel: topic,
		Symbol:  utils.Symbol2OkInstId(symbol),
	}
	if len(instType) != 0 {
		arg.InstType = instType[0]
	}
	sub.Args = append(sub.Args, arg)
	log.WithField("req", sub).Info("subscribe ok general")
	if err := cli.Write(sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func UnSubscribeOk(cli *WsClient, symbol string, topic string, InstType ...string) error {
	type subscribeReq struct {
		OP   string           `json:"op"`
		Args []okSubscribeArg `json:"args"`
	}

	sub := &subscribeReq{}
	sub.OP = "unsubscribe"

	var channel string
	switch topic {
	case constant.TOPIC_KLINE5M:
		channel = "candle5m"
	case constant.TOPIC_BOOKTICKER:
		channel = "bbo-tbt"
	case constant.TOPIC_ORDER:
		channel = "orders"
	default:
		return fmt.Errorf("unsupported unscribte ok topic: %s", topic)
	}
	arg := okSubscribeArg{
		Channel: channel,
		Symbol:  utils.Symbol2OkInstId(symbol),
	}
	if len(InstType) != 0 {
		arg.InstType = InstType[0]
	}
	sub.Args = append(sub.Args, arg)
	return cli.Write(sub)
}
