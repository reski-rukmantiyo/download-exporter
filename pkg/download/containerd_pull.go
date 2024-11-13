package download

import (
	"context"
	"log"
	"time"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/namespaces"
)

func PullContainerdImage(ctx context.Context, imageToPull ImageDownload, location string) {
	startTime := time.Now()
	imageName := imageToPull.Image
	imageLabel := imageToPull.Label

	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		dockerPullFailure.WithLabelValues(imageName, err.Error(), imageLabel, location).Inc()
		log.Println(err)
		return
	}
	defer client.Close()

	ctx = namespaces.WithNamespace(ctx, "default")

	img, err := client.Pull(ctx, imageName)
	if err != nil {
		dockerPullFailure.WithLabelValues(imageName, err.Error(), imageLabel, location).Inc()
		log.Println(err)
		return
	}

	elapsedTime := time.Since(startTime).Seconds()
	dockerPullDuration.WithLabelValues(imageName, imageLabel, location).Observe(elapsedTime)
	dockerPullSuccess.WithLabelValues(imageName, imageLabel, location).Inc()
	log.Printf("Pulled %s in %.2f seconds\n", imageName, elapsedTime)

	// Delete the pulled image
	if err := deleteContainerdImage(ctx, client, img.Name()); err != nil {
		dockerImageDeletionFailure.WithLabelValues(imageName, err.Error(), imageLabel, location).Inc()
		log.Println(err)
	} else {
		dockerImageDeletionSuccess.WithLabelValues(imageName, imageLabel, location).Inc()
		log.Printf("Deleted image %s\n", img.Name())
	}
}

func deleteContainerdImage(ctx context.Context, client *containerd.Client, imageFile string) error {
	imageService := client.ImageService()

	err := imageService.Delete(ctx, imageFile)
	if err != nil && !errdefs.IsNotFound(err) {
		return err
	}

	return nil
}
