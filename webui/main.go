package main

import (
	"flag"
	"github.com/moetang-arch/cicd/webui/config"
	"github.com/moetang-arch/cicd/webui/controller"

	"github.com/gin-gonic/gin"
)

var addr string

func init() {
	flag.StringVar(&addr, "docker.addr", "tcp://192.168.31.201:2001", "-docker.addr=tcp://192.168.31.201:2001")
}

func main() {
	engine := gin.New()

	c := new(config.Config)

	c.DockerConfig.Addr = addr

	controller.InitTemplates(engine, c)

	controller.InitController(engine, c)

	engine.Run("0.0.0.0:3001")
}
