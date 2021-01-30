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

var postSrc = `
owner      t07568    t3xd3xmwl...  other  14.695889803059827041 FIL  
worker     t07568    t3xd3xmwl...  other  14.695889803059827041 FIL  
control-0  t0116299  t3qsi32gm...  post   5.00000000001 FIL
`

var postStr001 = `
owner      t07568    t3xd3xmwl...  other  14.695889803059827041 FIL  
worker     t07568    t3xd3xmwl...  other  14.695889803059827041 FIL  
control-0  t0116299  t3qsi32gm...  \033[32m post \033[0m   \033[31m 98.565454456329164202 FIL \033[0m
`


func TestPost002(t *testing.T) {
	postBalance := postBalanceTestReg.FindAllStringSubmatch(postStr001, 1)

	fmt.Println(postBalance)
	fmt.Println("PostBalance:", getRegexValue(postBalance))
}

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
TotalMemory Spendable:  1279.727 FIL

Sectors:
        TotalMemory: 8468
        Proving: 7777
        Packing: 3
        PreCommit1: 28
        PreCommit2: 73
        WaitSeed: 14
        Committing: 63
        CommitWait: 2
		PreCommitWait: 10
        FinalizeSector: 82
        Removed: 316
        FailedUnrecoverable: 323
        SealPreCommit2Failed: 6`

func TestShellMinerInfo(t *testing.T) {

	result := minerIdReg.FindAllStringSubmatch(src, 1)
	fmt.Println("MinerId: ", getRegexValue(result))

	minerBalance := minerBalanceReg.FindAllStringSubmatch(src, 1)
	fmt.Println("minerBalance:  ", getRegexValue(minerBalance))

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

	effectSectors := effectSectorReg.FindAllStringSubmatch(src, 2)

	fmt.Println("effectSectors: ", getRegexValueByIndex(effectSectors, 1, 1))

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

	expectBlock := expectBlockReg.FindAllStringSubmatch(src, 1)
	fmt.Println("expectBlock: ", getRegexValue(expectBlock))

	commitWait := commitWaitReg.FindAllStringSubmatch(src, 1)
	fmt.Println("commitWait: ", getRegexValue(commitWait))

	PreCommitWait := preCommitWaitReg.FindAllStringSubmatch(src, 1)
	fmt.Println("PreCommitWait: ", getRegexValue(PreCommitWait))

	available := availableReg.FindAllStringSubmatch(src, 1)
	fmt.Println("available: ", getRegexValue(available))

}

var jobsSrc = `
build info: 0dfa8a218452bca3e8ee97abd2a3bd06cbeb2c70
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

