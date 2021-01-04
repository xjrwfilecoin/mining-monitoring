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
}

func NewShellParse() *ShellParse {
	return &ShellParse{}
}

func (sp *ShellParse) getCurrentInfo() map[string]interface{} {
	// todo
	return nil
}

func (sp *ShellParse) getMinerJobs(res map[string]interface{}) (error) {
	cmd := exec.CommandContext(context.TODO(), "lotus-miner", "sealing jobs")
	data, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("exec lotus-miner sealing jobs: %v \n", err)
	}
	param := make(map[string]map[string]interface{})
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
		if canParse {
			 := strings.Split(line, " ")
		}
	}
	return nil
}






func (sp *ShellParse) getPostBalance(res map[string]interface{}) error {
	cmd := exec.CommandContext(context.TODO(), "lotus-miner", "actor control list")
	data, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("exec lotus-miner actor control list: %v \n", err)
	}
	postBalance := postBalanceReg.FindAllStringSubmatch(string(data), 1)
	res["postBalance"] = postBalance[0][1]
	return nil
}

func (sp *ShellParse) getMinerInfo(res map[string]interface{}) error {
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
