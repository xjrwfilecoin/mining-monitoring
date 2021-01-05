package shellParsing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mining-monitoring/log"
	"time"
)

type Manager struct {
	currentInfo map[string]interface{}
	shellParse  *ShellParse
	Workers     []WorkerInfo
}

func (m *Manager) GetCurrentMinerInfo() interface{} {
	return m.currentInfo
}

func (m *Manager) DoShell() (interface{}, error) {
	taskInfo, err := m.shellParse.getTaskInfo()
	if err != nil {
		return nil, err
	}
	fmt.Printf("minerInfo:  %v \n", taskInfo)
	// todo
	return taskInfo, nil

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
			result, err := m.DoShell()
			if err != nil {
				fmt.Printf("doShell error %v \n", err)
				continue
			}
			obj <- result
		default:

		}
	}
}

func NewManager(path string) (*Manager, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read workerhost error %v \n", err)
	}
	var workers []WorkerInfo
	err = json.Unmarshal(data, &workers)
	if err != nil {
		return nil, fmt.Errorf("parse json error: %v \n", err)
	}

	return &Manager{
		currentInfo: map[string]interface{}{},
		shellParse:  NewShellParse(),
		Workers:     workers,
	}, nil
}
