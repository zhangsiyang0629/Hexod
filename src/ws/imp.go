package ws

type WsImp interface {
	OnConnected(*WsClient, ConnectType)
	Handle(*WsClient, []byte, int)
}