var hardwareInfo = `
bnxt_en-pci-0201
Adapter: PCI adapter
temp1:        +54.0°C  

k10temp-pci-00c3
Adapter: PCI adapter
Tdie:         +26.5°C  (high = +70.0°C)
Tctl:         +26.5°C  

bnxt_en-pci-0200
Adapter: PCI adapter
temp1:        +54.0°C  

 03:58:55 up 5 days, 19:04,  4 users,  load average: 0.62, 2.75, 3.62
              total        used        free      shared  buff/cache   available
Mem:           503G        9.4G        378G        1.8M        116G        492G
Swap:          8.0G         63M        7.9G
Filesystem      Size  UseMemory Avail Use% Mounted on
udev            252G     0  252G   0% /dev
tmpfs            51G  2.3M   51G   1% /run
/dev/nvme0n1p2  1.5T  1.2T  221G  85% /
tmpfs           252G     0  252G   0% /dev/shm
tmpfs           5.0M  4.0K  5.0M   1% /run/lock
tmpfs           252G     0  252G   0% /sys/fs/cgroup
/dev/nvme0n1p1  511M  6.1M  505M   2% /boot/efi
tmpfs            51G     0   51G   0% /run/user/0
/dev/md126       59T   60G   59T   1% /opt/hdd_pool
Linux 5.4.0-59-generic (worker01) 	01/12/21 	_x86_64_	(48 CPU)

03:58:55        IFACE   rxpck/s   txpck/s    rxkB/s    txkB/s   rxcmp/s   txcmp/s  rxmcst/s   %ifutil
03:58:56    enp2s0f0np0      0.00      0.00      0.00      0.00      0.00      0.00      0.00      0.00
03:58:56         eno1      8.00     10.00      0.72      2.48      0.00      0.00      0.00      0.00
03:58:56         eno2      0.00      0.00      0.00      0.00      0.00      0.00      0.00      0.00
03:58:56           lo      0.00      0.00      0.00      0.00      0.00      0.00      0.00      0.00
03:58:56    enp2s0f1np1      0.00      0.00      0.00      0.00      0.00      0.00      0.00      0.00

03:58:56        IFACE   rxpck/s   txpck/s    rxkB/s    txkB/s   rxcmp/s   txcmp/s  rxmcst/s   %ifutil
03:58:57    enp2s0f0np0      1.00      0.00      0.12      0.00      0.00      0.00      1.00      0.00
03:58:57         eno1      4.00      3.00      0.26      0.95      0.00      0.00      0.00      0.00
03:58:57         eno2      0.00      0.00      0.00      0.00      0.00      0.00      0.00      0.00
03:58:57           lo      0.00      0.00      0.00      0.00      0.00      0.00      0.00      0.00
03:58:57    enp2s0f1np1      0.00      0.00      0.00      0.00      0.00      0.00      0.00      0.00

Average:        IFACE   rxpck/s   txpck/s    rxkB/s    txkB/s   rxcmp/s   txcmp/s  rxmcst/s   %ifutil
Average:    enp2s0f0np0      0.50      0.00      0.06      0.00      0.00      0.00      0.50      0.00
Average:         eno1      6.00      6.50      0.49      1.72      0.00      0.00      0.00      0.00
Average:         eno2      0.00      0.00      0.00      0.00      0.00      0.00      0.00      0.00
Average:           lo      0.00      0.00      0.00      0.00      0.00      0.00      0.00      0.00
Average:    enp2s0f1np1      0.00      0.00      0.00      0.00      0.00      0.00      0.00      0.00
TotalMemory DISK READ :       0.00 B/s | TotalMemory DISK WRITE :    1242.27 K/s
Actual DISK READ:       0.00 B/s | Actual DISK WRITE:      22.53 M/s
Tue Jan 12 03:58:58 2021       
+-----------------------------------------------------------------------------+
| NVIDIA-SMI 455.45.01    Driver Version: 455.45.01    CUDA Version: 11.1     |
|-------------------------------+----------------------+----------------------+
| GPU  Name        Persistence-M| Bus-Id        Disp.A | Volatile Uncorr. ECC |
| Fan  CpuTemp  Perf  Pwr:Usage/Cap|         Memory-Usage | GPU-Util  Compute M. |
|                               |                      |               MIG M. |
|===============================+======================+======================|
|   0  GeForce RTX 306...  Off  | 00000000:C4:00.0 Off |                  N/A |
|  0%   36C    P0    26W / 200W |      0MiB /  7982MiB |      0%      Default |
|                               |                      |                  N/A |
+-------------------------------+----------------------+----------------------+
                                                                               
+-----------------------------------------------------------------------------+
| Processes:                                                                  |
|  GPU   GI   CI        PID   Type   Process name                  GPU Memory |
|        ID   ID                                                   Usage      |
|=============================================================================|
|  No running processes found                                                 |
+-----------------------------------------------------------------------------+

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

var newMapStr = `{"deletedSectors":"1","effectivePower":"0","effectiveSectors":"0","errorSectors":"0","failSectors":"0",
"hardwareInfo":{"worker01":{"cpuLoad":"14.73","cpuTemper":"+41.1°C","diskR":"906.67M/s","diskW":"163.63M/s","gpuInfo":{"0":{"name":"0","temp":"91C","use":"100%"}},"hostName":"worker01","netIO":{"eno1":{"name":"eno1","rx":"1.27","tx":"2.90"},"eno2":{"name":"eno2","rx":"0.00","tx":"0.00"},"enp2s0f0np0":{"name":"enp2s0f0np0","rx":"0.00","tx":"0.00"},"enp2s0f1np1":{"name":"enp2s0f1np1","rx":"0.00","tx":"0.00"},"lo":{"name":"lo","rx":"0.00","tx":"0.00"}},"totalMemory":"503G","useDisk":"40%","useMemory":"319G"}},
"jobs":{"17":{"hostName":"worker01","id":"d7fd42c9","sector":"17","state":"running","task":"PC1","time":"17m48s","worker":"98c441ab"},"40":{"hostName":"worker01","id":"f5d60859","sector":"40","state":"running","task":"PC2","time":"20m17.7s","worker":"98c441ab"},"47":{"hostName":"worker01","id":"4fda8ea1","sector":"47","state":"running","task":"PC1","time":"13m17.7s","worker":"98c441ab"},"48":{"hostName":"worker01","id":"eba6ef90","sector":"48","state":"running","task":"PC1","time":"16m17.9s","worker":"98c441ab"},"49":{"hostName":"worker01","id":"d3283f2f","sector":"49","state":"running","task":"PC1","time":"19m18s","worker":"98c441ab"},"50":{"hostName":"worker01","id":"c1309585","sector":"50","state":"running","task":"PC1","time":"14m48s","worker":"98c441ab"},"51":{"hostName":"worker01","id":"53185af0","sector":"51","state":"running","task":"PC1","time":"11m47s","worker":"98c441ab"}},"messageNums":9,"minerBalance":"0","minerId":"t0114613","pledgeBalance":"0","postBalance":"0","recoverySectors":"0","timestamp":1610171829,"totalSectors":"52","workerBalance":"39.522FIL"}`

var oldMapStr = `{"deletedSectors":"1","effectivePower":"0","effectiveSectors":"0","errorSectors":"0","failSectors":"0",
"hardwareInfo":{"worker01":{"cpuLoad":"14.73","cpuTemper":"+41.1°C","diskR":"906.67M/s","diskW":"163.63M/s","gpuInfo":{"0":{"name":"0","temp":"91C","use":"100%"}},"hostName":"worker01","netIO":{"eno1":{"name":"eno1","rx":"1.27","tx":"2.90"},"eno2":{"name":"eno2","rx":"0.00","tx":"0.00"},"enp2s0f0np0":{"name":"enp2s0f0np0","rx":"0.00","tx":"0.00"},"enp2s0f1np1":{"name":"enp2s0f1np1","rx":"0.00","tx":"0.00"},"lo":{"name":"lo","rx":"0.00","tx":"0.00"}},"totalMemory":"503G","useDisk":"40%","useMemory":"319G"}},
"jobs":{"18":{"hostName":"worker01","id":"d7fd42c9","sector":"17","state":"running","task":"PC1","time":"17m48s","worker":"98c441ab"},"40":{"hostName":"worker01","id":"f5d60859","sector":"40","state":"running","task":"PC2","time":"20m17.7s","worker":"98c441ab"},"47":{"hostName":"worker01","id":"4fda8ea1","sector":"47","state":"running","task":"PC1","time":"13m17.7s","worker":"98c441ab"},"48":{"hostName":"worker01","id":"eba6ef90","sector":"48","state":"running","task":"PC1","time":"16m17.9s","worker":"98c441ab"},"49":{"hostName":"worker01","id":"d3283f2f","sector":"49","state":"running","task":"PC1","time":"19m18s","worker":"98c441ab"},"50":{"hostName":"worker01","id":"c1309585","sector":"50","state":"running","task":"PC1","time":"14m48s","worker":"98c441ab"},"51":{"hostName":"worker01","id":"53185af0","sector":"51","state":"running","task":"PC1","time":"11m47s","worker":"98c441ab"}},"messageNums":9,"minerBalance":"0","minerId":"t0114613","pledgeBalance":"0","postBalance":"0","recoverySectors":"0","timestamp":1610171829,"totalSectors":"52","workerBalance":"39.522FIL"}`

func TestDiffMap(t *testing.T) {
	oldMap := make(map[string]interface{})
	newMap := make(map[string]interface{})

	err := json.Unmarshal([]byte(oldMapStr), &oldMap)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = json.Unmarshal([]byte(newMapStr), &newMap)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	diffMap := DiffMap(oldMap, newMap)
	data, err := json.MarshalIndent(diffMap, " ", " ")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("diffMap:  ", string(data))
	jobs, ok := diffMap["jobs"]
	if !ok {
		return
	}
	tJobs := jobs.(map[string]interface{})

	//hardwareInfo,ok := diffMap["hardwareInfo"]
	//if !ok{
	//	return
	//}
	//tHardware := hardwareInfo.(map[string]interface{})

	//info := ParseJobsInfo(tJobs, nil)

	info := MapParse(tJobs, nil)

	resData, err := json.MarshalIndent(info, " ", " ")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("workerInfo: ", string(resData))
}

func TestParseMinerInfo(t *testing.T) {
	param := make(map[string]interface{})
	err := json.Unmarshal([]byte(testStr), &param)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	jobs := make(map[string]interface{})
	hardwareInfo := make(map[string]interface{})
	tJobs, ok := param["jobs"]
	if ok {
		jobs = tJobs.(map[string]interface{})
	}
	tHardwareInfo, ok := param["hardwareInfo"]
	if ok {
		hardwareInfo = tHardwareInfo.(map[string]interface{})
	}

	result := MapParse(jobs, hardwareInfo)
	data, err := json.MarshalIndent(result, "  ", "   ")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(data))
}

var testStr = `{"deletedSectors":"1","effectivePower":"0","effectiveSectors":"0","errorSectors":"0","failSectors":"0","hardwareInfo":{"worker01":{"cpuLoad":"14.73","cpuTemper":"+41.1°C","diskR":"906.67M/s","diskW":"163.63M/s","gpuInfo":{"0":{"name":"0","temp":"91C","use":"100%"}},"hostName":"worker01","netIO":{"eno1":{"name":"eno1","rx":"1.27","tx":"2.90"},"eno2":{"name":"eno2","rx":"0.00","tx":"0.00"},"enp2s0f0np0":{"name":"enp2s0f0np0","rx":"0.00","tx":"0.00"},"enp2s0f1np1":{"name":"enp2s0f1np1","rx":"0.00","tx":"0.00"},"lo":{"name":"lo","rx":"0.00","tx":"0.00"}},"totalMemory":"503G","useDisk":"40%","useMemory":"319G"}},"jobs":{"17":{"hostName":"worker01","id":"d7fd42c9","sector":"17","state":"running","task":"PC1","time":"17m48s","worker":"98c441ab"},"40":{"hostName":"worker01","id":"f5d60859","sector":"40","state":"running","task":"PC2","time":"20m17.7s","worker":"98c441ab"},"47":{"hostName":"worker01","id":"4fda8ea1","sector":"47","state":"running","task":"PC1","time":"13m17.7s","worker":"98c441ab"},"48":{"hostName":"worker01","id":"eba6ef90","sector":"48","state":"running","task":"PC1","time":"16m17.9s","worker":"98c441ab"},"49":{"hostName":"worker01","id":"d3283f2f","sector":"49","state":"running","task":"PC1","time":"19m18s","worker":"98c441ab"},"50":{"hostName":"worker01","id":"c1309585","sector":"50","state":"running","task":"PC1","time":"14m48s","worker":"98c441ab"},"51":{"hostName":"worker01","id":"53185af0","sector":"51","state":"running","task":"PC1","time":"11m47s","worker":"98c441ab"}},"messageNums":9,"minerBalance":"0","minerId":"t0114613","pledgeBalance":"0","postBalance":"0","recoverySectors":"0","timestamp":1610171829,"totalSectors":"52","workerBalance":"39.522FIL"}`

var minerWorkers = `
Worker 0042426f-fcb5-4872-a6e5-58108b61ea8b, host ya-node111 (disabled) tasks AP|C1|PC1|PC2-0
        CPU:  [|||||                                                           ] 4/48 core(s) in use
        RAM:  [||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||    ] 95% 481.8 GiB/503.4 GiB
        VMEM: [|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||| ] 108% 545.8 GiB/503.4 GiB
        GPU: GeForce RTX 3060 Ti, used
