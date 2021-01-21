package store

import (
	"fmt"
	"mining-monitoring/shellParsing"
	"mining-monitoring/utils"
	"sync"
)

type MinerId string

type DeviceId string

type MinerInfo struct {
	MinerId    MinerId
	MiningInfo map[shellParsing.CmdType]shellParsing.CmdData
	Hardware   map[DeviceId]shellParsing.CmdData
	ml         sync.Mutex
	hl         sync.RWMutex
}

func NewMinerInfo(minerId MinerId) *MinerInfo {
	return &MinerInfo{
		MinerId:    minerId,
		MiningInfo: make(map[shellParsing.CmdType]shellParsing.CmdData),
		Hardware:   make(map[DeviceId]shellParsing.CmdData),
		ml:         sync.Mutex{},
		hl:         sync.RWMutex{},
	}
}

func (m *MinerInfo) updateData(obj shellParsing.CmdData) map[string]interface{} {
	switch obj.State {
	case shellParsing.LotusState:
		return m.miningInfo(obj)
	case shellParsing.HardwareState:
		return m.deviceInfo(obj)
	default:
	}
	return nil
}

func (m *MinerInfo) deviceInfo(obj shellParsing.CmdData) map[string]interface{} {
	devId := DeviceId(fmt.Sprintf("%v%v", obj.MinerId, obj.HostName))
	m.hl.Lock()
	cmdData, ok := m.Hardware[DeviceId(devId)]
	if !ok {
		m.Hardware[devId] = cmdData
	} else {
		m.Hardware[devId] = obj
	}
	m.hl.Unlock()
	return DiffMapValue(obj, cmdData)
}

func (m *MinerInfo) miningInfo(obj shellParsing.CmdData) map[string]interface{} {
	m.ml.Lock()
	cmdData, ok := m.MiningInfo[obj.CmdType]
	if !ok {
		m.MiningInfo[obj.CmdType] = cmdData
	} else {
		m.MiningInfo[obj.CmdType] = obj
	}
	m.ml.Unlock()
	return m.DiffMap(obj, cmdData)
}

func (m *MinerInfo) DiffMap(new, old shellParsing.CmdData) map[string]interface{} {
	if new.CmdType == "" {
		return nil
	}
	res := make(map[string]interface{})
	switch new.CmdType {
	case shellParsing.LotusMinerInfoCmd:
		res = DiffMapValue(new, old)
		break
	case shellParsing.LotusControlList:
		res = DiffMapValue(new, old)
		break
	case shellParsing.LotusMpoolCmd:
		res = DiffMapValue(new, old)
		break
	case shellParsing.LotusMinerWorkers:

		break
	case shellParsing.LotusMinerJobs:
		return ParseJobs(new.Data.([]map[string]interface{}))
	default:

	}

	return res
}

func ParseJobs(jobs []map[string]interface{}) map[string]interface{} {

	mapByHostName := mapByHost(jobs) // hostName:[{}]

	mapByState := mapByState(mapByHostName)

	mapByType := mapByType(mapByState)

	return mapByType
}

func DiffMapValue(new, old shellParsing.CmdData) map[string]interface{} {
	newMap := utils.StructToMapByJson(new.Data)
	oldMap := utils.StructToMapByJson(old.Data)
	return utils.DiffMap(oldMap, newMap)
}
