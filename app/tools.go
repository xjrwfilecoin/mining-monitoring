package app

import (
	"encoding/json"
	"reflect"
)

func MapParse(workerInfo, workerHardwareInfo map[string]interface{}) interface{} {
	mapByHostName := mapByHostName(workerInfo)

	mapByState := mapByState(mapByHostName)

	mapByType := mapByType(mapByState)

	var res []interface{}
	if workerHardwareInfo == nil || len(workerHardwareInfo) == 0 {
		for _, workerInfo := range mapByType {
			res = append(res, workerInfo)
		}
		return res
	}

	// 结合硬件信息
	for hostName, hardwareInfo := range workerHardwareInfo {
		thInfo := hardwareInfo.(map[string]interface{})
		if hInfo, ok := mapByType[hostName]; ok { // jobs 中存在的主机在更新设备硬件信息
			info := hInfo.(map[string]interface{})
			param := mergeMaps(parseHardwareInfo(thInfo), info)
			res = append(res, param)
		} else {
			res = append(res, mergeMaps(thInfo))
		}
	}
	return res
}



func mapByType(data map[string]interface{}) map[string]interface{} {
	// 把按状态分组，在按照任务类型分组
	mapByTask := make(map[string]interface{})
	for hostName, taskQueues := range data {
		result := make(map[string]interface{})
		result["hostName"] = hostName
		tq := taskQueues.(map[string]interface{})
		for taskType, queue := range tq {
			q1 := queue.([]map[string]interface{})
			param := taskListByType(q1)
			result[taskType] = param
		}
		mapByTask[hostName] = result
	}
	return mapByTask
}

func taskListByType(res []map[string]interface{}) map[string]interface{} {
	param := make(map[string]interface{})
	for i := 0; i < len(res); i++ {
		task := res[i]
		tType, ok := task["task"]
		if !ok {
			continue
		}
		if taskList, ok := param[tType.(string)]; ok {
			tt := taskList.([]map[string]interface{})
			taskList = append(tt, task)
			param[tType.(string)] = taskList
		} else {
			param[tType.(string)] = []map[string]interface{}{task}
		}
	}
	return param
}

func mapByState(data map[string]interface{}) map[string]interface{} {
	// 根据任务运行状态分组
	mapByState := make(map[string]interface{})
	for host, taskList := range data {
		taskMap := make(map[string]interface{})
		var currentQueue, pendQueue []map[string]interface{}
		taskList := taskList.([]map[string]interface{})
		for i := 0; i < len(taskList); i++ {
			task := taskList[i]
			state, ok := task["state"]
			if !ok {
				continue
			}

			if state == "running" {
				currentQueue = append(currentQueue, task)
			} else {
				pendQueue = append(pendQueue, task)
			}
		}
		taskMap["currentQueue"] = currentQueue
		taskMap["pendingQueue"] = pendQueue
		mapByState[host] = taskMap
	}
	return mapByState
}

// 根据 hostName 进行分组
func mapByHostName(jobs map[string]interface{}) map[string]interface{} {
	mapByHostName := make(map[string]interface{})
	for _, task := range jobs {
		sectorInfo := task.(map[string]interface{})
		if _, ok := sectorInfo["hostName"]; !ok { // 判断扇区是否存在
			continue
		}
		hostName := sectorInfo["hostName"].(string)
		if taskList, ok := mapByHostName[hostName]; ok {
			taskMap := taskList.([]map[string]interface{})
			taskList = append(taskMap, sectorInfo)
			mapByHostName[hostName] = taskList
		} else {
			mapByHostName[hostName] = []map[string]interface{}{sectorInfo}
		}
	}
	return mapByHostName
}



func parseHardwareInfo(src map[string]interface{}) map[string]interface{} {
	var gpus []interface{}
	if gpuList, ok := src["gpuInfo"]; ok {
		gpuMap := gpuList.(map[string]interface{})
		for _, gpu := range gpuMap {
			gpus = append(gpus, gpu)
		}
	}
	var netIOes []interface{}
	if netioMap, ok := src["netIO"]; ok {
		ioMap := netioMap.(map[string]interface{})
		for _, io := range ioMap {
			netIOes = append(netIOes, io)
		}
	}
	src["gpuInfo"] = gpus
	src["netIO"] = netIOes
	return src
}

func mergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			tk := k
			tV := v
			result[tk] = tV
		}
	}
	return result
}


func DeepCopyMap(input map[string]interface{}) (map[string]interface{}, error) {
	param := make(map[string]interface{})
	data, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &param)
	if err != nil {
		return nil, err
	}
	return param, nil
}


// 比较求两个map得差集,
func DiffMap(oldMap, newMap map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range newMap {
		if reflect.TypeOf(value).Kind() == reflect.Map {
			if _, ok := oldMap[key]; !ok {
				result[key] = value
			} else {
				tempOldMap := oldMap[key].(map[string]interface{})
				tempNewMap := value.(map[string]interface{})
				diffMap := DiffMap(tempOldMap, tempNewMap)
				if diffMap != nil {
					result[key] = diffMap
				}
			}
		} else {
			if tv, ok := oldMap[key]; !ok || value != tv {
				result[key] = value
			}
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}