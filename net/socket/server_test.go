package socket

import (
	"encoding/hex"
	"fmt"
	"github.com/googollee/go-socket.io"
	engineio "github.com/zhouhui8915/engine.io-go"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)



func TestServer02(t *testing.T){
	server, err := engineio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			conn, _ := server.Accept()
			go func() {
				defer conn.Close()
				for i := 0; i < 10; i++ {
					t, r, _ := conn.NextReader()
					b, _ := ioutil.ReadAll(r)
					r.Close()
					if t == engineio.MessageText {
						log.Println(t, string(b))
					} else {
						log.Println(t, hex.EncodeToString(b))
					}
					w, _ := conn.NextWriter(t)
					w.Write([]byte("pong"))
					w.Close()
				}
			}()
		}
	}()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:5000...")
	log.Fatal(http.ListenAndServe(":8899", nil))
}



func TestServer01(t *testing.T){

}




func TestServer(t *testing.T){
	server, _ := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		s.Emit("message","ddd")
		return nil
	})

	server.OnEvent("/", "message", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})

	server.OnEvent("/chat", "message", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		fmt.Printf("ddddd")
		return "recv " + msg
	})

	server.OnEvent("/", "message", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":9090", nil))
}