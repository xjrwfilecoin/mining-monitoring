package shellParsing

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mining-monitoring/log"
	"mining-monitoring/utils"
	"os/exec"
	"strings"
	"time"
)

type ShellParse struct {
	Workers      []*WorkerInfo
	Miners       Miner
	cmdSign      chan CmdData
	CmdParseMap  map[CmdType]func(cmd ShellCmd, input string) CmdData
	cmdHeartTime time.Duration // 秒
	closing      chan struct{}
	CmdMap       map[CmdType]ShellCmd
	workerSign   chan Worker
}

func NewShellParse() *ShellParse {
	return &ShellParse{
		cmdSign:      make(chan CmdData, 1000),
		CmdParseMap:  make(map[CmdType]func(cmd ShellCmd, input string) CmdData),
		closing:      make(chan struct{}),
		CmdMap:       make(map[CmdType]ShellCmd),
		workerSign:   make(chan Worker),
		cmdHeartTime: 10,
	}
}

func (sp *ShellParse) Close() {
	close(sp.closing)
	close(sp.cmdSign)
}

func (sp *ShellParse) initCmdParse() {
	sp.CmdParseMap[IOCmd] = sp.ExecIOTopCmd
	sp.CmdParseMap[SarCmd] = sp.ExecSarNetIOCmd
	sp.CmdParseMap[DfHCMd] = sp.ExecDfHCmd
	sp.CmdParseMap[FreeHCmd] = sp.ExecFreeHCmd
	sp.CmdParseMap[SensorsCmd] = sp.ExecSensorsCmd
	sp.CmdParseMap[GpuCmd] = sp.ExecGPUCmd
	sp.CmdParseMap[UpTimeCmd] = sp.ExecUptimeCmd

	sp.CmdParseMap[LotusMinerInfoCmd] = sp.ExecLotusMinerInfo
	sp.CmdParseMap[LotusControlList] = sp.ExecLotusPostInfo
	sp.CmdParseMap[LotusMinerJobs] = sp.ExecLotusMinerJobs
	sp.CmdParseMap[LotusMinerWorkers] = sp.ExecLotusWorkers
	sp.CmdParseMap[LotusMpoolCmd] = sp.ExecLotusMpoolInfo
}

func (sp *ShellParse) getMinerCmdList(minerId string) []ShellCmd {
	var cmdList []ShellCmd
	cmdList = append(cmdList, NewLotusShellCmd(minerId, "lotus-miner", LotusMinerInfoCmd, []string{"info"}))
	//cmdList = append(cmdList, NewLotusShellCmd(minerId, "lotus-miner", LotusControlList, []string{"actor", "control", "list"}))
	cmdList = append(cmdList, NewLotusShellCmd(minerId, "lotus-miner", LotusMinerJobs, []string{"sealing", "jobs"}))
	//cmdList = append(cmdList, NewLotusShellCmd(minerId, "lotus-miner", LotusMinerWorkers, []string{"sealing", "workers"}))
	//cmdList = append(cmdList, NewLotusShellCmd(minerId, "lotus", LotusMpoolCmd, []string{"mpool", "pending"}))
	return cmdList
}

