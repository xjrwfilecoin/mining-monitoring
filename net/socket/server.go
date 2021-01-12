package socket

import (
	"encoding/json"
	"fmt"
	"github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"mining-monitoring/log"
	"mining-monitoring/model"
	"mining-monitoring/utils"
	"net/http"
)

// todo 通用

var SServer = NewServer()

// todo
func BroadCaseMsg(obj interface{}) {
	cmd := model.ResponseCmd{Code: 1, Url: SubMinerInfo, Message: "success", Data: obj,}
	bytes, err := json.Marshal(cmd)
	if err!=nil{
		log.Error(err.Error())
	}
	log.Info("broadCast: ", string(bytes))
	SServer.broadcastMessage(DefaultNamespace, DefaultRoom, SubMinerInfo, cmd)
}

type Server struct {
	server    *socketio.Server
	namespace string
}

func (ss *Server) GetServer() *socketio.Server {
	return ss.server
}

func (ss *Server) broadcastMessage(namespace, room, event string, obj interface{}) {
	ok := ss.server.BroadcastToRoom(namespace, room, event, obj)
	if !ok {
		log.Error("broadcast msg fail ", obj)
	}
}

func (ss *Server) RegisterRouter(namespace, event string, fn func(s socketio.Conn, msg string)) {
	if namespace == "" {
		namespace = "/"
	}
	if event == "" {
		panic(fmt.Errorf("socketIo event is empty"))
	}

	ss.server.OnEvent(namespace, event, fn)
}

func (ss *Server) JoinRoom(namespace, room string, s socketio.Conn) {
	ss.server.JoinRoom(namespace, room, s)
}

func (ss *Server) Close() error {
	if ss.server != nil {
		return ss.server.Close()
	}
	return nil
}

func (ss *Server) Run() error {
	ss.server.OnConnect(ss.namespace, func(s socketio.Conn) error {
		log.Debug("socketIO client connect ", s.ID(), s.RemoteAddr(), )
		s.Emit("message", "connected ")
		return nil
	})

	ss.server.OnError(ss.namespace, func(s socketio.Conn, e error) {
		log.Error("socketIo error ", e.Error())
	})

	ss.server.OnDisconnect(ss.namespace, func(s socketio.Conn, reason string) {
		log.Debug("socketIo client disConnect ", s.ID(), reason)
	})
	return ss.server.Serve()

}

type GenId struct {
}

func (g GenId) NewID() string {

	return utils.GetUUID()
}

func NewServer() (*Server) {
	server, err := socketio.NewServer(
		&engineio.Options{
			Transports: []transport.Transport{
				&websocket.Transport{
					CheckOrigin: func(r *http.Request) bool {
						return true
					},
				},
				&polling.Transport{
					CheckOrigin: func(r *http.Request) bool {
						return true
					},
				},
			},
		})
	//server, err := socketio.NewServer(nil)
	if err != nil {
		panic(fmt.Errorf("init socket-io server %v \n", err))
	}
	return &Server{
		server:    server,
		namespace: "/",
	}

}
