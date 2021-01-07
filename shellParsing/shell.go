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

func (sp *ShellParse) getTaskInfo() (map[string]interface{}, error) {
	minerInfo, err := sp.GetMinerInfo()
	if err != nil {
		return nil, err
	}
	minerInfoMap := structToMapByJson(minerInfo)

	log.Debug("minerInfo: %v \n", minerInfo)

	postBalance, err := sp.GetPostBalance()
	if err != nil {
		return nil, err
	}
	minerInfoMap["postBalance"] = postBalance
	log.Debug("PostBalance: %v \n", postBalance)

	msgNums, err := sp.MsgNums()
	if err != nil {
		return nil, err
	}
	minerInfoMap["messageNums"] = msgNums
	log.Debug("msgNums: %v \n", msgNums)

	minerJobs, err := sp.GetMinerJobs()
	if err != nil {
		return nil, err
	}
	log.Debug("minerJobs: %v \n", minerJobs)
	hardwareInfo, err := sp.hardwareInfo(sp.Workers)
	if err != nil {
		return nil, err
	}
	log.Debug("hardwareInfo: %v \n", hardwareInfo)

	workerInfo := mergeWorkerInfo(minerJobs, hardwareInfo)
	log.Debug("workerInfo: %v \n", workerInfo)

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
	//execInfo := fmt.Sprintf(`root@%v "sensors&&uptime&&free -h&&df -h&&sar -n DEV 1 2&& iotop -bn1|head -n 2"`, w.IP)
	execInfo := fmt.Sprintf(`root@%v`, w.IP)
	data, err := sp.ExecCmd("ssh", execInfo,"sensors","&&","uptime","&&","free -h","&&","df -h","&&","sar","-n","DEV","1","2","&&","iotop","-bn1","|","head","-n","2")
	hardwareInfo := HardwareInfo{}
	if err != nil {
		obj <- hardwareInfo
		return
	}
	hardwareInfo.HostName = w.HostName
	cpuTemperature := cpuTemperatureReg.FindAllStringSubmatch(data, 1)
	hardwareInfo.CpuTemper = getRegexValue(cpuTemperature)

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
			if task, ok := getHardwareInfo(line); ok {
				taskList = append(taskList, task)
			}

		}
	}
	return taskList, nil
}

func getHardwareInfo(line string) (Task, bool) {
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
	src := string(data)
	minerInfo := &MinerInfo{}
	minerId := minerIdReg.FindString(src)
	minerInfo.MinerId = minerId
	minerBalance := minerBalanceReg.FindAllStringSubmatch(src, 1)
	minerInfo.MinerBalance = getRegexValue(minerBalance)
	workerBalance := workerBalanceReg.FindAllStringSubmatch(src, 1)
	minerInfo.WorkerBalance = getRegexValue(workerBalance)
	pledgeBalance := pledgeBalanceReg.FindAllStringSubmatch(src, 1)
	minerInfo.PledgeBalance = getRegexValue(pledgeBalance)
	totalPower := totalPowerReg.FindAllStringSubmatch(src, 1)
	minerInfo.EffectivePower = getRegexValue(totalPower)
	effectPower := effectPowerReg.FindAllStringSubmatch(src, 1)
	minerInfo.EffectivePower = getRegexValue(effectPower)
	totalSectors := totalSectorsReg.FindAllStringSubmatch(src, 1)
	minerInfo.TotalSectors = getRegexValue(totalSectors)
	effectSectors := effectSectorReg.FindAllStringSubmatch(src, 1)
	minerInfo.EffectiveSectors = getRegexValue(effectSectors)
	errorsSectors := errorSectorReg.FindAllStringSubmatch(src, 1)
	minerInfo.ErrorSectors = getRegexValue(errorsSectors)
	recoverySectors := recoverySectorReg.FindAllStringSubmatch(src, 1)
	minerInfo.RecoverySectors = getRegexValue(recoverySectors)
	deletedSectors := deletedSectorReg.FindAllStringSubmatch(src, 1)
	minerInfo.DeletedSectors = getRegexValue(deletedSectors)
	failSectors := failSectorReg.FindAllStringSubmatch(src, 1)
	minerInfo.FailSectors = getRegexValue(failSectors)
	return minerInfo, nil
}
