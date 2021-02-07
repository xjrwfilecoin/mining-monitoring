package model

/*

  "Debug": false,
  "Config": {
    "HardWardSample": {
      "Interval": 120,
      "Timeout": 60
    },
    "MinerInfo": {
      "Interval": 120,
      "Timeout": 60
    },
    "MiningInfo": {
      "Interval": 120,
      "Timeout": 60
    },
    "CmdConcurrentMaxNum": 100,
    "ConnMaxNum": 1000
  }
*/
type RuntimeConfig struct {
	Debug               bool // 是否是debug
	LogPath             string
	HTTPListen          string
	LogLevel            string // 日志等级
	HardWardSample      Value  // 硬件信息采样配置
	MinerInfo           Value
	MiningInfo          Value
	CmdConcurrentMaxNum int
	ConnMaxNum          int
}

type Value struct {
	Interval int
	Timeout  int
}
