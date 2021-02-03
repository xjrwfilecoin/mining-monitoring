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
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var ShellManager *shellParsing.Manager

func Run(cfgPath string) error {

	// just test
	go func() {
		http.ListenAndServe(":6060", nil)
	}()

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
	defer ShellManager.Close()
	sign := make(chan shellParsing.CmdData, 100)
	manager := store.NewManager()
	defer manager.Close()
	for i := 0; i < 100; i++ {
		go manager.Recv(sign)
	}
	go manager.Send()
	go ShellManager.Run(sign)
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

	httpsvr.ListenAndServe(runtimeConfig, socket.SServer)
	return nil
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
