package cache

import (
	"fmt"
	"testing"
)

func TestMinerInfo(t *testing.T) {
	res1 := make(map[string]interface{})
	res2 := make(map[string]interface{})
	equal := DeepEqual(res1, res2)
	fmt.Println(equal)
}

type TestMiner struct {
	MinerId    string
	WorkerInfo []*WorkerInfo
}

func TestWorkerInfo(t *testing.T) {

	var workerList01 []*WorkerInfo
	var workerList02 []*WorkerInfo

	info := &WorkerInfo{
		HostName:     Value{Value: 1},
		CurrentQueue: Value{Value: 2},
		PendingQueue: Value{Value: 3},
		CpuTemper:    Value{Value: 4},
		CpuLoad:      Value{Value: 5},
		GpuInfo:      Value{Value: 5},
		TotalMemory:  Value{Value: 7},
		UseMemory:    Value{Value: 8},
		UseDisk:      Value{Value: 9},
		DiskR:        Value{Value: 10},
		DiskW:        Value{Value: 11},
		NetIO:        Value{Value: 12},

		NetState: Value{Value: 13},

		TaskState: Value{Value: 14},
		TaskType:  Value{Value: 15},
	}
	info01 := &WorkerInfo{
		HostName:     Value{Value: 1},
		CurrentQueue: Value{Value: 2},
		PendingQueue: Value{Value: 3},
		CpuTemper:    Value{Value: 4},
		CpuLoad:      Value{Value: 5},
		GpuInfo:      Value{Value: 5},
		TotalMemory:  Value{Value: 7},
		UseMemory:    Value{Value: 8},
		UseDisk:      Value{Value: 9},
		DiskR:        Value{Value: 10},
		DiskW:        Value{Value: 11},
		NetIO:        Value{Value: 12},

		NetState: Value{Value: 13},

		TaskState: Value{Value: 14},
		TaskType:  Value{Value: 10},
	}
	workerList01 = append(workerList01, info)
	workerList02 = append(workerList02, info01)

	equal := DeepEqual(workerList01, workerList02)

	fmt.Println(equal)
	testMiner01 := TestMiner{MinerId: "100", WorkerInfo: workerList01}
	testMiner02 := TestMiner{MinerId: "100", WorkerInfo: workerList02}

	deepEqual := DeepEqual(testMiner01, testMiner02)
	fmt.Println(deepEqual)
}

func TestWorkerInfo_GetDiff(t *testing.T) {
	minerInfo := MinerInfo{}
	minerDiff := minerInfo.GetDiff(true)
	fmt.Println(minerDiff)

}
