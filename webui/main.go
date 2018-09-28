package main

import (
	"github.com/moetang-arch/cicd/webui/config"
	"github.com/moetang-arch/cicd/webui/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.New()

	c := new(config.Config)

	c.DockerConfig.Addr = "tcp://192.168.31.201:2001"

	controller.InitTemplates(engine, c)

	controller.InitController(engine, c)

	engine.Run("0.0.0.0:3001")
}
