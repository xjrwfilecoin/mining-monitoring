package shellParsing

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"
)

var postSrc = `name       ID      key           use         balance                          
owner      t0100   t3qxt533a...  other post  49899998.999888415326285699 FIL  
worker     t0100   t3qxt533a...  other post  49899998.999888415326285699 FIL  
control-0  t0100   t3qxt533a...  other post  49899998.999888415326285699 FIL  
control-1  t01001  t3whckee7...  post        100000.999952762113703286 FIL`

var src = `build info: 0dfa8a218452bca3e8ee97abd2a3bd06cbeb2c70
localIP:  172.70.16.201
Chain: [sync ok] [basefee 3.138 nFIL]
Miner: f096920 (32 GiB sectors)
Power: 230 Ti / 1.71 Ei (0.0127%)
        Raw: 229.6 TiB / 1.713 EiB (0.0127%)
        Committed: 238.5 TiB
        Proving: 229.6 TiB
Expected block win rate: 1.8288/day (every 13h7m24s)

Deals: 3, 1.75 GiB
        Active: 0, 0 B (Verified: 0, 0 B)

Miner Balance:    2502.748 FIL
      PreCommit:  6.393 FIL
      Pledge:     1914.36 FIL
      Vesting:    391.076 FIL
      Available:  190.919 FIL
Market Balance:   4.075 mFIL
       Locked:    3.702 mFIL
       Available: 373.466 μFIL
Worker Balance:   968.058 FIL
       Control:   120.75 FIL
Total Spendable:  1279.727 FIL

Sectors:
        Total: 8468
        Proving: 7558
        Packing: 3
        PreCommit1: 28
        PreCommit2: 73
        WaitSeed: 14
        Committing: 63
        CommitWait: 2
        FinalizeSector: 82
        Removed: 316
        FailedUnrecoverable: 323
        SealPreCommit2Failed: 6`

func TestShellMinerInfo(t *testing.T) {

	result := minerIdReg.FindAllStringSubmatch(src, 1)
	fmt.Println("MinerId: ", getRegexValue(result))

	minerBalance := minerBalanceReg.FindAllStringSubmatch(src, 1)
	fmt.Println("minerBalance:  ", getRegexValue(minerBalance))

	postBalance := postBalanceReg.FindAllStringSubmatch(postSrc, 1)
	fmt.Println("PostBalance: ", getRegexValue(postBalance))

	workerBalance := workerBalanceReg.FindAllStringSubmatch(src, 1)
	fmt.Println("WorkerBalance:  ", getRegexValue(workerBalance))

	pledgeBalance := pledgeBalanceReg.FindAllStringSubmatch(src, 1)
	fmt.Println("PledgeBalance:  ", getRegexValue(pledgeBalance))

	totalPower := totalPowerReg.FindAllStringSubmatch(src, 1)
	fmt.Println("totalPower: ", getRegexValue(totalPower))

	effectPower := effectPowerReg.FindAllStringSubmatch(src, 1)
	fmt.Println("effectPower: ", getRegexValue(effectPower))

	totalSectors := totalSectorsReg.FindAllStringSubmatch(src, 1)
	fmt.Println("TotalSectors: ", getRegexValue(totalSectors))

	effectSectors := effectSectorReg.FindAllStringSubmatch(src, 1)
	fmt.Println("effectSectors: ", getRegexValue(effectSectors))

	errorsSectors := errorSectorReg.FindAllStringSubmatch(src, 1)
	fmt.Println("errorsSectors: ", getRegexValue(errorsSectors))

	recoverySectors := recoverySectorReg.FindAllStringSubmatch(src, 1)
	fmt.Println("RecoverySectors: ", getRegexValue(recoverySectors))

	deletedSectors := deletedSectorReg.FindAllStringSubmatch(src, 1)
	fmt.Println("DeletedSectors: ", getRegexValue(deletedSectors))

	failSectors := failSectorReg.FindAllStringSubmatch(src, 1)
	fmt.Println("FailSectors: ", getRegexValue(failSectors))

	preCommitFailed := preCommitFailedReg.FindAllStringSubmatch(src, 1)
	fmt.Println("preCommitFailed: ", getRegexValue(preCommitFailed))

}

