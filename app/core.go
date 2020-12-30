package app

import (
	"encoding/json"
	"io/ioutil"
	"mining-monitoring/log"
	"mining-monitoring/model"
	httpsvr "mining-monitoring/net/http"
	"mining-monitoring/net/socket"
	"mining-monitoring/processmanager"
	"os"
	"os/signal"
	"syscall"
)

func Run(path string) error {
	processmanager.Daemon()
	processmanager.CheckPid("mining-monitoring")
	runtimeConfig, err := ReadCfg(path)
	if err != nil {
		return err
	}
	_, err = log.MyLogicLogger(runtimeConfig.LogPath)
	if err != nil {
		return err
	}

	server, err := socket.NewServer()
	if err != nil {
		return err
	}

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
	defer server.Close()
	go func() {
		err := server.Run()
		if err != nil {
			panic(err)
		}
	}()
	// todo db heartbeat
	//// 初始化mongodb
	//err = db.MongodbInit(runtimeConfig)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	log.Logger.Fatal("mongodb start error " + err.Error())
	//	panic(err)
	//}
	httpsvr.ListenAndServe(runtimeConfig,server)
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
