package shell

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mining-monitoring/log"
	"mining-monitoring/utils"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Parse struct {
	Workers      []*WorkerInfo
	Miners       Miner
	HostName     string
	cmdSign      chan CmdData
	CmdParseMap  map[CmdType]func(cmd ShellCmd, input string) CmdData
	cmdHeartTime time.Duration // 秒
	closing      chan struct{}
	CmdMap       map[CmdType]ShellCmd
	workerSign   chan Worker
	sync.Mutex
	cmdCount int64 //任务总数控制
}

func NewShellParse() *Parse {
	return &Parse{
		cmdSign:      make(chan CmdData, 150),
		CmdParseMap:  make(map[CmdType]func(cmd ShellCmd, input string) CmdData),
		closing:      make(chan struct{}),
		CmdMap:       make(map[CmdType]ShellCmd),
		workerSign:   make(chan Worker, 1),
		cmdHeartTime: 10,
	}
}

func (sp *Parse) Close() {
	close(sp.closing)
	close(sp.cmdSign)
}

func (sp *Parse) initCmdParse() {
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

func (sp *Parse) getMinerCmdList(minerId string) []ShellCmd {
	var cmdList []ShellCmd
	cmdList = append(cmdList, NewLotusShellCmd(minerId, "lotus-miner", LotusControlList, []string{"actor", "control", "list"}))
	cmdList = append(cmdList, NewLotusShellCmd(minerId, "lotus-miner", LotusMinerJobs, []string{"sealing", "jobs"}))
	return cmdList
}

func (sp *Parse) GetMinerHardwareCmdList(hostName string) []ShellCmd {
	var cmdList []ShellCmd
	cmdList = append(cmdList, NewHardwareShellCmd(hostName, "sensors", SensorsCmd, []string{}))
	cmdList = append(cmdList, NewHardwareShellCmd(hostName, "uptime", UpTimeCmd, []string{}))
	cmdList = append(cmdList, NewHardwareShellCmd(hostName, "free", FreeHCmd, []string{"-h"}))
	cmdList = append(cmdList, NewHardwareShellCmd(hostName, "df", DfHCMd, []string{"-h"}))
	cmdList = append(cmdList, NewHardwareShellCmd(hostName, "sar", SarCmd, []string{"-n", "DEV", "1", "2"}))
	cmdList = append(cmdList, NewHardwareShellCmd(hostName, "iotop", IOCmd, []string{"-bn1", "|", "head", "-n", "2"}))
	cmdList = append(cmdList, NewHardwareShellCmd(hostName, "nvidia-smi", GpuCmd, []string{}))
	return cmdList
}

// sshpass -p 1 ssh root@xjrw_node01 "free -h"
func (sp *Parse) getWorkCmdList(hostName string, gpuEnable bool) []ShellCmd {
	execInfo := fmt.Sprintf(`root@%v`, hostName)
	var cmdList []ShellCmd
	cmdList = append(cmdList, NewHardwareShellCmd(hostName, "sshpass", SensorsCmd, []string{"-p", "", "ssh", "-p", SSHPORT, execInfo, "sensors"}))
	cmdList = append(cmdList, NewHardwareShellCmd(hostName, "sshpass", UpTimeCmd, []string{"-p", "", "ssh", "-p", SSHPORT, execInfo, "uptime"}))
	cmdList = append(cmdList, NewHardwareShellCmd(hostName, "sshpass", FreeHCmd, []string{"-p", "", "ssh", "-p", SSHPORT, execInfo, "free", "-h"}))
	cmdList = append(cmdList, NewHardwareShellCmd(hostName, "sshpass", DfHCMd, []string{"-p", "", "ssh", "-p", SSHPORT, execInfo, "df", "-h"}))
	cmdList = append(cmdList, NewHardwareShellCmd(hostName, "sshpass", SarCmd, []string{"-p", "", "ssh", "-p", SSHPORT, execInfo, "sar", "-n", "DEV", "1", "2"}))
	cmdList = append(cmdList, NewHardwareShellCmd(hostName, "sshpass", IOCmd, []string{"-p", "", "ssh", "-p", SSHPORT, execInfo, "iotop", "-bn1", "|", "head", "-n", "2"}))
	if gpuEnable {
		cmdList = append(cmdList, NewHardwareShellCmd(hostName, "sshpass", GpuCmd, []string{"-p", "", "ssh", "-p", SSHPORT, execInfo, "nvidia-smi"}))
	}
	return cmdList
}

func (sp *Parse) doMinerInfo() {
	ticker := time.NewTicker(120 * time.Second)
	cmd := NewLotusShellCmd(sp.Miners.MinerId, "lotus-miner", LotusMinerInfoCmd, []string{"info"})
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if fn, ok := sp.CmdParseMap[cmd.CmdType]; ok {
				if !sp.CanAddTask() {
					continue
				}
				ctx, _ := context.WithTimeout(context.Background(), 40*time.Second)
				sp.processTask(ctx, cmd, sp.cmdSign, fn)
			}
		case <-sp.closing:
			return

		}
	}
}

