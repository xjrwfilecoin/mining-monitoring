package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	cache "mining-monitoring/cache"
	"mining-monitoring/config"
	"mining-monitoring/log"
	"mining-monitoring/model"
	httpsvr "mining-monitoring/net/http"
	"mining-monitoring/net/socket"
	"mining-monitoring/service"
	"mining-monitoring/shell"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var ShellManager *shell.Manager

func Run(cfgPath string) error {

	// todo just test
	go func() {
		http.ListenAndServe(":6060", nil)
	}()

	runtimeCfg, err := ReadCfg(cfgPath)
	if err != nil {
		return err
	}

	//if runtimeCfg.CpuNum != 0 {
	//	runtime.GOMAXPROCS(int(runtimeCfg.CpuNum))
	//}

	_, err = log.MyLogicLogger(runtimeCfg.LogPath, runtimeCfg.LogLevel)
	if err != nil {
		return err
	}
	ShellManager, err = shell.NewManager()
	if err != nil {
		return fmt.Errorf("init shell shellManager %v \n", err)
	}
	defer ShellManager.Close()
	sign := make(chan shell.CmdData, 100)
	manager := cache.NewManager()
	defer manager.Close()
	for i := 0; i < 100; i++ {
		go manager.Rec(sign)
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
		option := socket.Option{ConnMaxNum: runtimeCfg.ConnMaxNum}
		err := socket.SServer.Run(option)
		if err != nil {
			panic(err)
		}
	}()

	httpsvr.ListenAndServe(runtimeCfg, socket.SServer)
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
