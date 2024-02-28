package entity

import (
	"github.com/gorilla/websocket"
)

// WebSocketWrapper は WebSocket の操作をラップするためのインターフェースです。
type WebSocketWrapper interface {
	WriteMessage(messageType int, data []byte) error
}

// WebSocketConn は WebSocket コネクションをラップする構造体です。
type WebSocketConn struct {
	Conn *websocket.Conn
}

// WriteMessage は WebSocket コネクションにメッセージを書き込むためのメソッドです。
func (w *WebSocketConn) WriteMessage(messageType int, data []byte) error {
	return w.Conn.WriteMessage(messageType, data)
}

//func (w *WebSocketConn) ReadMessage() (messageType int, p []byte, err error) {
//	var r io.Reader
//	messageType, r, err = c.NextReader()
//	if err != nil {
//		return messageType, nil, err
//	}
//	p, err = io.ReadAll(r)
//	return messageType, p, err
//}
