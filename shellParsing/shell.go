package shellParsing

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type ShellParse struct {
	Workers []WorkerInfo
}

func NewShellParse() *ShellParse {
	return &ShellParse{}
}

func (sp *ShellParse) getTaskInfo() map[string]interface{} {
	// todo
	return nil
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

func (sp *ShellParse) GetMinerJobs(res map[string]interface{}) error {
	cmd := exec.CommandContext(context.TODO(), "lotus-miner", "sealing jobs")
	data, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("exec lotus-miner sealing jobs: %v \n", err)
	}
	param := make(map[string]map[string][]interface{})
	canParse := false
	reader := bufio.NewReader(bytes.NewBuffer(data))
	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		if strings.HasPrefix(line, "ID") {
			canParse = true
		}
		// todo 优化
		if canParse {
			task := strings.Split(line, " ")

			// 判断 hostname是否存在
			if vMap, ok := param[task[3]]; ok {
				// 判断某一任务类型是否存在
				if t, ok := vMap[task[4]]; ok {
					t = append(t, map[string]interface{}{
						"type":      task[4],
						"sectorId":  task[1],
						"status":    task[5],
						"spendTime": task[6],
					})
				} else {
					vMap[task[4]] = []interface{}{
						map[string]interface{}{
							"type":      task[4],
							"sectorId":  task[1],
							"status":    task[5],
							"spendTime": task[6],
						},
					}
				}
			} else {
				param[task[3]] = map[string][]interface{}{
					task[3]: []interface{}{
						map[string]interface{}{
							"type":      task[4],
							"sectorId":  task[1],
							"status":    task[5],
							"spendTime": task[6],
						},
					},
				}
			}

		}
	}
	res["workerInfo"] = param
	return nil
}

func (sp *ShellParse) GetPostBalance(res map[string]interface{}) error {
	cmd := exec.CommandContext(context.TODO(), "lotus-miner", "actor control list")
	data, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("exec lotus-miner actor control list: %v \n", err)
	}
	postBalance := postBalanceReg.FindAllStringSubmatch(string(data), 1)
	res["postBalance"] = postBalance[0][1]
	return nil
}

func (sp *ShellParse) GetMinerInfo(res map[string]interface{}) error {
	cmd := exec.CommandContext(context.TODO(), "lotus-miner", "info")
	data, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("exec lotus-miner info  %v \n", err)
	}
	src := string(data)
	minerId := minerIdReg.FindString(src)
	res["minerId"] = minerId
	minerBalance := minerBalanceReg.FindAllStringSubmatch(src, 1)
	res["minerBalance"] = minerBalance[0][1]
	workerBalance := workerBalanceReg.FindAllStringSubmatch(src, 1)
	res["workerBalance"] = workerBalance[0][1]
	pledgeBalance := pledgeBalanceReg.FindAllStringSubmatch(src, 1)
	res["pledgeBalance"] = pledgeBalance[0][1]
	totalPower := totalPowerReg.FindAllStringSubmatch(src, 1)
	res["totalPower"] = totalPower[0][1]
	effectPower := effectPowerReg.FindAllStringSubmatch(src, 1)
	res["effectPower"] = effectPower[0][1]
	totalSectors := totalSectorsReg.FindAllStringSubmatch(src, 1)
	res["totalSectors"] = totalSectors[0][1]
	effectSectors := effectSectorReg.FindAllStringSubmatch(src, 1)
	res["effectSectors"] = effectSectors[0][1]
	errorsSectors := errorSectorReg.FindAllStringSubmatch(src, 1)
	res["errorsSectors"] = errorsSectors[0][1]
	recoverySectors := recoverySectorReg.FindAllStringSubmatch(src, 1)
	res["recoverySectors"] = recoverySectors[0][1]
	deletedSectors := deletedSectorReg.FindAllStringSubmatch(src, 1)
	res["deletedSectors"] = deletedSectors[0][1]
	failSectors := failSectorReg.FindAllStringSubmatch(src, 1)
	res["failSectors"] = failSectors[0][1]
	return nil
}