// sshpass -p 1 ssh root@xjrw_node01 "free -h"
func (sp *ShellParse) getWorkCmdList(hostName string, gpuEnable bool) []ShellCmd {
	execInfo := fmt.Sprintf(`root@%v`, hostName)
	var cmdList []ShellCmd
	cmdList = append(cmdList, NewHardwareShellCmd(hostName, "sshpass", SensorsCmd, []string{"-p", "", "ssh", execInfo, "sensors"}))
	cmdList = append(cmdList, NewHardwareShellCmd(hostName, "sshpass", UpTimeCmd, []string{"-p", "", "ssh", execInfo, "uptime"}))
	cmdList = append(cmdList, NewHardwareShellCmd(hostName, "sshpass", FreeHCmd, []string{"-p", "", "ssh", execInfo, "free", "-h"}))
	cmdList = append(cmdList, NewHardwareShellCmd(hostName, "sshpass", DfHCMd, []string{"-p", "", "ssh", execInfo, "df", "-h"}))
	cmdList = append(cmdList, NewHardwareShellCmd(hostName, "sshpass", SarCmd, []string{"-p", "", "ssh", execInfo, "sar", "-n", "DEV", "1", "2"}))
	cmdList = append(cmdList, NewHardwareShellCmd(hostName, "sshpass", IOCmd, []string{"-p", "", "ssh", execInfo, "iotop", "-bn1", "|", "head", "-n", "2"}))
	if gpuEnable {
		cmdList = append(cmdList, NewHardwareShellCmd(hostName, "sshpass", GpuCmd, []string{"-p", "", "ssh", execInfo, "nvidia-smi"}))
	}
	return cmdList
}

func (sp *ShellParse) getWorkCmdListV1(hostName string, gpuEnable bool) []ShellCmd {
	execInfo := fmt.Sprintf(`root@%v`, hostName)
	var cmdList []ShellCmd
	cmdList = append(cmdList, NewHardwareShellCmd(hostName, "ssh", SensorsCmd, []string{execInfo, "sensors"}))
	cmdList = append(cmdList, NewHardwareShellCmd(hostName, "ssh", UpTimeCmd, []string{execInfo, "uptime"}))
	cmdList = append(cmdList, NewHardwareShellCmd(hostName, "ssh", FreeHCmd, []string{execInfo, "free", "-h"}))
	cmdList = append(cmdList, NewHardwareShellCmd(hostName, "ssh", DfHCMd, []string{execInfo, "df", "-h"}))
	cmdList = append(cmdList, NewHardwareShellCmd(hostName, "ssh", SarCmd, []string{execInfo, "sar", "-n", "DEV", "1", "2"}))
	cmdList = append(cmdList, NewHardwareShellCmd(hostName, "ssh", IOCmd, []string{execInfo, "iotop", "-bn1", "|", "head", "-n", "2"}))
	if gpuEnable {
		cmdList = append(cmdList, NewHardwareShellCmd(hostName, "ssh", GpuCmd, []string{execInfo, "nvidia-smi"}))
	}
	return cmdList
}

func (sp *ShellParse) doMinerInfo() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			cmdList := sp.getMinerCmdList(sp.Miners.MinerId)
			for i := 0; i < len(cmdList); i++ {
				cmd := cmdList[i]
				fn := sp.CmdParseMap[cmd.CmdType]
				sp.processTask(cmd, sp.cmdSign, fn)
			}
		case <-sp.closing:
			return
		default:


		}
	}
}

func (sp *ShellParse) doWorkers() {
	if err := recover(); err != nil {
		log.Error(err)
	}
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := sp.getWorkerList()
			if err != nil {
				log.Error("get worker info error: ", err.Error())
			}
		case <-sp.closing:
			return
		default:

		}
	}
}

func (sp *ShellParse) PingWorkers() {
	var workers []*WorkerInfo
	for i := 0; i < len(sp.Workers); i++ {
		worker := sp.Workers[i]
		execInfo := fmt.Sprintf(`root@%v`, worker.HostName)
		cmd := NewHardwareShellCmd(worker.HostName, "sshpass", SensorsCmd, []string{"-p", "", "ssh", execInfo, "free", "-h"})
		err := sp.execShellCmd(cmd, func(input string) {
			if !strings.Contains(input, "exit") {
				workers = append(workers, worker)
			}
		})
		if err != nil {
			log.Error("ping workers: ", err.Error())
		}

	}
	sp.Workers = workers
}

