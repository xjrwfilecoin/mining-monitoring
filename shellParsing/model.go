package shellParsing

import "time"

type Job struct {
	Id       string `json:"id"`
	Sector   string `json:"sector"` //扇区Id
	Worker   string `json:"worker"`
	HostName string `json:"hostName"`
	Task     string `json:"task"`  //任务类型
	State    string `json:"state"` // 任务状态
	Time     string `json:"time"`  // 耗时
}

type Worker struct {
	Hostname string
	Id       string
	Gpu      int
	CmdList  []ShellCmd
	close    chan struct{}
}

type Miner struct {
	MinerId string
	CmdList []ShellCmd
}

type ShellCmd struct {
	HostName string
	MinerId  string
	Name     string
	State    CmdState
	CmdType  CmdType
	Interval time.Duration
	Params   []string
	close    chan struct{}
}

func NewHardwareShellCmd(hostName, name string, cmdType CmdType, params []string) ShellCmd {
	return ShellCmd{HostName: hostName, Name: name, State: HardwareState, CmdType: cmdType, Params: params}
}

func NewLotusShellCmd(hostName, name string, cmdType CmdType, params []string) ShellCmd {
	return ShellCmd{HostName: hostName, Name: name, State: LotusState, CmdType: cmdType, Params: params}
}

type CmdData struct {
	HostName string
	MinerId  string
	CmdType  CmdType
	State    CmdState
	Data     interface{}
}

func NewCmdData(hostName, minerId string, cmdType CmdType, state CmdState, data interface{}) CmdData {
	return CmdData{HostName: hostName, MinerId: minerId, State: state, CmdType: cmdType, Data: data}
}

type CpuTemp struct {
	CpuTemp string `json:"cpuTemper"`
}

type GpuInfo struct {
	Name string `json:"name"`
	Temp string `json:"temp"`
	Use  string `json:"use"`
}

type Disk struct {
	UseDisk string `json:"useDisk"`
}

type Memory struct {
	UseMemory   string `json:"useMemory"`
	TotalMemory string `json:"totalMemory"`
}

type CpuLoad struct {
	CpuLoad string `json:"cpuLoad"`
}

type IoInfo struct {
	DiskR string `json:"diskR"`
	DiskW string `json:"diskW"`
}

type NetIO struct {
	Name string `json:"name"`
	Tx   string `json:"tx"`
	Rx   string `json:"rx"`
}

type WorkerInfo struct {
	HostName  string    `json:"hostName"`
	TaskType  []string  `json:"taskType"`
	TaskState TaskState `json:"taskState"`
	NetState  NetState  `json:"netState"`
	//IP        string    `json:"ip"`
	//Id        string    `json:"id"`
	GPU       int       `json:"gpu"`
}

type P map[string]interface{}

type PostBalance struct {
	PostBalance string `json:"postBalance"`
}

type MinerInfo struct {
	MinerId       string `json:"minerId"`       // MinerId
	MinerBalance  string `json:"minerBalance"`  // miner余额
	WorkerBalance string `json:"workerBalance"` // worker余额
	PledgeBalance string `json:"pledgeBalance"` // 抵押

	EffectivePower   string `json:"effectivePower"`   // 有效算力
	TotalSectors     string `json:"totalSectors"`     // 总扇区数
	EffectiveSectors string `json:"effectiveSectors"` // 有效扇区
	ErrorSectors     string `json:"errorSectors"`     // 错误扇区
	RecoverySectors  string `json:"recoverySectors"`  // 恢复中扇区
	DeletedSectors   string `json:"deletedSectors"`   // 删除扇区
	FailSectors      string `json:"failSectors"`      // 失败扇区

	ExpectBlock    string `json:"expectBlock"`    //  期望出块
	MinerAvailable string `json:"minerAvailable"` //  miner可用余额
	PreCommitWait  string `json:"preCommitWait"`  //  preCommitWait
	CommitWait     string `json:"commitWait"`     //  commitWait

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
	HostName    string      `json:"hostName"`
	CpuTemper   string      `json:"cpuTemper"` // cpu问题
	CpuLoad     string      `json:"cpuLoad"`   // cupu负载
	UseMemory   string      `json:"useMemory"` // 内存信息
	TotalMemory string      `json:"totalMemory"`
	UseDisk     string      `json:"useDisk"` // 磁盘使用率
	DiskR       string      `json:"diskR"`   //磁盘IO
	DiskW       string      `json:"diskW"`   //磁盘IO
	NetIO       interface{} `json:"netIO"`   //网络IO
	GpuInfo     interface{} `json:"gpu"`
}

func (hd *HardwareInfo) IsValid() bool {
	return hd.HostName != ""
}

type Sign struct {
	Type string
	Obj  interface{}
}

type NetCardIO struct {
	Name string `json:"name"`
	Rx   string `json:"rx"`
	TX   string `json:"tx"`
}

type GraphicsCardInfo struct {
	Name string `json:"name"`
	Temp string `json:"temp"`
	Use  string `json:"use"`
}
