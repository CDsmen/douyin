package main

import (
	"github.com/CDsmen/douyin/dal"
	"github.com/CDsmen/douyin/service"
	"github.com/gin-gonic/gin"
)

func main() {
	go service.RunMessageServer()

	r := gin.Default()

	dal.InitDB()

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
