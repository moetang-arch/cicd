package docker_test

import (
	"testing"

	"github.com/moetang-arch/cicd/webui/docker"
)

const (
	skipTest = true

	DOCKER_ENDPOINT = "tcp://192.168.31.201:2001"
)

func TestEndpoint_ListAllImages(t *testing.T) {
	if skipTest {
		return
	}

	ep, err := docker.NewDockerEndpoint(DOCKER_ENDPOINT)
	if err != nil {
		t.Fatal(err)
	}
	defer ep.Close()

	t.Log(ep.ListAllImages())

}

func TestEndpoint_ListAllContainers(t *testing.T) {
	if skipTest {
		return
	}

	ep, err := docker.NewDockerEndpoint(DOCKER_ENDPOINT)
	if err != nil {
		t.Fatal(err)
	}
	defer ep.Close()

	ep.ListAllContainers()
}
