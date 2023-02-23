package main

import (
	"github.com/CDsmen/douyin/controller"
	"github.com/CDsmen/douyin/dal"
	"github.com/CDsmen/douyin/service"
	"github.com/gin-gonic/gin"
)

func main() {
	controller.SeverIp = "127.0.0.1"

	go service.RunMessageServer()

	r := gin.Default()

	dal.InitDB()

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
