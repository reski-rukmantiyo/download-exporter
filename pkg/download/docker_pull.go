package download

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

func PullDockerImage(ctx context.Context, imageToPull ImageDownload, location string) {
	startTime := time.Now()
	imageName := imageToPull.Image
	imageLabel := imageToPull.Label

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		dockerPullFailure.WithLabelValues(imageName, err.Error(), imageLabel, location).Inc()
		log.Println(err)
		return
	}

	options := image.PullOptions{}
	rc, err := cli.ImagePull(ctx, imageToPull.Image, options)
	if err != nil {
		dockerPullFailure.WithLabelValues(imageName, err.Error(), imageLabel, location).Inc()
		log.Println(err)
		return
	}
	defer rc.Close()

	var (
		previousTimestamp time.Time
		previousSize      int64
		currentSpeed      float64
	)
	buf := make([]byte, 1024*1024) // 1MB buffer
	for {
		n, err := rc.Read(buf)
		if err != nil {
			if err != io.EOF {
				dockerPullFailure.WithLabelValues(imageName, err.Error(), imageLabel, location).Inc()
				log.Println(err)
			}
			break
		}

		currentTime := time.Now()
		if !previousTimestamp.IsZero() {
			sizeDiff := int64(n) + previousSize
			timeDiff := currentTime.Sub(previousTimestamp).Seconds()
			if timeDiff > 0 {
				currentSpeed = float64(sizeDiff) / timeDiff
				dockerPullSpeed.WithLabelValues(imageName, imageLabel, location).Set(currentSpeed)
			}
		}
		previousTimestamp = currentTime
		previousSize = int64(n)
	}

	elapsedTime := time.Since(startTime).Seconds()
	dockerPullDuration.WithLabelValues(imageName, imageLabel, location).Observe(elapsedTime)
	dockerPullSuccess.WithLabelValues(imageName, imageLabel, location).Inc()
	log.Printf("Pulled %s in %.2f seconds\n", imageName, elapsedTime)

	// Delete the pulled image
	if err := deleteDockerImage(cli, imageName); err != nil {
		dockerImageDeletionFailure.WithLabelValues(imageName, err.Error(), imageLabel, location).Inc()
		log.Println(err)
	} else {
		dockerImageDeletionSuccess.WithLabelValues(imageName, imageLabel, location).Inc()
		log.Printf("Deleted image %s\n", imageToPull)
	}
}

func deleteDockerImage(cli *client.Client, imageFile string) error {
	options := image.RemoveOptions{
		Force:         true,
		PruneChildren: true,
	}
	_, err := cli.ImageRemove(context.Background(), imageFile, options)
	return err
}
