package service

import (
	"fmt"
	"mining-monitoring/log"
	"mining-monitoring/net/socket"
	"mining-monitoring/shellParsing"
)

type MinerInfoService struct {
	shellManager *shellParsing.Manager
}

func (m *MinerInfoService) MinerInfo(c *socket.Context) {
	log.Debug("get minerInfo: ", c.Body)
	minerFrom := &MinerInfoForm{}
	err := c.BindJson(minerFrom)
	if err != nil {
		c.FailResp(fmt.Errorf("param is error: %v \n",err.Error()).Error())
		return
	}
	info := m.shellManager.GetCurrentMinerInfo()
	log.Info("minerInfo result: ", info)
	c.SuccessResp(info)

}

func NewMinerInfoService(sm *shellParsing.Manager) IMinerInfo {
	return &MinerInfoService{
		shellManager: sm,
	}
}