func (sp *Parse) doWorkers() {
	if err := recover(); err != nil {
		log.Error(err)
	}
	ticker := time.NewTicker(120 * time.Second)
	defer ticker.Stop()
	minerWorkersCmd := NewLotusShellCmd("", "lotus-miner", LotusMinerWorkers, []string{"sealing", "workers"})
	for {
		select {
		case <-ticker.C:
			err := sp.getWorkerList(minerWorkersCmd)
			if err != nil {
				log.Error("get worker info error: ", err.Error())
			}
		case <-sp.closing:
			return

		}
	}
}

func (sp *Parse) PingWorkersV1() {
	var res []map[string]interface{}
	for i := 0; i < len(sp.Workers); i++ {
		worker := sp.Workers[i]
		execInfo := fmt.Sprintf(`root@%v`, worker.HostName)
		cmd := NewHardwareShellCmd(worker.HostName, "sshpass", FreeHCmd, []string{"-p", "", "ssh", "-p", SSHPORT, execInfo})
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		err := sp.execShellCmd(ctx, cmd, func(input string) {

		})
		if err != nil {
			if strings.Contains(err.Error(), "exit") {
				worker.NetState = NetDisabled
			} else {
				worker.NetState = NetNormal
			}
			log.Error("ping workers: ", err.Error())
		}
		res = append(res, utils.StructToMapByJson(worker))
	}
	sp.cmdSign <- NewCmdData(sp.Miners.MinerId, sp.Miners.MinerId, LotusMinerWorkers, LotusState, res)
}

func (sp *Parse) getWorkerList(cmd ShellCmd) error {
	ctx, _ := context.WithTimeout(context.Background(), 6*time.Second)
	err := sp.execShellCmd(ctx, cmd, func(input string) {
		workers := sp.GetMinerWorkersV2(input)
		sp.Lock()
		sp.Workers = workers
		sp.Unlock()
	})
	if err != nil {
		return err
	}
	sp.PingWorkersV1()
	return nil
}

func (sp *Parse) doHardwareInfoV1() {
	for {
		select {
		case <-sp.closing:
			return
		default:
		}
		sp.Lock()
		workerList := sp.Workers
		sp.Unlock()
		for i := 0; i < len(workerList); i++ {
			worker := workerList[i]
			if i != 0 && i%3 == 0 {
				time.Sleep(300 * time.Millisecond)
			}
			sp.runWorkerCmdList(worker, sp.cmdSign)
		}
		time.Sleep(4 * time.Second)

	}
}

func (sp *Parse) doHardWareInfo() {
	ticker := time.NewTicker(6 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			sp.Lock()
			workerList := sp.Workers
			sp.Unlock()
			for i := 0; i < len(workerList); i++ {
				worker := workerList[i]
				sp.runWorkerCmdList(worker, sp.cmdSign)
			}
		case <-sp.closing:
			return

		}
	}
}

func (sp *Parse) runWorkerCmdList(worker *WorkerInfo, sing chan CmdData) {
	if worker == nil {
		return
	}
	cmdList := sp.getWorkCmdList(worker.HostName, worker.GPU != 0)
	for i := 0; i < len(cmdList); i++ {
		cmd := cmdList[i]
		if fn, ok := sp.CmdParseMap[cmd.CmdType]; ok {
			if !sp.CanAddTask() {
				continue
			}
			ctx, _ := context.WithTimeout(context.Background(), 6*time.Second)
			go sp.processTask(ctx,cmd, sing, fn)
		}

	}
}

