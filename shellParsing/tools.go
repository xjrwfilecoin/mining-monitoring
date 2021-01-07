package shellParsing

import "reflect"

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
			q1:=queue.([]Task)
			param := tasksByType(q1)
			result[taskType]=param
		}
		mapByTask[tHost] =result

	}

	// 结合硬件信息
	for i := 0; i < len(hardwareList); i++ {
		hardware := hardwareList[i]
		if info, ok := mapByTask[hardware.HostName]; ok {
			tp := info.(map[string]interface{})
			toMap := structToMap(&hardware)
			mapByTask[hardware.HostName] = mergeMaps(tp, toMap)
		}
	}
	return mapByTask
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

func structToMap(obj interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return m
	}
	elem := reflect.ValueOf(obj).Elem()
	relType := elem.Type()
	for i := 0; i < relType.NumField(); i++ {
		m[relType.Field(i).Name] = elem.Field(i).Interface()
	}
	return m
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

func getRegexValue(src [][]string) string {
	if len(src) == 0 || len(src[0]) == 0 {
		return ""
	}
	return src[0][1]
}

func getRegexValueById(src [][]string, id int) string {
	if len(src) == 0 || len(src[0]) < id {
		return ""
	}
	return src[0][id]
}
