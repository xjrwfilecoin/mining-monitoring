package cache

import (
	"mining-monitoring/shellParsing"
	"reflect"
)

type Value struct {
	Value interface{}
	Flag  bool
}

func (v *Value) Update(flag bool, value interface{}) {
	v.Value = value
	v.Flag = flag
}

type WorkerId struct {
	MinerId  MinerId
	HostName string
}

type MinerId string

type WorkerInfo struct {
	HostName     Value `json:"hostName"`
	CurrentQueue Value `json:"currentQueue"`
	PendingQueue Value `json:"pendingQueue"`
	CpuTemper    Value `json:"cpuTemper"`
	CpuLoad      Value `json:"cpuLoad"`
	GpuInfo      Value `json:"gpuInfo"`
	TotalMemory  Value `json:"totalMemory"`
	UseMemory    Value `json:"useMemory"`
	UseDisk      Value `json:"useDisk"`
	DiskR        Value `json:"diskR"`
	DiskW        Value `json:"diskW"`
	NetIO        Value `json:"netIo"`

	NetState Value `json:"netState"` // ping心跳

	TaskState Value `json:"taskState"` // lotus-miner sealing  workers
	TaskType  Value `json:"taskType"`
}

// todo
func (w WorkerInfo) UpdateNetState(obj interface{}) {
	if netStateMap, ok := obj.(map[string]interface{}); ok {
		if netState, ok := netStateMap["netState"]; ok {
			if reflect.DeepEqual(w.NetState, netState) {
				w.NetState.Update(true, netState)
			}
		}
	}
}

// todo
func (w *WorkerInfo) UpdateTaskType(obj interface{}) {
	if taskStateInfo, ok := obj.(map[string]interface{}); ok {
		if taskState, ok := taskStateInfo["taskState"]; ok {
			if reflect.DeepEqual(w.TaskState, taskState) {
				w.TaskState.Update(true, taskState)
			}
		}
		if taskType, ok := taskStateInfo["taskType"]; ok {
			if reflect.DeepEqual(w.TaskType, taskType) {
				w.TaskType.Update(true, taskType)
			}
		}
	}
}

// todo
func (w *WorkerInfo) ChangeState(state bool) {
	w.HostName.Flag = true
	w.CurrentQueue.Flag = state
	w.PendingQueue.Flag = state
	w.CpuTemper.Flag = state
	w.CpuLoad.Flag = state
	w.GpuInfo.Flag = state
	w.TotalMemory.Flag = state
	w.UseMemory.Flag = state
	w.UseDisk.Flag = state
	w.DiskR.Flag = state
	w.DiskW.Flag = state
	w.NetIO.Flag = state

	w.NetState.Flag = state
	w.TaskState.Flag = state
	w.TaskType.Flag = state
}

// todo
func (w *WorkerInfo) updateDiskIO(obj interface{}) {
	if resMap, ok := obj.(map[string]interface{}); ok {
		if diskR, ok := resMap["diskR"]; ok {
			if reflect.DeepEqual(w.DiskR, diskR) {
				w.DiskR.Update(true, diskR)
			}
		}
		if diskW, ok := resMap["diskW"]; ok {
			if reflect.DeepEqual(w.DiskR, diskW) {
				w.DiskW.Update(true, diskW)
			}
		}

	}
}

//todo
func (w *WorkerInfo) updateMemory(obj interface{}) {
	if resMap, ok := obj.(map[string]interface{}); ok {
		if totalMemory, ok := resMap["totalMemory"]; ok {
			if reflect.DeepEqual(w.TotalMemory, totalMemory) {
				w.TotalMemory.Update(true, totalMemory)
			}
		}
		if userMemory, ok := resMap["useMemory"]; ok {
			if reflect.DeepEqual(w.UseMemory, userMemory) {
				w.UseMemory.Update(true, userMemory)
			}
		}

	}
}

// todo
func (w *WorkerInfo) updateJobQueue(obj interface{}) {
	if resMap, ok := obj.(map[string]interface{}); ok {
		if currentQueue, ok := resMap["currentQueue"]; ok {
			if reflect.DeepEqual(w.CurrentQueue, currentQueue) {
				w.CurrentQueue.Update(true, currentQueue)
			}
		}
		if pendingQueue, ok := resMap["pendingQueue"]; ok {
			if reflect.DeepEqual(w.PendingQueue, pendingQueue) {
				w.PendingQueue.Update(true, pendingQueue)
			}
		}

	}
}