func (sp *Parse) needGPU(worker WorkerInfo01) bool {
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

func (sp *Parse) getMiner() error {
	minerInfoCmd := NewLotusShellCmd("", "lotus-miner", LotusMinerInfoCmd, []string{"info"})
	ctx, _ := context.WithTimeout(context.Background(), 40*time.Second)
	err := sp.execShellCmd(ctx, minerInfoCmd, func(input string) {
		minerInfo := sp.getMinerInfo(input)
		sp.cmdSign <- NewCmdData(minerInfo.MinerId, minerInfo.MinerId, LotusMinerInfoCmd, LotusState, utils.StructToMapByJson(minerInfo))
		sp.Miners = Miner{MinerId: minerInfo.MinerId,}
	})
	return err
}

func (sp *Parse) Send() {
	err := sp.Init()
	if err != nil {
		panic(fmt.Errorf("get miner info error: %v ", err))
	}
	go sp.doWorkers()
	go sp.miningInfo()
	go sp.doMinerInfo()
	//go sp.doHardWareInfo()
	go sp.doHardwareInfoV1()
	//go sp.getMinerHardwareInfo()

}

// miner ssh本地，如果不做，自己获取
func (sp *Parse) getMinerHardwareInfo() {
	ticker := time.NewTicker(5 * time.Second)
	cmdList := sp.GetMinerHardwareCmdList(sp.HostName)
	//defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for i := 0; i < len(cmdList); i++ {
				cmd := cmdList[i]
				if fn, ok := sp.CmdParseMap[cmd.CmdType]; ok {
					if !sp.CanAddTask() {
						continue
					}
					ctx, _ := context.WithTimeout(context.Background(), 6*time.Second)
					go sp.processTask(ctx,cmd, sp.cmdSign, fn)
				}
			}
		case <-sp.closing:
			return

		}
	}
}

func (sp *Parse) Init() error {
	sp.initCmdParse()
	err := sp.getMiner()
	if err != nil {
		return err
	}
	minerWorkersCmd := NewLotusShellCmd("", "lotus-miner", LotusMinerWorkers, []string{"sealing", "workers"})
	err = sp.getWorkerList(minerWorkersCmd)
	if err != nil {
		return err
	}
	err = sp.InitMinerInfo()
	if err != nil {
		return err
	}

	return nil
}

func (sp *Parse) InitMinerInfo() error {
	var cmdList []ShellCmd
	minerId := sp.Miners.MinerId
	cmdList = append(cmdList, NewLotusShellCmd(minerId, "lotus-miner", LotusControlList, []string{"actor", "control", "list"}))
	cmdList = append(cmdList, NewLotusShellCmd(minerId, "lotus-miner", LotusMinerJobs, []string{"sealing", "jobs"}))
	for i := 0; i < len(cmdList); i++ {
		cmd := cmdList[i]
		if fn, ok := sp.CmdParseMap[cmd.CmdType]; ok {
			ctx, _ := context.WithTimeout(context.Background(), 6*time.Second)
			err := sp.processTask(ctx,cmd, sp.cmdSign, fn)
			if err != nil {
				return fmt.Errorf("get miner info : %v ", err)
			}
		}
	}
	return nil
}

func (sp *Parse) miningInfo() {
	ticker := time.NewTicker(180 * time.Second)
	cmdList := sp.getMinerCmdList(sp.Miners.MinerId)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for i := 0; i < len(cmdList); i++ {
				cmd := cmdList[i]
				if fn, ok := sp.CmdParseMap[cmd.CmdType]; ok {
					if !sp.CanAddTask() {
						continue
					}
					ctx, _ := context.WithTimeout(context.Background(), 6*time.Second)
					go sp.processTask(ctx,cmd, sp.cmdSign, fn)
				}
			}

		case <-sp.closing:
			return
		}
	}
}

func (sp *Parse) Receiver(recv chan CmdData) {
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

		}
	}
}

func (sp *Parse) execShellCmd(ctx context.Context, cmd ShellCmd, fn func(input string)) error {
	data, err := sp.ExecCmd(ctx, cmd.Name, cmd.Params...)
	if err != nil {
		return err
	}
	fn(data)
	return nil
}

func (sp *Parse) WrapProcessTask(ctx context.Context, shellCmd <-chan ShellCmd, fn func(cmd ShellCmd, input string) CmdData) error {
	for {
		select {
		case cmd := <-shellCmd:
			err := sp.processTask(ctx, cmd, sp.cmdSign, fn)
			if err != nil {
				log.Error("wrapProcessTask: ", err.Error())
			}
		case <-ctx.Done():
			return nil

		}
	}
}

func (sp *Parse) CanAddTask() bool {
	return atomic.LoadInt64(&sp.cmdCount) < 350
}

