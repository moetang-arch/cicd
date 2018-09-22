package docker

import (
	"context"
	"time"
	"log"
	"os"
	"errors"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types/container"
)

var (
	logger *log.Logger = log.New(os.Stderr, "", log.LstdFlags)
)

func timeout(timeInSecond int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(timeInSecond)*time.Second)
}

type Endpoint struct {
	client *client.Client
}

func (this *Endpoint) Close() error {
	return this.client.Close()
}

func NewDockerEndpoint(addr string) (*Endpoint, error) {
	c, err := client.NewClientWithOpts(
		client.WithHost(addr),
		client.WithVersion("1.37"),
	)

	if err != nil {
		return nil, err
	}

	return &Endpoint{
		client: c,
	}, nil
}

func (this *Endpoint) RemoveImage(img string) error {
	tenSec, _ := timeout(10)
	option := types.ImageRemoveOptions{}
	resp, err := this.client.ImageRemove(tenSec, img, option)

	if err != nil {
		return err
	}

	logger.Println("delete result:", resp)

	return nil
}

func (this *Endpoint) PullImage() {
	//FIXME need to implement
	panic("need to implement")
}

type Image struct {
	UseName string
	Tags    []string
}

func (this *Endpoint) ListAllImages() ([]Image, error) {
	tenSec, _ := timeout(10)
	option := types.ImageListOptions{
		All: false,
	}
	summaries, err := this.client.ImageList(tenSec, option)
	if err != nil {
		return nil, err
	}

	var result []Image
	for _, v := range summaries {
		result = append(result, Image{
			UseName: v.RepoTags[0], // use first tag
			Tags:    v.RepoTags,
		})
	}

	return result, nil
}

type Container struct {
	Id      string
	Image   string
	Command string
	State   string
}

func (this *Endpoint) CreateAnonymousContainer(img string) error {
	tenSec, _ := timeout(10)
	config := &container.Config{
		Image: img,
	}
	r, err := this.client.ContainerCreate(tenSec, config, nil, nil, "")

	if err != nil {
		return err
	}

	logger.Println("create image:", r)
	return nil
}

func (this *Endpoint) ListAllContainers() ([]Container, error) {
	tenSec, _ := timeout(10)
	option := types.ContainerListOptions{
		All: true,
	}
	containers, err := this.client.ContainerList(tenSec, option)

	if err != nil {
		return nil, err
	}

	var result []Container
	for _, v := range containers {
		result = append(result, Container{
			Id:      v.ID,
			Image:   v.Image,
			Command: v.Command,
			State:   v.State,
		})
	}

	return result, nil
}

func (this *Endpoint) DeleteContainer(id string) error {
	tenSec, _ := timeout(10)
	option := types.ContainerRemoveOptions{}
	return this.client.ContainerRemove(tenSec, id, option)
}

func (this *Endpoint) DeleteContainersByImage(img string) error {
	containers, err := this.ListAllContainers()
	if err != nil {
		return err
	}

	var rmIds []string
	for _, v := range containers {
		if v.Image == img {
			rmIds = append(rmIds, v.Id)
		}
	}

	if len(rmIds) == 0 {
		return errors.New("no containers found with the specific image")
	}

	logger.Println("removing containers list:", rmIds)

	for _, v := range rmIds {
		err := this.DeleteContainer(v)
		if err != nil {
			return err
		}
	}

	return nil
}
