package shellParsing

type WorkerInfo struct {
	HostName string
	IP       string
}

type Workers []WorkerInfo

type P map[string]interface{}

type MinerInfo struct {
	MinerId          string // MinerId
	MinerBalance     string // miner余额
	PostBalance      string // post余额
	WorkerBalance    string // worker余额
	PledgeBalance    string // 抵押
	TotalMessages    string // 消息总数
	RawBytePower     string // 字节算力
	AdjustedPower    string // 原值算力
	EffectivePower   string // 有效算力
	TotalSectors     string // 总扇区数
	EffectiveSectors string // 有效扇区
	ErrorSectors     string // 错误扇区
	RecoverySectors  string // 恢复中扇区
	DeletedSectors   string // 删除扇区
	FailSectors      string // 失败扇区
	Timestamp        string // 此次统计时间
}

// ID        Sector  Worker    Hostname       Task  State        Time
//c71e05fc  8598    74d84e37  ya_amd_node36  PC1   running      2h12m29.5s
//
type Task struct {
	Id       string
	Sector   string //扇区Id
	Worker   string
	HostName string
	Task     string //任务类型
	State    string // 任务状态
	Time     string // 耗时
}

type HardwareInfo struct {
	HostName    string
	CpuTemper   string // cpu问题
	CpuLoad     string // cupu负载
	GpuTemper   string // gpu温度
	GpuLoad     string // gpu负载
	UseMemory   string // 内存信息
	TotalMemory string
	UseDisk     string // 磁盘使用率
	DiskR       string //磁盘IO
	DiskW       string //磁盘IO
	NetRW       string //网络IO
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