// todo 锁
func (sp *ShellParse) getWorkerList() error {
	minerWorkersCmd := NewLotusShellCmd("", "lotus-miner", LotusMinerWorkers, []string{"sealing", "workers"})
	err := sp.execShellCmd(minerWorkersCmd, func(input string) {
		workers := sp.GetMinerWorkersV2(input)
		sp.Workers = workers

	})
	if err != nil {
		return err
	}
	sp.PingWorkers()
	return nil
}

func (sp *ShellParse) doHardWareInfo() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for i := 0; i < len(sp.Workers); i++ {
				worker := sp.Workers[i]
				sp.runWorkerCmdList(worker, sp.cmdSign)
			}
		case <-sp.closing:
			return

		default:

		}
	}
}

func (sp *ShellParse) runWorkerCmdList(worker *WorkerInfo, sing chan CmdData) {
	if worker == nil {
		return
	}
	cmdList := sp.getWorkCmdList(worker.HostName, worker.GPU != 0)
	for i := 0; i < len(cmdList); i++ {
		cmd := cmdList[i]
		fn := sp.CmdParseMap[cmd.CmdType]
		go sp.processTask(cmd, sing, fn)
	}
}

func (sp *ShellParse) needGPU(worker WorkerInfo01) bool {
	if worker.HostName == "" {
		return false
	}
	for i := 0; i < len(worker.Jobs); i++ {
		job := worker.Jobs[i]
		if job.Task == "PC1" || job.Task == "C2" {
			return true
		}
	}
	return false

}

func (sp *ShellParse) getMiner() error {
	minerInfoCmd := NewLotusShellCmd("", "lotus-miner", LotusMinerInfoCmd, []string{"info"})
	err := sp.execShellCmd(minerInfoCmd, func(input string) {
		minerInfo := sp.getMinerInfo(input)
		minerCmdList := sp.getMinerCmdList(minerInfo.MinerId)
		sp.Miners = Miner{MinerId: minerInfo.MinerId, CmdList: minerCmdList}
	})
	return err
}

func (sp *ShellParse) Send() {
	sp.initCmdParse()
	err := sp.getMiner()
	if err != nil {
		panic(fmt.Errorf("lotus-miner cmd not available %v ", err.Error()))
	}
	err = sp.getWorkerList()
	if err != nil {
		panic(fmt.Errorf("check worker is avaibale: %v ", err.Error()))
	}
	go sp.doWorkers()
	go sp.doMinerInfo()
	go sp.doHardWareInfo()
}

func (sp *ShellParse) Receiver(recv chan CmdData) {
	for {
		select {
		case obj := <-sp.cmdSign:
			data, err := json.Marshal(obj)
			if err != nil {
				log.Error("json Marshal ", err.Error())
				continue
			}
			log.Debug("receiver info ", string(data))
			recv <- obj

		default:

		}
	}
}

func (sp *ShellParse) execShellCmd(cmd ShellCmd, fn func(input string)) error {
	data, err := sp.ExecCmd(cmd.Name, cmd.Params...)
	if err != nil {
		return err
	}
	fn(data)
	return nil
}

func (sp *ShellParse) processTask(cmd ShellCmd, sign chan CmdData, fn func(cmd ShellCmd, input string) CmdData) {
	output, err := sp.ExecCmd(cmd.Name, cmd.Params...)
	if err != nil {
		log.Error("process task error ", cmd.CmdType, cmd.Name, cmd.HostName, err)
		return
	}
	cmdData := fn(cmd, output)
	sign <- cmdData

}

func (sp *ShellParse) ExecLotusWorkers(cmd ShellCmd, data string) CmdData {
	workerInfos := sp.GetMinerWorkersV2(data)
	return NewCmdData(cmd.HostName, sp.Miners.MinerId, cmd.CmdType, cmd.State, workerInfos)
}

func (sp *ShellParse) ExecLotusMinerJobs(cmd ShellCmd, data string) CmdData {
	tasks := sp.GetMinerJobsCV2(data)
	return NewCmdData(cmd.HostName, sp.Miners.MinerId, cmd.CmdType, cmd.State, tasks)
}

