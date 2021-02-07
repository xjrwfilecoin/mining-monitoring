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
		m.updateWorkerState(obj.MinerId, obj.Data)
	} else if obj.CmdType == shell.LotusMinerJobs { //  jobs列表
		m.updateMinerJobs(obj.MinerId, obj.Data)
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

func (m *Manager) updateMinerJobsV1(minerId string, obj interface{}) {
	log.Error("checkJobs: rec ", obj)
	for _, workerInfo := range m.WorkerInfoTable {
		workerInfo.clearQueue()
	}
	if jobsMap, ok := obj.([]map[string]interface{}); ok {
		for i := 0; i < len(jobsMap); i++ {
			job := jobsMap[i]
			if hostName, ok := job["hostName"]; ok {
				if host, ok := hostName.(string); ok {
					workerId := WorkerId{MinerId: MinerId(minerId), HostName: host}
					workerInfo, ok := m.WorkerInfoTable[workerId]
					if ok {
						workerInfo.updateMinerJob(job)
					}
					log.Error("checkJobs: taskQueue: ", workerId, job)
				}
			}
		}
	}
}

func (m *Manager) updateMinerJobs(minerId string, obj interface{}) {
	log.Error("checkJobs: rec ", obj)
	if jobsMap, ok := obj.([]map[string]interface{}); ok {
		mapByHost := mapByHost(jobsMap)
		mapByState := mapByState(mapByHost)
		mapByType := mapByType(mapByState)
		log.Error("checkJobs: mapByType: ", mapByType)
		for hostName, taskQueue := range mapByType {
			workerId := WorkerId{MinerId: MinerId(minerId), HostName: hostName}
			workerInfo, ok := m.WorkerInfoTable[workerId]
			if ok {
				workerInfo.updateJobQueue(taskQueue)
			}
			log.Error("checkJobs: taskQueue: ", workerId, taskQueue)
		}
	}
}

func (m *Manager) updateWorkerState(minerId string, obj interface{}) {
	log.Debug("workers: ", obj)
	if workerList, ok := obj.([]map[string]interface{}); ok {
		for i := 0; i < len(workerList); i++ {
			worker := workerList[i]
			log.Warn("workers detail: ", worker)
			if hostName, ok := worker["hostName"]; ok {
				if host, ok := hostName.(string); ok {
					log.Warn("workers host", host)
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
		// todo
		// 清除不存在的worker
		for workerId, _ := range m.WorkerInfoTable {
			contain := false
			for i := 0; i < len(workerList); i++ {
				worker := workerList[i]
				if hostName, ok := worker["hostName"]; ok {
					if workerId.HostName == hostName {
						contain = true
						break
					}
				}
			}
			if !contain {
				delete(m.WorkerInfoTable, workerId)
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
	if !ok {
		return nil
	}
	result = minerInfo.GetDiff(false)
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
			log.Debug("cache rec: ", data)
			m.Update(data)
		case <-m.closing:
			return
		}
	}
}
