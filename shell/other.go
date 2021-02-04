package shell

func (sp *Parse) initLotusCmdData(minerId string) {
	sp.CmdMap[LotusMinerInfoCmd] = NewLotusShellCmd(minerId, "lotus-miner", LotusMinerInfoCmd, []string{"info"})
	sp.CmdMap[LotusControlList] = NewLotusShellCmd(minerId, "lotus-miner", LotusControlList, []string{"actor", "control", "list"})
	sp.CmdMap[LotusMinerJobs] = NewLotusShellCmd(minerId, "lotus-miner", LotusMinerJobs, []string{"sealing", "jobs"})
	sp.CmdMap[LotusMinerWorkers] = NewLotusShellCmd(minerId, "lotus-miner", LotusMinerWorkers, []string{"sealing", "workers"})
	sp.CmdMap[LotusMpoolCmd] = NewLotusShellCmd(minerId, "lotus", LotusMpoolCmd, []string{"mpool", "pending"})
}

func (sp *Parse) initHardwareCmd(hostName, execInfo string) {
	sp.CmdMap[IOCmd] = NewHardwareShellCmd(hostName, "ssh", IOCmd, []string{execInfo, "iotop", "-bn1", "|", "head", "-n", "2"})
	sp.CmdMap[SarCmd] = NewHardwareShellCmd(hostName, "ssh", SarCmd, []string{execInfo, "sar", "-n", "DEV", "1", "2"})
	sp.CmdMap[DfHCMd] = NewHardwareShellCmd(hostName, "ssh", DfHCMd, []string{execInfo, "df", "-h"})
	sp.CmdMap[FreeHCmd] = NewHardwareShellCmd(hostName, "ssh", FreeHCmd, []string{execInfo, "free", "-h"})
	sp.CmdMap[SensorsCmd] = NewHardwareShellCmd(hostName, "ssh", SensorsCmd, []string{execInfo, "sensors"})
	sp.CmdMap[GpuCmd] = NewHardwareShellCmd(hostName, "ssh", GpuCmd, []string{execInfo, "nvidia-smi"})
	sp.CmdMap[UpTimeCmd] = NewHardwareShellCmd(hostName, "ssh", UpTimeCmd, []string{execInfo, "uptime"})
}
