package model

type RuntimeConfig struct {
	Debug        bool         // 是否是debug
	LogPath      string
	HTTPListen   string
	LogLevel     string       // 日志等级
	ConnMaxNum   int
	MinerConfigs []MinerConfig // 集群miner配置
}

type MinerConfig struct {
	APIUrl              string
	APIToken            string
	Workers             [] Worker
	HardWardSampleTime  Value // 硬件信息采样配置
	MinerInfoTime       Value
	MiningInfoTime      Value
	CmdConcurrentMaxNum int
}

type Worker struct {
	HostName string
	Ip       string
	HostType string
	TaskType string
}

type Value struct {
	Interval int
	Timeout  int
}