var jobsSrc = `build info: 0dfa8a218452bca3e8ee97abd2a3bd06cbeb2c70
localIP:  172.70.16.201
ID        Sector  Worker    Hostname       Task  State        Time
c71e05fc  8598    74d84e37  ya_amd_node36  PC1   running      2h12m29.5s
b17ec3eb  8599    6a38fdf0  ya_amd_node18  PC1   running      2h11m56.6s
46118a65  8600    72f03062  ya_amd_node22  PC1   running      2h9m46s
c235c6fc  8553    fe77e2ff  ya_gpu_node06  PC2   running      1h26m50.2s
6476e4c6  8531    03892a81  ya_gpu_node09  C2    running      1h16m2.8s
0c964e5e  8601    7c221333  ya_amd_node16  PC1   running      1h13m37.4s
018c2e11  8602    d0ee84bc  ya_amd_node25  PC1   running      1h6m59.5s
07ba1a71  8603    a601c886  ya_amd_node33  PC1   running      1h6m54.9s
f4cd98b8  8604    2f9c7eb8  ya_amd_node14  PC1   running      1h6m5.6s
c7f267a1  8339    ba329dcb  ya_gpu_node10  C2    running      1h1m12.8s
bd48e86b  8346    ba329dcb  ya_gpu_node10  C2    running      1h1m12s
f4f924ee  8556    3e3bd7ad  ya_amd_node34  PC1   running      57m28s
d5859116  8366    2668f692  ya_amd_node01  PC1   running      56m38.1s
2e202937  8372    0277fc73  ya_amd_node27  PC1   running      56m33.7s
b04a2ae0  8581    c9f19254  ya_amd_node31  PC1   running      52m46s
10023a79  8571    2668f692  ya_amd_node01  PC1   running      52m43.7s
d836c5ca  8580    72f03062  ya_amd_node22  PC1   running      52m41.5s
f5a8cd29  8572    efc31ace  ya_amd_node08  PC1   running      52m41.2s
2e9c8dd3  8605    0877ce65  ya_amd_node05  PC1   running      51m49.5s
2eee566a  8575    3e8fc241  ya_amd_node21  PC1   running      51m48.2s
a9bed652  8515    88bfb751  ya_gpu_node03  C2    running      40m28s
21851c28  8526    88bfb751  ya_gpu_node03  C2    running      40m27.3s
34a9edd1  8590    27cfe0c3  ya_amd_node02  PC1   running      39m37.9s
b6f7ef87  8362    0877ce65  ya_amd_node05  PC1   running      24m28.9s
88571edc  8606    30219202  ya_amd_node10  PC1   running      24m8.8s
13221e93  8473    b7f2929c  ya_amd_node06  PC1   running      23m51.2s
`

func TestWorkerJobs(t *testing.T) {

	//result :=make(map[string]map[string]interface{})
	reader := bufio.NewReader(bytes.NewBuffer([]byte(jobsSrc)))
	var canParse bool
	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		if !canParse && strings.HasPrefix(line, "ID") {
			canParse = true
			continue
		}
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
		fmt.Println(task)
	}
}

