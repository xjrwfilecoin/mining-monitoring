package store

import (
	"encoding/json"
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

func (m *MinerInfo) updateData(obj shellParsing.CmdData) interface{} {
	m.Lock()
	defer m.Unlock()
	id := Id{MinerId: obj.MinerId, HostName: obj.HostName, CmdType: obj.CmdType, CmdState: obj.State}
	m.DataMap[id] = utils.StructToMapByJson(obj.Data)
	return nil
}

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
				temp["netIO"]=netIO
				value=temp
			} else if keyId.CmdType == shellParsing.GpuCmd {
				gpuInfo := TraverseMap(value)
				temp["gpuInfo"]=gpuInfo
				value=temp
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
	bytes, _ := json.Marshal(response)
	log.Error(string(bytes))
	return nil
}

func DiffMapValue(new, old shellParsing.CmdData) map[string]interface{} {
	newMap := utils.StructToMapByJson(new.Data)
	oldMap := utils.StructToMapByJson(old.Data)
	return utils.DiffMap(oldMap, newMap)
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