func (sp *ShellParse) ExecLotusPostInfo(cmd ShellCmd, data string) CmdData {
	postBalance := postBalanceTestReg.FindAllStringSubmatch(data, 1)
	postValue := getRegexValue(postBalance)
	return NewCmdData(cmd.HostName, sp.Miners.MinerId, cmd.CmdType, cmd.State, PostBalance{PostBalance: postValue})
}

func (sp *ShellParse) ExecLotusMpoolInfo(cmd ShellCmd, data string) CmdData {
	count := strings.Count(data, "Message")
	return NewCmdData(cmd.HostName, sp.Miners.MinerId, cmd.CmdType, cmd.State, count)
}

func (sp *ShellParse) ExecLotusMinerInfo(cmd ShellCmd, data string) CmdData {
	minerInfo := sp.getMinerInfo(data)
	return NewCmdData(cmd.HostName, sp.Miners.MinerId, cmd.CmdType, cmd.State, minerInfo)
}

func (sp *ShellParse) getMinerInfo(data string) MinerInfo {
	minerInfo := MinerInfo{}
	minerId := minerIdReg.FindAllStringSubmatch(data, 1)
	minerInfo.MinerId = getRegexValue(minerId)
	minerBalance := minerBalanceReg.FindAllStringSubmatch(data, 1)
	minerInfo.MinerBalance = getRegexValue(minerBalance)
	workerBalance := workerBalanceReg.FindAllStringSubmatch(data, 1)
	minerInfo.WorkerBalance = getRegexValue(workerBalance)
	pledgeBalance := pledgeBalanceReg.FindAllStringSubmatch(data, 1)
	minerInfo.PledgeBalance = getRegexValue(pledgeBalance)
	totalPower := totalPowerReg.FindAllStringSubmatch(data, 1)
	minerInfo.EffectivePower = getRegexValue(totalPower)
	effectPower := effectPowerReg.FindAllStringSubmatch(data, 1)
	minerInfo.EffectivePower = getRegexValue(effectPower)
	totalSectors := totalSectorsReg.FindAllStringSubmatch(data, 1)
	minerInfo.TotalSectors = getRegexValue(totalSectors)
	effectSectors := effectSectorReg.FindAllStringSubmatch(data, 2)
	minerInfo.EffectiveSectors = getRegexValueByIndex(effectSectors, 1, 1)
	errorsSectors := errorSectorReg.FindAllStringSubmatch(data, 1)
	minerInfo.ErrorSectors = getRegexValue(errorsSectors)
	recoverySectors := recoverySectorReg.FindAllStringSubmatch(data, 1)
	minerInfo.RecoverySectors = getRegexValue(recoverySectors)
	deletedSectors := deletedSectorReg.FindAllStringSubmatch(data, 1)
	minerInfo.DeletedSectors = getRegexValue(deletedSectors)
	failSectors := failSectorReg.FindAllStringSubmatch(data, 1)
	minerInfo.FailSectors = getRegexValue(failSectors)
	return minerInfo
}

func (sp *ShellParse) ExecGPUCmd(cmd ShellCmd, input string) CmdData {
	gpuInfos := getGraphicsCardInfoV2(input)
	return NewCmdData(cmd.HostName, sp.Miners.MinerId, cmd.CmdType, cmd.State, gpuInfos)
}

func (sp *ShellParse) ExecSensorsCmd(cmd ShellCmd, output string) CmdData {
	cpuTemp := getCpuTemperV2(output)
	return NewCmdData(cmd.HostName, sp.Miners.MinerId, cmd.CmdType, cmd.State, cpuTemp)

}

func (sp *ShellParse) ExecDfHCmd(cmd ShellCmd, input string) CmdData {
	diskUsed := diskUsedRateReg.FindAllStringSubmatch(input, 1)
	diskInfo := Disk{UseDisk: getRegexValue(diskUsed)}
	return NewCmdData(cmd.HostName, sp.Miners.MinerId, cmd.CmdType, cmd.State, diskInfo)

}

