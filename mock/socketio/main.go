package main

import (
	"flag"
	"fmt"
	"github.com/googollee/go-socket.io/engineio/base"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"time"
)

var addr = flag.String("addr", "0.0.0.0:8899", "http service address")

type packet struct {
	Type         string
	NSP          string
	Id           int
	Data         interface{}
	attachNumber int
}

func main() {
	log.SetFlags(0)
	u := url.URL{Scheme: "ws", Host: *addr, Path: `/socket.io/`}
	query := u.Query()
	query.Set("transport", "websocket")
	query.Set("t", base.Timestamp())
	u.RawQuery = query.Encode()
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
			src := `2/,456["message:message",123]`
			fmt.Printf(fmt.Sprintf("%v \n", src))
			err = c.WriteMessage(websocket.TextMessage, []byte(src))
			if err != nil {
				log.Println("write:", err)
				return
			}

		}
	}

}