var hardwareInfo = `power_meter-acpi-0
Adapter: ACPI interface
power1:      184.00 W  (interval =   1.00 s)

k10temp-pci-00c3
Adapter: PCI adapter
Tdie:         +51.6°C  (high = +70.0°C)
Tctl:         +51.6°C  
coretemp-isa-0000
Adapter: ISA adapter
Package id 0:  +43.0 C  (high = +81.0 C, crit = +91.0 C)
Core 0:        +42.0 C  (high = +81.0 C, crit = +91.0 C)
Core 1:        +40.0 C  (high = +81.0 C, crit = +91.0 C)
Core 2:        +41.0 C  (high = +81.0 C, crit = +91.0 C)
Core 3:        +40.0 C  (high = +81.0 C, crit = +91.0 C)
Core 4:        +41.0 C  (high = +81.0 C, crit = +91.0 C)
Core 5:        +40.0 C  (high = +81.0 C, crit = +91.0 C)
Core 6:        +38.0 C  (high = +81.0 C, crit = +91.0 C)
Core 7:        +40.0 C  (high = +81.0 C, crit = +91.0 C)
Core 8:        +41.0 C  (high = +81.0 C, crit = +91.0 C)
Core 9:        +41.0 C  (high = +81.0 C, crit = +91.0 C)
Core 10:       +40.0 C  (high = +81.0 C, crit = +91.0 C)
Core 11:       +40.0 C  (high = +81.0 C, crit = +91.0 C)
Core 12:       +41.0 C  (high = +81.0 C, crit = +91.0 C)
Core 13:       +40.0 C  (high = +81.0 C, crit = +91.0 C)
Core 14:       +42.0 C  (high = +81.0 C, crit = +91.0 C)
Core 15:       +40.0 C  (high = +81.0 C, crit = +91.0 C)

 11:24:15 up 5 days, 21:51,  1 user,  load average: 20.45, 17.81, 17.44
              total        used        free      shared  buff/cache   available
Mem:           251G        101G        1.9G         18M        147G        148G
Swap:            0B          0B          0B
Filesystem               Size  Used Avail Use% Mounted on
devtmpfs                 126G     0  126G   0% /dev
tmpfs                    126G     0  126G   0% /dev/shm
tmpfs                    126G   19M  126G   1% /run
tmpfs                    126G     0  126G   0% /sys/fs/cgroup
/dev/mapper/centos-root  372G  6.5G  365G   2% /
/dev/sda2               1014M  147M  868M  15% /boot
/dev/sda1                200M   12M  189M   6% /boot/efi
/dev/md127                51T  592G   51T   2% /opt/hdd_pool
tmpfs                     26G     0   26G   0% /run/user/1000
tmpfs                     26G     0   26G   0% /run/user/0
Linux 3.10.0-1127.el7.x86_64 (localhost.localdomain)    01/05/21        _x86_64_        (32 CPU)

11:24:15        IFACE   rxpck/s   txpck/s    rxkB/s    txkB/s   rxcmp/s   txcmp/s  rxmcst/s
11:24:16    enp180s0f0      5.00      0.00      0.29      0.00      0.00      0.00      0.00
11:24:16         eno1  44111.00  44003.00  28076.70   7590.19      0.00      0.00      0.00
11:24:16         eno2      0.00      0.00      0.00      0.00      0.00      0.00      0.00
11:24:16    enp180s0f1      5.00      0.00      0.29      0.00      0.00      0.00      0.00
11:24:16           lo      0.00      0.00      0.00      0.00      0.00      0.00      0.00

11:24:16        IFACE   rxpck/s   txpck/s    rxkB/s    txkB/s   rxcmp/s   txcmp/s  rxmcst/s
11:24:17    enp180s0f0      5.00      0.00      0.29      0.00      0.00      0.00      0.00
11:24:17         eno1  86274.00  86230.00  54972.77  14884.79      0.00      0.00      0.00
11:24:17         eno2      0.00      0.00      0.00      0.00      0.00      0.00      0.00
11:24:17    enp180s0f1      5.00      0.00      0.29      0.00      0.00      0.00      0.00
11:24:17           lo      0.00      0.00      0.00      0.00      0.00      0.00      0.00

Average:        IFACE   rxpck/s   txpck/s    rxkB/s    txkB/s   rxcmp/s   txcmp/s  rxmcst/s
Average:    enp180s0f0      5.00      0.00      0.29      0.00      0.00      0.00      0.00
Average:         eno1  65192.50  65116.50  41524.73  11237.49      0.00      0.00      0.00
Average:         eno2      0.00      0.00      0.00      0.00      0.00      0.00      0.00
Average:    enp180s0f1      5.00      0.00      0.29      0.00      0.00      0.00      0.00
Average:           lo      0.00      0.00      0.00      0.00      0.00      0.00      0.00
unable to set locale, falling back to the default locale
Total DISK READ :       0.00 B/s | Total DISK WRITE :      27.47 M/s
Actual DISK READ:      57.36 M/s | Actual DISK WRITE:       0.00 B/s
Fri Jan  8 11:46:08 2021
+-----------------------------------------------------------------------------+
| NVIDIA-SMI 455.23.04    Driver Version: 455.23.04    CUDA Version: 11.1     |
|-------------------------------+----------------------+----------------------+
| GPU  Name        Persistence-M| Bus-Id        Disp.A | Volatile Uncorr. ECC |
| Fan  Temp  Perf  Pwr:Usage/Cap|         Memory-Usage | GPU-Util  Compute M. |
|                               |                      |               MIG M. |
|===============================+======================+======================|
|   0  GeForce RTX 3080    Off  | 00000000:02:00.0 Off |                  N/A |
|  0%   21C    P8    11W / 320W |      2MiB / 10018MiB |      90%      Default |
|                               |                      |                  N/A |
+-------------------------------+----------------------+----------------------+
|   1  GeForce RTX 3080    Off  | 00000000:C1:00.0 Off |                  N/A |
|  0%   24C    P8     3W / 320W |      2MiB / 10018MiB |      0%      Default |
|                               |                      |                  N/A |
+-------------------------------+----------------------+----------------------+

+-----------------------------------------------------------------------------+
| Processes:                                                                  |
|  GPU   GI   CI        PID   Type   Process name                  GPU Memory |
|        ID   ID                                                   Usage      |
|=============================================================================|
|  No running processes found                                                 |
`

