package cache

import (
	"mining-monitoring/config"
	"mining-monitoring/log"
	"mining-monitoring/net/socket"
	"mining-monitoring/shell"
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

func NewManager() *Manager {
	return &Manager{
		MinerInfoTable:  make(map[MinerId]*MinerInfo),
		WorkerInfoTable: make(map[WorkerId]*WorkerInfo),
		closing:         make(chan struct{}, 1),
	}
}

func (m *Manager) Close() {
	close(m.closing)
}

func (m *Manager) Update(obj shell.CmdData) {
	m.Lock()
	defer m.Unlock()
	if obj.CmdType == shell.LotusMinerInfoCmd || obj.CmdType == shell.LotusControlList { // info , post
		minerId := MinerId(obj.MinerId)
		minerInfo, ok := m.MinerInfoTable[minerId]
		if !ok {
			minerInfo = NewMinerInfo(obj.MinerId)
			m.MinerInfoTable[minerId] = minerInfo
		}
		minerInfo.Update(obj)

	} else if obj.CmdType == shell.LotusMinerWorkers { // workers列表
		m.updateTaskState(obj.MinerId, obj.Data)
	} else if obj.CmdType == shell.LotusMinerJobs { //  jobs列表
		m.updateWorkerTask(obj.MinerId, obj.Data)
	} else if obj.State == shell.HardwareState { // 硬件信息
		workerId := WorkerId{MinerId: MinerId(obj.MinerId), HostName: obj.HostName}
		workerInfo, ok := m.WorkerInfoTable[workerId]
		if !ok {
			workerInfo = NewWorkerInfo(obj.HostName)
			m.WorkerInfoTable[workerId] = workerInfo
		}
		workerInfo.Update(obj)
	}
}

func (m *Manager) updateWorkerTask(minerId string, obj interface{}) {
	if jobsMap, ok := obj.([]map[string]interface{}); ok {
		mapByHost := mapByHost(jobsMap)
		mapByState := mapByState(mapByHost)
		mapByType := mapByType(mapByState)
		for hostName, taskQueue := range mapByType {
			workerId := WorkerId{MinerId: MinerId(minerId), HostName: hostName}
			workerInfo, ok := m.WorkerInfoTable[workerId]
			if !ok {
				workerInfo = NewWorkerInfo(hostName)
				m.WorkerInfoTable[workerId] = workerInfo
			}
			workerInfo.updateJobQueue(taskQueue)
		}
	}
}

// todo
func (m *Manager) updateTaskState(minerId string, obj interface{}) {
	log.Debug("workers: ", obj)
	if workerList, ok := obj.([]map[string]interface{}); ok {
		for i := 0; i < len(workerList); i++ {
			worker := workerList[i]
			log.Debug("workers detail: ", worker)
			if hostName, ok := worker["hostName"]; ok {
				if host, ok := hostName.(string); ok {
					log.Debug("workers host", host)
					workerId := WorkerId{MinerId: MinerId(minerId), HostName: host}
					workerInfo, ok := m.WorkerInfoTable[workerId]
					if !ok {
						workerInfo = NewWorkerInfo(host)
						m.WorkerInfoTable[workerId] = workerInfo
					}
					workerInfo.UpdateTaskType(worker)
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
	if len(workerList) != 0 {
		result["workerInfo"] = workerList
	}
	if len(result) == 1 {
		return nil
	}
	return result
}

func (m *Manager) GetMinerInfo() interface{} {
	m.Lock()
	result := make(map[string]interface{})
	defer m.Unlock()
	for _, minerInfo := range m.MinerInfoTable {
		result = minerInfo.GetDiff(true)
	}
	var workerList []interface{}
	for _, workerInfo := range m.WorkerInfoTable {
		worker := workerInfo.GetDiff(true)
		workerList = append(workerList, worker)
	}
	result["workerInfo"] = workerList
	return result

}

func (m *Manager) Send() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			diffData := m.diffData()
			if diffData == nil {
				continue
			}
			if socket.SServer.GetServer().Count() > 0 {
				socket.BroadCaseMsg(config.DefaultNamespace, config.DefaultRoom, config.SubMinerInfo, diffData)
			}
		case <-m.closing:
			return
		}
	}
}

func (m *Manager) Rec(obj <-chan shell.CmdData) {
	for {
		select {
		case data := <-obj:
			if len(m.MinerId) == 0 {
				m.MinerId = MinerId(data.MinerId)
			}
			m.Update(data)
		case <-m.closing:
			return
		}
	}
}
