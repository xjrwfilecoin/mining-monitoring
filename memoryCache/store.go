package cache

type Value struct {
	Value interface{}
	Flag  bool
}

type WorkerId struct {
	MinerId  MinerId
	HostName string
}

type MinerId string

// miner info 表
type MinerInfoTable map[MinerId]MinerInfo

// worker信息表
type WorkerInfoTable map[WorkerId]WorkerInfo

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
}

type MinerInfo struct {
	MinerId          Value `json:"minerId"`          // MinerId
	MinerBalance     Value `json:"minerBalance"`     // miner余额
	WorkerBalance    Value `json:"workerBalance"`    // worker余额
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
