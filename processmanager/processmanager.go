package processmanager

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

const (
	// SIGUSR1 linux SIGUSR1
	SIGUSR1 = syscall.Signal(0xa)

	// SIGUSR2 linux SIGUSR2
	SIGUSR2 = syscall.Signal(0xc)
)

// CheckPid ...
func CheckPid(processName string) {
	pidFile := "." + processName + ".pid"

	lastPidChar, fileErr := ioutil.ReadFile(pidFile)
	if fileErr == nil {
		cmd := exec.Command("kill", "-9", string(lastPidChar))
		_ = cmd.Run()
	}
	pid := os.Getpid()
	strPid := strconv.Itoa(pid)
	bytePid := []byte(strPid)
	_ = ioutil.WriteFile(pidFile, bytePid, os.ModePerm)
}

// Daemon ...
func Daemon() {
	if os.Getppid() != 1 {
		sysType := runtime.GOOS
		if sysType == "linux" {
			filePath, _ := filepath.Abs(os.Args[0])
			cmd := exec.Command(filePath, os.Args[1:]...)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			_ = cmd.Start()
			os.Exit(0)
			return
		}
	}
}

// DefExitFunc 缺省信号退出函数
func DefExitFunc() {
	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%v Default Sig Exit ... \n", now)
	os.Exit(0)
}
