package service

import (
	"fmt"
	"mining-monitoring/config"
	"mining-monitoring/log"
	"mining-monitoring/net/socket"
	"mining-monitoring/shellParsing"
)

type MinerInfoService struct {
	shellManager *shellParsing.Manager
	socketServer *socket.Server
}

func (m *MinerInfoService) SuMinerInfo(c *socket.Context) {
	m.socketServer.JoinRoom(config.DefaultNamespace, config.DefaultRoom, c.Conn)
	log.Debug("join room ")
	c.SuccessResp(nil)
}

func (m *MinerInfoService) MinerInfo(c *socket.Context) {
	log.Debug("get minerInfo: ", c.Body)
	minerFrom := &MinerInfoForm{}
	err := c.BindJson(minerFrom)
	if err != nil {
		c.FailResp(fmt.Errorf("param is error: %v \n", err.Error()).Error())
		return
	}
	info := m.shellManager.GetCurrentMinerInfo()
	log.Info("minerInfo result: ", info)
	c.SuccessResp(info)

}

func NewMinerInfoService(sm *shellParsing.Manager, server *socket.Server) IMinerInfo {
	return &MinerInfoService{
		socketServer: server,
		shellManager: sm,
	}
}
