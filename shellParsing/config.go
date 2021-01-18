package shellParsing

type CmdType string

const (
	IOCmd      CmdType = "ioCmd"
	SarCmd     CmdType = "sarCmd"
	DfHCMd     CmdType = "dfhCmd"
	FreeHCmd   CmdType = "freeHCmd"
	SensorsCmd CmdType = "sensorsCmd"
	GpuCmd     CmdType = "gpuCmd"
	UpTimeCmd  CmdType = "upTimeCmd"
	//---------------------------------
	LotusMinerInfoCmd CmdType = "lotusMinerInfoCmd"
	LotusMpoolCmd     CmdType = "LotusMpoolCmd"
	LotusControlList  CmdType = "lotusMinerControlCmd"
	LotusMinerJobs    CmdType = "lotusMinerJobsCmd"
	LotusMinerWorkers CmdType = "lotusMinerWorkersCmd"
)

type CmdState int

const (
	LotusState = iota
	HardwareState
)
