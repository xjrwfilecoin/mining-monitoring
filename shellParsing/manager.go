package shellParsing

import (
	"fmt"
	"mining-monitoring/log"
)

type Manager struct {
	currentInfo map[string]interface{}
	shellParse  *ShellParse
	Workers     []WorkerInfo
}

func (m *Manager) GetCurrentMinerInfo() interface{} {
	return m.currentInfo
}

func (m *Manager) UpdateCurrentMinerInfo(info map[string]interface{}) {
	m.currentInfo = info
}

func (m *Manager) DoShell() (map[string]interface{}, error) {
	if e := recover(); e != nil {
		log.Error("doShell error: ", e)
	}
	taskInfo, err := m.shellParse.getTaskInfo()
	if err != nil {
		return nil, err
	}
	return taskInfo, nil

}

func (m *Manager) Run(obj chan map[string]interface{}) {
	if e := recover(); e != nil {
		log.Error("manager shell error ", e)
	}
	result, err := m.DoShell()
	if err != nil {
		log.Error("manager do shell error: %v \n", err)
	} else {
		m.currentInfo = result
		obj <- result
	}
	for {
		log.Debug("start timer get minerInfo ")
		result, err = m.DoShell()
		if err != nil {
			fmt.Printf("doShell error %v \n", err)
		}
		obj <- result
	}

	//ticker := time.NewTicker(30 * time.Second)
	//defer ticker.Stop()
	//for {
	//	select {
	//	case <-ticker.C:
	//		log.Debug("start timer get minerInfo ")
	//		result, err := m.DoShell()
	//		if err != nil {
	//			fmt.Printf("doShell error %v \n", err)
	//			continue
	//		}
	//		obj <- result
	//	default:
	//
	//	}
	//}
}

func NewManager() (*Manager, error) {

	return &Manager{
		currentInfo: map[string]interface{}{},
		shellParse:  NewShellParse(),
	}, nil
}
