package service

import (
	"github.com/googollee/go-socket.io"
	"mining-monitoring/log"
	"mining-monitoring/shellParsing"
)

type MinerInfoService struct {
}

func (m *MinerInfoService) MinerInfo(s socketio.Conn, msg string) {
	log.Debug(s.ID(), s.LocalAddr(), "get minerInfo: ", msg)
	currentMinerInfo := shellParsing.MinerInfoManager.GetCurrentMinerInfo()
	// todo
	s.Emit("minerInfo", currentMinerInfo)

}


func NewMinerInfoService() IMinerInfo {
	return &MinerInfoService{}
}