func (sp *ShellParse) ExecFreeHCmd(cmd ShellCmd, input string) CmdData {

	memoryUsed := memoryUsedReg.FindAllStringSubmatch(input, 1)
	memory := Memory{UseMemory: getRegexValueById(memoryUsed, 2), TotalMemory: getRegexValueById(memoryUsed, 1)}
	return NewCmdData(cmd.HostName, sp.Miners.MinerId, cmd.CmdType, cmd.State, memory)

}

func (sp *ShellParse) ExecUptimeCmd(cmd ShellCmd, input string) CmdData {
	cpuLoad := cpuLoadReg.FindAllStringSubmatch(input, 1)
	load := CpuLoad{CpuLoad: getRegexValue(cpuLoad)}
	return NewCmdData(cmd.HostName, sp.Miners.MinerId, cmd.CmdType, cmd.State, load)
}

func (sp *ShellParse) ExecSarNetIOCmd(cmd ShellCmd, input string) CmdData {
	netIOS := getNetIOV2(input)
	return NewCmdData(cmd.HostName, sp.Miners.MinerId, cmd.CmdType, cmd.State, netIOS)

}

func (sp *ShellParse) ExecIOTopCmd(cmd ShellCmd, input string) CmdData {
	diskRead := diskReadReg.FindAllStringSubmatch(input, 1)
	diskWrite := diskWriteReg.FindAllStringSubmatch(input, 1)
	info := IoInfo{DiskR: getRegexValue(diskRead), DiskW: getRegexValue(diskWrite)}
	return NewCmdData(cmd.HostName, sp.Miners.MinerId, cmd.CmdType, cmd.State, info)
}

func getCpuTemperV2(data string) CpuTemp {
	tdieValue := cpuTemperatureRTdieReg.FindAllStringSubmatch(data, 1)
	value := getRegexValue(tdieValue)
	if value != "0" {
		return CpuTemp{CpuTemp: value}
	}
	coreValue := cpuTemperatureCoreReg.FindAllStringSubmatch(data, 1)
	return CpuTemp{CpuTemp: getRegexValue(coreValue)}
}

func getNetIOV2(data string) map[string]interface{} {
	param := make(map[string]interface{})
	allSubStr := netIOAverageReg.FindAllStringSubmatch(data, -1)
	for i := 0; i < len(allSubStr); i++ {
		if len(allSubStr[i]) == 0 {
			continue
		}
		temp := allSubStr[i][0]
		if strings.Contains(temp, "IFACE") {
			continue
		}
		fields := strings.Fields(temp)
		if len(fields) < 9 {
			continue
		}
		netIO := NetIO{
			Name: fields[1],
			Rx:   fields[4],
			Tx:   fields[5],
		}
		param[fields[1]] = utils.StructToMapByJson(netIO)

	}
	return param
}

func (sp *ShellParse) GetMinerJobsCV2(data string) interface{} {
	canParse := false
	jobsMap := make(map[string]interface{})
	reader := bufio.NewReader(bytes.NewBuffer([]byte(data)))
	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		if !canParse && strings.HasPrefix(line, "ID") {
			canParse = true
			continue
		}
		if canParse {
			arrs := strings.Fields(line)
			if len(arrs) < 7 {
				continue
			}
			task := Task{
				Id:       arrs[0],
				Sector:   arrs[1],
				Worker:   arrs[2],
				HostName: arrs[3],
				Task:     arrs[4],
				State:    arrs[5],
				Time:     arrs[6],
			}
			jobsMap[arrs[1]] = utils.StructToMapByJson(task)
		}
	}
	return jobsMap
}

