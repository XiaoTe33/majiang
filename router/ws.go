package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"majiang/model"
	"net/http"
)

func joinRoom(c *gin.Context) {
	conn, err := (&websocket.Upgrader{
		ReadBufferSize:  10240,
		WriteBufferSize: 10240,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("enter the chat room failed, err: %v", err)
		return
	}
	id := c.GetInt64("id")
	myLog.Println(id)
	model.NewPlayer(id, c.GetString("username"), c.Param("room"), conn)
	//go model.ShowClients()

}
