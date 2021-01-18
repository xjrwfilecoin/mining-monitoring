package shellParsing

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)


// jobs "jobs":{"18":{"hostName":"worker01","id":"d7fd42c9","sector":"17","state":"running","task":"PC1","time":"17m48s","worker":"98c441ab"},}
// hardwareInfo  hardwareInfo":{"worker01":{"cpuLoad":"14.73","cpuTemper":"+41.1°C","diskR":"906.67M/s","diskW":"163.63M/s","gpuInfo":{"0":{"name":"0","temp":"91C","use":"100%"}},"hostName":"worker01","netIO":{"eno1":{"name":"eno1","rx":"1.27","tx":"2.90"},"eno2":{"name":"eno2","rx":"0.00","tx":"0.00"},"enp2s0f0np0":{"name":"enp2s0f0np0","rx":"0.00","tx":"0.00"},"enp2s0f1np1":{"name":"enp2s0f1np1","rx":"0.00","tx":"0.00"},"lo":{"name":"lo","rx":"0.00","tx":"0.00"}},"totalMemory":"503G","useDisk":"40%","useMemory":"319G"}},
// 根据 hostName 分组归纳信息
// jobs task信息 ; workerHardwareInfo 硬件列表信息
func ParseJobsInfo(jobs, workerHardwareInfo map[string]interface{}) interface{} {

	// job 根据 hostName 分组 // {"hostName":{"sector":"111","state":""}}
	mapByHostName := make(map[string][]Task)
	for _, job := range jobs {
		sectorInfo := job.(map[string]interface{})

		if _, ok := sectorInfo["hostName"]; !ok {
			continue
		}
		hostName := sectorInfo["hostName"].(string)
		if taskList, ok := mapByHostName[hostName]; ok {
			taskList = append(taskList, mapToTask(sectorInfo))
			mapByHostName[hostName] = taskList
		} else {
			mapByHostName[hostName] = []Task{mapToTask(sectorInfo)}
		}
	}

	// 根据任务运行状态分组
	mapByState := make(map[string]interface{})
	for host, taskList := range mapByHostName {
		tHost := host
		taskMap := make(map[string]interface{})
		var currentQueue, pendQueue []Task

		for i := 0; i < len(taskList); i++ {
			task := taskList[i]
			if task.State == "running" {
				currentQueue = append(currentQueue, task)
			} else {
				pendQueue = append(pendQueue, task)
			}
		}
		taskMap["currentQueue"] = currentQueue
		taskMap["pendingQueue"] = pendQueue
		mapByState[tHost] = taskMap
	}

	// 把按状态分组，在按照任务类型分组
	mapByTask := make(map[string]interface{})
	for hostName, taskQueue := range mapByState {
		tHost := hostName
		result := make(map[string]interface{})
		result["hostName"] = tHost
		tq := taskQueue.(map[string]interface{})
		for taskType, queue := range tq {
			q1 := queue.([]Task)
			param := tasksByType(q1)
			result[taskType] = param
		}
		mapByTask[tHost] = result
	}

	var res []interface{}
	if workerHardwareInfo == nil || len(workerHardwareInfo) == 0 {
		for _, tasks := range mapByTask {
			res = append(res, tasks)
		}
		return res
	}
	// 结合硬件信息
	for hostName, hardwareInfo := range workerHardwareInfo {
		thInfo := hardwareInfo.(map[string]interface{})
		if hInfo, ok := mapByTask[hostName]; ok {
			info := hInfo.(map[string]interface{})
			param := mergeMaps(parseHardwareInfo(thInfo), info)
			res = append(res, param)
		}
	}
	return res

}

func parseHardwareInfo(src map[string]interface{}) map[string]interface{} {
	var gpus []interface{}
	if gpuList, ok := src["gpuInfo"]; ok {
		if gpuList != nil {
			gpuMap := gpuList.(map[string]interface{})
			for _, gpu := range gpuMap {
				gpus = append(gpus, gpu)
			}
		}

	}
	var netIOes []interface{}
	if netioMap, ok := src["netIO"]; ok {
		if netioMap != nil {
			ioMap := netioMap.(map[string]interface{})
			for _, io := range ioMap {
				netIOes = append(netIOes, io)
			}
		}
	}
	src["gpuInfo"] = gpus
	src["netIO"] = netIOes
	return src
}

