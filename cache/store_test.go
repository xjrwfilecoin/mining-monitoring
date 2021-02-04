package cache

import (
	"fmt"
	"reflect"
	"testing"
)

func TestWorkerInfo(t *testing.T) {
	info := WorkerInfo{}
	diff := info.GetDiff(true)
	fmt.Println(diff)
	minerInfo := MinerInfo{}
	minerDiff := minerInfo.GetDiff(true)
	fmt.Println(minerDiff)
}

func TestWorkerInfo_GetDiff(t *testing.T) {
	info := WorkerInfo{}
	v := reflect.ValueOf(info)
	tf := reflect.TypeOf(info)
	count := v.NumField()
	for i := 0; i < count; i++ {
		f := v.Field(i)
		fmt.Println(tf.Field(i).Name)
		value := f.FieldByName("Value").Interface()
		flag := f.FieldByName("Flag").Interface().(bool)
		fmt.Println(value, flag)

	}
}
