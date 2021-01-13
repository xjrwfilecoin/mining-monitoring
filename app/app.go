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
	"mining-monitoring/processmanager"
	"mining-monitoring/service"
	"mining-monitoring/shellParsing"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

var ShellManager *shellParsing.Manager

func Run(cfgPath, workerHost string) error {
	processmanager.Daemon()
	processmanager.CheckPid("mining-monitoring")
	runtimeConfig, err := ReadCfg(cfgPath)
	if err != nil {
		return err
	}
	_, err = log.MyLogicLogger(runtimeConfig.LogPath)
	if err != nil {
		return err
	}

	ShellManager, err = shellParsing.NewManager(workerHost)
	if err != nil {
		return fmt.Errorf("init shell shellManager %v \n", err)
	}

	// 注册socketIo路由
	minerInfo := service.NewMinerInfoService(ShellManager, socket.SServer)
	socket.SServer.RegisterRouterV1(config.DefaultNamespace, config.MinerInfo, minerInfo.MinerInfo)
	socket.SServer.RegisterRouterV1(config.DefaultNamespace, config.SubMinerInfo, minerInfo.SuMinerInfo)

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, processmanager.SIGUSR1, processmanager.SIGUSR2)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				log.Logger.Printf("recv sig %+v and exit process\n", s)
				processmanager.DefExitFunc()
			case processmanager.SIGUSR1:
				log.Logger.Printf("recv sig %+v and exit process\n", s)
				processmanager.DefExitFunc()
			case processmanager.SIGUSR2:
				log.Logger.Printf("recv sig %+v and exit process\n", s)
				processmanager.DefExitFunc()
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

	minerObjSign := make(chan interface{}, 1)
	go func() {
		for {
			select {
			case result := <-minerObjSign:
				output := ParseMinerInfo(result)
				socket.BroadCaseMsg(config.DefaultNamespace, config.DefaultRoom, config.SubMinerInfo, output)
			default:

			}
		}
	}()

	// todo
	go ShellManager.Run(minerObjSign)

	// todo db heartbeat
	//// 初始化mongodb
	//err = db.MongodbInit(runtimeConfig)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	log.Logger.Fatal("mongodb start error " + err.Error())
	//	panic(err)
	//}
	httpsvr.ListenAndServe(runtimeConfig, socket.SServer)
	return nil
}

func ParseMinerInfo(input interface{}) interface{} {
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
	delete(tempMap, "jobs")
	delete(tempMap, "hardwareInfo")
	tempMap["workerInfo"] = workerInfo
	return tempMap
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
