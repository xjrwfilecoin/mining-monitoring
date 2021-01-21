package app

import (
	"encoding/json"
	"fmt"
	"testing"
)

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
	fmt.Println(string(data))
}