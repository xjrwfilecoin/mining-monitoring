package socket

import (
	"github.com/googollee/go-socket.io"
	"mining-monitoring/log"
	"mining-monitoring/service"
)

const (
	MinerInfo    = "minerInfo"
	SubMinerInfo = "subMinerInfo"
)
const (
	DefaultNamespace = "/"
	DefaultRoom      = "miner-info"
)

func Router(server *Server) {
	minerInfo := service.NewMinerInfoService()
	server.RegisterRouter(DefaultNamespace, MinerInfo, minerInfo.MinerInfo)
	server.RegisterRouter(DefaultNamespace, SubMinerInfo, func(s socketio.Conn, msg string) {
		SServer.JoinRoom(DefaultNamespace, DefaultRoom, s)
		log.Info(s.ID(), "join room ", DefaultRoom)
		s.Emit(SubMinerInfo, "sub....")
	})
}
