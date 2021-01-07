package shellParsing

type WorkerInfo struct {
	HostName string
	IP       string
}

type Workers []WorkerInfo

type P map[string]interface{}


type MinerInfo struct {
	MinerId          string `json:"minerId"`          // MinerId
	MinerBalance     string `json:"minerBalance"`     // miner余额
	PostBalance      string `json:"postBalance"`      // post余额
	WorkerBalance    string `json:"workerBalance"`    // worker余额
	PledgeBalance    string `json:"pledgeBalance"`    // 抵押
	TotalMessages    string `json:"totalMessages"`    // 消息总数
	RawBytePower     string `json:"rawBytePower"`     // 字节算力
	AdjustedPower    string `json:"adjustedPower"`    // 原值算力
	EffectivePower   string `json:"effectivePower"`   // 有效算力
	TotalSectors     string `json:"totalSectors"`     // 总扇区数
	EffectiveSectors string `json:"effectiveSectors"` // 有效扇区
	ErrorSectors     string `json:"errorSectors"`     // 错误扇区
	RecoverySectors  string `json:"recoverySectors"`  // 恢复中扇区
	DeletedSectors   string `json:"deletedSectors"`   // 删除扇区
	FailSectors      string `json:"failSectors"`      // 失败扇区
	Timestamp        string `json:"timestamp"`        // 此次统计时间
}

// ID        Sector  Worker    Hostname       Task  State        Time
//c71e05fc  8598    74d84e37  ya_amd_node36  PC1   running      2h12m29.5s
//
type Task struct {
	Id       string `json:"id"`
	Sector   string `json:"sector"` //扇区Id
	Worker   string `json:"worker"`
	HostName string `json:"hostName"`
	Task     string `json:"task"`  //任务类型
	State    string `json:"state"` // 任务状态
	Time     string `json:"time"`  // 耗时
}

type HardwareInfo struct {
	HostName    string `json:"hostName"`
	CpuTemper   string `json:"cpuTemper"` // cpu问题
	CpuLoad     string `json:"cpuLoad"`   // cupu负载
	GpuTemper   string `json:"gpuTemper"` // gpu温度
	GpuLoad     string `json:"gpuLoad"`   // gpu负载
	UseMemory   string `json:"useMemory"` // 内存信息
	TotalMemory string `json:"totalMemory"`
	UseDisk     string `json:"useDisk"` // 磁盘使用率
	DiskR       string `json:"diskR"`   //磁盘IO
	DiskW       string `json:"diskW"`   //磁盘IO
	NetRW       string `json:"netRW"`   //网络IO
}

func (hd *HardwareInfo) IsValid() bool {
	return hd.HostName != ""
}

type Sign struct {
	Type string
	Obj  interface{}
}

type Worker struct {
	Hostname     string
	CurrentQueue []Task // 当前任务
	PendingQueue []Task // 队列中任务

}
