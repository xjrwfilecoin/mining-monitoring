package main

import (
	"bufio"
	"fmt"
	socketio_client "github.com/zhouhui8915/go-socket.io-client"
	"log"
	"os"
)



func main() {
	opts := &socketio_client.Options{
		//Transport:"polling",
		Transport:"websocket",
		Query:     make(map[string]string),
	}
	opts.Query["user"] = "user"
	opts.Query["pwd"] = "pass"
	uri := "http://127.0.0.1:9090"

	client, err := socketio_client.NewClient(uri, opts)
	if err != nil {
		log.Printf("NewClient error:%v\n", err)
		return
	}

	err = client.On("error", func() {
		log.Printf("on error\n")
	})
	err = client.On("connection", func() {
		log.Printf("on connect\n")
	})

	err = client.On("minerInfo", func(msg string) {
		log.Printf("on minerInfo:%v\n", msg)
	})

	err = client.On("subMinerInfo", func(msg string) {
		log.Printf("on subMinerInfo:%v\n", msg)
	})

	err = client.On("disconnection", func() {
		log.Printf("on disconnect\n")
	})
	if err!=nil{
		fmt.Println(err.Error())
	}
	err= client.Emit("minerInfo","getMinerInfo" )
	if err!=nil{
		fmt.Println(err.Error())
	}
	err=client.Emit("subMinerInfo", "get subMinerInfo result")
	if err!=nil{
		fmt.Println(err.Error())
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		data, _, _ := reader.ReadLine()
		command := string(data)
		err:= client.Emit("minerInfo", command)
		if err!=nil{
			fmt.Println(err.Error())
		}
		log.Printf("send message:%v\n", command)
	}
}