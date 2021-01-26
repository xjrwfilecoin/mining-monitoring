package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mining-monitoring/config"
	"mining-monitoring/log"
	"mining-monitoring/model"
	httpsvr "mining-monitoring/net/http"
	"mining-monitoring/net/socket"
	"mining-monitoring/service"
	"mining-monitoring/shellParsing"
	"mining-monitoring/store"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"
)

var ShellManager *shellParsing.Manager

func Run(cfgPath string) error {
	runtimeConfig, err := ReadCfg(cfgPath)
	if err != nil {
		return err
	}
	_, err = log.MyLogicLogger(runtimeConfig.LogPath)
	if err != nil {
		return err
	}

	ShellManager, err = shellParsing.NewManager()
	if err != nil {
		return fmt.Errorf("init shell shellManager %v \n", err)
	}
	sign := make(chan shellParsing.CmdData, 100)
	manager := store.NewManager()

	go manager.Recv(sign)
	go manager.Send()
	go ShellManager.RunV1(sign)
	// 注册路由
	minerInfo := service.NewMinerInfoService(manager, socket.SServer)
	socket.SServer.RegisterRouterV1(config.DefaultNamespace, config.MinerInfo, minerInfo.MinerInfo)
	socket.SServer.RegisterRouterV1(config.DefaultNamespace, config.SubMinerInfo, minerInfo.SuMinerInfo)

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, config.SIGUSR1, config.SIGUSR2)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				log.Logger.Printf("recv sig %+v and exit process\n", s)
				DefExitFunc()
			case config.SIGUSR1:
				log.Logger.Printf("recv sig %+v and exit process\n", s)
				DefExitFunc()
			case config.SIGUSR2:
				log.Logger.Printf("recv sig %+v and exit process\n", s)
				DefExitFunc()
			default:
				log.Logger.Println("other", s)
			}
		}
	}()
	defer socket.SServer.Close()
	go func() {
		err := socket.SServer.Run()
		if err != nil {
			panic(err)
		}
	}()

	//minerObjSign := make(chan map[string]interface{}, 1)
	//go timerMinerInfo(minerObjSign)
	//go broadCastMessage(ShellManager, minerObjSign)
	//go ShellManager.Run(sign)

	httpsvr.ListenAndServe(runtimeConfig, socket.SServer)
	return nil
}

func timerMinerInfo(minerObjSign chan map[string]interface{}) {
	for {
		select {
		case minerInfo := <-minerObjSign:
			output := ParseMinerInfo(minerInfo)
			socket.BroadCaseMsg(config.DefaultNamespace, config.DefaultRoom, config.SubMinerInfo, output)
		default:

		}
	}
}

func broadCastMessage(shellManger *shellParsing.Manager, sign chan map[string]interface{}) {
	previousMap := make(map[string]interface{})
	result := make(map[string]interface{})
	var err error
	for {
		select {
		case minerInfo := <-sign:
			if len(previousMap) == 0 {
				previousMap = minerInfo
				result = minerInfo
			} else {
				previousMap, err = DeepCopyMap(minerInfo)
				if err != nil {
					log.Error("deepCopyMap: ", err.Error())
					continue
				}

				info := mergeMinerInfo(previousMap)
				shellManger.UpdateCurrentMinerInfo(info)

				result = DiffMap(previousMap, minerInfo)
			}
			bytes, err := json.Marshal(result)
			if err == nil {
				log.Warn("diffMinerInfo: ", string(bytes))

			}
			output := mergeMinerInfo(result)
			socket.BroadCaseMsg(config.DefaultNamespace, config.DefaultRoom, config.SubMinerInfo, output)
		default:

		}
	}
}

func mergeMinerInfo(input map[string]interface{}) map[string]interface{} {
	jobs := make(map[string]interface{})
	hardwareInfo := make(map[string]interface{})
	if reflect.TypeOf(input).Kind() != reflect.Map {
		return nil
	}
	tJobs, ok := input["jobs"]
	if ok {
		jobs = tJobs.(map[string]interface{})
	}
	tHardwareInfo, ok := input["hardwareInfo"]
	if ok {
		hardwareInfo = tHardwareInfo.(map[string]interface{})
	}

	workerInfo := MapParse(jobs, hardwareInfo)
	input["workerInfo"] = workerInfo
	delete(input, "jobs")
	delete(input, "hardwareInfo")
	return input
}

func ParseMinerInfo(input interface{}) map[string]interface{} {
	param := make(map[string]interface{})
	if reflect.TypeOf(input).Kind() != reflect.Map {
		return param
	}
	tempMap := input.(map[string]interface{})
	jobs := tempMap["jobs"]
	hardwareInfo := tempMap["hardwareInfo"]
	tJobs := jobs.(map[string]interface{})
	tHardware := hardwareInfo.(map[string]interface{})
	workerInfo := MapParse(tJobs, tHardware)
	tempMap["workerInfo"] = workerInfo
	delete(tempMap, "jobs")
	delete(tempMap, "hardwareInfo")
	return tempMap
}

// DefExitFunc 缺省信号退出函数
func DefExitFunc() {
	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%v Default Sig Exit ... \n", now)
	os.Exit(0)
}

func ReadCfg(path string) (*model.RuntimeConfig, error) {
	c := &model.RuntimeConfig{}
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(d, c)
	if err != nil {
		return nil, err
	}
	fmt.Printf("config info: %v \n", string(d))
	return c, nil
}