func (sp *Parse) AddTaskCount() {
	atomic.AddInt64(&sp.cmdCount, 1)
}

func (sp *Parse) DelTaskCount() {
	atomic.AddInt64(&sp.cmdCount, -1)
}

func (sp *Parse) processTask(ctx context.Context, cmd ShellCmd, sign chan CmdData, fn func(cmd ShellCmd, input string) CmdData) error {
	sp.AddTaskCount()
	defer sp.DelTaskCount()

	output, err := sp.ExecCmd(ctx, cmd.Name, cmd.Params...)
	if err != nil {
		log.Error("process task error ", cmd.CmdType, cmd.Name, cmd.HostName, err)
		return err
	}
	cmdData := fn(cmd, output)
	sign <- cmdData
	return nil

}

func (sp *Parse) ExecLotusWorkers(cmd ShellCmd, data string) CmdData {
	workerInfos := sp.GetMinerWorkersV2(data)
	return NewCmdData(cmd.HostName, sp.Miners.MinerId, cmd.CmdType, cmd.State, workerInfos)
}

func (sp *Parse) ExecLotusMinerJobs(cmd ShellCmd, data string) CmdData {
	tasks := sp.GetMinerJobsV3(data)
	return NewCmdData(cmd.HostName, sp.Miners.MinerId, cmd.CmdType, cmd.State, tasks)
}

func (sp *Parse) ExecLotusPostInfo(cmd ShellCmd, data string) CmdData {
	postBalance := postBalanceTestReg.FindAllStringSubmatch(data, 1)
	postValue := getRegexValue(postBalance)
	balance := PostBalance{PostBalance: postValue}
	return NewCmdData(cmd.HostName, sp.Miners.MinerId, cmd.CmdType, cmd.State, utils.StructToMapByJson(balance))
}

func (sp *Parse) ExecLotusMpoolInfo(cmd ShellCmd, data string) CmdData {
	count := strings.Count(data, "Message")
	return NewCmdData(cmd.HostName, sp.Miners.MinerId, cmd.CmdType, cmd.State, count)
}

func (sp *Parse) ExecLotusMinerInfo(cmd ShellCmd, data string) CmdData {
	minerInfo := sp.getMinerInfo(data)
	return NewCmdData(cmd.HostName, sp.Miners.MinerId, cmd.CmdType, cmd.State, minerInfo)
}

func (sp *Parse) getMinerInfo(data string) MinerInfo {
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

	expectBlock := expectBlockReg.FindAllStringSubmatch(data, 1)
	minerInfo.ExpectBlock = getRegexValue(expectBlock)

	commitWait := commitWaitReg.FindAllStringSubmatch(data, 1)
	minerInfo.CommitWait = getRegexValue(commitWait)

	PreCommitWait := preCommitWaitReg.FindAllStringSubmatch(data, 1)
	minerInfo.PreCommitWait = getRegexValue(PreCommitWait)

	available := availableReg.FindAllStringSubmatch(data, 1)
	minerInfo.MinerAvailable = getRegexValue(available)

	PreCommit1 := PreCommit1Reg.FindAllStringSubmatch(data, 1)
	minerInfo.PreCommit1 = getRegexValue(PreCommit1)
	PreCommit2 := PreCommit2Reg.FindAllStringSubmatch(data, 1)
	minerInfo.PreCommit2 = getRegexValue(PreCommit2)
	WaitSeed := WaitSeedReg.FindAllStringSubmatch(data, 1)
	minerInfo.WaitSeed = getRegexValue(WaitSeed)
	Committing := CommittingReg.FindAllStringSubmatch(data, 1)
	minerInfo.Committing = getRegexValue(Committing)

	FinalizeSector := FinalizeSectorReg.FindAllStringSubmatch(data, 1)
	minerInfo.FinalizeSector = getRegexValue(FinalizeSector)

	return minerInfo
}

func (sp *Parse) ExecGPUCmd(cmd ShellCmd, input string) CmdData {
	gpuInfos := getGraphicsCardInfoV3(input)
	return NewCmdData(cmd.HostName, sp.Miners.MinerId, cmd.CmdType, cmd.State, gpuInfos)
}

func (sp *Parse) ExecSensorsCmd(cmd ShellCmd, output string) CmdData {
	cpuTemp := getCpuTemperV2(output)
	return NewCmdData(cmd.HostName, sp.Miners.MinerId, cmd.CmdType, cmd.State, utils.StructToMapByJson(cpuTemp))

}

