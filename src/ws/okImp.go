package ws

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hexod/src/config"
	"hexod/src/constant"
	"hexod/src/router"
	"hexod/src/types"
	"hexod/src/utils"
	"strconv"
	"strings"
	"time"

	"github.com/bytedance/sonic"
)

type OkWsData struct {
	Event string `json:"event"`
	Code  string `json:"code"`
	Msg   string `json:"msg"`
	Arg   struct {
		Channel string `json:"channel"`
		InstId  string `json:"instId"`
	} `json:"arg"`
	Data   json.RawMessage `json:"data"`
	ConnId string          `json:"connId"`
}

type OkImp struct {
	isPrivate  bool
	pingTimer  *time.Timer
	debugEpoch int64
	parent     *WsClient
}

func (ok *OkImp) OnConnected(cli *WsClient, typ ConnectType) {
	log.Info("ok ws connected")
	if !ok.isPrivate {
		return
	}
	ok.pingTimer = time.NewTimer(20 * time.Second)
	go func() {
		for range ok.pingTimer.C {
			cli.WriteBytes([]byte("ping"))
			ok.pingTimer.Reset(20 * time.Second)
		}
	}()
	ok.Login()
}

func (ok *OkImp) Handle(cli *WsClient, bs []byte, messageType int) {
	ok.debugEpoch++
	if ok.debugEpoch%10000 == 0 {
		log.WithField("bs", string(bs)).Info("ok ws message")
	}

	if len(bs) == 4 && string(bs) == "pong" {
		cli.recvPongTime = time.Now()
		return
	}

	var dat OkWsData
	if err := sonic.Unmarshal(bs, &dat); err != nil {
		log.WithError(err).Error("unmarshal ok ws data failed")
		return
	}

	if (dat.Code != "" && dat.Code != "0") || dat.Event == "error" {
		err := fmt.Errorf("code:%s, msg:%s", dat.Code, dat.Msg)
		log.WithError(err).Error("ok ws data error")
		return
	}

	if dat.Event == "subscribe" {
		log.WithField("arg", dat.Arg).Info("ok subscribe success")
		return
	} else if dat.Event == "unsubscribe" {
		log.WithField("arg", dat.Arg).Info("ok unsubscribe success")
		return
	} else if dat.Event == "login" {
		log.WithField("connId", dat.ConnId).Info("ok login success")
	}

	switch dat.Arg.Channel {
	case "bbo-tbt":
		ok.onBboTbtRecv(dat.Arg.InstId, dat.Data)
	default:
		log.WithField("dat", string(dat.Data)).Warn("unknown ok message")
	}
}

type okLoginReq struct {
	OP   string       `json:"op"`
	Args []okLoginArg `json:"args"`
}

type okLoginArg struct {
	ApiKey     string `json:"apiKey"`
	Passphrase string `json:"passphrase"`
	Timestamp  string `json:"timestamp"`
	Sign       string `json:"sign"`
}

func (ok *OkImp) Login() {
	ts := time.Now().Unix()
	path := "/users/self/verify"
	method := "GET"
	message := fmt.Sprintf("%d%s%s", ts, method, path)

	signature := okWsSign([]byte(message), []byte(config.Cfg.Ok.Secret))
	req := &okLoginReq{}
	req.OP = "login"
	req.Args = append(req.Args, okLoginArg{
		ApiKey:     config.Cfg.Ok.Key,
		Passphrase: config.Cfg.Ok.Passphase,
		Timestamp:  fmt.Sprint(ts),
		Sign:       signature,
	})
	log.WithField("req", req).Info("ok ws login")
	ok.parent.Write(req)
}

func (ok *OkImp) onBboTbtRecv(instId string, dat json.RawMessage) {
	type bookTicker struct {
		Asks      [][]string `json:"asks"`
		Bids      [][]string `json:"bids"`
		Ts        string     `json:"ts"`
		Checksum  int        `json:"checksum"`
		PrevSeqID int        `json:"prevSeqId"`
		SeqID     int        `json:"seqId"`
	}

	tickers := make([]bookTicker, 0, 1)
	if err := sonic.Unmarshal(dat, &tickers); err != nil {
		log.WithError(err).Error("unmarshal ok tbt failed")
		return
	}

	if len(tickers) == 0 {
		log.Warn("empty ok tbt")
		return
	}

	var (
		ticker      = tickers[0]
		ask1        = ticker.Asks[0]
		askPrice, _ = strconv.ParseFloat(ask1[0], 64)
		askQty, _   = strconv.ParseFloat(ask1[1], 64)
		bid1        = ticker.Bids[0]
		bidPrice, _ = strconv.ParseFloat(bid1[0], 64)
		bidQty, _   = strconv.ParseFloat(bid1[1], 64)
		ts, _       = strconv.ParseInt(ticker.Ts, 10, 64)
		exchange    = constant.OkSpot
	)

	if strings.Contains(instId, "-SWAP") {
		exchange = constant.OkFuture
	}

	evt := &types.BookTickerEvent{
		Symbol:     utils.OkInstId2Symbol(instId),
		Exchange:   exchange,
		AskPrice:   askPrice,
		AskQty:     askQty,
		BidPrice:   bidPrice,
		BidQty:     bidQty,
		ExchangeTs: ts,
		TraceId:    utils.RandomString(8),
		Ts:         utils.Millisec(time.Now()),
	}
	router.InputBookTicker(evt)
}

func okWsSign(body []byte, secretKey []byte) string {
	mac := hmac.New(sha256.New, secretKey)
	mac.Write(body)
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
