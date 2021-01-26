package store

import (
	"mining-monitoring/log"
	"mining-monitoring/shellParsing"
	"mining-monitoring/utils"
	"sync"
)

type MinerId string

type Id struct {
	MinerId  string
	HostName string
	CmdType  shellParsing.CmdType
	CmdState shellParsing.CmdState
}

type MinerInfo struct {
	DataMap map[Id]map[string]interface{}
	sync.RWMutex
}

func NewMinerInfo() *MinerInfo {
	return &MinerInfo{
		DataMap: make(map[Id]map[string]interface{}),
	}
}

func Parse(obj shellParsing.CmdData) interface{} {

	switch obj.CmdType {
	case shellParsing.LotusMinerJobs:
		return LotusJobsToArray(obj)
	case shellParsing.GpuCmd:
		return NewCommArrayMap(obj.HostName, "gpuInfo", MapToArray(obj))
	case shellParsing.SarCmd:
		return NewCommArrayMap(obj.HostName, "netIO", MapToArray(obj))
	default:
	}

	switch obj.State {
	case shellParsing.LotusState:
		return NewCommonMap(obj)
	case shellParsing.HardwareState:
		return NewWorkerInfoMap(obj)
	}
	return nil
}

func (m *MinerInfo) updateData(obj shellParsing.CmdData) interface{} {
	id := Id{MinerId: obj.MinerId, HostName: obj.HostName, CmdType: obj.CmdType, CmdState: obj.State}
	m.Lock()
	defer m.Unlock()
	oldMap, ok := m.DataMap[id]
	newMap := utils.StructToMapByJson(obj.Data)
	diffMap := newMap
	m.DataMap[id] = newMap
	if obj.CmdType == shellParsing.LotusMinerWorkers { // workerList信息不用处理后面陆续上报
		return nil
	}
	if ok {
		if obj.CmdType != shellParsing.LotusMinerJobs { // jobs命令时间每时每刻都在变化，不用处理
			diffMap = utils.DeepDiffMap(oldMap, newMap)
		}

	}

	obj.Data = diffMap
	if len(diffMap) == 0 {
		return nil
	}
	diffResult := Parse(obj)
	log.Debug("check diff result: ", "type: ", obj.CmdType, "src: ", newMap, "diff: ", diffMap, "diffResult: ", diffResult)
	return diffResult

}

// 把map中数据找出来，封装成指定数据格式返回前端
func (m *MinerInfo) getMinerInfo(minerId string) interface{} {
	m.Lock()
	defer m.Unlock()
	hostHardwareMap := make(map[string]map[string]interface{}) // 根据hostName进行分组
	jobsInfo := make(map[string]interface{})
	minerInfo := make(map[string]interface{})
	for keyId, value := range m.DataMap {
		//log.Error(keyId, value)
		if keyId.CmdState == shellParsing.LotusState {
			if keyId.CmdType == shellParsing.LotusMinerJobs {
				jobsInfo = JobsToArrayV1(value)
			} else {
				minerInfo = utils.MergeMaps(minerInfo, value)
			}
		} else {
			info, ok := hostHardwareMap[keyId.HostName]
			temp := make(map[string]interface{})
			temp["hostName"] = keyId.HostName
			if keyId.CmdType == shellParsing.SarCmd {
				netIO := TraverseMap(value)
				temp["netIO"] = netIO
				value = temp
			} else if keyId.CmdType == shellParsing.GpuCmd {
				gpuInfo := TraverseMap(value)
				temp["gpuInfo"] = gpuInfo
				value = temp
			}

			if ok {
				hostHardwareMap[keyId.HostName] = utils.MergeMaps(info, value)
			} else {
				hostHardwareMap[keyId.HostName] = value
			}
		}
	}
	result := MergeJobsAndHardware(jobsInfo, hostHardwareMap)
	response := utils.MergeMaps(minerInfo, result)
	//bytes, _ := json.Marshal(response)
	//log.Error(string(bytes))
	return response
}

func TraverseMap(param map[string]interface{}) interface{} {
	var res []interface{}
	for _, value := range param {
		res = append(res, value)
	}
	return res
}

func JobsToArrayV1(param map[string]interface{}) map[string]interface{} {
	mapByHost := mapByHost(param)
	mapByState := mapByState(mapByHost)
	mapByType := mapByType(mapByState)

	return mapByType
}

func MergeJobsAndHardware(jobs map[string]interface{}, hadrMap map[string]map[string]interface{}) map[string]interface{} {
	param := make(map[string]interface{})
	var workerList []map[string]interface{}
	for hostName, job := range jobs {
		hardinfo, ok := hadrMap[hostName]
		if ok {
			tJob, ok1 := job.(map[string]interface{})
			if ok1 {
				workerList = append(workerList, utils.MergeMaps(tJob, hardinfo))
			}
		}
	}
	param["workerInfo"] = workerList
	return param
}
