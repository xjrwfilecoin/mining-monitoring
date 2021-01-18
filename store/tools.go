package store



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
func mapByHost(jobs []map[string]interface{}) map[string]interface{} {
	mapByHost := make(map[string]interface{})
	for i := 0; i < len(jobs); i++ {
		task := jobs[i]
		hostName, ok := task["hostName"]
		if !ok {
			continue
		}
		host := hostName.(string)
		taskList, ok := mapByHost[host]
		if ok {
			tasks := taskList.([]map[string]interface{})
			tasks = append(tasks, task)
		} else {
			mapByHost[host] = []map[string]interface{}{task}
		}
	}
	return mapByHost
}
