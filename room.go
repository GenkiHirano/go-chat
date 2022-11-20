package main

import "golang.org/x/net/websocket"

type room struct {
	// 他のクライアントに転送するためのメッセージを保持するチャネル
	forward chan []byte
	// チャットに参加しようとしているクライアントのためのチャネル
	join chan *client
	// チャットに退室しようとしているクライアントのためのチャネル
	leave   chan *client
	clients map[*client]bool
}

func (r *room) run() {
	// 無限ループで終了されるまで実行
	// 無限ループは良くないが、goroutineとして実行しているので OK!
	for {
		select {
		case client := <-r.join:
			// 参加
			r.clients[client] = true
		case client := <-r.leave:
			// 退室
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			// 全てのクライアントにメッセージを送信
			for client := range r.clients {
				select {
				case client.send <- msg:
				default:
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}

const (
	socketBufferSize = 1024
	messageBufferSize = 256
)

// TODO: 後で修正する
var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}