func (sp *ShellParse) GetMinerWorkersV2(input string) []*WorkerInfo {
	reader := bufio.NewReader(bytes.NewBuffer([]byte(input)))
	var res []*WorkerInfo
	preHostIndex := -1
	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		if strings.HasPrefix(line, "Worker") {
			fields := strings.Fields(line)
			if len(fields) < 4 {
				continue
			}
			preHostIndex = preHostIndex + 1
			res = append(res, &WorkerInfo{HostName: fields[3], Id: fields[1]})
		} else if strings.Contains(line, "GPU") {
			if len(res) != 0 && len(res) == preHostIndex+1 {
				workerInfo := res[preHostIndex]
				workerInfo.GPU = 1
			}
		}
	}
	return res
}

func getGraphicsCardInfoV2(data string) interface{} {

	param := make(map[string]interface{})
	idAllStrs := gpuIdReg.FindAllStringSubmatch(data, -1)
	gpInfoAllStrs := gpuInfoReg.FindAllStringSubmatch(data, -1)
	if len(idAllStrs) < 1 || len(gpInfoAllStrs) < 1 {
		return nil
	}
	for i := 0; i < len(idAllStrs); i++ {
		if len(idAllStrs[i]) < 1 || len(gpInfoAllStrs) < 1 {
			continue
		}
		temp, used := getGpuInfo(gpInfoAllStrs[i][0])
		gpuId := getGpuId(idAllStrs[i][0])
		gpu := GpuInfo{
			Name: gpuId,
			Temp: temp,
			Use:  used,
		}
		mapByJson := utils.StructToMapByJson(gpu)
		param[gpuId] = mapByJson

	}
	return param
}

func (sp *ShellParse) runHardware(w *WorkerInfo, obj chan HardwareInfo) {
	execInfo := fmt.Sprintf(`root@%v`, w.HostName)
	hardwareInfo := HardwareInfo{}
	var data string
	var err error
	if w.GPU == 0 {
		data, err = sp.ExecCmd("ssh", execInfo, "sensors", "&&", "uptime", "&&", "free -h", "&&", "df -h", "&&", "sar", "-n", "DEV", "1", "2", "&&", "iotop", "-bn1", "|", "head", "-n", "2")
		if err != nil {
			obj <- hardwareInfo
			return
		}
	} else {
		data, err = sp.ExecCmd("ssh", execInfo, "sensors", "&&", "uptime", "&&", "free -h", "&&", "df -h", "&&", "sar", "-n", "DEV", "1", "2", "&&", "iotop", "-bn1", "|", "head", "-n", "2", "&&", "nvidia-smi")
		if err != nil {
			obj <- hardwareInfo
			return
		}
	}

	hardwareInfo.HostName = w.HostName
	hardwareInfo.CpuTemper = getCpuTemper(data)
	hardwareInfo.NetIO = getNetIOV1(data)

	if w.GPU == 1 {
		hardwareInfo.GpuInfo = getGraphicsCardInfoV1(data)
	}

	cpuLoad := cpuLoadReg.FindAllStringSubmatch(data, 1)
	hardwareInfo.CpuLoad = getRegexValue(cpuLoad)

	memoryUsed := memoryUsedReg.FindAllStringSubmatch(data, 1)
	hardwareInfo.UseMemory = getRegexValueById(memoryUsed, 2)
	hardwareInfo.TotalMemory = getRegexValueById(memoryUsed, 1)

	diskUsed := diskUsedRateReg.FindAllStringSubmatch(data, 1)
	hardwareInfo.UseDisk = getRegexValue(diskUsed)
	diskRead := diskReadReg.FindAllStringSubmatch(data, 1)
	hardwareInfo.DiskR = getRegexValue(diskRead)

	diskWrite := diskWriteReg.FindAllStringSubmatch(data, 1)
	hardwareInfo.DiskW = getRegexValue(diskWrite)
	obj <- hardwareInfo
	return
}

