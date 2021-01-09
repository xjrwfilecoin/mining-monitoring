package service

import (
	"github.com/googollee/go-socket.io"
	"mining-monitoring/log"
	"mining-monitoring/shellParsing"
)

type MinerInfoService struct {
	shellManager *shellParsing.Manager
}

func (m *MinerInfoService) MinerInfo(s socketio.Conn, msg string) {
	log.Debug(s.ID(), s.LocalAddr(), "get minerInfo: ", msg)
	info := m.shellManager.GetCurrentMinerInfo()
	log.Info("minerInfo result: ", info)
	s.Emit("minerInfo", "nihao")

}

func NewMinerInfoService(sm *shellParsing.Manager) IMinerInfo {
	return &MinerInfoService{
		shellManager:sm,
	}
}
