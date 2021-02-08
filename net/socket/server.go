package socket

import (
	"fmt"
	"github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"mining-monitoring/log"
	"mining-monitoring/utils"
	"net/http"
	"sync/atomic"
)

var SServer = NewServer()

func BroadCaseMsg(namespace, room, event string, obj interface{}) {
	cmd := ResponseCmd{Code: 1, Url: event, Message: "success", Body: obj,}
	SServer.broadcastMessage(namespace, room, event, cmd)
}

type Server struct {
	server     *socketio.Server
	namespace  string
	connMaxNum int64 // 连接总数限制
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

func (ss *Server) RegisterRouterV1(namespace, event string, fn func(c *Context)) {
	if namespace == "" {
		namespace = "/"
	}
	if event == "" {
		panic(fmt.Errorf("socketIo event is empty"))
	}

	ss.server.OnEvent(namespace, event, func(s socketio.Conn, data string) {
		log.Debug("client request: ", s.ID(), s.RemoteAddr(), data)
		uri := utils.GetJsonValue(data, "uri")
		rEvent := utils.GetJsonValue(data, "event")
		body := utils.GetJsonValue(data, "body")
		msgId := utils.GetJsonValue(data, "msgId")
		if (uri == "" && rEvent == "") || msgId == "" {
			s.Emit(rEvent, NewFailResp("Uri or MsgId is empty"))
		} else {
			tempUri := uri
			if uri == "" {
				tempUri = rEvent
			}
			context := NewContext(s, tempUri, msgId, body)
			fn(context)
		}

	})
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

func (ss *Server) AddConnCount() {
	atomic.AddInt64(&ss.connMaxNum, 1)
}

func (ss *Server) CanConn() (int64, bool) {
	connCount := atomic.LoadInt64(&ss.connMaxNum)
	if connCount < 3000 { // todo
		return connCount, true
	}
	return connCount, false
}

func (ss *Server) DelConnCount() {
	atomic.AddInt64(&ss.connMaxNum, -1)
}

func (ss *Server) Close() error {
	if ss.server != nil {
		return ss.server.Close()
	}
	return nil
}

func (ss *Server) Run(option Option) error {
	if option.ConnMaxNum != 0 {
		ss.connMaxNum = option.ConnMaxNum
	} else {
		ss.connMaxNum = DefaultConnMaxNum
	}
	ss.server.OnConnect(ss.namespace, func(s socketio.Conn) error {
		connCount, ok := ss.CanConn()
		if !ok {
			log.Warn("server conn count is max: please wait  ", connCount)
			return nil
		}
		ss.AddConnCount()
		log.Warn("socketIO client connect ", "connMaxNum: ", connCount, s.ID(), s.LocalAddr(), s.RemoteAddr(), )
		s.Emit("message", "connected ")
		return nil
	})

	ss.server.OnError(ss.namespace, func(s socketio.Conn, e error) {
		log.Error("socketIo error ", e.Error())
		if s != nil {
			log.Error("socketIo error： info: ", s.ID(), s.LocalAddr(), s.RemoteAddr())
			s.LeaveAll()
			err := s.Close()
			if err != nil {
				log.Error(err.Error())

			}
			ss.DelConnCount()
		}

	})

	ss.server.OnDisconnect(ss.namespace, func(s socketio.Conn, reason string) {
		log.Error("socketIo client disConnect ", reason)
		if s != nil {
			log.Error("socketIO client disConnect: ", s.ID(), s.LocalAddr(), s.RemoteAddr(), )
			s.LeaveAll()
			err := s.Close()
			if err != nil {
				log.Error(err.Error())

			}
			ss.DelConnCount()
		}
	})
	return ss.server.Serve()

}

type GenId struct {
}

func (g GenId) NewID() string {

	return utils.GetUUID()
}

func NewServer() *Server {

	server, err := socketio.NewServer(
		&engineio.Options{
			ConnInitor: func(request *http.Request, conn engineio.Conn) {
				log.Info("socketIO request Conn: ", request.URL.Host, conn.ID())
			},
			Transports: []transport.Transport{
				&websocket.Transport{
					CheckOrigin: func(r *http.Request) bool {
						return true
					},
				},
			},
		})
	if err != nil {
		panic(fmt.Errorf("init socket-io server %v \n", err))
	}
	return &Server{
		server:    server,
		namespace: "/",
	}

}
