package store

import (
	"mining-monitoring/config"
	"mining-monitoring/log"
	"mining-monitoring/net/socket"
	"mining-monitoring/shellParsing"
	"sync"
)

type Manager struct {
	Miners   map[MinerId]*MinerInfo
	sendSign chan interface{}
	sync.Mutex
	MinerId string
	closing chan struct{}
}

func NewManager() *Manager {
	return &Manager{
		Miners:   make(map[MinerId]*MinerInfo),
		sendSign: make(chan interface{}, 100),
		closing:  make(chan struct{}, 1),
	}
}

func (m *Manager) Close() {
	close(m.closing)
	close(m.sendSign)
}

func (m *Manager) GetMinerInfo() interface{} {
	m.Lock()
	defer m.Unlock()
	minerInfo, ok := m.Miners[MinerId(m.MinerId)]
	if !ok {
		return nil
	}
	info := minerInfo.getMinerInfo(m.MinerId)
	return info
}

func (m *Manager) Recv(obj chan shellParsing.CmdData) {
	for {
		select {
		case data := <-obj:
			m.MinerId = data.MinerId
			log.Debug("store rec data: ", data)
			minerId := MinerId(data.MinerId)
			minerInfo, ok := m.Miners[minerId]
			if !ok {
				minerInfo = NewMinerInfo()
				m.Miners[minerId] = minerInfo
			}
			diffData := minerInfo.updateData(data)
			m.sendSign <- diffData
		case <-m.closing:
			return
		}
	}
}

func (m *Manager) Send() {
	for {
		select {
		case diffData := <-m.sendSign:
			if diffData != nil {
				log.Debug("send diff map:  ", diffData)
				go socket.BroadCaseMsg(config.DefaultNamespace, config.DefaultRoom, config.SubMinerInfo, diffData)
			}

		case <-m.closing:
			return

		}
	}
}