func (sp *Parse) ExecDfHCmd(cmd ShellCmd, input string) CmdData {
	diskUsed := diskUsedRateReg.FindAllStringSubmatch(input, 1)
	diskInfo := Disk{UseDisk: getRegexValue(diskUsed)}
	return NewCmdData(cmd.HostName, sp.Miners.MinerId, cmd.CmdType, cmd.State, utils.StructToMapByJson(diskInfo))

}

func (sp *Parse) ExecFreeHCmd(cmd ShellCmd, input string) CmdData {

	memoryUsed := memoryUsedReg.FindAllStringSubmatch(input, 1)
	memory := Memory{UseMemory: getRegexValueById(memoryUsed, 2), TotalMemory: getRegexValueById(memoryUsed, 1)}
	return NewCmdData(cmd.HostName, sp.Miners.MinerId, cmd.CmdType, cmd.State, utils.StructToMapByJson(memory))

}

func (sp *Parse) ExecUptimeCmd(cmd ShellCmd, input string) CmdData {
	cpuLoad := cpuLoadReg.FindAllStringSubmatch(input, 1)
	load := CpuLoad{CpuLoad: getRegexValue(cpuLoad)}
	return NewCmdData(cmd.HostName, sp.Miners.MinerId, cmd.CmdType, cmd.State, utils.StructToMapByJson(load))
}

func (sp *Parse) ExecSarNetIOCmd(cmd ShellCmd, input string) CmdData {
	netIOS := getNetIOV3(input)
	return NewCmdData(cmd.HostName, sp.Miners.MinerId, cmd.CmdType, cmd.State, netIOS)

}

func (sp *Parse) ExecIOTopCmd(cmd ShellCmd, input string) CmdData {
	diskRead := diskReadReg.FindAllStringSubmatch(input, 1)
	diskWrite := diskWriteReg.FindAllStringSubmatch(input, 1)
	info := IoInfo{DiskR: getRegexValue(diskRead), DiskW: getRegexValue(diskWrite)}
	return NewCmdData(cmd.HostName, sp.Miners.MinerId, cmd.CmdType, cmd.State, utils.StructToMapByJson(info))
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

func getNetIOV3(data string) interface{} {
	var res []map[string]interface{}
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
		if strings.HasPrefix(fields[1], "lo") || strings.HasPrefix(fields[1], "bond") {
			continue
		}
		netIO := NetIO{
			Name:  fields[1],
			Rxpck: fields[2],
			Txpck: fields[3],
			Rx:    fields[4],
			Tx:    fields[5],
		}
		res = append(res, utils.StructToMapByJson(netIO))
	}
	return res
}

func (sp *Parse) GetMinerJobsV3(data string) interface{} {
	canParse := false
	var jobsMap []map[string]interface{}
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
			jobsMap = append(jobsMap, utils.StructToMapByJson(task))
		}
	}
	log.Error("checkJobs: shell ", jobsMap)
	return jobsMap
}

func (sp *Parse) GetMinerWorkersV2(input string) []*WorkerInfo {
	reader := bufio.NewReader(bytes.NewBuffer([]byte(input)))
	var res []*WorkerInfo
	preHostIndex := -1
	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		if strings.HasPrefix(line, "Worker") {
			taskState := Normal
			if strings.Contains(line, "disabled") {
				taskState = TaskDisabled
			}
			line = strings.ReplaceAll(line, "(disabled)", "")
			fields := strings.Fields(line)
			if len(fields) < 6 {
				continue
			}
			hostType := strings.Split(fields[5], "|")
			preHostIndex = preHostIndex + 1
			res = append(res, &WorkerInfo{HostName: fields[3], TaskState: taskState, TaskType: hostType})
		} else if strings.Contains(line, "GPU") {
			if len(res) != 0 && len(res) == preHostIndex+1 {
				workerInfo := res[preHostIndex]
				workerInfo.GPU = 1
			}
		}
	}
	return res
}

func getGraphicsCardInfoV3(data string) interface{} {
	var res []map[string]interface{}
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
		res = append(res, mapByJson)

	}
	return res
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

func (sp *Parse) ExecCmd(ctx context.Context, cmdName string, args ...string) (string, error) {
	log.Debug("exec cmd: ", cmdName, args)
	cmd := exec.CommandContext(ctx, cmdName, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
