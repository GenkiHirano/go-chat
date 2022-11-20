package main

import (
	"golang.org/x/net/websocket"
)

// client チャットを行っている1人のユーザーを表す
type client struct {
	socket *websocket.Conn
	// メッセージが送られるチャネル
	// 受信したメッセージをwebsocketを通じてユーザーのブラウザに送られるのを待機する
	send chan []byte
	room *room
}

func (c *client) read() {
	for {
		// TODO: 後ほど定義する
		var msg *message
		if _, err := c.socket.Read(&msg); err == nil {
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		_, err := c.socket.Write(msg)
		if err != nil {
			break
		}
	}
	c.socket.Close()
}
