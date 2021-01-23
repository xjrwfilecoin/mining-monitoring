package shellParsing

import (
	"sync"
)

type WorkerInfo01 struct {
	MinerId     string
	HostName    string
	Gpu         int
	DiskR       string
	DiskW       string
	UseDisk     string
	CpuLoad     string
	UseMemory   string
	TotalMemory string
	CpuTemp     string
	GpuInfo     []GpuInfo
	NetIO       []NetIO
	Jobs        []Task
}

// todo 兼容,比对差异更新
func (w *WorkerInfo01) updateInfo(data CmdData) {
	switch data.CmdType {
	case IOCmd:
		info := data.Data.(IoInfo)
		w.DiskR = info.DiskR
		w.DiskW = info.DiskW
		break
	case SarCmd:
		w.NetIO = data.Data.([]NetIO)
		break
	case DfHCMd:
		w.UseDisk = data.Data.(Disk).UseDisk
		break
	case FreeHCmd:
		memory := data.Data.(Memory)
		w.TotalMemory = memory.TotalMemory
		w.UseMemory = memory.UseMemory
		break
	case SensorsCmd:
		w.CpuTemp = data.Data.(CpuTemp).CpuTemp
		break
	case GpuCmd:
		w.GpuInfo = data.Data.([]GpuInfo)
		break
	case UpTimeCmd:
		w.CpuLoad = data.Data.(CpuLoad).CpuLoad
		break
	case GpuEnable:
		w.Gpu = data.Data.(int)

	case LotusMinerJobs:
		w.Jobs = data.Data.([]Task)
	default:

	}
}

func NewStore() *Store {
	return &Store{
		WorkerInfoMap: make(map[string]*WorkerInfo01),
		sign:          make(chan interface{}),
		minerInfo:     make(map[string]interface{}),
	}
}

type Store struct {
	WorkerInfoMap map[string]*WorkerInfo01 // hostName
	minerInfo     map[string]interface{}
	sign          chan interface{}
	wl            sync.RWMutex
	ml            sync.RWMutex
}

func (s *Store) Update(cmdData CmdData) {
	hostName := cmdData.HostName
	s.wl.Lock()
	defer s.wl.Unlock()
	workerInfo01, ok := s.WorkerInfoMap[hostName]
	if !ok {
		workerInfo01 = &WorkerInfo01{HostName: hostName}
		s.WorkerInfoMap[hostName] = workerInfo01
	}
	workerInfo01.updateInfo(cmdData)

}

func (s *Store) Get(hostName string) WorkerInfo01 {
	s.wl.Lock()
	defer s.wl.Unlock()
	workerInfo01, ok := s.WorkerInfoMap[hostName]
	if ok {
		return *workerInfo01
	}
	return WorkerInfo01{}
}
