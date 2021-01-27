package store

import (
	"mining-monitoring/shellParsing"
	"mining-monitoring/utils"
)

func LotusJobsToArrayV1(workers, obj shellParsing.CmdData) interface{} {

	workersMap := utils.StructToMapByJson(workers.Data)

	var workerList []map[string]interface{}
	response := NewMap()
	param := utils.StructToMapByJson(obj.Data)
	mapByHost := mapByHost(param)
	mapByState := mapByState(mapByHost)
	mapByType := mapByType(mapByState)

	for hostName, _ := range workersMap {
		if _, ok := mapByType[hostName]; !ok {
			newMap := NewMap()
			newMap["hostName"] = hostName
			newMap["currentQueue"] = nil
			newMap["pendingQueue"] = nil
			workerList = append(workerList, newMap)
		}
	}
	for _, value := range mapByType {
		if tv, ok := value.(map[string]interface{}); ok {
			workerList = append(workerList, tv)
		}
	}
	if len(workerList) == 0 {
		return nil
	}
	response["workerInfo"] = workerList
	return response
}

// lotus jobs 返回前端指定格式
func LotusJobsToArray(obj shellParsing.CmdData) interface{} {
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
	if len(workerList) == 0 {
		return nil
	}
	response["workerInfo"] = workerList
	return response
}

func MapToArray(obj shellParsing.CmdData) interface{} {
	param := utils.StructToMapByJson(obj.Data)
	var res []map[string]interface{}
	for key, value := range param {
		if item, ok := value.(map[string]interface{}); ok {
			item["name"] = key
			res = append(res, item)
		}
	}
	return res
}

// obj 为array
func NewCommArrayMap(hostName, key string, obj interface{}) interface{} {
	var workerList []map[string]interface{}
	response := NewMap()
	worker := NewMap()
	worker["hostName"] = hostName
	worker[key] = obj
	workerList = append(workerList, worker)
	response["workerInfo"] = workerList
	return response

}

func NewCommonMap(obj shellParsing.CmdData) interface{} {
	mapByJson := utils.StructToMapByJson(obj.Data)
	return mapByJson
}

func NewWorkerInfoMap(obj shellParsing.CmdData) interface{} {
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

func mapByHost(jobs map[string]interface{}) map[string]interface{} {
	mapByHost := make(map[string]interface{})
	for _, task := range jobs {
		if tvalue, ok := task.(map[string]interface{}); ok {
			hostName, ok := tvalue["hostName"]
			if !ok {
				continue
			}
			host := hostName.(string)
			taskList, ok1 := mapByHost[host]
			if ok1 {
				tasks := taskList.([]map[string]interface{})
				tasks = append(tasks, tvalue)
				mapByHost[host] = tasks
			} else {
				mapByHost[host] = []map[string]interface{}{tvalue}
			}
		}
	}
	return mapByHost
}
