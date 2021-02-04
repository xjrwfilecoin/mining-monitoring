package cache

import (
	"mining-monitoring/log"
	"mining-monitoring/shellParsing"
	"sync"
	"time"
)

type Manager struct {
	// miner info 表
	MinerInfoTable map[MinerId]*MinerInfo
	// worker信息表
	WorkerInfoTable map[WorkerId]*WorkerInfo

	MinerId MinerId
	closing chan struct{}
	sync.RWMutex
}

func (m *Manager) close() {
	close(m.closing)
}

func (m *Manager) Update(obj shellParsing.CmdData) {
	m.Lock()
	defer m.Unlock()
	if obj.CmdType == shellParsing.LotusMinerInfoCmd || obj.CmdType == shellParsing.LotusControlList { // info , post
		minerId := MinerId(obj.MinerId)
		minerInfo, ok := m.MinerInfoTable[minerId]
		if !ok {
			m.MinerInfoTable[minerId] = &MinerInfo{MinerId: Value{Value: minerId}}
		}
		minerInfo.Update(obj)

	} else if obj.CmdType == shellParsing.LotusMinerWorkers { // workers列表
		m.updateWorkerNetState(obj.MinerId, obj.Data)
	} else if obj.CmdType == shellParsing.LotusMinerJobs { //  jobs列表
		m.updateTaskState(obj.MinerId, obj.Data)
	} else if obj.State == shellParsing.HardwareState { // 硬件信息
		workerId := WorkerId{MinerId: MinerId(obj.MinerId), HostName: obj.HostName}
		workerInfo, ok := m.WorkerInfoTable[workerId]
		if !ok {
			m.WorkerInfoTable[workerId] = &WorkerInfo{HostName: Value{Value: obj.HostName}}
		}
		workerInfo.Update(obj)
	}
}

// todo
func (m *Manager) updateTaskState(minerId string, obj interface{}) {
	if workerList, ok := obj.([]map[string]interface{}); ok {
		for i := 0; i < len(workerList); i++ {
			worker := workerList[i]
			if hostName, ok := worker["hostName"]; ok {
				if host, ok := hostName.(string); ok {
					workerId := WorkerId{MinerId: MinerId(minerId), HostName: host}
					workerInfo, ok := m.WorkerInfoTable[workerId]
					if !ok {
						m.WorkerInfoTable[workerId] = &WorkerInfo{HostName: Value{Value: host}}
					}
					workerInfo.UpdateNetState(worker)
				}
			}
		}
	}
}

// todo
func (m *Manager) updateWorkerNetState(minerId string, obj interface{}) {
	if workerList, ok := obj.([]map[string]interface{}); ok {
		for i := 0; i < len(workerList); i++ {
			worker := workerList[i]
			if hostName, ok := worker["hostName"]; ok {
				if host, ok := hostName.(string); ok {
					workerId := WorkerId{MinerId: MinerId(minerId), HostName: host}
					workerInfo, ok := m.WorkerInfoTable[workerId]
					if !ok {
						m.WorkerInfoTable[workerId] = &WorkerInfo{HostName: Value{Value: host}}
					}
					workerInfo.UpdateNetState(worker)
				}
			}
		}
	}
}

func (m *Manager) diffData() interface{} {
	m.Lock()
	defer m.Unlock()
	result := make(map[string]interface{})
	var workerList []map[string]interface{}
	minerInfo, ok := m.MinerInfoTable[m.MinerId]
	if ok {
		result = minerInfo.GetDiff(false)
	}
	for _, workerInfo := range m.WorkerInfoTable {
		workerDiff := workerInfo.GetDiff(false)
		if len(workerDiff) > 1 {
			workerList = append(workerList, workerDiff)
		}
		workerInfo.ChangeState(false)
	}
	minerInfo.ChangeState(false)
	result["workerInfo"] = workerList
	return result
}

func (m *Manager) send() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			diffData := m.diffData()
			log.Warn("send diff: ", diffData)
		case <-m.closing:
			return
		}
	}
}

func (m *Manager) rec(obj <-chan shellParsing.CmdData) {
	for {
		select {
		case data := <-obj:
			if len(m.MinerId) == 0 {
				m.MinerId = MinerId(data.MinerId)
			}
			//m.Update(data)
			log.Info("rec info:  ", data.Data)
		}
	}
}