func TestHardwareInfo(t *testing.T) {

	cpuTemperature := getCpuTemper(hardwareInfo)
	fmt.Println("cpuTemperature: ", cpuTemperature)

	cpuLoad := cpuLoadReg.FindAllStringSubmatch(hardwareInfo, 1)
	fmt.Println("CpuLoad: ", getRegexValue(cpuLoad))

	memoryUsed := memoryUsedReg.FindAllStringSubmatch(hardwareInfo, 1)
	fmt.Println("memoryUsed: ", getRegexValueById(memoryUsed, 2))
	fmt.Println("memoryTotal: ", getRegexValueById(memoryUsed, 1))

	diskUsed := diskUsedRateReg.FindAllStringSubmatch(hardwareInfo, 1)
	fmt.Println("diskUsed: ", getRegexValue(diskUsed))

	diskRead := diskReadReg.FindAllStringSubmatch(hardwareInfo, 1)
	fmt.Println("diskRead: ", getRegexValue(diskRead))

	diskWrite := diskWriteReg.FindAllStringSubmatch(hardwareInfo, 1)
	fmt.Println("diskWrite: ", getRegexValue(diskWrite))

	netIO := getNetIO(hardwareInfo)
	fmt.Println("netIO: ", netIO)

	graphicsCardInfo := getGraphicsCardInfo(hardwareInfo)
	fmt.Println("graphicsCardInfo: ", graphicsCardInfo)
}

func TestWorkerInfo(t *testing.T) {
	tasks := []Task{
		{
			Id:       "ddddd",
			Sector:   "1",
			Worker:   "dddd",
			HostName: "worker01",
			Task:     "PC1",
			State:    "running",
			Time:     "1h52m",
		},
		{
			Id:       "ddddd",
			Sector:   "2",
			Worker:   "dddd",
			HostName: "worker01",
			Task:     "PC1",
			State:    "running",
			Time:     "1h52m",
		},
		{
			Id:       "ddddd",
			Sector:   "3",
			Worker:   "dddd",
			HostName: "worker01",
			Task:     "PC2",
			State:    "wait",
			Time:     "1h52m",
		},
		{
			Id:       "ddddd",
			Sector:   "3",
			Worker:   "dddd",
			HostName: "worker02",
			Task:     "C2",
			State:    "running",
			Time:     "1h52m",
		},
	}
	hardwareInfo := []HardwareInfo{
		{
			HostName:    "worker01",
			CpuTemper:   "100",
			CpuLoad:     "100",
			GpuTemper:   "100",
			GpuLoad:     "100",
			UseMemory:   "100",
			TotalMemory: "100",
			UseDisk:     "100",
			DiskR:       "100",
			DiskW:       "100",
			NetIO:       "100",
		},
		{
			HostName:    "worker02",
			CpuTemper:   "100",
			CpuLoad:     "100",
			GpuTemper:   "100",
			GpuLoad:     "100",
			UseMemory:   "100",
			TotalMemory: "100",
			UseDisk:     "100",
			DiskR:       "100",
			DiskW:       "100",
			NetIO:       "100",
		},
	}
	info := mergeWorkerInfo(tasks, hardwareInfo)
	data, _ := json.MarshalIndent(info, "   ", "   ")
	fmt.Printf("result: %v \n", string(data))
}

func TestShellParse(t *testing.T) {

}
