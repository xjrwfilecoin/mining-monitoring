package service

import (
	"fmt"
	"mining-monitoring/config"
	"mining-monitoring/log"
	"mining-monitoring/net/socket"
	"mining-monitoring/store"
)

type MinerInfoService struct {
	storageManager *store.Manager
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
	info := m.storageManager.GetMinerInfo()
	log.Info("minerInfo result: ", info)
	c.SuccessResp(info)

}

func NewMinerInfoService(sm *store.Manager, server *socket.Server) IMinerInfo {
	return &MinerInfoService{
		socketServer: server,
		storageManager: sm,
	}
}
