package ws

import (
	"fmt"
	"hexod/src/constant"
	"hexod/src/utils"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type WsClient struct {
	url      string
	conn     *websocket.Conn
	wch      chan []byte
	imp      WsImp
	exchange constant.ExchangeType
	priv     bool

	recvPingTime time.Time
	recvPongTime time.Time
	pingInterval time.Duration
	pongTimeout  time.Duration

	mutex  sync.Mutex
	quit   chan struct{}
	closed bool
	epoch  int64
}

func (ws *WsClient) SetPingInterval(t time.Duration) {
	ws.pingInterval = t
}

func (ws *WsClient) SetpPongTimeout(t time.Duration) {
	ws.pongTimeout = t
}

// 用于监控ws是否断开
func (ws *WsClient) WatchClosed() {
	<-ws.quit
}

func (ws *WsClient) Dial(typ ConnectType) error {
	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: time.Second * 10,
	}
	conn, _, err := dialer.Dial(ws.url, nil)
	if err != nil {
		return fmt.Errorf("ws.Dial:%v", err)
	}

	now := time.Now()
	ws.conn = conn
	ws.closed = false
	ws.recvPingTime = now
	ws.recvPongTime = now
	ws.quit = make(chan struct{})
	ws.conn.SetPingHandler(func(message string) error {
		ws.recvPingTime = time.Now()
		return conn.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(time.Second*10))
	})

	ws.conn.SetPongHandler(func(message string) error {
		ws.recvPongTime = time.Now()
		return nil
	})

	atomic.AddInt64(&ws.epoch, 1)
	go ws.readLoop()
	go ws.pingLoop()
	go ws.writeLoop()
	ws.imp.OnConnected(ws, typ)
	return nil
}

func (ws *WsClient) Close() {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	if ws.closed {
		log.Warn("already closed")
		return
	}
	log.WithFields(logrus.Fields{
		"exchange": ws.exchange.Name(),
		"priv":     ws.priv,
	}).Warn("ws closed")
	ws.closed = true
	ws.conn.Close()
	close(ws.quit)
}

func (ws *WsClient) pingLoop() {
	defer ws.Close()

	epoch := ws.epoch
	for !ws.closed || epoch != atomic.LoadInt64(&ws.epoch) {
		if time.Since(ws.recvPongTime) > ws.pongTimeout {
			log.Info("Recv PONG timeout")
			return
		}

		select {
		case <-time.After(ws.pingInterval):
		case <-ws.quit:
			log.Info("pingLoop quit")
			return
		}

		now := time.Now()
		bs := []byte(fmt.Sprintf("%d", utils.Millisec(now)))
		if err := ws.conn.WriteControl(websocket.PingMessage, bs, now.Add(time.Second*10)); err != nil {
			log.WithError(err).Errorln("control ping failed")
			return
		}
	}
}

func (ws *WsClient) readLoop() {
	conn := ws.conn
	log.Println("Start WS read loop")
	epoch := ws.epoch

	defer ws.Close()
	for !ws.closed || epoch != atomic.LoadInt64(&ws.epoch) {
		deadline := time.Now().Add(time.Second * 120)
		conn.SetReadDeadline(deadline)
		messageType, body, err := conn.ReadMessage()
		if err != nil {
			log.WithError(err).Errorf("websocket conn read timeout in 120 seconds")
			break
		}
		func() {
			defer func() {
				if err := recover(); err != nil {
					log.WithField("body", string(body)).WithField("err", err).Error("handle panic")
				}
			}()
			ws.imp.Handle(ws, body, messageType)
		}()
	}
}

func (ws *WsClient) writeLoop() {
	log.Println("Start WS write loop")
	epoch := ws.epoch

	defer ws.Close()
	for !ws.closed || epoch != atomic.LoadInt64(&ws.epoch) {
		select {
		case <-ws.quit:
			log.Info("writeLoop quit")
			return
		case bs := <-ws.wch:
			if err := ws.conn.WriteMessage(websocket.TextMessage, bs); err != nil {
				log.WithError(err).Errorln("write failed")
				return
			}
		}
	}
}

func (ws *WsClient) Write(req interface{}) error {
	bs, err := sonic.Marshal(req)
	if err != nil {
		return err
	}
	ws.wch <- bs
	return nil
}

func (ws *WsClient) WriteBytes(bs []byte) {
	ws.wch <- bs
}

func newCli(url string, imp WsImp) *WsClient {
	return &WsClient{
		url:          url,
		imp:          imp,
		wch:          make(chan []byte, 1024),
		pingInterval: 20 * time.Second,
		pongTimeout:  30 * time.Second,
	}
}
