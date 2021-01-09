package app

import (
	"encoding/json"
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"io/ioutil"
	"mining-monitoring/log"
	"mining-monitoring/model"
	httpsvr "mining-monitoring/net/http"
	"mining-monitoring/net/socket"
	"mining-monitoring/processmanager"
	"mining-monitoring/service"
	"mining-monitoring/shellParsing"
	"os"
	"os/signal"
	"syscall"
)

var ShellManager *shellParsing.Manager

func Run(config, workerHost string) error {
	processmanager.Daemon()
	processmanager.CheckPid("mining-monitoring")
	runtimeConfig, err := ReadCfg(config)
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
	// todo
	// 注册socketIo路由
	minerInfo := service.NewMinerInfoService(ShellManager)
	socket.SServer.RegisterRouter(DefaultNamespace, MinerInfo, minerInfo.MinerInfo)
	socket.SServer.RegisterRouter(DefaultNamespace, SubMinerInfo, func(s socketio.Conn, msg string) {
		socket.SServer.JoinRoom(DefaultNamespace, DefaultRoom, s)
		log.Debug(s.ID(), s.LocalAddr(),"join room ", DefaultRoom)
		s.Emit(SubMinerInfo, "info")
	})


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
				fmt.Println()
				bytes, err := json.Marshal(result)
				if err!=nil{
					continue
				}
				log.Info("send subMinerInfo:  ", string(bytes))
				fmt.Println()
				socket.BroadCaseMsg(result)
			default:

			}
		}
	}()

	// todo
	//go ShellManager.Run(minerObjSign)

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
	return c, nil
}
