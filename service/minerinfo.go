package service

import (
	"github.com/googollee/go-socket.io"
	"mining-monitoring/log"
	"mining-monitoring/shellParsing"
	"mining-monitoring/utils"
)

type MinerInfoService struct {
	shellManager *shellParsing.Manager
}

// todo 封装中间件
func (m *MinerInfoService) MinerInfo(s socketio.Conn, data string) {
	log.Debug(s.ID(), s.LocalAddr(), "get minerInfo: ", data)
	minerInfoForm := &BaseFrom{}
	err := BindJson(data, minerInfoForm)
	if err != nil {
		utils.FailResp(minerInfoForm.Url, minerInfoForm.MsgId, s, err.Error())
		return
	}
	info := m.shellManager.GetCurrentMinerInfo()
	log.Info("minerInfo result: ", info)
	utils.SuccResp(minerInfoForm.Url, minerInfoForm.MsgId, s, info)

}

func NewMinerInfoService(sm *shellParsing.Manager) IMinerInfo {
	return &MinerInfoService{
		shellManager: sm,
	}
}
