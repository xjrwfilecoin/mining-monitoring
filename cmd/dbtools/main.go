package main

import "mining-monitoring/db"

func main() {
	manager := db.NewTableManager()
	err := manager.ReadCfg("./tablecfg.json")
	if err != nil {
		panic(err)
	}
	err = manager.Run()
	if err != nil {
		panic(err)
	}

}
