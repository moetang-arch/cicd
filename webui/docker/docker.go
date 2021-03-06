package docker

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
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
	ID      string
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
		var useName string
		if len(v.RepoTags) > 0 {
			useName = v.RepoTags[0]
		}
		result = append(result, Image{
			UseName: useName, // use first tag
			Tags:    v.RepoTags,
			ID:      v.ID,
		})
	}

	return result, nil
}

type Container struct {
	Id      string
	Image   string
	Command string
	State   string
	Name    string
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
			Name:    v.Names[0], // using first name
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
