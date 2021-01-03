package service

import socketio "github.com/googollee/go-socket.io"

type IMinerInfo interface {
	MinerInfo(s socketio.Conn, msg string)
}
