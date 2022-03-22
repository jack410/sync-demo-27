package ws

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wshandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	//升级请求
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	//升级成功后初始化一个client信息
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	//将用户注册到事件中心
	//新用户往regster管道里去放，然后hub.go那边去注册
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
func HttpController(c *gin.Context, hub *Hub) {
	wshandler(hub, c.Writer, c.Request)
}
