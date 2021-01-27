package store

import (
	"mining-monitoring/config"
	"mining-monitoring/log"
	"mining-monitoring/net/socket"
	"mining-monitoring/shellParsing"
	"sync"
	"time"
)

type Manager struct {
	Miners   map[MinerId]*MinerInfo
	sendSign chan interface{}
	sync.Mutex
	MinerId string
}

func NewManager() *Manager {
	return &Manager{
		Miners:   make(map[MinerId]*MinerInfo),
		sendSign: make(chan interface{}, 100),
	}
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

func (m *Manager) test() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			m.GetMinerInfo()
		}
	}
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
		default:

		}
	}
}

func (m *Manager) Send() {
	for {
		select {
		case diffData := <-m.sendSign:
			if diffData != nil {
				log.Debug("send diff map:  ", diffData)
				socket.BroadCaseMsg(config.DefaultNamespace, config.DefaultRoom, config.SubMinerInfo, diffData)
			}

		default:

		}
	}
}
