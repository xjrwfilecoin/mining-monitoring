package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"time"
)

var addr = flag.String("addr", "127.0.0.1:9090", "http service address")

func main() {

	log.SetFlags(0)
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/api/v1.0/socket.io"}
	log.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	done := make(chan struct{})

	defer c.Close()
	go func() {
		for {
			select {
			default:
				_, message, err := c.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					return
				}
				log.Printf("recv: %s", message)
			}
		}
	}()

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			return
		case <-ticker.C:

			time.Sleep(3*time.Second)
			blockInfoCmd := fmt.Sprintf(`{"uri":"blockInfo","body":{"type":"S","chainKey":"SFF01","height":1,"hash":"sdfsdfs"},"msgId":"msgId%d"}`, time.Now().Unix())
			err = c.WriteMessage(websocket.TextMessage, []byte(blockInfoCmd))
			if err != nil {
				fmt.Println(err.Error())
				return
			}



		}
	}

}
