package store

import (
	"mining-monitoring/log"
	"mining-monitoring/shellParsing"
	"sync"
)

type Manager struct {
	Miners   map[MinerId]*MinerInfo
	sendSign chan map[string]interface{}
	ml       sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		Miners:   make(map[MinerId]*MinerInfo),
		sendSign: make(chan map[string]interface{}, 100),
	}
}

func (m *Manager) Recv(obj chan shellParsing.CmdData) {
	go m.Send()
	for {
		select {
		case data := <-obj:
			log.Debug("store rec data: %v ", data)
			minerId := MinerId(data.MinerId)
			m.ml.Lock()
			minerInfo, ok := m.Miners[minerId]
			m.ml.Unlock()
			if !ok {
				minerInfo = &MinerInfo{MinerId: minerId}
				m.Miners[minerId] = minerInfo
			}
			go m.UpdateMinerInfo(minerInfo, data)
		}
	}
}

func (m *Manager) Send() {
	for {
		select {
		case diffData := <-m.sendSign:
			log.Debug("send diff map:  %v ", diffData)
			// todo 广播
		}
	}
}

func (m *Manager) UpdateMinerInfo(minerInfo *MinerInfo, obj shellParsing.CmdData) {
	diffMap := minerInfo.updateData(obj)
	m.sendSign <- diffMap
}
