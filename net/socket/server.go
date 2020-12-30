package socket

import (
	"fmt"
	"github.com/googollee/go-socket.io"
	"mining-monitoring/log"
)

type Options struct {
	namespace string
	event     string
	root      string
}

// todo 通用
type Server struct {
	server    *socketio.Server
	options   *Options
	namespace string
	event     string
	room      string
}

func (ss *Server) GetServer() *socketio.Server{
	return ss.server
}


func (ss *Server) BroadcastMsg(obj interface{}) {
	ss.server.BroadcastToRoom(ss.namespace, ss.room, ss.event)
}

func (ss *Server) Close() error {
	if ss.server != nil {
		return ss.server.Close()
	}
	return nil
}

func (ss *Server) Run() error {
	ss.server.OnConnect(ss.namespace, func(s socketio.Conn) error {
		log.Debug("websocket client connect ", s.ID(), s.LocalAddr(),)
		s.Emit("message","test")
		s.Join(ss.room)
		return nil
	})

	ss.server.OnEvent(ss.namespace, ss.event, func(s socketio.Conn, msg string) {
		log.Debug("socketIo onEvent ", msg)
		s.Emit("reply", "have "+msg)
	})
	ss.server.OnError(ss.namespace, func(s socketio.Conn, e error) {
		log.Error("socketIo error ", s.ID(), e.Error())
	})

	ss.server.OnDisconnect(ss.namespace, func(s socketio.Conn, reason string) {
		log.Debug("socketIo client disConnect ", s.ID(), reason)
	})
	return ss.server.Serve()

}

func NewServer() (*Server, error) {
	server, err := socketio.NewServer(nil)
	if err != nil {
		return nil, fmt.Errorf("init socket-io server %v \n", err)
	}
	return &Server{
		server:    server,
		namespace: "/",
		room:      "lotus-miner",
		event:     "message",
	}, nil

}
