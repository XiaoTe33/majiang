package router

import (
	"github.com/gin-gonic/gin"
	"majiang/dao"
	"majiang/log"
)

var (
	myLog = log.Log
	db    = dao.DB
)

func InitRouters() {
	r := gin.Default()
	r.Use(Cors())

	r.POST("/user/register", register)
	r.POST("/user/login", login)
	r.GET("/user/token/refresh", refreshToken)

	r.GET("/join/:room", JWT(), joinRoom)

	_ = r.Run()
}
