package shell

var SSHPORT = "22521"
//var SSHPORT = "22"

var SSHUSER = "root"

type TaskState int

type NetState int

const (
	Normal TaskState = iota
	TaskDisabled
)

const (
	NetNormal NetState = iota
	NetDisabled
)

type CmdType string

const (
	IOCmd      CmdType = "ioCmd"
	SarCmd     CmdType = "sarCmd"
	DfHCMd     CmdType = "dfhCmd"
	FreeHCmd   CmdType = "freeHCmd"
	SensorsCmd CmdType = "sensorsCmd"
	GpuCmd     CmdType = "gpuCmd"
	UpTimeCmd  CmdType = "upTimeCmd"

	GpuEnable CmdType = "gpuEnable"
	//---------------------------------
	LotusMinerInfoCmd CmdType = "lotusMinerInfoCmd"
	LotusMpoolCmd     CmdType = "LotusMpoolCmd"
	LotusControlList  CmdType = "lotusMinerControlCmd"
	LotusMinerJobs    CmdType = "lotusMinerJobsCmd"
	LotusMinerWorkers CmdType = "lotusMinerWorkersCmd"
)

type CmdState string

const (
	LotusState    = "lotus-state"
	HardwareState = "hardware-state"
)
