package store

import (
	"encoding/json"
	"mining-monitoring/config"
	"mining-monitoring/log"
	"mining-monitoring/net/socket"
	"mining-monitoring/shellParsing"
	"mining-monitoring/utils"
	"sync"
	"time"
)

type Manager struct {
	Miners   map[MinerId]*MinerInfo
	sendSign chan interface{}
	sync.Mutex
	socket  socket.Server
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
	go m.Send()
	go m.test()
	for {
		select {
		case data := <-obj:
			m.MinerId = data.MinerId
			log.Debug("store rec data: ", data)
			recvInfo := m.Parse(data)
			bytes, err := json.Marshal(recvInfo)
			if err != nil {
				log.Error(err.Error())
			}
			//m.sendSign <- data
			log.Info("send info: ", string(bytes))
			socket.BroadCaseMsg(config.DefaultNamespace, config.DefaultRoom, config.SubMinerInfo, recvInfo)
			minerId := MinerId(data.MinerId)
			m.Lock()
			minerInfo, ok := m.Miners[minerId]
			m.Unlock()
			if !ok {
				minerInfo = NewMinerInfo()
				m.Miners[minerId] = minerInfo
			}
			go m.UpdateMinerInfo(minerInfo, data)
		}
	}
}

func (m *Manager) Parse(obj shellParsing.CmdData) interface{} {

	switch obj.CmdType {
	case shellParsing.LotusMinerJobs:
		return JobsToArray(obj)
	case shellParsing.GpuCmd:
		return NewCommMap(obj.HostName, "gpuInfo", MapToArray(obj))
	case shellParsing.SarCmd:
		return NewCommMap(obj.HostName, "netIO", MapToArray(obj))
	default:
	}
	switch obj.State {
	case shellParsing.LotusState:
		return MinerResponse(obj)
	case shellParsing.HardwareState:
		return WorkerResponse(obj)
	}
	return nil
}

func JobsToArray(obj shellParsing.CmdData) interface{} {
	var workerList []map[string]interface{}
	response := NewMap()
	param := utils.StructToMapByJson(obj.Data)
	mapByHost := mapByHost(param)
	mapByState := mapByState(mapByHost)
	mapByType := mapByType(mapByState)
	for _, value := range mapByType {
		if tv, ok := value.(map[string]interface{}); ok {
			workerList = append(workerList, tv)
		}
	}
	response["workerInfo"] = workerList
	return response
}

func MapToArray(obj shellParsing.CmdData) interface{} {
	param := utils.StructToMapByJson(obj.Data)
	var res []map[string]interface{}
	for _, value := range param {
		if item, ok := value.(map[string]interface{}); ok {
			res = append(res, item)
		}
	}
	return res
}

func NewCommMap(hostName, key string, obj interface{}) interface{} {
	var workerList []map[string]interface{}
	response := NewMap()
	worker := NewMap()
	worker["hostName"] = hostName
	worker[key] = obj
	workerList = append(workerList, worker)
	response["workerInfo"] = workerList
	return response

}

func MinerResponse(obj shellParsing.CmdData) interface{} {
	mapByJson := utils.StructToMapByJson(obj.Data)
	return mapByJson
}

func WorkerResponse(obj shellParsing.CmdData) interface{} {
	var workerList []map[string]interface{}
	response := NewMap()
	worker := NewMap()
	worker["hostName"] = obj.HostName
	mapByJson := utils.StructToMapByJson(obj.Data)
	mergeMaps := utils.MergeMaps(worker, mapByJson)
	workerList = append(workerList, mergeMaps)
	response["workerInfo"] = workerList
	return response
}

func NewMap() map[string]interface{} {
	param := make(map[string]interface{})
	return param
}

func (m *Manager) Send() {
	//for {
	//	select {
	//	case diffData := <-m.sendSign:
	//		//log.Debug("send diff map:  ", diffData)
	//		// todo 广播
	//		socket.BroadCaseMsg(config.DefaultNamespace, config.DefaultRoom, config.SubMinerInfo, diffData)
	//	}
	//}
}

func (m *Manager) UpdateMinerInfo(minerInfo *MinerInfo, obj shellParsing.CmdData) {
	diffMap := minerInfo.updateData(obj)
	m.sendSign <- diffMap
}
