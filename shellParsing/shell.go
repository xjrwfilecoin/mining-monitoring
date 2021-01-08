package shellParsing

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"mining-monitoring/log"
	"os/exec"
	"strings"
	"time"
)

type ShellParse struct {
	Workers []WorkerInfo
}

func NewShellParse(workers []WorkerInfo) *ShellParse {
	return &ShellParse{
		Workers: workers,
	}
}

// todo 隔离分开
func (sp *ShellParse) getTaskInfo() (map[string]interface{}, error) {
	minerInfo, err := sp.GetMinerInfo()
	if err != nil {
		return nil, err
	}
	minerInfoMap := structToMapByJson(minerInfo)

	log.Debug("minerInfo: ", *minerInfo)

	postBalance, err := sp.GetPostBalance()
	if err != nil {
		return nil, err
	}
	minerInfoMap["postBalance"] = postBalance
	log.Debug("PostBalance: ", postBalance)

	msgNums, err := sp.MsgNums()
	if err != nil {
		return nil, err
	}
	minerInfoMap["messageNums"] = msgNums
	log.Debug("msgNums: ", msgNums)

	minerJobs, err := sp.GetMinerJobs()
	if err != nil {
		return nil, err
	}
	log.Debug("minerJobs: ", minerJobs)
	hardwareInfo, err := sp.hardwareInfo(sp.Workers)
	if err != nil {
		return nil, err
	}
	log.Debug("hardwareInfo: ", hardwareInfo)

	workerInfo := mergeWorkerInfo(minerJobs, hardwareInfo)
	log.Debug("workerInfo: ", workerInfo)

	minerInfoMap["workerInfo"] = workerInfo
	return minerInfoMap, nil
}

func (sp *ShellParse) MsgNums() (interface{}, error) {
	data, err := sp.ExecCmd("lotus", `mpool`, "pending", )
	if err != nil {
		return "", fmt.Errorf("exec mpool pending: %v \n", err)
	}
	count := strings.Count(data, "Message")
	return count, nil
}

// 获取所有worker硬件信息
func (sp *ShellParse) hardwareInfo(workers []WorkerInfo) ([]HardwareInfo, error) {
	if len(workers) == 0 {
		return nil, nil
	}
	obj := make(chan HardwareInfo, 10)
	for i := 0; i < len(workers); i++ {
		wInfo := workers[i]
		go sp.runHardware(wInfo, obj)
	}
	ctx, _ := context.WithTimeout(context.TODO(), 60*time.Second)

	total := 0
	var resInfo []HardwareInfo
	for {
		select {
		case res := <-obj:
			if res.IsValid() {
				resInfo = append(resInfo, res)
			}
			total = total + 1
			if total == len(workers) {
				return resInfo, nil
			}
		case <-ctx.Done():
			return resInfo, nil
		}
	}
}

func (sp *ShellParse) runHardware(w WorkerInfo, obj chan HardwareInfo) {
	execInfo := fmt.Sprintf(`root@%v`, w.IP)
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
	hardwareInfo.NetIO = getNetIO(data)

	if w.GPU == 1 {
		hardwareInfo.GpuInfo = getGraphicsCardInfo(data)
	}

	cpuLoad := cpuLoadReg.FindAllStringSubmatch(data, 1)
	hardwareInfo.CpuLoad = getRegexValue(cpuLoad)

	memoryUsed := memoryUsedReg.FindAllStringSubmatch(data, 1)
	hardwareInfo.UseMemory = getRegexValue(memoryUsed)
	hardwareInfo.TotalMemory = getRegexValueById(memoryUsed, 2)

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

func (sp *ShellParse) GetMinerJobs() ([]Task, error) {
	data, err := sp.ExecCmd("lotus-miner", "sealing", "jobs")
	if err != nil {
		return nil, fmt.Errorf("exec lotus-miner sealing jobs: %v \n", err)
	}
	canParse := false
	var taskList []Task
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
			if task, ok := parseTaskByStr(line); ok {
				taskList = append(taskList, task)
			}

		}
	}
	return taskList, nil
}

func parseTaskByStr(line string) (Task, bool) {
	arrs := strings.Fields(line)
	if len(arrs) < 7 {
		return Task{}, false
	}
	return Task{
		Id:       arrs[0],
		Sector:   arrs[1],
		Worker:   arrs[2],
		HostName: arrs[3],
		Task:     arrs[4],
		State:    arrs[5],
		Time:     arrs[6],
	}, true
}

func (sp *ShellParse) ExecCmd(cmdName string, args ...string) (string, error) {
	cmd := exec.CommandContext(context.TODO(), cmdName, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (sp *ShellParse) GetPostBalance() (string, error) {
	data, err := sp.ExecCmd("lotus-miner", "actor", "control", "list")
	if err != nil {
		return "", fmt.Errorf("exec lotus-miner actor control list: %v \n", err)
	}
	postBalance := postBalanceReg.FindAllStringSubmatch(data, 1)
	pb := getRegexValue(postBalance)
	return pb, nil
}

func (sp *ShellParse) GetMinerInfo() (*MinerInfo, error) {
	data, err := sp.ExecCmd("lotus-miner", "info")
	if err != nil {
		return nil, fmt.Errorf("exec lotus-miner info  %v \n", err)
	}
	minerInfo := &MinerInfo{}
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
	effectSectors := effectSectorReg.FindAllStringSubmatch(data, 1)
	minerInfo.EffectiveSectors = getRegexValue(effectSectors)
	errorsSectors := errorSectorReg.FindAllStringSubmatch(data, 1)
	minerInfo.ErrorSectors = getRegexValue(errorsSectors)
	recoverySectors := recoverySectorReg.FindAllStringSubmatch(data, 1)
	minerInfo.RecoverySectors = getRegexValue(recoverySectors)
	deletedSectors := deletedSectorReg.FindAllStringSubmatch(data, 1)
	minerInfo.DeletedSectors = getRegexValue(deletedSectors)
	failSectors := failSectorReg.FindAllStringSubmatch(data, 1)
	minerInfo.FailSectors = getRegexValue(failSectors)
	minerInfo.Timestamp = time.Now().Unix()
	return minerInfo, nil
}
