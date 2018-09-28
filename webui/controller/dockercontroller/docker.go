package dockercontroller

import (
	"html/template"
	"net/http"

	"github.com/moetang-arch/cicd/webui/config"
	"github.com/moetang-arch/cicd/webui/docker"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

var (
	dockerTemplateRender render.HTMLRender
	endpoint             *docker.Endpoint
)

func InitDockerTemplate(engine *gin.Engine, config *config.Config) {
	left := "{{"
	right := "}}"
	templ := template.Must(template.New("").Delims(left, right).Funcs(engine.FuncMap).ParseGlob("templates/docker/*"))

	dockerTemplateRender = render.HTMLProduction{Template: templ.Funcs(engine.FuncMap)}
}

func InitDockerController(engine *gin.Engine, config *config.Config) {
	ep, err := docker.NewDockerEndpoint(config.DockerConfig.Addr)
	if err != nil {
		panic(err)
	}
	endpoint = ep

	engine.GET("/docker", index)
	engine.GET("/docker/", index)
	engine.GET("/docker/delete_image/:name", deleteImage)
	engine.GET("/docker/create_temp_container/:name", createTempContainer)
	engine.GET("/docker/delete_container/:name", deleteContainer)
	engine.GET("/docker/delete_containers_by_image/:name", deleteContainersByImage)
}

type IndexResult struct {
	Images     []docker.Image
	Containers []docker.Container
}

func index(c *gin.Context) {
	ir := new(IndexResult)

	images, err := endpoint.ListAllImages()
	if err != nil {
		instance := dockerTemplateRender.Instance("error.html", err.Error())
		c.Render(http.StatusInternalServerError, instance)
		return
	}
	ir.Images = images

	containers, err := endpoint.ListAllContainers()
	if err != nil {
		instance := dockerTemplateRender.Instance("error.html", err.Error())
		c.Render(http.StatusInternalServerError, instance)
		return
	}
	ir.Containers = containers

	instance := dockerTemplateRender.Instance("index.html", ir)
	c.Render(http.StatusOK, instance)
}

func deleteImage(c *gin.Context) {
	img := c.Param("name")

	if len(img) == 0 {
		instance := dockerTemplateRender.Instance("error.html", "image name is empty")
		c.Render(http.StatusBadRequest, instance)
	}

	err := endpoint.RemoveImage(img)
	if err != nil {
		instance := dockerTemplateRender.Instance("error.html", err.Error())
		c.Render(http.StatusInternalServerError, instance)
		return
	}

	c.Redirect(http.StatusFound, "/docker")
}

func createTempContainer(c *gin.Context) {
	img := c.Param("name")

	if len(img) == 0 {
		instance := dockerTemplateRender.Instance("error.html", "image name is empty")
		c.Render(http.StatusBadRequest, instance)
	}

	err := endpoint.CreateAnonymousContainer(img)
	if err != nil {
		instance := dockerTemplateRender.Instance("error.html", err.Error())
		c.Render(http.StatusInternalServerError, instance)
		return
	}

	c.Redirect(http.StatusFound, "/docker")
}

func deleteContainer(c *gin.Context) {
	container := c.Param("name")

	if len(container) == 0 {
		instance := dockerTemplateRender.Instance("error.html", "container id is empty")
		c.Render(http.StatusBadRequest, instance)
	}

	err := endpoint.DeleteContainer(container)
	if err != nil {
		instance := dockerTemplateRender.Instance("error.html", err.Error())
		c.Render(http.StatusInternalServerError, instance)
		return
	}

	c.Redirect(http.StatusFound, "/docker")
}

func deleteContainersByImage(c *gin.Context) {
	img := c.Param("name")

	if len(img) == 0 {
		instance := dockerTemplateRender.Instance("error.html", "image name is empty")
		c.Render(http.StatusBadRequest, instance)
	}

	err := endpoint.DeleteContainersByImage(img)
	if err != nil {
		instance := dockerTemplateRender.Instance("error.html", err.Error())
		c.Render(http.StatusInternalServerError, instance)
		return
	}

	c.Redirect(http.StatusFound, "/docker")
}
