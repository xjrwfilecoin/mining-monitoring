package shell

type Option struct {
	HardWardSampleTime  Item   //硬件采样率
	MinerInfoTime       Item   // minerInfo 比较耗时，需要单独配置
	MiningInfoTime      Item   // 其它通用配置
	CmdConcurrentMaxNum uint64 // 最大协程并发数
}

type Item struct {
	Interval uint64 // 间隔时间
	Timeout  uint64 //超时间
}
