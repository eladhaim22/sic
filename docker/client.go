package docker

import(
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

func getClient() *client.Client {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}

	return cli
}

func ListImages() []types.ImageSummary {
	images, err := getClient().ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		log.Error("Error listing images")
		panic(err)
	}
	return images
}

func ListActiveContainersImageIds() []string {
	containers,err := getClient().ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		log.Error("Error listing active containers", err)
		panic(err)
	}

	var activeContainers []string
	for _, container := range containers {
		if container.State == "running" {
			activeContainers = append(activeContainers, container.ImageID)
		}
	}
	return activeContainers
}

func RemoveImage(imageId string) bool {
	_, err := getClient().ImageRemove(context.Background(), imageId, types.ImageRemoveOptions{Force: true} )
	if err != nil {
		return false
	}
	return true
}