func (w *WorkerInfo) Update(obj shellParsing.CmdData) {
	switch obj.CmdType {
	case shellParsing.IOCmd:
		w.updateDiskIO(obj.Data)
		break

	case shellParsing.SarCmd:
		isDiff := IsDiff(w.NetIO, obj.Data)
		w.NetIO.Update(isDiff, obj.Data)
		break

	case shellParsing.DfHCMd:
		isDiff := IsDiff(w.UseDisk, obj.Data)
		w.UseDisk.Update(isDiff, obj.Data)
		break

	case shellParsing.FreeHCmd:
		w.updateMemory(obj.Data)
		break

	case shellParsing.SensorsCmd:
		isDiff := IsDiff(w.CpuTemper, obj.Data)
		w.CpuTemper.Update(isDiff, isDiff)
		break

	case shellParsing.GpuCmd:
		isDiff := IsDiff(w.GpuInfo.Value, obj.Data)
		w.GpuInfo.Update(isDiff, obj.Data)
		break

	case shellParsing.UpTimeCmd:
		isDiff := IsDiff(w.CpuLoad, obj.Data)
		w.CpuLoad.Update(isDiff, obj.Data)
		break
	case shellParsing.LotusMinerJobs:
		w.updateJobQueue(obj.Data)
		break
	default:

	}
}

func (w WorkerInfo) GetDiff(all bool) map[string]interface{} {
	param := make(map[string]interface{})
	vw := reflect.ValueOf(w)
	tw := reflect.TypeOf(w)
	for i := 0; i < vw.NumField(); i++ {
		f := vw.Field(i)
		value := f.FieldByName("Value").Interface()
		if flag, ok := f.FieldByName("Flag").Interface().(bool); ok && flag || all {
			param[tw.Field(i).Name] = value
		}

	}
	return param
}

type MinerInfo struct {
	MinerId       Value `json:"minerId"`       // MinerId
	MinerBalance  Value `json:"minerBalance"`  // miner余额
	WorkerBalance Value `json:"workerBalance"` // worker余额
	PostBalance   Value `json:"postBalance"`

	PledgeBalance    Value `json:"pledgeBalance"`    // 抵押
	EffectivePower   Value `json:"effectivePower"`   // 有效算力
	TotalSectors     Value `json:"totalSectors"`     // 总扇区数
	EffectiveSectors Value `json:"effectiveSectors"` // 有效扇区
	ErrorSectors     Value `json:"errorSectors"`     // 错误扇区
	RecoverySectors  Value `json:"recoverySectors"`  // 恢复中扇区
	DeletedSectors   Value `json:"deletedSectors"`   // 删除扇区
	FailSectors      Value `json:"failSectors"`      // 失败扇区
	ExpectBlock      Value `json:"expectBlock"`      //  期望出块
	MinerAvailable   Value `json:"minerAvailable"`   //  miner可用余额
	PreCommitWait    Value `json:"preCommitWait"`    //  preCommitWait
	CommitWait       Value `json:"commitWait"`       //  commitWait
	PreCommit1       Value `json:"preCommit1"`       //  PreCommit1
	PreCommit2       Value `json:"preCommit2"`       //  PreCommit2
	WaitSeed         Value `json:"waitSeed"`         //  WaitSeed
	Committing       Value `json:"committing"`       //  Committing
	FinalizeSector   Value `json:"finalizeSector"`   //  finalizeSector

}

func (w *MinerInfo) Update(obj shellParsing.CmdData) {
	switch obj.CmdType {
	case shellParsing.LotusControlList:
		isDiff := IsDiff(w.PostBalance, obj.Data)
		w.PostBalance.Update(isDiff, obj.Data)
		break

	case shellParsing.LotusMinerInfoCmd:
		w.UpdateMinerInfo(obj.Data)
		break

	default:

	}
}

func (w MinerInfo) GetDiff(all bool) map[string]interface{} {
	param := make(map[string]interface{})
	vw := reflect.ValueOf(w)
	tw := reflect.TypeOf(w)
	for i := 0; i < vw.NumField(); i++ {
		f := vw.Field(i)
		value := f.FieldByName("Value").Interface()
		if flag, ok := f.FieldByName("Flag").Interface().(bool); ok && flag || all {
			param[tw.Field(i).Name] = value
		}
	}
	return param
}

// todo
func (w *MinerInfo) ChangeState(state bool) {
	w.MinerId.Flag = true
	w.MinerBalance.Flag = state
	w.WorkerBalance.Flag = state
	w.PostBalance.Flag = state

	w.PledgeBalance.Flag = state
	w.EffectivePower.Flag = state
	w.TotalSectors.Flag = state
	w.EffectiveSectors.Flag = state
	w.ErrorSectors.Flag = state
	w.RecoverySectors.Flag = state
	w.DeletedSectors.Flag = state
	w.FailSectors.Flag = state
	w.ExpectBlock.Flag = state
	w.MinerAvailable.Flag = state
	w.PreCommitWait.Flag = state
	w.CommitWait.Flag = state
	w.PreCommit1.Flag = state
	w.PreCommit2.Flag = state
	w.WaitSeed.Flag = state
	w.Committing.Flag = state
	w.FinalizeSector.Flag = state
}