Worker 0cb8f1a3-ed67-4b0b-86fd-a5535b3ae9fe, host ya-node102 (disabled) tasks AP|C1|PC1|PC2-0
        CPU:  [|||||                                                           ] 4/48 core(s) in use
        RAM:  [||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||    ] 95% 481.7 GiB/503.4 GiB
        VMEM: [|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||| ] 108% 545.7 GiB/503.4 GiB
        GPU: GeForce RTX 3060 Ti, used
Worker 44fdd58f-0c5d-4a99-9d9b-69e2537029c5, host ya-node98-miner tasks RD
        CPU:  [                                                                ] 0/48 core(s) in use
        RAM:  [||||||                                                          ] 9% 37.66 GiB/377.6 GiB
        VMEM: [||||||                                                          ] 9% 37.66 GiB/379.6 GiB
        GPU: GeForce RTX 2080 Ti, not used
Worker 48d674a7-62a2-4044-b7e5-0b8f4195ebc7, host ya-node87 (disabled) tasks C2
        CPU:  [                                                                ] 0/32 core(s) in use
        RAM:  [                                                                ] 0% 3.153 GiB/377.6 GiB
        VMEM: [                                                                ] 0% 3.153 GiB/377.6 GiB
        GPU: GeForce RTX 3080, not used
        GPU: GeForce RTX 3080, not used
