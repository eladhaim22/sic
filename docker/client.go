package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	utils "github.com/eladhaim22/sic/utils"
	"github.com/eladhaim22/sic/utils/env"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"net"
	"strings"
)

type DockerClient interface {
	ListImages() []DockerImage
	ListActiveContainersImageIds() []string
	RemoveImage(image DockerImage) bool
}

type DockerClientImpl struct {
	clients map[string]*client.Client
}

type DockerImage struct {
	Summery  types.ImageSummary
	ClientId string
}

func NewDockerClient() *DockerClientImpl {
	clients := make(map[string]*client.Client)

	if env.SwarmMode {
		if len(env.NodesAgentsIp) == 0 {
			log.Error("NODES_AGENTS is mandatory on swarm mode")
			panic("")
		}
		ips := make([]string,0)
		success, ips := parsePlainIps(env.NodesAgentsIp)
		if !success {
			if len(env.NodesAgentsPort) == 0 {
				log.Error("NODE_AGENTS_PORT is mandatory if NODES_AGENTS_IPS is not explicit ips")
			}
			ips = ipLookUp(env.NodesAgentsIp, env.NodesAgentsPort)
		}

		for _, ip := range ips {
			cli, err := client.NewClientWithOpts(client.WithHost("http://" + ip),client.WithAPIVersionNegotiation())
			if err != nil {
				panic(err)
			}

			clientUuid := uuid.NewV4()

			if err != nil {
				panic(err)
			}

			clients[clientUuid.String()] = cli
		}
	} else {
		cli, err := client.NewClientWithOpts()
		if err != nil {
			panic(err)
		}

		clientUuid := uuid.NewV4()
		clients[clientUuid.String()] = cli
	}

	return &DockerClientImpl{
		clients: clients,
	}
}

func (dockerClients *DockerClientImpl) ListImages() []DockerImage {
	images := make([]DockerImage,0)
	for clientId, client := range dockerClients.clients {
		clientImages, err := client.ImageList(context.Background(), types.ImageListOptions{All: true})
		if err != nil {
			log.Error("Error listing images")
			panic(err)
		}
		images = append(images, addClientIdToImageArray(clientImages, clientId)...)
	}
	return images
}

func (dockerClients *DockerClientImpl) ListActiveContainersImageIds() []string {
	var containers []string
	filterArgs := filters.NewArgs()
	filterArgs.Add("status", "running")
	for _, client := range dockerClients.clients {
		clientContainers, err := client.ContainerList(context.Background(), types.ContainerListOptions{Filters: filterArgs})
		if err != nil {
			log.Error("Error listing containers")
			panic(err)
		}
		containers = append(containers, Map(clientContainers, func(m types.Container) string { return m.ImageID })...)
	}

	return containers
}

func (dockerClients *DockerClientImpl) RemoveImage (image DockerImage) bool {
	_, err := dockerClients.clients[image.ClientId].ImageRemove(context.Background(), image.Summery.ID, types.ImageRemoveOptions{Force: true} )
	if err != nil {
		return false
	}
	return true
}

func Map(vs []types.Container, f func(container types.Container) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func addClientIdToImageArray(images []types.ImageSummary, clientId string) []DockerImage {
	imagesWIthClientId := make([]DockerImage,0)
	for _, image := range images {
		dockerImage:= DockerImage{image, clientId}
		imagesWIthClientId = append(imagesWIthClientId, dockerImage)
	}
	return imagesWIthClientId
}

func parsePlainIps(ipString string) (bool, []string){
	ips := strings.Split(ipString, ",")
	parsedIps := make([]string,0)
	for _, ip := range ips {
		if ok := utils.IsIpv4WithPort(ip); !ok {
			return false, make([]string,0)
		}
		parsedIps = append(parsedIps, ip)
	}
	return true, parsedIps
}

func ipLookUp(hostName string, port string) []string {
	netIps, err := net.LookupIP(hostName)
	if err != nil {
		panic(err)
	}

	ips := make([]string,0)
	for _,ip := range netIps {
		ips = append(ips, fmt.Sprintf("%s:%s", ip, port))
	}

	return ips
}