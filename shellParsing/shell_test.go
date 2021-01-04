package shellParsing

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"testing"
)

//var minerIdReg = regexp.MustCompile(` t[\d]* `)
//var minerBalanceReg = regexp.MustCompile(`Miner Balance:    ([\d]*.*[\d]*.*FIL)`)
//var postBalanceReg = regexp.MustCompile(`\.\.\.  post        ([\d]*.*[\d]*.*FIL)`)
//var workerBalanceReg = regexp.MustCompile(`Worker Balance:   ([\d]*.*[\d]*.*FIL)`)
//var pledgeBalanceReg = regexp.MustCompile(`Pledge:     ([\d]*.*[\d]*.*FIL)`)
//var effectPowerReg = regexp.MustCompile(`Power: ([\d]*.*[\d].*) /`)
//var totalPowerReg = regexp.MustCompile(`Committed: ([\d]*.*[\d].*)`)
//var totalSectorsReg = regexp.MustCompile(`Total: ([\d]*)`)
//var effectSectorReg = regexp.MustCompile(`Proving: ([\d]*)`)
//var errorSectorReg = regexp.MustCompile(`FailedUnrecoverable: ([\d]*)`)
//var recoverySectorReg = regexp.MustCompile(`SealPreCommit2Failed: ([\d]*)`)
//var deletedSectorReg = regexp.MustCompile(`Removed: ([\d]*)`)
//var failSectorReg = regexp.MustCompile(`SealPreCommit2Failed: ([\d]*)`)



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

	result := minerIdReg.FindString(src)
	fmt.Println("minerId: ",result)

	minerBalance := minerBalanceReg.FindAllStringSubmatch(src, 1)
	fmt.Println("minerBalance:  ",minerBalance[0][1])

	postBalance := postBalanceReg.FindAllStringSubmatch(postSrc, 1)
	fmt.Println("postBalance: ",postBalance[0][1])

	workerBalance := workerBalanceReg.FindAllStringSubmatch(src, 1)
	fmt.Println("workerBalance:  ",workerBalance[0][1])

	pledgeBalance := pledgeBalanceReg.FindAllStringSubmatch(src, 1)
	fmt.Println("pledgeBalance:  ",pledgeBalance[0][1])

	totalPower := totalPowerReg.FindAllStringSubmatch(src, 1)
	fmt.Println("totalPower: ",totalPower[0][1])

	effectPower := effectPowerReg.FindAllStringSubmatch(src, 1)
	fmt.Println("effectPower: ",effectPower[0][1])

	totalSectors := totalSectorsReg.FindAllStringSubmatch(src, 1)
	fmt.Println("totalSectors: ",totalSectors[0][1])

	effectSectors := effectSectorReg.FindAllStringSubmatch(src, 1)
	fmt.Println("effectSectors: ",effectSectors[0][1])

	errorsSectors := errorSectorReg.FindAllStringSubmatch(src, 1)
	fmt.Println("errorsSectors: ",errorsSectors[0][1])

	recoverySectors := recoverySectorReg.FindAllStringSubmatch(src, 1)
	fmt.Println("recoverySectors: ",recoverySectors[0][1])

	deletedSectors := deletedSectorReg.FindAllStringSubmatch(src, 1)
	fmt.Println("deletedSectors: ",deletedSectors[0][1])

	failSectors := failSectorReg.FindAllStringSubmatch(src, 1)
	fmt.Println("failSectors: ",failSectors[0][1])

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
func TestWorkerJobs(t *testing.T){

	//result :=make(map[string]map[string]interface{})

	reader := bufio.NewReader(bytes.NewBuffer([]byte(jobsSrc)))
	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		fmt.Println(line)
	}
}

var hardwareInfo =`172.70.16.101 | CHANGED | rc=0 >>
k10temp-pci-00c3
Adapter: PCI adapter
temp1:         +0.0°C  (high = +70.0°C)

nouveau-pci-0400
Adapter: PCI adapter
GPU core:     +0.91 V  (min =  +0.80 V, max =  +1.19 V)
temp1:        +27.0°C  (high = +95.0°C, hyst =  +3.0°C)
                       (crit = +105.0°C, hyst =  +5.0°C)
                       (emerg = +135.0°C, hyst =  +5.0°C)

 15:27:19 up 25 days, 12:53,  2 users,  load average: 4.94, 4.19, 4.35
              total        used        free      shared  buff/cache   available
Mem:           125G        114G        732M        2.4M         10G         10G
Swap:            0B          0B          0B
Filesystem      Size  Used Avail Use% Mounted on
udev             63G     0   63G   0% /dev
tmpfs            13G  1.2M   13G   1% /run
/dev/nvme0n1p2  1.8T  528G  1.2T  32% /
tmpfs            63G     0   63G   0% /dev/shm
tmpfs           5.0M     0  5.0M   0% /run/lock
tmpfs            63G     0   63G   0% /sys/fs/cgroup
/dev/nvme0n1p1  511M  6.1M  505M   2% /boot/efi
/dev/md127      7.3T  2.5T  4.9T  34% /opt/hdd_pool
tmpfs            13G     0   13G   0% /run/user/0
tmpfs            13G     0   13G   0% /run/user/1000
Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
enp10s0: 515659456584 4469525963    0 83152    0     0          0   1103181 133117032642438 89567040730    0    0    0     0       0          0
docker0: 14297868050 2605023    0    0    0     0          0         0 392878806 5683481    0    0    0     0       0          0
    lo:   86178     854    0    0    0     0          0         0    86178     854    0    0    0     0       0          0
  eno1: 1322093904 14310224    0 427497    0     0          0   1174643 15325372882 16834889    0    0    0     0       0          0`

var cpuTemperatureReg = regexp.MustCompile(``)
var cpuLoadReg = regexp.MustCompile(``)
var gpuTemperatureReg = regexp.MustCompile(``)
var gpuLoadReg = regexp.MustCompile(``)
var memoryUsedReg = regexp.MustCompile(``)
var memoryTotalReg = regexp.MustCompile(``)
var diskUsedRateReg =regexp.MustCompile(``)


func TestHardwareInfo(t *testing.T){

	cpuTemperature := cpuTemperatureReg.FindAllStringSubmatch(hardwareInfo, 1)
	fmt.Println("cpuTemperature: ",cpuTemperature[0][1])

	cpuLoad := cpuLoadReg.FindAllStringSubmatch(hardwareInfo, 1)
	fmt.Println("cpuLoad: ",cpuLoad[0][1])


	gpuLoad := gpuLoadReg.FindAllStringSubmatch(hardwareInfo, 1)
	fmt.Println("gpuLoad: ",gpuLoad[0][1])

	memoryUsed := memoryUsedReg.FindAllStringSubmatch(hardwareInfo, 1)
	fmt.Println("memoryUsed: ",memoryUsed[0][1])

	memoryTotal := memoryTotalReg.FindAllStringSubmatch(hardwareInfo, 1)
	fmt.Println("memoryTotal: ",memoryTotal[0][1])

	diskUsed := diskUsedRateReg.FindAllStringSubmatch(hardwareInfo, 1)
	fmt.Println("diskUsed: ",diskUsed[0][1])


}



