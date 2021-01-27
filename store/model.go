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
	DataMap    map[Id]map[string]interface{}
	WorkerList shellParsing.CmdData
	sync.RWMutex
}

func NewMinerInfo() *MinerInfo {
	return &MinerInfo{
		DataMap: make(map[Id]map[string]interface{}),
	}
}

func (m *MinerInfo) Parse(obj shellParsing.CmdData) interface{} {

	switch obj.CmdType {
	case shellParsing.LotusMinerJobs:
		return LotusJobsToArrayV1(m.WorkerList, obj)
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
	if obj.CmdType == shellParsing.LotusMinerWorkers {
		m.WorkerList = obj
		return nil
	}

	m.Lock()
	defer m.Unlock()
	oldMap, ok := m.DataMap[id]
	newMap := utils.StructToMapByJson(obj.Data)
	diffMap := newMap
	m.DataMap[id] = newMap
	if ok {
		if obj.CmdType == shellParsing.LotusMinerJobs || obj.CmdType == shellParsing.SarCmd || obj.CmdType == shellParsing.GpuCmd {
			if !utils.MapIsDiff(oldMap, newMap) { // 数据格式前端为array，没有变化，不用推数据
				return nil
			}
		} else {
			diffMap = utils.DeepDiffMap(oldMap, newMap)
		}

	}
	obj.Data = diffMap
	if len(diffMap) == 0 {
		return nil
	}
	diffResult := m.Parse(obj)
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
			} else if keyId.CmdType == shellParsing.LotusMinerWorkers {

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

	workerList := utils.StructToMapByJson(m.WorkerList.Data)
	result := MergeJobsAndHardwareV1(workerList, jobsInfo, hostHardwareMap)

	response := utils.MergeMaps(minerInfo, result)
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

func MergeJobsAndHardwareV1(workers, jobs map[string]interface{}, hardwareInfoMap map[string]map[string]interface{}) map[string]interface{} {
	param := make(map[string]interface{})
	var workerList []map[string]interface{}
	for hostName, _ := range workers {
		job, jOk := jobs[hostName]
		hardware, hOk := hardwareInfoMap[hostName]
		if jOk {
			tJob, ok1 := job.(map[string]interface{})
			if !ok1 {
				continue
			}

			if hOk {
				workerList = append(workerList, utils.MergeMaps(tJob, hardware))
			} else {
				workerList = append(workerList, fixHardWare(tJob))
			}
		} else {
			newEmptyInfo := NewEmptyInfo(hostName)
			if hOk {
				workerList = append(workerList, utils.MergeMaps(newEmptyInfo, hardware))
			} else {
				workerList = append(workerList, newEmptyInfo)
			}

		}
	}

	//for hostName, deviceInfo := range hardwareInfoMap {
	//	job, ok := jobs[hostName]
	//	if ok {
	//		if tJob, ok1 := job.(map[string]interface{}); ok1 {
	//			workerList = append(workerList, utils.MergeMaps(tJob, deviceInfo))
	//		}
	//	} else {
	//		worker := NewMap()
	//		worker["hostName"] = hostName
	//		workerList = append(workerList, utils.MergeMaps(worker, deviceInfo))
	//	}
	//}
	param["workerInfo"] = workerList
	return param
}

func fixHardWare(job map[string]interface{}) map[string]interface{} {
	job["cpuLoad"] = "0"
	job["cpuTemper"] = "0"
	job["diskR"] = "0"
	job["diskW"] = "0"
	job["gpuInfo"] = []interface{}{}
	job["netIO"] = []interface{}{}
	job["totalMemory"] = "0"
	job["useDisk"] = "0"
	job["useMemory"] = "0"
	return job
}

func NewEmptyInfo(hostName string) map[string]interface{} {
	newMap := NewMap()
	newMap["hostName"] = hostName
	newMap["currentQueue"] = []interface{}{}
	newMap["pendingQueue"] = []interface{}{}
	newMap["cpuLoad"] = "0"
	newMap["cpuTemper"] = "0"
	newMap["diskR"] = "0"
	newMap["diskW"] = "0"
	newMap["gpuInfo"] = []interface{}{}
	newMap["netIO"] = []interface{}{}
	newMap["totalMemory"] = "0"
	newMap["useDisk"] = "0"
	newMap["useMemory"] = "0"
	return newMap
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
