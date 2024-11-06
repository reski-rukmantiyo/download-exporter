package download

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	dockerPullSuccess = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "docker_pull_success",
			Help: "Number of successful Docker image pulls.",
		},
		[]string{"image"},
	)

	dockerPullFailure = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "docker_pull_failure",
			Help: "Number of failed Docker image pulls.",
		},
		[]string{"image", "reason"},
	)

	dockerPullDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "docker_pull_duration_seconds",
			Help:    "Histogram of Docker image pull durations in seconds.",
			Buckets: []float64{1, 5, 10, 30, 60, 120}, // Adjust buckets as needed
		},
		[]string{"image"},
	)

	dockerPullSpeed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "docker_pull_speed_bytes_per_second",
			Help: "Speed of Docker image pull in bytes per second.",
		},
		[]string{"image"},
	)

	dockerImageDeletionSuccess = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "docker_image_deletion_success",
			Help: "Number of successful Docker image deletions.",
		},
		[]string{"image"},
	)

	dockerImageDeletionFailure = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "docker_image_deletion_failure",
			Help: "Number of failed Docker image deletions.",
		},
		[]string{"image", "reason"},
	)
	// imageToPull = "YOUR_DOCKER_IMAGE_HERE" // Update this (e.g., "nginx:latest")
)

func init() {
	prometheus.MustRegister(
		dockerPullSuccess,
		dockerPullFailure,
		dockerPullDuration,
		dockerPullSpeed,
		dockerImageDeletionSuccess,
		dockerImageDeletionFailure,
	)
}

func PullDockerImage(ctx context.Context, imageToPull string) {
	startTime := time.Now()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		dockerPullFailure.WithLabelValues(imageToPull, err.Error()).Inc()
		log.Println(err)
		return
	}

	options := image.PullOptions{}
	rc, err := cli.ImagePull(ctx, imageToPull, options)
	if err != nil {
		dockerPullFailure.WithLabelValues(imageToPull, err.Error()).Inc()
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
				dockerPullFailure.WithLabelValues(imageToPull, err.Error()).Inc()
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
				dockerPullSpeed.WithLabelValues(imageToPull).Set(currentSpeed)
			}
		}
		previousTimestamp = currentTime
		previousSize = int64(n)
	}

	elapsedTime := time.Since(startTime).Seconds()
	dockerPullDuration.WithLabelValues(imageToPull).Observe(elapsedTime)
	dockerPullSuccess.WithLabelValues(imageToPull).Inc()
	log.Printf("Pulled %s in %.2f seconds\n", imageToPull, elapsedTime)

	// Delete the pulled image
	if err := deleteDockerImage(cli, imageToPull); err != nil {
		dockerImageDeletionFailure.WithLabelValues(imageToPull, err.Error()).Inc()
		log.Println(err)
	} else {
		dockerImageDeletionSuccess.WithLabelValues(imageToPull).Inc()
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