// todo
func (w *MinerInfo) UpdateMinerInfo(obj interface{}) {
	if resMap, ok := obj.(map[string]interface{}); ok {
		if minerId, ok := resMap["minerId"]; ok {
			if reflect.DeepEqual(w.MinerId, minerId) {
				w.MinerId.Update(true, minerId)
			}
		}
		if minerBalance, ok := resMap["minerBalance"]; ok {
			if reflect.DeepEqual(w.MinerBalance, minerBalance) {
				w.MinerBalance.Update(true, minerBalance)
			}
		}
		if workerBalance, ok := resMap["workerBalance"]; ok {
			if reflect.DeepEqual(w.MinerBalance, workerBalance) {
				w.WorkerBalance.Update(true, workerBalance)
			}
		}
		if postBalance, ok := resMap["postBalance"]; ok {
			if reflect.DeepEqual(w.MinerBalance, postBalance) {
				w.PostBalance.Update(true, postBalance)
			}
		}
		if pledgeBalance, ok := resMap["pledgeBalance"]; ok {
			if reflect.DeepEqual(w.MinerBalance, pledgeBalance) {
				w.PledgeBalance.Update(true, pledgeBalance)
			}
		}
		if effectivePower, ok := resMap["effectivePower"]; ok {
			if reflect.DeepEqual(w.MinerBalance, effectivePower) {
				w.EffectivePower.Update(true, effectivePower)
			}
		}
		if totalSectors, ok := resMap["totalSectors"]; ok {
			if reflect.DeepEqual(w.MinerBalance, totalSectors) {
				w.TotalSectors.Update(true, totalSectors)
			}
		}
		if effectiveSectors, ok := resMap["effectiveSectors"]; ok {
			if reflect.DeepEqual(w.MinerBalance, effectiveSectors) {
				w.EffectiveSectors.Update(true, effectiveSectors)
			}
		}
		if errorSectors, ok := resMap["errorSectors"]; ok {
			if reflect.DeepEqual(w.MinerBalance, errorSectors) {
				w.ErrorSectors.Update(true, errorSectors)
			}
		}
		if recoverySectors, ok := resMap["recoverySectors"]; ok {
			if reflect.DeepEqual(w.MinerBalance, recoverySectors) {
				w.RecoverySectors.Update(true, recoverySectors)
			}
		}
		if deletedSectors, ok := resMap["deletedSectors"]; ok {
			if reflect.DeepEqual(w.MinerBalance, deletedSectors) {
				w.DeletedSectors.Update(true, deletedSectors)
			}
		}
		if failSectors, ok := resMap["failSectors"]; ok {
			if reflect.DeepEqual(w.MinerBalance, failSectors) {
				w.FailSectors.Update(true, failSectors)
			}
		}
		if deletedSectors, ok := resMap["deletedSectors"]; ok {
			if reflect.DeepEqual(w.MinerBalance, deletedSectors) {
				w.DeletedSectors.Update(true, deletedSectors)
			}
		}
		if expectBlock, ok := resMap["expectBlock"]; ok {
			if reflect.DeepEqual(w.MinerBalance, expectBlock) {
				w.ExpectBlock.Update(true, expectBlock)
			}
		}
		if minerAvailable, ok := resMap["minerAvailable"]; ok {
			if reflect.DeepEqual(w.MinerBalance, minerAvailable) {
				w.MinerAvailable.Update(true, minerAvailable)
			}
		}
		if preCommitWait, ok := resMap["preCommitWait"]; ok {
			if reflect.DeepEqual(w.MinerBalance, preCommitWait) {
				w.PreCommitWait.Update(true, preCommitWait)
			}
		}
		if commitWait, ok := resMap["commitWait"]; ok {
			if reflect.DeepEqual(w.MinerBalance, commitWait) {
				w.CommitWait.Update(true, commitWait)
			}
		}
		if preCommit1, ok := resMap["preCommit1"]; ok {
			if reflect.DeepEqual(w.MinerBalance, preCommit1) {
				w.PreCommit1.Update(true, preCommit1)
			}
		}
		if preCommit2, ok := resMap["preCommit2"]; ok {
			if reflect.DeepEqual(w.MinerBalance, preCommit2) {
				w.PreCommit2.Update(true, preCommit2)
			}
		}
		if waitSeed, ok := resMap["waitSeed"]; ok {
			if reflect.DeepEqual(w.MinerBalance, waitSeed) {
				w.WaitSeed.Update(true, waitSeed)
			}
		}
		if committing, ok := resMap["committing"]; ok {
			if reflect.DeepEqual(w.MinerBalance, committing) {
				w.Committing.Update(true, committing)
			}
		}
		if finalizeSector, ok := resMap["finalizeSector"]; ok {
			if reflect.DeepEqual(w.MinerBalance, finalizeSector) {
				w.FinalizeSector.Update(true, finalizeSector)
			}
		}
	}
}
