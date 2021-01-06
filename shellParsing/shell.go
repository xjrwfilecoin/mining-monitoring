package shellParsing

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"reflect"
	"strings"
	"time"
)

type ShellParse struct {
	Workers []WorkerInfo
}

func NewShellParse() *ShellParse {
	return &ShellParse{}
}

func (sp *ShellParse) getTaskInfo() (map[string]interface{}, error) {
	minerInfo, err := sp.GetMinerInfo()
	if err != nil {
		return nil, err
	}
	postBalance, err := sp.GetPostBalance()
	if err != nil {
		return nil, err
	}
	msgNums, err := sp.MsgNums()
	if err != nil {
		return nil, err
	}
	minerJobs, err := sp.GetMinerJobs()
	if err != nil {
		return nil, err
	}

	hardwareInfo, err := sp.BatchHardwareInfo()
	if err != nil {
		return nil, err
	}



	return result, nil
}

// 将硬件信息和任务信息合并起来
// hardwareList 硬件信息列表
func mergeWorkerInfo(src []map[string]interface{}, hardwareList map[string]interface{}) map[string]interface{} {
	for i := 0; i < len(src); i++ {
		workerInfo := src[i]
		workerInfo["hardwareInfo"] = hardwareList[workerInfo["hostname"]]
	}
}

func (sp *ShellParse) MsgNums() (string, error) {
	cmd := exec.CommandContext(context.TODO(), "lotus", `mpool pending | grep -a "Version" |wc -l`)
	data, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("exec lotus-miner sealing jobs: %v \n", err)
	}
	return string(data), nil
}

func (sp *ShellParse) BatchHardwareInfo() (map[string]interface{}, error) {
	cmd := exec.CommandContext(context.TODO(), "bash", "sensors&&uptime&&free -h&&df -h&&sar -n DEV 1 2&& iotop -bn1|head -n 2")
	data, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("exec lotus-miner sealing jobs: %v \n", err)
	}
	// todo
	param := make(map[string]interface{})
	hardwareInfo := string(data)
	cpuTemperature := cpuTemperatureReg.FindAllStringSubmatch(hardwareInfo, 1)
	param["cpuTemperature"] = cpuTemperature[0][1]
	cpuLoad := cpuLoadReg.FindAllStringSubmatch(hardwareInfo, 1)
	param["cpuLoad"] = cpuLoad[0][1]
	gpuLoad := gpuLoadReg.FindAllStringSubmatch(hardwareInfo, 1)
	param["gpuLoad"] = gpuLoad[0][1]
	memoryUsed := memoryUsedReg.FindAllStringSubmatch(hardwareInfo, 1)
	param["memoryUsed"] = memoryUsed[0][1]
	memoryTotal := memoryTotalReg.FindAllStringSubmatch(hardwareInfo, 1)
	param["memoryTotal"] = memoryTotal[0][1]
	diskUsed := diskUsedRateReg.FindAllStringSubmatch(hardwareInfo, 1)
	param["diskUsed"] = diskUsed[0][1]
	return param, nil
}

// 获取所有worker硬件信息
func (sp *ShellParse) hardwareInfo(workers []WorkerInfo) ([]map[string]interface{}, error) {
	if len(workers) == 0 {
		return nil, nil
	}
	obj := make(chan map[string]interface{}, 10)
	for i := 0; i < len(workers); i++ {
		wInfo := workers[i]
		go sp.runHardware(wInfo, obj)
	}
	ctx, _ := context.WithTimeout(context.TODO(), 60*time.Second)

	total := 0
	var resInfo []map[string]interface{}
	for {
		select {
		case res := <-obj:
			resInfo = append(resInfo, res)
			total = total + 1
			if total == len(workers) {
				return resInfo, nil
			}
		case <-ctx.Done():
			return resInfo, nil
		}
	}
}

func (sp *ShellParse) runHardware(w WorkerInfo, obj chan map[string]interface{}) {

	// todo
}


func (sp *ShellParse) GetMinerJobs() ([]Task, error) {
	cmd := exec.CommandContext(context.TODO(), "lotus-miner", "sealing jobs")
	data, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("exec lotus-miner sealing jobs: %v \n", err)
	}
	canParse := false
	var taskList []Task
	reader := bufio.NewReader(bytes.NewBuffer(data))
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
	arrs := strings.Split(line, " ")
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

func (sp *ShellParse) GetPostBalance() (string, error) {
	cmd := exec.CommandContext(context.TODO(), "lotus-miner", "actor control list")
	data, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("exec lotus-miner actor control list: %v \n", err)
	}
	postBalance := postBalanceReg.FindAllStringSubmatch(string(data), 1)
	pb := postBalance[0][1]
	return pb, nil
}

func (sp *ShellParse) GetMinerInfo() (*MinerInfo, error) {
	cmd := exec.CommandContext(context.TODO(), "lotus-miner", "info")
	data, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("exec lotus-miner info  %v \n", err)
	}
	src := string(data)
	minerInfo := &MinerInfo{}
	minerId := minerIdReg.FindString(src)
	minerInfo.minerId = minerId
	minerBalance := minerBalanceReg.FindAllStringSubmatch(src, 1)
	minerInfo.MinerBalance = minerBalance[0][1]
	workerBalance := workerBalanceReg.FindAllStringSubmatch(src, 1)
	minerInfo.workerBalance = workerBalance[0][1]
	pledgeBalance := pledgeBalanceReg.FindAllStringSubmatch(src, 1)
	minerInfo.minerId = pledgeBalance[0][1]
	totalPower := totalPowerReg.FindAllStringSubmatch(src, 1)
	minerInfo.minerId = totalPower[0][1]
	effectPower := effectPowerReg.FindAllStringSubmatch(src, 1)
	minerInfo.minerId = effectPower[0][1]
	totalSectors := totalSectorsReg.FindAllStringSubmatch(src, 1)
	minerInfo.minerId = totalSectors[0][1]
	effectSectors := effectSectorReg.FindAllStringSubmatch(src, 1)
	minerInfo.minerId = effectSectors[0][1]
	errorsSectors := errorSectorReg.FindAllStringSubmatch(src, 1)
	minerInfo.minerId = errorsSectors[0][1]
	recoverySectors := recoverySectorReg.FindAllStringSubmatch(src, 1)
	minerInfo.minerId = recoverySectors[0][1]
	deletedSectors := deletedSectorReg.FindAllStringSubmatch(src, 1)
	minerInfo.minerId = deletedSectors[0][1]
	failSectors := failSectorReg.FindAllStringSubmatch(src, 1)
	minerInfo.minerId = failSectors[0][1]
	return minerInfo, nil
}

func structToMap(obj interface{}) map[string]interface{} {
	elem := reflect.ValueOf(&obj).Elem()
	m := make(map[string]interface{})
	relType := elem.Type()
	for i := 0; i < relType.NumField(); i++ {
		m[relType.Field(i).Name] = elem.Field(i).Interface()
	}
	return m
}

func mergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			tk := k
			tV := v
			result[tk] = tV
		}
	}
	return result
}
