package model

import "fmt"

// 默认时间配置

const DefaultConnMaxNum = 1000

const DefaultMinerInfoInterval = 60
const DefaultMinerInfoTimeout = 60

const CommonMiningInterval = 45
const CommonMiningTimeout = 10

const HardWardInterval = 30
const HardWardTimeout = 15

const CmdConcurrentMaxNum = 300

type RuntimeConfig struct {
	Debug        bool // 是否是debug
	LogPath      string
	HTTPListen   string
	LogLevel     string // 日志等级
	ConnMaxNum   int64  // websocketIo 最大连接数
	CpuNum       int
	MinerConfigs []MinerConfig // 集群miner配置
}

type MinerConfig struct {
	MinerIp             string
	MinerId             string
	HardWardSampleTime  TimeValue // 硬件信息采样配置
	MinerInfoTime       TimeValue // minerInfo
	MiningInfoTime      TimeValue
	CmdConcurrentMaxNum int // 命令最大并发数
}

type TimeValue struct {
	Interval int // 单位秒
	Timeout  int // 单位秒
}

func (t *TimeValue) Check() error {
	if t.Interval <= 0 {
		return fmt.Errorf("inerval must more than zero")
	}

	if t.Timeout <= 0 {
		return fmt.Errorf("timeout must more than zero")
	}
	return nil

}

func (m *MinerConfig) Check() error {
	if err := m.HardWardSampleTime.Check(); err != nil {
		return fmt.Errorf("%v hardwardSampleTime  %v", m.MinerId, err)
	}
	if err := m.HardWardSampleTime.Check(); err != nil {
		return fmt.Errorf("%v MinerInfoTime %v", m.MinerId, err)
	}
	if err := m.HardWardSampleTime.Check(); err != nil {
		return fmt.Errorf("%v MiningInfoTime %v", m.MinerId, err)
	}
	if m.CmdConcurrentMaxNum <= 0 {
		return fmt.Errorf("%v CmdConcurrentMaxNum less than zero", m.MinerId)
	}
	return nil
}

func (r *RuntimeConfig) Check() error {
	if r.ConnMaxNum <= 0 {
		return fmt.Errorf("ConnMaxNum must more than zero")
	}
	for i := 0; i < len(r.MinerConfigs); i++ {
		minerConfig := r.MinerConfigs[i]
		if err := minerConfig.Check(); err != nil {
			return err
		}
	}
	return nil
}
