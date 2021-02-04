package cache

import "reflect"

func IsDiff(old, new interface{}) bool {
	if old == nil || new == nil {
		return true
	}
	return reflect.DeepEqual(old, new)
}
