package cache

import (
	"fmt"
	"mining-monitoring/log"
	"mining-monitoring/shell"
	"reflect"
	"strings"
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

func NewWorkerInfo(hostName string) *WorkerInfo {
	return &WorkerInfo{
		HostName:     Value{Value: hostName},
		CurrentQueue: Value{},
		PendingQueue: Value{},
		CpuTemper:    Value{},
		CpuLoad:      Value{},
		GpuInfo:      Value{},
		TotalMemory:  Value{},
		UseMemory:    Value{},
		UseDisk:      Value{},
		DiskR:        Value{},
		DiskW:        Value{},
		NetIO:        Value{},
		NetState:     Value{},
		TaskState:    Value{},
		TaskType:     Value{},
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
func (w *WorkerInfo) UpdateTaskType(obj interface{}) {
	log.Debug("workers state: ", obj)
	if taskStateInfo, ok := obj.(map[string]interface{}); ok {
		if taskState, ok := taskStateInfo["taskState"]; ok {
			if !DeepEqual(w.TaskState.Value, taskState) {
				w.TaskState.Update(true, taskState)
			}
		}
		if taskType, ok := taskStateInfo["taskType"]; ok {
			if !DeepEqual(w.TaskType.Value, taskType) {
				w.TaskType.Update(true, taskType)
			}
		}
		if netState, ok := taskStateInfo["netState"]; ok {
			if !DeepEqual(w.NetState.Value, netState) {
				w.NetState.Update(true, netState)
			}
		}
	}
}

// todo
func (w *WorkerInfo) updateDiskIO(obj interface{}) {
	if resMap, ok := obj.(map[string]interface{}); ok {
		if diskR, ok := resMap["diskR"]; ok {
			if !DeepEqual(w.DiskR.Value, diskR) {
				w.DiskR.Update(true, diskR)
			}
		}
		if diskW, ok := resMap["diskW"]; ok {
			if !DeepEqual(w.DiskW.Value, diskW) {
				w.DiskW.Update(true, diskW)
			}
		}

	}
}

//todo
func (w *WorkerInfo) updateMemory(obj interface{}) {
	if resMap, ok := obj.(map[string]interface{}); ok {
		if totalMemory, ok := resMap["totalMemory"]; ok {
			if !DeepEqual(w.TotalMemory.Value, totalMemory) {
				w.TotalMemory.Update(true, totalMemory)
			}
		}
		if userMemory, ok := resMap["useMemory"]; ok {
			if !DeepEqual(w.UseMemory.Value, userMemory) {
				w.UseMemory.Update(true, userMemory)
			}
		}

	}
}

// todo
func (w *WorkerInfo) updateJobQueue(obj interface{}) {
	if resMap, ok := obj.(map[string]interface{}); ok {
		if currentQueue, ok := resMap["currentQueue"]; ok {
			if !DeepEqual(w.CurrentQueue.Value, currentQueue) {
				w.CurrentQueue.Update(true, currentQueue)
			}
		}
		if pendingQueue, ok := resMap["pendingQueue"]; ok {
			if !DeepEqual(w.PendingQueue.Value, pendingQueue) {
				w.PendingQueue.Update(true, pendingQueue)
			}
		}

	}
}
func (w *WorkerInfo) updateCpuLoad(obj interface{}) {
	if resMap, ok := obj.(map[string]interface{}); ok {
		if cpuLoad, ok := resMap["cpuLoad"]; ok {
			if !DeepEqual(w.CpuLoad.Value, cpuLoad) {
				w.CpuLoad.Update(true, cpuLoad)
			}
		}
	}
}

func (w *WorkerInfo) updateCpuTemper(obj interface{}) {
	if resMap, ok := obj.(map[string]interface{}); ok {
		if cpuTemper, ok := resMap["cpuTemper"]; ok {
			if !DeepEqual(w.CpuTemper.Value, cpuTemper) {
				w.CpuTemper.Update(true, cpuTemper)
			}
		}
	}
}

func (w *WorkerInfo) updateUseDisk(obj interface{}) {
	if resMap, ok := obj.(map[string]interface{}); ok {
		if useDisk, ok := resMap["useDisk"]; ok {
			if !DeepEqual(w.UseDisk.Value, useDisk) {
				w.UseDisk.Update(true, useDisk)
			}
		}
	}
}

func (w *WorkerInfo) Update(obj shell.CmdData) {
	switch obj.CmdType {
	case shell.IOCmd:
		w.updateDiskIO(obj.Data)
		break

	case shell.SarCmd:
		isDiff := DeepEqual(w.NetIO, obj.Data)
		w.NetIO.Update(isDiff, obj.Data)
		break

	case shell.DfHCMd:
		w.updateUseDisk(obj.Data)
		break

	case shell.FreeHCmd:
		w.updateMemory(obj.Data)
		break

	case shell.SensorsCmd:
		w.updateCpuTemper(obj.Data)
		break

	case shell.GpuCmd:
		isDiff := DeepEqual(w.GpuInfo.Value, obj.Data)
		w.GpuInfo.Update(isDiff, obj.Data)
		break

	case shell.UpTimeCmd:
		w.updateCpuLoad(obj.Data)
		break
	case shell.LotusMinerJobs:
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
		if flag, ok := f.FieldByName("Flag").Interface().(bool); ok && flag && value != nil || all {
			keyName := tw.Field(i).Name
			if len(keyName) < 2 {
				continue
			}
			newKey := fmt.Sprintf("%v%v", strings.ToLower(keyName[0:1]), keyName[1:])
			param[newKey] = value
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

func NewMinerInfo(minerId string) *MinerInfo {
	return &MinerInfo{
		MinerId:       Value{Value: minerId},
		MinerBalance:  Value{},
		WorkerBalance: Value{},
		PostBalance:   Value{},

		PledgeBalance:    Value{},
		EffectivePower:   Value{},
		TotalSectors:     Value{},
		EffectiveSectors: Value{},
		ErrorSectors:     Value{},
		RecoverySectors:  Value{},
		DeletedSectors:   Value{},
		FailSectors:      Value{},
		ExpectBlock:      Value{},
		MinerAvailable:   Value{},
		PreCommitWait:    Value{},
		CommitWait:       Value{},
		PreCommit1:       Value{},
		PreCommit2:       Value{},
		WaitSeed:         Value{},
		Committing:       Value{},
		FinalizeSector:   Value{},
	}
}

func (w *MinerInfo) Update(obj shell.CmdData) {
	switch obj.CmdType {
	case shell.LotusControlList:
		w.UpdatePostBalance(obj.Data)
		break

	case shell.LotusMinerInfoCmd:
		w.UpdateMinerInfo(obj.Data)
		break

	default:

	}
}

func (w *MinerInfo) UpdatePostBalance(obj interface{}) {
	if resMap, ok := obj.(map[string]interface{}); ok {
		if postBalance, ok := resMap["postBalance"]; ok {
			if !DeepEqual(w.PostBalance.Value, postBalance) {
				w.PostBalance.Update(true, postBalance)
			}
		}
	}
}

func (w MinerInfo) GetDiff(all bool) map[string]interface{} {
	param := make(map[string]interface{})
	vw := reflect.ValueOf(w)
	tw := reflect.TypeOf(w)
	for i := 0; i < vw.NumField(); i++ {
		f := vw.Field(i)
		value := f.FieldByName("Value").Interface()
		if flag, ok := f.FieldByName("Flag").Interface().(bool); ok && flag && value != nil || all {
			keyName := tw.Field(i).Name
			if len(keyName) < 2 {
				continue
			}
			newKey := fmt.Sprintf("%v%v", strings.ToLower(keyName[0:1]), keyName[1:])
			param[newKey] = value
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
			if !DeepEqual(w.MinerId.Value, minerId) {
				w.MinerId.Update(true, minerId)
			}
		}
		if minerBalance, ok := resMap["minerBalance"]; ok {
			log.Debug("minerInfo: minerBalance ", w.MinerBalance, minerBalance)
			if !DeepEqual(w.MinerBalance.Value, minerBalance) {
				w.MinerBalance.Update(true, minerBalance)
			}
		}
		if workerBalance, ok := resMap["workerBalance"]; ok {
			if !DeepEqual(w.WorkerBalance.Value, workerBalance) {
				w.WorkerBalance.Update(true, workerBalance)
			}
		}
		if pledgeBalance, ok := resMap["pledgeBalance"]; ok {
			if !DeepEqual(w.PledgeBalance.Value, pledgeBalance) {
				w.PledgeBalance.Update(true, pledgeBalance)
			}
		}
		if effectivePower, ok := resMap["effectivePower"]; ok {
			if !DeepEqual(w.EffectivePower.Value, effectivePower) {
				w.EffectivePower.Update(true, effectivePower)
			}
		}
		if totalSectors, ok := resMap["totalSectors"]; ok {
			if !DeepEqual(w.TotalSectors.Value, totalSectors) {
				w.TotalSectors.Update(true, totalSectors)
			}
		}
		if effectiveSectors, ok := resMap["effectiveSectors"]; ok {
			if !DeepEqual(w.EffectiveSectors.Value, effectiveSectors) {
				w.EffectiveSectors.Update(true, effectiveSectors)
			}
		}
		if errorSectors, ok := resMap["errorSectors"]; ok {
			if !DeepEqual(w.ErrorSectors.Value, errorSectors) {
				w.ErrorSectors.Update(true, errorSectors)
			}
		}
		if recoverySectors, ok := resMap["recoverySectors"]; ok {
			if !DeepEqual(w.RecoverySectors.Value, recoverySectors) {
				w.RecoverySectors.Update(true, recoverySectors)
			}
		}
		if deletedSectors, ok := resMap["deletedSectors"]; ok {
			if !DeepEqual(w.DeletedSectors.Value, deletedSectors) {
				w.DeletedSectors.Update(true, deletedSectors)
			}
		}
		if failSectors, ok := resMap["failSectors"]; ok {
			if !DeepEqual(w.FailSectors.Value, failSectors) {
				w.FailSectors.Update(true, failSectors)
			}
		}
		if deletedSectors, ok := resMap["deletedSectors"]; ok {
			if !DeepEqual(w.DeletedSectors.Value, deletedSectors) {
				w.DeletedSectors.Update(true, deletedSectors)
			}
		}
		if expectBlock, ok := resMap["expectBlock"]; ok {
			if !DeepEqual(w.ExpectBlock.Value, expectBlock) {
				w.ExpectBlock.Update(true, expectBlock)
			}
		}
		if minerAvailable, ok := resMap["minerAvailable"]; ok {
			if !DeepEqual(w.MinerAvailable.Value, minerAvailable) {
				w.MinerAvailable.Update(true, minerAvailable)
			}
		}
		if preCommitWait, ok := resMap["preCommitWait"]; ok {
			if !DeepEqual(w.PreCommitWait.Value, preCommitWait) {
				w.PreCommitWait.Update(true, preCommitWait)
			}
		}
		if commitWait, ok := resMap["commitWait"]; ok {
			if !DeepEqual(w.CommitWait.Value, commitWait) {
				w.CommitWait.Update(true, commitWait)
			}
		}
		if preCommit1, ok := resMap["preCommit1"]; ok {
			if !DeepEqual(w.PreCommit1.Value, preCommit1) {
				w.PreCommit1.Update(true, preCommit1)
			}
		}
		if preCommit2, ok := resMap["preCommit2"]; ok {
			if !DeepEqual(w.PreCommit2.Value, preCommit2) {
				w.PreCommit2.Update(true, preCommit2)
			}
		}
		if waitSeed, ok := resMap["waitSeed"]; ok {
			if !DeepEqual(w.WaitSeed.Value, waitSeed) {
				w.WaitSeed.Update(true, waitSeed)
			}
		}
		if committing, ok := resMap["committing"]; ok {
			if !DeepEqual(w.Committing.Value, committing) {
				w.Committing.Update(true, committing)
			}
		}
		if finalizeSector, ok := resMap["finalizeSector"]; ok {
			if !DeepEqual(w.FinalizeSector.Value, finalizeSector) {
				w.FinalizeSector.Update(true, finalizeSector)
			}
		}
	}
}
