package tasks

import (
	"fmt"
	"github.com/eladhaim22/sic/docker"
	"github.com/eladhaim22/sic/utils"
	"github.com/eladhaim22/sic/utils/env"
	log "github.com/sirupsen/logrus"
	"strings"
)

func CleaningTask(){
	dockerSwarmClient := docker.NewDockerClient()
	log.Info("Task started")
	log.Info("Searching for unused images")

	images := dockerSwarmClient.ListImages()
	activeContainersImageIds := dockerSwarmClient.ListActiveContainersImageIds()
	var unusedImages []docker.DockerImage

	for _, image := range images {
		if ok := !utils.FindInSlice(activeContainersImageIds, image.Summery.ID) && !excludeImage(image); ok {
			unusedImages = append(unusedImages, image)
		}
	}

	if len(unusedImages) > 0 {
		log.Info(fmt.Sprintf("%d unused images can be delete", len(unusedImages)))
		for _, imageToDelete := range unusedImages {
			if deleted := dockerSwarmClient.RemoveImage(imageToDelete); deleted {
				log.Info(fmt.Sprintf("The image: %s deleted",imageToDelete.Summery.RepoDigests[0]))
			}
		}
	} else {
		log.Info("All images are used.")
	}
	log.Info("Task completed")
}

func excludeImage(image docker.DockerImage) bool {
	if len(env.ExcludeImages) != 0 {
		imagesToExclude := strings.Split(env.ExcludeImages, ",")
		if ok := utils.FindInSlice(imagesToExclude, image.Summery.RepoDigests[0]); ok {
			return true
		}
		return false
	}
	return false
}