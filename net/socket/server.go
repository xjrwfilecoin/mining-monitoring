package socket

import (
	"fmt"
	"github.com/googollee/go-socket.io"
)

type Options struct {
	namespace string
	event     string
	root      string
}

type Server struct {
	server    *socketio.Server
	options   *Options
	namespace string
	room      string
}


func (ss *Server) Close() error {
	return ss.server.Close()
}

func (ss *Server) Run() error {
	ss.server.OnConnect(ss.namespace, func(s socketio.Conn) error {
		fmt.Printf("socket connected  %v \n", s.ID())
		s.Join(ss.room)
		return nil
	})
	ss.server.OnEvent(ss.namespace, "message", func(s socketio.Conn, msg string) {
		fmt.Printf("socket on event %v \n", msg)
		s.Emit("reply", "have "+msg)
	})
	ss.server.OnError(ss.namespace, func(s socketio.Conn, e error) {
		fmt.Printf("socket on error %v \n", e)
	})

	ss.server.OnDisconnect(ss.namespace, func(s socketio.Conn, reason string) {
		fmt.Printf("socket disconnect %v  %v \n", s.ID(), reason)
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
	}, nil

}