func getCpuTemper(data string) string {
	tdieValue := cpuTemperatureRTdieReg.FindAllStringSubmatch(data, 1)
	value := getRegexValue(tdieValue)
	if value != "0" {
		return value
	}
	coreValue := cpuTemperatureCoreReg.FindAllStringSubmatch(data, 1)
	return getRegexValue(coreValue)
}

func getGraphicsCardInfo(data string) interface{} {
	var graphCardList []GraphicsCardInfo
	idAllStrs := gpuIdReg.FindAllStringSubmatch(data, -1)
	gpInfoAllStrs := gpuInfoReg.FindAllStringSubmatch(data, -1)
	if len(idAllStrs) < 1 || len(gpInfoAllStrs) < 1 {
		return graphCardList
	}
	for i := 0; i < len(idAllStrs); i++ {
		if len(idAllStrs[i]) < 1 || len(gpInfoAllStrs) < 1 {
			continue
		}
		temp, used := getGpuInfo(gpInfoAllStrs[i][0])
		graphCardList = append(graphCardList, GraphicsCardInfo{
			Name: getGpuId(idAllStrs[i][0]),
			Temp: temp,
			Use:  used,
		})
	}
	return graphCardList
}

func getGraphicsCardInfoV1(data string) interface{} {
	graphCardList := make(map[string]interface{})
	idAllStrs := gpuIdReg.FindAllStringSubmatch(data, -1)
	gpInfoAllStrs := gpuInfoReg.FindAllStringSubmatch(data, -1)
	if len(idAllStrs) < 1 || len(gpInfoAllStrs) < 1 {
		return graphCardList
	}
	for i := 0; i < len(idAllStrs); i++ {
		if len(idAllStrs[i]) < 1 || len(gpInfoAllStrs) < 1 {
			continue
		}
		temp, used := getGpuInfo(gpInfoAllStrs[i][0])
		graphCardList[getGpuId(idAllStrs[i][0])] = GraphicsCardInfo{
			Name: getGpuId(idAllStrs[i][0]),
			Temp: temp,
			Use:  used,
		}
	}
	return graphCardList
}

func getGpuInfo(src string) (string, string) {
	fields := strings.Fields(src)
	if len(fields) < 15 {
		return "-", "-"
	}
	return fields[2], fields[12]
}

func getGpuId(src string) string {
	fields := strings.Fields(src)
	if len(fields) < 12 {
		return "-"
	}
	return fields[1]
}

func getNetIOV1(data string) interface{} {
	allSubStr := netIOAverageReg.FindAllStringSubmatch(data, -1)
	NetIOes := make(map[string]interface{})
	for i := 0; i < len(allSubStr); i++ {
		if len(allSubStr[i]) == 0 {
			continue
		}
		temp := allSubStr[i][0]
		if strings.Contains(temp, "IFACE") {
			continue
		}
		fields := strings.Fields(temp)
		if len(fields) < 9 {
			continue
		}
		NetIOes[fields[1]] = NetCardIO{
			Name: fields[1],
			Rx:   fields[4],
			TX:   fields[5],
		}
	}
	return NetIOes
}

func getNetIO(data string) interface{} {
	allSubStr := netIOAverageReg.FindAllStringSubmatch(data, -1)
	var NetIOes []NetCardIO
	for i := 0; i < len(allSubStr); i++ {
		if len(allSubStr[i]) == 0 {
			continue
		}
		temp := allSubStr[i][0]
		if strings.Contains(temp, "IFACE") {
			continue
		}
		fields := strings.Fields(temp)
		if len(fields) < 9 {
			continue
		}
		NetIOes = append(NetIOes, NetCardIO{
			Name: fields[1],
			Rx:   fields[4],
			TX:   fields[5],
		})
	}
	return NetIOes
}

func (sp *ShellParse) ExecCmd(cmdName string, args ...string) (string, error) {
	log.Debug("exec cmd: ", cmdName, args)
	cmd := exec.Command(cmdName, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
