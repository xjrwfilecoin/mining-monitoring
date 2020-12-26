package socket

import (
	"bufio"
	"fmt"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"github.com/zhouhui8915/go-socket.io-client"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

func TestSocketClient(t *testing.T) {
	dialer := engineio.Dialer{
		Transports: []transport.Transport{polling.Default, websocket.Default},
	}
	conn, err := dialer.Dial("http://localhost:8899/socket.io/", nil)
	if err != nil {
		log.Fatalln("dial error:", err)
	}
	defer conn.Close()
	fmt.Println(conn.ID(), conn.LocalAddr(), "->", conn.RemoteAddr(), "with", conn.RemoteHeader())

	go func() {
		defer conn.Close()

		for {
			ft, r, err := conn.NextReader()
			if err != nil {
				log.Println("next reader error:", err)
				return
			}
			b, err := ioutil.ReadAll(r)
			if err != nil {
				r.Close()
				log.Println("read all error:", err)
				return
			}
			if err := r.Close(); err != nil {
				log.Println("read close:", err)
				return
			}
			fmt.Println("read:", ft, b)
		}
	}()

	for {
		fmt.Println("write text hello")
		w, err := conn.NextWriter(engineio.TEXT)
		if err != nil {
			log.Println("next writer error:", err)
			return
		}
		if _, err := w.Write([]byte("hello")); err != nil {
			w.Close()
			log.Println("write error:", err)
			return
		}
		if err := w.Close(); err != nil {
			log.Println("write close error:", err)
			return
		}
		fmt.Println("write binary 1234")
		w, err = conn.NextWriter(engineio.BINARY)
		if err != nil {
			log.Println("next writer error:", err)
			return
		}
		if _, err := w.Write([]byte{1, 2, 3, 4}); err != nil {
			w.Close()
			log.Println("write error:", err)
			return
		}
		if err := w.Close(); err != nil {
			log.Println("write close error:", err)
			return
		}
		time.Sleep(time.Second * 5)
	}
}



func TestSocketIOClient(t *testing.T){
	opts := &socketio_client.Options{
		Transport: "websocket",
		Query:     make(map[string]string),
	}
	opts.Query["user"] = "user"
	opts.Query["pwd"] = "pass"
	uri := "http://localhost:8899/socket.io/"

	client, err := socketio_client.NewClient(uri, opts)
	if err != nil {
		log.Printf("NewClient error:%v\n", err)
		return
	}

	client.On("error", func() {
		log.Printf("on error\n")
	})
	client.On("connection", func() {
		log.Printf("on connect\n")
	})
	client.On("message", func(msg string) {
		log.Printf("on message:%v\n", msg)
	})
	client.On("disconnection", func() {
		log.Printf("on disconnect\n")
	})

	reader := bufio.NewReader(os.Stdin)
	for {
		data, _, _ := reader.ReadLine()
		command := string(data)
		client.Emit("message", "test")
		log.Printf("send message:%v\n", command)
	}
}
