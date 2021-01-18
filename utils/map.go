package utils

import (
	"encoding/json"
	"reflect"
)

// todo 优化 o^2
func DiffArrMap(oldMaps, newMaps []map[string]interface{}, diffKey string) []map[string]interface{} {
	var res []map[string]interface{}
	for i := 0; i < len(oldMaps); i++ {
		oldTempMap := oldMaps[i]
		oldTempValue, ok := oldTempMap[diffKey]
		if !ok {
			continue
		}
		for j := 0; j < len(newMaps); j++ {
			newTempMap := newMaps[j]
			newTempValue, ok := newTempMap[diffKey]
			if !ok {
				continue
			}
			if oldTempValue == newTempValue {
				diffMap := DiffMap(oldTempMap, newTempMap)
				if len(diffMap) > 0 {
					res = append(res, diffMap)
				}
			}
		}
	}
	return res
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
