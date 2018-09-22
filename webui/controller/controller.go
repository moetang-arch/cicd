package controller

import (
	"github.com/moetang-arch/cicd/webui/controller/dockercontroller"
	"github.com/moetang-arch/cicd/webui/config"

	"github.com/gin-gonic/gin"
)

func InitTemplates(engine *gin.Engine, config *config.Config) {
	dockercontroller.InitDockerTemplate(engine, config)
}

func InitController(engine *gin.Engine, config *config.Config) {
	dockercontroller.InitDockerController(engine, config)
}
