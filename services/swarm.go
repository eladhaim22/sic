package services

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/eladhaim22/sic/docker"
	"github.com/eladhaim22/sic/utils"
	"github.com/eladhaim22/sic/utils/env"
	log "github.com/sirupsen/logrus"
	"strings"
)

func CleaningTask(){
	log.Info("Task started")
	log.Info("Searching for unused images")

	images := docker.ListImages()
	activeContainersImageIds := docker.ListActiveContainersImageIds()

	var unusedImages []types.ImageSummary

	for _, image := range images {
		if _, ok := utils.FindInSlice(activeContainersImageIds, image.ID); !ok {
			unusedImages = append(unusedImages, image)
		}
	}

	var unusedImagesWithoutExcluded []types.ImageSummary
	if len(env.ExcludeImages) != 0 {
		imagesToExclude := strings.Split(env.ExcludeImages, ",")
		for _, image := range unusedImages {
			if _, ok := utils.FindInSlice(imagesToExclude, image.RepoDigests[0]); ok {
				unusedImagesWithoutExcluded = append(unusedImagesWithoutExcluded, image)
			}
		}
	} else {
		unusedImagesWithoutExcluded = unusedImages
	}

	if len(unusedImagesWithoutExcluded) > 0 {
		log.Info(fmt.Sprintf("%d unused images can be delete", len(unusedImagesWithoutExcluded)))
		for _, imageToDelete := range unusedImagesWithoutExcluded {
			if deleted := docker.RemoveImage(imageToDelete.ID); deleted {
				var imageSha = ""
				imageDigest := imageToDelete.RepoDigests[0]
				imageName := strings.Split(imageDigest,"@")[0]
				if count := len(strings.Split(imageDigest, "@")); count > 1 {
					imageSha = strings.Split(imageDigest, "@")[1]
				}
				log.Info(fmt.Sprintf("The image: %s%s %s deleted",imageName, imageDigest, imageSha))
			}
		}
	} else {
		log.Info("All images are used.")
	}
	log.Info("Task completed")
}