// 整理worker任务信息
func mergeWorkerInfo(tasks []Task, hardwareList []HardwareInfo) interface{} {
	// 根据 hostName分组
	mapByHostName := make(map[string][]Task)
	for i := 0; i < len(tasks); i++ {
		task := tasks[i]
		if taskList, ok := mapByHostName[task.HostName]; ok {
			taskList = append(taskList, task)
			mapByHostName[task.HostName] = taskList
		} else {
			mapByHostName[task.HostName] = []Task{task}
		}
	}
	// 根据任务运行状态分组
	mapByState := make(map[string]interface{})
	for host, taskList := range mapByHostName {
		tHost := host
		taskMap := make(map[string]interface{})
		var currentQueue, pendQueue []Task

		for i := 0; i < len(taskList); i++ {
			task := taskList[i]
			if task.State == "running" {
				currentQueue = append(currentQueue, task)
			} else {
				pendQueue = append(pendQueue, task)
			}
		}
		taskMap["currentQueue"] = currentQueue
		taskMap["pendingQueue"] = pendQueue
		mapByState[tHost] = taskMap
	}

	// 把按状态分组，在按照任务类型分组
	mapByTask := make(map[string]interface{})
	for hostName, taskQueue := range mapByState {
		tHost := hostName
		result := make(map[string]interface{})
		tq := taskQueue.(map[string]interface{})
		for taskType, queue := range tq {
			q1 := queue.([]Task)
			param := tasksByType(q1)
			result[taskType] = param
		}
		mapByTask[tHost] = result

	}

	// 结合硬件信息
	var res []interface{}
	for i := 0; i < len(hardwareList); i++ {
		hardware := hardwareList[i]
		if info, ok := mapByTask[hardware.HostName]; ok {
			tp := info.(map[string]interface{})
			toMap := structToMapByJson(&hardware)
			mapByTask[hardware.HostName] = mergeMaps(tp, toMap)
			res = append(res, mergeMaps(tp, toMap))
		}
	}
	return res
}

// 根据任务类型分组
func tasksByType(res []Task) map[string]interface{} {
	param := make(map[string]interface{})
	for i := 0; i < len(res); i++ {
		task := res[i]
		if taskList, ok := param[task.Task]; ok {
			tt := taskList.([]Task)
			taskList = append(tt, task)
			param[task.Task] = taskList
		} else {
			param[task.Task] = []Task{task}
		}
	}
	return param
}

func structToMapByJson(obj interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	bytes, err := json.Marshal(obj)
	if err != nil {
		return m
	}
	err = json.Unmarshal(bytes, &m)
	if err != nil {
		return m
	}
	return m
}

func structToMapByReflect(obj interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return m
	}
	elem := reflect.ValueOf(obj).Elem()
	relType := elem.Type()
	for i := 0; i < relType.NumField(); i++ {
		m[getHumpName(relType.Field(i).Name)] = elem.Field(i).Interface()
	}
	return m
}

func getHumpName(name string) string {
	if len(name) < 1 {
		return name
	}
	return fmt.Sprintf("%v%v", strings.ToLower(name[0:1]), name[1:])
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

func DeleteMapNull(src *map[string]interface{}) *map[string]interface{} {
	for key, value := range *src {
		if value == nil {
			delete(*src, key)
		} else {
			if reflect.TypeOf(value).Kind() == reflect.Map {
				tValue := value.(map[string]interface{})
				if len(tValue) == 0 {
					delete(*src, key)
				} else {
					DeleteMapNull(&tValue)
				}
			}
		}
	}
	return src
}

func getRegexValue(src [][]string) string {
	if len(src) == 0 || len(src[0]) == 0 {
		return "0"
	}
	return strings.ReplaceAll(src[0][1], " ", "")
}

func getRegexValueById(src [][]string, id int) string {
	if len(src) == 0 || len(src[0]) < id {
		return "0"
	}
	return strings.ReplaceAll(src[0][id], " ", "")
}

func getRegexValueByIndex(src [][]string, fid, sid int) string {
	if len(src) < fid+1 || len(src[fid]) < sid {
		return "0"
	}
	return strings.ReplaceAll(src[fid][sid], " ", "")
}

func mapToTask(src map[string]interface{}) Task {
	return Task{
		Id:       src["id"].(string),
		Sector:   src["sector"].(string),
		Worker:   src["worker"].(string),
		HostName: src["hostName"].(string),
		Task:     src["task"].(string),
		State:    src["state"].(string),
		Time:     src["time"].(string),
	}
}

// jobs "jobs":{"18":{"hostName":"worker01","id":"d7fd42c9","sector":"17","state":"running","task":"PC1","time":"17m48s","worker":"98c441ab"},}
// hardwareInfo  hardwareInfo":{"worker01":{"cpuLoad":"14.73","cpuTemper":"+41.1°C","diskR":"906.67M/s","diskW":"163.63M/s","gpuInfo":{"0":{"name":"0","temp":"91C","use":"100%"}},"hostName":"worker01","netIO":{"eno1":{"name":"eno1","rx":"1.27","tx":"2.90"},"eno2":{"name":"eno2","rx":"0.00","tx":"0.00"},"enp2s0f0np0":{"name":"enp2s0f0np0","rx":"0.00","tx":"0.00"},"enp2s0f1np1":{"name":"enp2s0f1np1","rx":"0.00","tx":"0.00"},"lo":{"name":"lo","rx":"0.00","tx":"0.00"}},"totalMemory":"503G","useDisk":"40%","useMemory":"319G"}},
// 根据 hostName 分组归纳信息
// jobs task信息 ; workerHardwareInfo 硬件列表信息

func MapParse(jobs, hardwareInfos map[string]interface{}) interface{} {

	mapByHostName := mapByHostName(jobs)

	mapByState := mapByState(mapByHostName)

	mapByType := mapByType(mapByState)

	var res []interface{}
	if hardwareInfos == nil || len(hardwareInfos) == 0 {
		for _, workerInfo := range mapByType {
			res = append(res, workerInfo)
		}
		return res
	}

	// 结合硬件信息
	for hostName, hardwareInfo := range hardwareInfos {
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
	if len(jobs) == 0 {
		return nil
	}
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
