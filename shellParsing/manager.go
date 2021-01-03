package shellParsing

import (
	"mining-monitoring/log"
	"time"
)

var MinerInfoManager = NewManager()

type Manager struct {
	currentInfo map[string]interface{}
	shellParse  *ShellParse
}

func (m *Manager) GetCurrentMinerInfo() interface{} {
	return m.currentInfo
}

func (m *Manager) doShell() interface{} {
	currentInfo := m.shellParse.getCurrentInfo()
	return m.getEffectiveInfo(currentInfo)

}

func (m *Manager) getEffectiveInfo(param map[string]interface{}) map[string]interface{} {
	firstMap := make(map[string]interface{})
	secondMap := make(map[string]interface{})
	if len(param) > len(m.currentInfo) {
		firstMap = param
		secondMap = m.currentInfo
	} else {
		firstMap = m.currentInfo
		secondMap = param
	}
	tmpMap := make(map[string]interface{})
	for key, value := range firstMap {
		tk := key
		tv := value
		// todo 深度
		if v, ok := secondMap[tk]; !ok || v != tv {
			tmpMap[tk] = tv
		}
	}
	m.currentInfo = param
	return tmpMap
}

func (m *Manager) Run(obj chan interface{}) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Debug("start timer get minerInfo ")
			result := m.doShell()
			obj <- result
		default:

		}
	}
}

func NewManager() *Manager {
	return &Manager{
		currentInfo: map[string]interface{}{},
		shellParse:  &ShellParse{},
	}
}
