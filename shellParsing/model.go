package shellParsing

type Worker struct {
	Hostname string
	Id       string
	Gpu      int
	CmdList  []ShellCmd
}

type Miner struct {
	MinerId string
	CmdList []ShellCmd
}

type ShellCmd struct {
	HostName string
	Name     string
	State    CmdState
	CmdType  CmdType
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

func NewCmdData(hostName string, cmdType CmdType, state CmdState, data interface{}) CmdData {
	return CmdData{HostName: hostName, MinerId: hostName, State: state, CmdType: cmdType, Data: data}
}

type CpuTemp struct {
	Temp string `json:"temp"`
}

type GpuInfo struct {
	Name string `json:"name"`
	Temp string `json:"temp"`
	Used string `json:"used"`
}

type Disk struct {
	Used string `json:"used"`
}

type Memory struct {
	Used  string `json:"used"`
	Total string `json:"total"`
}

type CpuLoad struct {
	Load string `json:"load"`
}

type IoInfo struct {
	ReadIO  string `json:"readIo"`
	WriteIO string `json:"writeIo"`
}

type NetIO struct {
	Name string `json:"name"`
	Tx   string `json:"tx"`
	Rx   string `json:"rx"`
}

type WorkerInfo struct {
	HostName string
	IP       string
	Id       string
	GPU      int
}



type P map[string]interface{}

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
	Timestamp        int64  `json:"timestamp"`        // 此次统计时间

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