Worker 5403b7bd-b126-4aa6-96bf-5e92ec2ad80a, host ya-node109 tasks AP|C1|PC1|PC2-0
        CPU:  [|||||                                                           ] 4/48 core(s) in use
        RAM:  [||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||    ] 95% 481.7 GiB/503.4 GiB
        VMEM: [|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||| ] 108% 545.7 GiB/503.4 GiB
        GPU: GeForce RTX 3060 Ti, used
Worker 6770dc35-a179-4c4e-9736-d2da32b36be9, host ya-node93 tasks C2
        CPU:  [                                                                ] 0/32 core(s) in use
        RAM:  [                                                                ] 0% 3.125 GiB/377.6 GiB
        VMEM: [                                                                ] 0% 3.125 GiB/377.6 GiB
        GPU: GeForce RTX 3080, not used
        GPU: GeForce RTX 3080, not used`

func TestMinerWorkers(t *testing.T) {
	reader := bufio.NewReader(bytes.NewBuffer([]byte(minerWorkers)))
	param := make(map[string]*WorkerInfo)
	preHostName := ""
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
			fmt.Println(line)
			fields := strings.Fields(line)
			fmt.Println(len(fields))
			if len(fields) < 6 {
				continue
			}
			preHostName = fields[3]
			hostType := strings.Split(fields[5], "|")
			param[fields[3]] = &WorkerInfo{HostName: fields[3], TaskState: taskState, TaskType: hostType}

		} else if strings.Contains(line, "GPU") {
			workerInfo, ok := param[preHostName]
			if ok {
				workerInfo.GPU = 1
			}
		}

	}
	for key, value := range param {
		fmt.Println(key, value)
	}
}

func DeepCopyMap(input map[string]interface{}) (map[string]interface{}, error) {
	param := make(map[string]interface{})
	data, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &param)
	if err != nil {
		return nil, err
	}
	return param, nil
}

func Test03(t *testing.T) {
	test := make([]int, 0, 100)
	fmt.Println(test[0])
}

func Test02(t *testing.T) {
	users := []string{"01", "02", "03", "04"}
	for index, info := range users[1:] {
		fmt.Println(index, info)
	}
	users01 := []string{"aa", "bb"}
	copy(users[1:], users[3:])
	fmt.Println(users)
	fmt.Println(users01)

}

func Test01(t *testing.T) {
	param := make(map[string]interface{})
	param["test"] = 01

	copyMap, err := DeepCopyMap(param)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	param["test"] = 1000
	fmt.Println(copyMap)
	fmt.Println(param)
}
