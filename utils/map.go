package utils

import (
	"encoding/json"
	"reflect"
)

// 判断 map是否有差异
func MapIsDiff(oldMap, newMap map[string]interface{}) bool {
	for key, value := range newMap {
		if reflect.TypeOf(value).Kind() == reflect.Map {
			if _, ok := oldMap[key]; !ok {
				return true
			} else {
				tempOldMap, oOk := oldMap[key].(map[string]interface{})
				if !oOk {
					continue
				}
				tempNewMap, nOK := value.(map[string]interface{})
				if !nOK {
					continue
				}
				diffMap := DeepDiffMap(tempOldMap, tempNewMap)
				if diffMap != nil && len(diffMap) != 0 {
					return true
				}
			}
		} else {
			if tv, ok := oldMap[key]; !ok || value != tv {
				return true
			}
		}
	}
	return false
}

func SimpleDiffMap(oldMap, newMap map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range newMap {
		if tv, ok := oldMap[key]; !ok || value != tv {
			result[key] = value
		}
	}
	return result
}

// 两个map的差集
func DeepDiffMap(oldMap, newMap map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range newMap {

		if reflect.TypeOf(value).Kind() == reflect.Map {
			if _, ok := oldMap[key]; !ok {
				result[key] = value
			} else {
				tempOldMap, oOk := oldMap[key].(map[string]interface{})
				if !oOk {
					continue
				}
				tempNewMap, nOK := value.(map[string]interface{})
				if !nOK {
					continue
				}
				diffMap := DeepDiffMap(tempOldMap, tempNewMap)
				if diffMap != nil && len(diffMap) != 0 {
					result[key] = diffMap
				}
			}
		} else {
			if tv, ok := oldMap[key]; !ok || value != tv {
				result[key] = value
			}
		}
	}
	return result
}

func StructToMapByJson(obj interface{}) map[string]interface{} {
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

func MergeMaps(maps ...map[string]interface{}) map[string]interface{} {
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
