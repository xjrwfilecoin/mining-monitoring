package ws

import (
	"mining-monitoring/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

func WebSocketConnHandler(c *gin.Context) {
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if error != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	client := &Client{ID: utils.GetUUID(), Socket: conn, Send: make(chan []byte), ClientIp: c.ClientIP(), Dispatch: NewDisPatch()}
	Manager.Register <- client
	go client.Read()
	go client.Write()
}
